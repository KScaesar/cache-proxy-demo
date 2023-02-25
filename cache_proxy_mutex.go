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
