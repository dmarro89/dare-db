// Hash defines a hash interface for objects.

/*
This software is distributed under an MIT license. You should have received a copy of the license along with this software. If not, see <https://opensource.org/licenses/MIT>.
*/

package hash

import (
	"math"
)

/////////////////////////////////////////////////////////////////////////
// Public functions
/////////////////////////////////////////////////////////////////////////

// Combine returns a new hash value given by combining h1 and h2.
func Combine(h1 uint32, h2 uint32) uint32 {
	return h1 ^ (h2 + 0x9e3779b9 + (h1 << 6) + (h1 >> 2))
}

// Bool returns the hash of the given boolean.
func Bool(c bool) uint32 {
	if c {
		return 1
	}
	return 0
}

// BoolSlice returns the hash of the given slice.
func BoolSlice(S []bool) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Bool(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Bool(S[i]))
	}
	return h
}

// Int8 returns the hash of the given integer.
func Int8(c int8) uint32 {
	return uint32(c)
}

// Int8Slice returns the hash of the given slice.
func Int8Slice(S []int8) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Int8(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Int8(S[i]))
	}
	return h
}

// Int16 returns the hash of the given integer.
func Int16(c int16) uint32 {
	return uint32(c)
}

// Int16Slice returns the hash of the given slice.
func Int16Slice(S []int16) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Int16(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Int16(S[i]))
	}
	return h
}

// Int32 returns the hash of the given integer.
func Int32(c int32) uint32 {
	return uint32(c)
}

// Int32Slice returns the hash of the given slice.
func Int32Slice(S []int32) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Int32(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Int32(S[i]))
	}
	return h
}

// Int64 returns the hash of the given integer.
func Int64(c int64) uint32 {
	return Combine(uint32(c), uint32(c>>32))
}

// Int64Slice returns the hash of the given slice.
func Int64Slice(S []int64) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Int64(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Int64(S[i]))
	}
	return h
}

// Int returns the hash of the given integer.
func Int(c int) uint32 {
	return Int64(int64(c))
}

// IntSlice returns the hash of the given slice.
func IntSlice(S []int) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Int(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Int(S[i]))
	}
	return h
}

// Uint8 returns the hash of the given unsigned integer.
func Uint8(c uint8) uint32 {
	return uint32(c)
}

// Uint8Slice returns the hash of the given slice.
func Uint8Slice(S []uint8) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Uint8(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Uint8(S[i]))
	}
	return h
}

// Uint16 returns the hash of the given unsigned integer.
func Uint16(c uint16) uint32 {
	return uint32(c)
}

// Uint16Slice returns the hash of the given slice.
func Uint16Slice(S []uint16) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Uint16(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Uint16(S[i]))
	}
	return h
}

// Uint32 returns the hash of the given unsigned integer.
func Uint32(c uint32) uint32 {
	return c
}

// Uint32Slice returns the hash of the given slice.
func Uint32Slice(S []uint32) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Uint32(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Uint32(S[i]))
	}
	return h
}

// Uint64 returns the hash of the given unsigned integer.
func Uint64(c uint64) uint32 {
	return Combine(uint32(c), uint32(c>>32))
}

// Uint64Slice returns the hash of the given slice.
func Uint64Slice(S []uint64) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Uint64(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Uint64(S[i]))
	}
	return h
}

// Uint returns the hash of the given unsigned integer.
func Uint(c uint) uint32 {
	return Uint64(uint64(c))
}

// UintSlice returns the hash of the given slice.
func UintSlice(S []uint) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Uint(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Uint(S[i]))
	}
	return h
}

// Float32 returns the hash of the given float.
func Float32(f float32) uint32 {
	return Uint32(math.Float32bits(f))
}

// Float32Slice returns the hash of the given slice.
func Float32Slice(S []float32) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Float32(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Float32(S[i]))
	}
	return h
}

// Float64 returns the hash of the given float.
func Float64(f float64) uint32 {
	return Uint64(math.Float64bits(f))
}

// Float64Slice returns the hash of the given slice.
func Float64Slice(S []float64) uint32 {
	n := len(S)
	if n == 0 {
		return 0
	}
	h := Float64(S[0])
	for i := 1; i < n; i++ {
		h = Combine(h, Float64(S[i]))
	}
	return h
}

// String returns the hash of the given string.
func String(s string) uint32 {
	return ByteSlice([]byte(s))
}

// StringSlice returns the hash of the given slice.
func StringSlice(S []string) uint32 {
	if len(S) == 0 {
		return 0
	}
	h := New()
	for _, s := range S {
		h.Write([]byte(s))
	}
	v := h.Hash()
	Reuse(h)
	return v
}

// Byte returns the hash of the given byte.
func Byte(c byte) uint32 {
	return uint32(c)
}

// ByteSlice returns the hash of the given slice.
func ByteSlice(S []byte) uint32 {
	switch len(S) {
	case 0:
		return 0
	case 1:
		return uint32(S[0])
	case 2:
		return uint32(S[0]) | (uint32(S[1]) << 8)
	case 3:
		return uint32(S[0]) | (uint32(S[1]) << 8) | (uint32(S[2]) << 16)
	case 4:
		return uint32(S[0]) | (uint32(S[1]) << 8) | (uint32(S[2]) << 16) | (uint32(S[3]) << 24)
	}
	h := New()
	h.Write(S)
	v := h.Hash()
	Reuse(h)
	return v
}
