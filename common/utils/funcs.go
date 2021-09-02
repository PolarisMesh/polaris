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

package utils

import (
	"encoding/hex"
	uuid "github.com/google/uuid"
	api "github.com/polarismesh/polaris-server/common/api/v1"
	"github.com/polarismesh/polaris-server/common/model"
	"strings"
)

// 创建存储层服务实例模型
func CreateInstanceModel(serviceID string, req *api.Instance) *model.Instance {
	// 默认为健康的
	healthy := true
	if req.GetHealthy() != nil {
		healthy = req.GetHealthy().GetValue()
	}

	// 默认为不隔离的
	isolate := false
	if req.GetIsolate() != nil {
		isolate = req.GetIsolate().GetValue()
	}

	// 权重默认是100
	var weight uint32 = 100
	if req.GetWeight() != nil {
		weight = req.GetWeight().GetValue()
	}

	instance := &model.Instance{
		ServiceID: serviceID,
	}
	uuidBytes := uuid.New()
	protoIns := &api.Instance{
		Id:       req.GetId(),
		Host:     NewStringValue(strings.TrimSpace(req.GetHost().GetValue())),
		VpcId:    req.GetVpcId(),
		Port:     req.GetPort(),
		Protocol: req.GetProtocol(),
		Version:  req.GetVersion(),
		Priority: req.GetPriority(),
		Weight:   NewUInt32Value(weight),
		Healthy:  NewBoolValue(healthy),
		Isolate:  NewBoolValue(isolate),
		Location: req.Location,
		Metadata: req.Metadata,
		LogicSet: req.GetLogicSet(),
		Revision: NewStringValue(hex.EncodeToString(uuidBytes[:])), // 更新版本号
	}

	// health Check，healthCheck不能为空，且没有显示把enable_health_check置为false
	// 如果create的时候，打开了healthCheck，那么实例模式是unhealthy，必须要一次心跳才会healthy
	if req.GetHealthCheck().GetHeartbeat() != nil &&
		(req.GetEnableHealthCheck() == nil || req.GetEnableHealthCheck().GetValue()) {
		protoIns.EnableHealthCheck = NewBoolValue(true)
		protoIns.HealthCheck = req.HealthCheck
		protoIns.HealthCheck.Type = api.HealthCheck_HEARTBEAT
		// ttl range: (0, 60]
		ttl := protoIns.GetHealthCheck().GetHeartbeat().GetTtl().GetValue()
		if ttl == 0 || ttl > 60 {
			if protoIns.HealthCheck.Heartbeat.Ttl == nil {
				protoIns.HealthCheck.Heartbeat.Ttl = NewUInt32Value(5)
			}
			protoIns.HealthCheck.Heartbeat.Ttl.Value = 5
		}
		// 开启健康检查，且没有代入健康状态，则健康状态默认都是false
		protoIns.Healthy.Value = false
	}

	instance.Proto = protoIns
	return instance
}
