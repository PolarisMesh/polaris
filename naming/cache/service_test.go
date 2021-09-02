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

package cache

import (
	"fmt"
	"github.com/polarismesh/polaris-server/common/model"
	"github.com/polarismesh/polaris-server/store/mock"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

// 生成一个测试的serviceCache和对应的mock对象
func newTestServiceCache(t *testing.T) (*gomock.Controller, *mock.MockStore, *serviceCache) {
	ctl := gomock.NewController(t)

	storage := mock.NewMockStore(ctl)
	sc := newServiceCache(storage, make(chan *revisionNotify, 1024))
	opt := map[string]interface{}{
		"disableBusiness": false,
		"needMeta":        true,
	}
	_ = sc.initialize(opt)

	return ctl, storage, sc
}

// 获取当前缓存中的services总数
func getServiceCacheCount(sc *serviceCache) int {
	sum := 0
	_ = sc.IteratorServices(func(key string, value *model.Service) (bool, error) {
		sum++
		return true, nil
	})
	return sum
}

// 生成一些测试的services
func genModelService(total int) map[string]*model.Service {
	out := make(map[string]*model.Service)
	for i := 0; i < total; i++ {
		item := &model.Service{
			ID:         fmt.Sprintf("ID-%d", i),
			Namespace:  fmt.Sprintf("Namespace-%d", i),
			Name:       fmt.Sprintf("Name-%d", i),
			Valid:      true,
			ModifyTime: time.Unix(int64(i), 0),
		}
		out[item.ID] = item
	}

	return out
}

// 测试缓存更新函数
func TestServiceUpdate(t *testing.T) {
	ctl, storage, sc := newTestServiceCache(t)
	defer ctl.Finish()

	t.Run("所有数据为空，可以正常获取数据", func(t *testing.T) {
		gomock.InOrder(
			storage.EXPECT().
				GetMoreServices(sc.lastMtime.Add(DefaultTimeDiff), sc.firstUpdate, sc.disableBusiness, sc.needMeta).
				Return(nil, nil),
		)

		if err := sc.update(); err != nil {
			t.Fatalf("error: %s", err.Error())
		}

		if sum := getServiceCacheCount(sc); sum != 0 {
			t.Fatalf("error: %d", sum)
		}
	})
	t.Run("有数据更新，数据正常", func(t *testing.T) {
		_ = sc.clear()
		services := genModelService(100)
		gomock.InOrder(
			storage.EXPECT().GetMoreServices(sc.lastMtime.Add(DefaultTimeDiff), sc.firstUpdate, sc.disableBusiness, sc.needMeta).
				Return(services, nil),
		)

		if err := sc.update(); err != nil {
			t.Fatalf("error: %s", err.Error())
		}

		if sum := getServiceCacheCount(sc); sum != 100 {
			t.Fatalf("error: %d", sum)
		}
	})
	t.Run("有数据更新，重复更新，数据更新正常", func(t *testing.T) {
		_ = sc.clear()
		services1 := genModelService(100)
		services2 := genModelService(300)
		gomock.InOrder(
			storage.EXPECT().GetMoreServices(sc.lastMtime.Add(DefaultTimeDiff), sc.firstUpdate, sc.disableBusiness, sc.needMeta).
				Return(services1, nil),
		)

		if err := sc.update(); err != nil {
			t.Fatalf("error: %s", err.Error())
		}

		gomock.InOrder(
			storage.EXPECT().GetMoreServices(sc.lastMtime.Add(DefaultTimeDiff), sc.firstUpdate, sc.disableBusiness, sc.needMeta).
				Return(services2, nil),
		)
		_ = sc.update()
		if sum := getServiceCacheCount(sc); sum != 300 {
			t.Fatalf("error: %d", sum)
		}
	})
}

// 测试缓存更新函数1
func TestServiceUpdate1(t *testing.T) {
	ctl, storage, sc := newTestServiceCache(t)
	defer ctl.Finish()

	t.Run("服务全部被删除，会被清除掉", func(t *testing.T) {
		_ = sc.clear()
		services := genModelService(100)
		gomock.InOrder(storage.EXPECT().
			GetMoreServices(sc.lastMtime.Add(DefaultTimeDiff), sc.firstUpdate, sc.disableBusiness, sc.needMeta).Return(services, nil))
		_ = sc.update()

		// 把所有的都置为false
		for _, service := range services {
			service.Valid = false
		}

		gomock.InOrder(storage.EXPECT().
			GetMoreServices(sc.lastMtime.Add(DefaultTimeDiff), sc.firstUpdate, sc.disableBusiness, sc.needMeta).Return(services, nil))
		_ = sc.update()

		if sum := getServiceCacheCount(sc); sum != 0 {
			t.Fatalf("error: %d", sum)
		}
	})

	t.Run("服务部分被删除，缓存内容正常", func(t *testing.T) {
		_ = sc.clear()
		services := genModelService(100)
		gomock.InOrder(storage.EXPECT().
			GetMoreServices(sc.lastMtime.Add(DefaultTimeDiff), sc.firstUpdate, sc.disableBusiness, sc.needMeta).Return(services, nil))
		_ = sc.update()

		// 把所有的都置为false
		idx := 0
		for _, service := range services {
			if idx%2 == 0 {
				service.Valid = false
			}
			idx++
		}

		gomock.InOrder(storage.EXPECT().
			GetMoreServices(sc.lastMtime.Add(DefaultTimeDiff), sc.firstUpdate, sc.disableBusiness, sc.needMeta).Return(services, nil))
		_ = sc.update()

		if sum := getServiceCacheCount(sc); sum != 50 { // remain half
			t.Fatalf("error: %d", sum)
		}
	})
}

// 测试缓存更新
func TestServiceUpdate2(t *testing.T) {
	ctl, storage, sc := newTestServiceCache(t)
	defer ctl.Finish()

	t.Run("store返回失败，update会返回失败", func(t *testing.T) {
		_ = sc.clear()
		gomock.InOrder(
			storage.EXPECT().GetMoreServices(sc.lastMtime.Add(DefaultTimeDiff), sc.firstUpdate, sc.disableBusiness, sc.needMeta).
				Return(nil, fmt.Errorf("store error")),
		)

		if err := sc.update(); err != nil {
			t.Logf("pass: %s", err.Error())
		} else {
			t.Fatalf("error")
		}
	})
}

// 根据服务名获取服务缓存信息
func TestGetServiceByName(t *testing.T) {
	ctl, _, sc := newTestServiceCache(t)
	defer ctl.Finish()
	t.Run("可以根据服务名和命名空间，正常获取缓存服务信息", func(t *testing.T) {
		_ = sc.clear()
		services := genModelService(20)
		sc.setServices(services)

		for _, entry := range services {
			service := sc.GetServiceByName(entry.Name, entry.Namespace)
			if service == nil {
				t.Fatalf("error")
			}
		}
	})
	t.Run("服务不存在，返回为空", func(t *testing.T) {
		_ = sc.clear()
		services := genModelService(20)
		sc.setServices(services)
		if service := sc.GetServiceByName("aaa", "bbb"); service != nil {
			t.Fatalf("error")
		}
	})
}

// 根据服务ID获取服务缓存信息
func TestServiceCache_GetServiceByID(t *testing.T) {
	ctl, _, sc := newTestServiceCache(t)
	defer ctl.Finish()

	t.Run("可以根据服务ID，正常获取缓存的服务信息", func(t *testing.T) {
		_ = sc.clear()
		services := genModelService(30)
		sc.setServices(services)

		for _, entry := range services {
			service := sc.GetServiceByID(entry.ID)
			if service == nil {
				t.Fatalf("error")
			}
		}
	})

	t.Run("缓存内容为空，根据ID获取数据，会返回为空", func(t *testing.T) {
		_ = sc.clear()
		services := genModelService(30)
		sc.setServices(services)

		if service := sc.GetServiceByID("123456789"); service != nil {
			t.Fatalf("error")
		}
	})
}
