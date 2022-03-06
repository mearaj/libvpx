package vp8

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

type uc uint8

func vp8_signed_char_clamp(t int) int8 {
	if t < math.MinInt8 {
		t = math.MinInt8
	} else {
		t = t
	}
	if t > math.MaxInt8 {
		t = math.MaxInt8
	} else {
		t = t
	}
	return int8(t)
}
func vp8_filter_mask(limit uc, blimit uc, p3 uc, p2 uc, p1 uc, p0 uc, q0 uc, q1 uc, q2 uc, q3 uc) int8 {
	var mask int8 = 0
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(p3-p2)) > int64(limit)))
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(p2-p1)) > int64(limit)))
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(p1-p0)) > int64(limit)))
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(q1-q0)) > int64(limit)))
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(q2-q1)) > int64(limit)))
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(q3-q2)) > int64(limit)))
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(p0-q0))*2+cmath.Abs(int64(p1-q1))/2 > int64(blimit)))
	return int8(int(mask) - 1)
}
func vp8_hevmask(thresh uc, p1 uc, p0 uc, q0 uc, q1 uc) int8 {
	var hev int8 = 0
	hev |= int8(libc.BoolToInt(cmath.Abs(int64(p1-p0)) > int64(thresh))) * (-1)
	hev |= int8(libc.BoolToInt(cmath.Abs(int64(q1-q0)) > int64(thresh))) * (-1)
	return hev
}
func vp8_filter(mask int8, hev uc, op1 *uc, op0 *uc, oq0 *uc, oq1 *uc) {
	var (
		ps0          int8
		qs0          int8
		ps1          int8
		qs1          int8
		filter_value int8
		Filter1      int8
		Filter2      int8
		u            int8
	)
	ps1 = int8(int(int8(*op1)) ^ 128)
	ps0 = int8(int(int8(*op0)) ^ 128)
	qs0 = int8(int(int8(*oq0)) ^ 128)
	qs1 = int8(int(int8(*oq1)) ^ 128)
	filter_value = vp8_signed_char_clamp(int(ps1) - int(qs1))
	filter_value &= int8(hev)
	filter_value = vp8_signed_char_clamp(int(filter_value) + (int(qs0)-int(ps0))*3)
	filter_value &= mask
	Filter1 = vp8_signed_char_clamp(int(filter_value) + 4)
	Filter2 = vp8_signed_char_clamp(int(filter_value) + 3)
	Filter1 >>= 3
	Filter2 >>= 3
	u = vp8_signed_char_clamp(int(qs0) - int(Filter1))
	*oq0 = uc(int8(int(u) ^ 128))
	u = vp8_signed_char_clamp(int(ps0) + int(Filter2))
	*op0 = uc(int8(int(u) ^ 128))
	filter_value = Filter1
	filter_value += 1
	filter_value >>= 1
	filter_value &= int8(int(^hev))
	u = vp8_signed_char_clamp(int(qs1) - int(filter_value))
	*oq1 = uc(int8(int(u) ^ 128))
	u = vp8_signed_char_clamp(int(ps1) + int(filter_value))
	*op1 = uc(int8(int(u) ^ 128))
}
func loop_filter_horizontal_edge_c(s *uint8, p int, blimit *uint8, limit *uint8, thresh *uint8, count int) {
	var (
		hev  int  = 0
		mask int8 = 0
		i    int  = 0
	)
	for {
		mask = vp8_filter_mask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(limit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(blimit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-4)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-3)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-2)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-1)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*1))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*2))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*3))))
		hev = int(vp8_hevmask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(thresh), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-2)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-1)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*1)))))
		vp8_filter(mask, uc(int8(hev)), (*uc)(unsafe.Add(unsafe.Pointer(s), -(p*2))), (*uc)(unsafe.Add(unsafe.Pointer(s), -(p*1))), (*uc)(unsafe.Pointer(s)), (*uc)(unsafe.Add(unsafe.Pointer(s), p*1)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), 1))
		if func() int {
			p := &i
			*p++
			return *p
		}() >= count*8 {
			break
		}
	}
}
func loop_filter_vertical_edge_c(s *uint8, p int, blimit *uint8, limit *uint8, thresh *uint8, count int) {
	var (
		hev  int  = 0
		mask int8 = 0
		i    int  = 0
	)
	for {
		mask = vp8_filter_mask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(limit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(blimit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -4))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -3))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -2))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -1))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 1))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 2))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 3))))
		hev = int(vp8_hevmask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(thresh), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -2))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -1))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 1)))))
		vp8_filter(mask, uc(int8(hev)), (*uc)(unsafe.Add(unsafe.Pointer(s), -2)), (*uc)(unsafe.Add(unsafe.Pointer(s), -1)), (*uc)(unsafe.Pointer(s)), (*uc)(unsafe.Add(unsafe.Pointer(s), 1)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), p))
		if func() int {
			p := &i
			*p++
			return *p
		}() >= count*8 {
			break
		}
	}
}
func vp8_mbfilter(mask int8, hev uc, op2 *uc, op1 *uc, op0 *uc, oq0 *uc, oq1 *uc, oq2 *uc) {
	var (
		s            int8
		u            int8
		filter_value int8
		Filter1      int8
		Filter2      int8
		ps2          int8 = int8(int(int8(*op2)) ^ 128)
		ps1          int8 = int8(int(int8(*op1)) ^ 128)
		ps0          int8 = int8(int(int8(*op0)) ^ 128)
		qs0          int8 = int8(int(int8(*oq0)) ^ 128)
		qs1          int8 = int8(int(int8(*oq1)) ^ 128)
		qs2          int8 = int8(int(int8(*oq2)) ^ 128)
	)
	filter_value = vp8_signed_char_clamp(int(ps1) - int(qs1))
	filter_value = vp8_signed_char_clamp(int(filter_value) + (int(qs0)-int(ps0))*3)
	filter_value &= mask
	Filter2 = filter_value
	Filter2 &= int8(hev)
	Filter1 = vp8_signed_char_clamp(int(Filter2) + 4)
	Filter2 = vp8_signed_char_clamp(int(Filter2) + 3)
	Filter1 >>= 3
	Filter2 >>= 3
	qs0 = vp8_signed_char_clamp(int(qs0) - int(Filter1))
	ps0 = vp8_signed_char_clamp(int(ps0) + int(Filter2))
	filter_value &= int8(int(^hev))
	Filter2 = filter_value
	u = vp8_signed_char_clamp((int(Filter2)*27 + 63) >> 7)
	s = vp8_signed_char_clamp(int(qs0) - int(u))
	*oq0 = uc(int8(int(s) ^ 128))
	s = vp8_signed_char_clamp(int(ps0) + int(u))
	*op0 = uc(int8(int(s) ^ 128))
	u = vp8_signed_char_clamp((int(Filter2)*18 + 63) >> 7)
	s = vp8_signed_char_clamp(int(qs1) - int(u))
	*oq1 = uc(int8(int(s) ^ 128))
	s = vp8_signed_char_clamp(int(ps1) + int(u))
	*op1 = uc(int8(int(s) ^ 128))
	u = vp8_signed_char_clamp((int(Filter2)*9 + 63) >> 7)
	s = vp8_signed_char_clamp(int(qs2) - int(u))
	*oq2 = uc(int8(int(s) ^ 128))
	s = vp8_signed_char_clamp(int(ps2) + int(u))
	*op2 = uc(int8(int(s) ^ 128))
}
func mbloop_filter_horizontal_edge_c(s *uint8, p int, blimit *uint8, limit *uint8, thresh *uint8, count int) {
	var (
		hev  int8 = 0
		mask int8 = 0
		i    int  = 0
	)
	for {
		mask = vp8_filter_mask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(limit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(blimit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-4)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-3)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-2)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-1)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*1))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*2))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*3))))
		hev = vp8_hevmask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(thresh), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-2)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*(-1)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), p*1))))
		vp8_mbfilter(mask, uc(hev), (*uc)(unsafe.Add(unsafe.Pointer(s), -(p*3))), (*uc)(unsafe.Add(unsafe.Pointer(s), -(p*2))), (*uc)(unsafe.Add(unsafe.Pointer(s), -(p*1))), (*uc)(unsafe.Pointer(s)), (*uc)(unsafe.Add(unsafe.Pointer(s), p*1)), (*uc)(unsafe.Add(unsafe.Pointer(s), p*2)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), 1))
		if func() int {
			p := &i
			*p++
			return *p
		}() >= count*8 {
			break
		}
	}
}
func mbloop_filter_vertical_edge_c(s *uint8, p int, blimit *uint8, limit *uint8, thresh *uint8, count int) {
	var (
		hev  int8 = 0
		mask int8 = 0
		i    int  = 0
	)
	for {
		mask = vp8_filter_mask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(limit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(blimit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -4))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -3))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -2))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -1))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 1))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 2))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 3))))
		hev = vp8_hevmask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(thresh), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -2))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), -1))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(s), 1))))
		vp8_mbfilter(mask, uc(hev), (*uc)(unsafe.Add(unsafe.Pointer(s), -3)), (*uc)(unsafe.Add(unsafe.Pointer(s), -2)), (*uc)(unsafe.Add(unsafe.Pointer(s), -1)), (*uc)(unsafe.Pointer(s)), (*uc)(unsafe.Add(unsafe.Pointer(s), 1)), (*uc)(unsafe.Add(unsafe.Pointer(s), 2)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), p))
		if func() int {
			p := &i
			*p++
			return *p
		}() >= count*8 {
			break
		}
	}
}
func vp8_simple_filter_mask(blimit uc, p1 uc, p0 uc, q0 uc, q1 uc) int8 {
	var mask int8 = int8(libc.BoolToInt(cmath.Abs(int64(p0-q0))*2+cmath.Abs(int64(p1-q1))/2 <= int64(blimit))) * (-1)
	return mask
}
func vp8_simple_filter(mask int8, op1 *uc, op0 *uc, oq0 *uc, oq1 *uc) {
	var (
		filter_value int8
		Filter1      int8
		Filter2      int8
		p1           int8 = int8(int(int8(*op1)) ^ 128)
		p0           int8 = int8(int(int8(*op0)) ^ 128)
		q0           int8 = int8(int(int8(*oq0)) ^ 128)
		q1           int8 = int8(int(int8(*oq1)) ^ 128)
		u            int8
	)
	filter_value = vp8_signed_char_clamp(int(p1) - int(q1))
	filter_value = vp8_signed_char_clamp(int(filter_value) + (int(q0)-int(p0))*3)
	filter_value &= mask
	Filter1 = vp8_signed_char_clamp(int(filter_value) + 4)
	Filter1 >>= 3
	u = vp8_signed_char_clamp(int(q0) - int(Filter1))
	*oq0 = uc(int8(int(u) ^ 128))
	Filter2 = vp8_signed_char_clamp(int(filter_value) + 3)
	Filter2 >>= 3
	u = vp8_signed_char_clamp(int(p0) + int(Filter2))
	*op0 = uc(int8(int(u) ^ 128))
}
func Vp8LoopFilterSimpleHorizontalEdgeC(y_ptr *uint8, y_stride int, blimit *uint8) {
	var (
		mask int8 = 0
		i    int  = 0
	)
	for {
		mask = vp8_simple_filter_mask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(blimit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*(-2)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*(-1)))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*1))))
		vp8_simple_filter(mask, (*uc)(unsafe.Add(unsafe.Pointer(y_ptr), -(y_stride*2))), (*uc)(unsafe.Add(unsafe.Pointer(y_ptr), -(y_stride*1))), (*uc)(unsafe.Pointer(y_ptr)), (*uc)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*1)))
		y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 1))
		if func() int {
			p := &i
			*p++
			return *p
		}() >= 16 {
			break
		}
	}
}
func Vp8LoopFilterSimpleVerticalEdgeC(y_ptr *uint8, y_stride int, blimit *uint8) {
	var (
		mask int8 = 0
		i    int  = 0
	)
	for {
		mask = vp8_simple_filter_mask(uc(*(*uint8)(unsafe.Add(unsafe.Pointer(blimit), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), -2))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), -1))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 0))), uc(*(*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 1))))
		vp8_simple_filter(mask, (*uc)(unsafe.Add(unsafe.Pointer(y_ptr), -2)), (*uc)(unsafe.Add(unsafe.Pointer(y_ptr), -1)), (*uc)(unsafe.Pointer(y_ptr)), (*uc)(unsafe.Add(unsafe.Pointer(y_ptr), 1)))
		y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride))
		if func() int {
			p := &i
			*p++
			return *p
		}() >= 16 {
			break
		}
	}
}
func vp8_loop_filter_mbh_c(y_ptr *uint8, u_ptr *uint8, v_ptr *uint8, y_stride int, uv_stride int, lfi *loop_filter_info) {
	mbloop_filter_horizontal_edge_c(y_ptr, y_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr, 2)
	if u_ptr != nil {
		mbloop_filter_horizontal_edge_c(u_ptr, uv_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr, 1)
	}
	if v_ptr != nil {
		mbloop_filter_horizontal_edge_c(v_ptr, uv_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr, 1)
	}
}
func vp8_loop_filter_mbv_c(y_ptr *uint8, u_ptr *uint8, v_ptr *uint8, y_stride int, uv_stride int, lfi *loop_filter_info) {
	mbloop_filter_vertical_edge_c(y_ptr, y_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr, 2)
	if u_ptr != nil {
		mbloop_filter_vertical_edge_c(u_ptr, uv_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr, 1)
	}
	if v_ptr != nil {
		mbloop_filter_vertical_edge_c(v_ptr, uv_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr, 1)
	}
}
func vp8_loop_filter_bh_c(y_ptr *uint8, u_ptr *uint8, v_ptr *uint8, y_stride int, uv_stride int, lfi *loop_filter_info) {
	loop_filter_horizontal_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*4)), y_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 2)
	loop_filter_horizontal_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*8)), y_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 2)
	loop_filter_horizontal_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*12)), y_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 2)
	if u_ptr != nil {
		loop_filter_horizontal_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), uv_stride*4)), uv_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 1)
	}
	if v_ptr != nil {
		loop_filter_horizontal_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), uv_stride*4)), uv_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 1)
	}
}
func vp8_loop_filter_bhs_c(y_ptr *uint8, y_stride int, blimit *uint8) {
	Vp8LoopFilterSimpleHorizontalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*4)), y_stride, blimit)
	Vp8LoopFilterSimpleHorizontalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*8)), y_stride, blimit)
	Vp8LoopFilterSimpleHorizontalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*12)), y_stride, blimit)
}
func vp8_loop_filter_bv_c(y_ptr *uint8, u_ptr *uint8, v_ptr *uint8, y_stride int, uv_stride int, lfi *loop_filter_info) {
	loop_filter_vertical_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 4)), y_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 2)
	loop_filter_vertical_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 8)), y_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 2)
	loop_filter_vertical_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 12)), y_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 2)
	if u_ptr != nil {
		loop_filter_vertical_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), 4)), uv_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 1)
	}
	if v_ptr != nil {
		loop_filter_vertical_edge_c((*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), 4)), uv_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 1)
	}
}
func vp8_loop_filter_bvs_c(y_ptr *uint8, y_stride int, blimit *uint8) {
	Vp8LoopFilterSimpleVerticalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 4)), y_stride, blimit)
	Vp8LoopFilterSimpleVerticalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 8)), y_stride, blimit)
	Vp8LoopFilterSimpleVerticalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 12)), y_stride, blimit)
}
