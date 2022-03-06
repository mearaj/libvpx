package dsp

import "unsafe"

type vpx_rb_error_handler func(data unsafe.Pointer)
type vpx_read_bit_buffer struct {
	Bit_buffer         *uint8
	Bit_buffer_end     *uint8
	Bit_offset         uint64
	Error_handler_data unsafe.Pointer
	Error_handler      vpx_rb_error_handler
}

func vpx_rb_bytes_read(rb *vpx_read_bit_buffer) uint64 {
	return (rb.Bit_offset + 7) >> 3
}
func vpx_rb_read_bit(rb *vpx_read_bit_buffer) int {
	var (
		off uint64 = rb.Bit_offset
		p   uint64 = off >> 3
		q   int    = 7 - int(off&7)
	)
	if uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(rb.Bit_buffer), p)))) < uintptr(unsafe.Pointer(rb.Bit_buffer_end)) {
		var bit int = (int(*(*uint8)(unsafe.Add(unsafe.Pointer(rb.Bit_buffer), p))) >> q) & 1
		rb.Bit_offset = off + 1
		return bit
	} else {
		if rb.Error_handler != nil {
			rb.Error_handler(rb.Error_handler_data)
		}
		return 0
	}
}
func vpx_rb_read_literal(rb *vpx_read_bit_buffer, bits int) int {
	var (
		value int = 0
		bit   int
	)
	for bit = bits - 1; bit >= 0; bit-- {
		value |= vpx_rb_read_bit(rb) << bit
	}
	return value
}
func vpx_rb_read_signed_literal(rb *vpx_read_bit_buffer, bits int) int {
	var value int = vpx_rb_read_literal(rb, bits)
	if vpx_rb_read_bit(rb) != 0 {
		return -value
	}
	return value
}
func vpx_rb_read_inv_signed_literal(rb *vpx_read_bit_buffer, bits int) int {
	return vpx_rb_read_signed_literal(rb, bits)
}
