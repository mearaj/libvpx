package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const MAX_PSNR = 100.0

type vpx_psnr_pkt struct {
	Samples [4]uint
	Sse     [4]uint64
	Psnr    [4]float64
}
type PSNR_STATS vpx_psnr_pkt

func vpx_sse_to_psnr(samples float64, peak float64, sse float64) float64 {
	if sse > 0.0 {
		var psnr float64 = log10(samples*peak*peak/sse) * 10.0
		if psnr > MAX_PSNR {
			return MAX_PSNR
		}
		return psnr
	} else {
		return MAX_PSNR
	}
}
func encoder_variance(a *uint8, a_stride int, b *uint8, b_stride int, w int, h int, sse *uint, sum *int) {
	var (
		i int
		j int
	)
	*sum = 0
	*sse = 0
	for i = 0; i < h; i++ {
		for j = 0; j < w; j++ {
			var diff int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(a), j)) - *(*uint8)(unsafe.Add(unsafe.Pointer(b), j)))
			*sum += diff
			*sse += uint(diff * diff)
		}
		a = (*uint8)(unsafe.Add(unsafe.Pointer(a), a_stride))
		b = (*uint8)(unsafe.Add(unsafe.Pointer(b), b_stride))
	}
}
func get_sse(a *uint8, a_stride int, b *uint8, b_stride int, width int, height int) int64 {
	var (
		dw        int   = width % 16
		dh        int   = height % 16
		total_sse int64 = 0
		sse       uint  = 0
		sum       int   = 0
		x         int
		y         int
	)
	if dw > 0 {
		encoder_variance((*uint8)(unsafe.Add(unsafe.Pointer(a), width-dw)), a_stride, (*uint8)(unsafe.Add(unsafe.Pointer(b), width-dw)), b_stride, dw, height, &sse, &sum)
		total_sse += int64(sse)
	}
	if dh > 0 {
		encoder_variance((*uint8)(unsafe.Add(unsafe.Pointer(a), (height-dh)*a_stride)), a_stride, (*uint8)(unsafe.Add(unsafe.Pointer(b), (height-dh)*b_stride)), b_stride, width-dw, dh, &sse, &sum)
		total_sse += int64(sse)
	}
	for y = 0; y < height/16; y++ {
		var (
			pa *uint8 = a
			pb *uint8 = b
		)
		for x = 0; x < width/16; x++ {
			vpx_mse16x16(pa, a_stride, pb, b_stride, &sse)
			total_sse += int64(sse)
			pa = (*uint8)(unsafe.Add(unsafe.Pointer(pa), 16))
			pb = (*uint8)(unsafe.Add(unsafe.Pointer(pb), 16))
		}
		a = (*uint8)(unsafe.Add(unsafe.Pointer(a), a_stride*16))
		b = (*uint8)(unsafe.Add(unsafe.Pointer(b), b_stride*16))
	}
	return total_sse
}
func vpx_get_y_sse(a *YV12_BUFFER_CONFIG, b *YV12_BUFFER_CONFIG) int64 {
	if a.Y_crop_width == b.Y_crop_width {
	} else {
		__assert_fail(libc.CString("a->y_crop_width == b->y_crop_width"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	if a.Y_crop_height == b.Y_crop_height {
	} else {
		__assert_fail(libc.CString("a->y_crop_height == b->y_crop_height"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	return get_sse(a.Y_buffer, a.Y_stride, b.Y_buffer, b.Y_stride, a.Y_crop_width, a.Y_crop_height)
}
func vpx_calc_psnr(a *YV12_BUFFER_CONFIG, b *YV12_BUFFER_CONFIG, psnr *PSNR_STATS) {
	var (
		peak          float64   = 255.0
		widths        [3]int    = [3]int{a.Y_crop_width, a.Uv_crop_width, a.Uv_crop_width}
		heights       [3]int    = [3]int{a.Y_crop_height, a.Uv_crop_height, a.Uv_crop_height}
		a_planes      [3]*uint8 = [3]*uint8{a.Y_buffer, a.U_buffer, a.V_buffer}
		a_strides     [3]int    = [3]int{a.Y_stride, a.Uv_stride, a.Uv_stride}
		b_planes      [3]*uint8 = [3]*uint8{b.Y_buffer, b.U_buffer, b.V_buffer}
		b_strides     [3]int    = [3]int{b.Y_stride, b.Uv_stride, b.Uv_stride}
		i             int
		total_sse     uint64 = 0
		total_samples uint32 = 0
	)
	for i = 0; i < 3; i++ {
		var (
			w       int    = widths[i]
			h       int    = heights[i]
			samples uint32 = uint32(w * h)
			sse     uint64 = uint64(get_sse(a_planes[i], a_strides[i], b_planes[i], b_strides[i], w, h))
		)
		psnr.Sse[i+1] = sse
		psnr.Samples[i+1] = uint(samples)
		psnr.Psnr[i+1] = vpx_sse_to_psnr(float64(samples), peak, float64(sse))
		total_sse += sse
		total_samples += samples
	}
	psnr.Sse[0] = total_sse
	psnr.Samples[0] = uint(total_samples)
	psnr.Psnr[0] = vpx_sse_to_psnr(float64(total_samples), peak, float64(total_sse))
}
