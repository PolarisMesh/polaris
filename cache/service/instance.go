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
	"time"

	apimodel "github.com/polarismesh/specification/source/go/api/v1/model"
	apiservice "github.com/polarismesh/specification/source/go/api/v1/service_manage"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	types "github.com/polarismesh/polaris/cache/api"
	"github.com/polarismesh/polaris/common/model"
	"github.com/polarismesh/polaris/common/utils"
	"github.com/polarismesh/polaris/store"
)

const (
	// 定时全量对账
	checkAllIntervalSec = 60
)

// instanceCache 实例缓存的类
type instanceCache struct {
	*types.BaseCache

	svcCache        *serviceCache
	storage         store.Store
	lastMtimeLogged int64
	// instanceid -> instance
	ids *utils.SyncMap[string, *model.Instance]
	// service id -> [instanceid ->instance]
	services *utils.SyncMap[string, *utils.SyncMap[string, *model.Instance]]
	// service id -> [instanceCount]
	instanceCounts     *utils.SyncMap[string, *model.InstanceCount]
	servicePortsBucket *servicePortsBucket
	disableBusiness    bool
	needMeta           bool
	systemServiceID    []string
	singleFlight       *singleflight.Group
	instanceCount      int64
	lastCheckAllTime   int64
}

// NewInstanceCache 新建一个instanceCache
func NewInstanceCache(storage store.Store, cacheMgr types.CacheManager) types.InstanceCache {
	return &instanceCache{
		BaseCache:    types.NewBaseCache(storage, cacheMgr),
		storage:      storage,
		singleFlight: new(singleflight.Group),
	}
}

// Initialize 初始化函数
func (ic *instanceCache) Initialize(opt map[string]interface{}) error {
	ic.svcCache = ic.BaseCache.CacheMgr.GetCacher(types.CacheService).(*serviceCache)
	ic.ids = utils.NewSyncMap[string, *model.Instance]()
	ic.services = utils.NewSyncMap[string, *utils.SyncMap[string, *model.Instance]]()
	ic.instanceCounts = utils.NewSyncMap[string, *model.InstanceCount]()
	ic.servicePortsBucket = newServicePortsBucket()
	if opt == nil {
		return nil
	}
	ic.disableBusiness, _ = opt["disableBusiness"].(bool)
	ic.needMeta, _ = opt["needMeta"].(bool)
	// 只加载系统服务
	if ic.disableBusiness {
		services, err := ic.getSystemServices()
		if err != nil {
			return err
		}
		ic.systemServiceID = make([]string, 0, len(services))
		for _, service := range services {
			if service.IsAlias() {
				continue
			}
			ic.systemServiceID = append(ic.systemServiceID, service.ID)
		}
	}
	ic.InitBaseOptions(opt)
	return nil
}

// Update 更新缓存函数
func (ic *instanceCache) Update() error {
	err, _ := ic.singleUpdate()
	return err
}

func (ic *instanceCache) singleUpdate() (error, bool) {
	// 多个线程竞争，只有一个线程进行更新
	_, err, shared := ic.singleFlight.Do(ic.Name(), func() (interface{}, error) {
		return nil, ic.DoCacheUpdate(ic.Name(), ic.realUpdate)
	})
	return err, shared
}

func (ic *instanceCache) LastMtime() time.Time {
	return ic.BaseCache.LastMtime(ic.Name())
}

// fetchStartTime 获取数据增量更新起始时间
func (ic *instanceCache) fetchStartTime() time.Time {
	if ic.GetFetchStartTimeType() == types.FetchFromLastMtime {
		return ic.LastMtime()
	}
	return ic.LastFetchTime()
}

func (ic *instanceCache) checkAll(tx store.Tx) {
	curTimeSec := time.Now().Unix()
	if curTimeSec-ic.lastCheckAllTime < checkAllIntervalSec {
		return
	}
	defer func() {
		ic.lastCheckAllTime = curTimeSec
	}()
	count, err := ic.storage.GetInstancesCountTx(tx)
	if err != nil {
		log.Errorf("[Cache][Instance] get instance count from storage err: %s", err.Error())
		return
	}
	if ic.instanceCount == int64(count) {
		return
	}
	log.Infof(
		"[Cache][Instance] instance count not match, expect %d, actual %d, fallback to load all",
		count, ic.instanceCount)
	ic.ResetLastMtime(ic.Name())
	ic.ResetLastFetchTime()
}

const maxLoadTimeDuration = 1 * time.Second

func (ic *instanceCache) realUpdate() (map[string]time.Time, int64, error) {
	// 拉取diff前的所有数据
	start := time.Now()
	tx, err := ic.storage.StartReadTx()
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		log.Error("[Cache][Instance] begin transaction storage read tx", zap.Error(err))
		return nil, -1, err
	}

	var instanceChangeEvents []*cacheInstanceEvent
	defer func() {
		_ = tx.Rollback()
		for i := range instanceChangeEvents {
			ic.Manager.OnEvent(instanceChangeEvents[i].item, instanceChangeEvents[i].eventType)
		}
		ic.reportMetricsInfo()
	}()

	if err := tx.CreateReadView(); err != nil {
		log.Error("[Cache][Instance] create storage snapshot read view", zap.Error(err))
		return nil, -1, err
	}

	events, lastMtimes, total, err := ic.handleUpdate(start, tx)
	_ = tx.Commit()
	instanceChangeEvents = events
	return lastMtimes, total, err
}

func (ic *instanceCache) handleUpdate(start time.Time, tx store.Tx) ([]*cacheInstanceEvent, map[string]time.Time, int64, error) {
	defer func() {
		ic.lastMtimeLogged = types.LogLastMtime(ic.lastMtimeLogged, ic.LastMtime().Unix(), "Instance")
		ic.checkAll(tx)
	}()

	instances, err := ic.storage.GetMoreInstances(tx, ic.fetchStartTime(), ic.IsFirstUpdate(),
		ic.needMeta, ic.systemServiceID)

	if err != nil {
		log.Error("[Cache][Instance] update get storage more", zap.Error(err))
		return nil, nil, -1, err
	}

	events, lastMtimes, update, del := ic.setInstances(instances)
	log.Info("[Cache][Instance] get more instances",
		zap.Int("pull-from-store", len(instances)), zap.Int("update", update), zap.Int("delete", del),
		zap.Time("last", ic.LastMtime()), zap.Duration("used", time.Since(start)))
	return events, lastMtimes, int64(len(instances)), err
}

// Clear 清理内部缓存数据
func (ic *instanceCache) Clear() error {
	ic.BaseCache.Clear()
	ic.ids = utils.NewSyncMap[string, *model.Instance]()
	ic.services = utils.NewSyncMap[string, *utils.SyncMap[string, *model.Instance]]()
	ic.instanceCounts = utils.NewSyncMap[string, *model.InstanceCount]()
	ic.servicePortsBucket.reset()
	ic.instanceCount = 0
	return nil
}

// Name 获取资源名称
func (ic *instanceCache) Name() string {
	return types.InstanceName
}

// getSystemServices 获取系统服务ID
func (ic *instanceCache) getSystemServices() ([]*model.Service, error) {
	services, err := ic.storage.GetSystemServices()
	if err != nil {
		log.Errorf("[Cache][Instance] get system services err: %s", err.Error())
		return nil, err
	}
	return services, nil
}

// setInstances 保存instance到内存中
// 返回：更新个数，删除个数
func (ic *instanceCache) setInstances(ins map[string]*model.Instance) ([]*cacheInstanceEvent, map[string]time.Time, int, int) {
	if len(ins) == 0 {
		return nil, nil, 0, 0
	}
	events := make([]*cacheInstanceEvent, 0, len(ins))
	addInstances := map[string]string{}
	updateInstances := map[string]string{}
	deleteInstances := map[string]string{}

	lastMtime := ic.LastMtime().Unix()
	update := 0
	del := 0
	affect := make(map[string]bool)
	progress := 0
	instanceCount := ic.instanceCount

	for _, item := range ins {
		progress++
		if progress%50000 == 0 {
			log.Infof("[Cache][Instance] set instances progress: %d / %d", progress, len(ins))
		}
		modifyTime := item.ModifyTime.Unix()
		if lastMtime < modifyTime {
			lastMtime = modifyTime
		}
		affect[item.ServiceID] = true
		_, itemExist := ic.ids.Load(item.ID())
		// 待删除的instance
		if !item.Valid {
			deleteInstances[item.ID()] = item.Revision()
			del++
			ic.ids.Delete(item.ID())
			if itemExist {
				events = append(events, &cacheInstanceEvent{
					item:      item,
					eventType: types.EventDeleted,
				})
				instanceCount--
			}
			value, ok := ic.services.Load(item.ServiceID)
			if !ok {
				continue
			}

			value.Delete(item.ID())
			continue
		}
		// 有修改或者新增的数据
		update++
		// 缓存的instance map增加一个version和protocol字段
		if item.Proto.Metadata == nil {
			item.Proto.Metadata = make(map[string]string)
		}

		item = fillInternalLabels(item)

		ic.ids.Store(item.ID(), item)
		if !itemExist {
			addInstances[item.ID()] = item.Revision()
			instanceCount++
			events = append(events, &cacheInstanceEvent{
				item:      item,
				eventType: types.EventCreated,
			})
		} else {
			updateInstances[item.ID()] = item.Revision()
			events = append(events, &cacheInstanceEvent{
				item:      item,
				eventType: types.EventUpdated,
			})
		}
		value, ok := ic.services.Load(item.ServiceID)
		if !ok {
			value = utils.NewSyncMap[string, *model.Instance]()
			ic.services.Store(item.ServiceID, value)
		}

		ic.servicePortsBucket.appendPort(item.ServiceID, item.Protocol(), item.Port())
		value.Store(item.ID(), item)
	}

	if ic.instanceCount != instanceCount {
		log.Infof("[Cache][Instance] instance count update from %d to %d",
			ic.instanceCount, instanceCount)
		ic.instanceCount = instanceCount
	}

	log.Info("[Cache][Instance] instances change info", zap.Any("add", addInstances),
		zap.Any("update", updateInstances), zap.Any("delete", deleteInstances))

	ic.postProcessUpdatedServices(affect)
	ic.svcCache.notifyServiceCountReload(affect)
	return events, map[string]time.Time{
		ic.Name(): time.Unix(lastMtime, 0),
	}, update, del
}

func fillInternalLabels(item *model.Instance) *model.Instance {
	if len(item.Version()) > 0 {
		item.Proto.Metadata["version"] = item.Version()
	}
	if len(item.Protocol()) > 0 {
		item.Proto.Metadata["protocol"] = item.Protocol()
	}

	if item.Location() != nil {
		item.Proto.Metadata["region"] = item.Location().GetRegion().GetValue()
		item.Proto.Metadata["zone"] = item.Location().GetZone().GetValue()
		item.Proto.Metadata["campus"] = item.Location().GetCampus().GetValue()
	}
	return item
}

func (ic *instanceCache) postProcessUpdatedServices(affect map[string]bool) {
	progress := 0
	for serviceID := range affect {
		ic.svcCache.GetRevisionWorker().Notify(serviceID, true)
		progress++
		if progress%10000 == 0 {
			log.Infof("[Cache][Instance] revision notify progress(%d / %d)", progress, len(affect))
		}
		// 构建服务数量统计
		value, ok := ic.services.Load(serviceID)
		if !ok {
			ic.instanceCounts.Delete(serviceID)
			continue
		}
		count := &model.InstanceCount{}
		value.Range(func(key string, instance *model.Instance) bool {
			count.TotalInstanceCount++
			if isInstanceHealthy(instance) {
				count.HealthyInstanceCount++
			}
			if instance.Proto.GetIsolate().GetValue() {
				count.IsolateInstanceCount++
			}
			return true
		})
		if count.TotalInstanceCount == 0 {
			ic.instanceCounts.Delete(serviceID)
			continue
		}
		ic.instanceCounts.Store(serviceID, count)
	}
}

func isInstanceHealthy(instance *model.Instance) bool {
	return instance.Proto.GetHealthy().GetValue() && !instance.Proto.GetIsolate().GetValue()
}

// GetInstance 根据实例ID获取实例数据
func (ic *instanceCache) GetInstance(instanceID string) *model.Instance {
	if instanceID == "" {
		return nil
	}

	value, ok := ic.ids.Load(instanceID)
	if !ok {
		return nil
	}

	return value
}

// GetInstancesByServiceID 根据ServiceID获取实例数据
func (ic *instanceCache) GetInstancesByServiceID(serviceID string) []*model.Instance {
	if serviceID == "" {
		return nil
	}

	value, ok := ic.services.Load(serviceID)
	if !ok {
		return nil
	}

	var out []*model.Instance
	value.Range(func(k string, v *model.Instance) bool {
		out = append(out, v)
		return true
	})

	return out
}

// GetInstancesCountByServiceID 根据服务ID获取实例数
func (ic *instanceCache) GetInstancesCountByServiceID(serviceID string) model.InstanceCount {
	if serviceID == "" {
		return model.InstanceCount{}
	}

	value, ok := ic.instanceCounts.Load(serviceID)
	if !ok {
		return model.InstanceCount{}
	}
	return *value
}

// IteratorInstances 迭代所有的instance的函数
func (ic *instanceCache) IteratorInstances(iterProc types.InstanceIterProc) error {
	return iteratorInstancesProc(ic.ids, iterProc)
}

// IteratorInstancesWithService 根据服务ID进行迭代回调
func (ic *instanceCache) IteratorInstancesWithService(serviceID string, iterProc types.InstanceIterProc) error {
	if serviceID == "" {
		return nil
	}
	value, ok := ic.services.Load(serviceID)
	if !ok {
		return nil
	}

	return iteratorInstancesProc(value, iterProc)
}

// GetInstancesCount 获取实例的个数
func (ic *instanceCache) GetInstancesCount() int {
	count := 0
	ic.ids.Range(func(key string, value *model.Instance) bool {
		count++
		return true
	})

	return count
}

// GetInstanceLabels 获取某个服务下实例的所有标签信息集合
func (ic *instanceCache) GetInstanceLabels(serviceID string) *apiservice.InstanceLabels {
	if serviceID == "" {
		return &apiservice.InstanceLabels{}
	}

	value, ok := ic.services.Load(serviceID)
	if !ok {
		return &apiservice.InstanceLabels{}
	}

	ret := &apiservice.InstanceLabels{
		Labels: make(map[string]*apimodel.StringList),
	}

	tmp := make(map[string]map[string]struct{})
	_ = iteratorInstancesProc(value, func(key string, value *model.Instance) (bool, error) {
		metadata := value.Metadata()
		for k, v := range metadata {
			if _, ok := tmp[k]; !ok {
				tmp[k] = make(map[string]struct{})
			}
			tmp[k][v] = struct{}{}
		}
		return true, nil
	})

	for k, v := range tmp {
		if _, ok := ret.Labels[k]; !ok {
			ret.Labels[k] = &apimodel.StringList{Values: make([]string, 0, 4)}
		}

		for vv := range v {
			ret.Labels[k].Values = append(ret.Labels[k].Values, vv)
		}
	}

	return ret
}

func (ic *instanceCache) GetServicePorts(serviceID string) []*model.ServicePort {
	return ic.servicePortsBucket.listPort(serviceID)
}

// iteratorInstancesProc 迭代指定的instance数据，id->instance
func iteratorInstancesProc(data *utils.SyncMap[string, *model.Instance], iterProc types.InstanceIterProc) error {
	var (
		cont bool
		err  error
	)

	proc := func(k string, v *model.Instance) bool {
		cont, err = iterProc(k, v)
		if err != nil {
			return false
		}
		return cont
	}

	data.Range(proc)
	return err
}

type cacheInstanceEvent struct {
	item      *model.Instance
	eventType types.EventType
}
