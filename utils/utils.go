package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func GenerateRandomString(length int) string {

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
} // end func GenerateRandomString

func IsDigit(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
} // end func is Digit

func IsSpace(b byte) bool {
	return b < 33
}

func unprintable(b byte) bool {
	return b < 32
}

func Str2int(str string) int {
	if IsDigit(str) {
		aint, err := strconv.Atoi(str)
		if err == nil {
			return aint
		}
	}
	return 0
} // end func str2int

func Str2int64(str string) int64 {
	if IsDigit(str) {
		aint64, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			return aint64
		}
	}
	return 0
} // end func str2int64

func Str2uint64(str string) uint64 {
	if IsDigit(str) {
		auint64, err := strconv.ParseUint(str, 10, 64)
		if err == nil {
			return auint64
		}
	}
	return 0
} // end func str2int64

func Lines2Bytes(lines []string, delim string) []byte {
	var buf []byte
	for _, line := range lines {
		buf = append(buf, []byte(line+delim)...)
	}
	return buf
} // end func Lines2Bytes

func Bytes2Lines(data []byte, delim string) []string {
	return strings.Split(string(data), delim)
} // end func Bytes2Lines
