package vp8

import (
	"github.com/mearaj/libvpx/internal/scale"
	"unsafe"
)

func vp8_swap_yv12_buffer(new_frame *scale.Yv12BufferConfig, last_frame *scale.Yv12BufferConfig) {
	var temp *uint8
	temp = (*uint8)(unsafe.Pointer(last_frame.Buffer_alloc))
	last_frame.Buffer_alloc = new_frame.Buffer_alloc
	new_frame.Buffer_alloc = (*uint8)(unsafe.Pointer(temp))
	temp = (*uint8)(unsafe.Pointer(last_frame.Y_buffer))
	last_frame.Y_buffer = new_frame.Y_buffer
	new_frame.Y_buffer = (*uint8)(unsafe.Pointer(temp))
	temp = (*uint8)(unsafe.Pointer(last_frame.U_buffer))
	last_frame.U_buffer = new_frame.U_buffer
	new_frame.U_buffer = (*uint8)(unsafe.Pointer(temp))
	temp = (*uint8)(unsafe.Pointer(last_frame.V_buffer))
	last_frame.V_buffer = new_frame.V_buffer
	new_frame.V_buffer = (*uint8)(unsafe.Pointer(temp))
}
