package example

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

type CacheProxySingleflight struct {
	Cache          Cache
	singleDelivery singleflight.Group // key:ReadDataSource(func)

	TransformReadOption func(readDtoOption any) (key string)
	ReadDataSource      func(ctx context.Context, readDtoOption any) (readModel any, err error)

	IsAnNotFoundError                func(err error) bool
	CanIgnoreCacheError              bool
	CanIgnoreReadSourceErrorNotFound bool // source not found, 是否交給 caller 處理
	CacheTTL                         time.Duration
}

func (proxy *CacheProxySingleflight) ReadValue(ctx context.Context, readDtoOption any) (readModel any, err error) {
	return proxy.ReadValueV3(ctx, readDtoOption)
}

func (proxy *CacheProxySingleflight) withoutLockReadValue(ctx context.Context, readDtoOption any) (readModel any, Err error) {
	var empty any

	key := proxy.TransformReadOption(readDtoOption)
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

// v1: userCount = int(2e4)
// BenchmarkCacheProxySingleflight_ReadValue
//    cache_proxy_impl_test.go:236: db qry count = 1, b.N=1
//    cache_proxy_impl_test.go:236: db qry count = 100, b.N=100
//    cache_proxy_impl_test.go:236: db qry count = 3141, b.N=3141
//    cache_proxy_impl_test.go:236: db qry count = 4305, b.N=4305
// BenchmarkCacheProxySingleflight_ReadValue-8   	    4305	    277905 ns/op

func (proxy *CacheProxySingleflight) ReadValueV1(ctx context.Context, readDtoOption any) (readModel any, err error) {
	key := proxy.TransformReadOption(readDtoOption)
	readModel, err, _ = proxy.singleDelivery.Do(key, func() (interface{}, error) {
		return proxy.withoutLockReadValue(ctx, readDtoOption)
	})
	return readModel, err
}

// v2: userCount = int(2e4)
// BenchmarkCacheProxySingleflight_ReadValue
// cache_proxy_impl_test.go:236: db qry count = 1, b.N=1
// cache_proxy_impl_test.go:236: db qry count = 100, b.N=100
// cache_proxy_impl_test.go:236: db qry count = 2929, b.N=2893
// cache_proxy_impl_test.go:236: db qry count = 4456, b.N=4412
// BenchmarkCacheProxySingleflight_ReadValue-8   	    4412	    258161 ns/op

// v2: userCount = int(2e3)
// BenchmarkCacheProxySingleflight_ReadValue
//    cache_proxy_impl_test.go:236: db qry count = 1, b.N=1
//    cache_proxy_impl_test.go:236: db qry count = 100, b.N=100
//    cache_proxy_impl_test.go:236: db qry count = 2015, b.N=3279
//    cache_proxy_impl_test.go:236: db qry count = 2008, b.N=7022
//    cache_proxy_impl_test.go:236: db qry count = 2021, b.N=14248
//    cache_proxy_impl_test.go:236: db qry count = 2018, b.N=24469
//    cache_proxy_impl_test.go:236: db qry count = 2018, b.N=38866
//    cache_proxy_impl_test.go:236: db qry count = 2017, b.N=60798
//    cache_proxy_impl_test.go:236: db qry count = 2018, b.N=86976
// BenchmarkCacheProxySingleflight_ReadValue-8   	   86976	     13433 ns/op

// bug: double main read
//
// cache id[3]
// cache id[3]
// cache id[3]
// cache id[3]
// cache id[3]
// cache id[3]
// main id[3] <-
// cache id[3]
// cache id[3]
// main id[3] <-
//     cache_proxy_impl_test.go:236: db qry count = 2867, b.N=22153

func (proxy *CacheProxySingleflight) ReadValueV2(ctx context.Context, readDtoOption any) (ReadModel any, Err error) {
	var empty any
	key := proxy.TransformReadOption(readDtoOption)

	val, err := proxy.Cache.GetValue(ctx, key)
	if err == nil {
		return val, nil
	}

	ReadModel, Err, _ = proxy.singleDelivery.Do(key, func() (interface{}, error) {
		// 無法解決 bug: double main read
		// 套件自動刪除 main read record
		// 無法靠外力控制

		readModel, err := proxy.ReadDataSource(ctx, readDtoOption)
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
	})

	return ReadModel, Err
}

// v3: userCount = int(2e4)
// BenchmarkCacheProxySingleflight_ReadValue
//    cache_proxy_impl_test.go:236: db qry count = 1, b.N=1
//    cache_proxy_impl_test.go:236: db qry count = 100, b.N=100
//    cache_proxy_impl_test.go:236: db qry count = 2810, b.N=2810
//    cache_proxy_impl_test.go:236: db qry count = 4419, b.N=4419
// BenchmarkCacheProxySingleflight_ReadValue-8   	    4419	    269547 ns/op

func (proxy *CacheProxySingleflight) ReadValueV3(ctx context.Context, readDtoOption any) (readModel any, Err error) {
	var empty any
	key := proxy.TransformReadOption(readDtoOption)

	val, err := proxy.Cache.GetValue(ctx, key)
	if err == nil {
		return val, nil
	}

	readModel, Err, _ = proxy.singleDelivery.Do(key, func() (interface{}, error) {
		// main read 再次判斷 cache 是否存在
		// 確保 不會發生 重複進行 db read
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
	})

	return readModel, Err
}
