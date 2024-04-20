// database.go

package database

import (
	"sync"

	"github.com/dmarro89/go-redis-hashtable/datastr"
)

type Database struct {
	dict *datastr.Dict
	mu   sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		dict: datastr.NewDict(),
	}
}

func (db *Database) Get(key string) interface{} {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.dict.Get(key)
}

func (db *Database) Set(key string, value interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.dict.Set(key, value)
}

func (db *Database) Delete(key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.dict.Delete(key)
}
