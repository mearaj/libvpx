package dsp

import "unsafe"

const FILTER_BITS = 7
const FILTER_WEIGHT = 128

type vpx_sad_fn_t func(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint
type vpx_sad_avg_fn_t func(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, second_pred *uint8) uint
type vp8_copy32xn_fn_t func(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, n int)
type vpx_sad_multi_fn_t func(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sad_array *uint)
type vpx_sad_multi_d_fn_t func(src_ptr *uint8, src_stride int, b_array [0]*uint8, ref_stride int, sad_array *uint)
type vpx_variance_fn_t func(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint) uint
type vpx_subpixvariance_fn_t func(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint) uint
type vpx_subp_avg_variance_fn_t func(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint, second_pred *uint8) uint
type variance_vtable struct {
	Sdf     vpx_sad_fn_t
	Vf      vpx_variance_fn_t
	Svf     vpx_subpixvariance_fn_t
	Sdx3f   vpx_sad_multi_fn_t
	Sdx8f   vpx_sad_multi_fn_t
	Sdx4df  vpx_sad_multi_d_fn_t
	Copymem vp8_copy32xn_fn_t
}
type vp8_variance_fn_ptr_t variance_vtable
type vp9_variance_vtable struct {
	Sdf    vpx_sad_fn_t
	Sdaf   vpx_sad_avg_fn_t
	Vf     vpx_variance_fn_t
	Svf    vpx_subpixvariance_fn_t
	Svaf   vpx_subp_avg_variance_fn_t
	Sdx4df vpx_sad_multi_d_fn_t
	Sdx8f  vpx_sad_multi_fn_t
}
type vp9_variance_fn_ptr_t vp9_variance_vtable

var bilinear_filters [8][2]uint8 = [8][2]uint8{{128, 0}, {112, 16}, {96, 32}, {80, 48}, {64, 64}, {48, 80}, {32, 96}, {16, 112}}

func vpx_get4x4sse_cs_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int) uint32 {
	var (
		distortion int = 0
		r          int
		c          int
	)
	for r = 0; r < 4; r++ {
		for c = 0; c < 4; c++ {
			var diff int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), c))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), c)))
			distortion += diff * diff
		}
		src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), src_stride))
		ref_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), ref_stride))
	}
	return uint32(int32(distortion))
}
func vpx_get_mb_ss_c(src_ptr *int16) uint32 {
	var (
		i   uint
		sum uint = 0
	)
	for i = 0; i < 256; i++ {
		sum += uint(int(*(*int16)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int16(0))*uintptr(i)))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int16(0))*uintptr(i)))))
	}
	return uint32(sum)
}
func variance(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, w int, h int, sse *uint32, sum *int) {
	var (
		i int
		j int
	)
	*sum = 0
	*sse = 0
	for i = 0; i < h; i++ {
		for j = 0; j < w; j++ {
			var diff int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), j))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), j)))
			*sum += diff
			*sse += uint32(int32(diff * diff))
		}
		src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), src_stride))
		ref_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), ref_stride))
	}
}
func var_filter_block2d_bil_first_pass(src_ptr *uint8, ref_ptr *uint16, src_pixels_per_line uint, pixel_step int, output_height uint, output_width uint, filter *uint8) {
	var (
		i uint
		j uint
	)
	for i = 0; i < output_height; i++ {
		for j = 0; j < output_width; j++ {
			*(*uint16)(unsafe.Add(unsafe.Pointer(ref_ptr), unsafe.Sizeof(uint16(0))*uintptr(j))) = uint16(int16(((int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), 0)))*int(*(*uint8)(unsafe.Add(unsafe.Pointer(filter), 0))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), pixel_step)))*int(*(*uint8)(unsafe.Add(unsafe.Pointer(filter), 1)))) + (1 << (int(FILTER_BITS - 1)))) >> FILTER_BITS))
			src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), 1))
		}
		src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), src_pixels_per_line-output_width))
		ref_ptr = (*uint16)(unsafe.Add(unsafe.Pointer(ref_ptr), unsafe.Sizeof(uint16(0))*uintptr(output_width)))
	}
}
func var_filter_block2d_bil_second_pass(src_ptr *uint16, ref_ptr *uint8, src_pixels_per_line uint, pixel_step uint, output_height uint, output_width uint, filter *uint8) {
	var (
		i uint
		j uint
	)
	for i = 0; i < output_height; i++ {
		for j = 0; j < output_width; j++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), j)) = uint8(int8(((int(*(*uint16)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(uint16(0))*0)))*int(*(*uint8)(unsafe.Add(unsafe.Pointer(filter), 0))) + int(*(*uint16)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(uint16(0))*uintptr(pixel_step))))*int(*(*uint8)(unsafe.Add(unsafe.Pointer(filter), 1)))) + (1 << (int(FILTER_BITS - 1)))) >> FILTER_BITS))
			src_ptr = (*uint16)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(uint16(0))*1))
		}
		src_ptr = (*uint16)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(uint16(0))*uintptr(src_pixels_per_line-output_width)))
		ref_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(ref_ptr), output_width))
	}
}
func vpx_variance64x64_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 64, 64, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(64*64))))))
}
func vpx_sub_pixel_variance64x64_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [4160]uint16
		temp2  [4096]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 64+1, 64, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 64, 64, 64, 64, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance64x64_c(&temp2[0], 64, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance64x64_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [4160]uint16
		temp2  [4096]uint8
		temp3  [4096]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 64+1, 64, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 64, 64, 64, 64, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 64, 64, &temp2[0], 64)
	return uint32(vpx_variance64x64_c(&temp3[0], 64, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_variance64x32_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 64, 32, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(64*32))))))
}
func vpx_sub_pixel_variance64x32_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [2112]uint16
		temp2  [2048]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 32+1, 64, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 64, 64, 32, 64, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance64x32_c(&temp2[0], 64, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance64x32_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [2112]uint16
		temp2  [2048]uint8
		temp3  [2048]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 32+1, 64, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 64, 64, 32, 64, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 64, 32, &temp2[0], 64)
	return uint32(vpx_variance64x32_c(&temp3[0], 64, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_variance32x64_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 32, 64, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(32*64))))))
}
func vpx_sub_pixel_variance32x64_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [2080]uint16
		temp2  [2048]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 64+1, 32, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 32, 32, 64, 32, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance32x64_c(&temp2[0], 32, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance32x64_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [2080]uint16
		temp2  [2048]uint8
		temp3  [2048]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 64+1, 32, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 32, 32, 64, 32, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 32, 64, &temp2[0], 32)
	return uint32(vpx_variance32x64_c(&temp3[0], 32, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_variance32x32_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 32, 32, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(32*32))))))
}
func vpx_sub_pixel_variance32x32_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [1056]uint16
		temp2  [1024]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 32+1, 32, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 32, 32, 32, 32, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance32x32_c(&temp2[0], 32, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance32x32_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [1056]uint16
		temp2  [1024]uint8
		temp3  [1024]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 32+1, 32, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 32, 32, 32, 32, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 32, 32, &temp2[0], 32)
	return uint32(vpx_variance32x32_c(&temp3[0], 32, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_variance32x16_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 32, 16, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(32*16))))))
}
func vpx_sub_pixel_variance32x16_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [544]uint16
		temp2  [512]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 16+1, 32, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 32, 32, 16, 32, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance32x16_c(&temp2[0], 32, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance32x16_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [544]uint16
		temp2  [512]uint8
		temp3  [512]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 16+1, 32, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 32, 32, 16, 32, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 32, 16, &temp2[0], 32)
	return uint32(vpx_variance32x16_c(&temp3[0], 32, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_variance16x32_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 16, 32, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(16*32))))))
}
func vpx_sub_pixel_variance16x32_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [528]uint16
		temp2  [512]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 32+1, 16, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 16, 16, 32, 16, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance16x32_c(&temp2[0], 16, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance16x32_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [528]uint16
		temp2  [512]uint8
		temp3  [512]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 32+1, 16, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 16, 16, 32, 16, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 16, 32, &temp2[0], 16)
	return uint32(vpx_variance16x32_c(&temp3[0], 16, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func VpxVariance16x16C(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 16, 16, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(16*16))))))
}
func vpx_sub_pixel_variance16x16_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [272]uint16
		temp2  [256]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 16+1, 16, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 16, 16, 16, 16, &bilinear_filters[y_offset][0])
	return uint32(VpxVariance16x16C(&temp2[0], 16, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance16x16_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [272]uint16
		temp2  [256]uint8
		temp3  [256]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 16+1, 16, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 16, 16, 16, 16, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 16, 16, &temp2[0], 16)
	return uint32(VpxVariance16x16C(&temp3[0], 16, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_variance16x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 16, 8, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(16*8))))))
}
func vpx_sub_pixel_variance16x8_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [144]uint16
		temp2  [128]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 8+1, 16, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 16, 16, 8, 16, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance16x8_c(&temp2[0], 16, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance16x8_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [144]uint16
		temp2  [128]uint8
		temp3  [128]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 8+1, 16, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 16, 16, 8, 16, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 16, 8, &temp2[0], 16)
	return uint32(vpx_variance16x8_c(&temp3[0], 16, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_variance8x16_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 8, 16, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(8*16))))))
}
func vpx_sub_pixel_variance8x16_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [136]uint16
		temp2  [128]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 16+1, 8, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 8, 8, 16, 8, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance8x16_c(&temp2[0], 8, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance8x16_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [136]uint16
		temp2  [128]uint8
		temp3  [128]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 16+1, 8, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 8, 8, 16, 8, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 8, 16, &temp2[0], 8)
	return uint32(vpx_variance8x16_c(&temp3[0], 8, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func VpxVariance8x8C(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 8, 8, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(8*8))))))
}
func vpx_sub_pixel_variance8x8_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [72]uint16
		temp2  [64]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 8+1, 8, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 8, 8, 8, 8, &bilinear_filters[y_offset][0])
	return uint32(VpxVariance8x8C(&temp2[0], 8, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance8x8_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [72]uint16
		temp2  [64]uint8
		temp3  [64]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 8+1, 8, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 8, 8, 8, 8, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 8, 8, &temp2[0], 8)
	return uint32(VpxVariance8x8C(&temp3[0], 8, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_variance8x4_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 8, 4, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(8*4))))))
}
func vpx_sub_pixel_variance8x4_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [40]uint16
		temp2  [32]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 4+1, 8, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 8, 8, 4, 8, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance8x4_c(&temp2[0], 8, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance8x4_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [40]uint16
		temp2  [32]uint8
		temp3  [32]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 4+1, 8, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 8, 8, 4, 8, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 8, 4, &temp2[0], 8)
	return uint32(vpx_variance8x4_c(&temp3[0], 8, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_variance4x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 4, 8, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(4*8))))))
}
func vpx_sub_pixel_variance4x8_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [36]uint16
		temp2  [32]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 8+1, 4, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 4, 4, 8, 4, &bilinear_filters[y_offset][0])
	return uint32(vpx_variance4x8_c(&temp2[0], 4, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance4x8_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [36]uint16
		temp2  [32]uint8
		temp3  [32]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 8+1, 4, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 4, 4, 8, 4, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 4, 8, &temp2[0], 4)
	return uint32(vpx_variance4x8_c(&temp3[0], 4, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func VpxVariance4x4C(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 4, 4, sse, &sum)
	return uint32(int32(int(*sse) - int(uint32(int32((int64(sum)*int64(sum))/(4*4))))))
}
func vpx_sub_pixel_variance4x4_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var (
		fdata3 [20]uint16
		temp2  [16]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 4+1, 4, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 4, 4, 4, 4, &bilinear_filters[y_offset][0])
	return uint32(VpxVariance4x4C(&temp2[0], 4, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_sub_pixel_avg_variance4x4_c(src_ptr *uint8, src_stride int, x_offset int, y_offset int, ref_ptr *uint8, ref_stride int, sse *uint32, second_pred *uint8) uint32 {
	var (
		fdata3 [20]uint16
		temp2  [16]uint8
		temp3  [16]uint8
	)
	var_filter_block2d_bil_first_pass(src_ptr, &fdata3[0], uint(src_stride), 1, 4+1, 4, &bilinear_filters[x_offset][0])
	var_filter_block2d_bil_second_pass(&fdata3[0], &temp2[0], 4, 4, 4, 4, &bilinear_filters[y_offset][0])
	vpx_comp_avg_pred_c(&temp3[0], second_pred, 4, 4, &temp2[0], 4)
	return uint32(VpxVariance4x4C(&temp3[0], 4, ref_ptr, ref_stride, (*uint)(unsafe.Pointer(sse))))
}
func vpx_get16x16var_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32, sum *int) {
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 16, 16, sse, sum)
}
func vpx_get8x8var_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32, sum *int) {
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 8, 8, sse, sum)
}
func vpx_mse16x16_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 16, 16, sse, &sum)
	return *sse
}
func vpx_mse16x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 16, 8, sse, &sum)
	return *sse
}
func vpx_mse8x16_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 8, 16, sse, &sum)
	return *sse
}
func vpx_mse8x8_c(src_ptr *uint8, src_stride int, ref_ptr *uint8, ref_stride int, sse *uint32) uint32 {
	var sum int
	variance(src_ptr, src_stride, ref_ptr, ref_stride, 8, 8, sse, &sum)
	return *sse
}
func vpx_comp_avg_pred_c(comp_pred *uint8, pred *uint8, width int, height int, ref *uint8, ref_stride int) {
	var (
		i int
		j int
	)
	for i = 0; i < height; i++ {
		for j = 0; j < width; j++ {
			var tmp int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(pred), j))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(ref), j)))
			*(*uint8)(unsafe.Add(unsafe.Pointer(comp_pred), j)) = uint8(int8((tmp + (1 << (1 - 1))) >> 1))
		}
		comp_pred = (*uint8)(unsafe.Add(unsafe.Pointer(comp_pred), width))
		pred = (*uint8)(unsafe.Add(unsafe.Pointer(pred), width))
		ref = (*uint8)(unsafe.Add(unsafe.Pointer(ref), ref_stride))
	}
}
