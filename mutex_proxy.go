package cache_proxy_demo

import "sync"

func UseMutex(baseProxy *BaseCacheProxy) CacheProxy {
	return &MutexProxy{
		transform: baseProxy.Transform,
		cache:     baseProxy.Cache,

		baseProxy: baseProxy,
	}
}

type MutexProxy struct {
	transform TransformQryOptionToCacheKey
	cache     Cache

	mu sync.Mutex

	baseProxy *BaseCacheProxy
}

func (proxy *MutexProxy) Execute(qryOption any, readModelType any) (readModel any, err error) {
	return proxy.execute1(qryOption, readModelType)
}

func (proxy *MutexProxy) execute1(qryOption any, readModelType any) (readModel any, err error) {
	proxy.mu.Lock()
	defer proxy.mu.Unlock()
	return proxy.baseProxy.Execute(qryOption, readModelType)
}

func (proxy *MutexProxy) execute3(qryOption any, readModelType any) (readModel any, err error) {
	key := proxy.transform(qryOption)
	val, err := proxy.cache.GetValue(key, readModelType)
	if err == nil {
		return val, nil
	}

	return proxy.execute1(qryOption, readModelType)
}
