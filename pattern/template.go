package pattern

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type CacheProxyWorkflow func(proxy CacheProxy, ctx context.Context, msg any, ttl time.Duration) (result any, err error)

func CacheProxyWorkflow_GetStrategyA(proxy CacheProxy, ctx context.Context, msg any, ttl time.Duration) (result any, err error) {
	var empty any
	key := proxy.TransformMessage(msg)

	val, err := proxy.GetCacheValue(ctx, key)
	if err == nil {
		return val, nil
	}

	if !errors.Is(err, ErrNotFound) {
		return empty, err
	}

	result, err = proxy.GetDatabaseValue(ctx, msg)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return empty, err
		}
	}

	err = proxy.SetCacheValue(ctx, key, result, ttl)
	if err != nil {
		return empty, err
	}

	return result, nil
}

func CacheProxyWorkflow_GetStrategyB(proxy CacheProxy, ctx context.Context, msg any, ttl time.Duration) (result any, err error) {
	var empty any
	key := proxy.TransformMessage(msg)

	val, err := proxy.GetCacheValue(ctx, key)
	if err == nil {
		return val, nil
	}

	if !errors.Is(err, ErrNotFound) {
		// 	print log
	}

	result, err = proxy.GetDatabaseValue(ctx, msg)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return empty, err
		}
	}

	err = proxy.SetCacheValue(ctx, key, result, ttl)
	if err != nil {
		// 	print log
	}

	return result, nil
}

func WrapSingleflightStrategy1(flow CacheProxyWorkflow) CacheProxyWorkflow {
	var g singleflight.Group

	return func(proxy CacheProxy, ctx context.Context, msg any, ttl time.Duration) (result any, err error) {
		key := proxy.TransformMessage(msg)
		result, err, _ = g.Do(key, func() (interface{}, error) {
			return flow(proxy, ctx, msg, ttl)
		})
		return
	}
}

func WrapSyncMapStrategy1(flow CacheProxyWorkflow, localTTL time.Duration) CacheProxyWorkflow {
	var store sync.Map // key:func() (any, error)

	return func(proxy CacheProxy, ctx context.Context, msg any, ttl time.Duration) (result any, err error) {
		key := proxy.TransformMessage(msg)

		var wg sync.WaitGroup
		wg.Add(1)
		firstGoroutine := func() (any, error) {
			wg.Wait()
			return result, err // 注意不可以引用到 匿名函數的變數
		}

		firstFn, isFirst := store.LoadOrStore(key, firstGoroutine)
		if !isFirst {
			// 其他 goroutine 得到的是 first goroutine 的閉包函數
			// 利用閉包機制 會被其他 goroutine 取得 first goroutine 回傳值
			return firstFn.(func() (any, error))()
		}

		defer func() {
			wg.Done()
			go func() {
				time.Sleep(localTTL)
				store.Delete(key)
			}()
		}()
		return flow(proxy, ctx, msg, ttl)
	}
}

func WrapSyncMapStrategy3(flow CacheProxyWorkflow) CacheProxyWorkflow {
	var store sync.Map // key:func() (any, error)

	return func(proxy CacheProxy, ctx context.Context, msg any, ttl time.Duration) (result any, err error) {
		key := proxy.TransformMessage(msg)

		val, err := proxy.GetCacheValue(ctx, key)
		if err == nil {
			return val, nil
		}

		var wg sync.WaitGroup
		wg.Add(1)
		firstGoroutine := func() (any, error) {
			wg.Wait()
			return result, err // 注意不可以引用到 匿名函數的變數
		}

		firstFn, isFist := store.LoadOrStore(key, firstGoroutine)
		if !isFist {
			// 其他 goroutine 得到的是 first goroutine 的閉包函數
			// 利用閉包機制 會被其他 goroutine 取得 first goroutine 回傳值
			return firstFn.(func() (any, error))()
		}

		defer func() {
			wg.Done()
			store.Delete(key)
		}()
		return flow(proxy, ctx, msg, ttl)
	}
}
