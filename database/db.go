package database

import (
	//"sync"
	"github.com/dmarro89/dare-db/logger"
	"github.com/dmarro89/go-redis-hashtable/datastr"

)

type Database struct {
	Dict *datastr.Dict
	//mu   sync.RWMutex
}

func NewDatabase(logger *darelog.LOG) *Database {
	return &Database{
		Dict: datastr.NewDict(logger),
	}
}

func (db *Database) Get(key string) interface{} {
	return db.Dict.Get(key)
}

func (db *Database) Set(key string, value interface{}) error {
	return db.Dict.Set(key, value)
}

func (db *Database) Delete(key string) error {
	return db.Dict.Delete(key)
}
