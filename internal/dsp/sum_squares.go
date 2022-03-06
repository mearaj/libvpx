package dsp

import "unsafe"

func vpx_sum_squares_2d_i16_c(src *int16, stride int, size int) uint64 {
	var (
		r  int
		c  int
		ss uint64 = 0
	)
	for r = 0; r < size; r++ {
		for c = 0; c < size; c++ {
			var v int16 = *(*int16)(unsafe.Add(unsafe.Pointer(src), unsafe.Sizeof(int16(0))*uintptr(c)))
			ss += uint64(int(v) * int(v))
		}
		src = (*int16)(unsafe.Add(unsafe.Pointer(src), unsafe.Sizeof(int16(0))*uintptr(stride)))
	}
	return ss
}
