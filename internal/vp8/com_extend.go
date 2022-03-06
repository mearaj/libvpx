package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/scale"
	"unsafe"
)

func copy_and_extend_plane(s *uint8, sp int, d *uint8, dp int, h int, w int, et int, el int, eb int, er int, interleave_step int) {
	var (
		i         int
		j         int
		src_ptr1  *uint8
		src_ptr2  *uint8
		dest_ptr1 *uint8
		dest_ptr2 *uint8
		linesize  int
	)
	if interleave_step < 1 {
		interleave_step = 1
	}
	src_ptr1 = s
	src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(s), (w-1)*interleave_step))
	dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(d), -el))
	dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(d), w))
	for i = 0; i < h; i++ {
		libc.MemSet(unsafe.Pointer(dest_ptr1), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), 0))), el)
		if interleave_step == 1 {
			libc.MemCpy(unsafe.Add(unsafe.Pointer(dest_ptr1), el), unsafe.Pointer(src_ptr1), w)
		} else {
			for j = 0; j < w; j++ {
				*(*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr1), el+j)) = *(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), interleave_step*j))
			}
		}
		libc.MemSet(unsafe.Pointer(dest_ptr2), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), 0))), er)
		src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), sp))
		src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), sp))
		dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr1), dp))
		dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr2), dp))
	}
	src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(d), -el))
	src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(d), dp*(h-1)))), -el))
	dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(d), dp*(-et)))), -el))
	dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(d), dp*h))), -el))
	linesize = el + er + w
	for i = 0; i < et; i++ {
		libc.MemCpy(unsafe.Pointer(dest_ptr1), unsafe.Pointer(src_ptr1), linesize)
		dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr1), dp))
	}
	for i = 0; i < eb; i++ {
		libc.MemCpy(unsafe.Pointer(dest_ptr2), unsafe.Pointer(src_ptr2), linesize)
		dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr2), dp))
	}
}
func vp8_copy_and_extend_frame(src *scale.Yv12BufferConfig, dst *scale.Yv12BufferConfig) {
	var (
		et = dst.Border
		el = dst.Border
		eb     = dst.Border + dst.Y_height - src.Y_height
		er     = dst.Border + dst.Y_width - src.Y_width
		chroma_step int
	)
	if int64(uintptr(unsafe.Pointer(src.V_buffer))-uintptr(unsafe.Pointer(src.U_buffer))) == 1 {
		chroma_step = 2
	} else {
		chroma_step = 1
	}
	copy_and_extend_plane((*uint8)(unsafe.Pointer(src.Y_buffer)), src.Y_stride, (*uint8)(unsafe.Pointer(dst.Y_buffer)), dst.Y_stride, src.Y_height, src.Y_width, et, el, eb, er, 1)
	et = dst.Border >> 1
	el = dst.Border >> 1
	eb = (dst.Border >> 1) + dst.Uv_height - src.Uv_height
	er = (dst.Border >> 1) + dst.Uv_width - src.Uv_width
	copy_and_extend_plane((*uint8)(unsafe.Pointer(src.U_buffer)), src.Uv_stride, (*uint8)(unsafe.Pointer(dst.U_buffer)), dst.Uv_stride, src.Uv_height, src.Uv_width, et, el, eb, er, chroma_step)
	copy_and_extend_plane((*uint8)(unsafe.Pointer(src.V_buffer)), src.Uv_stride, (*uint8)(unsafe.Pointer(dst.V_buffer)), dst.Uv_stride, src.Uv_height, src.Uv_width, et, el, eb, er, chroma_step)
}
func vp8_copy_and_extend_frame_with_rect(src *scale.Yv12BufferConfig, dst *scale.Yv12BufferConfig, srcy int, srcx int, srch int, srcw int) {
	var (
		et = dst.Border
		el = dst.Border
		eb     = dst.Border + dst.Y_height - src.Y_height
		er     = dst.Border + dst.Y_width - src.Y_width
		src_y_offset     = srcy*src.Y_stride + srcx
		dst_y_offset     = srcy*dst.Y_stride + srcx
		src_uv_offset     = ((srcy * src.Uv_stride) >> 1) + (srcx >> 1)
		dst_uv_offset     = ((srcy * dst.Uv_stride) >> 1) + (srcx >> 1)
		chroma_step   int
	)
	if int64(uintptr(unsafe.Pointer(src.V_buffer))-uintptr(unsafe.Pointer(src.U_buffer))) == 1 {
		chroma_step = 2
	} else {
		chroma_step = 1
	}
	if srcy != 0 {
		et = 0
	}
	if srcx != 0 {
		el = 0
	}
	if srcy+srch != src.Y_height {
		eb = 0
	}
	if srcx+srcw != src.Y_width {
		er = 0
	}
	copy_and_extend_plane((*uint8)(unsafe.Add(unsafe.Pointer(src.Y_buffer), src_y_offset)), src.Y_stride, (*uint8)(unsafe.Add(unsafe.Pointer(dst.Y_buffer), dst_y_offset)), dst.Y_stride, srch, srcw, et, el, eb, er, 1)
	et = (et + 1) >> 1
	el = (el + 1) >> 1
	eb = (eb + 1) >> 1
	er = (er + 1) >> 1
	srch = (srch + 1) >> 1
	srcw = (srcw + 1) >> 1
	copy_and_extend_plane((*uint8)(unsafe.Add(unsafe.Pointer(src.U_buffer), src_uv_offset)), src.Uv_stride, (*uint8)(unsafe.Add(unsafe.Pointer(dst.U_buffer), dst_uv_offset)), dst.Uv_stride, srch, srcw, et, el, eb, er, chroma_step)
	copy_and_extend_plane((*uint8)(unsafe.Add(unsafe.Pointer(src.V_buffer), src_uv_offset)), src.Uv_stride, (*uint8)(unsafe.Add(unsafe.Pointer(dst.V_buffer), dst_uv_offset)), dst.Uv_stride, srch, srcw, et, el, eb, er, chroma_step)
}
func vp8_extend_mb_row(ybf *scale.Yv12BufferConfig, YPtr *uint8, UPtr *uint8, VPtr *uint8) {
	var i int
	YPtr = (*uint8)(unsafe.Add(unsafe.Pointer(YPtr), ybf.Y_stride*14))
	UPtr = (*uint8)(unsafe.Add(unsafe.Pointer(UPtr), ybf.Uv_stride*6))
	VPtr = (*uint8)(unsafe.Add(unsafe.Pointer(VPtr), ybf.Uv_stride*6))
	for i = 0; i < 4; i++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(YPtr), i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(YPtr), -1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(UPtr), i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(UPtr), -1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(VPtr), i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(VPtr), -1))
	}
	YPtr = (*uint8)(unsafe.Add(unsafe.Pointer(YPtr), ybf.Y_stride))
	UPtr = (*uint8)(unsafe.Add(unsafe.Pointer(UPtr), ybf.Uv_stride))
	VPtr = (*uint8)(unsafe.Add(unsafe.Pointer(VPtr), ybf.Uv_stride))
	for i = 0; i < 4; i++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(YPtr), i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(YPtr), -1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(UPtr), i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(UPtr), -1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(VPtr), i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(VPtr), -1))
	}
}
