package database

import (
	"fmt"
	pcas "github.com/go-while/nodare-db/pcas_hash"
)

//const MOD = 10 // last 1 digit
//const MOD = 100 // last 2 digits
//const MOD = 1000 // last 3 digits

// Get returns the value associated with the given key in the dictionary.
//
// Parameters:
// - key: the key to look up in the dictionary.
//
// Return:
// - interface{}: the value associated with the key, or nil if the key is not found.
func (d *XDICK) Get(key string) interface{} {
	idx := pcas.String(key) % d.SubCount // last N digit(s)
	d.logger.Debug("Get key='%s' idx='%v'", key, idx)
	d.SubDICKs[idx].submux.Lock()
	defer d.SubDICKs[idx].submux.Unlock()
	entry := d.get(idx, key)
	if entry == nil {
		return nil
	}
	retval := entry.value
	return retval // copy avoids race conditions
}

// Set sets the value of a key in the dictionary.
//
// Parameters:
//   - key: the key to set the value for.
//   - value: the value to set.
//
// Returns:
//   - error: an error if the key already exists in the dictionary.
func (d *XDICK) Set(key string, value interface{}) error {
	idx := pcas.String(key) % d.SubCount // last N digit(s)
	d.logger.Debug("Set key='%s' idx='%v'", key, idx)
	d.SubDICKs[idx].submux.Lock()
	defer d.SubDICKs[idx].submux.Unlock()
	entry := d.get(idx, key)
	if entry != nil {
		entry.value = value
		return nil
	}
	retval := d.add(idx, key, value) // copy avoids race conditions
	return retval
}

// Delete deletes an entry from the dictionary.
//
// Parameters:
// - key: the key of the entry to be deleted.
//
// Returns:
// - error: if the entry is not found.
func (d *XDICK) Del(key string) error {
	idx := pcas.String(key) % d.SubCount // last N digit(s)
	d.logger.Debug("Del key='%s' idx='%v'", key, idx)
	d.SubDICKs[idx].submux.Lock()
	defer d.SubDICKs[idx].submux.Unlock()
	dictEntry := d.del(idx, key)
	if dictEntry == nil {
		return fmt.Errorf(`entry not found`)
	}
	d.logger.Debug("deleted key='%s'", key)
	return nil
}
