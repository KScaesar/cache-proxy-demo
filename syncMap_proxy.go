package cache_proxy_demo

import (
	"sync"
)

func UseSyncMap(baseProxy *BaseCacheProxy) CacheProxy {
	return &SyncMapProxy{
		transform: baseProxy.Transform,
		cache:     baseProxy.Cache,

		baseProxy: baseProxy,
	}
}

type SyncMapProxy struct {
	transform TransformQryOptionToCacheKey
	cache     Cache

	// Map.LoadOrStore(key, value any) (actual any, loaded bool)
	shards sync.Map

	baseProxy *BaseCacheProxy
}

func (proxy *SyncMapProxy) Execute(qryOption any, readModelType any) (readModel any, err error) {
	return proxy.execute1(qryOption, readModelType)
}

func (proxy *SyncMapProxy) execute1(qryOption any, readModelType any) (readModel any, err error) {
	key := proxy.transform(qryOption)

	var wg sync.WaitGroup
	wg.Add(1)
	getReadModelFunc := func() (any, error) {
		wg.Wait()
		return readModel, err
	}

	fn, isSecond := proxy.shards.LoadOrStore(key, getReadModelFunc)
	if isSecond {
		// 其他 goroutine 拿到的是, first goroutine 的 閉包 func
		return fn.(func() (any, error))()
	}

	defer func() {
		wg.Done()
		proxy.shards.Delete(key)
	}()

	return proxy.baseProxy.Execute(qryOption, readModelType)
}

func (proxy *SyncMapProxy) execute3(qryOption any, readModelType any) (readModel any, err error) {
	key := proxy.transform(qryOption)
	val, err := proxy.cache.GetValue(key, readModelType)
	if err == nil {
		return val, nil
	}

	return proxy.execute1(qryOption, readModelType)
}
