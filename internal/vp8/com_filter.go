package vp8

import (
	"log"
	"math"
	"unsafe"
)

const BLOCK_HEIGHT_WIDTH = 4
const VP8_FILTER_WEIGHT = 128
const VP8_FILTER_SHIFT = 7

var vp8_bilinear_filters = [8][2]int16{{128, 0}, {112, 16}, {96, 32}, {80, 48}, {64, 64}, {48, 80}, {32, 96}, {16, 112}}
var vp8_sub_pel_filters = [8][6]int16{{0, 0, 128, 0, 0, 0}, {0, -6, 123, 12, -1, 0}, {2, -11, 108, 36, -8, 1}, {0, -9, 93, 50, -6, 0}, {3, -16, 77, 77, -16, 3}, {0, -6, 50, 93, -9, 0}, {1, -8, 36, 108, -11, 2}, {0, -1, 12, 123, -6, 0}}

func filter_block2d_first_pass(src_ptr *uint8, output_ptr *int, src_pixels_per_line uint, pixel_step uint, output_height uint, output_width uint, vp8_filter *int16) {
	var (
		i    uint
		j    uint
		Temp int
	)
	for i = 0; i < output_height; i++ {
		for j = 0; j < output_width; j++ {
			Temp = (int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), int(pixel_step)*(-2)))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*0)))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), int(pixel_step)*(-1))))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), 0)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*2))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), pixel_step)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*3))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), pixel_step*2)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*4))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), pixel_step*3)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*5))) + (int(VP8_FILTER_WEIGHT >> 1))
			Temp = Temp >> VP8_FILTER_SHIFT
			if Temp < 0 {
				Temp = 0
			} else if Temp > math.MaxUint8 {
				Temp = math.MaxUint8
			}
			*(*int)(unsafe.Add(unsafe.Pointer(output_ptr), unsafe.Sizeof(int(0))*uintptr(j))) = Temp
			src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), 1))
		}
		src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), src_pixels_per_line-output_width))
		output_ptr = (*int)(unsafe.Add(unsafe.Pointer(output_ptr), unsafe.Sizeof(int(0))*uintptr(output_width)))
	}
}
func filter_block2d_second_pass(src_ptr *int, output_ptr *uint8, output_pitch int, src_pixels_per_line uint, pixel_step uint, output_height uint, output_width uint, vp8_filter *int16) {
	var (
		i    uint
		j    uint
		Temp int
	)
	for i = 0; i < output_height; i++ {
		for j = 0; j < output_width; j++ {
			Temp = (*(*int)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int(0))*uintptr(int(pixel_step)*(-2)))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*0)))) + *(*int)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int(0))*uintptr(int(pixel_step)*(-1))))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*1))) + *(*int)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int(0))*0))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*2))) + *(*int)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int(0))*uintptr(pixel_step)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*3))) + *(*int)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int(0))*uintptr(pixel_step*2)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*4))) + *(*int)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int(0))*uintptr(pixel_step*3)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*5))) + (int(VP8_FILTER_WEIGHT >> 1))
			Temp = Temp >> VP8_FILTER_SHIFT
			if Temp < 0 {
				Temp = 0
			} else if Temp > math.MaxUint8 {
				Temp = math.MaxUint8
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(output_ptr), j)) = uint8(int8(Temp))
			src_ptr = (*int)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int(0))*1))
		}
		src_ptr = (*int)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(int(0))*uintptr(src_pixels_per_line-output_width)))
		output_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(output_ptr), output_pitch))
	}
}
func filter_block2d(src_ptr *uint8, output_ptr *uint8, src_pixels_per_line uint, output_pitch int, HFilter *int16, VFilter *int16) {
	var FData [36]int
	filter_block2d_first_pass((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), -int(src_pixels_per_line*2))), &FData[0], src_pixels_per_line, 1, 9, 4, HFilter)
	filter_block2d_second_pass(&FData[8], output_ptr, output_pitch, 4, 4, 4, 4, VFilter)
}
func Vp8SixtapPredict4x4C(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
	var (
		HFilter *int16
		VFilter *int16
	)
	HFilter = &vp8_sub_pel_filters[xoffset][0]
	VFilter = &vp8_sub_pel_filters[yoffset][0]
	filter_block2d(src_ptr, dst_ptr, uint(src_pixels_per_line), dst_pitch, HFilter, VFilter)
}
func Vp8SixtapPredict8x8C(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
	var (
		HFilter *int16
		VFilter *int16
		FData   [208]int
	)
	HFilter = &vp8_sub_pel_filters[xoffset][0]
	VFilter = &vp8_sub_pel_filters[yoffset][0]
	filter_block2d_first_pass((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), -(src_pixels_per_line*2))), &FData[0], uint(src_pixels_per_line), 1, 13, 8, HFilter)
	filter_block2d_second_pass(&FData[16], dst_ptr, dst_pitch, 8, 8, 8, 8, VFilter)
}
func Vp8SixtapPredict8x4C(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
	var (
		HFilter *int16
		VFilter *int16
		FData   [208]int
	)
	HFilter = &vp8_sub_pel_filters[xoffset][0]
	VFilter = &vp8_sub_pel_filters[yoffset][0]
	filter_block2d_first_pass((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), -(src_pixels_per_line*2))), &FData[0], uint(src_pixels_per_line), 1, 9, 8, HFilter)
	filter_block2d_second_pass(&FData[16], dst_ptr, dst_pitch, 8, 8, 4, 8, VFilter)
}
func Vp8SixtapPredict16x16C(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
	var (
		HFilter *int16
		VFilter *int16
		FData   [504]int
	)
	HFilter = &vp8_sub_pel_filters[xoffset][0]
	VFilter = &vp8_sub_pel_filters[yoffset][0]
	filter_block2d_first_pass((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), -(src_pixels_per_line*2))), &FData[0], uint(src_pixels_per_line), 1, 21, 16, HFilter)
	filter_block2d_second_pass(&FData[32], dst_ptr, dst_pitch, 16, 16, 16, 16, VFilter)
}
func filter_block2d_bil_first_pass(src_ptr *uint8, dst_ptr *uint16, src_stride uint, height uint, width uint, vp8_filter *int16) {
	var (
		i uint
		j uint
	)
	for i = 0; i < height; i++ {
		for j = 0; j < width; j++ {
			*(*uint16)(unsafe.Add(unsafe.Pointer(dst_ptr), unsafe.Sizeof(uint16(0))*uintptr(j))) = uint16(int16(((int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), 0))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*0)))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), 1)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*1))) + (int(VP8_FILTER_WEIGHT / 2))) >> VP8_FILTER_SHIFT))
			src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), 1))
		}
		src_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr), src_stride-width))
		dst_ptr = (*uint16)(unsafe.Add(unsafe.Pointer(dst_ptr), unsafe.Sizeof(uint16(0))*uintptr(width)))
	}
}
func filter_block2d_bil_second_pass(src_ptr *uint16, dst_ptr *uint8, dst_pitch int, height uint, width uint, vp8_filter *int16) {
	var (
		i    uint
		j    uint
		Temp int
	)
	for i = 0; i < height; i++ {
		for j = 0; j < width; j++ {
			Temp = (int(*(*uint16)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(uint16(0))*0))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*0)))) + int(*(*uint16)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(uint16(0))*uintptr(width))))*int(*(*int16)(unsafe.Add(unsafe.Pointer(vp8_filter), unsafe.Sizeof(int16(0))*1))) + (int(VP8_FILTER_WEIGHT / 2))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr), j)) = uint8(uint(Temp >> VP8_FILTER_SHIFT))
			src_ptr = (*uint16)(unsafe.Add(unsafe.Pointer(src_ptr), unsafe.Sizeof(uint16(0))*1))
		}
		dst_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr), dst_pitch))
	}
}
func filter_block2d_bil(src_ptr *uint8, dst_ptr *uint8, src_pitch uint, dst_pitch uint, HFilter *int16, VFilter *int16, Width int, Height int) {
	var FData [272]uint16
	filter_block2d_bil_first_pass(src_ptr, &FData[0], src_pitch, uint(Height+1), uint(Width), HFilter)
	filter_block2d_bil_second_pass(&FData[0], dst_ptr, int(dst_pitch), uint(Height), uint(Width), VFilter)
}
func Vp8BilinearPredict4x4C(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
	var (
		HFilter *int16
		VFilter *int16
	)
	if (xoffset | yoffset) != 0 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	HFilter = &vp8_bilinear_filters[xoffset][0]
	VFilter = &vp8_bilinear_filters[yoffset][0]
	filter_block2d_bil(src_ptr, dst_ptr, uint(src_pixels_per_line), uint(dst_pitch), HFilter, VFilter, 4, 4)
}
func Vp8BilinearPredict8x8C(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
	var (
		HFilter *int16
		VFilter *int16
	)
	if (xoffset | yoffset) != 0 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	HFilter = &vp8_bilinear_filters[xoffset][0]
	VFilter = &vp8_bilinear_filters[yoffset][0]
	filter_block2d_bil(src_ptr, dst_ptr, uint(src_pixels_per_line), uint(dst_pitch), HFilter, VFilter, 8, 8)
}
func Vp8BilinearPredict8x4C(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
	var (
		HFilter *int16
		VFilter *int16
	)
	if (xoffset | yoffset) != 0 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	HFilter = &vp8_bilinear_filters[xoffset][0]
	VFilter = &vp8_bilinear_filters[yoffset][0]
	filter_block2d_bil(src_ptr, dst_ptr, uint(src_pixels_per_line), uint(dst_pitch), HFilter, VFilter, 8, 4)
}
func Vp8BilinearPredict16x16C(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
	var (
		HFilter *int16
		VFilter *int16
	)
	if (xoffset | yoffset) != 0 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	HFilter = &vp8_bilinear_filters[xoffset][0]
	VFilter = &vp8_bilinear_filters[yoffset][0]
	filter_block2d_bil(src_ptr, dst_ptr, uint(src_pixels_per_line), uint(dst_pitch), HFilter, VFilter, 16, 16)
}
