package util

import "math"

const LOCAL_GCC_VERSION = 0

func BSwap16(x uint16_t) uint16_t {
	return (x >> 8) | (x&math.MaxUint8)<<8
}
func BSwap32(x uint32_t) uint32_t {
	return (x >> 24) | ((x >> 8) & 0xFF00) | ((x << 8) & 0xFF0000) | x<<24
}
func BSwap64(x uint64) uint64 {
	x = ((x & 0xFFFFFFFF00000000) >> 32) | (x&math.MaxUint32)<<32
	x = ((x & 0xFFFF0000FFFF0000) >> 16) | (x&0xFFFF0000FFFF)<<16
	x = ((x & 0xFF00FF00FF00FF00) >> 8) | (x&0xFF00FF00FF00FF)<<8
	return x
}
