package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func setup_intra_recon_left(y_buffer *uint8, u_buffer *uint8, v_buffer *uint8, y_stride int, uv_stride int) {
	var i int
	for i = 0; i < 16; i++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(y_buffer), y_stride*i)) = 129
	}
	for i = 0; i < 8; i++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(u_buffer), uv_stride*i)) = 129
	}
	for i = 0; i < 8; i++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(v_buffer), uv_stride*i)) = 129
	}
}
func vp8_setup_intra_recon(ybf *scale.Yv12BufferConfig) {
	var i int
	libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ybf.Y_buffer), -1))), -ybf.Y_stride), math.MaxInt8, ybf.Y_width+5)
	for i = 0; i < ybf.Y_height; i++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(ybf.Y_buffer), ybf.Y_stride*i-1)) = 129
	}
	libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ybf.U_buffer), -1))), -ybf.Uv_stride), math.MaxInt8, ybf.Uv_width+5)
	for i = 0; i < ybf.Uv_height; i++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(ybf.U_buffer), ybf.Uv_stride*i-1)) = 129
	}
	libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ybf.V_buffer), -1))), -ybf.Uv_stride), math.MaxInt8, ybf.Uv_width+5)
	for i = 0; i < ybf.Uv_height; i++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(ybf.V_buffer), ybf.Uv_stride*i-1)) = 129
	}
}
func vp8_setup_intra_recon_top_line(ybf *scale.Yv12BufferConfig) {
	libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ybf.Y_buffer), -1))), -ybf.Y_stride), math.MaxInt8, ybf.Y_width+5)
	libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ybf.U_buffer), -1))), -ybf.Uv_stride), math.MaxInt8, ybf.Uv_width+5)
	libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ybf.V_buffer), -1))), -ybf.Uv_stride), math.MaxInt8, ybf.Uv_width+5)
}
