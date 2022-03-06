package dsp

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"unsafe"
)

func sad(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, width int, height int) uint {
	var (
		y   int
		x   int
		sad uint = 0
	)
	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			sad += uint(cmath.Abs(int64(int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), x))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), x))))))
		}
		src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), src_stride))
		ref_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), ref_stride))
	}
	return sad
}
func vpx_sad64x64_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 64, 64)
}
func vpx_sad64x64_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [4096]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 64, 64, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 64, 64, 64)
}
func vpx_sad64x64x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad64x64_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad64x32_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 64, 32)
}
func vpx_sad64x32_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [2048]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 64, 32, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 64, 64, 32)
}
func vpx_sad64x32x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad64x32_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad32x64_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 32, 64)
}
func vpx_sad32x64_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [2048]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 32, 64, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 32, 32, 64)
}
func vpx_sad32x64x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad32x64_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad32x32_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 32, 32)
}
func vpx_sad32x32_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [1024]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 32, 32, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 32, 32, 32)
}
func vpx_sad32x32x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 8; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad32x32_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad32x32x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad32x32_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad32x16_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 32, 16)
}
func vpx_sad32x16_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [512]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 32, 16, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 32, 32, 16)
}
func vpx_sad32x16x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad32x16_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad16x32_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 16, 32)
}
func vpx_sad16x32_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [512]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 16, 32, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 16, 16, 32)
}
func vpx_sad16x32x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad16x32_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad16x16_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 16, 16)
}
func vpx_sad16x16_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [256]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 16, 16, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 16, 16, 16)
}
func vpx_sad16x16x3_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 3; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad16x16_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad16x16x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 8; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad16x16_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad16x16x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad16x16_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad16x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 16, 8)
}
func vpx_sad16x8_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [128]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 16, 8, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 16, 16, 8)
}
func vpx_sad16x8x3_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 3; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad16x8_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad16x8x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 8; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad16x8_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad16x8x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad16x8_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad8x16_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 8, 16)
}
func vpx_sad8x16_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [128]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 8, 16, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 8, 8, 16)
}
func vpx_sad8x16x3_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 3; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad8x16_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad8x16x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 8; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad8x16_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad8x16x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad8x16_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad8x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 8, 8)
}
func vpx_sad8x8_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [64]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 8, 8, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 8, 8, 8)
}
func vpx_sad8x8x3_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 3; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad8x8_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad8x8x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 8; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad8x8_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad8x8x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad8x8_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad8x4_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 8, 4)
}
func vpx_sad8x4_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [32]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 8, 4, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 8, 8, 4)
}
func vpx_sad8x4x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad8x4_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad4x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 4, 8)
}
func vpx_sad4x8_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [32]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 4, 8, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 4, 4, 8)
}
func vpx_sad4x8x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad4x8_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
func vpx_sad4x4_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint {
	return sad(src_ptr, src_stride, ref_ptr, ref_stride, 4, 4)
}
func vpx_sad4x4_avg_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint {
	var comp_pred [16]uint8
	vpx_comp_avg_pred_c(&comp_pred[0], second_pred, 4, 4, ref_ptr, ref_stride)
	return sad(src_ptr, src_stride, &comp_pred[0], 4, 4, 4)
}
func vpx_sad4x4x3_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 3; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad4x4_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad4x4x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 8; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad4x4_c(src_ptr, src_stride, (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), i)), ref_stride))
	}
}
func vpx_sad4x4x4d_c(src_ptr *uint8, src_stride int, ref_array [0]*uint8, ref_stride int, sad_array *uint32) {
	var i int
	for i = 0; i < 4; i++ {
		*(*uint32)(unsafe.Add(unsafe.Pointer(sad_array), unsafe.Sizeof(uint32(0))*uintptr(i))) = uint32(vpx_sad4x4_c(src_ptr, src_stride, ref_array[i], ref_stride))
	}
}
