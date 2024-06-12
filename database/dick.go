package database

import (
	"encoding/binary"
	"fmt"
	"github.com/dchest/siphash"
	"github.com/go-while/nodare-db/logger"
	//"log"
	"hash/fnv"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	INITIAL_SIZE   = int64(128)
	MAX_SIZE       = 1 << 63
	HASH_siphash   = 0xFF
	HASH_FNV32A    = 0x32A
	HASH_FNV64A    = 0x64A
)

// experiment! MOD can be 10, 100, 1000, 10000
//const DEFAULT_SUBDICKS uint32 = 10
//var AVAIL_SUBDICKS = []uint32{10,100,1000,10000,100000,1000000}
//var USE_SUBDICKS = DEFAULT_SUBDICKS
var key0, key1 uint64
var once sync.Once
var SALT [16]byte
var HASHER = HASH_siphash

type XDICK struct {
	//hashmode     int // mode of hashing
	mainmux  sync.RWMutex  // not used anywhere but exists for whatever reason we may find
	//SubDICKs [DEFAULT_SUBDICKS]*SubDICK
	SubDICKs []*SubDICK
	SubCount uint32
	logger   *ilog.LOG
}

type SubDICK struct {
	parent     *sync.RWMutex
	submux     sync.RWMutex
	hashTables [2]*DickTable
	rehashidx  int
	logger     *ilog.LOG
}

// NewXDICK returns a new instance of XDICK.
//
// The function does not take any parameters.
// It returns a pointer to XDICK.
func NewXDICK(logger *ilog.LOG, sdCh chan uint32, returnsubDICKs chan []*SubDICK) *XDICK {
	var mainmux sync.RWMutex
	//getSubDICKs := make(chan *SubDICK) // unbuffered

	go func(sdCh chan uint32, returnsubDICKs chan []*SubDICK) {
		logger.Info("NewXDICK: creating SubDICKs waits async for configs to load!")
		sub_dicks := <- sdCh // after NewFactory waits for sub_dicks
		sdCh <- sub_dicks // re-push in so NewDICK() can set SubCount
		created := 0
		var subDICKs []*SubDICK
		for i := uint32(0); i < sub_dicks; i++ {
			subDICK := &SubDICK{
				parent:     &mainmux,
				hashTables: [2]*DickTable{NewDickTable(0), NewDickTable(0)},
				rehashidx:  -1,
				logger:     logger,
			}
			subDICKs = append(subDICKs, subDICK)
			created++
		} // end for

		returnsubDICKs <- subDICKs
		close(returnsubDICKs)
		logger.Debug("Created subDICKs %d/%d ", sub_dicks, created)
	}(sdCh, returnsubDICKs) // end go func

	xdick := &XDICK{
		//hashmode: HASHER,
		mainmux:  mainmux,
		SubDICKs: nil, // sets later
		logger:   logger,
	}
	return xdick
}

// mainDICK returns the main hash table of the SubDICK.
//
// No parameters.
// Returns a pointer to a DickTable.
func (d *XDICK) mainDICK(idx uint32) *DickTable {
	return d.SubDICKs[idx].hashTables[0]
}

// rehashingTable returns the DickTable at index 1 of the SubDICK.
//
// No parameters.
// Returns *DickTable.
func (d *XDICK) rehashingTable(idx uint32) *DickTable {
	return d.SubDICKs[idx].hashTables[1]
}

// expand expands the dictionary to a new size if necessary.
//
// newSize: the new size to expand the dictionary to.
// The function does not return anything.
func (d *XDICK) expand(idx uint32, newSize int64) {
	isrehashing := d.isRehashing(idx)
	istablefull := d.mainDICK(idx).used > newSize
	if isrehashing || istablefull {
		d.logger.Info("SubDick [%d] expand return1 ! (newSize=%d isrehashing=%t istablefull=%t used=%d)", idx, newSize, isrehashing, istablefull, d.mainDICK(idx).used)
		return
	}

	d.logger.Info("SubDick [%d] expand newSize=%d isrehashing=%t istablefull=%t used=%d", idx, newSize, isrehashing, istablefull, d.mainDICK(idx).used)

	nextSize := nextPower(newSize)
	if d.mainDICK(idx).used >= nextSize {
		return
	}

	newDickTable := NewDickTable(nextSize)

	if d.mainDICK(idx) == nil || len(d.mainDICK(idx).table) == 0 {
		*d.mainDICK(idx) = *newDickTable
		return
	}

	*d.rehashingTable(idx) = *newDickTable
	d.SubDICKs[idx].rehashidx = 0
}

// expandIfNeeded checks if the dictionary needs to be expanded and performs the expansion if necessary.
//
// No parameters.
// No return values.
func (d *XDICK) expandIfNeeded(idx uint32) {
	if d.isRehashing(idx) {
		return
	}

	if d.mainDICK(idx) == nil || len(d.mainDICK(idx).table) == 0 {
		d.expand(idx, INITIAL_SIZE)
	} else if d.mainDICK(idx).used >= d.mainDICK(idx).size {
		newSize := int64(d.mainDICK(idx).used * 2)
		d.expand(idx, newSize)
	}
}

func (d *XDICK) split(key [16]byte) (uint64, uint64) {
	if len(key) == 0 || len(key) < 16 {
		d.logger.Error("ERROR split len(key)=%d", len(key))
		return 0, 0
	}
	key0 := binary.LittleEndian.Uint64(key[:8])
	key1 := binary.LittleEndian.Uint64(key[8:16])
	return key0, key1
}

// main hasher func
func (d *XDICK) hasher(any string) uint64 {
	switch HASHER {
		case HASH_siphash:
			// sipHashDigest calculates the SipHash-2-4 digest of the given message using the provided key.
			return siphash.Hash(key0, key1, []byte(any))

		case HASH_FNV32A:
			hash := fnv.New32a()
			hash.Write([]byte(any))
			return uint64(hash.Sum32())

		case HASH_FNV64A:
			hash := fnv.New64a()
			hash.Write([]byte(any))
			return hash.Sum64()
	}
	d.logger.Fatal("No HASHER defined! HASHER=%d %x", HASHER, HASHER)
	os.Exit(1)
	return 0
} // end func hasher

// keyIndex returns the index of the given key in the dictionary.
//
// It returns an integer representing the index of the key in the dictionary.
func (d *XDICK) keyIndex(idx uint32, key string) int {
	d.logger.Debug("keyIndex(key=len(%d)='%s'", len(key), key)
	d.expandIfNeeded(idx)
	hash := d.hasher(key) // TODO!FIXME: hash earlier?

	var index int
	loops1 := 0
	loops2 := 0
	for i := 0; i <= 1; i++ {
		loops1++
		hashTable := d.SubDICKs[idx].hashTables[i]
		index = int(hash & hashTable.sizemask)

		for entry := hashTable.table[index]; entry != nil; entry = entry.next {
			loops2++
			if entry.key == key {
				d.logger.Info("keyIndex [%d] entry.key==key='%s' loops1=%d loops2=%d return -1", idx, key, loops1, loops2)
				return -1
			}
		}

		if index == -1 || !d.isRehashing(idx) {
			break
		}
	}
	d.logger.Info("keyIndex [%d] key='%s' loops1=%d loops2=%d return index=%d", idx, key, loops1, loops2, index)
	return index
} // end func keyIndex

// add adds a key-value pair to the SubDICK.
//
// Parameters:
// - key: The key to add.
// - value: The value associated with the key.
//
// Returns:
// - error: An error if the key already exists in the SubDICK.
func (d *XDICK) add(idx uint32, key string, value interface{}) error {
	X := d.keyIndex(idx, key)
	d.logger.Debug("add(key=%d='%s' value='%#v' X=%d", len(key), key, value, X)

	if X == -1 {
		return fmt.Errorf(`unexpectedly found an entry with the same key when trying to add #{ %s } / #{ %s }`, key, value)
	}

	hashTable := d.mainDICK(idx)
	if d.isRehashing(idx) {
		d.rehashStep(idx)
		hashTable = d.mainDICK(idx)
		if d.isRehashing(idx) {
			hashTable = d.rehashingTable(idx)
		}
	}

	entry := hashTable.table[X]

	for entry != nil && entry.key != key {
		entry = entry.next
	}

	if entry == nil {
		entry = NewDickEntry(key, value)
		entry.next = hashTable.table[X]
		hashTable.table[X] = entry
		hashTable.used++
	}

	return nil
} // end func add

// rehashStep returns the result of calling the rehash function on the SubDICK object with an argument of 1.
//
// No parameters.
// Returns an integer.
func (d *XDICK) rehashStep(idx uint32) {
	d.rehash(idx, 1)
}

// rehash rehashes the dictionary with a new size.
//
// n is the new size of the dictionary.
// Returns 0 if the rehashing is not in progress.
// Returns 1 if the rehashing is in progress.
func (d *XDICK) rehash(idx uint32, n int) {
	emptyVisits := n * 10
	if !d.isRehashing(idx) {
		return
	}
	d.logger.Info("SubDick [%d] rehash used=%d", idx, d.mainDICK(idx).used)
	for n > 0 && d.mainDICK(idx).used != 0 {
		n--

		var entry *DickEntry

		for len(d.mainDICK(idx).table) == 0 || d.mainDICK(idx).table[d.SubDICKs[idx].rehashidx] == nil {
			d.SubDICKs[idx].rehashidx++
			emptyVisits--
			if emptyVisits == 0 {
				return
			}
		}

		entry = d.mainDICK(idx).table[d.SubDICKs[idx].rehashidx]

		for entry != nil {
			nextEntry := entry.next
			X := d.hasher(entry.key) & d.rehashingTable(idx).sizemask

			entry.next = d.rehashingTable(idx).table[X]
			d.rehashingTable(idx).table[X] = entry
			d.mainDICK(idx).used--
			d.rehashingTable(idx).used++
			entry = nextEntry
		}

		d.mainDICK(idx).table[d.SubDICKs[idx].rehashidx] = nil
		d.SubDICKs[idx].rehashidx++
	}

	if d.mainDICK(idx).used == 0 {
		d.SubDICKs[idx].hashTables[0] = d.rehashingTable(idx)
		d.SubDICKs[idx].hashTables[1] = NewDickTable(0)
		d.SubDICKs[idx].rehashidx = -1
		return
	}
} // end func rehash

// isRehashing checks if the rehash index of the SubDICK struct is not equal to -1.
//
// It does not take any parameters.
// It returns a boolean value indicating whether the rehash index is not equal to -1.
func (d *XDICK) isRehashing(idx uint32) bool {
	return d.SubDICKs[idx].rehashidx != -1
}

// get returns the DickEntry associated with the given key in the SubDICK.
//
// Parameters:
// - key: the key to search for in the SubDICK.
//
// Return:
// - *DickEntry: the DickEntry associated with the given key, or nil if not found.
func (d *XDICK) get(idx uint32, key string) *DickEntry {
	if d.mainDICK(idx).used == 0 && d.rehashingTable(idx).used == 0 {
		return nil
	}

	hash := d.hasher(key) // TODO!FIXME: hash earlier?

	for ind, hashTable := range []*DickTable{d.mainDICK(idx), d.rehashingTable(idx)} {
		if hashTable == nil || len(hashTable.table) == 0 || (ind == 1 && !d.isRehashing(idx)) {
			continue
		}

		index := hash & hashTable.sizemask
		entry := hashTable.table[index]

		for entry != nil {
			if entry.key == key {
				return entry
			}
			entry = entry.next
		}
	}

	return nil
} // end func get

// delete deletes a key from the dictionary and returns the corresponding value.
//
// Parameters:
// - key: the key to be deleted from the SubDICK.
//
// Return:
// - *DickEntry: the deleted DickEntry if found, otherwise nil.
func (d *XDICK) del(idx uint32, key string) *DickEntry {

	if d.mainDICK(idx).used == 0 && d.rehashingTable(idx).used == 0 {
		return nil
	}

	if d.isRehashing(idx) {
		d.rehashStep(idx)
	}

	hash := d.hasher(key) // TODO!FIXME: hash earlier!

	for i, hashTable := range []*DickTable{d.mainDICK(idx), d.rehashingTable(idx)} {
		if hashTable == nil || (i == 1 && !d.isRehashing(idx)) {
			continue
		}
		index := hash & hashTable.sizemask
		entry := hashTable.table[index]
		var previousEntry *DickEntry

		for entry != nil {
			if entry.key == key {
				if previousEntry != nil {
					previousEntry.next = entry.next
				} else {
					hashTable.table[index] = entry.next
				}
				hashTable.used--
				return entry
			}
			previousEntry = entry
			entry = entry.next
		}
	}

	return nil
} // end func del

func (d *XDICK) watchDog(idx uint32) {
	//log.Printf("Booted Watchdog [%d]", idx)

	for {
		time.Sleep(time.Second)  // TODO!FIXME setting: watchdog_timer

		if d == nil || d.SubCount == 0 || d.SubDICKs[idx] == nil {
			// not finished booting
			continue
		}

		if !d.SubDICKs[idx].logger.IfDebug() {
			time.Sleep(59 * time.Second)
			continue
		}

		// print some statistics
		d.SubDICKs[idx].submux.RLock()
		//ht := len(d.SubDICKs[idx].hashTables), // is always 2
		ht0 := d.SubDICKs[idx].hashTables[0].used
		ht1 := d.SubDICKs[idx].hashTables[1].used
		if ht0 == 0 && ht1 == 0 {
			// subdick is empty
			d.SubDICKs[idx].submux.RUnlock()
			continue
		}
		ht0cap := len(d.SubDICKs[idx].hashTables[0].table)
		ht1cap := len(d.SubDICKs[idx].hashTables[1].table)
		d.logger.Info("watchDog [%d] ht0=%d/%d ht1=%d/%d", idx, ht0, ht0cap, ht1, ht1cap)
		d.SubDICKs[idx].submux.RUnlock()
		//d.logger.Debug("watchDog [%d] SubDICKs='\n   ---> %#v", idx, d.SubDICKs[idx])
	}
} // end func watchDog


// GenerateSALT generates a fixed slice of 16 maybe NOT really random bytes
func (d *XDICK) GenerateSALT() {
	once.Do(func() {
		rand.Seed(time.Now().UnixNano())
		cs := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		for i := 0; i < 16; i++ {
			SALT[i] = cs[rand.Intn(len(cs))]
		}
		key0, key1 = d.split(SALT)
	})
}

// nextPower calculates the next power of 2 greater than the given size.
//
// Parameters:
// - size: the size for which we want to find the next power of 2.
//
// Return type:
// - int64: the next power of 2 greater than the given size.
func nextPower(size int64) int64 {
	if size <= INITIAL_SIZE {
		return INITIAL_SIZE
	}

	size--
	size |= size >> 1
	size |= size >> 2
	size |= size >> 4
	size |= size >> 8
	size |= size >> 16
	size |= size >> 32

	return size + 1
}
