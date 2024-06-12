package database

type DickTable struct {
	table    []*DickEntry
	size     int64
	sizemask uint64
	used     int64
}

// NewDickTable creates a new DickTable with the specified size.
//
// Parameters:
// - size: the size of the DickTable.
//
// Returns:
// - *DickTable: a pointer to the newly created DickTable.
func NewDickTable(size int64) *DickTable {
	var sizemask uint64
	table := []*DickEntry{}

	if size > 0 {
		table = make([]*DickEntry, size)
		sizemask = uint64(size - 1)
	}

	return &DickTable{
		table:    table,
		size:     size,
		sizemask: sizemask,
	}
}

// empty checks if the hash table is empty.
//
// Returns true if the hash table is empty, false otherwise.
func (ht *DickTable) empty() bool {
	return ht.size == 0
}
