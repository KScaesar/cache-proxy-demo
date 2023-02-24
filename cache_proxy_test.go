package example

import (
	"context"
	"errors"
	"testing"
)

func TestCacheProxyV1_ReadValue(t *testing.T) {
	// arrange
	users := NewUsers(100)
	db := NewDatabase(users)
	cache := NewMutexCache()

	proxy := CacheProxyV1{
		Cache: cache,

		TransformReadOption: func(readDtoOption any) (key string) {
			return readDtoOption.(string)
		},
		ReadSource: func(ctx context.Context, readDtoOption any) (readModel any, err error) {
			id := readDtoOption.(string)
			return db.QueryUserForShareMode(id)
		},
		IsNotFound: func(err error) bool {
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
		ids := GetUserIdAll(users)
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

func BenchmarkCacheProxyV1_ReadValue(b *testing.B) {
	// arrange
	users := NewUsers(1e3)
	db := NewDatabase(users)
	cache := NewMutexCache()

	proxy := CacheProxyV1{
		Cache: cache,

		TransformReadOption: func(readDtoOption any) (key string) {
			return readDtoOption.(string)
		},
		ReadSource: func(ctx context.Context, readDtoOption any) (readModel any, err error) {
			id := readDtoOption.(string)
			return db.QueryUserForShareMode(id)
		},
		IsNotFound: func(err error) bool {
			return errors.Is(err, ErrNotFound)
		},
		CanIgnoreCacheError:              false,
		CanIgnoreReadSourceErrorNotFound: true,
		CacheTTL:                         0,
	}

	for i := 0; i < b.N; i++ {

	}
}
