package dsp

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

type vpx_write_bit_buffer struct {
	Bit_buffer *uint8
	Bit_offset uint64
}

func vpx_wb_bytes_written(wb *vpx_write_bit_buffer) uint64 {
	return wb.Bit_offset/CHAR_BIT + uint64(libc.BoolToInt(wb.Bit_offset%CHAR_BIT > 0))
}
func vpx_wb_write_bit(wb *vpx_write_bit_buffer, bit int) {
	var (
		off int = int(wb.Bit_offset)
		p   int = off / CHAR_BIT
		q   int = int(CHAR_BIT-1) - off%CHAR_BIT
	)
	if q == int(CHAR_BIT-1) {
		*(*uint8)(unsafe.Add(unsafe.Pointer(wb.Bit_buffer), p)) = uint8(int8(bit << q))
	} else {
		*(*uint8)(unsafe.Add(unsafe.Pointer(wb.Bit_buffer), p)) &= uint8(int8(^(1 << q)))
		*(*uint8)(unsafe.Add(unsafe.Pointer(wb.Bit_buffer), p)) |= uint8(int8(bit << q))
	}
	wb.Bit_offset = uint64(off + 1)
}
func vpx_wb_write_literal(wb *vpx_write_bit_buffer, data int, bits int) {
	var bit int
	for bit = bits - 1; bit >= 0; bit-- {
		vpx_wb_write_bit(wb, (data>>bit)&1)
	}
}
func vpx_wb_write_inv_signed_literal(wb *vpx_write_bit_buffer, data int, bits int) {
	vpx_wb_write_literal(wb, int(cmath.Abs(int64(data))), bits)
	vpx_wb_write_bit(wb, int(libc.BoolToInt(data < 0)))
}
