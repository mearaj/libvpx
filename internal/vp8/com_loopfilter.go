package vp8

import "unsafe"

func vp8_loop_filter_mbh_sse2(y_ptr *uint8, u_ptr *uint8, v_ptr *uint8, y_stride int, uv_stride int, lfi *loop_filter_info) {
	vp8_mbloop_filter_horizontal_edge_sse2(y_ptr, y_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr)
	if u_ptr != nil {
		vp8_mbloop_filter_horizontal_edge_uv_sse2(u_ptr, uv_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr, v_ptr)
	}
}
func vp8_loop_filter_mbv_sse2(y_ptr *uint8, u_ptr *uint8, v_ptr *uint8, y_stride int, uv_stride int, lfi *loop_filter_info) {
	vp8_mbloop_filter_vertical_edge_sse2(y_ptr, y_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr)
	if u_ptr != nil {
		vp8_mbloop_filter_vertical_edge_uv_sse2(u_ptr, uv_stride, lfi.Mblim, lfi.Lim, lfi.Hev_thr, v_ptr)
	}
}
func vp8_loop_filter_bh_sse2(y_ptr *uint8, u_ptr *uint8, v_ptr *uint8, y_stride int, uv_stride int, lfi *loop_filter_info) {
	vp8_loop_filter_bh_y_sse2(y_ptr, y_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 2)
	if u_ptr != nil {
		vp8_loop_filter_horizontal_edge_uv_sse2((*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), uv_stride*4)), uv_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), uv_stride*4)))
	}
}
func vp8_loop_filter_bhs_sse2(y_ptr *uint8, y_stride int, blimit *uint8) {
	Vp8LoopFilterSimpleHorizontalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*4)), y_stride, blimit)
	Vp8LoopFilterSimpleHorizontalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*8)), y_stride, blimit)
	Vp8LoopFilterSimpleHorizontalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), y_stride*12)), y_stride, blimit)
}
func vp8_loop_filter_bv_sse2(y_ptr *uint8, u_ptr *uint8, v_ptr *uint8, y_stride int, uv_stride int, lfi *loop_filter_info) {
	vp8_loop_filter_bv_y_sse2(y_ptr, y_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, 2)
	if u_ptr != nil {
		vp8_loop_filter_vertical_edge_uv_sse2((*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), 4)), uv_stride, lfi.Blim, lfi.Lim, lfi.Hev_thr, (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), 4)))
	}
}
func vp8_loop_filter_bvs_sse2(y_ptr *uint8, y_stride int, blimit *uint8) {
	Vp8LoopFilterSimpleVerticalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 4)), y_stride, blimit)
	Vp8LoopFilterSimpleVerticalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 8)), y_stride, blimit)
	Vp8LoopFilterSimpleVerticalEdgeC((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 12)), y_stride, blimit)
}
