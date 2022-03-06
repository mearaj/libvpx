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
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(p3-p2)) > int64(limit))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(p2-p1)) > int64(limit))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(p1-p0)) > int64(limit))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(q1-q0)) > int64(limit))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(q2-q1)) > int64(limit))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(q3-q2)) > int64(limit))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(p0-q0))*2+cmath.Abs(int64(p1-q1))/2 > int64(blimit))) * (-1))
	return int8(int(^mask))
}
func flat_mask4(thresh uint8, p3 uint8, p2 uint8, p1 uint8, p0 uint8, q0 uint8, q1 uint8, q2 uint8, q3 uint8) int8 {
	var mask int8 = 0
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(p1-p0)) > int64(thresh))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(q1-q0)) > int64(thresh))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(p2-p0)) > int64(thresh))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(q2-q0)) > int64(thresh))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(p3-p0)) > int64(thresh))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(q3-q0)) > int64(thresh))) * (-1))
	return int8(int(^mask))
}
func flat_mask5(thresh uint8, p4 uint8, p3 uint8, p2 uint8, p1 uint8, p0 uint8, q0 uint8, q1 uint8, q2 uint8, q3 uint8, q4 uint8) int8 {
	var mask int8 = int8(int(^flat_mask4(thresh, p3, p2, p1, p0, q0, q1, q2, q3)))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(p4-p0)) > int64(thresh))) * (-1))
	mask |= int8(int8(libc.BoolToInt(cmath.Abs(int64(q4-q0)) > int64(thresh))) * (-1))
	return int8(int(^mask))
}
func hev_mask(thresh uint8, p1 uint8, p0 uint8, q0 uint8, q1 uint8) int8 {
	var hev int8 = 0
	hev |= int8(int8(libc.BoolToInt(cmath.Abs(int64(p1-p0)) > int64(thresh))) * (-1))
	hev |= int8(int8(libc.BoolToInt(cmath.Abs(int64(q1-q0)) > int64(thresh))) * (-1))
	return hev
}
func filter4(mask int8, thresh uint8, op1 *uint8, op0 *uint8, oq0 *uint8, oq1 *uint8) {
	var (
		filter1 int8
		filter2 int8
		ps1     int8 = int8(*op1 ^ 128)
		ps0     int8 = int8(*op0 ^ 128)
		qs0     int8 = int8(*oq0 ^ 128)
		qs1     int8 = int8(*oq1 ^ 128)
		hev     int8 = hev_mask(thresh, *op1, *op0, *oq0, *oq1)
		filter  int8 = signed_char_clamp(int(ps1-qs1)) & hev
	)
	filter = signed_char_clamp(int(filter+(qs0-ps0)*3)) & mask
	filter1 = signed_char_clamp(int(filter+4)) >> 3
	filter2 = signed_char_clamp(int(filter+3)) >> 3
	*oq0 = uint8(signed_char_clamp(int(qs0-filter1)) ^ (math.MinInt8))
	*op0 = uint8(signed_char_clamp(int(ps0+filter2)) ^ (math.MinInt8))
	filter = int8(int((filter1+(1<<(1-1)))>>1) & int(^hev))
	*oq1 = uint8(signed_char_clamp(int(qs1-filter)) ^ (math.MinInt8))
	*op1 = uint8(signed_char_clamp(int(ps1+filter)) ^ (math.MinInt8))
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
	if flat != 0 && mask != 0 {
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
		*op2 = ((p3 + p3 + p3 + p2*2 + p1 + p0 + q0) + (1 << (3 - 1))) >> 3
		*op1 = ((p3 + p3 + p2 + p1*2 + p0 + q0 + q1) + (1 << (3 - 1))) >> 3
		*op0 = ((p3 + p2 + p1 + p0*2 + q0 + q1 + q2) + (1 << (3 - 1))) >> 3
		*oq0 = ((p2 + p1 + p0 + q0*2 + q1 + q2 + q3) + (1 << (3 - 1))) >> 3
		*oq1 = ((p1 + p0 + q0 + q1*2 + q2 + q3 + q3) + (1 << (3 - 1))) >> 3
		*oq2 = ((p0 + q0 + q1 + q2*2 + q3 + q3 + q3) + (1 << (3 - 1))) >> 3
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
	if flat2 != 0 && flat != 0 && mask != 0 {
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
		*op6 = ((p7*7 + p6*2 + p5 + p4 + p3 + p2 + p1 + p0 + q0) + (1 << (4 - 1))) >> 4
		*op5 = ((p7*6 + p6 + p5*2 + p4 + p3 + p2 + p1 + p0 + q0 + q1) + (1 << (4 - 1))) >> 4
		*op4 = ((p7*5 + p6 + p5 + p4*2 + p3 + p2 + p1 + p0 + q0 + q1 + q2) + (1 << (4 - 1))) >> 4
		*op3 = ((p7*4 + p6 + p5 + p4 + p3*2 + p2 + p1 + p0 + q0 + q1 + q2 + q3) + (1 << (4 - 1))) >> 4
		*op2 = ((p7*3 + p6 + p5 + p4 + p3 + p2*2 + p1 + p0 + q0 + q1 + q2 + q3 + q4) + (1 << (4 - 1))) >> 4
		*op1 = ((p7*2 + p6 + p5 + p4 + p3 + p2 + p1*2 + p0 + q0 + q1 + q2 + q3 + q4 + q5) + (1 << (4 - 1))) >> 4
		*op0 = ((p7 + p6 + p5 + p4 + p3 + p2 + p1 + p0*2 + q0 + q1 + q2 + q3 + q4 + q5 + q6) + (1 << (4 - 1))) >> 4
		*oq0 = ((p6 + p5 + p4 + p3 + p2 + p1 + p0 + q0*2 + q1 + q2 + q3 + q4 + q5 + q6 + q7) + (1 << (4 - 1))) >> 4
		*oq1 = ((p5 + p4 + p3 + p2 + p1 + p0 + q0 + q1*2 + q2 + q3 + q4 + q5 + q6 + q7*2) + (1 << (4 - 1))) >> 4
		*oq2 = ((p4 + p3 + p2 + p1 + p0 + q0 + q1 + q2*2 + q3 + q4 + q5 + q6 + q7*3) + (1 << (4 - 1))) >> 4
		*oq3 = ((p3 + p2 + p1 + p0 + q0 + q1 + q2 + q3*2 + q4 + q5 + q6 + q7*4) + (1 << (4 - 1))) >> 4
		*oq4 = ((p2 + p1 + p0 + q0 + q1 + q2 + q3 + q4*2 + q5 + q6 + q7*5) + (1 << (4 - 1))) >> 4
		*oq5 = ((p1 + p0 + q0 + q1 + q2 + q3 + q4 + q5*2 + q6 + q7*6) + (1 << (4 - 1))) >> 4
		*oq6 = ((p0 + q0 + q1 + q2 + q3 + q4 + q5 + q6*2 + q7*7) + (1 << (4 - 1))) >> 4
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
