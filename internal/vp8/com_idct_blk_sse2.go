package vp8

import "unsafe"

func vp8_dequant_idct_add_y_block_sse2(q *int16, dq *int16, dst *uint8, stride int, eobs *byte) {
	var i int
	for i = 0; i < 4; i++ {
		if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*0))) != 0 {
			if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*0)))&0xFEFE != 0 {
				vp8_idct_dequant_full_2x_sse2(q, dq, dst, stride)
			} else {
				vp8_idct_dequant_0_2x_sse2(q, dq, dst, stride)
			}
		}
		if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*1))) != 0 {
			if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*1)))&0xFEFE != 0 {
				vp8_idct_dequant_full_2x_sse2((*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*32)), dq, (*uint8)(unsafe.Add(unsafe.Pointer(dst), 8)), stride)
			} else {
				vp8_idct_dequant_0_2x_sse2((*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*32)), dq, (*uint8)(unsafe.Add(unsafe.Pointer(dst), 8)), stride)
			}
		}
		q = (*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*64))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*4))
		eobs = (*byte)(unsafe.Add(unsafe.Pointer(eobs), 4))
	}
}
func vp8_dequant_idct_add_uv_block_sse2(q *int16, dq *int16, dst_u *uint8, dst_v *uint8, stride int, eobs *byte) {
	if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*0))) != 0 {
		if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*0)))&0xFEFE != 0 {
			vp8_idct_dequant_full_2x_sse2(q, dq, dst_u, stride)
		} else {
			vp8_idct_dequant_0_2x_sse2(q, dq, dst_u, stride)
		}
	}
	q = (*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*32))
	dst_u = (*uint8)(unsafe.Add(unsafe.Pointer(dst_u), stride*4))
	if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*1))) != 0 {
		if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*1)))&0xFEFE != 0 {
			vp8_idct_dequant_full_2x_sse2(q, dq, dst_u, stride)
		} else {
			vp8_idct_dequant_0_2x_sse2(q, dq, dst_u, stride)
		}
	}
	q = (*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*32))
	if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*2))) != 0 {
		if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*2)))&0xFEFE != 0 {
			vp8_idct_dequant_full_2x_sse2(q, dq, dst_v, stride)
		} else {
			vp8_idct_dequant_0_2x_sse2(q, dq, dst_v, stride)
		}
	}
	q = (*int16)(unsafe.Add(unsafe.Pointer(q), unsafe.Sizeof(int16(0))*32))
	dst_v = (*uint8)(unsafe.Add(unsafe.Pointer(dst_v), stride*4))
	if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*3))) != 0 {
		if int(*(*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Pointer(eobs))), unsafe.Sizeof(int16(0))*3)))&0xFEFE != 0 {
			vp8_idct_dequant_full_2x_sse2(q, dq, dst_v, stride)
		} else {
			vp8_idct_dequant_0_2x_sse2(q, dq, dst_v, stride)
		}
	}
}
