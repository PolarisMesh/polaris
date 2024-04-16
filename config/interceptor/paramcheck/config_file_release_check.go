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

package paramcheck

import (
	"context"
	"strconv"

	apiconfig "github.com/polarismesh/specification/source/go/api/v1/config_manage"
	apimodel "github.com/polarismesh/specification/source/go/api/v1/model"

	api "github.com/polarismesh/polaris/common/api/v1"
	"github.com/polarismesh/polaris/common/utils"
)

// PublishConfigFile 发布配置文件
func (s *Server) PublishConfigFile(ctx context.Context,
	configFileRelease *apiconfig.ConfigFileRelease) *apiconfig.ConfigResponse {

	return s.nextServer.PublishConfigFile(ctx, configFileRelease)
}

// GetConfigFileRelease 获取配置文件发布内容
func (s *Server) GetConfigFileRelease(ctx context.Context,
	req *apiconfig.ConfigFileRelease) *apiconfig.ConfigResponse {
	if errCode, errMsg := checkBaseReleaseParam(req, false); errCode != apimodel.Code_ExecuteSuccess {
		return api.NewConfigResponseWithInfo(errCode, errMsg)
	}
	return s.nextServer.GetConfigFileRelease(ctx, req)
}

// DeleteConfigFileReleases implements ConfigCenterServer.
func (s *Server) DeleteConfigFileReleases(ctx context.Context,
	reqs []*apiconfig.ConfigFileRelease) *apiconfig.ConfigBatchWriteResponse {
	responses := api.NewConfigBatchWriteResponse(apimodel.Code_ExecuteSuccess)
	chs := make([]chan *apiconfig.ConfigResponse, 0, len(reqs))
	for i, instance := range reqs {
		chs = append(chs, make(chan *apiconfig.ConfigResponse))
		go func(index int, ins *apiconfig.ConfigFileRelease) {
			chs[index] <- s.DeleteConfigFileRelease(ctx, ins)
		}(i, instance)
	}

	for _, ch := range chs {
		resp := <-ch
		api.ConfigCollect(responses, resp)
	}
	return responses
}

func (s *Server) DeleteConfigFileRelease(ctx context.Context, req *apiconfig.ConfigFileRelease) *apiconfig.ConfigResponse {
	if errCode, errMsg := checkBaseReleaseParam(req, true); errCode != apimodel.Code_ExecuteSuccess {
		return api.NewConfigResponseWithInfo(errCode, errMsg)
	}
	return s.nextServer.DeleteConfigFileRelease(ctx, req)
}

// GetConfigFileReleaseVersions implements ConfigCenterServer.
func (s *Server) GetConfigFileReleaseVersions(ctx context.Context,
	filters map[string]string) *apiconfig.ConfigBatchQueryResponse {

	searchFilters := map[string]string{}
	for k, v := range filters {
		if nk, ok := availableSearch["config_file_release"][k]; ok {
			searchFilters[nk] = v
		}
	}

	namespace := searchFilters["namespace"]
	group := searchFilters["group"]
	fileName := searchFilters["file_name"]
	if namespace == "" {
		return api.NewConfigBatchQueryResponseWithInfo(apimodel.Code_BadRequest, "invalid namespace")
	}
	if group == "" {
		return api.NewConfigBatchQueryResponseWithInfo(apimodel.Code_BadRequest, "invalid config group")
	}
	if fileName == "" {
		return api.NewConfigBatchQueryResponseWithInfo(apimodel.Code_BadRequest, "invalid config file name")
	}

	return s.nextServer.GetConfigFileReleaseVersions(ctx, searchFilters)
}

// GetConfigFileReleases implements ConfigCenterServer.
func (s *Server) GetConfigFileReleases(ctx context.Context,
	filters map[string]string) *apiconfig.ConfigBatchQueryResponse {

	offset, limit, err := utils.ParseOffsetAndLimit(filters)
	if err != nil {
		return api.NewConfigBatchQueryResponseWithInfo(apimodel.Code_BadRequest, err.Error())
	}

	searchFilters := map[string]string{
		"offset": strconv.Itoa(int(offset)),
		"limit":  strconv.Itoa(int(limit)),
	}
	for k, v := range filters {
		if nK, ok := availableSearch["config_file_release"][k]; ok {
			searchFilters[nK] = v
		}
	}

	return s.nextServer.GetConfigFileReleases(ctx, searchFilters)
}

// RollbackConfigFileReleases implements ConfigCenterServer.
func (s *Server) RollbackConfigFileReleases(ctx context.Context,
	reqs []*apiconfig.ConfigFileRelease) *apiconfig.ConfigBatchWriteResponse {

	responses := api.NewConfigBatchWriteResponse(apimodel.Code_ExecuteSuccess)
	chs := make([]chan *apiconfig.ConfigResponse, 0, len(reqs))
	for i, instance := range reqs {
		chs = append(chs, make(chan *apiconfig.ConfigResponse))
		go func(index int, ins *apiconfig.ConfigFileRelease) {
			chs[index] <- s.RollbackConfigFileRelease(ctx, ins)
		}(i, instance)
	}

	for _, ch := range chs {
		resp := <-ch
		api.ConfigCollect(responses, resp)
	}
	return responses
}

func (s *Server) RollbackConfigFileRelease(ctx context.Context,
	req *apiconfig.ConfigFileRelease) *apiconfig.ConfigResponse {
	if errCode, errMsg := checkBaseReleaseParam(req, true); errCode != apimodel.Code_ExecuteSuccess {
		return api.NewConfigResponseWithInfo(errCode, errMsg)
	}
	return s.nextServer.RollbackConfigFileRelease(ctx, req)
}

// UpsertAndReleaseConfigFile .
func (s *Server) UpsertAndReleaseConfigFile(ctx context.Context,
	req *apiconfig.ConfigFilePublishInfo) *apiconfig.ConfigResponse {

	return s.nextServer.UpsertAndReleaseConfigFile(ctx, req)
}

func (s *Server) StopGrayConfigFileReleases(ctx context.Context,
	reqs []*apiconfig.ConfigFileRelease) *apiconfig.ConfigBatchWriteResponse {

	return s.nextServer.StopGrayConfigFileReleases(ctx, reqs)
}

func checkBaseReleaseParam(req *apiconfig.ConfigFileRelease, checkRelease bool) (apimodel.Code, string) {
	namespace := req.GetNamespace().GetValue()
	group := req.GetGroup().GetValue()
	fileName := req.GetFileName().GetValue()
	releaseName := req.GetName().GetValue()
	if namespace == "" {
		return apimodel.Code_BadRequest, "invalid namespace"
	}
	if group == "" {
		return apimodel.Code_BadRequest, "invalid config group"
	}
	if fileName == "" {
		return apimodel.Code_BadRequest, "invalid config file name"
	}
	if checkRelease {
		if releaseName == "" {
			return apimodel.Code_BadRequest, "invalid config release name"
		}
	}
	return apimodel.Code_ExecuteSuccess, ""
}
