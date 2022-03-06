package dsp

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"math"
	"unsafe"
)

func vpx_avg_8x8_c(s *uint8, p int) uint {
	var (
		i   int
		j   int
		sum int = 0
	)
	for i = 0; i < 8; func() *uint8 {
		i++
		return func() *uint8 {
			s += (*uint8)(unsafe.Pointer(uintptr(p)))
			return s
		}()
	}() {
		for j = 0; j < 8; func() int {
			sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), j)))
			return func() int {
				p := &j
				*p++
				return *p
			}()
		}() {
		}
	}
	return uint((sum + 32) >> 6)
}
func vpx_avg_4x4_c(s *uint8, p int) uint {
	var (
		i   int
		j   int
		sum int = 0
	)
	for i = 0; i < 4; func() *uint8 {
		i++
		return func() *uint8 {
			s += (*uint8)(unsafe.Pointer(uintptr(p)))
			return s
		}()
	}() {
		for j = 0; j < 4; func() int {
			sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), j)))
			return func() int {
				p := &j
				*p++
				return *p
			}()
		}() {
		}
	}
	return uint((sum + 8) >> 4)
}
func hadamard_col8(src_diff *int16, src_stride int64, coeff *int16) {
	var (
		b0 int16 = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*0)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*1)))))
		b1 int16 = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*0)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*1)))))
		b2 int16 = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*2)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*3)))))
		b3 int16 = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*2)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*3)))))
		b4 int16 = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*4)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*5)))))
		b5 int16 = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*4)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*5)))))
		b6 int16 = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*6)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*7)))))
		b7 int16 = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*6)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr(src_stride*7)))))
		c0 int16 = int16(int(b0) + int(b2))
		c1 int16 = int16(int(b1) + int(b3))
		c2 int16 = int16(int(b0) - int(b2))
		c3 int16 = int16(int(b1) - int(b3))
		c4 int16 = int16(int(b4) + int(b6))
		c5 int16 = int16(int(b5) + int(b7))
		c6 int16 = int16(int(b4) - int(b6))
		c7 int16 = int16(int(b5) - int(b7))
	)
	*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*0)) = int16(int(c0) + int(c4))
	*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*7)) = int16(int(c1) + int(c5))
	*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*3)) = int16(int(c2) + int(c6))
	*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*4)) = int16(int(c3) + int(c7))
	*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*2)) = int16(int(c0) - int(c4))
	*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*6)) = int16(int(c1) - int(c5))
	*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*1)) = int16(int(c2) - int(c6))
	*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*5)) = int16(int(c3) - int(c7))
}
func vpx_hadamard_8x8_c(src_diff *int16, src_stride int64, coeff *int16) {
	var (
		idx     int
		buffer  [64]int16
		buffer2 [64]int16
		tmp_buf *int16 = &buffer[0]
	)
	for idx = 0; idx < 8; idx++ {
		hadamard_col8(src_diff, src_stride, tmp_buf)
		tmp_buf = (*int16)(unsafe.Add(unsafe.Pointer(tmp_buf), unsafe.Sizeof(int16(0))*8))
		src_diff = (*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*1))
	}
	tmp_buf = &buffer[0]
	for idx = 0; idx < 8; idx++ {
		hadamard_col8(tmp_buf, 8, &buffer2[idx*8])
		tmp_buf = (*int16)(unsafe.Add(unsafe.Pointer(tmp_buf), unsafe.Sizeof(int16(0))*1))
	}
	for idx = 0; idx < 64; idx++ {
		*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*uintptr(idx))) = buffer2[idx]
	}
}
func vpx_hadamard_16x16_c(src_diff *int16, src_stride int64, coeff *int16) {
	var idx int
	for idx = 0; idx < 4; idx++ {
		var src_ptr *int16 = (*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr((idx>>1)*8*int(src_stride))))), unsafe.Sizeof(int16(0))*uintptr((idx&1)*8)))
		vpx_hadamard_8x8_c(src_ptr, src_stride, (*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*uintptr(idx*64))))
	}
	for idx = 0; idx < 64; idx++ {
		var (
			a0 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*0))
			a1 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*64))
			a2 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*128))
			a3 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*192))
			b0 int16 = int16((int(a0) + int(a1)) >> 1)
			b1 int16 = int16((int(a0) - int(a1)) >> 1)
			b2 int16 = int16((int(a2) + int(a3)) >> 1)
			b3 int16 = int16((int(a2) - int(a3)) >> 1)
		)
		*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*0)) = int16(int(b0) + int(b2))
		*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*64)) = int16(int(b1) + int(b3))
		*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*128)) = int16(int(b0) - int(b2))
		*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*192)) = int16(int(b1) - int(b3))
		coeff = (*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*1))
	}
}
func vpx_hadamard_32x32_c(src_diff *int16, src_stride int64, coeff *int16) {
	var idx int
	for idx = 0; idx < 4; idx++ {
		var src_ptr *int16 = (*int16)(unsafe.Add(unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(src_diff), unsafe.Sizeof(int16(0))*uintptr((idx>>1)*16*int(src_stride))))), unsafe.Sizeof(int16(0))*uintptr((idx&1)*16)))
		vpx_hadamard_16x16_c(src_ptr, src_stride, (*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*uintptr(idx*256))))
	}
	for idx = 0; idx < 256; idx++ {
		var (
			a0 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*0))
			a1 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*256))
			a2 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*512))
			a3 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*768))
			b0 int16 = int16((int(a0) + int(a1)) >> 2)
			b1 int16 = int16((int(a0) - int(a1)) >> 2)
			b2 int16 = int16((int(a2) + int(a3)) >> 2)
			b3 int16 = int16((int(a2) - int(a3)) >> 2)
		)
		*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*0)) = int16(int(b0) + int(b2))
		*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*256)) = int16(int(b1) + int(b3))
		*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*512)) = int16(int(b0) - int(b2))
		*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*768)) = int16(int(b1) - int(b3))
		coeff = (*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*1))
	}
}
func vpx_satd_c(coeff *int16, length int) int {
	var (
		i    int
		satd int = 0
	)
	for i = 0; i < length; i++ {
		satd += int(cmath.Abs(int64(*(*int16)(unsafe.Add(unsafe.Pointer(coeff), unsafe.Sizeof(int16(0))*uintptr(i))))))
	}
	return satd
}
func vpx_int_pro_row_c(hbuf [16]int16, ref *uint8, ref_stride int, height int) {
	var (
		idx         int
		norm_factor int = height >> 1
	)
	for idx = 0; idx < 16; idx++ {
		var i int
		hbuf[idx] = 0
		for i = 0; i < height; i++ {
			hbuf[idx] += int16(*(*uint8)(unsafe.Add(unsafe.Pointer(ref), i*ref_stride)))
		}
		hbuf[idx] /= int16(norm_factor)
		ref = (*uint8)(unsafe.Add(unsafe.Pointer(ref), 1))
	}
}
func vpx_int_pro_col_c(ref *uint8, width int) int16 {
	var (
		idx int
		sum int16 = 0
	)
	for idx = 0; idx < width; idx++ {
		sum += int16(*(*uint8)(unsafe.Add(unsafe.Pointer(ref), idx)))
	}
	return sum
}
func vpx_vector_var_c(ref *int16, src *int16, bwl int) int {
	var (
		i     int
		width int = 4 << bwl
		sse   int = 0
		mean  int = 0
		var_  int
	)
	for i = 0; i < width; i++ {
		var diff int = int(*(*int16)(unsafe.Add(unsafe.Pointer(ref), unsafe.Sizeof(int16(0))*uintptr(i)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(src), unsafe.Sizeof(int16(0))*uintptr(i))))
		mean += diff
		sse += diff * diff
	}
	var_ = sse - ((mean * mean) >> (bwl + 2))
	return var_
}
func vpx_minmax_8x8_c(s *uint8, p int, d *uint8, dp int, min *int, max *int) {
	var (
		i int
		j int
	)
	*min = math.MaxUint8
	*max = 0
	for i = 0; i < 8; func() *uint8 {
		i++
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), p))
		return func() *uint8 {
			d += (*uint8)(unsafe.Pointer(uintptr(dp)))
			return d
		}()
	}() {
		for j = 0; j < 8; j++ {
			var diff int = int(cmath.Abs(int64(int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), j))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(d), j))))))
			if diff < *min {
				*min = diff
			} else {
				*min = *min
			}
			if diff > *max {
				*max = diff
			} else {
				*max = *max
			}
		}
	}
}
