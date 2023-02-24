package example

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type CacheProxy interface {
	ReadValue(ctx context.Context, readDtoOption any) (readModel any, err error)
}

type Cache interface {
	GetValue(ctx context.Context, key string) (val any, err error)
	PutValue(ctx context.Context, key string, val any, ttl time.Duration) error
	DeleteValue(ctx context.Context, key string) error
}

var ErrNotFound = errors.New("not found")

func WrapErrWithMsg(err error, msg string, args ...any) error {
	return fmt.Errorf("%v: %w", fmt.Sprintf(msg, args...), err)
}
