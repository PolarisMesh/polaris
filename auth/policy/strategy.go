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

package policy

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/wrappers"
	apimodel "github.com/polarismesh/specification/source/go/api/v1/model"
	apisecurity "github.com/polarismesh/specification/source/go/api/v1/security"
	apiservice "github.com/polarismesh/specification/source/go/api/v1/service_manage"
	"go.uber.org/zap"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/wrapperspb"

	cachetypes "github.com/polarismesh/polaris/cache/api"
	api "github.com/polarismesh/polaris/common/api/v1"
	"github.com/polarismesh/polaris/common/model"
	authcommon "github.com/polarismesh/polaris/common/model/auth"
	commonstore "github.com/polarismesh/polaris/common/store"
	commontime "github.com/polarismesh/polaris/common/time"
	"github.com/polarismesh/polaris/common/utils"
)

type (
	// StrategyDetail2Api strategy detail to *apisecurity.AuthStrategy func
	StrategyDetail2Api func(ctx context.Context, user *authcommon.StrategyDetail) *apisecurity.AuthStrategy
)

// CreateStrategy 创建鉴权策略
func (svr *Server) CreateStrategy(ctx context.Context, req *apisecurity.AuthStrategy) *apiservice.Response {
	req.Owner = utils.NewStringValue(utils.ParseOwnerID(ctx))
	if checkErrResp := svr.checkCreateStrategy(req); checkErrResp != nil {
		return checkErrResp
	}

	req.Resources = svr.normalizeResource(req.Resources)

	data := svr.createAuthStrategyModel(req)

	tx, err := svr.storage.StartTx()
	if err != nil {
		log.Error("[Auth][Strategy] start tx", utils.RequestID(ctx), zap.Error(err))
		return api.NewAuthResponse(commonstore.StoreCode2APICode(err))
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := svr.storage.AddStrategy(tx, data); err != nil {
		log.Error("[Auth][Strategy] create strategy into store", utils.RequestID(ctx),
			zap.Error(err))
		return api.NewAuthResponse(commonstore.StoreCode2APICode(err))
	}
	if err := tx.Commit(); err != nil {
		log.Error("[Auth][Strategy] create strategy  commit tx", utils.RequestID(ctx), zap.Error(err))
		return api.NewAuthResponse(commonstore.StoreCode2APICode(err))
	}

	log.Info("[Auth][Strategy] create strategy", utils.RequestID(ctx), zap.String("name", req.Name.GetValue()))
	svr.RecordHistory(authStrategyRecordEntry(ctx, req, data, model.OCreate))

	return api.NewAuthStrategyResponse(apimodel.Code_ExecuteSuccess, req)
}

// UpdateStrategies 批量修改鉴权
func (svr *Server) UpdateStrategies(
	ctx context.Context, reqs []*apisecurity.ModifyAuthStrategy) *apiservice.BatchWriteResponse {
	resp := api.NewAuthBatchWriteResponse(apimodel.Code_ExecuteSuccess)

	for index := range reqs {
		ret := svr.UpdateStrategy(ctx, reqs[index])
		api.Collect(resp, ret)
	}

	return resp
}

// UpdateStrategy 实现鉴权策略的变更
// Case 1. 修改的是默认鉴权策略的话，只能修改资源，不能添加、删除用户 or 用户组
// Case 2. 鉴权策略只能被自己的 owner 对应的用户修改
// Case 3. 主账户的默认策略不得修改
func (svr *Server) UpdateStrategy(ctx context.Context, req *apisecurity.ModifyAuthStrategy) *apiservice.Response {
	strategy, err := svr.storage.GetStrategyDetail(req.GetId().GetValue())
	if err != nil {
		log.Error("[Auth][Strategy] get strategy from store", utils.RequestID(ctx),
			zap.Error(err))
		return api.NewModifyAuthStrategyResponse(commonstore.StoreCode2APICode(err), req)
	}
	if strategy == nil {
		return api.NewModifyAuthStrategyResponse(apimodel.Code_NotFoundAuthStrategyRule, req)
	}

	if checkErrResp := svr.checkUpdateStrategy(ctx, req, strategy); checkErrResp != nil {
		return checkErrResp
	}

	req.AddResources = svr.normalizeResource(req.AddResources)
	data, needUpdate := svr.updateAuthStrategyAttribute(ctx, req, strategy)
	if !needUpdate {
		return api.NewModifyAuthStrategyResponse(apimodel.Code_NoNeedUpdate, req)
	}

	if err := svr.storage.UpdateStrategy(data); err != nil {
		log.Error("[Auth][Strategy] update strategy into store",
			utils.RequestID(ctx), zap.Error(err))
		return api.NewAuthResponseWithMsg(commonstore.StoreCode2APICode(err), err.Error())
	}

	log.Info("[Auth][Strategy] update strategy into store", utils.RequestID(ctx),
		zap.String("name", strategy.Name))
	svr.RecordHistory(authModifyStrategyRecordEntry(ctx, req, data, model.OUpdate))

	return api.NewModifyAuthStrategyResponse(apimodel.Code_ExecuteSuccess, req)
}

// DeleteStrategies 批量删除鉴权策略
func (svr *Server) DeleteStrategies(
	ctx context.Context, reqs []*apisecurity.AuthStrategy) *apiservice.BatchWriteResponse {
	resp := api.NewAuthBatchWriteResponse(apimodel.Code_ExecuteSuccess)
	for index := range reqs {
		ret := svr.DeleteStrategy(ctx, reqs[index])
		api.Collect(resp, ret)
	}

	return resp
}

// DeleteStrategy 删除鉴权策略
// Case 1. 只有该策略的 owner 账户可以删除策略
// Case 2. 默认策略不能被删除，默认策略只能随着账户的删除而被清理
func (svr *Server) DeleteStrategy(ctx context.Context, req *apisecurity.AuthStrategy) *apiservice.Response {
	strategy, err := svr.storage.GetStrategyDetail(req.GetId().GetValue())
	if err != nil {
		log.Error("[Auth][Strategy] get strategy from store", utils.RequestID(ctx),
			zap.Error(err))
		return api.NewAuthStrategyResponse(commonstore.StoreCode2APICode(err), req)
	}

	if strategy == nil {
		return api.NewAuthStrategyResponse(apimodel.Code_ExecuteSuccess, req)
	}

	if strategy.Default {
		log.Error("[Auth][Strategy] delete default strategy is denied", utils.RequestID(ctx))
		return api.NewAuthStrategyResponseWithMsg(apimodel.Code_BadRequest, "default strategy can't delete", req)
	}

	if strategy.Owner != utils.ParseUserID(ctx) {
		return api.NewAuthStrategyResponse(apimodel.Code_NotAllowedAccess, req)
	}

	if err := svr.storage.DeleteStrategy(req.GetId().GetValue()); err != nil {
		log.Error("[Auth][Strategy] delete strategy from store",
			utils.RequestID(ctx), zap.Error(err))
		return api.NewAuthResponse(commonstore.StoreCode2APICode(err))
	}

	log.Info("[Auth][Strategy] delete strategy from store", utils.RequestID(ctx),
		zap.String("name", req.Name.GetValue()))
	svr.RecordHistory(authStrategyRecordEntry(ctx, req, strategy, model.ODelete))

	return api.NewAuthStrategyResponse(apimodel.Code_ExecuteSuccess, req)
}

// GetStrategies 批量查询鉴权策略
// Case 1. 如果是以资源视角来查询鉴权策略，那么就会忽略自动根据账户类型进行数据查看的限制
//
//	eg. 比如当前子账户A想要查看资源R的相关的策略，那么不在会自动注入 principal_id 以及 principal_type 的查询条件
//
// Case 2. 如果是以用户视角来查询鉴权策略，如果没有带上 principal_id，那么就会根据账户类型自动注入 principal_id 以
//
//	及 principal_type 的查询条件，从而限制该账户的数据查看
//	eg.
//		a. 如果当前是超级管理账户，则按照传入的 query 进行查询即可
//		b. 如果当前是主账户，则自动注入 owner 字段，即只能查看策略的 owner 是自己的策略
//		c. 如果当前是子账户，则自动注入 principal_id 以及 principal_type 字段，即稚嫩查询与自己有关的策略
func (svr *Server) GetStrategies(ctx context.Context, filters map[string]string) *apiservice.BatchQueryResponse {
	filters = ParseStrategySearchArgs(ctx, filters)
	offset, limit, _ := utils.ParseOffsetAndLimit(filters)
	total, strategies, err := svr.cacheMgr.AuthStrategy().Query(ctx, cachetypes.PolicySearchArgs{
		Filters: filters,
		Offset:  offset,
		Limit:   limit,
	})
	if err != nil {
		log.Error("[Auth][Strategy] get strategies from store", zap.Any("query", filters),
			utils.RequestID(ctx), zap.Error(err))
		return api.NewAuthBatchQueryResponse(commonstore.StoreCode2APICode(err))
	}

	resp := api.NewAuthBatchQueryResponse(apimodel.Code_ExecuteSuccess)
	resp.Amount = utils.NewUInt32Value(total)
	resp.Size = utils.NewUInt32Value(uint32(len(strategies)))

	if strings.Compare(filters["show_detail"], "true") == 0 {
		log.Info("[Auth][Strategy] fill strategy detail", utils.RequestID(ctx))
		resp.AuthStrategies = enhancedAuthStrategy2Api(ctx, strategies, svr.authStrategyFull2Api)
	} else {
		resp.AuthStrategies = enhancedAuthStrategy2Api(ctx, strategies, svr.authStrategy2Api)
	}

	return resp
}

var (
	resTypeFilter = map[string]string{
		"namespace":    "0",
		"service":      "1",
		"config_group": "2",
	}

	principalTypeFilter = map[string]string{
		"user":   "1",
		"group":  "2",
		"groups": "2",
	}
)

// ParseStrategySearchArgs 处理鉴权策略的搜索参数
func ParseStrategySearchArgs(ctx context.Context, searchFilters map[string]string) map[string]string {
	if val, ok := searchFilters["res_type"]; ok {
		if v, exist := resTypeFilter[val]; exist {
			searchFilters["res_type"] = v
		} else {
			searchFilters["res_type"] = "0"
		}
	}

	if val, ok := searchFilters["principal_type"]; ok {
		if v, exist := principalTypeFilter[val]; exist {
			searchFilters["principal_type"] = v
		} else {
			searchFilters["principal_type"] = "1"
		}
	}

	if authcommon.ParseUserRole(ctx) != authcommon.AdminUserRole {
		// 如果当前账户不是 admin 角色，既不是走资源视角查看，也不是指定principal查看，那么只能查询当前操作用户被关联到的鉴权策略，
		if _, ok := searchFilters["res_id"]; !ok {
			// 设置 owner 参数，只能查看对应 owner 下的策略
			searchFilters["owner"] = utils.ParseOwnerID(ctx)
			if _, ok := searchFilters["principal_id"]; !ok {
				// 如果当前不是 owner 角色，那么只能查询与自己有关的策略
				if !utils.ParseIsOwner(ctx) {
					searchFilters["principal_id"] = utils.ParseUserID(ctx)
					searchFilters["principal_type"] = strconv.Itoa(int(authcommon.PrincipalUser))
				}
			}
		}
	}

	return searchFilters
}

// GetStrategy 根据策略ID获取详细的鉴权策略
// Case 1 如果当前操作者是该策略 principal 中的一员，则可以查看
// Case 2 如果当前操作者是该策略的 owner，则可以查看
// Case 3 如果当前操作者是admin角色，直接查看
func (svr *Server) GetStrategy(ctx context.Context, req *apisecurity.AuthStrategy) *apiservice.Response {
	userId := utils.ParseUserID(ctx)
	isOwner := utils.ParseIsOwner(ctx)

	if req.GetId().GetValue() == "" {
		return api.NewAuthResponse(apimodel.Code_EmptyQueryParameter)
	}

	ret, err := svr.storage.GetStrategyDetail(req.GetId().GetValue())
	if err != nil {
		log.Error("[Auth][Strategy] get strategt from store",
			utils.RequestID(ctx), zap.Error(err))
		return api.NewAuthResponse(commonstore.StoreCode2APICode(err))
	}
	if ret == nil {
		return api.NewAuthStrategyResponse(apimodel.Code_NotFoundAuthStrategyRule, req)
	}

	var canView bool
	if isOwner {
		// 是否是本鉴权策略的 owner 账户, 或者是否是超级管理员, 是的话则快速跳过下面的检查
		canView = (ret.Owner == userId) || authcommon.ParseUserRole(ctx) == authcommon.AdminUserRole
	}

	// 判断是否在该策略所属的成员列表中，如果自己在某个用户组，而该用户组又在这个策略的成员中，则也是可以查看的
	if !canView {
		curUser := &apisecurity.User{
			Id: wrapperspb.String(userId),
		}
		for index := range ret.Principals {
			principal := ret.Principals[index]
			if principal.PrincipalType == authcommon.PrincipalUser && principal.PrincipalID == userId {
				canView = true
				break
			}
			if principal.PrincipalType == authcommon.PrincipalGroup {
				group := &apisecurity.UserGroup{
					Id: wrapperspb.String(principal.PrincipalID),
				}
				if svr.userSvr.GetUserHelper().CheckUserInGroup(ctx, group, curUser) {
					canView = true
					break
				}
			}
		}
	}

	if !canView {
		log.Error("[Auth][Strategy] get strategy detail denied",
			utils.RequestID(ctx), zap.String("user", userId), zap.String("strategy", req.Id.Value),
			zap.Bool("is-owner", isOwner),
		)
		return api.NewAuthStrategyResponse(apimodel.Code_NotAllowedAccess, req)
	}

	return api.NewAuthStrategyResponse(apimodel.Code_ExecuteSuccess, svr.authStrategyFull2Api(ctx, ret))
}

// GetPrincipalResources 获取某个principal可以获取到的所有资源ID数据信息
func (svr *Server) GetPrincipalResources(ctx context.Context, query map[string]string) *apiservice.Response {
	if len(query) == 0 {
		return api.NewAuthResponse(apimodel.Code_EmptyRequest)
	}

	principalId := query["principal_id"]
	if principalId == "" {
		return api.NewAuthResponse(apimodel.Code_EmptyQueryParameter)
	}

	var principalType string
	if v, exist := principalTypeFilter[query["principal_type"]]; exist {
		principalType = v
	} else {
		principalType = "1"
	}

	principalRole, _ := strconv.ParseInt(principalType, 10, 64)
	if err := authcommon.CheckPrincipalType(int(principalRole)); err != nil {
		return api.NewAuthResponse(apimodel.Code_InvalidPrincipalType)
	}

	var (
		resources = make([]authcommon.StrategyResource, 0, 20)
		err       error
	)

	// 找这个用户所关联的用户组
	if authcommon.PrincipalType(principalRole) == authcommon.PrincipalUser {
		groups := svr.userSvr.GetUserHelper().GetUserOwnGroup(ctx, &apisecurity.User{
			Id: wrapperspb.String(principalId),
		})
		for i := range groups {
			item := groups[i]
			res, err := svr.storage.GetStrategyResources(item.GetId().GetValue(), authcommon.PrincipalGroup)
			if err != nil {
				log.Error("[Auth][Strategy] get principal link resource", utils.RequestID(ctx),
					zap.String("principal-id", principalId), zap.Any("principal-role", principalRole), zap.Error(err))
				return api.NewAuthResponse(commonstore.StoreCode2APICode(err))
			}
			resources = append(resources, res...)
		}
	}

	pResources, err := svr.storage.GetStrategyResources(principalId, authcommon.PrincipalType(principalRole))
	if err != nil {
		log.Error("[Auth][Strategy] get principal link resource", utils.RequestID(ctx),
			zap.String("principal-id", principalId), zap.Any("principal-role", principalRole), zap.Error(err))
		return api.NewAuthResponse(commonstore.StoreCode2APICode(err))
	}

	resources = append(resources, pResources...)
	tmp := &apisecurity.AuthStrategy{
		Resources: &apisecurity.StrategyResources{
			Namespaces:   make([]*apisecurity.StrategyResourceEntry, 0),
			Services:     make([]*apisecurity.StrategyResourceEntry, 0),
			ConfigGroups: make([]*apisecurity.StrategyResourceEntry, 0),
		},
	}

	svr.enrichResourceInfo(ctx, tmp, &authcommon.StrategyDetail{
		Resources: resourceDeduplication(resources),
	})

	return api.NewStrategyResourcesResponse(apimodel.Code_ExecuteSuccess, tmp.Resources)
}

// enhancedAuthStrategy2Api
func enhancedAuthStrategy2Api(ctx context.Context, s []*authcommon.StrategyDetail,
	fn StrategyDetail2Api) []*apisecurity.AuthStrategy {
	out := make([]*apisecurity.AuthStrategy, 0, len(s))
	for k := range s {
		out = append(out, fn(ctx, s[k]))
	}
	return out
}

// authStrategy2Api
func (svr *Server) authStrategy2Api(ctx context.Context, s *authcommon.StrategyDetail) *apisecurity.AuthStrategy {
	if s == nil {
		return nil
	}

	// note: 不包括token，token比较特殊
	out := &apisecurity.AuthStrategy{
		Id:              utils.NewStringValue(s.ID),
		Name:            utils.NewStringValue(s.Name),
		Owner:           utils.NewStringValue(s.Owner),
		Comment:         utils.NewStringValue(s.Comment),
		Ctime:           utils.NewStringValue(commontime.Time2String(s.CreateTime)),
		Mtime:           utils.NewStringValue(commontime.Time2String(s.ModifyTime)),
		Action:          apisecurity.AuthAction(apisecurity.AuthAction_value[s.Action]),
		DefaultStrategy: utils.NewBoolValue(s.Default),
	}

	return out
}

// authStrategyFull2Api
func (svr *Server) authStrategyFull2Api(ctx context.Context, data *authcommon.StrategyDetail) *apisecurity.AuthStrategy {
	if data == nil {
		return nil
	}

	users := make([]*wrappers.StringValue, 0, len(data.Principals))
	groups := make([]*wrappers.StringValue, 0, len(data.Principals))
	for index := range data.Principals {
		principal := data.Principals[index]
		if principal.PrincipalType == authcommon.PrincipalUser {
			users = append(users, utils.NewStringValue(principal.PrincipalID))
		} else {
			groups = append(groups, utils.NewStringValue(principal.PrincipalID))
		}
	}

	// note: 不包括token，token比较特殊
	out := &apisecurity.AuthStrategy{
		Id:              utils.NewStringValue(data.ID),
		Name:            utils.NewStringValue(data.Name),
		Owner:           utils.NewStringValue(data.Owner),
		Comment:         utils.NewStringValue(data.Comment),
		Ctime:           utils.NewStringValue(commontime.Time2String(data.CreateTime)),
		Mtime:           utils.NewStringValue(commontime.Time2String(data.ModifyTime)),
		Action:          apisecurity.AuthAction(apisecurity.AuthAction_value[data.Action]),
		DefaultStrategy: utils.NewBoolValue(data.Default),
		Functions:       data.CalleeMethods,
		Metadata:        data.Metadata,
	}

	svr.enrichPrincipalInfo(out, data)
	svr.enrichResourceInfo(ctx, out, data)
	return out
}

var (
	resourceFieldNames = map[string]apisecurity.ResourceType{
		"namespaces":           apisecurity.ResourceType_Namespaces,
		"service":              apisecurity.ResourceType_Services,
		"config_groups":        apisecurity.ResourceType_ConfigGroups,
		"route_rules":          apisecurity.ResourceType_RouteRules,
		"ratelimit_rules":      apisecurity.ResourceType_RateLimitRules,
		"circuitbreaker_rules": apisecurity.ResourceType_CircuitBreakerRules,
		"faultdetect_rules":    apisecurity.ResourceType_FaultDetectRules,
		"lane_rules":           apisecurity.ResourceType_LaneRules,
		"users":                apisecurity.ResourceType_Users,
		"user_groups":          apisecurity.ResourceType_UserGroups,
		"roles":                apisecurity.ResourceType_Roles,
		"auth_policies":        apisecurity.ResourceType_PolicyRules,
	}

	resourceToFieldNames = map[apisecurity.ResourceType]string{
		apisecurity.ResourceType_Namespaces:          "namespaces",
		apisecurity.ResourceType_Services:            "service",
		apisecurity.ResourceType_ConfigGroups:        "config_groups",
		apisecurity.ResourceType_RouteRules:          "route_rules",
		apisecurity.ResourceType_RateLimitRules:      "ratelimit_rules",
		apisecurity.ResourceType_CircuitBreakerRules: "circuitbreaker_rules",
		apisecurity.ResourceType_FaultDetectRules:    "faultdetect_rules",
		apisecurity.ResourceType_LaneRules:           "lane_rules",
		apisecurity.ResourceType_Users:               "users",
		apisecurity.ResourceType_UserGroups:          "user_groups",
		apisecurity.ResourceType_Roles:               "roles",
		apisecurity.ResourceType_PolicyRules:         "auth_policies",
	}
)

// createAuthStrategyModel 创建鉴权策略的存储模型
func (svr *Server) createAuthStrategyModel(strategy *apisecurity.AuthStrategy) *authcommon.StrategyDetail {
	ret := &authcommon.StrategyDetail{}
	ret.FromSpec(strategy)

	// 收集涉及的资源信息
	resEntry := make([]authcommon.StrategyResource, 0, 20)

	message := strategy.Resources.ProtoReflect()
	fields := strategy.Resources.ProtoReflect().Descriptor().Fields()
	for name, resType := range resourceFieldNames {
		field := fields.ByName(protoreflect.Name(resourceFieldNames[name]))
		resEntry = append(resEntry, svr.collectResourceEntry(ret.ID, resType, message.Get(field).List(), false)...)
	}

	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_Namespaces,
	// 	strategy.GetResources().GetNamespaces(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_Services,
	// 	strategy.GetResources().GetServices(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_ConfigGroups,
	// 	strategy.GetResources().GetConfigGroups(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_RateLimitRules,
	// 	strategy.GetResources().GetRatelimitRules(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_RouteRules,
	// 	strategy.GetResources().GetRouteRules(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_LaneRules,
	// 	strategy.GetResources().GetLaneRules(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_CircuitBreakerRules,
	// 	strategy.GetResources().GetCircuitbreakerRules(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_FaultDetectRules,
	// 	strategy.GetResources().GetFaultdetectRules(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_Users,
	// 	strategy.GetResources().GetUsers(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_UserGroups,
	// 	strategy.GetResources().GetUserGroups(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_Roles,
	// 	strategy.GetResources().GetRoles(), false)...)
	// resEntry = append(resEntry, svr.collectResEntry(ret.ID, apisecurity.ResourceType_PolicyRules,
	// 	strategy.GetResources().GetAuthPolicies(), false)...)

	// 收集涉及的 principal 信息
	principals := make([]authcommon.Principal, 0, 20)
	principals = append(principals, collectPrincipalEntry(ret.ID, authcommon.PrincipalUser,
		strategy.GetPrincipals().GetUsers())...)
	principals = append(principals, collectPrincipalEntry(ret.ID, authcommon.PrincipalGroup,
		strategy.GetPrincipals().GetGroups())...)
	principals = append(principals, collectPrincipalEntry(ret.ID, authcommon.PrincipalRole,
		strategy.GetPrincipals().GetRoles())...)

	ret.Resources = resEntry
	ret.Principals = principals

	return ret
}

// updateAuthStrategyAttribute 更新计算鉴权策略的属性
func (svr *Server) updateAuthStrategyAttribute(ctx context.Context, strategy *apisecurity.ModifyAuthStrategy,
	saved *authcommon.StrategyDetail) (*authcommon.ModifyStrategyDetail, bool) {
	var needUpdate bool
	ret := &authcommon.ModifyStrategyDetail{
		ID:         strategy.Id.GetValue(),
		Name:       saved.Name,
		Action:     saved.Action,
		Comment:    saved.Comment,
		ModifyTime: time.Now(),
	}

	// 只有 owner 可以修改的属性
	if utils.ParseIsOwner(ctx) {
		if strategy.GetComment() != nil && strategy.GetComment().GetValue() != saved.Comment {
			needUpdate = true
			ret.Comment = strategy.GetComment().GetValue()
		}

		if strategy.GetName().GetValue() != "" && strategy.GetName().GetValue() != saved.Name {
			needUpdate = true
			ret.Name = strategy.GetName().GetValue()
		}
	}

	if svr.computeResourceChange(ret, strategy) {
		needUpdate = true
	}
	if computePrincipalChange(ret, strategy) {
		needUpdate = true
	}
	needUpdate = true
	saved.CalleeMethods = strategy.Functions
	saved.Metadata = strategy.Metadata
	saved.Conditions = func() []authcommon.Condition {
		conditions := make([]authcommon.Condition, 0, len(strategy.GetResourceLabels()))
		for index := range strategy.GetResourceLabels() {
			conditions = append(conditions, authcommon.Condition{
				Key:   strategy.GetResourceLabels()[index].GetKey(),
				Value: strategy.GetResourceLabels()[index].GetValue(),
			})
		}
		return conditions
	}()

	return ret, needUpdate
}

// computeResourceChange 计算资源的变化情况，判断是否涉及变更
func (svr *Server) computeResourceChange(
	modify *authcommon.ModifyStrategyDetail, strategy *apisecurity.ModifyAuthStrategy) bool {
	var needUpdate bool
	addResEntry := make([]authcommon.StrategyResource, 0)
	addResEntry = append(addResEntry, svr.collectResEntry(modify.ID, apisecurity.ResourceType_Namespaces,
		strategy.GetAddResources().GetNamespaces(), false)...)
	addResEntry = append(addResEntry, svr.collectResEntry(modify.ID, apisecurity.ResourceType_Services,
		strategy.GetAddResources().GetServices(), false)...)
	addResEntry = append(addResEntry, svr.collectResEntry(modify.ID, apisecurity.ResourceType_ConfigGroups,
		strategy.GetAddResources().GetConfigGroups(), false)...)

	if len(addResEntry) != 0 {
		needUpdate = true
		modify.AddResources = addResEntry
	}

	removeResEntry := make([]authcommon.StrategyResource, 0)
	removeResEntry = append(removeResEntry, svr.collectResEntry(modify.ID, apisecurity.ResourceType_Namespaces,
		strategy.GetRemoveResources().GetNamespaces(), true)...)
	removeResEntry = append(removeResEntry, svr.collectResEntry(modify.ID, apisecurity.ResourceType_Services,
		strategy.GetRemoveResources().GetServices(), true)...)
	removeResEntry = append(removeResEntry, svr.collectResEntry(modify.ID, apisecurity.ResourceType_ConfigGroups,
		strategy.GetRemoveResources().GetConfigGroups(), true)...)

	if len(removeResEntry) != 0 {
		needUpdate = true
		modify.RemoveResources = removeResEntry
	}

	return needUpdate
}

// computePrincipalChange 计算 principal 的变化情况，判断是否涉及变更
func computePrincipalChange(modify *authcommon.ModifyStrategyDetail, strategy *apisecurity.ModifyAuthStrategy) bool {
	var needUpdate bool
	addPrincipals := make([]authcommon.Principal, 0)
	addPrincipals = append(addPrincipals, collectPrincipalEntry(modify.ID, authcommon.PrincipalUser,
		strategy.GetAddPrincipals().GetUsers())...)
	addPrincipals = append(addPrincipals, collectPrincipalEntry(modify.ID, authcommon.PrincipalGroup,
		strategy.GetAddPrincipals().GetGroups())...)
	addPrincipals = append(addPrincipals, collectPrincipalEntry(modify.ID, authcommon.PrincipalRole,
		strategy.GetAddPrincipals().GetRoles())...)

	if len(addPrincipals) != 0 {
		needUpdate = true
		modify.AddPrincipals = addPrincipals
	}

	removePrincipals := make([]authcommon.Principal, 0)
	removePrincipals = append(removePrincipals, collectPrincipalEntry(modify.ID, authcommon.PrincipalUser,
		strategy.GetRemovePrincipals().GetUsers())...)
	removePrincipals = append(removePrincipals, collectPrincipalEntry(modify.ID, authcommon.PrincipalGroup,
		strategy.GetRemovePrincipals().GetGroups())...)
	removePrincipals = append(removePrincipals, collectPrincipalEntry(modify.ID, authcommon.PrincipalRole,
		strategy.GetRemovePrincipals().GetRoles())...)

	if len(removePrincipals) != 0 {
		needUpdate = true
		modify.RemovePrincipals = removePrincipals
	}

	return needUpdate
}

type pbStringValue interface {
	GetValue() string
}

// collectResEntry 将资源ID转换为对应的 []authcommon.StrategyResource 数组
func (svr *Server) collectResourceEntry(ruleId string, resType apisecurity.ResourceType,
	res protoreflect.List, delete bool) []authcommon.StrategyResource {
	resEntries := make([]authcommon.StrategyResource, 0, res.Len())
	if res.Len() == 0 {
		return resEntries
	}

	for i := 0; i < res.Len(); i++ {
		resource := res.Get(i).Message()
		fields := resource.Descriptor().Fields()
		resId := resource.Get(fields.ByName("id")).Interface().(pbStringValue)
		resName := resource.Get(fields.ByName("name")).Interface().(pbStringValue)
		// 如果是添加的动作，那么需要进行归一化处理
		if !delete {
			// 归一化处理
			if resId.GetValue() == "*" || resName.GetValue() == "*" {
				return []authcommon.StrategyResource{
					{
						StrategyID: ruleId,
						ResType:    int32(resType),
						ResID:      "*",
					},
				}
			}
		}

		entry := authcommon.StrategyResource{
			StrategyID: ruleId,
			ResType:    int32(resType),
			ResID:      resId.GetValue(),
		}

		resEntries = append(resEntries, entry)
	}

	return resEntries
}

// collectResEntry 将资源ID转换为对应的 []authcommon.StrategyResource 数组
func (svr *Server) collectResEntry(ruleId string, resType apisecurity.ResourceType,
	res []*apisecurity.StrategyResourceEntry, delete bool) []authcommon.StrategyResource {
	resEntries := make([]authcommon.StrategyResource, 0, len(res)+1)
	if len(res) == 0 {
		return resEntries
	}

	for index := range res {
		// 如果是添加的动作，那么需要进行归一化处理
		if !delete {
			// 归一化处理
			if res[index].GetId().GetValue() == "*" || res[index].GetName().GetValue() == "*" {
				return []authcommon.StrategyResource{
					{
						StrategyID: ruleId,
						ResType:    int32(resType),
						ResID:      "*",
					},
				}
			}
		}

		entry := authcommon.StrategyResource{
			StrategyID: ruleId,
			ResType:    int32(resType),
			ResID:      res[index].GetId().GetValue(),
		}

		resEntries = append(resEntries, entry)
	}

	return resEntries
}

// collectPrincipalEntry 将 Principal 转换为对应的 []authcommon.Principal 数组
func collectPrincipalEntry(ruleID string, uType authcommon.PrincipalType, res []*apisecurity.Principal) []authcommon.Principal {
	principals := make([]authcommon.Principal, 0, len(res)+1)
	if len(res) == 0 {
		return principals
	}

	for index := range res {
		principals = append(principals, authcommon.Principal{
			StrategyID:    ruleID,
			PrincipalID:   res[index].GetId().GetValue(),
			PrincipalType: uType,
		})
	}

	return principals
}

// checkCreateStrategy 检查创建鉴权策略的请求
func (svr *Server) checkCreateStrategy(req *apisecurity.AuthStrategy) *apiservice.Response {
	// 检查名称信息
	if err := CheckName(req.GetName()); err != nil {
		return api.NewAuthStrategyResponse(apimodel.Code_InvalidUserName, req)
	}
	// 检查用户是否存在
	if err := svr.checkUserExist(convertPrincipalsToUsers(req.GetPrincipals())); err != nil {
		return api.NewAuthStrategyResponse(apimodel.Code_NotFoundUser, req)
	}
	// 检查用户组是否存在
	if err := svr.checkGroupExist(convertPrincipalsToGroups(req.GetPrincipals())); err != nil {
		return api.NewAuthStrategyResponse(apimodel.Code_NotFoundUserGroup, req)
	}
	// 检查资源是否存在
	if errResp := svr.checkResourceExist(req.GetResources()); errResp != nil {
		return errResp
	}
	return nil
}

// checkUpdateStrategy 检查更新鉴权策略的请求
// Case 1. 修改的是默认鉴权策略的话，只能修改资源，不能添加用户 or 用户组
// Case 2. 鉴权策略只能被自己的 owner 对应的用户修改
func (svr *Server) checkUpdateStrategy(ctx context.Context, req *apisecurity.ModifyAuthStrategy,
	saved *authcommon.StrategyDetail) *apiservice.Response {
	userId := utils.ParseUserID(ctx)
	if authcommon.ParseUserRole(ctx) != authcommon.AdminUserRole {
		if !utils.ParseIsOwner(ctx) || userId != saved.Owner {
			log.Error("[Auth][Strategy] modify strategy denied, current user not owner",
				utils.RequestID(ctx), zap.String("user", userId), zap.String("owner", saved.Owner),
				zap.String("strategy", saved.ID))
			return api.NewModifyAuthStrategyResponse(apimodel.Code_NotAllowedAccess, req)
		}
	}

	if saved.Default {
		if len(req.AddPrincipals.Users) != 0 ||
			len(req.AddPrincipals.Groups) != 0 ||
			len(req.RemovePrincipals.Groups) != 0 ||
			len(req.RemovePrincipals.Users) != 0 {
			return api.NewModifyAuthStrategyResponse(apimodel.Code_NotAllowModifyDefaultStrategyPrincipal, req)
		}

		// 主账户的默认策略禁止编辑
		if len(saved.Principals) == 1 && saved.Principals[0].PrincipalType == authcommon.PrincipalUser {
			if saved.Principals[0].PrincipalID == utils.ParseOwnerID(ctx) {
				return api.NewAuthResponse(apimodel.Code_NotAllowModifyOwnerDefaultStrategy)
			}
		}
	}

	// 检查用户是否存在
	if err := svr.checkUserExist(convertPrincipalsToUsers(req.GetAddPrincipals())); err != nil {
		return api.NewModifyAuthStrategyResponse(apimodel.Code_NotFoundUser, req)
	}

	// 检查用户组是否存
	if err := svr.checkGroupExist(convertPrincipalsToGroups(req.GetAddPrincipals())); err != nil {
		return api.NewModifyAuthStrategyResponse(apimodel.Code_NotFoundUserGroup, req)
	}

	// 检查资源是否存在
	if errResp := svr.checkResourceExist(req.GetAddResources()); errResp != nil {
		return errResp
	}

	return nil
}

// authStrategyRecordEntry 转换为鉴权策略的记录结构体
func authStrategyRecordEntry(ctx context.Context, req *apisecurity.AuthStrategy, md *authcommon.StrategyDetail,
	operationType model.OperationType) *model.RecordEntry {

	marshaler := jsonpb.Marshaler{}
	detail, _ := marshaler.MarshalToString(req)

	entry := &model.RecordEntry{
		ResourceType:  model.RAuthStrategy,
		ResourceName:  fmt.Sprintf("%s(%s)", md.Name, md.ID),
		OperationType: operationType,
		Operator:      utils.ParseOperator(ctx),
		Detail:        detail,
		HappenTime:    time.Now(),
	}

	return entry
}

// authModifyStrategyRecordEntry
func authModifyStrategyRecordEntry(
	ctx context.Context, req *apisecurity.ModifyAuthStrategy, md *authcommon.ModifyStrategyDetail,
	operationType model.OperationType) *model.RecordEntry {

	marshaler := jsonpb.Marshaler{}
	detail, _ := marshaler.MarshalToString(req)

	entry := &model.RecordEntry{
		ResourceType:  model.RAuthStrategy,
		ResourceName:  fmt.Sprintf("%s(%s)", md.Name, md.ID),
		OperationType: operationType,
		Operator:      utils.ParseOperator(ctx),
		Detail:        detail,
		HappenTime:    time.Now(),
	}

	return entry
}

func convertPrincipalsToUsers(principals *apisecurity.Principals) []*apisecurity.User {
	if principals == nil {
		return make([]*apisecurity.User, 0)
	}

	users := make([]*apisecurity.User, 0, len(principals.Users))
	for k := range principals.GetUsers() {
		user := principals.GetUsers()[k]
		users = append(users, &apisecurity.User{
			Id: user.Id,
		})
	}

	return users
}

func convertPrincipalsToGroups(principals *apisecurity.Principals) []*apisecurity.UserGroup {
	if principals == nil {
		return make([]*apisecurity.UserGroup, 0)
	}

	groups := make([]*apisecurity.UserGroup, 0, len(principals.Groups))
	for k := range principals.GetGroups() {
		group := principals.GetGroups()[k]
		groups = append(groups, &apisecurity.UserGroup{
			Id: group.Id,
		})
	}

	return groups
}

// checkUserExist 检查用户是否存在
func (svr *Server) checkUserExist(users []*apisecurity.User) error {
	if len(users) == 0 {
		return nil
	}
	return svr.userSvr.GetUserHelper().CheckUsersExist(context.TODO(), users)
}

// checkUserGroupExist 检查用户组是否存在
func (svr *Server) checkGroupExist(groups []*apisecurity.UserGroup) error {
	if len(groups) == 0 {
		return nil
	}
	return svr.userSvr.GetUserHelper().CheckGroupsExist(context.TODO(), groups)
}

// checkResourceExist 检查资源是否存在
func (svr *Server) checkResourceExist(resources *apisecurity.StrategyResources) *apiservice.Response {
	namespaces := resources.GetNamespaces()

	nsCache := svr.cacheMgr.Namespace()
	for index := range namespaces {
		val := namespaces[index]
		if val.GetId().GetValue() == "*" {
			break
		}
		ns := nsCache.GetNamespace(val.GetId().GetValue())
		if ns == nil {
			return api.NewAuthResponse(apimodel.Code_NotFoundNamespace)
		}
	}

	services := resources.GetServices()
	svcCache := svr.cacheMgr.Service()
	for index := range services {
		val := services[index]
		if val.GetId().GetValue() == "*" {
			break
		}
		svc := svcCache.GetServiceByID(val.GetId().GetValue())
		if svc == nil {
			return api.NewAuthResponse(apimodel.Code_NotFoundService)
		}
	}

	return nil
}

// normalizeResource 对于资源进行归一化处理, 如果出现 * 的话，则该资源访问策略就是 *
func (svr *Server) normalizeResource(resources *apisecurity.StrategyResources) *apisecurity.StrategyResources {
	namespaces := resources.GetNamespaces()
	for index := range namespaces {
		val := namespaces[index]
		if val.GetId().GetValue() == "*" {
			resources.Namespaces = []*apisecurity.StrategyResourceEntry{{
				Id: utils.NewStringValue("*"),
			}}
			break
		}
	}
	services := resources.GetServices()
	for index := range services {
		val := services[index]
		if val.GetId().GetValue() == "*" {
			resources.Services = []*apisecurity.StrategyResourceEntry{{
				Id: utils.NewStringValue("*"),
			}}
			break
		}
	}
	return resources
}

// enrichPrincipalInfo 填充 principal 摘要信息
func (svr *Server) enrichPrincipalInfo(resp *apisecurity.AuthStrategy, data *authcommon.StrategyDetail) {
	users := make([]*apisecurity.Principal, 0, len(data.Principals))
	groups := make([]*apisecurity.Principal, 0, len(data.Principals))
	roles := make([]*apisecurity.Principal, 0, len(data.Principals))
	for index := range data.Principals {
		principal := data.Principals[index]
		switch principal.PrincipalType {
		case authcommon.PrincipalUser:
			user := svr.userSvr.GetUserHelper().GetUser(context.TODO(), &apisecurity.User{
				Id: wrapperspb.String(principal.PrincipalID),
			})
			if user == nil {
				continue
			}
			users = append(users, &apisecurity.Principal{
				Id:   utils.NewStringValue(user.GetId().GetValue()),
				Name: utils.NewStringValue(user.GetName().GetValue()),
			})
		case authcommon.PrincipalGroup:
			group := svr.userSvr.GetUserHelper().GetGroup(context.TODO(), &apisecurity.UserGroup{
				Id: wrapperspb.String(principal.PrincipalID),
			})
			if group == nil {
				continue
			}
			groups = append(groups, &apisecurity.Principal{
				Id:   utils.NewStringValue(group.GetId().GetValue()),
				Name: utils.NewStringValue(group.GetName().GetValue()),
			})
		case authcommon.PrincipalRole:
			role := svr.PolicyHelper().GetRole(principal.PrincipalID)
			if role == nil {
				continue
			}
			roles = append(roles, &apisecurity.Principal{
				Id:   utils.NewStringValue(role.ID),
				Name: utils.NewStringValue(role.Name),
			})
		}
	}

	resp.Principals = &apisecurity.Principals{
		Users:  users,
		Groups: groups,
		Roles:  roles,
	}
}

// enrichResourceInfo 填充资源摘要信息
func (svr *Server) enrichResourceInfo(ctx context.Context, resp *apisecurity.AuthStrategy, data *authcommon.StrategyDetail) {
	allMatch := map[apisecurity.ResourceType]struct{}{}
	resp.Resources = &apisecurity.StrategyResources{
		Namespaces:          make([]*apisecurity.StrategyResourceEntry, 0, 4),
		ConfigGroups:        make([]*apisecurity.StrategyResourceEntry, 0, 4),
		Services:            make([]*apisecurity.StrategyResourceEntry, 0, 4),
		RouteRules:          make([]*apisecurity.StrategyResourceEntry, 0, 4),
		RatelimitRules:      make([]*apisecurity.StrategyResourceEntry, 0, 4),
		CircuitbreakerRules: make([]*apisecurity.StrategyResourceEntry, 0, 4),
		FaultdetectRules:    make([]*apisecurity.StrategyResourceEntry, 0, 4),
		LaneRules:           make([]*apisecurity.StrategyResourceEntry, 0, 4),
		Users:               make([]*apisecurity.StrategyResourceEntry, 0, 4),
		UserGroups:          make([]*apisecurity.StrategyResourceEntry, 0, 4),
		Roles:               make([]*apisecurity.StrategyResourceEntry, 0, 4),
		AuthPolicies:        make([]*apisecurity.StrategyResourceEntry, 0, 4),
	}

	for index := range data.Resources {
		res := data.Resources[index]
		svr.enrichResourceDetial(ctx, res, allMatch, resp)
	}
}

var (
	resourceConvert = map[apisecurity.ResourceType]func(context.Context, *Server, authcommon.StrategyResource) *apisecurity.StrategyResourceEntry{
		// 注册、配置、治理
		apisecurity.ResourceType_Namespaces: func(ctx context.Context, svr *Server, item authcommon.StrategyResource) *apisecurity.StrategyResourceEntry {
			user := svr.cacheMgr.Namespace().GetNamespace(item.ResID)
			if user == nil {
				log.Warn("[Auth][Strategy] not found namespace in fill-info",
					zap.String("id", item.StrategyID), zap.String("res-id", item.ResID), utils.RequestID(ctx))
				return nil
			}
			return &apisecurity.StrategyResourceEntry{
				Id:        utils.NewStringValue(item.ResID),
				Namespace: utils.NewStringValue(user.Name),
				Name:      utils.NewStringValue(user.Name),
			}
		},
		apisecurity.ResourceType_Services: func(ctx context.Context, svr *Server, item authcommon.StrategyResource) *apisecurity.StrategyResourceEntry {
			user := svr.cacheMgr.Namespace().GetNamespace(item.ResID)
			if user == nil {
				log.Warn("[Auth][Strategy] not found namespace in fill-info",
					zap.String("id", item.StrategyID), zap.String("res-id", item.ResID), utils.RequestID(ctx))
				return nil
			}
			return &apisecurity.StrategyResourceEntry{
				Id:        utils.NewStringValue(item.ResID),
				Namespace: utils.NewStringValue(user.Name),
				Name:      utils.NewStringValue(user.Name),
			}
		},
		apisecurity.ResourceType_RouteRules: func(ctx context.Context, svr *Server, item authcommon.StrategyResource) *apisecurity.StrategyResourceEntry {
			user := svr.cacheMgr.RoutingConfig().GetRule(item.ResID)
			if user == nil {
				log.Warn("[Auth][Strategy] not found route_rule in fill-info",
					zap.String("id", item.StrategyID), zap.String("res-id", item.ResID), utils.RequestID(ctx))
				return nil
			}
			return &apisecurity.StrategyResourceEntry{
				Id:        utils.NewStringValue(item.ResID),
				Namespace: utils.NewStringValue(user.Name),
				Name:      utils.NewStringValue(user.Name),
			}
		},
		apisecurity.ResourceType_LaneRules: func(ctx context.Context, svr *Server, item authcommon.StrategyResource) *apisecurity.StrategyResourceEntry {
			user := svr.cacheMgr.LaneRule().GetRule(item.ResID)
			if user == nil {
				log.Warn("[Auth][Strategy] not found lane_rule in fill-info",
					zap.String("id", item.StrategyID), zap.String("res-id", item.ResID), utils.RequestID(ctx))
				return nil
			}
			return &apisecurity.StrategyResourceEntry{
				Id:        utils.NewStringValue(item.ResID),
				Namespace: utils.NewStringValue(user.Name),
				Name:      utils.NewStringValue(user.Name),
			}
		},
		apisecurity.ResourceType_RateLimitRules: func(ctx context.Context, svr *Server, item authcommon.StrategyResource) *apisecurity.StrategyResourceEntry {
			user := svr.cacheMgr.RateLimit().GetRule(item.ResID)
			if user == nil {
				log.Warn("[Auth][Strategy] not found ratelimit_rule in fill-info",
					zap.String("id", item.StrategyID), zap.String("res-id", item.ResID), utils.RequestID(ctx))
				return nil
			}
			return &apisecurity.StrategyResourceEntry{
				Id:        utils.NewStringValue(item.ResID),
				Namespace: utils.NewStringValue(user.Name),
				Name:      utils.NewStringValue(user.Name),
			}
		},
		// 鉴权资源
		apisecurity.ResourceType_Users: func(ctx context.Context, svr *Server, item authcommon.StrategyResource) *apisecurity.StrategyResourceEntry {
			user := svr.cacheMgr.User().GetUserByID(item.ResID)
			if user == nil {
				log.Warn("[Auth][Strategy] not found user in fill-info",
					zap.String("id", item.StrategyID), zap.String("res-id", item.ResID), utils.RequestID(ctx))
				return nil
			}
			return &apisecurity.StrategyResourceEntry{
				Id:   utils.NewStringValue(item.ResID),
				Name: utils.NewStringValue(user.Name),
			}
		},
		apisecurity.ResourceType_UserGroups: func(ctx context.Context, svr *Server, item authcommon.StrategyResource) *apisecurity.StrategyResourceEntry {
			user := svr.cacheMgr.User().GetGroup(item.ResID)
			if user == nil {
				log.Warn("[Auth][Strategy] not found user_group in fill-info",
					zap.String("id", item.StrategyID), zap.String("res-id", item.ResID), utils.RequestID(ctx))
				return nil
			}
			return &apisecurity.StrategyResourceEntry{
				Id:   utils.NewStringValue(item.ResID),
				Name: utils.NewStringValue(user.Name),
			}
		},
		apisecurity.ResourceType_Roles: func(ctx context.Context, svr *Server, item authcommon.StrategyResource) *apisecurity.StrategyResourceEntry {
			user := svr.cacheMgr.Role().GetRole(item.ResID)
			if user == nil {
				log.Warn("[Auth][Strategy] not found role in fill-info",
					zap.String("id", item.StrategyID), zap.String("res-id", item.ResID), utils.RequestID(ctx))
				return nil
			}
			return &apisecurity.StrategyResourceEntry{
				Id:   utils.NewStringValue(item.ResID),
				Name: utils.NewStringValue(user.Name),
			}
		},
		apisecurity.ResourceType_PolicyRules: func(ctx context.Context, svr *Server, item authcommon.StrategyResource) *apisecurity.StrategyResourceEntry {
			user := svr.cacheMgr.AuthStrategy().GetPolicyRule(item.ResID)
			if user == nil {
				log.Warn("[Auth][Strategy] not found auth_policy in fill-info",
					zap.String("id", item.StrategyID), zap.String("res-id", item.ResID), utils.RequestID(ctx))
				return nil
			}
			return &apisecurity.StrategyResourceEntry{
				Id:   utils.NewStringValue(item.ResID),
				Name: utils.NewStringValue(user.Name),
			}
		},
	}
)

func (svr *Server) enrichResourceDetial(ctx context.Context,
	item authcommon.StrategyResource,
	allMatch map[apisecurity.ResourceType]struct{},
	resp *apisecurity.AuthStrategy) {

	resType := apisecurity.ResourceType(item.ResType)

	message := resp.Resources.ProtoReflect()
	fields := message.Descriptor().Fields()
	list := message.Get(fields.ByName(protoreflect.Name(resourceToFieldNames[resType]))).List()

	if item.ResID == "*" {
		allMatch[resType] = struct{}{}
		resp.Resources.Users = []*apisecurity.StrategyResourceEntry{
			{
				Id:        utils.NewStringValue("*"),
				Namespace: utils.NewStringValue("*"),
				Name:      utils.NewStringValue("*"),
			},
		}
		return
	}
	if _, ok := allMatch[resType]; !ok {
		if data := resourceConvert[resType](ctx, svr, item); data != nil {
			list.Append(protoreflect.ValueOf(data))
		}
	}
}

type resourceFilter struct {
	ns   map[string]struct{}
	svc  map[string]struct{}
	conf map[string]struct{}
}

func (f *resourceFilter) GetFilter(t apisecurity.ResourceType) (map[string]struct{}, bool) {
	switch t {
	case apisecurity.ResourceType_Namespaces:
		return f.ns, true
	case apisecurity.ResourceType_Services:
		return f.svc, true
	case apisecurity.ResourceType_ConfigGroups:
		return f.conf, true
	}
	return nil, false
}

// filter different types of Strategy resources
func resourceDeduplication(resources []authcommon.StrategyResource) []authcommon.StrategyResource {
	rLen := len(resources)
	ret := make([]authcommon.StrategyResource, 0, rLen)
	rf := resourceFilter{
		ns:   make(map[string]struct{}, rLen),
		svc:  make(map[string]struct{}, rLen),
		conf: make(map[string]struct{}, rLen),
	}

	est := struct{}{}
	for i := range resources {
		res := resources[i]
		filter, ok := rf.GetFilter(apisecurity.ResourceType(res.ResType))
		if !ok {
			continue
		}
		if _, exist := filter[res.ResID]; !exist {
			rf.ns[res.ResID] = est
			ret = append(ret, res)
		}
	}
	return ret
}

// ----- 以下代码为历史代码，当前留作备份

// func (svr *Server) enrichNamespaceResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_Namespaces] = struct{}{}
// 		resp.Resources.Namespaces = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_Namespaces]; !ok {
// 		ns := svr.cacheMgr.Namespace().GetNamespace(item.ResID)
// 		if ns == nil {
// 			log.Warn("[Auth][Strategy] not found namespace in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.Services = append(resp.Resources.Services, &apisecurity.StrategyResourceEntry{
// 			Id:        utils.NewStringValue(item.ResID),
// 			Namespace: utils.NewStringValue(ns.Name),
// 			Name:      utils.NewStringValue(ns.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichServiceResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_Services] = struct{}{}
// 		resp.Resources.Services = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_Services]; !ok {
// 		svc := svr.cacheMgr.Service().GetServiceByID(item.ResID)
// 		if svc == nil {
// 			log.Warn("[Auth][Strategy] not found service in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.Services = append(resp.Resources.Services, &apisecurity.StrategyResourceEntry{
// 			Id:        utils.NewStringValue(item.ResID),
// 			Namespace: utils.NewStringValue(svc.Namespace),
// 			Name:      utils.NewStringValue(svc.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichRateLimitRulesResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_RateLimitRules] = struct{}{}
// 		resp.Resources.RatelimitRules = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_RateLimitRules]; !ok {
// 		user := svr.cacheMgr.RateLimit().GetRule(item.ResID)
// 		if user == nil {
// 			log.Warn("[Auth][Strategy] not found ratelimit_rule in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.RatelimitRules = append(resp.Resources.RatelimitRules, &apisecurity.StrategyResourceEntry{
// 			Id:   utils.NewStringValue(item.ResID),
// 			Name: utils.NewStringValue(user.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichRouteRulesResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_RouteRules] = struct{}{}
// 		resp.Resources.RouteRules = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_RouteRules]; !ok {
// 		user := svr.cacheMgr.RoutingConfig().GetRule(item.ResID)
// 		if user == nil {
// 			log.Warn("[Auth][Strategy] not found route_rule in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.RouteRules = append(resp.Resources.RouteRules, &apisecurity.StrategyResourceEntry{
// 			Id:   utils.NewStringValue(item.ResID),
// 			Name: utils.NewStringValue(user.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichLaneRulesResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_LaneRules] = struct{}{}
// 		resp.Resources.LaneRules = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_LaneRules]; !ok {
// 		user := svr.cacheMgr.LaneRule().GetRule(item.ResID)
// 		if user == nil {
// 			log.Warn("[Auth][Strategy] not found lane_rule in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.LaneRules = append(resp.Resources.LaneRules, &apisecurity.StrategyResourceEntry{
// 			Id:   utils.NewStringValue(item.ResID),
// 			Name: utils.NewStringValue(user.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichCircuitBreakerRulesResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_CircuitBreakerRules] = struct{}{}
// 		resp.Resources.LaneRules = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_CircuitBreakerRules]; !ok {
// 		user := svr.cacheMgr.LaneRule().GetRule(item.ResID)
// 		if user == nil {
// 			log.Warn("[Auth][Strategy] not found circuitbreaker_rule in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.LaneRules = append(resp.Resources.LaneRules, &apisecurity.StrategyResourceEntry{
// 			Id:   utils.NewStringValue(item.ResID),
// 			Name: utils.NewStringValue(user.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichFaultDetectRulesResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_FaultDetectRules] = struct{}{}
// 		resp.Resources.FaultdetectRules = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_FaultDetectRules]; !ok {
// 		user := svr.cacheMgr.LaneRule().GetRule(item.ResID)
// 		if user == nil {
// 			log.Warn("[Auth][Strategy] not found faultdetect_rule in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.FaultdetectRules = append(resp.Resources.FaultdetectRules, &apisecurity.StrategyResourceEntry{
// 			Id:   utils.NewStringValue(item.ResID),
// 			Name: utils.NewStringValue(user.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichConfigGroupResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_ConfigGroups] = struct{}{}
// 		resp.Resources.ConfigGroups = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_ConfigGroups]; !ok {
// 		groupId, err := strconv.ParseUint(item.ResID, 10, 64)
// 		if err != nil {
// 			log.Warn("[Auth][Strategy] invalid resource id",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("config_file_group", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		group := svr.cacheMgr.ConfigGroup().GetGroupByID(groupId)
// 		if group == nil {
// 			log.Warn("[Auth][Strategy] not found config_file_group in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("config_file_group", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.ConfigGroups = append(resp.Resources.ConfigGroups, &apisecurity.StrategyResourceEntry{
// 			Id:        utils.NewStringValue(item.ResID),
// 			Namespace: utils.NewStringValue(group.Namespace),
// 			Name:      utils.NewStringValue(group.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichUserResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_Users] = struct{}{}
// 		resp.Resources.Users = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_Users]; !ok {
// 		user := svr.cacheMgr.User().GetUserByID(item.ResID)
// 		if user == nil {
// 			log.Warn("[Auth][Strategy] not found user in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.Users = append(resp.Resources.Users, &apisecurity.StrategyResourceEntry{
// 			Id:   utils.NewStringValue(item.ResID),
// 			Name: utils.NewStringValue(user.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichUserGroupsResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_UserGroups] = struct{}{}
// 		resp.Resources.UserGroups = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_UserGroups]; !ok {
// 		user := svr.cacheMgr.User().GetGroup(item.ResID)
// 		if user == nil {
// 			log.Warn("[Auth][Strategy] not found user_group in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.UserGroups = append(resp.Resources.UserGroups, &apisecurity.StrategyResourceEntry{
// 			Id:   utils.NewStringValue(item.ResID),
// 			Name: utils.NewStringValue(user.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichRolesResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_Roles] = struct{}{}
// 		resp.Resources.Roles = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_Roles]; !ok {
// 		user := svr.cacheMgr.Role().GetRole(item.ResID)
// 		if user == nil {
// 			log.Warn("[Auth][Strategy] not found role in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.Roles = append(resp.Resources.Roles, &apisecurity.StrategyResourceEntry{
// 			Id:   utils.NewStringValue(item.ResID),
// 			Name: utils.NewStringValue(user.Name),
// 		})
// 	}
// }

// func (svr *Server) enrichPolicyRulesResource(ctx context.Context, item authcommon.StrategyResource,
// 	allMatch map[apisecurity.ResourceType]struct{}, resp *apisecurity.AuthStrategy) {
// 	if item.ResID == "*" {
// 		allMatch[apisecurity.ResourceType_PolicyRules] = struct{}{}
// 		resp.Resources.AuthPolicies = []*apisecurity.StrategyResourceEntry{
// 			{
// 				Id:        utils.NewStringValue("*"),
// 				Namespace: utils.NewStringValue("*"),
// 				Name:      utils.NewStringValue("*"),
// 			},
// 		}
// 		return
// 	}
// 	if _, ok := allMatch[apisecurity.ResourceType_PolicyRules]; !ok {
// 		user := svr.cacheMgr.AuthStrategy().GetPolicyRule(item.ResID)
// 		if user == nil {
// 			log.Warn("[Auth][Strategy] not found auth_policy in fill-info",
// 				zap.String("id", resp.GetId().GetValue()), zap.String("res-id", item.ResID), utils.RequestID(ctx))
// 			return
// 		}
// 		resp.Resources.AuthPolicies = append(resp.Resources.AuthPolicies, &apisecurity.StrategyResourceEntry{
// 			Id:   utils.NewStringValue(item.ResID),
// 			Name: utils.NewStringValue(user.Name),
// 		})
// 	}
// }
