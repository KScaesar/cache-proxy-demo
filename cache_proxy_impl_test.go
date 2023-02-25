package example

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"testing"
	"time"
)

const userCount = int(2e3)

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

	CacheProxy_Benchmark(b, proxy, db)
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

	CacheProxy_Benchmark(b, proxy, db)
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

	CacheProxy_Benchmark(b, proxy, db)
}

func CacheProxy_Benchmark(b *testing.B, proxy CacheProxy, db *Database) {
	users := db.users
	ids := GetUserIdAll(users)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%userCount]
		// id := ids[rand.Intn(userCount)]
		start, end := concurrencyWorker(func() {
			proxy.ReadValue(nil, id)
		})
		close(start)
		<-end
	}

	b.Logf("db qry count = %v, b.N=%v", db.qryCount, b.N)
}

func concurrencyWorker(action func()) (
	start chan struct{},
	end chan struct{},
) {
	var wg sync.WaitGroup
	start = make(chan struct{})
	end = make(chan struct{})

	// workerCount := 1
	workerCount := 1 * runtime.GOMAXPROCS(0)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			<-start
			action()
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(end)
	}()
	return start, end
}
