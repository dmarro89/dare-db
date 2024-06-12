package database

type DickEntry struct {
	next  *DickEntry
	key   string
	value interface{}
}

// NewDickEntry creates a new DickEntry with the given key and value.
//
// Parameters:
// - key: a string representing the key of the entry.
// - value: an interface{} representing the value of the entry.
//
// Returns:
// - *DickEntry: a pointer to the newly created DickEntry.
func NewDickEntry(key string, value interface{}) *DickEntry {
	return &DickEntry{
		key:   key,
		value: value,
		next:  nil,
	}
}
