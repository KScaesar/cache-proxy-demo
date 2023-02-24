package example

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func Benchmark_ParallelCache_WriteRation(b *testing.B) {
	// for i := 0; i <= 10; i++ {
	// 	writeRatio := 0 + i*10
	// 	ParallelCache_WriteRation(b, "MutexCache", NewMutexCache, writeRatio)
	// 	ParallelCache_WriteRation(b, "RwMutexCache", NewRwMutexCache, writeRatio)
	// 	ParallelCache_WriteRation(b, "SyncMapCache", NewSyncMapCache, writeRatio)
	// 	ParallelCache_WriteRation(b, "ChanCache", NewChanCache, writeRatio)
	// }

	writeRatio := 5
	ParallelCache_WriteRation(b, "MutexCache", NewMutexCache, writeRatio)
	ParallelCache_WriteRation(b, "RwMutexCache", NewRwMutexCache, writeRatio)
	ParallelCache_WriteRation(b, "SyncMapCache", NewSyncMapCache, writeRatio)
	ParallelCache_WriteRation(b, "ChanCache", NewChanCache, writeRatio)

}

func ParallelCache_WriteRation(b *testing.B, fnName string, NewCache func() Cache, writeRatio int) {
	b.Run(fmt.Sprintf("%v-WriteRatio=%v%%", fnName, writeRatio), func(b *testing.B) {
		cache := NewCache()
		ctx := context.Background()
		dataCount := int(1e2)
		users := NewUsers(dataCount)
		userIds := GetUserIdAll(users)
		isWrite := func() bool {
			return rand.Intn(100) < writeRatio
		}

		var wg sync.WaitGroup
		workerSize := 16
		task := make(chan struct{}, workerSize)
		for i := 0; i < workerSize; i++ {
			wg.Add(1)
			go func() {
				for range task {
					id := userIds[rand.Intn(dataCount)]
					val := users[id]
					if isWrite() {
						cache.PutValue(ctx, id, val, 0)
					} else {
						cache.GetValue(ctx, id)
					}
				}
				wg.Done()
			}()
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			task <- struct{}{}
		}

		close(task)
		wg.Wait()
	})
}

//

func Benchmark_ParallelCache_WriteReadOnly(b *testing.B) {
	ParallelCache_WriteReadOnly(b, "MutexCache", NewMutexCache)
	ParallelCache_WriteReadOnly(b, "RwMutexCache", NewRwMutexCache)
	ParallelCache_WriteReadOnly(b, "SyncMapCache", NewSyncMapCache)
	ParallelCache_WriteReadOnly(b, "ChanCache", NewChanCache)
}

func ParallelCache_WriteReadOnly(b *testing.B, fnName string, NewCache func() Cache) {

	// b.Log("write-before")
	b.Run(fmt.Sprintf("%v-Write-Only", fnName), func(b *testing.B) {
		cache := NewCache()
		ctx := context.Background()
		dataCount := int(1e2)
		users := NewUsers(dataCount)
		userIds := GetUserIdAll(users)

		var wg sync.WaitGroup
		workerSize := 8
		// workerSize := int(1e2)
		task := make(chan struct{}, workerSize)
		for i := 0; i < workerSize; i++ {
			wg.Add(1)
			go func() {
				for range task {
					id := userIds[rand.Intn(dataCount)]
					val := users[id]
					cache.PutValue(ctx, id, val, 0)
				}
				wg.Done()
			}()
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			task <- struct{}{}
		}

		close(task)
		wg.Wait()
	})
	// b.Log("write-after")

	// b.Log("read-before")
	b.Run(fmt.Sprintf("%v-Read-Only", fnName), func(b *testing.B) {
		cache := NewCache()
		ctx := context.Background()
		dataCount := int(1e2)
		users := NewUsers(dataCount)
		userIds := GetUserIdAll(users)
		for id, user := range users {
			cache.PutValue(ctx, id, user, 0)
		}

		var wg sync.WaitGroup
		workerSize := 8
		// workerSize := int(1e2)
		task := make(chan struct{}, workerSize)
		for i := 0; i < workerSize; i++ {
			wg.Add(1)
			go func() {
				for range task {
					id := userIds[rand.Intn(dataCount)]
					cache.GetValue(ctx, id)
				}
				wg.Done()
			}()
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			task <- struct{}{}
		}

		close(task)
		wg.Wait()
	})
	// b.Log("read-after")
}

//

func Benchmark_SequenceCache(b *testing.B) {
	SequenceCache(b, "MutexCache", NewMutexCache)
	SequenceCache(b, "RwMutexCache", NewRwMutexCache)
	SequenceCache(b, "SyncMapCache", NewSyncMapCache)
	SequenceCache(b, "ChanCache", NewChanCache)
}

func SequenceCache(b *testing.B, fnName string, NewCache func() Cache) {

	b.Run(fmt.Sprintf("%v-Put-Sequence", fnName), func(b *testing.B) {
		cache := NewCache()
		ctx := context.Background()
		dataCount := int(1e2)
		users := NewUsers(dataCount)
		userIds := GetUserIdAll(users)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id := userIds[rand.Intn(dataCount)]
			val := users[id]
			cache.PutValue(ctx, id, val, 0)
		}
	})

	b.Run(fmt.Sprintf("%v-Get-Sequence", fnName), func(b *testing.B) {
		cache := NewCache()
		ctx := context.Background()
		dataCount := int(1e2)
		users := NewUsers(dataCount)
		userIds := GetUserIdAll(users)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id := userIds[rand.Intn(dataCount)]
			cache.GetValue(ctx, id)
		}
	})
}

func TestCache(t *testing.T) {
	// cache := NewMutexCache()
	// cache := NewRwMutexCache()
	// cache := NewSyncMapCache()
	cache := NewChanCache()

	cache.PutValue(nil, "x", 1, 0)
	t.Log(cache.GetValue(nil, "x"))
	t.Log(cache.GetValue(nil, "y"))
}
