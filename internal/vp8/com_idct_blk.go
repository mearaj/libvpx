package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func vp8_dequant_idct_add_y_block_c(q *int16, dq *int16, dst *uint8, stride int, eobs *byte) {
	var (
		i int
		j int
	)
	for i = 0; i < 4; i++ {
		for j = 0; j < 4; j++ {
			if *func() *byte {
				p := &eobs
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() > 1 {
				Vp8DequantIdctAddC(q, dq, dst, stride)
			} else {
				Vp8DcOnlyIdctAddC(int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*0)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(dq), unsafe.Sizeof(int16(0))*0)))), dst, stride, dst, stride)
				libc.MemSet(unsafe.Pointer(q), 0, int(2*unsafe.Sizeof(int16(0))))
			}
			q = (*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*16))
			dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), 4))
		}
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*4-16))
	}
}
func vp8_dequant_idct_add_uv_block_c(q *int16, dq *int16, dst_u *uint8, dst_v *uint8, stride int, eobs *byte) {
	var (
		i int
		j int
	)
	for i = 0; i < 2; i++ {
		for j = 0; j < 2; j++ {
			if *func() *byte {
				p := &eobs
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() > 1 {
				Vp8DequantIdctAddC(q, dq, dst_u, stride)
			} else {
				Vp8DcOnlyIdctAddC(int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*0)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(dq), unsafe.Sizeof(int16(0))*0)))), dst_u, stride, dst_u, stride)
				libc.MemSet(unsafe.Pointer(q), 0, int(2*unsafe.Sizeof(int16(0))))
			}
			q = (*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*16))
			dst_u = (*uint8)(unsafe.Add(unsafe.Pointer(dst_u), 4))
		}
		dst_u = (*uint8)(unsafe.Add(unsafe.Pointer(dst_u), stride*4-8))
	}
	for i = 0; i < 2; i++ {
		for j = 0; j < 2; j++ {
			if *func() *byte {
				p := &eobs
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() > 1 {
				Vp8DequantIdctAddC(q, dq, dst_v, stride)
			} else {
				Vp8DcOnlyIdctAddC(int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*0)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(dq), unsafe.Sizeof(int16(0))*0)))), dst_v, stride, dst_v, stride)
				libc.MemSet(unsafe.Pointer(q), 0, int(2*unsafe.Sizeof(int16(0))))
			}
			q = (*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*16))
			dst_v = (*uint8)(unsafe.Add(unsafe.Pointer(dst_v), 4))
		}
		dst_v = (*uint8)(unsafe.Add(unsafe.Pointer(dst_v), stride*4-8))
	}
}
