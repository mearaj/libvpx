package vp8

import (
	"github.com/mearaj/libvpx/internal/vpx"
	"math"
	"unsafe"
)

const VP8_LOTS_OF_BITS = 0x40000000

func vp8dx_decode_bool(br *BOOL_DECODER, probability int) int {
	var (
		bit      uint = 0
		value    uint64
		split    uint
		bigsplit uint64
		count    int
		range_   uint
	)
	split = (((br.Range - 1) * uint(probability)) >> 8) + 1
	if br.Count < 0 {
		vp8dx_bool_decoder_fill(br)
	}
	value = br.Value
	count = br.Count
	bigsplit = uint64(split) << uint64((CHAR_BIT*int(unsafe.Sizeof(uint64(0))))-8)
	range_ = split
	if value >= bigsplit {
		range_ = br.Range - split
		value = value - bigsplit
		bit = 1
	}
	{
		var shift uint8 = vp8_norm[uint8(range_)]
		range_ <<= uint(shift)
		value <<= uint64(shift)
		count -= int(shift)
	}
	br.Value = value
	br.Count = count
	br.Range = range_
	return int(bit)
}
func vp8_decode_value(br *BOOL_DECODER, bits int) int {
	var (
		z   int = 0
		bit int
	)
	for bit = bits - 1; bit >= 0; bit-- {
		z |= vp8dx_decode_bool(br, 128) << bit
	}
	return z
}
func vp8dx_bool_error(br *BOOL_DECODER) int {
	if br.Count > (CHAR_BIT*int(unsafe.Sizeof(uint64(0)))) && br.Count < 0x40000000 {
		return 1
	}
	return 0
}
func vp8dx_start_decode(br *BOOL_DECODER, source *uint8, source_sz uint, decrypt_cb vpx.DecryptCb, decrypt_state unsafe.Pointer) int {
	if source_sz != 0 && source == nil {
		return 1
	}
	if source != nil {
		br.User_buffer_end = (*uint8)(unsafe.Add(unsafe.Pointer(source), source_sz))
	} else {
		br.User_buffer_end = source
	}
	br.User_buffer = source
	br.Value = 0
	br.Count = -8
	br.Range = math.MaxUint8
	br.Decrypt_cb = decrypt_cb
	br.Decrypt_state = decrypt_state
	vp8dx_bool_decoder_fill(br)
	return 0
}
func vp8dx_bool_decoder_fill(br *BOOL_DECODER) {
	var (
		bufptr     *uint8 = br.User_buffer
		value      uint64 = br.Value
		count      int    = br.Count
		shift      int    = (CHAR_BIT * int(unsafe.Sizeof(uint64(0)))) - CHAR_BIT - (count + CHAR_BIT)
		bytes_left uint64 = uint64(int64(uintptr(unsafe.Pointer(br.User_buffer_end)) - uintptr(unsafe.Pointer(bufptr))))
		bits_left  uint64 = bytes_left * CHAR_BIT
		x          int    = shift + CHAR_BIT - int(bits_left)
		loop_end   int    = 0
		decrypted  [9]uint8
	)
	if br.Decrypt_cb != nil {
		var n uint64 = uint64(func() uintptr {
			if (unsafe.Sizeof([9]uint8{})) < uintptr(bytes_left) {
				return unsafe.Sizeof([9]uint8{})
			}
			return uintptr(bytes_left)
		}())
		br.Decrypt_cb(br.Decrypt_state, bufptr, &decrypted[0], int(n))
		bufptr = &decrypted[0]
	}
	if x >= 0 {
		count += 0x40000000
		loop_end = x
	}
	if x < 0 || bits_left != 0 {
		for shift >= loop_end {
			count += CHAR_BIT
			value |= uint64(*bufptr) << uint64(shift)
			bufptr = (*uint8)(unsafe.Add(unsafe.Pointer(bufptr), 1))
			br.User_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(br.User_buffer), 1))
			shift -= CHAR_BIT
		}
	}
	br.Value = value
	br.Count = count
}
