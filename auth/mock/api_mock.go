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
// Code generated by MockGen. DO NOT EDIT.
// Source: api.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	security "github.com/polarismesh/specification/source/go/api/v1/security"
	service_manage "github.com/polarismesh/specification/source/go/api/v1/service_manage"

	auth "github.com/polarismesh/polaris/auth"
	cache "github.com/polarismesh/polaris/cache"
	model "github.com/polarismesh/polaris/common/model"
	store "github.com/polarismesh/polaris/store"
)

// MockAuthServer is a mock of AuthServer interface.
type MockAuthServer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServerMockRecorder
}

// MockAuthServerMockRecorder is the mock recorder for MockAuthServer.
type MockAuthServerMockRecorder struct {
	mock *MockAuthServer
}

// NewMockAuthServer creates a new mock instance.
func NewMockAuthServer(ctrl *gomock.Controller) *MockAuthServer {
	mock := &MockAuthServer{ctrl: ctrl}
	mock.recorder = &MockAuthServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServer) EXPECT() *MockAuthServerMockRecorder {
	return m.recorder
}

// AfterResourceOperation mocks base method.
func (m *MockAuthServer) AfterResourceOperation(afterCtx *model.AcquireContext) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AfterResourceOperation", afterCtx)
	ret0, _ := ret[0].(error)
	return ret0
}

// AfterResourceOperation indicates an expected call of AfterResourceOperation.
func (mr *MockAuthServerMockRecorder) AfterResourceOperation(afterCtx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterResourceOperation", reflect.TypeOf((*MockAuthServer)(nil).AfterResourceOperation), afterCtx)
}

// CreateGroup mocks base method.
func (m *MockAuthServer) CreateGroup(ctx context.Context, group *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", ctx, group)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockAuthServerMockRecorder) CreateGroup(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockAuthServer)(nil).CreateGroup), ctx, group)
}

// CreateStrategy mocks base method.
func (m *MockAuthServer) CreateStrategy(ctx context.Context, strategy *security.AuthStrategy) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStrategy", ctx, strategy)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// CreateStrategy indicates an expected call of CreateStrategy.
func (mr *MockAuthServerMockRecorder) CreateStrategy(ctx, strategy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStrategy", reflect.TypeOf((*MockAuthServer)(nil).CreateStrategy), ctx, strategy)
}

// CreateUsers mocks base method.
func (m *MockAuthServer) CreateUsers(ctx context.Context, users []*security.User) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUsers", ctx, users)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// CreateUsers indicates an expected call of CreateUsers.
func (mr *MockAuthServerMockRecorder) CreateUsers(ctx, users interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUsers", reflect.TypeOf((*MockAuthServer)(nil).CreateUsers), ctx, users)
}

// DeleteGroups mocks base method.
func (m *MockAuthServer) DeleteGroups(ctx context.Context, group []*security.UserGroup) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroups", ctx, group)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// DeleteGroups indicates an expected call of DeleteGroups.
func (mr *MockAuthServerMockRecorder) DeleteGroups(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroups", reflect.TypeOf((*MockAuthServer)(nil).DeleteGroups), ctx, group)
}

// DeleteStrategies mocks base method.
func (m *MockAuthServer) DeleteStrategies(ctx context.Context, reqs []*security.AuthStrategy) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStrategies", ctx, reqs)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// DeleteStrategies indicates an expected call of DeleteStrategies.
func (mr *MockAuthServerMockRecorder) DeleteStrategies(ctx, reqs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStrategies", reflect.TypeOf((*MockAuthServer)(nil).DeleteStrategies), ctx, reqs)
}

// DeleteUsers mocks base method.
func (m *MockAuthServer) DeleteUsers(ctx context.Context, users []*security.User) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUsers", ctx, users)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// DeleteUsers indicates an expected call of DeleteUsers.
func (mr *MockAuthServerMockRecorder) DeleteUsers(ctx, users interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUsers", reflect.TypeOf((*MockAuthServer)(nil).DeleteUsers), ctx, users)
}

// GetAuthChecker mocks base method.
func (m *MockAuthServer) GetAuthChecker() auth.AuthChecker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthChecker")
	ret0, _ := ret[0].(auth.AuthChecker)
	return ret0
}

// GetAuthChecker indicates an expected call of GetAuthChecker.
func (mr *MockAuthServerMockRecorder) GetAuthChecker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthChecker", reflect.TypeOf((*MockAuthServer)(nil).GetAuthChecker))
}

// GetGroup mocks base method.
func (m *MockAuthServer) GetGroup(ctx context.Context, req *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", ctx, req)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetGroup indicates an expected call of GetGroup.
func (mr *MockAuthServerMockRecorder) GetGroup(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockAuthServer)(nil).GetGroup), ctx, req)
}

// GetGroupToken mocks base method.
func (m *MockAuthServer) GetGroupToken(ctx context.Context, group *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupToken", ctx, group)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetGroupToken indicates an expected call of GetGroupToken.
func (mr *MockAuthServerMockRecorder) GetGroupToken(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupToken", reflect.TypeOf((*MockAuthServer)(nil).GetGroupToken), ctx, group)
}

// GetGroups mocks base method.
func (m *MockAuthServer) GetGroups(ctx context.Context, query map[string]string) *service_manage.BatchQueryResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroups", ctx, query)
	ret0, _ := ret[0].(*service_manage.BatchQueryResponse)
	return ret0
}

// GetGroups indicates an expected call of GetGroups.
func (mr *MockAuthServerMockRecorder) GetGroups(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroups", reflect.TypeOf((*MockAuthServer)(nil).GetGroups), ctx, query)
}

// GetPrincipalResources mocks base method.
func (m *MockAuthServer) GetPrincipalResources(ctx context.Context, query map[string]string) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrincipalResources", ctx, query)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetPrincipalResources indicates an expected call of GetPrincipalResources.
func (mr *MockAuthServerMockRecorder) GetPrincipalResources(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrincipalResources", reflect.TypeOf((*MockAuthServer)(nil).GetPrincipalResources), ctx, query)
}

// GetStrategies mocks base method.
func (m *MockAuthServer) GetStrategies(ctx context.Context, query map[string]string) *service_manage.BatchQueryResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStrategies", ctx, query)
	ret0, _ := ret[0].(*service_manage.BatchQueryResponse)
	return ret0
}

// GetStrategies indicates an expected call of GetStrategies.
func (mr *MockAuthServerMockRecorder) GetStrategies(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStrategies", reflect.TypeOf((*MockAuthServer)(nil).GetStrategies), ctx, query)
}

// GetStrategy mocks base method.
func (m *MockAuthServer) GetStrategy(ctx context.Context, strategy *security.AuthStrategy) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStrategy", ctx, strategy)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetStrategy indicates an expected call of GetStrategy.
func (mr *MockAuthServerMockRecorder) GetStrategy(ctx, strategy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStrategy", reflect.TypeOf((*MockAuthServer)(nil).GetStrategy), ctx, strategy)
}

// GetUserToken mocks base method.
func (m *MockAuthServer) GetUserToken(ctx context.Context, user *security.User) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserToken", ctx, user)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetUserToken indicates an expected call of GetUserToken.
func (mr *MockAuthServerMockRecorder) GetUserToken(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserToken", reflect.TypeOf((*MockAuthServer)(nil).GetUserToken), ctx, user)
}

// GetUsers mocks base method.
func (m *MockAuthServer) GetUsers(ctx context.Context, query map[string]string) *service_manage.BatchQueryResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx, query)
	ret0, _ := ret[0].(*service_manage.BatchQueryResponse)
	return ret0
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockAuthServerMockRecorder) GetUsers(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockAuthServer)(nil).GetUsers), ctx, query)
}

// Initialize mocks base method.
func (m *MockAuthServer) Initialize(authOpt *auth.Config, storage store.Store, cacheMgn *cache.CacheManager) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", authOpt, storage, cacheMgn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockAuthServerMockRecorder) Initialize(authOpt, storage, cacheMgn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockAuthServer)(nil).Initialize), authOpt, storage, cacheMgn)
}

// Login mocks base method.
func (m *MockAuthServer) Login(req *security.LoginRequest) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", req)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// Login indicates an expected call of Login.
func (mr *MockAuthServerMockRecorder) Login(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthServer)(nil).Login), req)
}

// Name mocks base method.
func (m *MockAuthServer) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockAuthServerMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockAuthServer)(nil).Name))
}

// ResetGroupToken mocks base method.
func (m *MockAuthServer) ResetGroupToken(ctx context.Context, group *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetGroupToken", ctx, group)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// ResetGroupToken indicates an expected call of ResetGroupToken.
func (mr *MockAuthServerMockRecorder) ResetGroupToken(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetGroupToken", reflect.TypeOf((*MockAuthServer)(nil).ResetGroupToken), ctx, group)
}

// ResetUserToken mocks base method.
func (m *MockAuthServer) ResetUserToken(ctx context.Context, user *security.User) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetUserToken", ctx, user)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// ResetUserToken indicates an expected call of ResetUserToken.
func (mr *MockAuthServerMockRecorder) ResetUserToken(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetUserToken", reflect.TypeOf((*MockAuthServer)(nil).ResetUserToken), ctx, user)
}

// UpdateGroupToken mocks base method.
func (m *MockAuthServer) UpdateGroupToken(ctx context.Context, group *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroupToken", ctx, group)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// UpdateGroupToken indicates an expected call of UpdateGroupToken.
func (mr *MockAuthServerMockRecorder) UpdateGroupToken(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroupToken", reflect.TypeOf((*MockAuthServer)(nil).UpdateGroupToken), ctx, group)
}

// UpdateGroups mocks base method.
func (m *MockAuthServer) UpdateGroups(ctx context.Context, groups []*security.ModifyUserGroup) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroups", ctx, groups)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// UpdateGroups indicates an expected call of UpdateGroups.
func (mr *MockAuthServerMockRecorder) UpdateGroups(ctx, groups interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroups", reflect.TypeOf((*MockAuthServer)(nil).UpdateGroups), ctx, groups)
}

// UpdateStrategies mocks base method.
func (m *MockAuthServer) UpdateStrategies(ctx context.Context, reqs []*security.ModifyAuthStrategy) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStrategies", ctx, reqs)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// UpdateStrategies indicates an expected call of UpdateStrategies.
func (mr *MockAuthServerMockRecorder) UpdateStrategies(ctx, reqs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStrategies", reflect.TypeOf((*MockAuthServer)(nil).UpdateStrategies), ctx, reqs)
}

// UpdateUser mocks base method.
func (m *MockAuthServer) UpdateUser(ctx context.Context, user *security.User) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, user)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockAuthServerMockRecorder) UpdateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockAuthServer)(nil).UpdateUser), ctx, user)
}

// UpdateUserPassword mocks base method.
func (m *MockAuthServer) UpdateUserPassword(ctx context.Context, req *security.ModifyUserPassword) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", ctx, req)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockAuthServerMockRecorder) UpdateUserPassword(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockAuthServer)(nil).UpdateUserPassword), ctx, req)
}

// UpdateUserToken mocks base method.
func (m *MockAuthServer) UpdateUserToken(ctx context.Context, user *security.User) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserToken", ctx, user)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// UpdateUserToken indicates an expected call of UpdateUserToken.
func (mr *MockAuthServerMockRecorder) UpdateUserToken(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserToken", reflect.TypeOf((*MockAuthServer)(nil).UpdateUserToken), ctx, user)
}

// MockAuthChecker is a mock of AuthChecker interface.
type MockAuthChecker struct {
	ctrl     *gomock.Controller
	recorder *MockAuthCheckerMockRecorder
}

// MockAuthCheckerMockRecorder is the mock recorder for MockAuthChecker.
type MockAuthCheckerMockRecorder struct {
	mock *MockAuthChecker
}

// NewMockAuthChecker creates a new mock instance.
func NewMockAuthChecker(ctrl *gomock.Controller) *MockAuthChecker {
	mock := &MockAuthChecker{ctrl: ctrl}
	mock.recorder = &MockAuthCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthChecker) EXPECT() *MockAuthCheckerMockRecorder {
	return m.recorder
}

// CheckClientPermission mocks base method.
func (m *MockAuthChecker) CheckClientPermission(preCtx *model.AcquireContext) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckClientPermission", preCtx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckClientPermission indicates an expected call of CheckClientPermission.
func (mr *MockAuthCheckerMockRecorder) CheckClientPermission(preCtx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckClientPermission", reflect.TypeOf((*MockAuthChecker)(nil).CheckClientPermission), preCtx)
}

// CheckConsolePermission mocks base method.
func (m *MockAuthChecker) CheckConsolePermission(preCtx *model.AcquireContext) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckConsolePermission", preCtx)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckConsolePermission indicates an expected call of CheckConsolePermission.
func (mr *MockAuthCheckerMockRecorder) CheckConsolePermission(preCtx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckConsolePermission", reflect.TypeOf((*MockAuthChecker)(nil).CheckConsolePermission), preCtx)
}

// Initialize mocks base method.
func (m *MockAuthChecker) Initialize(options *auth.Config, storage store.Store, cacheMgn *cache.CacheManager) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", options, storage, cacheMgn)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockAuthCheckerMockRecorder) Initialize(options, storage, cacheMgn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockAuthChecker)(nil).Initialize), options, storage, cacheMgn)
}

// IsOpenClientAuth mocks base method.
func (m *MockAuthChecker) IsOpenClientAuth() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsOpenClientAuth")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsOpenClientAuth indicates an expected call of IsOpenClientAuth.
func (mr *MockAuthCheckerMockRecorder) IsOpenClientAuth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsOpenClientAuth", reflect.TypeOf((*MockAuthChecker)(nil).IsOpenClientAuth))
}

// IsOpenConsoleAuth mocks base method.
func (m *MockAuthChecker) IsOpenConsoleAuth() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsOpenConsoleAuth")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsOpenConsoleAuth indicates an expected call of IsOpenConsoleAuth.
func (mr *MockAuthCheckerMockRecorder) IsOpenConsoleAuth() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsOpenConsoleAuth", reflect.TypeOf((*MockAuthChecker)(nil).IsOpenConsoleAuth))
}

// VerifyCredential mocks base method.
func (m *MockAuthChecker) VerifyCredential(preCtx *model.AcquireContext) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyCredential", preCtx)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyCredential indicates an expected call of VerifyCredential.
func (mr *MockAuthCheckerMockRecorder) VerifyCredential(preCtx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyCredential", reflect.TypeOf((*MockAuthChecker)(nil).VerifyCredential), preCtx)
}

// MockUserOperator is a mock of UserOperator interface.
type MockUserOperator struct {
	ctrl     *gomock.Controller
	recorder *MockUserOperatorMockRecorder
}

// MockUserOperatorMockRecorder is the mock recorder for MockUserOperator.
type MockUserOperatorMockRecorder struct {
	mock *MockUserOperator
}

// NewMockUserOperator creates a new mock instance.
func NewMockUserOperator(ctrl *gomock.Controller) *MockUserOperator {
	mock := &MockUserOperator{ctrl: ctrl}
	mock.recorder = &MockUserOperatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserOperator) EXPECT() *MockUserOperatorMockRecorder {
	return m.recorder
}

// CreateUsers mocks base method.
func (m *MockUserOperator) CreateUsers(ctx context.Context, users []*security.User) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUsers", ctx, users)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// CreateUsers indicates an expected call of CreateUsers.
func (mr *MockUserOperatorMockRecorder) CreateUsers(ctx, users interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUsers", reflect.TypeOf((*MockUserOperator)(nil).CreateUsers), ctx, users)
}

// DeleteUsers mocks base method.
func (m *MockUserOperator) DeleteUsers(ctx context.Context, users []*security.User) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUsers", ctx, users)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// DeleteUsers indicates an expected call of DeleteUsers.
func (mr *MockUserOperatorMockRecorder) DeleteUsers(ctx, users interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUsers", reflect.TypeOf((*MockUserOperator)(nil).DeleteUsers), ctx, users)
}

// GetUserToken mocks base method.
func (m *MockUserOperator) GetUserToken(ctx context.Context, user *security.User) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserToken", ctx, user)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetUserToken indicates an expected call of GetUserToken.
func (mr *MockUserOperatorMockRecorder) GetUserToken(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserToken", reflect.TypeOf((*MockUserOperator)(nil).GetUserToken), ctx, user)
}

// GetUsers mocks base method.
func (m *MockUserOperator) GetUsers(ctx context.Context, query map[string]string) *service_manage.BatchQueryResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", ctx, query)
	ret0, _ := ret[0].(*service_manage.BatchQueryResponse)
	return ret0
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUserOperatorMockRecorder) GetUsers(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUserOperator)(nil).GetUsers), ctx, query)
}

// ResetUserToken mocks base method.
func (m *MockUserOperator) ResetUserToken(ctx context.Context, user *security.User) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetUserToken", ctx, user)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// ResetUserToken indicates an expected call of ResetUserToken.
func (mr *MockUserOperatorMockRecorder) ResetUserToken(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetUserToken", reflect.TypeOf((*MockUserOperator)(nil).ResetUserToken), ctx, user)
}

// UpdateUser mocks base method.
func (m *MockUserOperator) UpdateUser(ctx context.Context, user *security.User) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, user)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserOperatorMockRecorder) UpdateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserOperator)(nil).UpdateUser), ctx, user)
}

// UpdateUserPassword mocks base method.
func (m *MockUserOperator) UpdateUserPassword(ctx context.Context, req *security.ModifyUserPassword) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserPassword", ctx, req)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// UpdateUserPassword indicates an expected call of UpdateUserPassword.
func (mr *MockUserOperatorMockRecorder) UpdateUserPassword(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserPassword", reflect.TypeOf((*MockUserOperator)(nil).UpdateUserPassword), ctx, req)
}

// UpdateUserToken mocks base method.
func (m *MockUserOperator) UpdateUserToken(ctx context.Context, user *security.User) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserToken", ctx, user)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// UpdateUserToken indicates an expected call of UpdateUserToken.
func (mr *MockUserOperatorMockRecorder) UpdateUserToken(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserToken", reflect.TypeOf((*MockUserOperator)(nil).UpdateUserToken), ctx, user)
}

// MockGroupOperator is a mock of GroupOperator interface.
type MockGroupOperator struct {
	ctrl     *gomock.Controller
	recorder *MockGroupOperatorMockRecorder
}

// MockGroupOperatorMockRecorder is the mock recorder for MockGroupOperator.
type MockGroupOperatorMockRecorder struct {
	mock *MockGroupOperator
}

// NewMockGroupOperator creates a new mock instance.
func NewMockGroupOperator(ctrl *gomock.Controller) *MockGroupOperator {
	mock := &MockGroupOperator{ctrl: ctrl}
	mock.recorder = &MockGroupOperatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupOperator) EXPECT() *MockGroupOperatorMockRecorder {
	return m.recorder
}

// CreateGroup mocks base method.
func (m *MockGroupOperator) CreateGroup(ctx context.Context, group *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", ctx, group)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockGroupOperatorMockRecorder) CreateGroup(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockGroupOperator)(nil).CreateGroup), ctx, group)
}

// DeleteGroups mocks base method.
func (m *MockGroupOperator) DeleteGroups(ctx context.Context, group []*security.UserGroup) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroups", ctx, group)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// DeleteGroups indicates an expected call of DeleteGroups.
func (mr *MockGroupOperatorMockRecorder) DeleteGroups(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroups", reflect.TypeOf((*MockGroupOperator)(nil).DeleteGroups), ctx, group)
}

// GetGroup mocks base method.
func (m *MockGroupOperator) GetGroup(ctx context.Context, req *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", ctx, req)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetGroup indicates an expected call of GetGroup.
func (mr *MockGroupOperatorMockRecorder) GetGroup(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockGroupOperator)(nil).GetGroup), ctx, req)
}

// GetGroupToken mocks base method.
func (m *MockGroupOperator) GetGroupToken(ctx context.Context, group *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupToken", ctx, group)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetGroupToken indicates an expected call of GetGroupToken.
func (mr *MockGroupOperatorMockRecorder) GetGroupToken(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupToken", reflect.TypeOf((*MockGroupOperator)(nil).GetGroupToken), ctx, group)
}

// GetGroups mocks base method.
func (m *MockGroupOperator) GetGroups(ctx context.Context, query map[string]string) *service_manage.BatchQueryResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroups", ctx, query)
	ret0, _ := ret[0].(*service_manage.BatchQueryResponse)
	return ret0
}

// GetGroups indicates an expected call of GetGroups.
func (mr *MockGroupOperatorMockRecorder) GetGroups(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroups", reflect.TypeOf((*MockGroupOperator)(nil).GetGroups), ctx, query)
}

// ResetGroupToken mocks base method.
func (m *MockGroupOperator) ResetGroupToken(ctx context.Context, group *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetGroupToken", ctx, group)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// ResetGroupToken indicates an expected call of ResetGroupToken.
func (mr *MockGroupOperatorMockRecorder) ResetGroupToken(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetGroupToken", reflect.TypeOf((*MockGroupOperator)(nil).ResetGroupToken), ctx, group)
}

// UpdateGroupToken mocks base method.
func (m *MockGroupOperator) UpdateGroupToken(ctx context.Context, group *security.UserGroup) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroupToken", ctx, group)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// UpdateGroupToken indicates an expected call of UpdateGroupToken.
func (mr *MockGroupOperatorMockRecorder) UpdateGroupToken(ctx, group interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroupToken", reflect.TypeOf((*MockGroupOperator)(nil).UpdateGroupToken), ctx, group)
}

// UpdateGroups mocks base method.
func (m *MockGroupOperator) UpdateGroups(ctx context.Context, groups []*security.ModifyUserGroup) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroups", ctx, groups)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// UpdateGroups indicates an expected call of UpdateGroups.
func (mr *MockGroupOperatorMockRecorder) UpdateGroups(ctx, groups interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroups", reflect.TypeOf((*MockGroupOperator)(nil).UpdateGroups), ctx, groups)
}

// MockStrategyOperator is a mock of StrategyOperator interface.
type MockStrategyOperator struct {
	ctrl     *gomock.Controller
	recorder *MockStrategyOperatorMockRecorder
}

// MockStrategyOperatorMockRecorder is the mock recorder for MockStrategyOperator.
type MockStrategyOperatorMockRecorder struct {
	mock *MockStrategyOperator
}

// NewMockStrategyOperator creates a new mock instance.
func NewMockStrategyOperator(ctrl *gomock.Controller) *MockStrategyOperator {
	mock := &MockStrategyOperator{ctrl: ctrl}
	mock.recorder = &MockStrategyOperatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStrategyOperator) EXPECT() *MockStrategyOperatorMockRecorder {
	return m.recorder
}

// CreateStrategy mocks base method.
func (m *MockStrategyOperator) CreateStrategy(ctx context.Context, strategy *security.AuthStrategy) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStrategy", ctx, strategy)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// CreateStrategy indicates an expected call of CreateStrategy.
func (mr *MockStrategyOperatorMockRecorder) CreateStrategy(ctx, strategy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStrategy", reflect.TypeOf((*MockStrategyOperator)(nil).CreateStrategy), ctx, strategy)
}

// DeleteStrategies mocks base method.
func (m *MockStrategyOperator) DeleteStrategies(ctx context.Context, reqs []*security.AuthStrategy) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStrategies", ctx, reqs)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// DeleteStrategies indicates an expected call of DeleteStrategies.
func (mr *MockStrategyOperatorMockRecorder) DeleteStrategies(ctx, reqs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStrategies", reflect.TypeOf((*MockStrategyOperator)(nil).DeleteStrategies), ctx, reqs)
}

// GetPrincipalResources mocks base method.
func (m *MockStrategyOperator) GetPrincipalResources(ctx context.Context, query map[string]string) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPrincipalResources", ctx, query)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetPrincipalResources indicates an expected call of GetPrincipalResources.
func (mr *MockStrategyOperatorMockRecorder) GetPrincipalResources(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPrincipalResources", reflect.TypeOf((*MockStrategyOperator)(nil).GetPrincipalResources), ctx, query)
}

// GetStrategies mocks base method.
func (m *MockStrategyOperator) GetStrategies(ctx context.Context, query map[string]string) *service_manage.BatchQueryResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStrategies", ctx, query)
	ret0, _ := ret[0].(*service_manage.BatchQueryResponse)
	return ret0
}

// GetStrategies indicates an expected call of GetStrategies.
func (mr *MockStrategyOperatorMockRecorder) GetStrategies(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStrategies", reflect.TypeOf((*MockStrategyOperator)(nil).GetStrategies), ctx, query)
}

// GetStrategy mocks base method.
func (m *MockStrategyOperator) GetStrategy(ctx context.Context, strategy *security.AuthStrategy) *service_manage.Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStrategy", ctx, strategy)
	ret0, _ := ret[0].(*service_manage.Response)
	return ret0
}

// GetStrategy indicates an expected call of GetStrategy.
func (mr *MockStrategyOperatorMockRecorder) GetStrategy(ctx, strategy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStrategy", reflect.TypeOf((*MockStrategyOperator)(nil).GetStrategy), ctx, strategy)
}

// UpdateStrategies mocks base method.
func (m *MockStrategyOperator) UpdateStrategies(ctx context.Context, reqs []*security.ModifyAuthStrategy) *service_manage.BatchWriteResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStrategies", ctx, reqs)
	ret0, _ := ret[0].(*service_manage.BatchWriteResponse)
	return ret0
}

// UpdateStrategies indicates an expected call of UpdateStrategies.
func (mr *MockStrategyOperatorMockRecorder) UpdateStrategies(ctx, reqs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStrategies", reflect.TypeOf((*MockStrategyOperator)(nil).UpdateStrategies), ctx, reqs)
}
