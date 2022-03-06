package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func check_range(input int) int {
	return input
}
func dct_const_round_shift(input int) int {
	var rv int = ((input + (1 << (int(DCT_CONST_BITS - 1)))) >> DCT_CONST_BITS)
	return rv
}
func clip_pixel_add(dest uint8, trans int) uint8 {
	trans = check_range(trans)
	return clip_pixel(int(dest) + trans)
}
func vpx_iwht4x4_16_add_c(input *int16, dest *uint8, stride int) {
	var (
		i      int
		output [16]int16
		a1     int
		b1     int
		c1     int
		d1     int
		e1     int
		ip     *int16 = input
		op     *int16 = &output[0]
	)
	for i = 0; i < 4; i++ {
		a1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) >> UNIT_QUANT_SHIFT
		c1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))) >> UNIT_QUANT_SHIFT
		d1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*2))) >> UNIT_QUANT_SHIFT
		b1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*3))) >> UNIT_QUANT_SHIFT
		a1 += c1
		d1 -= b1
		e1 = (a1 - d1) >> 1
		b1 = e1 - b1
		c1 = e1 - c1
		a1 -= b1
		d1 += c1
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*0)) = int16(check_range(a1))
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*1)) = int16(check_range(b1))
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*2)) = int16(check_range(c1))
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*3)) = int16(check_range(d1))
		ip = (*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*4))
		op = (*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*4))
	}
	ip = &output[0]
	for i = 0; i < 4; i++ {
		a1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*(4*0))))
		c1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*(4*1))))
		d1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*(4*2))))
		b1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*(4*3))))
		a1 += c1
		d1 -= b1
		e1 = (a1 - d1) >> 1
		b1 = e1 - b1
		c1 = e1 - c1
		a1 -= b1
		d1 += c1
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*0)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*0)), check_range(a1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*1)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*1)), check_range(b1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*2)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*2)), check_range(c1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*3)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*3)), check_range(d1))
		ip = (*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))
		dest = (*uint8)(unsafe.Add(unsafe.Pointer(dest), 1))
	}
}
func vpx_iwht4x4_1_add_c(input *int16, dest *uint8, stride int) {
	var (
		i   int
		a1  int
		e1  int
		tmp [4]int16
		ip  *int16 = input
		op  *int16 = &tmp[0]
	)
	a1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) >> UNIT_QUANT_SHIFT
	e1 = a1 >> 1
	a1 -= e1
	*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*0)) = int16(check_range(a1))
	*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*1)) = func() int16 {
		p := (*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*2))
		*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*2)) = func() int16 {
			p := (*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*3))
			*(*int16)(unsafe.Add(unsafe.Pointer(op), unsafe.Sizeof(int16(0))*3)) = int16(check_range(e1))
			return *p
		}()
		return *p
	}()
	ip = &tmp[0]
	for i = 0; i < 4; i++ {
		e1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) >> 1
		a1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*0))) - e1
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*0)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*0)), a1)
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*1)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*1)), e1)
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*2)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*2)), e1)
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*3)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), stride*3)), e1)
		ip = (*int16)(unsafe.Add(unsafe.Pointer(ip), unsafe.Sizeof(int16(0))*1))
		dest = (*uint8)(unsafe.Add(unsafe.Pointer(dest), 1))
	}
}
func iadst4_c(input *int16, output *int16) {
	var (
		s0 int
		s1 int
		s2 int
		s3 int
		s4 int
		s5 int
		s6 int
		s7 int
		x0 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))
		x1 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1))
		x2 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*2))
		x3 int16 = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*3))
	)
	if (int(x0) | int(x1) | int(x2) | int(x3)) == 0 {
		libc.MemSet(unsafe.Pointer(output), 0, int(4*unsafe.Sizeof(int16(0))))
		return
	}
	s0 = int(sinpi_1_9) * int(x0)
	s1 = int(sinpi_2_9) * int(x0)
	s2 = int(sinpi_3_9) * int(x1)
	s3 = int(sinpi_4_9) * int(x2)
	s4 = int(sinpi_1_9) * int(x2)
	s5 = int(sinpi_2_9) * int(x3)
	s6 = int(sinpi_4_9) * int(x3)
	s7 = check_range(int(x0) - int(x2) + int(x3))
	s0 = s0 + s3 + s5
	s1 = s1 - s4 - s6
	s3 = s2
	s2 = int(sinpi_3_9) * s7
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(check_range(dct_const_round_shift(s0 + s3)))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*1)) = int16(check_range(dct_const_round_shift(s1 + s3)))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*2)) = int16(check_range(dct_const_round_shift(s2)))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*3)) = int16(check_range(dct_const_round_shift(s0 + s1 - s3)))
}
func idct4_c(input *int16, output *int16) {
	var (
		step  [4]int16
		temp1 int
		temp2 int
	)
	temp1 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*2)))) * int(cospi_16_64)
	temp2 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*2)))) * int(cospi_16_64)
	step[0] = int16(check_range(dct_const_round_shift(temp1)))
	step[1] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1)))*int(cospi_24_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*3)))*int(cospi_8_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1)))*int(cospi_8_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*3)))*int(cospi_24_64)
	step[2] = int16(check_range(dct_const_round_shift(temp1)))
	step[3] = int16(check_range(dct_const_round_shift(temp2)))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(check_range(int(step[0]) + int(step[3])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*1)) = int16(check_range(int(step[1]) + int(step[2])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*2)) = int16(check_range(int(step[1]) - int(step[2])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*3)) = int16(check_range(int(step[0]) - int(step[3])))
}
func vpx_idct4x4_16_add_c(input *int16, dest *uint8, stride int) {
	var (
		i        int
		j        int
		out      [16]int16
		outptr   *int16 = &out[0]
		temp_in  [4]int16
		temp_out [4]int16
	)
	for i = 0; i < 4; i++ {
		idct4_c(input, outptr)
		input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*4))
		outptr = (*int16)(unsafe.Add(unsafe.Pointer(outptr), unsafe.Sizeof(int16(0))*4))
	}
	for i = 0; i < 4; i++ {
		for j = 0; j < 4; j++ {
			temp_in[j] = out[j*4+i]
		}
		idct4_c(&temp_in[0], &temp_out[0])
		for j = 0; j < 4; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)), (int(temp_out[j])+(1<<(4-1)))>>4)
		}
	}
}
func vpx_idct4x4_1_add_c(input *int16, dest *uint8, stride int) {
	var (
		i   int
		a1  int
		out int16 = int16(check_range(dct_const_round_shift(int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))) * int(cospi_16_64))))
	)
	out = int16(check_range(dct_const_round_shift(int(out) * int(cospi_16_64))))
	a1 = (int(out) + (1 << (4 - 1))) >> 4
	for i = 0; i < 4; i++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), 0)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), 0)), a1)
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), 1)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), 1)), a1)
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), 2)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), 2)), a1)
		*(*uint8)(unsafe.Add(unsafe.Pointer(dest), 3)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), 3)), a1)
		dest = (*uint8)(unsafe.Add(unsafe.Pointer(dest), stride))
	}
}
func iadst8_c(input *int16, output *int16) {
	var (
		s0 int
		s1 int
		s2 int
		s3 int
		s4 int
		s5 int
		s6 int
		s7 int
		x0 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*7)))
		x1 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0)))
		x2 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*5)))
		x3 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*2)))
		x4 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*3)))
		x5 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*4)))
		x6 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1)))
		x7 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*6)))
	)
	if (x0 | x1 | x2 | x3 | x4 | x5 | x6 | x7) == 0 {
		libc.MemSet(unsafe.Pointer(output), 0, int(8*unsafe.Sizeof(int16(0))))
		return
	}
	s0 = int(cospi_2_64)*x0 + int(cospi_30_64)*x1
	s1 = int(cospi_30_64)*x0 - int(cospi_2_64)*x1
	s2 = int(cospi_10_64)*x2 + int(cospi_22_64)*x3
	s3 = int(cospi_22_64)*x2 - int(cospi_10_64)*x3
	s4 = int(cospi_18_64)*x4 + int(cospi_14_64)*x5
	s5 = int(cospi_14_64)*x4 - int(cospi_18_64)*x5
	s6 = int(cospi_26_64)*x6 + int(cospi_6_64)*x7
	s7 = int(cospi_6_64)*x6 - int(cospi_26_64)*x7
	x0 = check_range(dct_const_round_shift(s0 + s4))
	x1 = check_range(dct_const_round_shift(s1 + s5))
	x2 = check_range(dct_const_round_shift(s2 + s6))
	x3 = check_range(dct_const_round_shift(s3 + s7))
	x4 = check_range(dct_const_round_shift(s0 - s4))
	x5 = check_range(dct_const_round_shift(s1 - s5))
	x6 = check_range(dct_const_round_shift(s2 - s6))
	x7 = check_range(dct_const_round_shift(s3 - s7))
	s0 = x0
	s1 = x1
	s2 = x2
	s3 = x3
	s4 = int(cospi_8_64)*x4 + int(cospi_24_64)*x5
	s5 = int(cospi_24_64)*x4 - int(cospi_8_64)*x5
	s6 = int(-cospi_24_64)*x6 + int(cospi_8_64)*x7
	s7 = int(cospi_8_64)*x6 + int(cospi_24_64)*x7
	x0 = check_range(s0 + s2)
	x1 = check_range(s1 + s3)
	x2 = check_range(s0 - s2)
	x3 = check_range(s1 - s3)
	x4 = check_range(dct_const_round_shift(s4 + s6))
	x5 = check_range(dct_const_round_shift(s5 + s7))
	x6 = check_range(dct_const_round_shift(s4 - s6))
	x7 = check_range(dct_const_round_shift(s5 - s7))
	s2 = int(cospi_16_64) * (x2 + x3)
	s3 = int(cospi_16_64) * (x2 - x3)
	s6 = int(cospi_16_64) * (x6 + x7)
	s7 = int(cospi_16_64) * (x6 - x7)
	x2 = check_range(dct_const_round_shift(s2))
	x3 = check_range(dct_const_round_shift(s3))
	x6 = check_range(dct_const_round_shift(s6))
	x7 = check_range(dct_const_round_shift(s7))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(check_range(x0))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*1)) = int16(check_range(-x4))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*2)) = int16(check_range(x6))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*3)) = int16(check_range(-x2))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*4)) = int16(check_range(x3))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*5)) = int16(check_range(-x7))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*6)) = int16(check_range(x5))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*7)) = int16(check_range(-x1))
}
func idct8_c(input *int16, output *int16) {
	var (
		step1 [8]int16
		step2 [8]int16
		temp1 int
		temp2 int
	)
	step1[0] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))
	step1[2] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*4))
	step1[1] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*2))
	step1[3] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*6))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1)))*int(cospi_28_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*7)))*int(cospi_4_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1)))*int(cospi_4_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*7)))*int(cospi_28_64)
	step1[4] = int16(check_range(dct_const_round_shift(temp1)))
	step1[7] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*5)))*int(cospi_12_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*3)))*int(cospi_20_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*5)))*int(cospi_20_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*3)))*int(cospi_12_64)
	step1[5] = int16(check_range(dct_const_round_shift(temp1)))
	step1[6] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = (int(step1[0]) + int(step1[2])) * int(cospi_16_64)
	temp2 = (int(step1[0]) - int(step1[2])) * int(cospi_16_64)
	step2[0] = int16(check_range(dct_const_round_shift(temp1)))
	step2[1] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step1[1])*int(cospi_24_64) - int(step1[3])*int(cospi_8_64)
	temp2 = int(step1[1])*int(cospi_8_64) + int(step1[3])*int(cospi_24_64)
	step2[2] = int16(check_range(dct_const_round_shift(temp1)))
	step2[3] = int16(check_range(dct_const_round_shift(temp2)))
	step2[4] = int16(check_range(int(step1[4]) + int(step1[5])))
	step2[5] = int16(check_range(int(step1[4]) - int(step1[5])))
	step2[6] = int16(check_range(int(-step1[6]) + int(step1[7])))
	step2[7] = int16(check_range(int(step1[6]) + int(step1[7])))
	step1[0] = int16(check_range(int(step2[0]) + int(step2[3])))
	step1[1] = int16(check_range(int(step2[1]) + int(step2[2])))
	step1[2] = int16(check_range(int(step2[1]) - int(step2[2])))
	step1[3] = int16(check_range(int(step2[0]) - int(step2[3])))
	step1[4] = step2[4]
	temp1 = (int(step2[6]) - int(step2[5])) * int(cospi_16_64)
	temp2 = (int(step2[5]) + int(step2[6])) * int(cospi_16_64)
	step1[5] = int16(check_range(dct_const_round_shift(temp1)))
	step1[6] = int16(check_range(dct_const_round_shift(temp2)))
	step1[7] = step2[7]
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(check_range(int(step1[0]) + int(step1[7])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*1)) = int16(check_range(int(step1[1]) + int(step1[6])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*2)) = int16(check_range(int(step1[2]) + int(step1[5])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*3)) = int16(check_range(int(step1[3]) + int(step1[4])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*4)) = int16(check_range(int(step1[3]) - int(step1[4])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*5)) = int16(check_range(int(step1[2]) - int(step1[5])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*6)) = int16(check_range(int(step1[1]) - int(step1[6])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*7)) = int16(check_range(int(step1[0]) - int(step1[7])))
}
func vpx_idct8x8_64_add_c(input *int16, dest *uint8, stride int) {
	var (
		i        int
		j        int
		out      [64]int16
		outptr   *int16 = &out[0]
		temp_in  [8]int16
		temp_out [8]int16
	)
	for i = 0; i < 8; i++ {
		idct8_c(input, outptr)
		input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*8))
		outptr = (*int16)(unsafe.Add(unsafe.Pointer(outptr), unsafe.Sizeof(int16(0))*8))
	}
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			temp_in[j] = out[j*8+i]
		}
		idct8_c(&temp_in[0], &temp_out[0])
		for j = 0; j < 8; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)), (int(temp_out[j])+(1<<(5-1)))>>5)
		}
	}
}
func vpx_idct8x8_12_add_c(input *int16, dest *uint8, stride int) {
	var (
		i        int
		j        int
		out      [64]int16 = [64]int16{}
		outptr   *int16    = &out[0]
		temp_in  [8]int16
		temp_out [8]int16
	)
	for i = 0; i < 4; i++ {
		idct8_c(input, outptr)
		input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*8))
		outptr = (*int16)(unsafe.Add(unsafe.Pointer(outptr), unsafe.Sizeof(int16(0))*8))
	}
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			temp_in[j] = out[j*8+i]
		}
		idct8_c(&temp_in[0], &temp_out[0])
		for j = 0; j < 8; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)), (int(temp_out[j])+(1<<(5-1)))>>5)
		}
	}
}
func vpx_idct8x8_1_add_c(input *int16, dest *uint8, stride int) {
	var (
		i   int
		j   int
		a1  int
		out int16 = int16(check_range(dct_const_round_shift(int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))) * int(cospi_16_64))))
	)
	out = int16(check_range(dct_const_round_shift(int(out) * int(cospi_16_64))))
	a1 = (int(out) + (1 << (5 - 1))) >> 5
	for j = 0; j < 8; j++ {
		for i = 0; i < 8; i++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), i)), a1)
		}
		dest = (*uint8)(unsafe.Add(unsafe.Pointer(dest), stride))
	}
}
func iadst16_c(input *int16, output *int16) {
	var (
		s0  int
		s1  int
		s2  int
		s3  int
		s4  int
		s5  int
		s6  int
		s7  int
		s8  int
		s9  int
		s10 int
		s11 int
		s12 int
		s13 int
		s14 int
		s15 int
		x0  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*15)))
		x1  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0)))
		x2  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*13)))
		x3  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*2)))
		x4  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*11)))
		x5  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*4)))
		x6  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*9)))
		x7  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*6)))
		x8  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*7)))
		x9  int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*8)))
		x10 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*5)))
		x11 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*10)))
		x12 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*3)))
		x13 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*12)))
		x14 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1)))
		x15 int = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*14)))
	)
	if (x0 | x1 | x2 | x3 | x4 | x5 | x6 | x7 | x8 | x9 | x10 | x11 | x12 | x13 | x14 | x15) == 0 {
		libc.MemSet(unsafe.Pointer(output), 0, int(16*unsafe.Sizeof(int16(0))))
		return
	}
	s0 = x0*int(cospi_1_64) + x1*int(cospi_31_64)
	s1 = x0*int(cospi_31_64) - x1*int(cospi_1_64)
	s2 = x2*int(cospi_5_64) + x3*int(cospi_27_64)
	s3 = x2*int(cospi_27_64) - x3*int(cospi_5_64)
	s4 = x4*int(cospi_9_64) + x5*int(cospi_23_64)
	s5 = x4*int(cospi_23_64) - x5*int(cospi_9_64)
	s6 = x6*int(cospi_13_64) + x7*int(cospi_19_64)
	s7 = x6*int(cospi_19_64) - x7*int(cospi_13_64)
	s8 = x8*int(cospi_17_64) + x9*int(cospi_15_64)
	s9 = x8*int(cospi_15_64) - x9*int(cospi_17_64)
	s10 = x10*int(cospi_21_64) + x11*int(cospi_11_64)
	s11 = x10*int(cospi_11_64) - x11*int(cospi_21_64)
	s12 = x12*int(cospi_25_64) + x13*int(cospi_7_64)
	s13 = x12*int(cospi_7_64) - x13*int(cospi_25_64)
	s14 = x14*int(cospi_29_64) + x15*int(cospi_3_64)
	s15 = x14*int(cospi_3_64) - x15*int(cospi_29_64)
	x0 = check_range(dct_const_round_shift(s0 + s8))
	x1 = check_range(dct_const_round_shift(s1 + s9))
	x2 = check_range(dct_const_round_shift(s2 + s10))
	x3 = check_range(dct_const_round_shift(s3 + s11))
	x4 = check_range(dct_const_round_shift(s4 + s12))
	x5 = check_range(dct_const_round_shift(s5 + s13))
	x6 = check_range(dct_const_round_shift(s6 + s14))
	x7 = check_range(dct_const_round_shift(s7 + s15))
	x8 = check_range(dct_const_round_shift(s0 - s8))
	x9 = check_range(dct_const_round_shift(s1 - s9))
	x10 = check_range(dct_const_round_shift(s2 - s10))
	x11 = check_range(dct_const_round_shift(s3 - s11))
	x12 = check_range(dct_const_round_shift(s4 - s12))
	x13 = check_range(dct_const_round_shift(s5 - s13))
	x14 = check_range(dct_const_round_shift(s6 - s14))
	x15 = check_range(dct_const_round_shift(s7 - s15))
	s0 = x0
	s1 = x1
	s2 = x2
	s3 = x3
	s4 = x4
	s5 = x5
	s6 = x6
	s7 = x7
	s8 = x8*int(cospi_4_64) + x9*int(cospi_28_64)
	s9 = x8*int(cospi_28_64) - x9*int(cospi_4_64)
	s10 = x10*int(cospi_20_64) + x11*int(cospi_12_64)
	s11 = x10*int(cospi_12_64) - x11*int(cospi_20_64)
	s12 = -x12*int(cospi_28_64) + x13*int(cospi_4_64)
	s13 = x12*int(cospi_4_64) + x13*int(cospi_28_64)
	s14 = -x14*int(cospi_12_64) + x15*int(cospi_20_64)
	s15 = x14*int(cospi_20_64) + x15*int(cospi_12_64)
	x0 = check_range(s0 + s4)
	x1 = check_range(s1 + s5)
	x2 = check_range(s2 + s6)
	x3 = check_range(s3 + s7)
	x4 = check_range(s0 - s4)
	x5 = check_range(s1 - s5)
	x6 = check_range(s2 - s6)
	x7 = check_range(s3 - s7)
	x8 = check_range(dct_const_round_shift(s8 + s12))
	x9 = check_range(dct_const_round_shift(s9 + s13))
	x10 = check_range(dct_const_round_shift(s10 + s14))
	x11 = check_range(dct_const_round_shift(s11 + s15))
	x12 = check_range(dct_const_round_shift(s8 - s12))
	x13 = check_range(dct_const_round_shift(s9 - s13))
	x14 = check_range(dct_const_round_shift(s10 - s14))
	x15 = check_range(dct_const_round_shift(s11 - s15))
	s0 = x0
	s1 = x1
	s2 = x2
	s3 = x3
	s4 = x4*int(cospi_8_64) + x5*int(cospi_24_64)
	s5 = x4*int(cospi_24_64) - x5*int(cospi_8_64)
	s6 = -x6*int(cospi_24_64) + x7*int(cospi_8_64)
	s7 = x6*int(cospi_8_64) + x7*int(cospi_24_64)
	s8 = x8
	s9 = x9
	s10 = x10
	s11 = x11
	s12 = x12*int(cospi_8_64) + x13*int(cospi_24_64)
	s13 = x12*int(cospi_24_64) - x13*int(cospi_8_64)
	s14 = -x14*int(cospi_24_64) + x15*int(cospi_8_64)
	s15 = x14*int(cospi_8_64) + x15*int(cospi_24_64)
	x0 = check_range(s0 + s2)
	x1 = check_range(s1 + s3)
	x2 = check_range(s0 - s2)
	x3 = check_range(s1 - s3)
	x4 = check_range(dct_const_round_shift(s4 + s6))
	x5 = check_range(dct_const_round_shift(s5 + s7))
	x6 = check_range(dct_const_round_shift(s4 - s6))
	x7 = check_range(dct_const_round_shift(s5 - s7))
	x8 = check_range(s8 + s10)
	x9 = check_range(s9 + s11)
	x10 = check_range(s8 - s10)
	x11 = check_range(s9 - s11)
	x12 = check_range(dct_const_round_shift(s12 + s14))
	x13 = check_range(dct_const_round_shift(s13 + s15))
	x14 = check_range(dct_const_round_shift(s12 - s14))
	x15 = check_range(dct_const_round_shift(s13 - s15))
	s2 = (int(-cospi_16_64)) * (x2 + x3)
	s3 = int(cospi_16_64) * (x2 - x3)
	s6 = int(cospi_16_64) * (x6 + x7)
	s7 = int(cospi_16_64) * (-x6 + x7)
	s10 = int(cospi_16_64) * (x10 + x11)
	s11 = int(cospi_16_64) * (-x10 + x11)
	s14 = (int(-cospi_16_64)) * (x14 + x15)
	s15 = int(cospi_16_64) * (x14 - x15)
	x2 = check_range(dct_const_round_shift(s2))
	x3 = check_range(dct_const_round_shift(s3))
	x6 = check_range(dct_const_round_shift(s6))
	x7 = check_range(dct_const_round_shift(s7))
	x10 = check_range(dct_const_round_shift(s10))
	x11 = check_range(dct_const_round_shift(s11))
	x14 = check_range(dct_const_round_shift(s14))
	x15 = check_range(dct_const_round_shift(s15))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(check_range(x0))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*1)) = int16(check_range(-x8))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*2)) = int16(check_range(x12))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*3)) = int16(check_range(-x4))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*4)) = int16(check_range(x6))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*5)) = int16(check_range(x14))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*6)) = int16(check_range(x10))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*7)) = int16(check_range(x2))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*8)) = int16(check_range(x3))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*9)) = int16(check_range(x11))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*10)) = int16(check_range(x15))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*11)) = int16(check_range(x7))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*12)) = int16(check_range(x5))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*13)) = int16(check_range(-x13))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*14)) = int16(check_range(x9))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*15)) = int16(check_range(-x1))
}
func idct16_c(input *int16, output *int16) {
	var (
		step1 [16]int16
		step2 [16]int16
		temp1 int
		temp2 int
	)
	step1[0] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(0/2)))
	step1[1] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(16/2)))
	step1[2] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(8/2)))
	step1[3] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(24/2)))
	step1[4] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(4/2)))
	step1[5] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(20/2)))
	step1[6] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(12/2)))
	step1[7] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(28/2)))
	step1[8] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(2/2)))
	step1[9] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(18/2)))
	step1[10] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(10/2)))
	step1[11] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(26/2)))
	step1[12] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(6/2)))
	step1[13] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(22/2)))
	step1[14] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(14/2)))
	step1[15] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*(30/2)))
	step2[0] = step1[0]
	step2[1] = step1[1]
	step2[2] = step1[2]
	step2[3] = step1[3]
	step2[4] = step1[4]
	step2[5] = step1[5]
	step2[6] = step1[6]
	step2[7] = step1[7]
	temp1 = int(step1[8])*int(cospi_30_64) - int(step1[15])*int(cospi_2_64)
	temp2 = int(step1[8])*int(cospi_2_64) + int(step1[15])*int(cospi_30_64)
	step2[8] = int16(check_range(dct_const_round_shift(temp1)))
	step2[15] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step1[9])*int(cospi_14_64) - int(step1[14])*int(cospi_18_64)
	temp2 = int(step1[9])*int(cospi_18_64) + int(step1[14])*int(cospi_14_64)
	step2[9] = int16(check_range(dct_const_round_shift(temp1)))
	step2[14] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step1[10])*int(cospi_22_64) - int(step1[13])*int(cospi_10_64)
	temp2 = int(step1[10])*int(cospi_10_64) + int(step1[13])*int(cospi_22_64)
	step2[10] = int16(check_range(dct_const_round_shift(temp1)))
	step2[13] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step1[11])*int(cospi_6_64) - int(step1[12])*int(cospi_26_64)
	temp2 = int(step1[11])*int(cospi_26_64) + int(step1[12])*int(cospi_6_64)
	step2[11] = int16(check_range(dct_const_round_shift(temp1)))
	step2[12] = int16(check_range(dct_const_round_shift(temp2)))
	step1[0] = step2[0]
	step1[1] = step2[1]
	step1[2] = step2[2]
	step1[3] = step2[3]
	temp1 = int(step2[4])*int(cospi_28_64) - int(step2[7])*int(cospi_4_64)
	temp2 = int(step2[4])*int(cospi_4_64) + int(step2[7])*int(cospi_28_64)
	step1[4] = int16(check_range(dct_const_round_shift(temp1)))
	step1[7] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step2[5])*int(cospi_12_64) - int(step2[6])*int(cospi_20_64)
	temp2 = int(step2[5])*int(cospi_20_64) + int(step2[6])*int(cospi_12_64)
	step1[5] = int16(check_range(dct_const_round_shift(temp1)))
	step1[6] = int16(check_range(dct_const_round_shift(temp2)))
	step1[8] = int16(check_range(int(step2[8]) + int(step2[9])))
	step1[9] = int16(check_range(int(step2[8]) - int(step2[9])))
	step1[10] = int16(check_range(int(-step2[10]) + int(step2[11])))
	step1[11] = int16(check_range(int(step2[10]) + int(step2[11])))
	step1[12] = int16(check_range(int(step2[12]) + int(step2[13])))
	step1[13] = int16(check_range(int(step2[12]) - int(step2[13])))
	step1[14] = int16(check_range(int(-step2[14]) + int(step2[15])))
	step1[15] = int16(check_range(int(step2[14]) + int(step2[15])))
	temp1 = (int(step1[0]) + int(step1[1])) * int(cospi_16_64)
	temp2 = (int(step1[0]) - int(step1[1])) * int(cospi_16_64)
	step2[0] = int16(check_range(dct_const_round_shift(temp1)))
	step2[1] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step1[2])*int(cospi_24_64) - int(step1[3])*int(cospi_8_64)
	temp2 = int(step1[2])*int(cospi_8_64) + int(step1[3])*int(cospi_24_64)
	step2[2] = int16(check_range(dct_const_round_shift(temp1)))
	step2[3] = int16(check_range(dct_const_round_shift(temp2)))
	step2[4] = int16(check_range(int(step1[4]) + int(step1[5])))
	step2[5] = int16(check_range(int(step1[4]) - int(step1[5])))
	step2[6] = int16(check_range(int(-step1[6]) + int(step1[7])))
	step2[7] = int16(check_range(int(step1[6]) + int(step1[7])))
	step2[8] = step1[8]
	step2[15] = step1[15]
	temp1 = int(-step1[9])*int(cospi_8_64) + int(step1[14])*int(cospi_24_64)
	temp2 = int(step1[9])*int(cospi_24_64) + int(step1[14])*int(cospi_8_64)
	step2[9] = int16(check_range(dct_const_round_shift(temp1)))
	step2[14] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(-step1[10])*int(cospi_24_64) - int(step1[13])*int(cospi_8_64)
	temp2 = int(-step1[10])*int(cospi_8_64) + int(step1[13])*int(cospi_24_64)
	step2[10] = int16(check_range(dct_const_round_shift(temp1)))
	step2[13] = int16(check_range(dct_const_round_shift(temp2)))
	step2[11] = step1[11]
	step2[12] = step1[12]
	step1[0] = int16(check_range(int(step2[0]) + int(step2[3])))
	step1[1] = int16(check_range(int(step2[1]) + int(step2[2])))
	step1[2] = int16(check_range(int(step2[1]) - int(step2[2])))
	step1[3] = int16(check_range(int(step2[0]) - int(step2[3])))
	step1[4] = step2[4]
	temp1 = (int(step2[6]) - int(step2[5])) * int(cospi_16_64)
	temp2 = (int(step2[5]) + int(step2[6])) * int(cospi_16_64)
	step1[5] = int16(check_range(dct_const_round_shift(temp1)))
	step1[6] = int16(check_range(dct_const_round_shift(temp2)))
	step1[7] = step2[7]
	step1[8] = int16(check_range(int(step2[8]) + int(step2[11])))
	step1[9] = int16(check_range(int(step2[9]) + int(step2[10])))
	step1[10] = int16(check_range(int(step2[9]) - int(step2[10])))
	step1[11] = int16(check_range(int(step2[8]) - int(step2[11])))
	step1[12] = int16(check_range(int(-step2[12]) + int(step2[15])))
	step1[13] = int16(check_range(int(-step2[13]) + int(step2[14])))
	step1[14] = int16(check_range(int(step2[13]) + int(step2[14])))
	step1[15] = int16(check_range(int(step2[12]) + int(step2[15])))
	step2[0] = int16(check_range(int(step1[0]) + int(step1[7])))
	step2[1] = int16(check_range(int(step1[1]) + int(step1[6])))
	step2[2] = int16(check_range(int(step1[2]) + int(step1[5])))
	step2[3] = int16(check_range(int(step1[3]) + int(step1[4])))
	step2[4] = int16(check_range(int(step1[3]) - int(step1[4])))
	step2[5] = int16(check_range(int(step1[2]) - int(step1[5])))
	step2[6] = int16(check_range(int(step1[1]) - int(step1[6])))
	step2[7] = int16(check_range(int(step1[0]) - int(step1[7])))
	step2[8] = step1[8]
	step2[9] = step1[9]
	temp1 = (int(-step1[10]) + int(step1[13])) * int(cospi_16_64)
	temp2 = (int(step1[10]) + int(step1[13])) * int(cospi_16_64)
	step2[10] = int16(check_range(dct_const_round_shift(temp1)))
	step2[13] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = (int(-step1[11]) + int(step1[12])) * int(cospi_16_64)
	temp2 = (int(step1[11]) + int(step1[12])) * int(cospi_16_64)
	step2[11] = int16(check_range(dct_const_round_shift(temp1)))
	step2[12] = int16(check_range(dct_const_round_shift(temp2)))
	step2[14] = step1[14]
	step2[15] = step1[15]
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(check_range(int(step2[0]) + int(step2[15])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*1)) = int16(check_range(int(step2[1]) + int(step2[14])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*2)) = int16(check_range(int(step2[2]) + int(step2[13])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*3)) = int16(check_range(int(step2[3]) + int(step2[12])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*4)) = int16(check_range(int(step2[4]) + int(step2[11])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*5)) = int16(check_range(int(step2[5]) + int(step2[10])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*6)) = int16(check_range(int(step2[6]) + int(step2[9])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*7)) = int16(check_range(int(step2[7]) + int(step2[8])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*8)) = int16(check_range(int(step2[7]) - int(step2[8])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*9)) = int16(check_range(int(step2[6]) - int(step2[9])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*10)) = int16(check_range(int(step2[5]) - int(step2[10])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*11)) = int16(check_range(int(step2[4]) - int(step2[11])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*12)) = int16(check_range(int(step2[3]) - int(step2[12])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*13)) = int16(check_range(int(step2[2]) - int(step2[13])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*14)) = int16(check_range(int(step2[1]) - int(step2[14])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*15)) = int16(check_range(int(step2[0]) - int(step2[15])))
}
func vpx_idct16x16_256_add_c(input *int16, dest *uint8, stride int) {
	var (
		i        int
		j        int
		out      [256]int16
		outptr   *int16 = &out[0]
		temp_in  [16]int16
		temp_out [16]int16
	)
	for i = 0; i < 16; i++ {
		idct16_c(input, outptr)
		input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*16))
		outptr = (*int16)(unsafe.Add(unsafe.Pointer(outptr), unsafe.Sizeof(int16(0))*16))
	}
	for i = 0; i < 16; i++ {
		for j = 0; j < 16; j++ {
			temp_in[j] = out[j*16+i]
		}
		idct16_c(&temp_in[0], &temp_out[0])
		for j = 0; j < 16; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)), (int(temp_out[j])+(1<<(6-1)))>>6)
		}
	}
}
func vpx_idct16x16_38_add_c(input *int16, dest *uint8, stride int) {
	var (
		i        int
		j        int
		out      [256]int16 = [256]int16{}
		outptr   *int16     = &out[0]
		temp_in  [16]int16
		temp_out [16]int16
	)
	for i = 0; i < 8; i++ {
		idct16_c(input, outptr)
		input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*16))
		outptr = (*int16)(unsafe.Add(unsafe.Pointer(outptr), unsafe.Sizeof(int16(0))*16))
	}
	for i = 0; i < 16; i++ {
		for j = 0; j < 16; j++ {
			temp_in[j] = out[j*16+i]
		}
		idct16_c(&temp_in[0], &temp_out[0])
		for j = 0; j < 16; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)), (int(temp_out[j])+(1<<(6-1)))>>6)
		}
	}
}
func vpx_idct16x16_10_add_c(input *int16, dest *uint8, stride int) {
	var (
		i        int
		j        int
		out      [256]int16 = [256]int16{}
		outptr   *int16     = &out[0]
		temp_in  [16]int16
		temp_out [16]int16
	)
	for i = 0; i < 4; i++ {
		idct16_c(input, outptr)
		input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*16))
		outptr = (*int16)(unsafe.Add(unsafe.Pointer(outptr), unsafe.Sizeof(int16(0))*16))
	}
	for i = 0; i < 16; i++ {
		for j = 0; j < 16; j++ {
			temp_in[j] = out[j*16+i]
		}
		idct16_c(&temp_in[0], &temp_out[0])
		for j = 0; j < 16; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)), (int(temp_out[j])+(1<<(6-1)))>>6)
		}
	}
}
func vpx_idct16x16_1_add_c(input *int16, dest *uint8, stride int) {
	var (
		i   int
		j   int
		a1  int
		out int16 = int16(check_range(dct_const_round_shift(int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))) * int(cospi_16_64))))
	)
	out = int16(check_range(dct_const_round_shift(int(out) * int(cospi_16_64))))
	a1 = (int(out) + (1 << (6 - 1))) >> 6
	for j = 0; j < 16; j++ {
		for i = 0; i < 16; i++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), i)), a1)
		}
		dest = (*uint8)(unsafe.Add(unsafe.Pointer(dest), stride))
	}
}
func idct32_c(input *int16, output *int16) {
	var (
		step1 [32]int16
		step2 [32]int16
		temp1 int
		temp2 int
	)
	step1[0] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))
	step1[1] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*16))
	step1[2] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*8))
	step1[3] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*24))
	step1[4] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*4))
	step1[5] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*20))
	step1[6] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*12))
	step1[7] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*28))
	step1[8] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*2))
	step1[9] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*18))
	step1[10] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*10))
	step1[11] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*26))
	step1[12] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*6))
	step1[13] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*22))
	step1[14] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*14))
	step1[15] = *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*30))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1)))*int(cospi_31_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*31)))*int(cospi_1_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1)))*int(cospi_1_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*31)))*int(cospi_31_64)
	step1[16] = int16(check_range(dct_const_round_shift(temp1)))
	step1[31] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*17)))*int(cospi_15_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*15)))*int(cospi_17_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*17)))*int(cospi_17_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*15)))*int(cospi_15_64)
	step1[17] = int16(check_range(dct_const_round_shift(temp1)))
	step1[30] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*9)))*int(cospi_23_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*23)))*int(cospi_9_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*9)))*int(cospi_9_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*23)))*int(cospi_23_64)
	step1[18] = int16(check_range(dct_const_round_shift(temp1)))
	step1[29] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*25)))*int(cospi_7_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*7)))*int(cospi_25_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*25)))*int(cospi_25_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*7)))*int(cospi_7_64)
	step1[19] = int16(check_range(dct_const_round_shift(temp1)))
	step1[28] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*5)))*int(cospi_27_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*27)))*int(cospi_5_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*5)))*int(cospi_5_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*27)))*int(cospi_27_64)
	step1[20] = int16(check_range(dct_const_round_shift(temp1)))
	step1[27] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*21)))*int(cospi_11_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*11)))*int(cospi_21_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*21)))*int(cospi_21_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*11)))*int(cospi_11_64)
	step1[21] = int16(check_range(dct_const_round_shift(temp1)))
	step1[26] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*13)))*int(cospi_19_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*19)))*int(cospi_13_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*13)))*int(cospi_13_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*19)))*int(cospi_19_64)
	step1[22] = int16(check_range(dct_const_round_shift(temp1)))
	step1[25] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*29)))*int(cospi_3_64) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*3)))*int(cospi_29_64)
	temp2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*29)))*int(cospi_29_64) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*3)))*int(cospi_3_64)
	step1[23] = int16(check_range(dct_const_round_shift(temp1)))
	step1[24] = int16(check_range(dct_const_round_shift(temp2)))
	step2[0] = step1[0]
	step2[1] = step1[1]
	step2[2] = step1[2]
	step2[3] = step1[3]
	step2[4] = step1[4]
	step2[5] = step1[5]
	step2[6] = step1[6]
	step2[7] = step1[7]
	temp1 = int(step1[8])*int(cospi_30_64) - int(step1[15])*int(cospi_2_64)
	temp2 = int(step1[8])*int(cospi_2_64) + int(step1[15])*int(cospi_30_64)
	step2[8] = int16(check_range(dct_const_round_shift(temp1)))
	step2[15] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step1[9])*int(cospi_14_64) - int(step1[14])*int(cospi_18_64)
	temp2 = int(step1[9])*int(cospi_18_64) + int(step1[14])*int(cospi_14_64)
	step2[9] = int16(check_range(dct_const_round_shift(temp1)))
	step2[14] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step1[10])*int(cospi_22_64) - int(step1[13])*int(cospi_10_64)
	temp2 = int(step1[10])*int(cospi_10_64) + int(step1[13])*int(cospi_22_64)
	step2[10] = int16(check_range(dct_const_round_shift(temp1)))
	step2[13] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step1[11])*int(cospi_6_64) - int(step1[12])*int(cospi_26_64)
	temp2 = int(step1[11])*int(cospi_26_64) + int(step1[12])*int(cospi_6_64)
	step2[11] = int16(check_range(dct_const_round_shift(temp1)))
	step2[12] = int16(check_range(dct_const_round_shift(temp2)))
	step2[16] = int16(check_range(int(step1[16]) + int(step1[17])))
	step2[17] = int16(check_range(int(step1[16]) - int(step1[17])))
	step2[18] = int16(check_range(int(-step1[18]) + int(step1[19])))
	step2[19] = int16(check_range(int(step1[18]) + int(step1[19])))
	step2[20] = int16(check_range(int(step1[20]) + int(step1[21])))
	step2[21] = int16(check_range(int(step1[20]) - int(step1[21])))
	step2[22] = int16(check_range(int(-step1[22]) + int(step1[23])))
	step2[23] = int16(check_range(int(step1[22]) + int(step1[23])))
	step2[24] = int16(check_range(int(step1[24]) + int(step1[25])))
	step2[25] = int16(check_range(int(step1[24]) - int(step1[25])))
	step2[26] = int16(check_range(int(-step1[26]) + int(step1[27])))
	step2[27] = int16(check_range(int(step1[26]) + int(step1[27])))
	step2[28] = int16(check_range(int(step1[28]) + int(step1[29])))
	step2[29] = int16(check_range(int(step1[28]) - int(step1[29])))
	step2[30] = int16(check_range(int(-step1[30]) + int(step1[31])))
	step2[31] = int16(check_range(int(step1[30]) + int(step1[31])))
	step1[0] = step2[0]
	step1[1] = step2[1]
	step1[2] = step2[2]
	step1[3] = step2[3]
	temp1 = int(step2[4])*int(cospi_28_64) - int(step2[7])*int(cospi_4_64)
	temp2 = int(step2[4])*int(cospi_4_64) + int(step2[7])*int(cospi_28_64)
	step1[4] = int16(check_range(dct_const_round_shift(temp1)))
	step1[7] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step2[5])*int(cospi_12_64) - int(step2[6])*int(cospi_20_64)
	temp2 = int(step2[5])*int(cospi_20_64) + int(step2[6])*int(cospi_12_64)
	step1[5] = int16(check_range(dct_const_round_shift(temp1)))
	step1[6] = int16(check_range(dct_const_round_shift(temp2)))
	step1[8] = int16(check_range(int(step2[8]) + int(step2[9])))
	step1[9] = int16(check_range(int(step2[8]) - int(step2[9])))
	step1[10] = int16(check_range(int(-step2[10]) + int(step2[11])))
	step1[11] = int16(check_range(int(step2[10]) + int(step2[11])))
	step1[12] = int16(check_range(int(step2[12]) + int(step2[13])))
	step1[13] = int16(check_range(int(step2[12]) - int(step2[13])))
	step1[14] = int16(check_range(int(-step2[14]) + int(step2[15])))
	step1[15] = int16(check_range(int(step2[14]) + int(step2[15])))
	step1[16] = step2[16]
	step1[31] = step2[31]
	temp1 = int(-step2[17])*int(cospi_4_64) + int(step2[30])*int(cospi_28_64)
	temp2 = int(step2[17])*int(cospi_28_64) + int(step2[30])*int(cospi_4_64)
	step1[17] = int16(check_range(dct_const_round_shift(temp1)))
	step1[30] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(-step2[18])*int(cospi_28_64) - int(step2[29])*int(cospi_4_64)
	temp2 = int(-step2[18])*int(cospi_4_64) + int(step2[29])*int(cospi_28_64)
	step1[18] = int16(check_range(dct_const_round_shift(temp1)))
	step1[29] = int16(check_range(dct_const_round_shift(temp2)))
	step1[19] = step2[19]
	step1[20] = step2[20]
	temp1 = int(-step2[21])*int(cospi_20_64) + int(step2[26])*int(cospi_12_64)
	temp2 = int(step2[21])*int(cospi_12_64) + int(step2[26])*int(cospi_20_64)
	step1[21] = int16(check_range(dct_const_round_shift(temp1)))
	step1[26] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(-step2[22])*int(cospi_12_64) - int(step2[25])*int(cospi_20_64)
	temp2 = int(-step2[22])*int(cospi_20_64) + int(step2[25])*int(cospi_12_64)
	step1[22] = int16(check_range(dct_const_round_shift(temp1)))
	step1[25] = int16(check_range(dct_const_round_shift(temp2)))
	step1[23] = step2[23]
	step1[24] = step2[24]
	step1[27] = step2[27]
	step1[28] = step2[28]
	temp1 = (int(step1[0]) + int(step1[1])) * int(cospi_16_64)
	temp2 = (int(step1[0]) - int(step1[1])) * int(cospi_16_64)
	step2[0] = int16(check_range(dct_const_round_shift(temp1)))
	step2[1] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(step1[2])*int(cospi_24_64) - int(step1[3])*int(cospi_8_64)
	temp2 = int(step1[2])*int(cospi_8_64) + int(step1[3])*int(cospi_24_64)
	step2[2] = int16(check_range(dct_const_round_shift(temp1)))
	step2[3] = int16(check_range(dct_const_round_shift(temp2)))
	step2[4] = int16(check_range(int(step1[4]) + int(step1[5])))
	step2[5] = int16(check_range(int(step1[4]) - int(step1[5])))
	step2[6] = int16(check_range(int(-step1[6]) + int(step1[7])))
	step2[7] = int16(check_range(int(step1[6]) + int(step1[7])))
	step2[8] = step1[8]
	step2[15] = step1[15]
	temp1 = int(-step1[9])*int(cospi_8_64) + int(step1[14])*int(cospi_24_64)
	temp2 = int(step1[9])*int(cospi_24_64) + int(step1[14])*int(cospi_8_64)
	step2[9] = int16(check_range(dct_const_round_shift(temp1)))
	step2[14] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(-step1[10])*int(cospi_24_64) - int(step1[13])*int(cospi_8_64)
	temp2 = int(-step1[10])*int(cospi_8_64) + int(step1[13])*int(cospi_24_64)
	step2[10] = int16(check_range(dct_const_round_shift(temp1)))
	step2[13] = int16(check_range(dct_const_round_shift(temp2)))
	step2[11] = step1[11]
	step2[12] = step1[12]
	step2[16] = int16(check_range(int(step1[16]) + int(step1[19])))
	step2[17] = int16(check_range(int(step1[17]) + int(step1[18])))
	step2[18] = int16(check_range(int(step1[17]) - int(step1[18])))
	step2[19] = int16(check_range(int(step1[16]) - int(step1[19])))
	step2[20] = int16(check_range(int(-step1[20]) + int(step1[23])))
	step2[21] = int16(check_range(int(-step1[21]) + int(step1[22])))
	step2[22] = int16(check_range(int(step1[21]) + int(step1[22])))
	step2[23] = int16(check_range(int(step1[20]) + int(step1[23])))
	step2[24] = int16(check_range(int(step1[24]) + int(step1[27])))
	step2[25] = int16(check_range(int(step1[25]) + int(step1[26])))
	step2[26] = int16(check_range(int(step1[25]) - int(step1[26])))
	step2[27] = int16(check_range(int(step1[24]) - int(step1[27])))
	step2[28] = int16(check_range(int(-step1[28]) + int(step1[31])))
	step2[29] = int16(check_range(int(-step1[29]) + int(step1[30])))
	step2[30] = int16(check_range(int(step1[29]) + int(step1[30])))
	step2[31] = int16(check_range(int(step1[28]) + int(step1[31])))
	step1[0] = int16(check_range(int(step2[0]) + int(step2[3])))
	step1[1] = int16(check_range(int(step2[1]) + int(step2[2])))
	step1[2] = int16(check_range(int(step2[1]) - int(step2[2])))
	step1[3] = int16(check_range(int(step2[0]) - int(step2[3])))
	step1[4] = step2[4]
	temp1 = (int(step2[6]) - int(step2[5])) * int(cospi_16_64)
	temp2 = (int(step2[5]) + int(step2[6])) * int(cospi_16_64)
	step1[5] = int16(check_range(dct_const_round_shift(temp1)))
	step1[6] = int16(check_range(dct_const_round_shift(temp2)))
	step1[7] = step2[7]
	step1[8] = int16(check_range(int(step2[8]) + int(step2[11])))
	step1[9] = int16(check_range(int(step2[9]) + int(step2[10])))
	step1[10] = int16(check_range(int(step2[9]) - int(step2[10])))
	step1[11] = int16(check_range(int(step2[8]) - int(step2[11])))
	step1[12] = int16(check_range(int(-step2[12]) + int(step2[15])))
	step1[13] = int16(check_range(int(-step2[13]) + int(step2[14])))
	step1[14] = int16(check_range(int(step2[13]) + int(step2[14])))
	step1[15] = int16(check_range(int(step2[12]) + int(step2[15])))
	step1[16] = step2[16]
	step1[17] = step2[17]
	temp1 = int(-step2[18])*int(cospi_8_64) + int(step2[29])*int(cospi_24_64)
	temp2 = int(step2[18])*int(cospi_24_64) + int(step2[29])*int(cospi_8_64)
	step1[18] = int16(check_range(dct_const_round_shift(temp1)))
	step1[29] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(-step2[19])*int(cospi_8_64) + int(step2[28])*int(cospi_24_64)
	temp2 = int(step2[19])*int(cospi_24_64) + int(step2[28])*int(cospi_8_64)
	step1[19] = int16(check_range(dct_const_round_shift(temp1)))
	step1[28] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(-step2[20])*int(cospi_24_64) - int(step2[27])*int(cospi_8_64)
	temp2 = int(-step2[20])*int(cospi_8_64) + int(step2[27])*int(cospi_24_64)
	step1[20] = int16(check_range(dct_const_round_shift(temp1)))
	step1[27] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = int(-step2[21])*int(cospi_24_64) - int(step2[26])*int(cospi_8_64)
	temp2 = int(-step2[21])*int(cospi_8_64) + int(step2[26])*int(cospi_24_64)
	step1[21] = int16(check_range(dct_const_round_shift(temp1)))
	step1[26] = int16(check_range(dct_const_round_shift(temp2)))
	step1[22] = step2[22]
	step1[23] = step2[23]
	step1[24] = step2[24]
	step1[25] = step2[25]
	step1[30] = step2[30]
	step1[31] = step2[31]
	step2[0] = int16(check_range(int(step1[0]) + int(step1[7])))
	step2[1] = int16(check_range(int(step1[1]) + int(step1[6])))
	step2[2] = int16(check_range(int(step1[2]) + int(step1[5])))
	step2[3] = int16(check_range(int(step1[3]) + int(step1[4])))
	step2[4] = int16(check_range(int(step1[3]) - int(step1[4])))
	step2[5] = int16(check_range(int(step1[2]) - int(step1[5])))
	step2[6] = int16(check_range(int(step1[1]) - int(step1[6])))
	step2[7] = int16(check_range(int(step1[0]) - int(step1[7])))
	step2[8] = step1[8]
	step2[9] = step1[9]
	temp1 = (int(-step1[10]) + int(step1[13])) * int(cospi_16_64)
	temp2 = (int(step1[10]) + int(step1[13])) * int(cospi_16_64)
	step2[10] = int16(check_range(dct_const_round_shift(temp1)))
	step2[13] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = (int(-step1[11]) + int(step1[12])) * int(cospi_16_64)
	temp2 = (int(step1[11]) + int(step1[12])) * int(cospi_16_64)
	step2[11] = int16(check_range(dct_const_round_shift(temp1)))
	step2[12] = int16(check_range(dct_const_round_shift(temp2)))
	step2[14] = step1[14]
	step2[15] = step1[15]
	step2[16] = int16(check_range(int(step1[16]) + int(step1[23])))
	step2[17] = int16(check_range(int(step1[17]) + int(step1[22])))
	step2[18] = int16(check_range(int(step1[18]) + int(step1[21])))
	step2[19] = int16(check_range(int(step1[19]) + int(step1[20])))
	step2[20] = int16(check_range(int(step1[19]) - int(step1[20])))
	step2[21] = int16(check_range(int(step1[18]) - int(step1[21])))
	step2[22] = int16(check_range(int(step1[17]) - int(step1[22])))
	step2[23] = int16(check_range(int(step1[16]) - int(step1[23])))
	step2[24] = int16(check_range(int(-step1[24]) + int(step1[31])))
	step2[25] = int16(check_range(int(-step1[25]) + int(step1[30])))
	step2[26] = int16(check_range(int(-step1[26]) + int(step1[29])))
	step2[27] = int16(check_range(int(-step1[27]) + int(step1[28])))
	step2[28] = int16(check_range(int(step1[27]) + int(step1[28])))
	step2[29] = int16(check_range(int(step1[26]) + int(step1[29])))
	step2[30] = int16(check_range(int(step1[25]) + int(step1[30])))
	step2[31] = int16(check_range(int(step1[24]) + int(step1[31])))
	step1[0] = int16(check_range(int(step2[0]) + int(step2[15])))
	step1[1] = int16(check_range(int(step2[1]) + int(step2[14])))
	step1[2] = int16(check_range(int(step2[2]) + int(step2[13])))
	step1[3] = int16(check_range(int(step2[3]) + int(step2[12])))
	step1[4] = int16(check_range(int(step2[4]) + int(step2[11])))
	step1[5] = int16(check_range(int(step2[5]) + int(step2[10])))
	step1[6] = int16(check_range(int(step2[6]) + int(step2[9])))
	step1[7] = int16(check_range(int(step2[7]) + int(step2[8])))
	step1[8] = int16(check_range(int(step2[7]) - int(step2[8])))
	step1[9] = int16(check_range(int(step2[6]) - int(step2[9])))
	step1[10] = int16(check_range(int(step2[5]) - int(step2[10])))
	step1[11] = int16(check_range(int(step2[4]) - int(step2[11])))
	step1[12] = int16(check_range(int(step2[3]) - int(step2[12])))
	step1[13] = int16(check_range(int(step2[2]) - int(step2[13])))
	step1[14] = int16(check_range(int(step2[1]) - int(step2[14])))
	step1[15] = int16(check_range(int(step2[0]) - int(step2[15])))
	step1[16] = step2[16]
	step1[17] = step2[17]
	step1[18] = step2[18]
	step1[19] = step2[19]
	temp1 = (int(-step2[20]) + int(step2[27])) * int(cospi_16_64)
	temp2 = (int(step2[20]) + int(step2[27])) * int(cospi_16_64)
	step1[20] = int16(check_range(dct_const_round_shift(temp1)))
	step1[27] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = (int(-step2[21]) + int(step2[26])) * int(cospi_16_64)
	temp2 = (int(step2[21]) + int(step2[26])) * int(cospi_16_64)
	step1[21] = int16(check_range(dct_const_round_shift(temp1)))
	step1[26] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = (int(-step2[22]) + int(step2[25])) * int(cospi_16_64)
	temp2 = (int(step2[22]) + int(step2[25])) * int(cospi_16_64)
	step1[22] = int16(check_range(dct_const_round_shift(temp1)))
	step1[25] = int16(check_range(dct_const_round_shift(temp2)))
	temp1 = (int(-step2[23]) + int(step2[24])) * int(cospi_16_64)
	temp2 = (int(step2[23]) + int(step2[24])) * int(cospi_16_64)
	step1[23] = int16(check_range(dct_const_round_shift(temp1)))
	step1[24] = int16(check_range(dct_const_round_shift(temp2)))
	step1[28] = step2[28]
	step1[29] = step2[29]
	step1[30] = step2[30]
	step1[31] = step2[31]
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(check_range(int(step1[0]) + int(step1[31])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*1)) = int16(check_range(int(step1[1]) + int(step1[30])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*2)) = int16(check_range(int(step1[2]) + int(step1[29])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*3)) = int16(check_range(int(step1[3]) + int(step1[28])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*4)) = int16(check_range(int(step1[4]) + int(step1[27])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*5)) = int16(check_range(int(step1[5]) + int(step1[26])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*6)) = int16(check_range(int(step1[6]) + int(step1[25])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*7)) = int16(check_range(int(step1[7]) + int(step1[24])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*8)) = int16(check_range(int(step1[8]) + int(step1[23])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*9)) = int16(check_range(int(step1[9]) + int(step1[22])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*10)) = int16(check_range(int(step1[10]) + int(step1[21])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*11)) = int16(check_range(int(step1[11]) + int(step1[20])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*12)) = int16(check_range(int(step1[12]) + int(step1[19])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*13)) = int16(check_range(int(step1[13]) + int(step1[18])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*14)) = int16(check_range(int(step1[14]) + int(step1[17])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*15)) = int16(check_range(int(step1[15]) + int(step1[16])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*16)) = int16(check_range(int(step1[15]) - int(step1[16])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*17)) = int16(check_range(int(step1[14]) - int(step1[17])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*18)) = int16(check_range(int(step1[13]) - int(step1[18])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*19)) = int16(check_range(int(step1[12]) - int(step1[19])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*20)) = int16(check_range(int(step1[11]) - int(step1[20])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*21)) = int16(check_range(int(step1[10]) - int(step1[21])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*22)) = int16(check_range(int(step1[9]) - int(step1[22])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*23)) = int16(check_range(int(step1[8]) - int(step1[23])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*24)) = int16(check_range(int(step1[7]) - int(step1[24])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*25)) = int16(check_range(int(step1[6]) - int(step1[25])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*26)) = int16(check_range(int(step1[5]) - int(step1[26])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*27)) = int16(check_range(int(step1[4]) - int(step1[27])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*28)) = int16(check_range(int(step1[3]) - int(step1[28])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*29)) = int16(check_range(int(step1[2]) - int(step1[29])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*30)) = int16(check_range(int(step1[1]) - int(step1[30])))
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*31)) = int16(check_range(int(step1[0]) - int(step1[31])))
}
func vpx_idct32x32_1024_add_c(input *int16, dest *uint8, stride int) {
	var (
		i        int
		j        int
		out      [1024]int16
		outptr   *int16 = &out[0]
		temp_in  [32]int16
		temp_out [32]int16
	)
	for i = 0; i < 32; i++ {
		var zero_coeff int16 = 0
		for j = 0; j < 32; j++ {
			zero_coeff |= *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(j)))
		}
		if int(zero_coeff) != 0 {
			idct32_c(input, outptr)
		} else {
			libc.MemSet(unsafe.Pointer(outptr), 0, int(unsafe.Sizeof(int16(0))*32))
		}
		input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*32))
		outptr = (*int16)(unsafe.Add(unsafe.Pointer(outptr), unsafe.Sizeof(int16(0))*32))
	}
	for i = 0; i < 32; i++ {
		for j = 0; j < 32; j++ {
			temp_in[j] = out[j*32+i]
		}
		idct32_c(&temp_in[0], &temp_out[0])
		for j = 0; j < 32; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)), (int(temp_out[j])+(1<<(6-1)))>>6)
		}
	}
}
func vpx_idct32x32_135_add_c(input *int16, dest *uint8, stride int) {
	var (
		i        int
		j        int
		out      [1024]int16 = [1024]int16{}
		outptr   *int16      = &out[0]
		temp_in  [32]int16
		temp_out [32]int16
	)
	for i = 0; i < 16; i++ {
		idct32_c(input, outptr)
		input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*32))
		outptr = (*int16)(unsafe.Add(unsafe.Pointer(outptr), unsafe.Sizeof(int16(0))*32))
	}
	for i = 0; i < 32; i++ {
		for j = 0; j < 32; j++ {
			temp_in[j] = out[j*32+i]
		}
		idct32_c(&temp_in[0], &temp_out[0])
		for j = 0; j < 32; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)), (int(temp_out[j])+(1<<(6-1)))>>6)
		}
	}
}
func vpx_idct32x32_34_add_c(input *int16, dest *uint8, stride int) {
	var (
		i        int
		j        int
		out      [1024]int16 = [1024]int16{}
		outptr   *int16      = &out[0]
		temp_in  [32]int16
		temp_out [32]int16
	)
	for i = 0; i < 8; i++ {
		idct32_c(input, outptr)
		input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*32))
		outptr = (*int16)(unsafe.Add(unsafe.Pointer(outptr), unsafe.Sizeof(int16(0))*32))
	}
	for i = 0; i < 32; i++ {
		for j = 0; j < 32; j++ {
			temp_in[j] = out[j*32+i]
		}
		idct32_c(&temp_in[0], &temp_out[0])
		for j = 0; j < 32; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), j*stride+i)), (int(temp_out[j])+(1<<(6-1)))>>6)
		}
	}
}
func vpx_idct32x32_1_add_c(input *int16, dest *uint8, stride int) {
	var (
		i   int
		j   int
		a1  int
		out int16 = int16(check_range(dct_const_round_shift(int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*0))) * int(cospi_16_64))))
	)
	out = int16(check_range(dct_const_round_shift(int(out) * int(cospi_16_64))))
	a1 = (int(out) + (1 << (6 - 1))) >> 6
	for j = 0; j < 32; j++ {
		for i = 0; i < 32; i++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dest), i)) = clip_pixel_add(*(*uint8)(unsafe.Add(unsafe.Pointer(dest), i)), a1)
		}
		dest = (*uint8)(unsafe.Add(unsafe.Pointer(dest), stride))
	}
}
