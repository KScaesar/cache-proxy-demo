package infra

import (
	"context"
	"time"
)

// 內容和 root 目錄一樣, 只是用來示範 cache proxy 如何使用

type Cache interface {
	GetValue(ctx context.Context, key string) (val any, err error)
	PutValue(ctx context.Context, key string, val any, ttl time.Duration) error
	DeleteValue(ctx context.Context, key string) error
}

func NewCacheRedis() *CacheRedis {
	return &CacheRedis{}
}

type CacheRedis struct{}

func (cache *CacheRedis) DeleteValue(ctx context.Context, key string) error {
	// TODO implement me
	panic("implement me")
}

func (cache *CacheRedis) GetValue(ctx context.Context, key string) (any, error) {
	// TODO implement me
	panic("implement me")
}

func (cache *CacheRedis) PutValue(ctx context.Context, key string, val any, ttl time.Duration) error {
	// TODO implement me
	panic("implement me")
}
