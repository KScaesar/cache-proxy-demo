package example

import "context"

type Cache interface {
	GetValue(ctx context.Context, key string) (any, error)
	PutValue(ctx context.Context, key string, val any) error
	DeleteValue(ctx context.Context, key string) error
}
