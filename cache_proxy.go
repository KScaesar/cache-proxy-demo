package example

import (
	"context"
	"time"
)

type CacheProxy interface {
	ReadValue(ctx context.Context, readDtoOption any) (readModel any, err error)
}

type CacheProxyV1 struct {
	Cache Cache

	TransformReadOption func(readDtoOption any) (key string)
	ReadSource          func(ctx context.Context, readDtoOption any) (readModel any, err error)

	IsNotFound                       func(err error) bool
	CanIgnoreCacheError              bool
	CanIgnoreReadSourceErrorNotFound bool // source not found, 是否交給 caller 處理
	CacheTTL                         time.Duration
}

func (proxy *CacheProxyV1) ReadValue(ctx context.Context, readDtoOption any) (readModel any, err error) {
	var empty any

	key := proxy.TransformReadOption(readDtoOption)
	val, err := proxy.Cache.GetValue(ctx, key)
	if err == nil {
		return val, nil
	}
	if !proxy.IsNotFound(err) && !proxy.CanIgnoreCacheError {
		return empty, err
	}

	readModel, err = proxy.ReadSource(ctx, readDtoOption)
	if err != nil && !proxy.CanIgnoreReadSourceErrorNotFound {
		return empty, err
	}

	err = proxy.Cache.PutValue(ctx, key, readModel, proxy.CacheTTL)
	if err != nil && !proxy.CanIgnoreCacheError {
		return empty, err
	}

	return readModel, nil
}
