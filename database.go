package example

import (
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func NewDatabase(users map[string]*gofakeit.PersonInfo) *Database {
	return &Database{
		users: users,
	}
}

type Database struct {
	qryCount int
	mu       sync.Mutex
	users    map[string]*gofakeit.PersonInfo
}

func (db *Database) QueryUserForShareMode(id string) (users *gofakeit.PersonInfo, err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.qryCount++
	time.Sleep(10 * time.Millisecond)
	v, ok := db.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func GetUserIdAll(users map[string]*gofakeit.PersonInfo) []string {
	idList := make([]string, 0, len(users))
	for key := range users {
		idList = append(idList, key)
	}
	return idList
}

func NewUsers(size int) map[string]*gofakeit.PersonInfo {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	users := make(map[string]*gofakeit.PersonInfo, size)
	for i := 0; i < size; i++ {
		id := uuid.New().String()
		users[id] = gofakeit.Person()
	}
	return users
}
