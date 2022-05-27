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

	"go.uber.org/zap"

	"github.com/polarismesh/polaris-server/cache"
	api "github.com/polarismesh/polaris-server/common/api/v1"
	"github.com/polarismesh/polaris-server/common/log"
	"github.com/polarismesh/polaris-server/common/utils"
	utils2 "github.com/polarismesh/polaris-server/config/utils"
)

type checkFunc func(clientConfigFile *api.ClientConfigFileInfo, cacheEntry *cache.Entry) bool

// CheckClientConfigFileByVersion 通过比较版本号检查客户端使用的配置文件是否版本落后
func (cs *Impl) CheckClientConfigFileByVersion(ctx context.Context, configFiles []*api.ClientConfigFileInfo) *api.ConfigClientResponse {
	return cs.doCheckClientConfigFile(ctx, configFiles, func(clientConfigFile *api.ClientConfigFileInfo, cacheEntry *cache.Entry) bool {
		return !cacheEntry.Empty && clientConfigFile.Version.GetValue() < cacheEntry.Version
	})
}

// CheckClientConfigFileByMd5 通过比较md5检查客户端使用的配置文件是否版本落后
func (cs *Impl) CheckClientConfigFileByMd5(ctx context.Context, configFiles []*api.ClientConfigFileInfo) *api.ConfigClientResponse {
	return cs.doCheckClientConfigFile(ctx, configFiles, func(clientConfigFile *api.ClientConfigFileInfo, cacheEntry *cache.Entry) bool {
		return clientConfigFile.Md5.GetValue() != cacheEntry.Md5
	})
}

func (cs *Impl) doCheckClientConfigFile(ctx context.Context, configFiles []*api.ClientConfigFileInfo,
	checkFunc checkFunc) *api.ConfigClientResponse {
	if len(configFiles) == 0 {
		return api.NewConfigClientResponse(api.InvalidWatchConfigFileFormat, nil)
	}

	requestID, _ := ctx.Value(utils.StringContext("request-id")).(string)

	for _, configFile := range configFiles {
		namespace := configFile.Namespace.GetValue()
		group := configFile.Group.GetValue()
		fileName := configFile.FileName.GetValue()

		if namespace == "" || group == "" || fileName == "" {
			return api.NewConfigClientResponseWithMessage(api.BadRequest, "namespace & group & fileName can not be empty")
		}

		// 从缓存中获取最新的配置文件信息
		entry, err := cs.cache.GetOrLoadIfAbsent(namespace, group, fileName)

		if err != nil {
			log.ConfigScope().Error("[Config][Service] get or load config file from cache error.",
				zap.String("requestId", requestID),
				zap.String("fileName", fileName),
				zap.Error(err))

			return api.NewConfigClientResponse(api.ExecuteException, nil)
		}

		if checkFunc(configFile, entry) {
			return utils2.GenConfigFileResponse(namespace, group, fileName, "", entry.Md5, entry.Version)
		}
	}

	return api.NewConfigClientResponse(api.DataNoChange, nil)
}

// GetConfigFileForClient 从缓存中获取配置文件，如果客户端的版本号大于服务端，则服务端重新加载缓存
func (cs *Impl) GetConfigFileForClient(ctx context.Context, namespace, group, fileName string, clientVersion uint64) *api.ConfigClientResponse {
	if namespace == "" || group == "" || fileName == "" {
		return api.NewConfigClientResponseWithMessage(api.BadRequest, "namespace & group & fileName can not be empty")
	}

	requestID, _ := ctx.Value(utils.StringContext("request-id")).(string)

	log.ConfigScope().Info("[Config][Service] load config file from cache.",
		zap.String("requestId", requestID), zap.String("namespace", namespace), zap.String("group", group), zap.String("file", fileName))

	// 从缓存中获取配置内容
	entry, err := cs.cache.GetOrLoadIfAbsent(namespace, group, fileName)

	if err != nil {
		log.ConfigScope().Error("[Config][Service] get or load config file from cache error.",
			zap.String("requestId", requestID),
			zap.Error(err))

		return api.NewConfigClientResponseWithMessage(api.ExecuteException, "load config file error")
	}

	if entry.Empty {
		return api.NewConfigClientResponse(api.NotFoundResource, nil)
	}

	// 客户端版本号大于服务端版本号，服务端需要重新加载缓存
	if clientVersion > entry.Version {
		entry, err = cs.cache.ReLoad(namespace, group, fileName)
		if err != nil {
			log.ConfigScope().Error("[Config][Service] reload config file error.",
				zap.String("requestId", requestID),
				zap.Error(err))

			return api.NewConfigClientResponseWithMessage(api.ExecuteException, "load config file error")
		}
	}

	return utils2.GenConfigFileResponse(namespace, group, fileName, entry.Content, entry.Md5, entry.Version)
}
