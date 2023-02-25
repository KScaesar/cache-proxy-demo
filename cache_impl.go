package example

import (
	"context"
	"sync"
	"time"
)

func NewMutexCache() Cache {
	return &MutexCache{
		data: make(map[string]any),
	}
}

type MutexCache struct {
	mu   sync.Mutex
	data map[string]any
}

func (cache *MutexCache) DeleteValue(ctx context.Context, key string) error {
	// TODO implement me
	panic("implement me")
}

func (cache *MutexCache) GetValue(ctx context.Context, key string) (any, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if v, ok := cache.data[key]; ok {
		return v, nil
	}
	var empty any
	return empty, ErrNotFound
}

func (cache *MutexCache) PutValue(ctx context.Context, key string, val any, ttl time.Duration) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.data[key] = val
	return nil
}

//

func NewRwMutexCache() Cache {
	return &RwMutexCache{
		data: make(map[string]any),
	}
}

type RwMutexCache struct {
	mu   sync.RWMutex
	data map[string]any
}

func (cache *RwMutexCache) DeleteValue(ctx context.Context, key string) error {
	// TODO implement me
	panic("implement me")
}

func (cache *RwMutexCache) GetValue(ctx context.Context, key string) (any, error) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	if v, ok := cache.data[key]; ok {
		return v, nil
	}
	var empty any
	return empty, ErrNotFound
}

func (cache *RwMutexCache) PutValue(ctx context.Context, key string, val any, ttl time.Duration) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.data[key] = val
	return nil
}

//

func NewSyncMapCache() Cache {
	return &SyncMapCache{}
}

type SyncMapCache struct {
	data sync.Map
}

func (cache *SyncMapCache) DeleteValue(ctx context.Context, key string) error {
	// TODO implement me
	panic("implement me")
}

func (cache *SyncMapCache) GetValue(ctx context.Context, key string) (any, error) {
	if v, ok := cache.data.Load(key); ok {
		return v, nil
	}
	var empty any
	return empty, ErrNotFound
}

func (cache *SyncMapCache) PutValue(ctx context.Context, key string, val any, ttl time.Duration) error {
	cache.data.Store(key, val)
	return nil
}

//

func NewChanCache() Cache {
	cache := &ChanCache{
		getAction: make(chan CommandGet),
		putAction: make(chan CommandPut),
		data:      make(map[string]any),
	}
	go cache.manager()
	return cache
}

type ChanCache struct {
	getAction chan CommandGet
	putAction chan CommandPut
	data      map[string]any
}

func (cache *ChanCache) DeleteValue(ctx context.Context, key string) error {
	// TODO implement me
	panic("implement me")
}

func (cache *ChanCache) GetValue(ctx context.Context, key string) (any, error) {
	cmd := CommandGet{
		key:   key,
		reply: make(chan cacheResult, 1),
	}
	cache.getAction <- cmd

	r := <-cmd.reply
	return r.val, r.err
}

func (cache *ChanCache) PutValue(ctx context.Context, key string, val any, ttl time.Duration) error {
	cmd := CommandPut{
		key:   key,
		val:   val,
		reply: make(chan cacheResult, 1),
	}
	cache.putAction <- cmd

	r := <-cmd.reply
	return r.err
}

func (cache *ChanCache) manager() {
	for {
		select {
		case cmd := <-cache.getAction:
			if v, ok := cache.data[cmd.key]; ok {
				cmd.reply <- cacheResult{
					val: v,
					err: nil,
				}
			} else {
				cmd.reply <- cacheResult{
					val: nil,
					err: ErrNotFound,
				}
			}

		case cmd := <-cache.putAction:
			cache.data[cmd.key] = cmd.val
			cmd.reply <- cacheResult{}
		}
	}
}

type cacheResult struct {
	val any
	err error
}

type CommandGet struct {
	ctx   context.Context
	key   string
	reply chan cacheResult
}

type CommandPut struct {
	key   string
	val   any
	reply chan cacheResult
}
