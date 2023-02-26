package example

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"
)

const userCount = int(2e2)

func TestCacheProxyMutex_ReadValue(t *testing.T) {
	// arrange
	users := NewUsers(userCount)
	ids := GetUserIdAll(users)
	db := NewDatabase(users)
	cache := NewMutexCache()

	proxy := CacheProxyMutex{
		Cache: cache,

		TransformReadOption: func(readDtoOption any) (key string) {
			return readDtoOption.(string)
		},
		ReadDataSource: func(ctx context.Context, readDtoOption any) (readModel any, err error) {
			id := readDtoOption.(string)
			return db.QueryUserForShareMode(id)
		},
		IsAnNotFoundError: func(err error) bool {
			return errors.Is(err, ErrNotFound)
		},
		CanIgnoreCacheError:              false,
		CanIgnoreReadSourceErrorNotFound: true,
		CacheTTL:                         0,
	}

	var err error

	// assert cache
	{
		maxCnt := 10
		for i := 0; i < maxCnt; i++ {
			_, err := proxy.ReadValue(nil, ids[i])
			if err != nil {
				t.Error(err)
				return
			}
		}

		_, err = proxy.ReadValue(nil, ids[4])
		if err != nil {
			t.Error(err)
			return
		}

		resp, err := proxy.ReadValue(nil, ids[20])
		if err != nil {
			t.Error(err)
			return
		}
		_ = resp
		// t.Log(resp)

		expectedQueryCount := 11
		if db.qryCount != expectedQueryCount {
			t.Errorf("expected = %v, but actual = %v", expectedQueryCount, db.qryCount)
			return
		}
	}

	// assert CanIgnoreReadSourceErrorNotFound
	{
		proxy.CanIgnoreReadSourceErrorNotFound = true
		_, err = proxy.ReadValue(nil, "")
		if err != nil {
			t.Error(err)
			return
		}

		_, err = proxy.ReadValue(nil, "")
		if err != nil {
			t.Error(err)
			return
		}

		expectedQueryCount := 12
		if db.qryCount != expectedQueryCount {
			t.Errorf("expected = %v, but actual = %v", expectedQueryCount, db.qryCount)
			return
		}
	}

	// assert CanIgnoreReadSourceErrorNotFound
	{
		proxy.CanIgnoreReadSourceErrorNotFound = false
		_, err = proxy.ReadValue(nil, "caesar")
		if err != nil && !errors.Is(err, ErrNotFound) {
			t.Errorf("expected = ErrNotFoound, but actual = %v", err)
			return
		}

		_, err = proxy.ReadValue(nil, "caesar")
		if err != nil && !errors.Is(err, ErrNotFound) {
			t.Errorf("expected = ErrNotFoound, but actual = %v", err)
			return
		}

		expectedQueryCount := 14
		if db.qryCount != expectedQueryCount {
			t.Errorf("expected = %v, but actual = %v", expectedQueryCount, db.qryCount)
			return
		}
	}
}

func BenchmarkCacheProxyMutex_ReadValue(b *testing.B) {
	users := NewUsers(userCount)
	db := NewDatabase(users)
	cache := NewMutexCache()

	proxy := &CacheProxyMutex{
		Cache: cache,

		TransformReadOption: func(readDtoOption any) (key string) {
			return readDtoOption.(string)
		},
		ReadDataSource: func(ctx context.Context, readDtoOption any) (readModel any, err error) {
			id := readDtoOption.(string)
			return db.QueryUserForShareMode(id)
		},
		IsAnNotFoundError: func(err error) bool {
			return errors.Is(err, ErrNotFound)
		},
		CanIgnoreCacheError:              false,
		CanIgnoreReadSourceErrorNotFound: true,
		CacheTTL:                         0,
	}

	CacheProxyBenchmark(b, proxy, db)
}

func BenchmarkCacheProxyChannel_ReadValue(b *testing.B) {
	users := NewUsers(userCount)
	db := NewDatabase(users)
	cache := NewMutexCache()

	proxy := NewCacheProxyChannel(
		cache,
		func(readDtoOption any) (key string) {
			return readDtoOption.(string)
		},
		func(ctx context.Context, readDtoOption any) (readModel any, err error) {
			id := readDtoOption.(string)
			return db.QueryUserForShareMode(id)
		},
		func(err error) bool {
			return errors.Is(err, ErrNotFound)
		},
		false,
		false,
		time.Hour,
	)

	CacheProxyBenchmark(b, proxy, db)
}

func BenchmarkCacheProxySyncMap_ReadValue(b *testing.B) {
	users := NewUsers(userCount)
	db := NewDatabase(users)
	cache := NewMutexCache()

	proxy := &CacheProxySyncMap{
		Cache: cache,

		TransformReadOption: func(readDtoOption any) (key string) {
			return readDtoOption.(string)
		},
		ReadDataSource: func(ctx context.Context, readDtoOption any) (readModel any, err error) {
			id := readDtoOption.(string)
			return db.QueryUserForShareMode(id)
		},
		IsAnNotFoundError: func(err error) bool {
			return errors.Is(err, ErrNotFound)
		},
		CanIgnoreCacheError:              false,
		CanIgnoreReadSourceErrorNotFound: true,
		CacheTTL:                         0,
	}

	CacheProxyBenchmark(b, proxy, db)
}

func BenchmarkCacheProxySingleflight_ReadValue(b *testing.B) {
	users := NewUsers(userCount)
	db := NewDatabase(users)
	cache := NewMutexCache()

	proxy := &CacheProxySingleflight{
		Cache: cache,

		TransformReadOption: func(readDtoOption any) (key string) {
			return readDtoOption.(string)
		},
		ReadDataSource: func(ctx context.Context, readDtoOption any) (readModel any, err error) {
			id := readDtoOption.(string)
			return db.QueryUserForShareMode(id)
		},
		IsAnNotFoundError: func(err error) bool {
			return errors.Is(err, ErrNotFound)
		},
		CanIgnoreCacheError:              false,
		CanIgnoreReadSourceErrorNotFound: true,
		CacheTTL:                         0,
	}

	CacheProxyBenchmark(b, proxy, db)
}

func CacheProxyBenchmark(b *testing.B, proxy CacheProxy, db *Database) {
	users := db.users
	ids := GetUserIdAll(users)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%userCount]
		// id := ids[rand.Intn(userCount)]
		start, wait := ConcurrentTester(1, func() {
			proxy.ReadValue(nil, id)
		})
		start()
		wait()
	}

	b.Logf("db qry count = %v, b.N=%v", db.qryCount, b.N)
}

func ConcurrentTester(goroutinePower uint8, fn func()) (
	start func(),
	wait func(),
) {
	var wg sync.WaitGroup
	ready := make(chan struct{})
	done := make(chan struct{})

	var workerCount int
	if goroutinePower == 0 {
		workerCount = 1 // sequential
	} else {
		workerCount = int(goroutinePower) * runtime.GOMAXPROCS(0)
	}
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			<-ready
			fn()
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	return func() { close(ready) }, func() { <-done }
}
