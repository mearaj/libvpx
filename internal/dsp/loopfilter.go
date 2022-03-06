package dsp

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func signed_char_clamp(t int) int8 {
	return int8(clamp(t, math.MinInt8, math.MaxInt8))
}
func filter_mask(limit uint8, blimit uint8, p3 uint8, p2 uint8, p1 uint8, p0 uint8, q0 uint8, q1 uint8, q2 uint8, q3 uint8) int8 {
	var mask int8 = 0
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(p3)-int(p2))) > int64(limit))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(p2)-int(p1))) > int64(limit))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(p1)-int(p0))) > int64(limit))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(q1)-int(q0))) > int64(limit))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(q2)-int(q1))) > int64(limit))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(q3)-int(q2))) > int64(limit))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(p0)-int(q0)))*2+cmath.Abs(int64(int(p1)-int(q1)))/2 > int64(blimit))) * (-1)
	return ^mask
}
func flat_mask4(thresh uint8, p3 uint8, p2 uint8, p1 uint8, p0 uint8, q0 uint8, q1 uint8, q2 uint8, q3 uint8) int8 {
	var mask int8 = 0
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(p1)-int(p0))) > int64(thresh))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(q1)-int(q0))) > int64(thresh))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(p2)-int(p0))) > int64(thresh))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(q2)-int(q0))) > int64(thresh))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(p3)-int(p0))) > int64(thresh))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(q3)-int(q0))) > int64(thresh))) * (-1)
	return ^mask
}
func flat_mask5(thresh uint8, p4 uint8, p3 uint8, p2 uint8, p1 uint8, p0 uint8, q0 uint8, q1 uint8, q2 uint8, q3 uint8, q4 uint8) int8 {
	var mask int8 = ^flat_mask4(thresh, p3, p2, p1, p0, q0, q1, q2, q3)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(p4)-int(p0))) > int64(thresh))) * (-1)
	mask |= int8(libc.BoolToInt(cmath.Abs(int64(int(q4)-int(q0))) > int64(thresh))) * (-1)
	return ^mask
}
func hev_mask(thresh uint8, p1 uint8, p0 uint8, q0 uint8, q1 uint8) int8 {
	var hev int8 = 0
	hev |= int8(libc.BoolToInt(cmath.Abs(int64(int(p1)-int(p0))) > int64(thresh))) * (-1)
	hev |= int8(libc.BoolToInt(cmath.Abs(int64(int(q1)-int(q0))) > int64(thresh))) * (-1)
	return hev
}
func filter4(mask int8, thresh uint8, op1 *uint8, op0 *uint8, oq0 *uint8, oq1 *uint8) {
	var (
		filter1 int8
		filter2 int8
		ps1     int8 = int8(int(*op1) ^ 128)
		ps0     int8 = int8(int(*op0) ^ 128)
		qs0     int8 = int8(int(*oq0) ^ 128)
		qs1     int8 = int8(int(*oq1) ^ 128)
		hev     int8 = hev_mask(thresh, *op1, *op0, *oq0, *oq1)
		filter  int8 = int8(int(signed_char_clamp(int(ps1)-int(qs1))) & int(hev))
	)
	filter = int8(int(signed_char_clamp(int(filter)+(int(qs0)-int(ps0))*3)) & int(mask))
	filter1 = int8(int(signed_char_clamp(int(filter)+4)) >> 3)
	filter2 = int8(int(signed_char_clamp(int(filter)+3)) >> 3)
	*oq0 = uint8(int8(int(signed_char_clamp(int(qs0)-int(filter1))) ^ 128))
	*op0 = uint8(int8(int(signed_char_clamp(int(ps0)+int(filter2))) ^ 128))
	filter = int8(((int(filter1) + (1 << (1 - 1))) >> 1) & int(^hev))
	*oq1 = uint8(int8(int(signed_char_clamp(int(qs1)-int(filter))) ^ 128))
	*op1 = uint8(int8(int(signed_char_clamp(int(ps1)+int(filter))) ^ 128))
}
func vpx_lpf_horizontal_4_c(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8) {
	var i int
	for i = 0; i < 8; i++ {
		var (
			p3   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-4)))
			p2   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-3)))
			p1   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-2)))
			p0   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -pitch))
			q0   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*0))
			q1   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*1))
			q2   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*2))
			q3   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*3))
			mask int8  = filter_mask(*limit, *blimit, p3, p2, p1, p0, q0, q1, q2, q3)
		)
		filter4(mask, *thresh, (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*2))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*1))), s, (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*1)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), 1))
	}
}
func vpx_lpf_horizontal_4_dual_c(s *uint8, pitch int, blimit0 *uint8, limit0 *uint8, thresh0 *uint8, blimit1 *uint8, limit1 *uint8, thresh1 *uint8) {
	vpx_lpf_horizontal_4_c(s, pitch, blimit0, limit0, thresh0)
	vpx_lpf_horizontal_4_c((*uint8)(unsafe.Add(unsafe.Pointer(s), 8)), pitch, blimit1, limit1, thresh1)
}
func vpx_lpf_vertical_4_c(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8) {
	var i int
	for i = 0; i < 8; i++ {
		var (
			p3   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -4))
			p2   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -3))
			p1   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -2))
			p0   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -1))
			q0   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))
			q1   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 1))
			q2   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 2))
			q3   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 3))
			mask int8  = filter_mask(*limit, *blimit, p3, p2, p1, p0, q0, q1, q2, q3)
		)
		filter4(mask, *thresh, (*uint8)(unsafe.Add(unsafe.Pointer(s), -2)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -1)), s, (*uint8)(unsafe.Add(unsafe.Pointer(s), 1)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch))
	}
}
func vpx_lpf_vertical_4_dual_c(s *uint8, pitch int, blimit0 *uint8, limit0 *uint8, thresh0 *uint8, blimit1 *uint8, limit1 *uint8, thresh1 *uint8) {
	vpx_lpf_vertical_4_c(s, pitch, blimit0, limit0, thresh0)
	vpx_lpf_vertical_4_c((*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*8)), pitch, blimit1, limit1, thresh1)
}
func filter8(mask int8, thresh uint8, flat uint8, op3 *uint8, op2 *uint8, op1 *uint8, op0 *uint8, oq0 *uint8, oq1 *uint8, oq2 *uint8, oq3 *uint8) {
	if int(flat) != 0 && int(mask) != 0 {
		var (
			p3 uint8 = *op3
			p2 uint8 = *op2
			p1 uint8 = *op1
			p0 uint8 = *op0
			q0 uint8 = *oq0
			q1 uint8 = *oq1
			q2 uint8 = *oq2
			q3 uint8 = *oq3
		)
		*op2 = uint8(int8(((int(p3) + int(p3) + int(p3) + int(p2)*2 + int(p1) + int(p0) + int(q0)) + (1 << (3 - 1))) >> 3))
		*op1 = uint8(int8(((int(p3) + int(p3) + int(p2) + int(p1)*2 + int(p0) + int(q0) + int(q1)) + (1 << (3 - 1))) >> 3))
		*op0 = uint8(int8(((int(p3) + int(p2) + int(p1) + int(p0)*2 + int(q0) + int(q1) + int(q2)) + (1 << (3 - 1))) >> 3))
		*oq0 = uint8(int8(((int(p2) + int(p1) + int(p0) + int(q0)*2 + int(q1) + int(q2) + int(q3)) + (1 << (3 - 1))) >> 3))
		*oq1 = uint8(int8(((int(p1) + int(p0) + int(q0) + int(q1)*2 + int(q2) + int(q3) + int(q3)) + (1 << (3 - 1))) >> 3))
		*oq2 = uint8(int8(((int(p0) + int(q0) + int(q1) + int(q2)*2 + int(q3) + int(q3) + int(q3)) + (1 << (3 - 1))) >> 3))
	} else {
		filter4(mask, thresh, op1, op0, oq0, oq1)
	}
}
func vpx_lpf_horizontal_8_c(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8) {
	var i int
	for i = 0; i < 8; i++ {
		var (
			p3   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-4)))
			p2   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-3)))
			p1   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-2)))
			p0   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -pitch))
			q0   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*0))
			q1   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*1))
			q2   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*2))
			q3   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*3))
			mask int8  = filter_mask(*limit, *blimit, p3, p2, p1, p0, q0, q1, q2, q3)
			flat int8  = flat_mask4(1, p3, p2, p1, p0, q0, q1, q2, q3)
		)
		filter8(mask, *thresh, uint8(flat), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*4))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*3))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*2))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*1))), s, (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*1)), (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*2)), (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*3)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), 1))
	}
}
func vpx_lpf_horizontal_8_dual_c(s *uint8, pitch int, blimit0 *uint8, limit0 *uint8, thresh0 *uint8, blimit1 *uint8, limit1 *uint8, thresh1 *uint8) {
	vpx_lpf_horizontal_8_c(s, pitch, blimit0, limit0, thresh0)
	vpx_lpf_horizontal_8_c((*uint8)(unsafe.Add(unsafe.Pointer(s), 8)), pitch, blimit1, limit1, thresh1)
}
func vpx_lpf_vertical_8_c(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8) {
	var i int
	for i = 0; i < 8; i++ {
		var (
			p3   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -4))
			p2   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -3))
			p1   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -2))
			p0   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -1))
			q0   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))
			q1   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 1))
			q2   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 2))
			q3   uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 3))
			mask int8  = filter_mask(*limit, *blimit, p3, p2, p1, p0, q0, q1, q2, q3)
			flat int8  = flat_mask4(1, p3, p2, p1, p0, q0, q1, q2, q3)
		)
		filter8(mask, *thresh, uint8(flat), (*uint8)(unsafe.Add(unsafe.Pointer(s), -4)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -3)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -2)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -1)), s, (*uint8)(unsafe.Add(unsafe.Pointer(s), 1)), (*uint8)(unsafe.Add(unsafe.Pointer(s), 2)), (*uint8)(unsafe.Add(unsafe.Pointer(s), 3)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch))
	}
}
func vpx_lpf_vertical_8_dual_c(s *uint8, pitch int, blimit0 *uint8, limit0 *uint8, thresh0 *uint8, blimit1 *uint8, limit1 *uint8, thresh1 *uint8) {
	vpx_lpf_vertical_8_c(s, pitch, blimit0, limit0, thresh0)
	vpx_lpf_vertical_8_c((*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*8)), pitch, blimit1, limit1, thresh1)
}
func filter16(mask int8, thresh uint8, flat uint8, flat2 uint8, op7 *uint8, op6 *uint8, op5 *uint8, op4 *uint8, op3 *uint8, op2 *uint8, op1 *uint8, op0 *uint8, oq0 *uint8, oq1 *uint8, oq2 *uint8, oq3 *uint8, oq4 *uint8, oq5 *uint8, oq6 *uint8, oq7 *uint8) {
	if int(flat2) != 0 && int(flat) != 0 && int(mask) != 0 {
		var (
			p7 uint8 = *op7
			p6 uint8 = *op6
			p5 uint8 = *op5
			p4 uint8 = *op4
			p3 uint8 = *op3
			p2 uint8 = *op2
			p1 uint8 = *op1
			p0 uint8 = *op0
			q0 uint8 = *oq0
			q1 uint8 = *oq1
			q2 uint8 = *oq2
			q3 uint8 = *oq3
			q4 uint8 = *oq4
			q5 uint8 = *oq5
			q6 uint8 = *oq6
			q7 uint8 = *oq7
		)
		*op6 = uint8(int8(((int(p7)*7 + int(p6)*2 + int(p5) + int(p4) + int(p3) + int(p2) + int(p1) + int(p0) + int(q0)) + (1 << (4 - 1))) >> 4))
		*op5 = uint8(int8(((int(p7)*6 + int(p6) + int(p5)*2 + int(p4) + int(p3) + int(p2) + int(p1) + int(p0) + int(q0) + int(q1)) + (1 << (4 - 1))) >> 4))
		*op4 = uint8(int8(((int(p7)*5 + int(p6) + int(p5) + int(p4)*2 + int(p3) + int(p2) + int(p1) + int(p0) + int(q0) + int(q1) + int(q2)) + (1 << (4 - 1))) >> 4))
		*op3 = uint8(int8(((int(p7)*4 + int(p6) + int(p5) + int(p4) + int(p3)*2 + int(p2) + int(p1) + int(p0) + int(q0) + int(q1) + int(q2) + int(q3)) + (1 << (4 - 1))) >> 4))
		*op2 = uint8(int8(((int(p7)*3 + int(p6) + int(p5) + int(p4) + int(p3) + int(p2)*2 + int(p1) + int(p0) + int(q0) + int(q1) + int(q2) + int(q3) + int(q4)) + (1 << (4 - 1))) >> 4))
		*op1 = uint8(int8(((int(p7)*2 + int(p6) + int(p5) + int(p4) + int(p3) + int(p2) + int(p1)*2 + int(p0) + int(q0) + int(q1) + int(q2) + int(q3) + int(q4) + int(q5)) + (1 << (4 - 1))) >> 4))
		*op0 = uint8(int8(((int(p7) + int(p6) + int(p5) + int(p4) + int(p3) + int(p2) + int(p1) + int(p0)*2 + int(q0) + int(q1) + int(q2) + int(q3) + int(q4) + int(q5) + int(q6)) + (1 << (4 - 1))) >> 4))
		*oq0 = uint8(int8(((int(p6) + int(p5) + int(p4) + int(p3) + int(p2) + int(p1) + int(p0) + int(q0)*2 + int(q1) + int(q2) + int(q3) + int(q4) + int(q5) + int(q6) + int(q7)) + (1 << (4 - 1))) >> 4))
		*oq1 = uint8(int8(((int(p5) + int(p4) + int(p3) + int(p2) + int(p1) + int(p0) + int(q0) + int(q1)*2 + int(q2) + int(q3) + int(q4) + int(q5) + int(q6) + int(q7)*2) + (1 << (4 - 1))) >> 4))
		*oq2 = uint8(int8(((int(p4) + int(p3) + int(p2) + int(p1) + int(p0) + int(q0) + int(q1) + int(q2)*2 + int(q3) + int(q4) + int(q5) + int(q6) + int(q7)*3) + (1 << (4 - 1))) >> 4))
		*oq3 = uint8(int8(((int(p3) + int(p2) + int(p1) + int(p0) + int(q0) + int(q1) + int(q2) + int(q3)*2 + int(q4) + int(q5) + int(q6) + int(q7)*4) + (1 << (4 - 1))) >> 4))
		*oq4 = uint8(int8(((int(p2) + int(p1) + int(p0) + int(q0) + int(q1) + int(q2) + int(q3) + int(q4)*2 + int(q5) + int(q6) + int(q7)*5) + (1 << (4 - 1))) >> 4))
		*oq5 = uint8(int8(((int(p1) + int(p0) + int(q0) + int(q1) + int(q2) + int(q3) + int(q4) + int(q5)*2 + int(q6) + int(q7)*6) + (1 << (4 - 1))) >> 4))
		*oq6 = uint8(int8(((int(p0) + int(q0) + int(q1) + int(q2) + int(q3) + int(q4) + int(q5) + int(q6)*2 + int(q7)*7) + (1 << (4 - 1))) >> 4))
	} else {
		filter8(mask, thresh, flat, op3, op2, op1, op0, oq0, oq1, oq2, oq3)
	}
}
func mb_lpf_horizontal_edge_w(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8, count int) {
	var i int
	for i = 0; i < count*8; i++ {
		var (
			p3    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-4)))
			p2    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-3)))
			p1    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-2)))
			p0    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -pitch))
			q0    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*0))
			q1    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*1))
			q2    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*2))
			q3    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*3))
			mask  int8  = filter_mask(*limit, *blimit, p3, p2, p1, p0, q0, q1, q2, q3)
			flat  int8  = flat_mask4(1, p3, p2, p1, p0, q0, q1, q2, q3)
			flat2 int8  = flat_mask5(1, *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-8))), *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-7))), *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-6))), *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*(-5))), p0, q0, *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*4)), *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*5)), *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*6)), *(*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*7)))
		)
		filter16(mask, *thresh, uint8(flat), uint8(flat2), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*8))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*7))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*6))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*5))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*4))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*3))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*2))), (*uint8)(unsafe.Add(unsafe.Pointer(s), -(pitch*1))), s, (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*1)), (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*2)), (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*3)), (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*4)), (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*5)), (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*6)), (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch*7)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), 1))
	}
}
func vpx_lpf_horizontal_16_c(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8) {
	mb_lpf_horizontal_edge_w(s, pitch, blimit, limit, thresh, 1)
}
func vpx_lpf_horizontal_16_dual_c(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8) {
	mb_lpf_horizontal_edge_w(s, pitch, blimit, limit, thresh, 2)
}
func mb_lpf_vertical_edge_w(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8, count int) {
	var i int
	for i = 0; i < count; i++ {
		var (
			p3    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -4))
			p2    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -3))
			p1    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -2))
			p0    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), -1))
			q0    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 0))
			q1    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 1))
			q2    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 2))
			q3    uint8 = *(*uint8)(unsafe.Add(unsafe.Pointer(s), 3))
			mask  int8  = filter_mask(*limit, *blimit, p3, p2, p1, p0, q0, q1, q2, q3)
			flat  int8  = flat_mask4(1, p3, p2, p1, p0, q0, q1, q2, q3)
			flat2 int8  = flat_mask5(1, *(*uint8)(unsafe.Add(unsafe.Pointer(s), -8)), *(*uint8)(unsafe.Add(unsafe.Pointer(s), -7)), *(*uint8)(unsafe.Add(unsafe.Pointer(s), -6)), *(*uint8)(unsafe.Add(unsafe.Pointer(s), -5)), p0, q0, *(*uint8)(unsafe.Add(unsafe.Pointer(s), 4)), *(*uint8)(unsafe.Add(unsafe.Pointer(s), 5)), *(*uint8)(unsafe.Add(unsafe.Pointer(s), 6)), *(*uint8)(unsafe.Add(unsafe.Pointer(s), 7)))
		)
		filter16(mask, *thresh, uint8(flat), uint8(flat2), (*uint8)(unsafe.Add(unsafe.Pointer(s), -8)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -7)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -6)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -5)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -4)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -3)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -2)), (*uint8)(unsafe.Add(unsafe.Pointer(s), -1)), s, (*uint8)(unsafe.Add(unsafe.Pointer(s), 1)), (*uint8)(unsafe.Add(unsafe.Pointer(s), 2)), (*uint8)(unsafe.Add(unsafe.Pointer(s), 3)), (*uint8)(unsafe.Add(unsafe.Pointer(s), 4)), (*uint8)(unsafe.Add(unsafe.Pointer(s), 5)), (*uint8)(unsafe.Add(unsafe.Pointer(s), 6)), (*uint8)(unsafe.Add(unsafe.Pointer(s), 7)))
		s = (*uint8)(unsafe.Add(unsafe.Pointer(s), pitch))
	}
}
func vpx_lpf_vertical_16_c(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8) {
	mb_lpf_vertical_edge_w(s, pitch, blimit, limit, thresh, 8)
}
func vpx_lpf_vertical_16_dual_c(s *uint8, pitch int, blimit *uint8, limit *uint8, thresh *uint8) {
	mb_lpf_vertical_edge_w(s, pitch, blimit, limit, thresh, 16)
}
