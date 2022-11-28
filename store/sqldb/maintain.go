/**
 * Tencent is pleased to support the open source community by making Polaris available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package sqldb

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/polarismesh/polaris/common/utils"
	"github.com/polarismesh/polaris/store"
)

const (
	DefaultElectKey = "polaris-server"
	TickTime        = 2
	LeaseTime       = 10
)

// maintainStore implement MaintainStore interface
type maintainStore struct {
	master *BaseDB
	le     *leaderElectionStateMachine
}

// LeaderElectionStore store inteface
type LeaderElectionStore interface {
	// GetVersion get current version
	GetVersion(key string) (int64, error)
	// CompareAndSwapVersion cas version
	CompareAndSwapVersion(key string, curVersion int64, newVersion int64, leader string) (bool, error)
	// CheckMtimeExpired check mtime expired
	CheckMtimeExpired(key string, leaseTime int32) (bool, error)
}

// leaderElectionStore
type leaderElectionStore struct {
	master *BaseDB
}

// GetVersion
func (l *leaderElectionStore) GetVersion(key string) (int64, error) {
	log.Debugf("[Store][database] get version (%s)", key)
	mainStr := "select version from leader_election where elect_key = ?"

	var count int64
	err := l.master.DB.QueryRow(mainStr, key).Scan(&count)
	if err != nil {
		log.Errorf("[Store][database] get version (%s), err: %s", key, err.Error())
	}
	return count, store.Error(err)
}

// CompareAndSwapVersion
func (l *leaderElectionStore) CompareAndSwapVersion(key string, curVersion int64, newVersion int64, leader string) (bool, error) {
	log.Debugf("[Store][database] compare and swap version (%s, %d, %d, %s)", key, curVersion, newVersion, leader)
	mainStr := "update leader_election set leader = ?, version = ? where elect_key = ? and version = ?"
	result, err := l.master.DB.Exec(mainStr, leader, newVersion, key, curVersion)
	if err != nil {
		log.Errorf("[Store][database] compare and swap version, err: %s", err.Error())
		return false, store.Error(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Errorf("[Store][database] compare and swap version, get RowsAffected err: %s", err.Error())
		return false, store.Error(err)
	}
	return (rows > 0), nil
}

// CheckMtimeExpired
func (l *leaderElectionStore) CheckMtimeExpired(key string, leaseTime int32) (bool, error) {
	log.Debugf("[Store][database] check mtime expired (%s, %d)", key, leaseTime)
	mainStr := "select count(1) from leader_election where elect_key = ? and mtime < FROM_UNIXTIME(UNIX_TIMESTAMP(SYSDATE()) - ?)"

	var count int32
	err := l.master.DB.QueryRow(mainStr, key, leaseTime).Scan(&count)
	if err != nil {
		log.Errorf("[Store][database] check mtime expired (%s), err: %s", key, err.Error())
	}
	return (count > 0), store.Error(err)
}

// leaderElectionStateMachine
type leaderElectionStateMachine struct {
	leStore    LeaderElectionStore
	leaderFlag int32
	version    int64
	ctx        context.Context
	cancel     context.CancelFunc
}

// isLeader
func isLeader(flag int32) bool {
	return flag > 0
}

// mainLoop
func (le *leaderElectionStateMachine) mainLoop() {
	log.Infof("[Store][database] leader election started")
	ticker := time.NewTicker(TickTime * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			le.tick()
		case <-le.ctx.Done():
			log.Infof("[Store][database] leader election stopped")
			return
		}
	}
}

// tick
func (le *leaderElectionStateMachine) tick() {
	if le.isLeader() {
		r, err := le.heartbeat()
		if err != nil {
			log.Errorf("[Store][database] leader heartbeat err (%s), change to follower state", err.Error())
			le.changeToFollower()
			return
		}
		if !r {
			le.changeToFollower()
		}
	} else {
		dead, err := le.checkLeaderDead()
		if err != nil {
			log.Errorf("[Store][database] check leader dead err (%s), stay follower state", err.Error())
			return
		}
		if !dead {
			return
		}
		r, err := le.elect()
		if err != nil {
			log.Errorf("[Store][database] elect leader err (%s), stay follower state", err.Error())
			return
		}
		if r {
			le.changeToLeader()
		}
	}
}

// changeToLeader
func (le *leaderElectionStateMachine) changeToLeader() {
	log.Infof("[Store][database] change from follower to leader")
	atomic.StoreInt32(&le.leaderFlag, 1)
}

// changeToFollower
func (le *leaderElectionStateMachine) changeToFollower() {
	log.Infof("[Store][database] change from leader to follower")
	atomic.StoreInt32(&le.leaderFlag, 0)
}

// checkLeaderDead
func (le *leaderElectionStateMachine) checkLeaderDead() (bool, error) {
	return le.leStore.CheckMtimeExpired(DefaultElectKey, LeaseTime)
}

// elect
func (le *leaderElectionStateMachine) elect() (bool, error) {
	curVersion, err := le.leStore.GetVersion(DefaultElectKey)
	if err != nil {
		return false, err
	}
	le.version = curVersion + 1
	return le.leStore.CompareAndSwapVersion(DefaultElectKey, curVersion, le.version, utils.LocalHost)
}

// heartbeat
func (le *leaderElectionStateMachine) heartbeat() (bool, error) {
	curVersion := le.version
	le.version = curVersion + 1
	return le.leStore.CompareAndSwapVersion(DefaultElectKey, curVersion, le.version, utils.LocalHost)
}

// isLeader
func (le *leaderElectionStateMachine) isLeader() bool {
	return isLeader(le.leaderFlag)
}

// isLeaderAtomic
func (le *leaderElectionStateMachine) isLeaderAtomic() bool {
	return isLeader(atomic.LoadInt32(&le.leaderFlag))
}

// StartLeaderElection
func (m *maintainStore) StartLeaderElection() error {
	ctx, cancel := context.WithCancel(context.TODO())
	le := &leaderElectionStateMachine{
		leStore:    &leaderElectionStore{master: m.master},
		leaderFlag: 0,
		version:    0,
		ctx:        ctx,
		cancel:     cancel,
	}
	m.le = le
	go le.mainLoop()
	return nil
}

// StopLeaderElection
func (m *maintainStore) StopLeaderElection() {
	if m.le != nil {
		m.le.cancel()
	}
	m.le = nil
}

// IsLeader
func (maintain *maintainStore) IsLeader() bool {
	return maintain.le.isLeaderAtomic()
}

// BatchCleanDeletedInstances batch clean soft deleted instances
func (maintain *maintainStore) BatchCleanDeletedInstances(batchSize uint32) (uint32, error) {
	log.Infof("[Store][database] batch clean soft deleted instances(%d)", batchSize)
	mainStr := "delete from instance where flag = 1 limit ?"
	result, err := maintain.master.Exec(mainStr, batchSize)
	if err != nil {
		log.Errorf("[Store][database] batch clean soft deleted instances(%d), err: %s", batchSize, err.Error())
		return 0, store.Error(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Warnf("[Store][database] batch clean soft deleted instances(%d), get RowsAffected err: %s", batchSize, err.Error())
		return 0, store.Error(err)
	}

	return uint32(rows), nil
}
