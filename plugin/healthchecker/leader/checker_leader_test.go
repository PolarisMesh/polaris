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

package leader

import (
	"context"
<<<<<<< HEAD
=======
	"sync/atomic"
>>>>>>> 2d0c2e97... feat:support heartbeat without redis in cluster
	"testing"
	"time"

	"github.com/golang/mock/gomock"
<<<<<<< HEAD
	"github.com/stretchr/testify/assert"

=======
>>>>>>> 2d0c2e97... feat:support heartbeat without redis in cluster
	"github.com/polarismesh/polaris/common/batchjob"
	"github.com/polarismesh/polaris/common/model"
	"github.com/polarismesh/polaris/common/utils"
	"github.com/polarismesh/polaris/store"
	"github.com/polarismesh/polaris/store/mock"
<<<<<<< HEAD
=======
	"github.com/stretchr/testify/assert"
>>>>>>> 2d0c2e97... feat:support heartbeat without redis in cluster
)

func TestLeaderHealthChecker_OnEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})
	mockStore := mock.NewMockStore(ctrl)

	checker := &LeaderHealthChecker{
<<<<<<< HEAD
		self: NewLocalPeerFunc(),
		s:    mockStore,
=======
		s: mockStore,
>>>>>>> 2d0c2e97... feat:support heartbeat without redis in cluster
		conf: &Config{
			SoltNum: 0,
			Batch: batchjob.CtrlConfig{
				QueueSize:     16,
				WaitTime:      32 * time.Millisecond,
				MaxBatchCount: 32,
				Concurrency:   1,
			},
		},
	}

	mockPort := uint32(28888)
	_, err := newMockPolarisGRPCSever(t, mockPort)
	assert.NoError(t, err)

	utils.LocalHost = "127.0.0.2"
	utils.LocalPort = int(mockPort)
	t.Cleanup(func() {
		utils.LocalPort = 8091
		utils.LocalHost = "127.0.0.1"
	})

	t.Run("initialize-self-is-follower", func(t *testing.T) {
<<<<<<< HEAD
=======
		callTimes := int64(0)
		oldNewRemotePeerFunc := NewRemotePeerFunc
		NewRemotePeerFunc = func() Peer {
			p := oldNewRemotePeerFunc()
			return &MockPeerImpl{
				OnServe: func(ctx context.Context, _ *MockPeerImpl, ip string, port uint32) error {
					atomic.AddInt64(&callTimes, 1)
					return p.Serve(ctx, ip, port)
				},
				OnGet: func(key string) (*ReadBeatRecord, error) {
					atomic.AddInt64(&callTimes, 1)
					return p.Get(key)
				},
				OnPut: func(r WriteBeatRecord) error {
					atomic.AddInt64(&callTimes, 1)
					return p.Put(r)
				},
				OnDel: func(key string) error {
					atomic.AddInt64(&callTimes, 1)
					return p.Del(key)
				},
				OnClose: func(*MockPeerImpl) error {
					atomic.AddInt64(&callTimes, 1)
					return p.Close()
				},
				OnHost: func() string {
					return p.Host()
				},
			}
		}

>>>>>>> 2d0c2e97... feat:support heartbeat without redis in cluster
		mockStore.EXPECT().ListLeaderElections().Times(1).Return([]*model.LeaderElection{
			{
				ElectKey: electionKey,
				Host:     "127.0.0.1",
			},
		}, nil)

		checker.OnEvent(context.Background(), store.LeaderChangeEvent{
			Key:        electionKey,
			Leader:     false,
			LeaderHost: "127.0.0.1",
		})

		assert.True(t, checker.isInitialize())
		assert.False(t, checker.isLeader())

		skipRet := checker.skipCheck(utils.NewUUID(), 15)
		assert.True(t, skipRet)
		time.Sleep(15 * time.Second)
		skipRet = checker.skipCheck(utils.NewUUID(), 15)
		assert.False(t, skipRet)

		peer := checker.findLeaderPeer()
		assert.NotNil(t, peer)
		_, ok := peer.(*RemotePeer)
		assert.True(t, ok)
	})

	t.Run("initialize-self-become-leader", func(t *testing.T) {
		checker.OnEvent(context.Background(), store.LeaderChangeEvent{
			Key:        electionKey,
			Leader:     true,
			LeaderHost: "127.0.0.2",
		})

		assert.True(t, checker.isInitialize())
		assert.True(t, checker.isLeader())
		assert.Nil(t, checker.remote)

		skipRet := checker.skipCheck(utils.NewUUID(), 15)
		assert.True(t, skipRet)
		time.Sleep(15 * time.Second)
		skipRet = checker.skipCheck(utils.NewUUID(), 15)
		assert.False(t, skipRet)

		peer := checker.findLeaderPeer()
		assert.NotNil(t, peer)
		_, ok := peer.(*LocalPeer)
		assert.True(t, ok)
	})
}
