package dsp

import "unsafe"

func vpx_subtract_block_c(rows int, cols int, diff_ptr *int16, diff_stride int64, src_ptr *uint8, src_stride int64, pred_ptr *uint8, pred_stride int64) {
	var (
		r int
		c int
	)
	for r = 0; r < rows; r++ {
		for c = 0; c < cols; c++ {
			*(*int16)(unsafe.Add(unsafe.Pointer(diff_ptr), unsafe.Sizeof(int16(0))*uintptr(c))) = int16(int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), c))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), c))))
		}
		diff_ptr = (*int16)(unsafe.Add(unsafe.Pointer(diff_ptr), unsafe.Sizeof(int16(0))*uintptr(diff_stride)))
		pred_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), pred_stride))
		src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), src_stride))
	}
}
