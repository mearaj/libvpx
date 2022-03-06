package scale

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func extend_plane(src *uint8, src_stride int, width int, height int, extend_top int, extend_left int, extend_bottom int, extend_right int) {
	var (
		i        int
		linesize int    = extend_left + extend_right + width
		src_ptr1 *uint8 = src
		src_ptr2 *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src), width))), -1))
		dst_ptr1 *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(src), -extend_left))
		dst_ptr2 *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(src), width))
	)
	for i = 0; i < height; i++ {
		libc.MemSet(unsafe.Pointer(dst_ptr1), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), 0))), extend_left)
		libc.MemSet(unsafe.Pointer(dst_ptr2), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), 0))), extend_right)
		src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), src_stride))
		src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), src_stride))
		dst_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr1), src_stride))
		dst_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr2), src_stride))
	}
	src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src), -extend_left))
	src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride*(height-1)))), -extend_left))
	dst_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride*(-extend_top)))), -extend_left))
	dst_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride*height))), -extend_left))
	for i = 0; i < extend_top; i++ {
		libc.MemCpy(unsafe.Pointer(dst_ptr1), unsafe.Pointer(src_ptr1), linesize)
		dst_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr1), src_stride))
	}
	for i = 0; i < extend_bottom; i++ {
		libc.MemCpy(unsafe.Pointer(dst_ptr2), unsafe.Pointer(src_ptr2), linesize)
		dst_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dst_ptr2), src_stride))
	}
}
func Vp8Yv12ExtendFrameBorders(ybf *Yv12BufferConfig) {
	var uv_border int = ybf.Border / 2
	if ybf.Border%2 == 0 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	if ybf.Y_height-ybf.Y_crop_height < 16 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	if ybf.Y_width-ybf.Y_crop_width < 16 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	if ybf.Y_height-ybf.Y_crop_height >= 0 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	if ybf.Y_width-ybf.Y_crop_width >= 0 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	extend_plane(ybf.Y_buffer, ybf.Y_stride, ybf.Y_crop_width, ybf.Y_crop_height, ybf.Border, ybf.Border, ybf.Border+ybf.Y_height-ybf.Y_crop_height, ybf.Border+ybf.Y_width-ybf.Y_crop_width)
	extend_plane(ybf.U_buffer, ybf.Uv_stride, ybf.Uv_crop_width, ybf.Uv_crop_height, uv_border, uv_border, uv_border+ybf.Uv_height-ybf.Uv_crop_height, uv_border+ybf.Uv_width-ybf.Uv_crop_width)
	extend_plane(ybf.V_buffer, ybf.Uv_stride, ybf.Uv_crop_width, ybf.Uv_crop_height, uv_border, uv_border, uv_border+ybf.Uv_height-ybf.Uv_crop_height, uv_border+ybf.Uv_width-ybf.Uv_crop_width)
}
func extend_frame(ybf *Yv12BufferConfig, ext_size int) {
	var (
		c_w  int = ybf.Uv_crop_width
		c_h  int = ybf.Uv_crop_height
		ss_x int = int(libc.BoolToInt(ybf.Uv_width < ybf.Y_width))
		ss_y int = int(libc.BoolToInt(ybf.Uv_height < ybf.Y_height))
		c_et int = ext_size >> ss_y
		c_el int = ext_size >> ss_x
		c_eb int = c_et + ybf.Uv_height - ybf.Uv_crop_height
		c_er int = c_el + ybf.Uv_width - ybf.Uv_crop_width
	)
	if ybf.Y_height-ybf.Y_crop_height < 16 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	if ybf.Y_width-ybf.Y_crop_width < 16 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	if ybf.Y_height-ybf.Y_crop_height >= 0 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	if ybf.Y_width-ybf.Y_crop_width >= 0 {
	} else {
		// Todo:
		log.Fatal("error")

	}
	extend_plane(ybf.Y_buffer, ybf.Y_stride, ybf.Y_crop_width, ybf.Y_crop_height, ext_size, ext_size, ext_size+ybf.Y_height-ybf.Y_crop_height, ext_size+ybf.Y_width-ybf.Y_crop_width)
	extend_plane(ybf.U_buffer, ybf.Uv_stride, c_w, c_h, c_et, c_el, c_eb, c_er)
	extend_plane(ybf.V_buffer, ybf.Uv_stride, c_w, c_h, c_et, c_el, c_eb, c_er)
}
func vpx_extend_frame_borders_c(ybf *Yv12BufferConfig) {
	extend_frame(ybf, ybf.Border)
}
func vpx_extend_frame_inner_borders_c(ybf *Yv12BufferConfig) {
	var inner_bw int
	if ybf.Border > VP9INNERBORDERINPIXELS {
		inner_bw = VP9INNERBORDERINPIXELS
	} else {
		inner_bw = ybf.Border
	}
	extend_frame(ybf, inner_bw)
}
func Vp8Yv12CopyFrameC(src_ybc *Yv12BufferConfig, dst_ybc *Yv12BufferConfig) {
	var (
		row int
		src *uint8 = src_ybc.Y_buffer
		dst *uint8 = dst_ybc.Y_buffer
	)
	for row = 0; row < src_ybc.Y_height; row++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), src_ybc.Y_width)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_ybc.Y_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_ybc.Y_stride))
	}
	src = src_ybc.U_buffer
	dst = dst_ybc.U_buffer
	for row = 0; row < src_ybc.Uv_height; row++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), src_ybc.Uv_width)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_ybc.Uv_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_ybc.Uv_stride))
	}
	src = src_ybc.V_buffer
	dst = dst_ybc.V_buffer
	for row = 0; row < src_ybc.Uv_height; row++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), src_ybc.Uv_width)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_ybc.Uv_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_ybc.Uv_stride))
	}
	Vp8Yv12ExtendFrameBorders(dst_ybc)
}
func vpx_yv12_copy_frame_c(src_ybc *Yv12BufferConfig, dst_ybc *Yv12BufferConfig) {
	var (
		row int
		src *uint8 = src_ybc.Y_buffer
		dst *uint8 = dst_ybc.Y_buffer
	)
	for row = 0; row < src_ybc.Y_height; row++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), src_ybc.Y_width)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_ybc.Y_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_ybc.Y_stride))
	}
	src = src_ybc.U_buffer
	dst = dst_ybc.U_buffer
	for row = 0; row < src_ybc.Uv_height; row++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), src_ybc.Uv_width)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_ybc.Uv_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_ybc.Uv_stride))
	}
	src = src_ybc.V_buffer
	dst = dst_ybc.V_buffer
	for row = 0; row < src_ybc.Uv_height; row++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), src_ybc.Uv_width)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_ybc.Uv_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_ybc.Uv_stride))
	}
	vpx_extend_frame_borders_c(dst_ybc)
}
func vpx_yv12_copy_y_c(src_ybc *Yv12BufferConfig, dst_ybc *Yv12BufferConfig) {
	var (
		row int
		src *uint8 = src_ybc.Y_buffer
		dst *uint8 = dst_ybc.Y_buffer
	)
	for row = 0; row < src_ybc.Y_height; row++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), src_ybc.Y_width)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_ybc.Y_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_ybc.Y_stride))
	}
}
