package cache_proxy_demo

import (
	"golang.org/x/sync/singleflight"
)

func UseSingleflight(baseProxy *BaseCacheProxy) CacheProxy {
	return &SingleflightProxy{
		transform: baseProxy.Transform,
		cache:     baseProxy.Cache,

		baseProxy: baseProxy,
	}
}

type SingleflightProxy struct {
	transform TransformQryOptionToCacheKey
	cache     Cache

	// Group.Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)
	workflow singleflight.Group

	baseProxy *BaseCacheProxy
}

func (proxy *SingleflightProxy) Execute(qryOption any, readModelType any) (readModel any, err error) {
	return proxy.execute1(qryOption, readModelType)
}

func (proxy *SingleflightProxy) execute1(qryOption any, readModelType any) (readModel any, err error) {
	key := proxy.transform(qryOption)
	readModel, err, _ = proxy.workflow.Do(key, func() (interface{}, error) {
		return proxy.baseProxy.Execute(qryOption, readModelType)
	})
	return
}

func (proxy *SingleflightProxy) execute3(qryOption any, readModelType any) (readModel any, err error) {
	key := proxy.transform(qryOption)
	val, err := proxy.cache.GetValue(key, readModelType)
	if err == nil {
		return val, nil
	}

	return proxy.execute1(qryOption, readModelType)
}
