package cache_proxy_demo

import (
	"runtime"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func ConcurrentTester(goroutinePower uint8, fn func()) (start func(), wait func()) {
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

	start = func() { close(ready) }
	wait = func() { <-done }
	return start, wait
}

func CacheProxyBenchmarkConcurrentOneKey(b *testing.B, proxy CacheProxy, db *UserDatabase) {
	ids := db.GetUserIds()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id := ids[i%db.total]
		start, wait := ConcurrentTester(1, func() {
			proxy.Execute(id, &gofakeit.PersonInfo{})
		})
		start()
		wait()
	}

	b.Logf("one key: db qry count = %v, b.N=%v", db.qryCount, b.N)
}

func CacheProxyBenchmarkConcurrentMultiKey(b *testing.B, proxy CacheProxy, db *UserDatabase) {
	ids := db.GetUserIds()
	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		id := ids[i%db.total]
		go func() {
			start, wait := ConcurrentTester(1, func() {
				proxy.Execute(id, &gofakeit.PersonInfo{})
			})
			start()
			wait()
			wg.Done()
		}()
	}
	wg.Wait()

	b.Logf("multi key: db qry count = %v, b.N=%v", db.qryCount, b.N)
}

func dependency() (string, *BaseCacheProxy, *UserDatabase) {
	// concurrentKeys := "one"
	concurrentKeys := "multi"
	dataSize := 2e4

	db := NewUserDatabase(int(dataSize))
	baseProxy := &BaseCacheProxy{
		Transform: TransformQryOptionToCacheKey(func(qryOption any) (key string) {
			return qryOption.(string)
		}),

		Cache: NewMutexCache(),

		GetDB: DatabaseGetFunc(func(qryOption any) (result any, err error) {
			id := qryOption.(string)
			return db.QueryUserById(id)
		}),
	}

	return concurrentKeys, baseProxy, db
}

func BenchmarkMutexProxy(b *testing.B) {
	concurrentKeys, baseProxy, db := dependency()
	switch concurrentKeys {
	case "one":
		CacheProxyBenchmarkConcurrentOneKey(b, UseMutex(baseProxy), db)
	case "multi":
		CacheProxyBenchmarkConcurrentMultiKey(b, UseMutex(baseProxy), db)
	}
}

func BenchmarkChannelProxy(b *testing.B) {
	concurrentKeys, baseProxy, db := dependency()
	switch concurrentKeys {
	case "one":
		CacheProxyBenchmarkConcurrentOneKey(b, UseChannel(baseProxy), db)
	case "multi":
		CacheProxyBenchmarkConcurrentMultiKey(b, UseChannel(baseProxy), db)
	}
}

func BenchmarkSyncMapProxy(b *testing.B) {
	concurrentKeys, baseProxy, db := dependency()
	switch concurrentKeys {
	case "one":
		CacheProxyBenchmarkConcurrentOneKey(b, UseSyncMap(baseProxy), db)
	case "multi":
		CacheProxyBenchmarkConcurrentMultiKey(b, UseSyncMap(baseProxy), db)
	}
}

func BenchmarkSingleflightProxy(b *testing.B) {
	concurrentKeys, baseProxy, db := dependency()
	switch concurrentKeys {
	case "one":
		CacheProxyBenchmarkConcurrentOneKey(b, UseSingleflight(baseProxy), db)
	case "multi":
		CacheProxyBenchmarkConcurrentMultiKey(b, UseSingleflight(baseProxy), db)
	}
}
