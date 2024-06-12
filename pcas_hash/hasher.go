// Sequence implments a simple hasher for generating hash values for sequence of objects.

/*
This software is distributed under an MIT license. You should have received a copy of the license along with this software. If not, see <https://opensource.org/licenses/MIT>.
*/

package hash

import (
	"sync"
)

// pool is a pool of hashers available for (re)use.
var pool = &sync.Pool{
	New: func() interface{} { return &Hash{} },
}

// Hash defines a hasher. A zero-valued Hash is ready for use. Note that this is NOT intended for cryptographic use.
type Hash struct {
	h        uint32  // The current hash value
	idx      int     // The index of the next entry in the bytes buffer
	buf      [4]byte // The bytes buffer
	hasFirst bool    // Has the initial hash been computed?
}

/////////////////////////////////////////////////////////////////////////
// Hash functions
/////////////////////////////////////////////////////////////////////////

// New returns a new Hash from the pool ready for use.
func New() *Hash {
	h := pool.Get().(*Hash)
	h.Reset()
	return h
}

// Reuse returns the given Hash to the pool ready for future use.
func Reuse(h *Hash) {
	if h != nil {
		pool.Put(h)
	}
}

// Write appends the given bytes to the hash. Write never returns an error.
func (h *Hash) Write(b []byte) (int, error) {
	for _, c := range b {
		if h.idx < 4 {
			h.buf[h.idx] = c
			h.idx++
		} else {
			v := uint32(h.buf[0]) | (uint32(h.buf[1]) << 8) | (uint32(h.buf[2]) << 16) | (uint32(h.buf[3]) << 24)
			if h.hasFirst {
				h.h = Combine(h.h, v)
			} else {
				h.h = v
				h.hasFirst = true
			}
			h.idx = 0
		}
	}
	return len(b), nil
}

// Sum appends the current hash to b and returns the resulting slice. It does not change the underlying hash state.
func (h *Hash) Sum(b []byte) []byte {
	var v uint32
	if len(b) == 0 {
		v = h.Hash()
	} else {
		h2 := *h
		h2.Write(b)
		v = h2.Hash()
	}
	return []byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)}
}

// Sum32 returns the current hash value.
func (h *Hash) Sum32() uint32 {
	return uint32(h.Hash())
}

// Hash returns the current hash value.
func (h *Hash) Hash() uint32 {
	if h == nil {
		return 0
	}
	v := h.h
	if h.idx != 0 {
		var vv uint32
		switch h.idx {
		case 1:
			vv = uint32(h.buf[0])
		case 2:
			vv = uint32(h.buf[0]) | (uint32(h.buf[1]) << 8)
		case 3:
			vv = uint32(h.buf[0]) | (uint32(h.buf[1]) << 8) | (uint32(h.buf[2]) << 16)
		case 4:
			vv = uint32(h.buf[0]) | (uint32(h.buf[1]) << 8) | (uint32(h.buf[2]) << 16) | (uint32(h.buf[3]) << 24)
		}
		if h.hasFirst {
			v = Combine(v, vv)
		} else {
			v = vv
		}
	}
	return v
}

// Reset resets the Hash to its initial state.
func (h *Hash) Reset() {
	if h != nil {
		h.h = 0
		h.idx = 0
		h.hasFirst = false
	}
}

// Size returns the number of bytes Sum will return, which is 4.
func (*Hash) Size() int {
	return 4
}

// BlockSize is irrelevant for a Hash. This is implemented purely to satisfy the hash.Hash interface, and will always return 512.
func (*Hash) BlockSize() int {
	return 512
}
