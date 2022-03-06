package dsp

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

var vpx_rv [440]int16 = [440]int16{8, 5, 2, 2, 8, 12, 4, 9, 8, 3, 0, 3, 9, 0, 0, 0, 8, 3, 14, 4, 10, 1, 11, 14, 1, 14, 9, 6, 12, 11, 8, 6, 10, 0, 0, 8, 9, 0, 3, 14, 8, 11, 13, 4, 2, 9, 0, 3, 9, 6, 1, 2, 3, 14, 13, 1, 8, 2, 9, 7, 3, 3, 1, 13, 13, 6, 6, 5, 2, 7, 11, 9, 11, 8, 7, 3, 2, 0, 13, 13, 14, 4, 12, 5, 12, 10, 8, 10, 13, 10, 4, 14, 4, 10, 0, 8, 11, 1, 13, 7, 7, 14, 6, 14, 13, 2, 13, 5, 4, 4, 0, 10, 0, 5, 13, 2, 12, 7, 11, 13, 8, 0, 4, 10, 7, 2, 7, 2, 2, 5, 3, 4, 7, 3, 3, 14, 14, 5, 9, 13, 3, 14, 3, 6, 3, 0, 11, 8, 13, 1, 13, 1, 12, 0, 10, 9, 7, 6, 2, 8, 5, 2, 13, 7, 1, 13, 14, 7, 6, 7, 9, 6, 10, 11, 7, 8, 7, 5, 14, 8, 4, 4, 0, 8, 7, 10, 0, 8, 14, 11, 3, 12, 5, 7, 14, 3, 14, 5, 2, 6, 11, 12, 12, 8, 0, 11, 13, 1, 2, 0, 5, 10, 14, 7, 8, 0, 4, 11, 0, 8, 0, 3, 10, 5, 8, 0, 11, 6, 7, 8, 10, 7, 13, 9, 2, 5, 1, 5, 10, 2, 4, 3, 5, 6, 10, 8, 9, 4, 11, 14, 0, 10, 0, 5, 13, 2, 12, 7, 11, 13, 8, 0, 4, 10, 7, 2, 7, 2, 2, 5, 3, 4, 7, 3, 3, 14, 14, 5, 9, 13, 3, 14, 3, 6, 3, 0, 11, 8, 13, 1, 13, 1, 12, 0, 10, 9, 7, 6, 2, 8, 5, 2, 13, 7, 1, 13, 14, 7, 6, 7, 9, 6, 10, 11, 7, 8, 7, 5, 14, 8, 4, 4, 0, 8, 7, 10, 0, 8, 14, 11, 3, 12, 5, 7, 14, 3, 14, 5, 2, 6, 11, 12, 12, 8, 0, 11, 13, 1, 2, 0, 5, 10, 14, 7, 8, 0, 4, 11, 0, 8, 0, 3, 10, 5, 8, 0, 11, 6, 7, 8, 10, 7, 13, 9, 2, 5, 1, 5, 10, 2, 4, 3, 5, 6, 10, 8, 9, 4, 11, 14, 3, 8, 3, 7, 8, 5, 11, 4, 12, 3, 11, 9, 14, 8, 14, 13, 4, 3, 1, 2, 14, 6, 5, 4, 4, 11, 4, 6, 2, 1, 5, 8, 8, 12, 13, 5, 14, 10, 12, 13, 0, 9, 5, 5, 11, 10, 13, 9, 10, 13}

func VpxPostProcDownAndAcrossMbRowC(src *uint8, dst *uint8, src_pitch int, dst_pitch int, cols int, flimits *uint8, size int) {
	var (
		p_src *uint8
		p_dst *uint8
		row   int
		col   int
		v     uint8
		d     [4]uint8
	)
	if size >= 8 {
	} else {
		__assert_fail(libc.CString("size >= 8"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	if cols >= 8 {
	} else {
		__assert_fail(libc.CString("cols >= 8"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	for row = 0; row < size; row++ {
		p_src = src
		p_dst = dst
		for col = 0; col < cols; col++ {
			var (
				p_above2 uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col-src_pitch*2))
				p_above1 uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col-src_pitch))
				p_below1 uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col+src_pitch))
				p_below2 uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col+src_pitch*2))
			)
			v = *(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col))
			if cmath.Abs(int64(int(v)-int(p_above2))) < int64(*(*uint8)(unsafe.Add(unsafe.Pointer(flimits), col))) && cmath.Abs(int64(int(v)-int(p_above1))) < int64(*(*uint8)(unsafe.Add(unsafe.Pointer(flimits), col))) && cmath.Abs(int64(int(v)-int(p_below1))) < int64(*(*uint8)(unsafe.Add(unsafe.Pointer(flimits), col))) && cmath.Abs(int64(int(v)-int(p_below2))) < int64(*(*uint8)(unsafe.Add(unsafe.Pointer(flimits), col))) {
				var (
					k1 uint8
					k2 uint8
					k3 uint8
				)
				k1 = uint8(int8((int(p_above2) + int(p_above1) + 1) >> 1))
				k2 = uint8(int8((int(p_below2) + int(p_below1) + 1) >> 1))
				k3 = uint8(int8((int(k1) + int(k2) + 1) >> 1))
				v = uint8(int8((int(k3) + int(v) + 1) >> 1))
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(p_dst), col)) = v
		}
		p_src = dst
		p_dst = dst
		*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), -2)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(p_src), -1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), -1)) = *(*uint8)(unsafe.Add(unsafe.Pointer(p_src), 0))
			return *p
		}()
		*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), cols)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(p_src), cols+1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), cols+1)) = *(*uint8)(unsafe.Add(unsafe.Pointer(p_src), cols-1))
			return *p
		}()
		for col = 0; col < cols; col++ {
			v = *(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col))
			if cmath.Abs(int64(int(v)-int(*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col-2))))) < int64(*(*uint8)(unsafe.Add(unsafe.Pointer(flimits), col))) && cmath.Abs(int64(int(v)-int(*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col-1))))) < int64(*(*uint8)(unsafe.Add(unsafe.Pointer(flimits), col))) && cmath.Abs(int64(int(v)-int(*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col+1))))) < int64(*(*uint8)(unsafe.Add(unsafe.Pointer(flimits), col))) && cmath.Abs(int64(int(v)-int(*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col+2))))) < int64(*(*uint8)(unsafe.Add(unsafe.Pointer(flimits), col))) {
				var (
					k1 uint8
					k2 uint8
					k3 uint8
				)
				k1 = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col-2))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col-1))) + 1) >> 1))
				k2 = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col+2))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(p_src), col+1))) + 1) >> 1))
				k3 = uint8(int8((int(k1) + int(k2) + 1) >> 1))
				v = uint8(int8((int(k3) + int(v) + 1) >> 1))
			}
			d[col&3] = v
			if col >= 2 {
				*(*uint8)(unsafe.Add(unsafe.Pointer(p_dst), col-2)) = d[(col-2)&3]
			}
		}
		*(*uint8)(unsafe.Add(unsafe.Pointer(p_dst), col-2)) = d[(col-2)&3]
		*(*uint8)(unsafe.Add(unsafe.Pointer(p_dst), col-1)) = d[(col-1)&3]
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_pitch))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_pitch))
	}
}
func VpxMbPostProcAcrossIpC(src *uint8, pitch int, rows int, cols int, flimit int) {
	var (
		r int
		c int
		i int
		s *uint8 = src
		d [16]uint8
	)
	for r = 0; r < rows; r++ {
		var (
			sumsq int = 16
			sum   int = 0
		)
		for i = -8; i < 0; i++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(s), i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))
		}
		for i = 0; i < 17; i++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(s), i+cols)) = *(*uint8)(unsafe.Add(unsafe.Pointer(s), cols-1))
		}
		for i = -8; i <= 6; i++ {
			sumsq += int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), i))) * int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), i)))
			sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), i)))
			d[i+8] = 0
		}
		for c = 0; c < cols+8; c++ {
			var (
				x int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), c+7))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), c-8)))
				y int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), c+7))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), c-8)))
			)
			sum += x
			sumsq += x * y
			d[c&15] = *(*uint8)(unsafe.Add(unsafe.Pointer(s), c))
			if sumsq*15-sum*sum < flimit {
				d[c&15] = uint8(int8((sum + 8 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), c)))) >> 4))
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(s), c-8)) = d[(c-8)&15]
		}
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch))
	}
}
func VpxMbPostProcDownC(dst *uint8, pitch int, rows int, cols int, flimit int) {
	var (
		r int
		c int
		i int
	)
	for c = 0; c < cols; c++ {
		var (
			s     *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(dst), c))
			sumsq int    = 0
			sum   int    = 0
			d     [16]uint8
		)
		for i = -8; i < 0; i++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(s), i*pitch)) = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))
		}
		for i = 0; i < 17; i++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(s), (i+rows)*pitch)) = *(*uint8)(unsafe.Add(unsafe.Pointer(s), (rows-1)*pitch))
		}
		for i = -8; i <= 6; i++ {
			sumsq += int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), i*pitch))) * int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), i*pitch)))
			sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), i*pitch)))
		}
		for r = 0; r < rows+8; r++ {
			sumsq += int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*7)))*int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*7))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-8))))*int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-8))))
			sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*7))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-8))))
			d[r&15] = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))
			if sumsq*15-sum*sum < flimit {
				d[r&15] = uint8(int8((int(vpx_rv[(r&math.MaxInt8)+(c&7)]) + sum + int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 0)))) >> 4))
			}
			if r >= 8 {
				*(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-8))) = d[(r-8)&15]
			}
			s = (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch))
		}
	}
}
