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
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	api "github.com/polarismesh/polaris-server/common/api/v1"
	"github.com/polarismesh/polaris-server/common/utils"
)

// 检查routingConfig前后是否一致
func checkSameRoutingConfig(t *testing.T, lhs *api.Routing, rhs *api.Routing) {
	if lhs.GetService().GetValue() != rhs.GetService().GetValue() ||
		lhs.GetNamespace().GetValue() != rhs.GetNamespace().GetValue() {
		t.Fatalf("error: (%s), (%s)", lhs, rhs)
	}

	checkFunc := func(in []*api.Route, out []*api.Route) bool {
		if in == nil && out == nil {
			return true
		}

		if in == nil || out == nil {
			t.Fatalf("error: empty")
			return false
		}

		if len(in) != len(out) {
			t.Fatalf("error: %d, %d", len(in), len(out))
			return false
		}

		inStr, err := json.Marshal(in)
		if err != nil {
			t.Fatalf("error: %s", err.Error())
			return false
		}

		outStr, err := json.Marshal(out)
		if err != nil {
			t.Fatalf("error: %s", err.Error())
			return false
		}

		if string(inStr) != string(outStr) {
			t.Fatalf("error: (%s), (%s)", string(inStr), string(outStr))
			return false
		}

		return true
	}

	checkFunc(lhs.Inbounds, rhs.Inbounds)
	checkFunc(lhs.Outbounds, rhs.Outbounds)
}

// 测试创建路由配置
func TestCreateRoutingConfig(t *testing.T) {

	discoverSuit := &DiscoverTestSuit{}
	if err := discoverSuit.initialize(); err != nil {
		t.Fatal(err)
	}
	defer discoverSuit.Destroy()

	Convey("正常创建路由配置配置请求", t, func() {
		_, serviceResp := discoverSuit.createCommonService(t, 200)
		defer discoverSuit.cleanServiceName(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())
		req, _ := discoverSuit.createCommonRoutingConfig(t, serviceResp, 3, 0)

		// 对写进去的数据进行查询
		time.Sleep(discoverSuit.updateCacheInterval)
		out := discoverSuit.server.GetRoutingConfigWithCache(discoverSuit.defaultCtx, serviceResp)
		defer discoverSuit.cleanCommonRoutingConfig(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())
		if !respSuccess(out) {
			t.Fatalf("error: %+v", out)
		}
		checkSameRoutingConfig(t, req, out.GetRouting())
	})

	Convey("同一个服务重复创建路由配置，报错", t, func() {
		_, serviceResp := discoverSuit.createCommonService(t, 10)
		defer discoverSuit.cleanServiceName(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())

		req, _ := discoverSuit.createCommonRoutingConfig(t, serviceResp, 1, 0)
		defer discoverSuit.cleanCommonRoutingConfig(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())

		resp := discoverSuit.server.CreateRoutingConfigs(discoverSuit.defaultCtx, []*api.Routing{req})
		So(respSuccess(resp), ShouldEqual, false)
		t.Logf("%s", resp.GetInfo().GetValue())
	})
}

// 测试创建路由配置
func TestCreateRoutingConfig2(t *testing.T) {

	discoverSuit := &DiscoverTestSuit{}
	if err := discoverSuit.initialize(); err != nil {
		t.Fatal(err)
	}
	defer discoverSuit.Destroy()

	Convey("参数缺失，报错", t, func() {
		_, serviceResp := discoverSuit.createCommonService(t, 20)
		defer discoverSuit.cleanServiceName(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())

		req := &api.Routing{}
		resp := discoverSuit.server.CreateRoutingConfigs(discoverSuit.defaultCtx, []*api.Routing{req})
		So(respSuccess(resp), ShouldEqual, false)
		t.Logf("%s", resp.GetInfo().GetValue())

		req.Service = serviceResp.Name
		resp = discoverSuit.server.CreateRoutingConfigs(discoverSuit.defaultCtx, []*api.Routing{req})
		So(respSuccess(resp), ShouldEqual, false)
		t.Logf("%s", resp.GetInfo().GetValue())

		req.Namespace = serviceResp.Namespace
		resp = discoverSuit.server.CreateRoutingConfigs(discoverSuit.defaultCtx, []*api.Routing{req})
		defer discoverSuit.cleanCommonRoutingConfig(req.GetService().GetValue(), req.GetNamespace().GetValue())
		So(respSuccess(resp), ShouldEqual, true)
		t.Logf("%s", resp.GetInfo().GetValue())
	})

	Convey("服务不存在，创建路由配置，报错", t, func() {
		_, serviceResp := discoverSuit.createCommonService(t, 120)
		discoverSuit.cleanServiceName(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())

		req := &api.Routing{}
		req.Service = serviceResp.Name
		req.Namespace = serviceResp.Namespace
		req.ServiceToken = serviceResp.Token
		resp := discoverSuit.server.CreateRoutingConfigs(discoverSuit.defaultCtx, []*api.Routing{req})
		So(respSuccess(resp), ShouldEqual, false)
		t.Logf("%s", resp.GetInfo().GetValue())
	})
}

// 测试删除路由配置
func TestDeleteRoutingConfig(t *testing.T) {

	discoverSuit := &DiscoverTestSuit{}
	if err := discoverSuit.initialize(); err != nil {
		t.Fatal(err)
	}
	defer discoverSuit.Destroy()

	Convey("可以正常删除路由配置", t, func() {
		_, serviceResp := discoverSuit.createCommonService(t, 100)
		defer discoverSuit.cleanServiceName(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())

		_, resp := discoverSuit.createCommonRoutingConfig(t, serviceResp, 3, 0)
		discoverSuit.deleteCommonRoutingConfig(t, resp)
		defer discoverSuit.cleanCommonRoutingConfig(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())

		// 删除之后，数据不见
		time.Sleep(discoverSuit.updateCacheInterval)
		out := discoverSuit.server.GetRoutingConfigWithCache(discoverSuit.defaultCtx, serviceResp)
		So(out.GetRouting(), ShouldBeNil)
	})
}

// 测试更新路由配置
func TestUpdateRoutingConfig(t *testing.T) {

	discoverSuit := &DiscoverTestSuit{}
	if err := discoverSuit.initialize(); err != nil {
		t.Fatal(err)
	}
	defer discoverSuit.Destroy()

	Convey("可以正常更新路由配置", t, func() {
		_, serviceResp := discoverSuit.createCommonService(t, 50)
		serviceName := serviceResp.GetName().GetValue()
		namespace := serviceResp.GetNamespace().GetValue()
		defer discoverSuit.cleanServiceName(serviceName, namespace)

		_, resp := discoverSuit.createCommonRoutingConfig(t, serviceResp, 2, 0)
		defer discoverSuit.cleanCommonRoutingConfig(serviceName, namespace)
		resp.ServiceToken = utils.NewStringValue(serviceResp.GetToken().GetValue())

		resp.Outbounds = resp.Inbounds
		resp.Inbounds = make([]*api.Route, 0)
		discoverSuit.updateCommonRoutingConfig(t, resp)

		time.Sleep(discoverSuit.updateCacheInterval)
		out := discoverSuit.server.GetRoutingConfigWithCache(discoverSuit.defaultCtx, serviceResp)
		checkSameRoutingConfig(t, resp, out.GetRouting())
	})
}

// 测试缓存获取路由配置
func TestGetRoutingConfigWithCache(t *testing.T) {

	discoverSuit := &DiscoverTestSuit{}
	if err := discoverSuit.initialize(); err != nil {
		t.Fatal(err)
	}
	defer discoverSuit.Destroy()

	Convey("多个服务的，多个路由配置，都可以查询到", t, func() {
		total := 20
		serviceResps := make([]*api.Service, 0, total)
		routingResps := make([]*api.Routing, 0, total)
		for i := 0; i < total; i++ {
			_, resp := discoverSuit.createCommonService(t, i)
			defer discoverSuit.cleanServiceName(resp.GetName().GetValue(), resp.GetNamespace().GetValue())
			serviceResps = append(serviceResps, resp)

			_, routingResp := discoverSuit.createCommonRoutingConfig(t, resp, 2, 0)
			defer discoverSuit.cleanCommonRoutingConfig(resp.GetName().GetValue(), resp.GetNamespace().GetValue())
			routingResps = append(routingResps, routingResp)
		}

		time.Sleep(discoverSuit.updateCacheInterval)
		for i := 0; i < total; i++ {
			out := discoverSuit.server.GetRoutingConfigWithCache(discoverSuit.defaultCtx, serviceResps[i])
			checkSameRoutingConfig(t, routingResps[i], out.GetRouting())
		}
	})
	Convey("服务路由数据不改变，传递了路由revision，不返回数据", t, func() {
		_, serviceResp := discoverSuit.createCommonService(t, 10)
		defer discoverSuit.cleanServiceName(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())

		_, routingResp := discoverSuit.createCommonRoutingConfig(t, serviceResp, 2, 0)
		defer discoverSuit.cleanCommonRoutingConfig(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())

		time.Sleep(discoverSuit.updateCacheInterval)
		firstResp := discoverSuit.server.GetRoutingConfigWithCache(discoverSuit.defaultCtx, serviceResp)
		checkSameRoutingConfig(t, routingResp, firstResp.GetRouting())

		serviceResp.Revision = firstResp.Service.Revision
		secondResp := discoverSuit.server.GetRoutingConfigWithCache(discoverSuit.defaultCtx, serviceResp)
		if secondResp.GetService().GetRevision().GetValue() != serviceResp.GetRevision().GetValue() {
			t.Fatalf("error")
		}
		if secondResp.GetRouting() != nil {
			t.Fatalf("error: %+v", secondResp.GetRouting())
		}
		t.Logf("%+v", secondResp)
	})
	Convey("路由不存在，不会出异常", t, func() {
		_, serviceResp := discoverSuit.createCommonService(t, 10)
		defer discoverSuit.cleanServiceName(serviceResp.GetName().GetValue(), serviceResp.GetNamespace().GetValue())
		time.Sleep(discoverSuit.updateCacheInterval)
		if resp := discoverSuit.server.GetRoutingConfigWithCache(discoverSuit.defaultCtx, serviceResp); !respSuccess(resp) {
			t.Fatalf("error: %s", resp.GetInfo().GetValue())
		}
	})
}

// 测试直接从数据库读取路由配置数据
func TestGetRoutings(t *testing.T) {

	discoverSuit := &DiscoverTestSuit{}
	if err := discoverSuit.initialize(); err != nil {
		t.Fatal(err)
	}
	defer discoverSuit.Destroy()

	Convey("直接从数据库查询数据，可以查询成功", t, func() {
		total := 10
		var serviceResp *api.Service
		for i := 0; i < total; i++ {
			tmp, resp := discoverSuit.createCommonService(t, i)
			serviceResp = tmp
			defer discoverSuit.cleanServiceName(resp.GetName().GetValue(), resp.GetNamespace().GetValue())

			discoverSuit.createCommonRoutingConfig(t, resp, 2, 0)
			defer discoverSuit.cleanCommonRoutingConfig(resp.GetName().GetValue(), resp.GetNamespace().GetValue())
		}

		resp := discoverSuit.server.GetRoutingConfigs(discoverSuit.defaultCtx, nil)
		So(api.CalcCode(resp), ShouldEqual, 200)
		So(len(resp.GetRoutings()), ShouldBeGreaterThanOrEqualTo, total)

		resp = discoverSuit.server.GetRoutingConfigs(discoverSuit.defaultCtx, map[string]string{"limit": "5"})
		So(api.CalcCode(resp), ShouldEqual, 200)
		So(len(resp.GetRoutings()), ShouldEqual, 5)

		resp = discoverSuit.server.GetRoutingConfigs(discoverSuit.defaultCtx, map[string]string{"namespace": DefaultNamespace})
		So(api.CalcCode(resp), ShouldEqual, 200)
		So(len(resp.GetRoutings()), ShouldBeGreaterThanOrEqualTo, total)

		// 按命名空间和名字过滤，得到一个
		filter := map[string]string{
			"namespace": DefaultNamespace,
			"service":   serviceResp.GetName().GetValue(),
		}
		resp = discoverSuit.server.GetRoutingConfigs(discoverSuit.defaultCtx, filter)
		So(api.CalcCode(resp), ShouldEqual, 200)
		So(len(resp.GetRoutings()), ShouldEqual, 1)
	})
}

// test对routing字段进行校验
func TestCheckRoutingFieldLen(t *testing.T) {

	discoverSuit := &DiscoverTestSuit{}
	if err := discoverSuit.initialize(); err != nil {
		t.Fatal(err)
	}
	defer discoverSuit.Destroy()

	req := &api.Routing{
		ServiceToken: utils.NewStringValue("test"),
		Service:      utils.NewStringValue("test"),
		Namespace:    utils.NewStringValue("default"),
	}

	t.Run("创建路由规则，服务名超长", func(t *testing.T) {
		str := genSpecialStr(129)
		oldName := req.Service
		req.Service = utils.NewStringValue(str)
		resp := discoverSuit.server.CreateRoutingConfigs(discoverSuit.defaultCtx, []*api.Routing{req})
		req.Service = oldName
		if resp.Code.Value != api.InvalidServiceName {
			t.Fatalf("%+v", resp)
		}
	})
	t.Run("创建路由规则，命名空间超长", func(t *testing.T) {
		str := genSpecialStr(129)
		oldNamespace := req.Namespace
		req.Namespace = utils.NewStringValue(str)
		resp := discoverSuit.server.CreateRoutingConfigs(discoverSuit.defaultCtx, []*api.Routing{req})
		req.Namespace = oldNamespace
		if resp.Code.Value != api.InvalidNamespaceName {
			t.Fatalf("%+v", resp)
		}
	})
	t.Run("创建路由规则，toeken超长", func(t *testing.T) {
		str := genSpecialStr(2049)
		oldServiceToken := req.ServiceToken
		req.ServiceToken = utils.NewStringValue(str)
		resp := discoverSuit.server.CreateRoutingConfigs(discoverSuit.defaultCtx, []*api.Routing{req})
		req.ServiceToken = oldServiceToken
		if resp.Code.Value != api.InvalidServiceToken {
			t.Fatalf("%+v", resp)
		}
	})
}
