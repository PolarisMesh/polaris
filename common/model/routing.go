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

package model

import (
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	apitraffic "github.com/polarismesh/specification/source/go/api/v1/traffic_manage"
	protoV2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	commontime "github.com/polarismesh/polaris/common/time"
	"github.com/polarismesh/polaris/common/utils"
)

const (
	// V2RuleIDKey v2 版本的规则路由 ID
	V2RuleIDKey = "__routing_v2_id__"
	// V1RuleIDKey v1 版本的路由规则 ID
	V1RuleIDKey = "__routing_v1_id__"
	// V1RuleRouteIndexKey v1 版本 route 规则在自己 route 链中的 index 信息
	V1RuleRouteIndexKey = "__routing_v1_route_index__"
	// V1RuleRouteTypeKey 标识当前 v2 路由规则在 v1 的 inBound 还是 outBound
	V1RuleRouteTypeKey = "__routing_v1_route_type__"
	// V1RuleInRoute inBound 类型
	V1RuleInRoute = "in"
	// V1RuleOutRoute outBound 类型
	V1RuleOutRoute = "out"
)

var (
	// RuleRoutingTypeUrl 记录 anypb.Any 中关于 RuleRoutingConfig 的 url 信息
	RuleRoutingTypeUrl string
	// MetaRoutingTypeUrl 记录 anypb.Any 中关于 MetadataRoutingConfig 的 url 信息
	MetaRoutingTypeUrl string
	// RuleRoutingTypeUrlV2 记录 anypb.Any 中关于 RuleRoutingConfigV2 的 url 信息
	// RuleRoutingTypeUrlV2 string
	// MetaRoutingTypeUrlV2 记录 anypb.Any 中关于 MetadataRoutingConfigV2 的 url 信息
	// MetaRoutingTypeUrlV2 string
)

func init() {
	ruleAny, _ := ptypes.MarshalAny(&apitraffic.RuleRoutingConfig{})
	metaAny, _ := ptypes.MarshalAny(&apitraffic.MetadataRoutingConfig{})

	RuleRoutingTypeUrl = ruleAny.GetTypeUrl()
	MetaRoutingTypeUrl = metaAny.GetTypeUrl()

	// ruleAnyV2, _ := ptypes.MarshalAny(&v2.RuleRoutingConfig{})
	// metaAnyV2, _ := ptypes.MarshalAny(&v2.MetadataRoutingConfig{})
	// RuleRoutingTypeUrlV2 = ruleAnyV2.GetTypeUrl()
	// MetaRoutingTypeUrlV2 = metaAnyV2.GetTypeUrl()

}

// ExtendRouterConfig 路由信息的扩展
type ExtendRouterConfig struct {
	*RouterConfig
	// MetadataRouting 元数据路由配置
	MetadataRouting *apitraffic.MetadataRoutingConfig
	// RuleRouting 规则路由配置
	RuleRouting *apitraffic.RuleRoutingConfig
	// ExtendInfo 额外信息数据
	ExtendInfo map[string]string
}

// ToApi 转为 api 对象
func (r *ExtendRouterConfig) ToApi() (*apitraffic.RouteRule, error) {
	var (
		anyValue *anypb.Any
		err      error
	)

	if r.GetRoutingPolicy() == apitraffic.RoutingPolicy_MetadataPolicy {
		anyValue, err = ptypes.MarshalAny(r.MetadataRouting)
		if err != nil {
			return nil, err
		}
	} else {
		anyValue, err = ptypes.MarshalAny(r.RuleRouting)
		if err != nil {
			return nil, err
		}
	}

	return &apitraffic.RouteRule{
		Id:            r.ID,
		Name:          r.Name,
		Namespace:     r.Namespace,
		Enable:        r.Enable,
		RoutingPolicy: r.GetRoutingPolicy(),
		RoutingConfig: anyValue,
		Revision:      r.Revision,
		Ctime:         commontime.Time2String(r.CreateTime),
		Mtime:         commontime.Time2String(r.ModifyTime),
		Etime:         commontime.Time2String(r.EnableTime),
		Priority:      r.Priority,
		Description:   r.Description,
	}, nil
}

// RouterConfig 路由规则
type RouterConfig struct {
	// ID 规则唯一标识
	ID string `json:"id"`
	// namespace 所属的命名空间
	Namespace string `json:"namespace"`
	// name 规则名称
	Name string `json:"name"`
	// policy 规则类型
	Policy string `json:"policy"`
	// config 具体的路由规则内容
	Config string `json:"config"`
	// enable 路由规则是否启用
	Enable bool `json:"enable"`
	// priority 规则优先级
	Priority uint32 `json:"priority"`
	// revision 路由规则的版本信息
	Revision string `json:"revision"`
	// Description 规则简单描述
	Description string `json:"description"`
	// valid 路由规则是否有效，没有被逻辑删除
	Valid bool `json:"flag"`
	// createtime 规则创建时间
	CreateTime time.Time `json:"ctime"`
	// modifytime 规则修改时间
	ModifyTime time.Time `json:"mtime"`
	// enabletime 规则最近一次启用时间
	EnableTime time.Time `json:"etime"`
}

// GetRoutingPolicy 查询路由规则类型
func (r *RouterConfig) GetRoutingPolicy() apitraffic.RoutingPolicy {
	v, ok := apitraffic.RoutingPolicy_value[r.Policy]

	if !ok {
		return apitraffic.RoutingPolicy(-1)
	}

	return apitraffic.RoutingPolicy(v)
}

// ToExpendRoutingConfig 转为扩展对象，提前序列化出相应的 pb struct
func (r *RouterConfig) ToExpendRoutingConfig() (*ExtendRouterConfig, error) {
	ret := &ExtendRouterConfig{
		RouterConfig: r,
	}

	configText := r.Config
	if len(configText) == 0 {
		return ret, nil
	}
	policy := r.GetRoutingPolicy()
	var err error
	if strings.HasPrefix(configText, "{") {
		// process with json
		switch policy {
		case apitraffic.RoutingPolicy_RulePolicy:
			rule := &apitraffic.RuleRoutingConfig{}
			if err = utils.UnmarshalFromJsonString(rule, configText); nil != err {
				return nil, err
			}
			ret.RuleRouting = rule
			break
		case apitraffic.RoutingPolicy_MetadataPolicy:
			rule := &apitraffic.MetadataRoutingConfig{}
			if err = utils.UnmarshalFromJsonString(rule, configText); nil != err {
				return nil, err
			}
			ret.MetadataRouting = rule
			break
		}
		return ret, nil
	}

	err = r.parseBinaryAnyMessage(policy, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (r *RouterConfig) parseBinaryAnyMessage(
	policy apitraffic.RoutingPolicy, ret *ExtendRouterConfig) error {
	// parse v1 binary
	switch policy {
	case apitraffic.RoutingPolicy_RulePolicy:
		rule := &apitraffic.RuleRoutingConfig{}
		anyMsg := &anypb.Any{
			TypeUrl: RuleRoutingTypeUrl,
			Value:   []byte(r.Config),
		}
		if err := unmarshalToAny(anyMsg, rule); nil != err {
			// parse v2 binary
			// ruleV2 := &v2.RuleRoutingConfig{}
			// anyMsg = &anypb.Any{
			//	 TypeUrl: RuleRoutingTypeUrlV2,
			//	 Value:   []byte(r.Config),
			// }
			// if err = unmarshalToAny(anyMsg, ruleV2); nil != err {
			//	 return err
			// }
			// if err = utils.ConvertSameStructureMessage(ruleV2, rule); nil != err {
			//	 return err
			// }
			return err
		}
		ret.RuleRouting = rule
	case apitraffic.RoutingPolicy_MetadataPolicy:
		rule := &apitraffic.MetadataRoutingConfig{}
		anyMsg := &anypb.Any{
			TypeUrl: MetaRoutingTypeUrl,
			Value:   []byte(r.Config),
		}
		if err := unmarshalToAny(anyMsg, rule); nil != err {
			// parse v2 binary
			// ruleV2 := &v2.MetadataRoutingConfig{}
			// anyMsg = &anypb.Any{
			//	 TypeUrl: MetaRoutingTypeUrlV2,
			//	 Value:   []byte(r.Config),
			// }
			// if err = unmarshalToAny(anyMsg, ruleV2); nil != err {
			//	 return err
			// }
			// if err = utils.ConvertSameStructureMessage(ruleV2, rule); nil != err {
			//	 return err
			// }
			return err
		}
		ret.MetadataRouting = rule
	}
	return nil
}

// ParseRouteRuleFromAPI 从 API 对象中转换出内部对象
func (r *RouterConfig) ParseRouteRuleFromAPI(routing *apitraffic.RouteRule) error {
	ruleMessage, err := ParseRouteRuleAnyToMessage(routing.RoutingPolicy, routing.RoutingConfig)
	if nil != err {
		return err
	}
	if r.Config, err = utils.MarshalToJsonString(ruleMessage); nil != err {
		return err
	}
	r.ID = routing.Id
	r.Revision = routing.Revision
	r.Name = routing.Name
	r.Namespace = routing.Namespace
	r.Enable = routing.Enable
	r.Policy = routing.GetRoutingPolicy().String()
	r.Priority = routing.Priority
	r.Description = routing.Description

	// 优先级区间范围 [0, 10]
	if r.Priority > 10 {
		r.Priority = 10
	}

	return nil
}

func unmarshalToAny(anyMessage *anypb.Any, message proto.Message) error {
	return anypb.UnmarshalTo(anyMessage, proto.MessageV2(message),
		protoV2.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true})
}

// ParseRouteRuleAnyToMessage convert the any routing proto to message object
func ParseRouteRuleAnyToMessage(policy apitraffic.RoutingPolicy, anyMessage *anypb.Any) (proto.Message, error) {
	var rule proto.Message
	switch policy {
	case apitraffic.RoutingPolicy_RulePolicy:
		rule = &apitraffic.RuleRoutingConfig{}
		if err := unmarshalToAny(anyMessage, rule); err != nil {
			return nil, err
		}
		break
	case apitraffic.RoutingPolicy_MetadataPolicy:
		rule = &apitraffic.MetadataRoutingConfig{}
		if err := unmarshalToAny(anyMessage, rule); err != nil {
			return nil, err
		}
		break
	default:
		break
	}
	return rule, nil
}
