package example

import (
	"context"
	"sync"
	"time"
)

type CacheProxyMutex struct {
	Cache Cache
	mu    sync.Mutex

	TransformReadOption func(readDtoOption any) (key string)
	ReadDataSource      func(ctx context.Context, readDtoOption any) (readModel any, err error)

	IsAnNotFoundError                func(err error) bool
	CanIgnoreCacheError              bool
	CanIgnoreReadSourceErrorNotFound bool // source not found, 是否交給 caller 處理
	CacheTTL                         time.Duration
}

func (proxy *CacheProxyMutex) ReadValue(ctx context.Context, readDtoOption any) (readModel any, err error) {
	return proxy.ReadValueV3(ctx, readDtoOption)
}

func (proxy *CacheProxyMutex) ReadValueV1(ctx context.Context, readDtoOption any) (readModel any, err error) {
	var empty any

	proxy.mu.Lock()
	defer proxy.mu.Unlock()

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

func (proxy *CacheProxyMutex) ReadValueV2(ctx context.Context, readDtoOption any) (readModel any, err error) {
	// v2 沒有解決 bug: double main read
	// v3 才能解決

	var empty any

	key := proxy.TransformReadOption(readDtoOption)
	val, err := proxy.Cache.GetValue(ctx, key)
	if err == nil {
		return val, nil
	}

	proxy.mu.Lock()
	defer proxy.mu.Unlock()

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

func (proxy *CacheProxyMutex) ReadValueV3(ctx context.Context, readDtoOption any) (readModel any, err error) {
	// 和 v2 的差異:
	// 再次到 cache 檢查, 確認 cache 是否有資料
	// 如此可確保 不會發生 重複進行 db read

	var empty any

	key := proxy.TransformReadOption(readDtoOption)
	val, err := proxy.Cache.GetValue(ctx, key)
	if err == nil {
		return val, nil
	}

	proxy.mu.Lock()
	defer proxy.mu.Unlock()

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
