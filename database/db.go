package database

import (
	//"time"
	"github.com/go-while/nodare-db/logger"
)

type XDatabase struct {
	XDICK *XDICK
}

func NewDICK(logger *ilog.LOG, sdCh chan uint32, waitCh chan struct{}) *XDatabase {
	returnsubDICKs := make(chan []*SubDICK, 1)
	xdick := NewXDICK(logger, sdCh, returnsubDICKs)
	db := &XDatabase{
		XDICK: xdick,
	}
	go func(returnsubDICKs <- chan []*SubDICK, db *XDatabase) {
		db.XDICK.logger.Debug("NewDICK waits async to return subDICKs")
		subDICKs := <- returnsubDICKs

		db.XDICK.SubDICKs = subDICKs

		// reads re-pushed value from NewXDICK
		// which has been read from config and passed through sdCh
		db.XDICK.SubCount = <- sdCh

		for j := range db.XDICK.SubDICKs {
			go db.XDICK.watchDog(uint32(j))
		}

		db.XDICK.logger.Debug("NewDICK set subDICKs=%d/%d notify waitCh", len(subDICKs), len(db.XDICK.SubDICKs))
		waitCh <- struct{}{}
	}(returnsubDICKs, db)

	return db
}

func (db *XDatabase) Get(key string) interface{} {
	return db.XDICK.Get(key)
}

func (db *XDatabase) Set(key string, value interface{}) error {
	return db.XDICK.Set(key, value)
}

func (db *XDatabase) Del(key string) error {
	return db.XDICK.Del(key)
}
