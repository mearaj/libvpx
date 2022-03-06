package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func fdct_round_shift(input int) int {
	var rv int = ((input + (1 << (int(DCT_CONST_BITS - 1)))) >> DCT_CONST_BITS)
	return rv
}
func vpx_fdct4x4_c(input *int16, output *int16, stride int) {
	var (
		pass         int
		intermediate [16]int16
		in_low       *int16 = nil
		out          *int16 = &intermediate[0]
	)
	for pass = 0; pass < 2; pass++ {
		var (
			in_high [4]int
			step    [4]int
			temp1   int
			temp2   int
			i       int
		)
		for i = 0; i < 4; i++ {
			if pass == 0 {
				in_high[0] = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*0)))) * 16
				in_high[1] = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*1)))) * 16
				in_high[2] = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*2)))) * 16
				in_high[3] = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*3)))) * 16
				if i == 0 && in_high[0] != 0 {
					in_high[0]++
				}
			} else {
				if in_low != nil {
				} else {
					__assert_fail(libc.CString("in_low != NULL"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
				}
				in_high[0] = int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(0*4))))
				in_high[1] = int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(1*4))))
				in_high[2] = int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(2*4))))
				in_high[3] = int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(3*4))))
				in_low = (*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*1))
			}
			step[0] = in_high[0] + in_high[3]
			step[1] = in_high[1] + in_high[2]
			step[2] = in_high[1] - in_high[2]
			step[3] = in_high[0] - in_high[3]
			temp1 = (step[0] + step[1]) * int(cospi_16_64)
			temp2 = (step[0] - step[1]) * int(cospi_16_64)
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*0)) = int16(fdct_round_shift(temp1))
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*2)) = int16(fdct_round_shift(temp2))
			temp1 = step[2]*int(cospi_24_64) + step[3]*int(cospi_8_64)
			temp2 = -step[2]*int(cospi_8_64) + step[3]*int(cospi_24_64)
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*1)) = int16(fdct_round_shift(temp1))
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*3)) = int16(fdct_round_shift(temp2))
			input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1))
			out = (*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*4))
		}
		in_low = &intermediate[0]
		out = output
	}
	{
		var (
			i int
			j int
		)
		for i = 0; i < 4; i++ {
			for j = 0; j < 4; j++ {
				*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*uintptr(j+i*4))) = int16((int(*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*uintptr(j+i*4)))) + 1) >> 2)
			}
		}
	}
}
func vpx_fdct4x4_1_c(input *int16, output *int16, stride int) {
	var (
		r   int
		c   int
		sum int16 = 0
	)
	for r = 0; r < 4; r++ {
		for c = 0; c < 4; c++ {
			sum += *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(r*stride+c)))
		}
	}
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(int(sum) * 2)
}
func vpx_fdct8x8_c(input *int16, output *int16, stride int) {
	var (
		i            int
		j            int
		intermediate [64]int16
		pass         int
		out          *int16 = &intermediate[0]
		in           *int16 = nil
	)
	for pass = 0; pass < 2; pass++ {
		var (
			s0 int
			s1 int
			s2 int
			s3 int
			s4 int
			s5 int
			s6 int
			s7 int
			t0 int
			t1 int
			t2 int
			t3 int
			x0 int
			x1 int
			x2 int
			x3 int
		)
		for i = 0; i < 8; i++ {
			if pass == 0 {
				s0 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*0)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*7))))) * 4
				s1 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*1)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*6))))) * 4
				s2 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*2)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*5))))) * 4
				s3 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*3)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*4))))) * 4
				s4 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*3)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*4))))) * 4
				s5 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*2)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*5))))) * 4
				s6 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*1)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*6))))) * 4
				s7 = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*0)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*7))))) * 4
				input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1))
			} else {
				s0 = int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(0*8)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(7*8))))
				s1 = int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(1*8)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(6*8))))
				s2 = int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(2*8)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(5*8))))
				s3 = int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(3*8)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(4*8))))
				s4 = int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(3*8)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(4*8))))
				s5 = int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(2*8)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(5*8))))
				s6 = int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(1*8)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(6*8))))
				s7 = int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(0*8)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*(7*8))))
				in = (*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*1))
			}
			x0 = s0 + s3
			x1 = s1 + s2
			x2 = s1 - s2
			x3 = s0 - s3
			t0 = (x0 + x1) * int(cospi_16_64)
			t1 = (x0 - x1) * int(cospi_16_64)
			t2 = x2*int(cospi_24_64) + x3*int(cospi_8_64)
			t3 = -x2*int(cospi_8_64) + x3*int(cospi_24_64)
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*0)) = int16(fdct_round_shift(t0))
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*2)) = int16(fdct_round_shift(t2))
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*4)) = int16(fdct_round_shift(t1))
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*6)) = int16(fdct_round_shift(t3))
			t0 = (s6 - s5) * int(cospi_16_64)
			t1 = (s6 + s5) * int(cospi_16_64)
			t2 = fdct_round_shift(t0)
			t3 = fdct_round_shift(t1)
			x0 = s4 + t2
			x1 = s4 - t2
			x2 = s7 - t3
			x3 = s7 + t3
			t0 = x0*int(cospi_28_64) + x3*int(cospi_4_64)
			t1 = x1*int(cospi_12_64) + x2*int(cospi_20_64)
			t2 = x2*int(cospi_12_64) + x1*int(-cospi_20_64)
			t3 = x3*int(cospi_28_64) + x0*int(-cospi_4_64)
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*1)) = int16(fdct_round_shift(t0))
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*3)) = int16(fdct_round_shift(t2))
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*5)) = int16(fdct_round_shift(t1))
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*7)) = int16(fdct_round_shift(t3))
			out = (*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*8))
		}
		in = &intermediate[0]
		out = output
	}
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*uintptr(j+i*8))) /= 2
		}
	}
}
func vpx_fdct8x8_1_c(input *int16, output *int16, stride int) {
	var (
		r   int
		c   int
		sum int16 = 0
	)
	for r = 0; r < 8; r++ {
		for c = 0; c < 8; c++ {
			sum += *(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(r*stride+c)))
		}
	}
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = sum
}
func vpx_fdct16x16_c(input *int16, output *int16, stride int) {
	var (
		pass         int
		intermediate [256]int16
		in_low       *int16 = nil
		out          *int16 = &intermediate[0]
	)
	for pass = 0; pass < 2; pass++ {
		var (
			step1   [8]int
			step2   [8]int
			step3   [8]int
			in_high [8]int
			temp1   int
			temp2   int
			i       int
		)
		for i = 0; i < 16; i++ {
			if pass == 0 {
				in_high[0] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*0)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*15))))) * 4
				in_high[1] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*1)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*14))))) * 4
				in_high[2] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*2)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*13))))) * 4
				in_high[3] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*3)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*12))))) * 4
				in_high[4] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*4)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*11))))) * 4
				in_high[5] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*5)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*10))))) * 4
				in_high[6] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*6)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*9))))) * 4
				in_high[7] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*7)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*8))))) * 4
				step1[0] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*7)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*8))))) * 4
				step1[1] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*6)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*9))))) * 4
				step1[2] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*5)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*10))))) * 4
				step1[3] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*4)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*11))))) * 4
				step1[4] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*3)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*12))))) * 4
				step1[5] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*2)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*13))))) * 4
				step1[6] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*1)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*14))))) * 4
				step1[7] = (int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*0)))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(stride*15))))) * 4
			} else {
				if in_low != nil {
				} else {
					__assert_fail(libc.CString("in_low != NULL"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
				}
				in_high[0] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(0*16)))) + 1) >> 2) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(15*16)))) + 1) >> 2)
				in_high[1] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(1*16)))) + 1) >> 2) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(14*16)))) + 1) >> 2)
				in_high[2] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(2*16)))) + 1) >> 2) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(13*16)))) + 1) >> 2)
				in_high[3] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(3*16)))) + 1) >> 2) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(12*16)))) + 1) >> 2)
				in_high[4] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(4*16)))) + 1) >> 2) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(11*16)))) + 1) >> 2)
				in_high[5] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(5*16)))) + 1) >> 2) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(10*16)))) + 1) >> 2)
				in_high[6] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(6*16)))) + 1) >> 2) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(9*16)))) + 1) >> 2)
				in_high[7] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(7*16)))) + 1) >> 2) + ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(8*16)))) + 1) >> 2)
				step1[0] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(7*16)))) + 1) >> 2) - ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(8*16)))) + 1) >> 2)
				step1[1] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(6*16)))) + 1) >> 2) - ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(9*16)))) + 1) >> 2)
				step1[2] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(5*16)))) + 1) >> 2) - ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(10*16)))) + 1) >> 2)
				step1[3] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(4*16)))) + 1) >> 2) - ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(11*16)))) + 1) >> 2)
				step1[4] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(3*16)))) + 1) >> 2) - ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(12*16)))) + 1) >> 2)
				step1[5] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(2*16)))) + 1) >> 2) - ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(13*16)))) + 1) >> 2)
				step1[6] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(1*16)))) + 1) >> 2) - ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(14*16)))) + 1) >> 2)
				step1[7] = ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(0*16)))) + 1) >> 2) - ((int(*(*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*(15*16)))) + 1) >> 2)
				in_low = (*int16)(unsafe.Add(unsafe.Pointer(in_low), unsafe.Sizeof(int16(0))*1))
			}
			{
				var (
					s0 int
					s1 int
					s2 int
					s3 int
					s4 int
					s5 int
					s6 int
					s7 int
					t0 int
					t1 int
					t2 int
					t3 int
					x0 int
					x1 int
					x2 int
					x3 int
				)
				s0 = in_high[0] + in_high[7]
				s1 = in_high[1] + in_high[6]
				s2 = in_high[2] + in_high[5]
				s3 = in_high[3] + in_high[4]
				s4 = in_high[3] - in_high[4]
				s5 = in_high[2] - in_high[5]
				s6 = in_high[1] - in_high[6]
				s7 = in_high[0] - in_high[7]
				x0 = s0 + s3
				x1 = s1 + s2
				x2 = s1 - s2
				x3 = s0 - s3
				t0 = (x0 + x1) * int(cospi_16_64)
				t1 = (x0 - x1) * int(cospi_16_64)
				t2 = x3*int(cospi_8_64) + x2*int(cospi_24_64)
				t3 = x3*int(cospi_24_64) - x2*int(cospi_8_64)
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*0)) = int16(fdct_round_shift(t0))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*4)) = int16(fdct_round_shift(t2))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*8)) = int16(fdct_round_shift(t1))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*12)) = int16(fdct_round_shift(t3))
				t0 = (s6 - s5) * int(cospi_16_64)
				t1 = (s6 + s5) * int(cospi_16_64)
				t2 = fdct_round_shift(t0)
				t3 = fdct_round_shift(t1)
				x0 = s4 + t2
				x1 = s4 - t2
				x2 = s7 - t3
				x3 = s7 + t3
				t0 = x0*int(cospi_28_64) + x3*int(cospi_4_64)
				t1 = x1*int(cospi_12_64) + x2*int(cospi_20_64)
				t2 = x2*int(cospi_12_64) + x1*int(-cospi_20_64)
				t3 = x3*int(cospi_28_64) + x0*int(-cospi_4_64)
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*2)) = int16(fdct_round_shift(t0))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*6)) = int16(fdct_round_shift(t2))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*10)) = int16(fdct_round_shift(t1))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*14)) = int16(fdct_round_shift(t3))
			}
			{
				temp1 = (step1[5] - step1[2]) * int(cospi_16_64)
				temp2 = (step1[4] - step1[3]) * int(cospi_16_64)
				step2[2] = fdct_round_shift(temp1)
				step2[3] = fdct_round_shift(temp2)
				temp1 = (step1[4] + step1[3]) * int(cospi_16_64)
				temp2 = (step1[5] + step1[2]) * int(cospi_16_64)
				step2[4] = fdct_round_shift(temp1)
				step2[5] = fdct_round_shift(temp2)
				step3[0] = step1[0] + step2[3]
				step3[1] = step1[1] + step2[2]
				step3[2] = step1[1] - step2[2]
				step3[3] = step1[0] - step2[3]
				step3[4] = step1[7] - step2[4]
				step3[5] = step1[6] - step2[5]
				step3[6] = step1[6] + step2[5]
				step3[7] = step1[7] + step2[4]
				temp1 = step3[1]*int(-cospi_8_64) + step3[6]*int(cospi_24_64)
				temp2 = step3[2]*int(cospi_24_64) + step3[5]*int(cospi_8_64)
				step2[1] = fdct_round_shift(temp1)
				step2[2] = fdct_round_shift(temp2)
				temp1 = step3[2]*int(cospi_8_64) - step3[5]*int(cospi_24_64)
				temp2 = step3[1]*int(cospi_24_64) + step3[6]*int(cospi_8_64)
				step2[5] = fdct_round_shift(temp1)
				step2[6] = fdct_round_shift(temp2)
				step1[0] = step3[0] + step2[1]
				step1[1] = step3[0] - step2[1]
				step1[2] = step3[3] + step2[2]
				step1[3] = step3[3] - step2[2]
				step1[4] = step3[4] - step2[5]
				step1[5] = step3[4] + step2[5]
				step1[6] = step3[7] - step2[6]
				step1[7] = step3[7] + step2[6]
				temp1 = step1[0]*int(cospi_30_64) + step1[7]*int(cospi_2_64)
				temp2 = step1[1]*int(cospi_14_64) + step1[6]*int(cospi_18_64)
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*1)) = int16(fdct_round_shift(temp1))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*9)) = int16(fdct_round_shift(temp2))
				temp1 = step1[2]*int(cospi_22_64) + step1[5]*int(cospi_10_64)
				temp2 = step1[3]*int(cospi_6_64) + step1[4]*int(cospi_26_64)
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*5)) = int16(fdct_round_shift(temp1))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*13)) = int16(fdct_round_shift(temp2))
				temp1 = step1[3]*int(-cospi_26_64) + step1[4]*int(cospi_6_64)
				temp2 = step1[2]*int(-cospi_10_64) + step1[5]*int(cospi_22_64)
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*3)) = int16(fdct_round_shift(temp1))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*11)) = int16(fdct_round_shift(temp2))
				temp1 = step1[1]*int(-cospi_18_64) + step1[6]*int(cospi_14_64)
				temp2 = step1[0]*int(-cospi_2_64) + step1[7]*int(cospi_30_64)
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*7)) = int16(fdct_round_shift(temp1))
				*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*15)) = int16(fdct_round_shift(temp2))
			}
			input = (*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*1))
			out = (*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*16))
		}
		in_low = &intermediate[0]
		out = output
	}
}
func vpx_fdct16x16_1_c(input *int16, output *int16, stride int) {
	var (
		r   int
		c   int
		sum int = 0
	)
	for r = 0; r < 16; r++ {
		for c = 0; c < 16; c++ {
			sum += int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(r*stride+c))))
		}
	}
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(sum >> 1)
}
func dct_32_round(input int) int {
	var rv int = ((input + (1 << (int(DCT_CONST_BITS - 1)))) >> DCT_CONST_BITS)
	return rv
}
func half_round_shift(input int) int {
	var rv int = (input + 1 + int(libc.BoolToInt(input < 0))) >> 2
	return rv
}
func vpx_fdct32(input *int, output *int, round int) {
	var step [32]int
	step[0] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*0)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-1)))
	step[1] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*1)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-2)))
	step[2] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*2)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-3)))
	step[3] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*3)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-4)))
	step[4] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*4)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-5)))
	step[5] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*5)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-6)))
	step[6] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*6)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-7)))
	step[7] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*7)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-8)))
	step[8] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*8)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-9)))
	step[9] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*9)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-10)))
	step[10] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*10)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-11)))
	step[11] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*11)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-12)))
	step[12] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*12)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-13)))
	step[13] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*13)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-14)))
	step[14] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*14)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-15)))
	step[15] = *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*15)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-16)))
	step[16] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*16)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-17)))
	step[17] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*17)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-18)))
	step[18] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*18)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-19)))
	step[19] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*19)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-20)))
	step[20] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*20)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-21)))
	step[21] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*21)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-22)))
	step[22] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*22)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-23)))
	step[23] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*23)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-24)))
	step[24] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*24)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-25)))
	step[25] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*25)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-26)))
	step[26] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*26)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-27)))
	step[27] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*27)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-28)))
	step[28] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*28)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-29)))
	step[29] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*29)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-30)))
	step[30] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*30)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-31)))
	step[31] = -*(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*31)) + *(*int)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int(0))*(32-32)))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0)) = step[0] + step[16-1]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1)) = step[1] + step[16-2]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2)) = step[2] + step[16-3]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3)) = step[3] + step[16-4]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4)) = step[4] + step[16-5]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5)) = step[5] + step[16-6]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6)) = step[6] + step[16-7]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7)) = step[7] + step[16-8]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8)) = -step[8] + step[16-9]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9)) = -step[9] + step[16-10]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10)) = -step[10] + step[16-11]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11)) = -step[11] + step[16-12]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12)) = -step[12] + step[16-13]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13)) = -step[13] + step[16-14]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14)) = -step[14] + step[16-15]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15)) = -step[15] + step[16-16]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16)) = step[16]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17)) = step[17]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18)) = step[18]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19)) = step[19]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20)) = dct_32_round((-step[20] + step[27]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21)) = dct_32_round((-step[21] + step[26]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22)) = dct_32_round((-step[22] + step[25]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23)) = dct_32_round((-step[23] + step[24]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24)) = dct_32_round((step[24] + step[23]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25)) = dct_32_round((step[25] + step[22]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26)) = dct_32_round((step[26] + step[21]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27)) = dct_32_round((step[27] + step[20]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28)) = step[28]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29)) = step[29]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30)) = step[30]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31)) = step[31]
	if round != 0 {
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30)))
		*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31)) = half_round_shift(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31)))
	}
	step[0] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*(8-1)))
	step[1] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*(8-2)))
	step[2] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*(8-3)))
	step[3] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*(8-4)))
	step[4] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*(8-5)))
	step[5] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*(8-6)))
	step[6] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*(8-7)))
	step[7] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*(8-8)))
	step[8] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8))
	step[9] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9))
	step[10] = dct_32_round((-*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13))) * int(cospi_16_64))
	step[11] = dct_32_round((-*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12))) * int(cospi_16_64))
	step[12] = dct_32_round((*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11))) * int(cospi_16_64))
	step[13] = dct_32_round((*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10))) * int(cospi_16_64))
	step[14] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14))
	step[15] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15))
	step[16] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23))
	step[17] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22))
	step[18] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21))
	step[19] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20))
	step[20] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19))
	step[21] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18))
	step[22] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17))
	step[23] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16))
	step[24] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31))
	step[25] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30))
	step[26] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29))
	step[27] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28))
	step[28] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27))
	step[29] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26))
	step[30] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25))
	step[31] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0)) = step[0] + step[3]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1)) = step[1] + step[2]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2)) = -step[2] + step[1]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3)) = -step[3] + step[0]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4)) = step[4]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5)) = dct_32_round((-step[5] + step[6]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6)) = dct_32_round((step[6] + step[5]) * int(cospi_16_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7)) = step[7]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8)) = step[8] + step[11]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9)) = step[9] + step[10]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10)) = -step[10] + step[9]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11)) = -step[11] + step[8]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12)) = -step[12] + step[15]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13)) = -step[13] + step[14]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14)) = step[14] + step[13]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15)) = step[15] + step[12]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16)) = step[16]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17)) = step[17]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18)) = dct_32_round(step[18]*int(-cospi_8_64) + step[29]*int(cospi_24_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19)) = dct_32_round(step[19]*int(-cospi_8_64) + step[28]*int(cospi_24_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20)) = dct_32_round(step[20]*int(-cospi_24_64) + step[27]*int(-cospi_8_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21)) = dct_32_round(step[21]*int(-cospi_24_64) + step[26]*int(-cospi_8_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22)) = step[22]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23)) = step[23]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24)) = step[24]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25)) = step[25]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26)) = dct_32_round(step[26]*int(cospi_24_64) + step[21]*int(-cospi_8_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27)) = dct_32_round(step[27]*int(cospi_24_64) + step[20]*int(-cospi_8_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28)) = dct_32_round(step[28]*int(cospi_8_64) + step[19]*int(cospi_24_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29)) = dct_32_round(step[29]*int(cospi_8_64) + step[18]*int(cospi_24_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30)) = step[30]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31)) = step[31]
	step[0] = dct_32_round((*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1))) * int(cospi_16_64))
	step[1] = dct_32_round((-*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0))) * int(cospi_16_64))
	step[2] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2))*int(cospi_24_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3))*int(cospi_8_64))
	step[3] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3))*int(cospi_24_64) - *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2))*int(cospi_8_64))
	step[4] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5))
	step[5] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4))
	step[6] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7))
	step[7] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6))
	step[8] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8))
	step[9] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9))*int(-cospi_8_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14))*int(cospi_24_64))
	step[10] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10))*int(-cospi_24_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13))*int(-cospi_8_64))
	step[11] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11))
	step[12] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12))
	step[13] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13))*int(cospi_24_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10))*int(-cospi_8_64))
	step[14] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14))*int(cospi_8_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9))*int(cospi_24_64))
	step[15] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15))
	step[16] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19))
	step[17] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18))
	step[18] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17))
	step[19] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16))
	step[20] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23))
	step[21] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22))
	step[22] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21))
	step[23] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20))
	step[24] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27))
	step[25] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26))
	step[26] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25))
	step[27] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24))
	step[28] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31))
	step[29] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30))
	step[30] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29))
	step[31] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0)) = step[0]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1)) = step[1]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2)) = step[2]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3)) = step[3]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4)) = dct_32_round(step[4]*int(cospi_28_64) + step[7]*int(cospi_4_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5)) = dct_32_round(step[5]*int(cospi_12_64) + step[6]*int(cospi_20_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6)) = dct_32_round(step[6]*int(cospi_12_64) + step[5]*int(-cospi_20_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7)) = dct_32_round(step[7]*int(cospi_28_64) + step[4]*int(-cospi_4_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8)) = step[8] + step[9]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9)) = -step[9] + step[8]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10)) = -step[10] + step[11]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11)) = step[11] + step[10]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12)) = step[12] + step[13]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13)) = -step[13] + step[12]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14)) = -step[14] + step[15]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15)) = step[15] + step[14]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16)) = step[16]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17)) = dct_32_round(step[17]*int(-cospi_4_64) + step[30]*int(cospi_28_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18)) = dct_32_round(step[18]*int(-cospi_28_64) + step[29]*int(-cospi_4_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19)) = step[19]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20)) = step[20]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21)) = dct_32_round(step[21]*int(-cospi_20_64) + step[26]*int(cospi_12_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22)) = dct_32_round(step[22]*int(-cospi_12_64) + step[25]*int(-cospi_20_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23)) = step[23]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24)) = step[24]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25)) = dct_32_round(step[25]*int(cospi_12_64) + step[22]*int(-cospi_20_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26)) = dct_32_round(step[26]*int(cospi_20_64) + step[21]*int(cospi_12_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27)) = step[27]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28)) = step[28]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29)) = dct_32_round(step[29]*int(cospi_28_64) + step[18]*int(-cospi_4_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30)) = dct_32_round(step[30]*int(cospi_4_64) + step[17]*int(cospi_28_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31)) = step[31]
	step[0] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0))
	step[1] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1))
	step[2] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2))
	step[3] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3))
	step[4] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4))
	step[5] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5))
	step[6] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6))
	step[7] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7))
	step[8] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8))*int(cospi_30_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15))*int(cospi_2_64))
	step[9] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9))*int(cospi_14_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14))*int(cospi_18_64))
	step[10] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10))*int(cospi_22_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13))*int(cospi_10_64))
	step[11] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11))*int(cospi_6_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12))*int(cospi_26_64))
	step[12] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12))*int(cospi_6_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11))*int(-cospi_26_64))
	step[13] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13))*int(cospi_22_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10))*int(-cospi_10_64))
	step[14] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14))*int(cospi_14_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9))*int(-cospi_18_64))
	step[15] = dct_32_round(*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15))*int(cospi_30_64) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8))*int(-cospi_2_64))
	step[16] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17))
	step[17] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16))
	step[18] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19))
	step[19] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18))
	step[20] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21))
	step[21] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20))
	step[22] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23))
	step[23] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22))
	step[24] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25))
	step[25] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24))
	step[26] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27))
	step[27] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26))
	step[28] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29))
	step[29] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28))
	step[30] = -*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31))
	step[31] = *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31)) + *(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*0)) = step[0]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*16)) = step[1]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*8)) = step[2]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*24)) = step[3]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*4)) = step[4]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*20)) = step[5]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*12)) = step[6]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*28)) = step[7]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*2)) = step[8]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*18)) = step[9]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*10)) = step[10]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*26)) = step[11]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*6)) = step[12]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*22)) = step[13]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*14)) = step[14]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*30)) = step[15]
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*1)) = dct_32_round(step[16]*int(cospi_31_64) + step[31]*int(cospi_1_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*17)) = dct_32_round(step[17]*int(cospi_15_64) + step[30]*int(cospi_17_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*9)) = dct_32_round(step[18]*int(cospi_23_64) + step[29]*int(cospi_9_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*25)) = dct_32_round(step[19]*int(cospi_7_64) + step[28]*int(cospi_25_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*5)) = dct_32_round(step[20]*int(cospi_27_64) + step[27]*int(cospi_5_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*21)) = dct_32_round(step[21]*int(cospi_11_64) + step[26]*int(cospi_21_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*13)) = dct_32_round(step[22]*int(cospi_19_64) + step[25]*int(cospi_13_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*29)) = dct_32_round(step[23]*int(cospi_3_64) + step[24]*int(cospi_29_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*3)) = dct_32_round(step[24]*int(cospi_3_64) + step[23]*int(-cospi_29_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*19)) = dct_32_round(step[25]*int(cospi_19_64) + step[22]*int(-cospi_13_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*11)) = dct_32_round(step[26]*int(cospi_11_64) + step[21]*int(-cospi_21_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*27)) = dct_32_round(step[27]*int(cospi_27_64) + step[20]*int(-cospi_5_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*7)) = dct_32_round(step[28]*int(cospi_7_64) + step[19]*int(-cospi_25_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*23)) = dct_32_round(step[29]*int(cospi_23_64) + step[18]*int(-cospi_9_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*15)) = dct_32_round(step[30]*int(cospi_15_64) + step[17]*int(-cospi_17_64))
	*(*int)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int(0))*31)) = dct_32_round(step[31]*int(cospi_31_64) + step[16]*int(-cospi_1_64))
}
func vpx_fdct32x32_c(input *int16, output *int16, stride int) {
	var (
		i   int
		j   int
		out [1024]int
	)
	for i = 0; i < 32; i++ {
		var (
			temp_in  [32]int
			temp_out [32]int
		)
		for j = 0; j < 32; j++ {
			temp_in[j] = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(j*stride+i)))) * 4
		}
		vpx_fdct32(&temp_in[0], &temp_out[0], 0)
		for j = 0; j < 32; j++ {
			out[j*32+i] = (temp_out[j] + 1 + int(libc.BoolToInt(temp_out[j] > 0))) >> 2
		}
	}
	for i = 0; i < 32; i++ {
		var (
			temp_in  [32]int
			temp_out [32]int
		)
		for j = 0; j < 32; j++ {
			temp_in[j] = out[j+i*32]
		}
		vpx_fdct32(&temp_in[0], &temp_out[0], 0)
		for j = 0; j < 32; j++ {
			*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*uintptr(j+i*32))) = int16((temp_out[j] + 1 + int(libc.BoolToInt(temp_out[j] < 0))) >> 2)
		}
	}
}
func vpx_fdct32x32_rd_c(input *int16, output *int16, stride int) {
	var (
		i   int
		j   int
		out [1024]int
	)
	for i = 0; i < 32; i++ {
		var (
			temp_in  [32]int
			temp_out [32]int
		)
		for j = 0; j < 32; j++ {
			temp_in[j] = int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(j*stride+i)))) * 4
		}
		vpx_fdct32(&temp_in[0], &temp_out[0], 0)
		for j = 0; j < 32; j++ {
			out[j*32+i] = (temp_out[j] + 1 + int(libc.BoolToInt(temp_out[j] > 0))) >> 2
		}
	}
	for i = 0; i < 32; i++ {
		var (
			temp_in  [32]int
			temp_out [32]int
		)
		for j = 0; j < 32; j++ {
			temp_in[j] = out[j+i*32]
		}
		vpx_fdct32(&temp_in[0], &temp_out[0], 1)
		for j = 0; j < 32; j++ {
			*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*uintptr(j+i*32))) = int16(temp_out[j])
		}
	}
}
func vpx_fdct32x32_1_c(input *int16, output *int16, stride int) {
	var (
		r   int
		c   int
		sum int = 0
	)
	for r = 0; r < 32; r++ {
		for c = 0; c < 32; c++ {
			sum += int(*(*int16)(unsafe.Add(unsafe.Pointer(input), unsafe.Sizeof(int16(0))*uintptr(r*stride+c))))
		}
	}
	*(*int16)(unsafe.Add(unsafe.Pointer(output), unsafe.Sizeof(int16(0))*0)) = int16(sum >> 3)
}
