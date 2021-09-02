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
	"time"

	v1 "github.com/polarismesh/polaris-server/common/api/v1"
	"github.com/golang/protobuf/ptypes/wrappers"
)

/**
 * @brief 命名空间结构体
 */
type Namespace struct {
	Name       string
	Comment    string
	Token      string
	Owner      string
	Valid      bool
	CreateTime time.Time
	ModifyTime time.Time
}

/**
 * @brief 业务集
 */
type Business struct {
	ID         string
	Name       string
	Token      string
	Owner      string
	Valid      bool
	CreateTime time.Time
	ModifyTime time.Time
}

/**
 * @brief 服务数据
 */
type Service struct {
	ID          string
	Name        string
	Namespace   string
	Business    string
	Ports       string
	Meta        map[string]string
	Comment     string
	Department  string
	CmdbMod1    string
	CmdbMod2    string
	CmdbMod3    string
	Token       string
	Owner       string
	Revision    string
	Reference   string
	ReferFilter string
	PlatformID  string
	Valid       bool
	CreateTime  time.Time
	ModifyTime  time.Time
}

/**
 * @brief 服务名
 */
type ServiceKey struct {
	Namespace string
	Name      string
}

// 便捷函数封装
func (s *Service) IsAlias() bool {
	if s.Reference != "" {
		return true
	}

	return false
}

// 服务别名结构体
type ServiceAlias struct {
	ID         string
	Alias      string
	ServiceID  string
	Service    string
	Namespace  string
	Owner      string
	Comment    string
	CreateTime time.Time
	ModifyTime time.Time
}

/**
 * @brief 服务下实例的权重类型
 */
type WeightType uint32

const (
	// 动态权重
	WEIGHTDYNAMIC WeightType = 0

	// 静态权重
	WEIGHTSTATIC WeightType = 1
)

var WeightString = map[WeightType]string{
	WEIGHTDYNAMIC: "dynamic",
	WEIGHTSTATIC:  "static",
}

var WeightEnum = map[string]WeightType{
	"dynamic": WEIGHTDYNAMIC,
	"static":  WEIGHTSTATIC,
}

/**
 * @brief 地域信息，对应数据库字段
 */
type LocationStore struct {
	IP         string
	Region     string
	Zone       string
	Campus     string
	RegionID   uint32
	ZoneID     uint32
	CampusID   uint32
	Flag       int
	ModifyTime int64
}

// cmdb信息，对应内存结构体
type Location struct {
	Proto    *v1.Location
	RegionID uint32
	ZoneID   uint32
	CampusID uint32
	Valid    bool
}

// 转成内存数据结构
func Store2Location(s *LocationStore) *Location {
	return &Location{
		Proto: &v1.Location{
			Region: &wrappers.StringValue{Value: s.Region},
			Zone:   &wrappers.StringValue{Value: s.Zone},
			Campus: &wrappers.StringValue{Value: s.Campus},
		},
		RegionID: s.RegionID,
		ZoneID:   s.ZoneID,
		CampusID: s.CampusID,
		Valid:    flag2valid(s.Flag),
	}
}

/**
 * @brief 客户端上报信息表
 */
type Client struct {
	VpcID      string
	Host       string
	Typ        int
	Version    string
	CreateTime time.Time
	ModifyTime time.Time
}

/**
 * @brief 路由配置
 */
type RoutingConfig struct {
	ID         string
	InBounds   string
	OutBounds  string
	Revision   string
	Valid      bool
	CreateTime time.Time
	ModifyTime time.Time
}

/**
 * @brief 路由配置的扩展结构体
 */
type ExtendRoutingConfig struct {
	ServiceName   string
	NamespaceName string
	Config        *RoutingConfig
}

/**
 * @brief 限流规则
 */
type RateLimit struct {
	ID         string
	ServiceID  string
	ClusterID  string
	Labels     string
	Priority   uint32
	Rule       string
	Revision   string
	Valid      bool
	CreateTime time.Time
	ModifyTime time.Time
}

/**
 * @brief 包含服务信息的限流规则
 */
type ExtendRateLimit struct {
	ServiceName   string
	NamespaceName string
	RateLimit     *RateLimit
}

/**
 * @brief 包含最新版本号的限流规则
 */
type RateLimitRevision struct {
	ServiceID    string
	LastRevision string
}

/**
 * @brief 熔断规则
 */
type CircuitBreaker struct {
	ID         string
	Version    string
	Name       string
	Namespace  string
	Business   string
	Department string
	Comment    string
	Inbounds   string
	Outbounds  string
	Token      string
	Owner      string
	Revision   string
	Valid      bool
	CreateTime time.Time
	ModifyTime time.Time
}

/**
 * @brief 与服务关系绑定的熔断规则
 */
type ServiceWithCircuitBreaker struct {
	ServiceID      string
	CircuitBreaker *CircuitBreaker
	Valid          bool
	CreateTime     time.Time
	ModifyTime     time.Time
}

/**
 * @brief 熔断规则绑定关系
 */
type CircuitBreakerRelation struct {
	ServiceID   string
	RuleID      string
	RuleVersion string
	Valid       bool
	CreateTime  time.Time
	ModifyTime  time.Time
}

/**
 * @brief 返回给控制台的熔断规则及服务数据
 */
type CircuitBreakerDetail struct {
	Total               uint32
	CircuitBreakerInfos []*CircuitBreakerInfo
}

/**
 * @brief 熔断规则及绑定服务
 */
type CircuitBreakerInfo struct {
	CircuitBreaker *CircuitBreaker
	Services       []*Service
}

/**
 * @brief 平台信息
 */
type Platform struct {
	ID         string
	Name       string
	Domain     string
	QPS        uint32
	Token      string
	Owner      string
	Department string
	Comment    string
	Valid      bool
	CreateTime time.Time
	ModifyTime time.Time
}

// 整数转换为bool值
func int2bool(entry int) bool {
	if entry == 0 {
		return false
	}
	return true
}

// store的flag转换为valid
// flag==1为无效，其他情况为有效
func flag2valid(flag int) bool {
	if flag == 1 {
		return false
	}
	return true

}

// int64的时间戳转为字符串时间
func int64Time2String(t int64) string {
	return time.Unix(t, 0).Format("2006-01-02 15:04:05")
}

// 操作类型
type OperationType string

// 定义包含的操作类型
const (
	// 新建
	OCreate OperationType = "Create"

	// 删除
	ODelete OperationType = "Delete"

	// 更新
	OUpdate OperationType = "Update"

	// 更新隔离状态
	OUpdateIsolate OperationType = "UpdateIsolate"

	// 查看token
	OGetToken OperationType = "GetToken" // nolint

	// 更新token
	OUpdateToken OperationType = "UpdateToken" // nolint
)

// 操作资源
type Resource string

// 定义包含的资源类型
const (
	RNamespace     Resource = "Namespace"
	RService       Resource = "Service"
	RRouting       Resource = "Routing"
	RInstance      Resource = "Instance"
	RRateLimit     Resource = "RateLimit"
	RMeshResource  Resource = "MeshResource"
	RMesh          Resource = "Mesh"
	RMeshService   Resource = "MeshService"
	RFluxRateLimit Resource = "FluxRateLimit"
)

// 资源类型
type ResourceType int

const (
	// 网格类型资源
	MeshType ResourceType = 0
	// 北极星服务类型资源
	ServiceType ResourceType = 1
)

var ResourceTypeMap = map[Resource]ResourceType{
	RNamespace:    ServiceType,
	RService:      ServiceType,
	RRouting:      ServiceType,
	RInstance:     ServiceType,
	RRateLimit:    ServiceType,
	RMesh:         MeshType,
	RMeshResource: MeshType,
	RMeshService:  MeshType,
}

// 获取资源的大类型
func GetResourceType(r Resource) ResourceType {
	return ResourceTypeMap[r]
}

// 操作记录entry
type RecordEntry struct {
	ResourceType  Resource
	OperationType OperationType
	Namespace     string
	Service       string
	MeshID        string
	MeshName      string
	Context       string
	Operator      string
	Revision      string
	CreateTime    time.Time
}
