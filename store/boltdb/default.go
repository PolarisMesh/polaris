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

package boltdb

import (
	"errors"
	"github.com/polarismesh/polaris-server/store"
)

const (
	storeName = "boltdbStore"
)

type boltStore struct {
	*namespaceStore
	*businessStore
	*serviceStore
	*instanceStore
	*l5Store
	*routingStore
	*rateLimitStore
	*platformStore
	*circuitBreakerStore
	*toolStore

	handler BoltHandler
	start   bool
}

// Name store name
func (m *boltStore) Name() string {
	return storeName
}

// Initialize init store
func (m *boltStore) Initialize(c *store.Config) error {
	if m.start {
		return errors.New("store has been Initialize")
	}
	boltConfig := &BoltConfig{}
	boltConfig.Parse(c.Option)
	handler, err := NewBoltHandler(boltConfig)
	if nil != err {
		return err
	}
	m.handler = handler
	if err = m.newStore(); nil != err {
		_ = handler.Close()
		return err
	}
	m.start = true
	return nil
}

func (m *boltStore) newStore() error {
	m.l5Store = &l5Store{handler: m.handler}
	if err := m.l5Store.InitL5Data(); nil != err {
		return err
	}

	m.namespaceStore = &namespaceStore{handler: m.handler}
	if err := m.namespaceStore.InitData(); nil != err {
		return err
	}
	m.businessStore = &businessStore{handler: m.handler}

	m.serviceStore = &serviceStore{handler: m.handler}

	m.instanceStore = &instanceStore{handler: m.handler}

	m.routingStore = &routingStore{handler: m.handler}

	m.rateLimitStore = &rateLimitStore{handler: m.handler}

	m.circuitBreakerStore = &circuitBreakerStore{handler: m.handler}

	m.platformStore = &platformStore{handler: m.handler}

	m.toolStore = &toolStore{handler: m.handler}

	return nil
}

// Destroy destroy store
func (m *boltStore) Destroy() error {
	if nil != m.handler {
		return m.handler.Close()
	}
	return nil
}

//CreateTransaction create store transaction
func (m *boltStore) CreateTransaction() (store.Transaction, error) {
	tx, err := m.handler.Transaction()
	if nil != err {
		return nil, err
	}
	return &transaction{tx: tx}, nil
}

func init() {
	s := &boltStore{}
	_ = store.RegisterStore(s)
}
