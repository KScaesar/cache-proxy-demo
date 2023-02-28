package pattern

import (
	"context"
	"time"
)

type CacheProxy interface {
	TransformMessage(msg any) (key string)
	Cache
	DatabaseGetter
}

type Cache interface {
	GetCacheValue(ctx context.Context, key string) (val any, err error)
	SetCacheValue(ctx context.Context, key string, val any, ttl time.Duration) (err error)
	DeleteCacheValue(ctx context.Context, key string) (err error)
}

type DatabaseGetter interface {
	GetDatabaseValue(ctx context.Context, msg any) (result any, err error)
}

type DatabaseGetFunc func(ctx context.Context, msg any) (result any, err error)

func (fn DatabaseGetFunc) GetDatabaseValue(ctx context.Context, msg any) (result any, err error) {
	return fn(ctx, msg)
}
