package example

import (
	"context"
	"sync"
	"time"
)

type ReadDataSource func(ctx context.Context, readDtoOption any) (readModel any, err error)

type CacheProxySyncMap struct {
	Cache         Cache
	firstDelivery sync.Map // key:ReadDataSource(func)

	TransformReadOption func(readDtoOption any) (key string)
	ReadDataSource      func(ctx context.Context, readDtoOption any) (readModel any, err error)

	IsAnNotFoundError                func(err error) bool
	CanIgnoreCacheError              bool
	CanIgnoreReadSourceErrorNotFound bool // source not found, 是否交給 caller 處理
	CacheTTL                         time.Duration
}

func (proxy *CacheProxySyncMap) ReadValue(ctx context.Context, readDtoOption any) (readModel any, err error) {
	return proxy.ReadValueV3(ctx, readDtoOption)
}

// v1: userCount = int(2e4)
//	BenchmarkCacheProxySyncMap_ReadValue
//    cache_proxy_impl_test.go:210: db qry count = 1, b.N=1
//    cache_proxy_impl_test.go:210: db qry count = 100, b.N=100
//    cache_proxy_impl_test.go:210: db qry count = 4345, b.N=4345
//    cache_proxy_impl_test.go:210: db qry count = 8176, b.N=8176
// BenchmarkCacheProxySyncMap_ReadValue-8   	    8176	    147555 ns/op

func (proxy *CacheProxySyncMap) ReadValueV1(ctx context.Context, readDtoOption any) (readModel any, err error) {
	var empty any

	key := proxy.TransformReadOption(readDtoOption)
	readFn, exist := proxy.firstDelivery.Load(key)
	if exist {
		return readFn.(ReadDataSource)(ctx, readDtoOption)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	readFn = ReadDataSource(func(context.Context, any) (any, error) {
		wg.Wait()
		return readModel, err
	})

	mainReadFn, ok := proxy.firstDelivery.LoadOrStore(key, readFn)
	if ok {
		// 其他 thread 拿到的是, main read 的 閉包 func, 包含回傳值
		return mainReadFn.(ReadDataSource)(ctx, readDtoOption)
	}

	// main read
	defer func() {
		wg.Done()
		proxy.firstDelivery.Delete(key)
	}()

	val, err := proxy.Cache.GetValue(ctx, key)
	if err != nil {
		if !(proxy.IsAnNotFoundError(err) || proxy.CanIgnoreCacheError) {
			return empty, err
		}

		readModel, err = proxy.ReadDataSource(ctx, readDtoOption)
		if err != nil {
			if !(proxy.IsAnNotFoundError(err) && proxy.CanIgnoreReadSourceErrorNotFound) {
				return empty, err
			}
		}

		err = proxy.Cache.PutValue(ctx, key, readModel, proxy.CacheTTL)
		if err != nil && !proxy.CanIgnoreCacheError {
			return empty, err
		}

		return readModel, nil
	}

	return val, nil
}

// v2: userCount = int(2e4)
// BenchmarkCacheProxySyncMap_ReadValue
//    cache_proxy_impl_test.go:210: db qry count = 1, b.N=1
//    cache_proxy_impl_test.go:210: db qry count = 100, b.N=100
//    cache_proxy_impl_test.go:210: db qry count = 4887, b.N=4887
//    cache_proxy_impl_test.go:210: db qry count = 12674, b.N=12674
// BenchmarkCacheProxySyncMap_ReadValue-8   	   12674	     79351 ns/op

func (proxy *CacheProxySyncMap) ReadValueV2(ctx context.Context, readDtoOption any) (readModel any, err error) {
	var empty any
	key := proxy.TransformReadOption(readDtoOption)

	val, err := proxy.Cache.GetValue(ctx, key)
	if err == nil {
		return val, nil
	}

	readFn, exist := proxy.firstDelivery.Load(key)
	if exist {
		return readFn.(ReadDataSource)(ctx, readDtoOption)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	readFn = ReadDataSource(func(context.Context, any) (any, error) {
		wg.Wait()
		return readModel, err
	})

	mainReadFn, ok := proxy.firstDelivery.LoadOrStore(key, readFn)
	if ok {
		return mainReadFn.(ReadDataSource)(ctx, readDtoOption)
	}

	// main read
	defer func() {
		proxy.firstDelivery.Store(key, ReadDataSource(func(context.Context, any) (any, error) {
			return readModel, nil
		}))
		wg.Done()

		go func() {
			time.Sleep(time.Second)
			proxy.firstDelivery.Delete(key)
		}()
	}()

	if proxy.IsAnNotFoundError(err) || proxy.CanIgnoreCacheError {
		readModel, err = proxy.ReadDataSource(ctx, readDtoOption)
		if err != nil {
			if !(proxy.IsAnNotFoundError(err) && proxy.CanIgnoreReadSourceErrorNotFound) {
				return empty, err
			}
		}

		err = proxy.Cache.PutValue(ctx, key, readModel, proxy.CacheTTL)
		if err != nil && !proxy.CanIgnoreCacheError {
			return empty, err
		}

		return readModel, nil
	}

	return empty, err
}

func (proxy *CacheProxySyncMap) ReadValueV3(ctx context.Context, readDtoOption any) (readModel any, err error) {
	var empty any
	key := proxy.TransformReadOption(readDtoOption)

	val, err := proxy.Cache.GetValue(ctx, key)
	if err == nil {
		return val, nil
	}

	readFn, exist := proxy.firstDelivery.Load(key)
	if exist {
		return readFn.(ReadDataSource)(ctx, readDtoOption)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	readFn = ReadDataSource(func(context.Context, any) (any, error) {
		wg.Wait()
		return readModel, err
	})

	mainReadFn, ok := proxy.firstDelivery.LoadOrStore(key, readFn)
	if ok {
		return mainReadFn.(ReadDataSource)(ctx, readDtoOption)
	}

	// main read
	defer func() {
		wg.Done()
		proxy.firstDelivery.Delete(key)
	}()

	// 和 v2 的差異: 再次到 cache 檢查, 確認 cache 是否資料
	val, err = proxy.Cache.GetValue(ctx, key)
	if err != nil {
		if !(proxy.IsAnNotFoundError(err) || proxy.CanIgnoreCacheError) {
			return empty, err
		}

		readModel, err = proxy.ReadDataSource(ctx, readDtoOption)
		if err != nil {
			if !(proxy.IsAnNotFoundError(err) && proxy.CanIgnoreReadSourceErrorNotFound) {
				return empty, err
			}
		}

		err = proxy.Cache.PutValue(ctx, key, readModel, proxy.CacheTTL)
		if err != nil && !proxy.CanIgnoreCacheError {
			return empty, err
		}

		return readModel, nil
	}

	return val, nil
}
