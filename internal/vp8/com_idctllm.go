package vp8

import (
	"math"
	"unsafe"
)

var cospi8sqrt2minus1 int = 0x4E7B
var sinpi8sqrt2 int = 0x8A8C

func vp8_short_idct4x4llm_c(input *int16, pred_ptr *uint8, pred_stride int, dst_ptr *uint8, dst_stride int) {
	var (
		i          int
		r          int
		c          int
		a1         int
		b1         int
		c1         int
		d1         int
		output     [16]int16
		ip         *int16 = input
		op         *int16 = &output[0]
		temp1      int
		temp2      int
		shortpitch int = 4
	)
	for i = 0; i < 4; i++ {
		a1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*8)))
		b1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*8)))
		temp1 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*4))) * sinpi8sqrt2) >> 16
		temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*12))) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*12))) * cospi8sqrt2minus1) >> 16)
		c1 = temp1 - temp2
		temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*4))) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*4))) * cospi8sqrt2minus1) >> 16)
		temp2 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*12))) * sinpi8sqrt2) >> 16
		d1 = temp1 + temp2
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*uintptr(shortpitch*0))) = int16(a1 + d1)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*uintptr(shortpitch*3))) = int16(a1 - d1)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*uintptr(shortpitch*1))) = int16(b1 + c1)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*uintptr(shortpitch*2))) = int16(b1 - c1)
		ip = (*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))
		op = (*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*1))
	}
	ip = &output[0]
	op = &output[0]
	for i = 0; i < 4; i++ {
		a1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*2)))
		b1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*2)))
		temp1 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))) * sinpi8sqrt2) >> 16
		temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*3))) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*3))) * cospi8sqrt2minus1) >> 16)
		c1 = temp1 - temp2
		temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))) * cospi8sqrt2minus1) >> 16)
		temp2 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*3))) * sinpi8sqrt2) >> 16
		d1 = temp1 + temp2
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*0)) = int16((a1 + d1 + 4) >> 3)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*3)) = int16((a1 - d1 + 4) >> 3)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*1)) = int16((b1 + c1 + 4) >> 3)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*2)) = int16((b1 - c1 + 4) >> 3)
		ip = (*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*uintptr(shortpitch)))
		op = (*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*uintptr(shortpitch)))
	}
	ip = &output[0]
	for r = 0; r < 4; r++ {
		for c = 0; c < 4; c++ {
			var a int = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*uintptr(c)))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), c)))
			if a < 0 {
				a = 0
			}
			if a > math.MaxUint8 {
				a = math.MaxUint8
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr), c)) = uint8(int8(a))
		}
		ip = (*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*4))
		dst_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr), dst_stride))
		pred_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), pred_stride))
	}
}
func Vp8DcOnlyIdctAddC(input_dc int16, pred_ptr *uint8, pred_stride int, dst_ptr *uint8, dst_stride int) {
	var (
		a1 int = ((int(input_dc) + 4) >> 3)
		r  int
		c  int
	)
	for r = 0; r < 4; r++ {
		for c = 0; c < 4; c++ {
			var a int = a1 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), c)))
			if a < 0 {
				a = 0
			}
			if a > math.MaxUint8 {
				a = math.MaxUint8
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr), c)) = uint8(int8(a))
		}
		dst_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr), dst_stride))
		pred_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), pred_stride))
	}
}
func Vp8ShortInvWalsh4x4C(input *int16, mb_dqcoeff *int16) {
	var (
		output [16]int16
		i      int
		a1     int
		b1     int
		c1     int
		d1     int
		a2     int
		b2     int
		c2     int
		d2     int
		ip     *int16 = input
		op     *int16 = &output[0]
	)
	for i = 0; i < 4; i++ {
		a1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*12)))
		b1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*4))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*8)))
		c1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*4))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*8)))
		d1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*12)))
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*0)) = int16(a1 + b1)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*4)) = int16(c1 + d1)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*8)) = int16(a1 - b1)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*12)) = int16(d1 - c1)
		ip = (*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))
		op = (*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*1))
	}
	ip = &output[0]
	op = &output[0]
	for i = 0; i < 4; i++ {
		a1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*3)))
		b1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*2)))
		c1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*2)))
		d1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*3)))
		a2 = a1 + b1
		b2 = c1 + d1
		c2 = a1 - b1
		d2 = d1 - c1
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*0)) = int16((a2 + 3) >> 3)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*1)) = int16((b2 + 3) >> 3)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*2)) = int16((c2 + 3) >> 3)
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*3)) = int16((d2 + 3) >> 3)
		ip = (*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*4))
		op = (*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*4))
	}
	for i = 0; i < 16; i++ {
		*(*int16)(unsafe.Add(unsafe.Pointer(mb_dqcoeff), unsafe.Sizeof(int16(0))*uintptr(i*16))) = output[i]
	}
}
func vp8_short_inv_walsh4x4_1_c(input *int16, mb_dqcoeff *int16) {
	var (
		i  int
		a1 int
	)
	a1 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))) + 3) >> 3
	for i = 0; i < 16; i++ {
		*(*int16)(unsafe.Add(unsafe.Pointer(mb_dqcoeff), unsafe.Sizeof(int16(0))*uintptr(i*16))) = int16(a1)
	}
}
