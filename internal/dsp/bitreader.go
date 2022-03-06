package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const LOTS_OF_BITS = 0x40000000

type BD_VALUE uint64
type vpx_reader struct {
	Value         BD_VALUE
	Range         uint
	Count         int
	Buffer_end    *uint8
	Buffer        *uint8
	Decrypt_cb    vpx_decrypt_cb
	Decrypt_state unsafe.Pointer
	Clear_buffer  [9]uint8
}

func vpx_reader_has_error(r *vpx_reader) int {
	return int(libc.BoolToInt(r.Count > (CHAR_BIT*int(unsafe.Sizeof(BD_VALUE(0)))) && r.Count < LOTS_OF_BITS))
}
func vpx_read(r *vpx_reader, prob int) int {
	var (
		bit      uint = 0
		value    BD_VALUE
		bigsplit BD_VALUE
		count    int
		range_   uint
		split    uint = (r.Range*uint(prob) + uint(256-prob)) >> CHAR_BIT
	)
	if r.Count < 0 {
		vpx_reader_fill(r)
	}
	value = r.Value
	count = r.Count
	bigsplit = BD_VALUE(split) << BD_VALUE((CHAR_BIT*int(unsafe.Sizeof(BD_VALUE(0))))-CHAR_BIT)
	range_ = split
	if value >= bigsplit {
		range_ = r.Range - split
		value = value - bigsplit
		bit = 1
	}
	{
		var shift uint8 = uint8(vpx_norm[uint8(range_)])
		range_ <<= uint(shift)
		value <<= BD_VALUE(shift)
		count -= int(shift)
	}
	r.Value = value
	r.Count = count
	r.Range = range_
	return int(bit)
}
func vpx_read_bit(r *vpx_reader) int {
	return vpx_read(r, 128)
}
func vpx_read_literal(r *vpx_reader, bits int) int {
	var (
		literal int = 0
		bit     int
	)
	for bit = bits - 1; bit >= 0; bit-- {
		literal |= vpx_read_bit(r) << bit
	}
	return literal
}
func vpx_read_tree(r *vpx_reader, tree *int8, probs *uint8) int {
	var i int8 = 0
	for (func() int8 {
		i = *(*int8)(unsafe.Add(unsafe.Pointer(tree), int(i)+vpx_read(r, int(*(*uint8)(unsafe.Add(unsafe.Pointer(probs), i>>1))))))
		return i
	}()) > 0 {
		continue
	}
	return int(-i)
}
func vpx_reader_init(r *vpx_reader, buffer *uint8, size uint64, decrypt_cb vpx_decrypt_cb, decrypt_state unsafe.Pointer) int {
	if size != 0 && buffer == nil {
		return 1
	} else {
		r.Buffer_end = (*uint8)(unsafe.Add(unsafe.Pointer(buffer), size))
		r.Buffer = buffer
		r.Value = 0
		r.Count = -8
		r.Range = math.MaxUint8
		r.Decrypt_cb = decrypt_cb
		r.Decrypt_state = decrypt_state
		vpx_reader_fill(r)
		return int(libc.BoolToInt(vpx_read_bit(r) != 0))
	}
}
func vpx_reader_fill(r *vpx_reader) {
	var (
		buffer_end   *uint8   = r.Buffer_end
		buffer       *uint8   = r.Buffer
		buffer_start *uint8   = buffer
		value        BD_VALUE = r.Value
		count        int      = r.Count
		bytes_left   uint64   = uint64(int64(uintptr(unsafe.Pointer(buffer_end)) - uintptr(unsafe.Pointer(buffer))))
		bits_left    uint64   = bytes_left * CHAR_BIT
		shift        int      = (CHAR_BIT * int(unsafe.Sizeof(BD_VALUE(0)))) - CHAR_BIT - (count + CHAR_BIT)
	)
	if r.Decrypt_cb != nil {
		var n uint64 = uint64(func() uintptr {
			if (unsafe.Sizeof([9]uint8{})) < uintptr(bytes_left) {
				return unsafe.Sizeof([9]uint8{})
			}
			return uintptr(bytes_left)
		}())
		r.Decrypt_cb(r.Decrypt_state, (*uint8)(unsafe.Pointer(buffer)), (*uint8)(unsafe.Pointer(&r.Clear_buffer[0])), int(n))
		buffer = &r.Clear_buffer[0]
		buffer_start = &r.Clear_buffer[0]
	}
	if bits_left > uint64(CHAR_BIT*int(unsafe.Sizeof(BD_VALUE(0)))) {
		var (
			bits              int = (shift & 0xFFFFFFF8) + CHAR_BIT
			nv                BD_VALUE
			big_endian_values BD_VALUE
		)
		libc.MemCpy(unsafe.Pointer(&big_endian_values), unsafe.Pointer(buffer), int(unsafe.Sizeof(BD_VALUE(0))))
		big_endian_values = BD_VALUE(BSwap32(uint32(big_endian_values)))
		nv = big_endian_values >> BD_VALUE((CHAR_BIT*int(unsafe.Sizeof(BD_VALUE(0))))-bits)
		count += bits
		buffer = (*uint8)(unsafe.Add(unsafe.Pointer(buffer), bits>>3))
		value = r.Value | nv<<BD_VALUE(shift&7)
	} else {
		var (
			bits_over int = (shift + CHAR_BIT - int(bits_left))
			loop_end  int = 0
		)
		if bits_over >= 0 {
			count += LOTS_OF_BITS
			loop_end = bits_over
		}
		if bits_over < 0 || bits_left != 0 {
			for shift >= loop_end {
				count += CHAR_BIT
				value |= BD_VALUE(*func() *uint8 {
					p := &buffer
					x := *p
					*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}()) << BD_VALUE(shift)
				shift -= CHAR_BIT
			}
		}
	}
	r.Buffer = (*uint8)(unsafe.Add(unsafe.Pointer(r.Buffer), int64(uintptr(unsafe.Pointer(buffer))-uintptr(unsafe.Pointer(buffer_start)))))
	r.Value = value
	r.Count = count
}
func vpx_reader_find_end(r *vpx_reader) *uint8 {
	for r.Count > CHAR_BIT && r.Count < (CHAR_BIT*int(unsafe.Sizeof(BD_VALUE(0)))) {
		r.Count -= CHAR_BIT
		r.Buffer = (*uint8)(unsafe.Add(unsafe.Pointer(r.Buffer), -1))
	}
	return r.Buffer
}
