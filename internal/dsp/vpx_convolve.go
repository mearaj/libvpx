package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

type convolve_fn_t func(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int)

func convolve_horiz(src *uint8, src_stride int64, dst *uint8, dst_stride int64, x_filters *InterpKernel, x0_q4 int, x_step_q4 int, w int, h int) {
	var (
		x int
		y int
	)
	src = (*uint8)(unsafe.Add(unsafe.Pointer(src), -(int(SUBPEL_TAPS/2) - 1)))
	for y = 0; y < h; y++ {
		var x_q4 int = x0_q4
		for x = 0; x < w; x++ {
			var (
				src_x    *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(src), x_q4>>SUBPEL_BITS))
				x_filter *int16 = (*int16)(unsafe.Pointer((*InterpKernel)(unsafe.Add(unsafe.Pointer(x_filters), unsafe.Sizeof(InterpKernel{})*uintptr(x_q4&((int(1<<SUBPEL_BITS))-1))))))
				k        int
				sum      int = 0
			)
			for k = 0; k < SUBPEL_TAPS; k++ {
				sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_x), k))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(x_filter), unsafe.Sizeof(int16(0))*uintptr(k))))
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), x)) = clip_pixel((sum + (1 << (int(FILTER_BITS - 1)))) >> FILTER_BITS)
			x_q4 += x_step_q4
		}
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_stride))
	}
}
func convolve_avg_horiz(src *uint8, src_stride int64, dst *uint8, dst_stride int64, x_filters *InterpKernel, x0_q4 int, x_step_q4 int, w int, h int) {
	var (
		x int
		y int
	)
	src = (*uint8)(unsafe.Add(unsafe.Pointer(src), -(int(SUBPEL_TAPS/2) - 1)))
	for y = 0; y < h; y++ {
		var x_q4 int = x0_q4
		for x = 0; x < w; x++ {
			var (
				src_x    *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(src), x_q4>>SUBPEL_BITS))
				x_filter *int16 = (*int16)(unsafe.Pointer((*InterpKernel)(unsafe.Add(unsafe.Pointer(x_filters), unsafe.Sizeof(InterpKernel{})*uintptr(x_q4&((int(1<<SUBPEL_BITS))-1))))))
				k        int
				sum      int = 0
			)
			for k = 0; k < SUBPEL_TAPS; k++ {
				sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_x), k))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(x_filter), unsafe.Sizeof(int16(0))*uintptr(k))))
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), x)) = uint8(int8(((int(*(*uint8)(unsafe.Add(unsafe.Pointer(dst), x))) + int(clip_pixel((sum+(1<<(int(FILTER_BITS-1))))>>FILTER_BITS))) + (1 << (1 - 1))) >> 1))
			x_q4 += x_step_q4
		}
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_stride))
	}
}
func convolve_vert(src *uint8, src_stride int64, dst *uint8, dst_stride int64, y_filters *InterpKernel, y0_q4 int, y_step_q4 int, w int, h int) {
	var (
		x int
		y int
	)
	src = (*uint8)(unsafe.Add(unsafe.Pointer(src), -(src_stride * int64(int(SUBPEL_TAPS/2)-1))))
	for x = 0; x < w; x++ {
		var y_q4 int = y0_q4
		for y = 0; y < h; y++ {
			var (
				src_y    *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(src), (y_q4>>SUBPEL_BITS)*int(src_stride)))
				y_filter *int16 = (*int16)(unsafe.Pointer((*InterpKernel)(unsafe.Add(unsafe.Pointer(y_filters), unsafe.Sizeof(InterpKernel{})*uintptr(y_q4&((int(1<<SUBPEL_BITS))-1))))))
				k        int
				sum      int = 0
			)
			for k = 0; k < SUBPEL_TAPS; k++ {
				sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_y), k*int(src_stride)))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(y_filter), unsafe.Sizeof(int16(0))*uintptr(k))))
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), y*int(dst_stride))) = clip_pixel((sum + (1 << (int(FILTER_BITS - 1)))) >> FILTER_BITS)
			y_q4 += y_step_q4
		}
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), 1))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), 1))
	}
}
func convolve_avg_vert(src *uint8, src_stride int64, dst *uint8, dst_stride int64, y_filters *InterpKernel, y0_q4 int, y_step_q4 int, w int, h int) {
	var (
		x int
		y int
	)
	src = (*uint8)(unsafe.Add(unsafe.Pointer(src), -(src_stride * int64(int(SUBPEL_TAPS/2)-1))))
	for x = 0; x < w; x++ {
		var y_q4 int = y0_q4
		for y = 0; y < h; y++ {
			var (
				src_y    *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(src), (y_q4>>SUBPEL_BITS)*int(src_stride)))
				y_filter *int16 = (*int16)(unsafe.Pointer((*InterpKernel)(unsafe.Add(unsafe.Pointer(y_filters), unsafe.Sizeof(InterpKernel{})*uintptr(y_q4&((int(1<<SUBPEL_BITS))-1))))))
				k        int
				sum      int = 0
			)
			for k = 0; k < SUBPEL_TAPS; k++ {
				sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_y), k*int(src_stride)))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(y_filter), unsafe.Sizeof(int16(0))*uintptr(k))))
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), y*int(dst_stride))) = uint8(int8(((int(*(*uint8)(unsafe.Add(unsafe.Pointer(dst), y*int(dst_stride)))) + int(clip_pixel((sum+(1<<(int(FILTER_BITS-1))))>>FILTER_BITS))) + (1 << (1 - 1))) >> 1))
			y_q4 += y_step_q4
		}
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), 1))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), 1))
	}
}
func vpx_convolve8_horiz_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	_ = y0_q4
	_ = y_step_q4
	convolve_horiz(src, src_stride, dst, dst_stride, filter, x0_q4, x_step_q4, w, h)
}
func vpx_convolve8_avg_horiz_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	_ = y0_q4
	_ = y_step_q4
	convolve_avg_horiz(src, src_stride, dst, dst_stride, filter, x0_q4, x_step_q4, w, h)
}
func vpx_convolve8_vert_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	_ = x0_q4
	_ = x_step_q4
	convolve_vert(src, src_stride, dst, dst_stride, filter, y0_q4, y_step_q4, w, h)
}
func vpx_convolve8_avg_vert_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	_ = x0_q4
	_ = x_step_q4
	convolve_avg_vert(src, src_stride, dst, dst_stride, filter, y0_q4, y_step_q4, w, h)
}
func vpx_convolve8_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	var (
		temp                [8640]uint8
		intermediate_height int = (((h-1)*y_step_q4 + y0_q4) >> SUBPEL_BITS) + SUBPEL_TAPS
	)
	libc.Assert(w <= 64)
	libc.Assert(h <= 64)
	libc.Assert(y_step_q4 <= 32 || y_step_q4 <= 64 && h <= 32)
	libc.Assert(x_step_q4 <= 64)
	convolve_horiz((*uint8)(unsafe.Add(unsafe.Pointer(src), -(src_stride*int64(int(SUBPEL_TAPS/2)-1)))), src_stride, &temp[0], 64, filter, x0_q4, x_step_q4, w, intermediate_height)
	convolve_vert(&temp[(int(SUBPEL_TAPS/2)-1)*64], 64, dst, dst_stride, filter, y0_q4, y_step_q4, w, h)
}
func vpx_convolve8_avg_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	var temp [4096]uint8
	libc.Assert(w <= 64)
	libc.Assert(h <= 64)
	vpx_convolve8_c(src, src_stride, &temp[0], 64, filter, x0_q4, x_step_q4, y0_q4, y_step_q4, w, h)
	vpx_convolve_avg_c(&temp[0], 64, dst, dst_stride, nil, 0, 0, 0, 0, w, h)
}
func vpx_convolve_copy_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	var r int
	_ = filter
	_ = x0_q4
	_ = x_step_q4
	_ = y0_q4
	_ = y_step_q4
	for r = h; r > 0; r-- {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), w)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_stride))
	}
}
func vpx_convolve_avg_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	var (
		x int
		y int
	)
	_ = filter
	_ = x0_q4
	_ = x_step_q4
	_ = y0_q4
	_ = y_step_q4
	for y = 0; y < h; y++ {
		for x = 0; x < w; x++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), x)) = uint8(int8(((int(*(*uint8)(unsafe.Add(unsafe.Pointer(dst), x))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(src), x)))) + (1 << (1 - 1))) >> 1))
		}
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_stride))
	}
}
func vpx_scaled_horiz_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	vpx_convolve8_horiz_c(src, src_stride, dst, dst_stride, filter, x0_q4, x_step_q4, y0_q4, y_step_q4, w, h)
}
func vpx_scaled_vert_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	vpx_convolve8_vert_c(src, src_stride, dst, dst_stride, filter, x0_q4, x_step_q4, y0_q4, y_step_q4, w, h)
}
func vpx_scaled_2d_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	vpx_convolve8_c(src, src_stride, dst, dst_stride, filter, x0_q4, x_step_q4, y0_q4, y_step_q4, w, h)
}
func vpx_scaled_avg_horiz_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	vpx_convolve8_avg_horiz_c(src, src_stride, dst, dst_stride, filter, x0_q4, x_step_q4, y0_q4, y_step_q4, w, h)
}
func vpx_scaled_avg_vert_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	vpx_convolve8_avg_vert_c(src, src_stride, dst, dst_stride, filter, x0_q4, x_step_q4, y0_q4, y_step_q4, w, h)
}
func vpx_scaled_avg_2d_c(src *uint8, src_stride int64, dst *uint8, dst_stride int64, filter *InterpKernel, x0_q4 int, x_step_q4 int, y0_q4 int, y_step_q4 int, w int, h int) {
	vpx_convolve8_avg_c(src, src_stride, dst, dst_stride, filter, x0_q4, x_step_q4, y0_q4, y_step_q4, w, h)
}
