package example

import (
	"fmt"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

type UserRepository interface {
	QueryUserForShareMode(id string) (users *gofakeit.PersonInfo, err error)
}

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
	time.Sleep(10 * time.Microsecond)
	v, ok := db.users[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func GetUserIdAll(users map[string]*gofakeit.PersonInfo) []string {
	idList := make([]string, 0, len(users))
	for i := 0; i < len(users); i++ {
		idList = append(idList, fmt.Sprintf("id[%v]", i))
	}
	return idList
}

func NewUsers(size int) map[string]*gofakeit.PersonInfo {
	faker := gofakeit.NewCrypto()
	gofakeit.SetGlobalFaker(faker)
	users := make(map[string]*gofakeit.PersonInfo, size)
	for i := 0; i < size; i++ {
		id := fmt.Sprintf("id[%v]", i)
		users[id] = gofakeit.Person()
	}
	return users
}
