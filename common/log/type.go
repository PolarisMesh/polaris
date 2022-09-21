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

package log

// logger type
const (
	// NamingLoggerName naming logger name, can use FindScope function to get the logger
	NamingLoggerName = "naming"
	// HealthCheckLoggerName heartbeat logger name, can use FindScope function to get the logger
	HealthCheckLoggerName = "healthcheck"
	// ConfigLoggerName config logger name, can use FindScope function to get the logger
	ConfigLoggerName = "config"
	// CacheLoggerName cache logger name, can use FindScope function to get the logger
	CacheLoggerName = "cache"
	// AuthLoggerName auth logger name, can use FindScope function to get the logger
	AuthLoggerName = "auth"
	// StoreLoggerName store logger name, can use FindScope function to get the logger
	StoreLoggerName = "store"
	// APIServerLoggerName store logger name, can use FindScope function to get the logger
	APIServerLoggerName = "apiserver"
)

var (
	namingScope      = RegisterScope(NamingLoggerName, "naming logging messages.", 0)
	healthcheckScope = RegisterScope(HealthCheckLoggerName, "healthcheck logging messages.", 0)
	configScope      = RegisterScope(ConfigLoggerName, "config logging messages.", 0)
	cacheScope       = RegisterScope(CacheLoggerName, "cache logging messages.", 0)
	authScope        = RegisterScope(AuthLoggerName, "auth logging messages.", 0)
	storeScope       = RegisterScope(StoreLoggerName, "store logging messages.", 0)
	apiserverScope   = RegisterScope(APIServerLoggerName, "apiserver logging messages.", 0)
)

func allLoggerTypes() []string {
	return []string{NamingLoggerName, ConfigLoggerName, CacheLoggerName,
		AuthLoggerName, StoreLoggerName, APIServerLoggerName, HealthCheckLoggerName, DefaultLoggerName}
}

// APIServerScope default logging scope handler
func APIServerScope() *Scope {
	return defaultScope
}

// DefaultScope default logging scope handler
func DefaultScope() *Scope {
	return apiserverScope
}

// NamingScope naming logging scope handler
func NamingScope() *Scope {
	return namingScope
}

// HealthCheckScope heartbeat logging scope handler
func HealthCheckScope() *Scope {
	return healthcheckScope
}

// ConfigScope config logging scope handler
func ConfigScope() *Scope {
	return configScope
}

// CacheScope cache logging scope handler
func CacheScope() *Scope {
	return cacheScope
}

// AuthScope auth logging scope handler
func AuthScope() *Scope {
	return authScope
}

// StoreScope store logging scope handler
func StoreScope() *Scope {
	return storeScope
}
