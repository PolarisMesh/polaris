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

package service

import (
	"context"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"

	"github.com/polarismesh/polaris-server/auth"
	"github.com/polarismesh/polaris-server/cache"
	"github.com/polarismesh/polaris-server/common/model"
	"github.com/polarismesh/polaris-server/namespace"
	"github.com/polarismesh/polaris-server/plugin"
	"github.com/polarismesh/polaris-server/service/batch"
	"github.com/polarismesh/polaris-server/service/healthcheck"
	"github.com/polarismesh/polaris-server/store"
)

const (
	// MaxBatchSize max batch size
	MaxBatchSize = 100
	// MaxQuerySize max query size
	MaxQuerySize = 100
)

const (
	// SystemNamespace polaris system namespace
	SystemNamespace = "Polaris"
	// DefaultNamespace default namespace
	DefaultNamespace = "default"
	// ProductionNamespace default namespace
	ProductionNamespace = "Production"
	// DefaultTLL default ttl
	DefaultTLL = 5
)

var (
	server       DiscoverServer
	namingServer *Server = new(Server)
	once                 = sync.Once{}
	finishInit           = false
)

// Config 核心逻辑层配置
type Config struct {
	Auth  map[string]interface{} `yaml:"auth"`
	Batch map[string]interface{} `yaml:"batch"`
}

// Initialize 初始化
func Initialize(ctx context.Context, namingOpt *Config, cacheOpt *cache.Config,
	bc *batch.Controller, polarisServiceSet map[model.ServiceKey]struct{}) error {
	var err error
	once.Do(func() {
		err = initialize(ctx, namingOpt, cacheOpt, bc, polarisServiceSet)
	})

	if err != nil {
		return err
	}

	finishInit = true
	return nil
}

// GetServer 获取已经初始化好的Server
func GetServer() (DiscoverServer, error) {
	if !finishInit {
		return nil, errors.New("server has not done InitializeServer")
	}

	return server, nil
}

// GetOriginServer 获取已经初始化好的Server
func GetOriginServer() (*Server, error) {
	if !finishInit {
		return nil, errors.New("server has not done InitializeServer")
	}

	return namingServer, nil
}

// 内部初始化函数
func initialize(ctx context.Context, namingOpt *Config, cacheOpt *cache.Config, bc *batch.Controller,
	polarisServiceSet map[model.ServiceKey]struct{}) error {

	// 获取存储层对象
	s, err := store.GetStore()
	if err != nil {
		log.Errorf("[Naming][Server] can not get store, err: %s", err.Error())
		return errors.New("can not get store")
	}
	if s == nil {
		log.Errorf("[Naming][Server] store is null")
		return errors.New("store is null")
	}

	healthServer, err := healthcheck.GetServer()
	if err != nil {
		log.Errorf("[Naming][Server] can not get store, err: %s", err.Error())
		return errors.New("can not get healthcheck server")
	}

	namingServer.healthServer = healthServer

	namingServer.storage = s

	// 注入命名空间管理模块
	namespaceSvr, err := namespace.GetOriginServer()
	if err != nil {
		return err
	}
	namingServer.namespaceSvr = namespaceSvr

	// cache模块，可以不开启
	// 对于控制台集群，只访问控制台接口的，可以不开启cache
	if cacheOpt.Open {
		caches, cacheErr := cache.GetCacheManager()
		if cacheErr != nil {
			log.Errorf("[Naming][Server] new naming cache err: %s", cacheErr.Error())
			return cacheErr
		}
		log.Infof("[Naming][Server] cache is open, can access the client api function")
		namingServer.caches = caches
	}

	namingServer.bc = bc

	// l5service
	namingServer.l5service = &l5service{}

	namingServer.createServiceSingle = &singleflight.Group{}
	namingServer.createNamespaceSingle = &singleflight.Group{}
	namingServer.polarisServiceSet = polarisServiceSet

	// 插件初始化
	pluginInitialize()

	authServer, err := auth.GetAuthServer()
	if err != nil {
		return err
	}

	server = newServerAuthAbility(namingServer, authServer)

	return nil
}

// 插件初始化
func pluginInitialize() {
	// 获取CMDB插件
	namingServer.cmdb = plugin.GetCMDB()
	if namingServer.cmdb == nil {
		log.Warnf("Not Found CMDB Plugin")
	}

	// 获取History插件，注意：插件的配置在bootstrap已经设置好
	namingServer.history = plugin.GetHistory()
	if namingServer.history == nil {
		log.Warnf("Not Found History Log Plugin")
	}

	// 获取限流插件
	namingServer.ratelimit = plugin.GetRatelimit()
	if namingServer.ratelimit == nil {
		log.Warnf("Not found Ratelimit Plugin")
	}

	// 获取DiscoverStatis插件
	namingServer.discoverStatis = plugin.GetDiscoverStatis()
	if namingServer.discoverStatis == nil {
		log.Warnf("Not Found Discover Statis Plugin")
	}

	// 获取服务事件插件
	namingServer.discoverEvent = plugin.GetDiscoverEvent()
	if namingServer.discoverEvent == nil {
		log.Warnf("Not found DiscoverEvent Plugin")
	}

	// 获取鉴权插件
	namingServer.auth = plugin.GetAuth()
	if namingServer.auth == nil {
		log.Warnf("Not found Auth Plugin")
	}
}

// NewUUID 返回一个随机的UUID
func NewUUID() string {
	uuidBytes := uuid.New()
	return hex.EncodeToString(uuidBytes[:])
}
