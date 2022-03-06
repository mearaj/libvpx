package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func vp8_dequantize_b_c(d *Blockd, DQC *int16) {
	var (
		i  int
		DQ *int16 = d.Dqcoeff
		Q  *int16 = d.Qcoeff
	)
	for i = 0; i < 16; i++ {
		*(*int16)(unsafe.Add(unsafe.Pointer(DQ), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(Q), unsafe.Sizeof(int16(0))*uintptr(i)))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(DQC), unsafe.Sizeof(int16(0))*uintptr(i)))))
	}
}
func Vp8DequantIdctAddC(input *int16, dq *int16, dest *uint8, stride int) {
	var i int
	for i = 0; i < 16; i++ {
		*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(dq), unsafe.Sizeof(int16(0))*uintptr(i)))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(i)))))
	}
	vp8_short_idct4x4llm_c(input, dest, stride, dest, stride)
	libc.MemSet(unsafe.Pointer(input), 0, 32)
}
