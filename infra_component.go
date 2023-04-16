package cache_proxy_demo

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

var ErrNotFound = errors.New("not found")

// cache

type Cache interface {
	GetValue(key string, valType any) (val any, err error)
	SetValue(key string, val any) error
}

func NewMutexCache() *MutexCache {
	return &MutexCache{
		data: make(map[string]any),
	}
}

type MutexCache struct {
	mu   sync.Mutex
	data map[string]any
}

func (cache *MutexCache) GetValue(key string, valType any) (any, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if v, ok := cache.data[key]; ok {
		return v, nil
	}
	return valType, ErrNotFound
}

func (cache *MutexCache) SetValue(key string, val any) error {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.data[key] = val
	return nil
}

// database

type UserRepository interface {
	QueryUserById(id string) (users *gofakeit.PersonInfo, err error)
}

func NewUserDatabase(size int) *UserDatabase {
	newUsers := func(size int) map[string]*gofakeit.PersonInfo {
		faker := gofakeit.NewCrypto()
		gofakeit.SetGlobalFaker(faker)
		users := make(map[string]*gofakeit.PersonInfo, size)
		for i := 0; i < size; i++ {
			id := fmt.Sprintf("id[%v]", i)
			users[id] = gofakeit.Person()
		}
		return users
	}

	return &UserDatabase{
		total: size,
		users: newUsers(size),
	}
}

type UserDatabase struct {
	total    int
	qryCount int
	mu       sync.Mutex
	users    map[string]*gofakeit.PersonInfo
}

func (db *UserDatabase) QueryUserById(id string) (users *gofakeit.PersonInfo, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.qryCount++
	time.Sleep(10 * time.Microsecond)
	v, ok := db.users[id]
	if !ok {
		return nil, ErrNotFound
	}

	return v, nil
}

func (db *UserDatabase) GetUserIds() []string {
	size := len(db.users)
	ids := make([]string, 0, size)
	for i := 0; i < size; i++ {
		ids = append(ids, fmt.Sprintf("id[%v]", i))
	}
	return ids
}
