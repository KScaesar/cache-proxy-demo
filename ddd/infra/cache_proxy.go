package infra

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

// 內容和 root 目錄一樣, 只是用來示範 cache proxy 如何使用

type CacheProxy interface {
	ReadValue(ctx context.Context, readDtoOption any) (readModel any, err error)
}

type CacheProxyImpl struct {
	Cache          Cache
	singleDelivery singleflight.Group

	TransformReadOption func(readDtoOption any) (key string)
	ReadDataSource      func(ctx context.Context, readDtoOption any) (readModel any, err error)

	IsAnNotFoundError                func(err error) bool
	CanIgnoreCacheError              bool
	CanIgnoreReadSourceErrorNotFound bool // source not found, 是否交給 caller 處理
	CacheTTL                         time.Duration
}

func (proxy *CacheProxyImpl) ReadValue(ctx context.Context, readDtoOption any) (readModel any, err error) {
	// TODO implement me
	panic("implement me")
}
