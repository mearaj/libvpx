package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func Vp8CopyMem16x16C(src *uint8, src_stride int, dst *uint8, dst_stride int) {
	var r int
	for r = 0; r < 16; r++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), 16)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_stride))
	}
}
func Vp8CopyMem8x8C(src *uint8, src_stride int, dst *uint8, dst_stride int) {
	var r int
	for r = 0; r < 8; r++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), 8)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_stride))
	}
}
func vp8_copy_mem8x4_c(src *uint8, src_stride int, dst *uint8, dst_stride int) {
	var r int
	for r = 0; r < 4; r++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(src), 8)
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_stride))
	}
}
func vp8_build_inter_predictors_b(d *Blockd, pitch int, base_pre *uint8, pre_stride int, sppf vp8_subpix_fn_t) {
	var (
		r        int
		pred_ptr *uint8 = d.Predictor
		ptr      *uint8
	)
	ptr = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(base_pre), d.Offset))), (int(d.Bmi.Mv.As_mv.Row)>>3)*pre_stride))), int(d.Bmi.Mv.As_mv.Col)>>3))
	if int(d.Bmi.Mv.As_mv.Row)&7 != 0 || int(d.Bmi.Mv.As_mv.Col)&7 != 0 {
		sppf(ptr, pre_stride, int(d.Bmi.Mv.As_mv.Col)&7, int(d.Bmi.Mv.As_mv.Row)&7, pred_ptr, pitch)
	} else {
		for r = 0; r < 4; r++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), 0)) = *(*uint8)(unsafe.Add(unsafe.Pointer(ptr), 0))
			*(*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), 1)) = *(*uint8)(unsafe.Add(unsafe.Pointer(ptr), 1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), 2)) = *(*uint8)(unsafe.Add(unsafe.Pointer(ptr), 2))
			*(*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), 3)) = *(*uint8)(unsafe.Add(unsafe.Pointer(ptr), 3))
			pred_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(pred_ptr), pitch))
			ptr = (*uint8)(unsafe.Add(unsafe.Pointer(ptr), pre_stride))
		}
	}
}
func build_inter_predictors4b(x *MacroBlockd, d *Blockd, dst *uint8, dst_stride int, base_pre *uint8, pre_stride int) {
	var ptr *uint8
	ptr = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(base_pre), d.Offset))), (int(d.Bmi.Mv.As_mv.Row)>>3)*pre_stride))), int(d.Bmi.Mv.As_mv.Col)>>3))
	if int(d.Bmi.Mv.As_mv.Row)&7 != 0 || int(d.Bmi.Mv.As_mv.Col)&7 != 0 {
		x.Subpixel_predict8x8(ptr, pre_stride, int(d.Bmi.Mv.As_mv.Col)&7, int(d.Bmi.Mv.As_mv.Row)&7, dst, dst_stride)
	} else {
		Vp8CopyMem8x8C(ptr, pre_stride, dst, dst_stride)
	}
}
func build_inter_predictors2b(x *MacroBlockd, d *Blockd, dst *uint8, dst_stride int, base_pre *uint8, pre_stride int) {
	var ptr *uint8
	ptr = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(base_pre), d.Offset))), (int(d.Bmi.Mv.As_mv.Row)>>3)*pre_stride))), int(d.Bmi.Mv.As_mv.Col)>>3))
	if int(d.Bmi.Mv.As_mv.Row)&7 != 0 || int(d.Bmi.Mv.As_mv.Col)&7 != 0 {
		x.Subpixel_predict8x4(ptr, pre_stride, int(d.Bmi.Mv.As_mv.Col)&7, int(d.Bmi.Mv.As_mv.Row)&7, dst, dst_stride)
	} else {
		vp8_copy_mem8x4_c(ptr, pre_stride, dst, dst_stride)
	}
}
func build_inter_predictors_b(d *Blockd, dst *uint8, dst_stride int, base_pre *uint8, pre_stride int, sppf vp8_subpix_fn_t) {
	var (
		r   int
		ptr *uint8
	)
	ptr = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(base_pre), d.Offset))), (int(d.Bmi.Mv.As_mv.Row)>>3)*pre_stride))), int(d.Bmi.Mv.As_mv.Col)>>3))
	if int(d.Bmi.Mv.As_mv.Row)&7 != 0 || int(d.Bmi.Mv.As_mv.Col)&7 != 0 {
		sppf(ptr, pre_stride, int(d.Bmi.Mv.As_mv.Col)&7, int(d.Bmi.Mv.As_mv.Row)&7, dst, dst_stride)
	} else {
		for r = 0; r < 4; r++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 0)) = *(*uint8)(unsafe.Add(unsafe.Pointer(ptr), 0))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 1)) = *(*uint8)(unsafe.Add(unsafe.Pointer(ptr), 1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 2)) = *(*uint8)(unsafe.Add(unsafe.Pointer(ptr), 2))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 3)) = *(*uint8)(unsafe.Add(unsafe.Pointer(ptr), 3))
			dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_stride))
			ptr = (*uint8)(unsafe.Add(unsafe.Pointer(ptr), pre_stride))
		}
	}
}
func vp8_build_inter16x16_predictors_mbuv(x *MacroBlockd) {
	var (
		uptr       *uint8
		vptr       *uint8
		upred_ptr  *uint8 = &x.Predictor[256]
		vpred_ptr  *uint8 = &x.Predictor[320]
		mv_row     int    = int(x.Mode_info_context.Mbmi.Mv.As_mv.Row)
		mv_col     int    = int(x.Mode_info_context.Mbmi.Mv.As_mv.Col)
		offset     int
		pre_stride int = x.Pre.Uv_stride
	)
	mv_row += (mv_row >> int(CHAR_BIT*unsafe.Sizeof(int(0))-1)) | 1
	mv_col += (mv_col >> int(CHAR_BIT*unsafe.Sizeof(int(0))-1)) | 1
	mv_row /= 2
	mv_col /= 2
	mv_row &= x.Fullpixel_mask
	mv_col &= x.Fullpixel_mask
	offset = (mv_row>>3)*pre_stride + (mv_col >> 3)
	uptr = (*uint8)(unsafe.Add(unsafe.Pointer(x.Pre.U_buffer), offset))
	vptr = (*uint8)(unsafe.Add(unsafe.Pointer(x.Pre.V_buffer), offset))
	if (mv_row|mv_col)&7 != 0 {
		x.Subpixel_predict8x8(uptr, pre_stride, mv_col&7, mv_row&7, upred_ptr, 8)
		x.Subpixel_predict8x8(vptr, pre_stride, mv_col&7, mv_row&7, vpred_ptr, 8)
	} else {
		Vp8CopyMem8x8C(uptr, pre_stride, upred_ptr, 8)
		Vp8CopyMem8x8C(vptr, pre_stride, vpred_ptr, 8)
	}
}
func vp8_build_inter4x4_predictors_mbuv(x *MacroBlockd) {
	var (
		i          int
		j          int
		pre_stride int = x.Pre.Uv_stride
		base_pre   *uint8
	)
	for i = 0; i < 2; i++ {
		for j = 0; j < 2; j++ {
			var (
				yoffset int = i*8 + j*2
				uoffset int = i*2 + 16 + j
				voffset int = i*2 + 20 + j
				temp    int
			)
			temp = int(x.Block[yoffset].Bmi.Mv.As_mv.Row) + int(x.Block[yoffset+1].Bmi.Mv.As_mv.Row) + int(x.Block[yoffset+4].Bmi.Mv.As_mv.Row) + int(x.Block[yoffset+5].Bmi.Mv.As_mv.Row)
			temp += ((temp >> int(CHAR_BIT*unsafe.Sizeof(int(0))-1)) * 8) + 4
			x.Block[uoffset].Bmi.Mv.As_mv.Row = int16((temp / 8) & x.Fullpixel_mask)
			temp = int(x.Block[yoffset].Bmi.Mv.As_mv.Col) + int(x.Block[yoffset+1].Bmi.Mv.As_mv.Col) + int(x.Block[yoffset+4].Bmi.Mv.As_mv.Col) + int(x.Block[yoffset+5].Bmi.Mv.As_mv.Col)
			temp += ((temp >> int(CHAR_BIT*unsafe.Sizeof(int(0))-1)) * 8) + 4
			x.Block[uoffset].Bmi.Mv.As_mv.Col = int16((temp / 8) & x.Fullpixel_mask)
			x.Block[voffset].Bmi.Mv.As_int = x.Block[uoffset].Bmi.Mv.As_int
		}
	}
	base_pre = (*uint8)(unsafe.Pointer(x.Pre.U_buffer))
	for i = 16; i < 20; i += 2 {
		var (
			d0 *Blockd = &x.Block[i]
			d1 *Blockd = &x.Block[i+1]
		)
		if d0.Bmi.Mv.As_int == d1.Bmi.Mv.As_int {
			build_inter_predictors2b(x, d0, d0.Predictor, 8, base_pre, pre_stride)
		} else {
			vp8_build_inter_predictors_b(d0, 8, base_pre, pre_stride, x.Subpixel_predict)
			vp8_build_inter_predictors_b(d1, 8, base_pre, pre_stride, x.Subpixel_predict)
		}
	}
	base_pre = (*uint8)(unsafe.Pointer(x.Pre.V_buffer))
	for i = 20; i < 24; i += 2 {
		var (
			d0 *Blockd = &x.Block[i]
			d1 *Blockd = &x.Block[i+1]
		)
		if d0.Bmi.Mv.As_int == d1.Bmi.Mv.As_int {
			build_inter_predictors2b(x, d0, d0.Predictor, 8, base_pre, pre_stride)
		} else {
			vp8_build_inter_predictors_b(d0, 8, base_pre, pre_stride, x.Subpixel_predict)
			vp8_build_inter_predictors_b(d1, 8, base_pre, pre_stride, x.Subpixel_predict)
		}
	}
}
func vp8_build_inter16x16_predictors_mby(x *MacroBlockd, dst_y *uint8, dst_ystride int) {
	var (
		ptr_base   *uint8
		ptr        *uint8
		mv_row     int = int(x.Mode_info_context.Mbmi.Mv.As_mv.Row)
		mv_col     int = int(x.Mode_info_context.Mbmi.Mv.As_mv.Col)
		pre_stride int = x.Pre.Y_stride
	)
	ptr_base = (*uint8)(unsafe.Pointer(x.Pre.Y_buffer))
	ptr = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ptr_base), (mv_row>>3)*pre_stride))), mv_col>>3))
	if (mv_row|mv_col)&7 != 0 {
		x.Subpixel_predict16x16(ptr, pre_stride, mv_col&7, mv_row&7, dst_y, dst_ystride)
	} else {
		Vp8CopyMem16x16C(ptr, pre_stride, dst_y, dst_ystride)
	}
}
func clamp_mv_to_umv_border(mv *MV, xd *MacroBlockd) {
	if int(mv.Col) < (xd.Mb_to_left_edge - (19 << 3)) {
		mv.Col = int16(xd.Mb_to_left_edge - (16 << 3))
	} else if int(mv.Col) > xd.Mb_to_right_edge+(18<<3) {
		mv.Col = int16(xd.Mb_to_right_edge + (16 << 3))
	}
	if int(mv.Row) < (xd.Mb_to_top_edge - (19 << 3)) {
		mv.Row = int16(xd.Mb_to_top_edge - (16 << 3))
	} else if int(mv.Row) > xd.Mb_to_bottom_edge+(18<<3) {
		mv.Row = int16(xd.Mb_to_bottom_edge + (16 << 3))
	}
}
func clamp_uvmv_to_umv_border(mv *MV, xd *MacroBlockd) {
	if int(mv.Col)*2 < (xd.Mb_to_left_edge - (19 << 3)) {
		mv.Col = int16((xd.Mb_to_left_edge - (16 << 3)) >> 1)
	} else {
		mv.Col = mv.Col
	}
	if int(mv.Col)*2 > xd.Mb_to_right_edge+(18<<3) {
		mv.Col = int16((xd.Mb_to_right_edge + (16 << 3)) >> 1)
	} else {
		mv.Col = mv.Col
	}
	if int(mv.Row)*2 < (xd.Mb_to_top_edge - (19 << 3)) {
		mv.Row = int16((xd.Mb_to_top_edge - (16 << 3)) >> 1)
	} else {
		mv.Row = mv.Row
	}
	if int(mv.Row)*2 > xd.Mb_to_bottom_edge+(18<<3) {
		mv.Row = int16((xd.Mb_to_bottom_edge + (16 << 3)) >> 1)
	} else {
		mv.Row = mv.Row
	}
}
func vp8_build_inter16x16_predictors_mb(x *MacroBlockd, dst_y *uint8, dst_u *uint8, dst_v *uint8, dst_ystride int, dst_uvstride int) {
	var (
		offset     int
		ptr        *uint8
		uptr       *uint8
		vptr       *uint8
		_16x16mv   int_mv
		ptr_base   *uint8 = (*uint8)(unsafe.Pointer(x.Pre.Y_buffer))
		pre_stride int    = x.Pre.Y_stride
	)
	_16x16mv.As_int = x.Mode_info_context.Mbmi.Mv.As_int
	if x.Mode_info_context.Mbmi.Need_to_clamp_mvs != 0 {
		clamp_mv_to_umv_border(&_16x16mv.As_mv, x)
	}
	ptr = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ptr_base), (int(_16x16mv.As_mv.Row)>>3)*pre_stride))), int(_16x16mv.As_mv.Col)>>3))
	if _16x16mv.As_int&0x70007 != 0 {
		x.Subpixel_predict16x16(ptr, pre_stride, int(_16x16mv.As_mv.Col)&7, int(_16x16mv.As_mv.Row)&7, dst_y, dst_ystride)
	} else {
		Vp8CopyMem16x16C(ptr, pre_stride, dst_y, dst_ystride)
	}
	_16x16mv.As_mv.Row += int16(uint16((uintptr(_16x16mv.As_mv.Row) >> (CHAR_BIT*unsafe.Sizeof(int(0)) - 1)) | 1))
	_16x16mv.As_mv.Col += int16(uint16((uintptr(_16x16mv.As_mv.Col) >> (CHAR_BIT*unsafe.Sizeof(int(0)) - 1)) | 1))
	_16x16mv.As_mv.Row /= 2
	_16x16mv.As_mv.Col /= 2
	_16x16mv.As_mv.Row &= int16(x.Fullpixel_mask)
	_16x16mv.As_mv.Col &= int16(x.Fullpixel_mask)
	if int(_16x16mv.As_mv.Col)*2 < (x.Mb_to_left_edge-(19<<3)) || int(_16x16mv.As_mv.Col)*2 > x.Mb_to_right_edge+(18<<3) || int(_16x16mv.As_mv.Row)*2 < (x.Mb_to_top_edge-(19<<3)) || int(_16x16mv.As_mv.Row)*2 > x.Mb_to_bottom_edge+(18<<3) {
		return
	}
	pre_stride >>= 1
	offset = (int(_16x16mv.As_mv.Row)>>3)*pre_stride + (int(_16x16mv.As_mv.Col) >> 3)
	uptr = (*uint8)(unsafe.Add(unsafe.Pointer(x.Pre.U_buffer), offset))
	vptr = (*uint8)(unsafe.Add(unsafe.Pointer(x.Pre.V_buffer), offset))
	if _16x16mv.As_int&0x70007 != 0 {
		x.Subpixel_predict8x8(uptr, pre_stride, int(_16x16mv.As_mv.Col)&7, int(_16x16mv.As_mv.Row)&7, dst_u, dst_uvstride)
		x.Subpixel_predict8x8(vptr, pre_stride, int(_16x16mv.As_mv.Col)&7, int(_16x16mv.As_mv.Row)&7, dst_v, dst_uvstride)
	} else {
		Vp8CopyMem8x8C(uptr, pre_stride, dst_u, dst_uvstride)
		Vp8CopyMem8x8C(vptr, pre_stride, dst_v, dst_uvstride)
	}
}
func build_inter4x4_predictors_mb(x *MacroBlockd) {
	var (
		i        int
		base_dst *uint8 = (*uint8)(unsafe.Pointer(x.Dst.Y_buffer))
		base_pre *uint8 = (*uint8)(unsafe.Pointer(x.Pre.Y_buffer))
	)
	if x.Mode_info_context.Mbmi.Partitioning < 3 {
		var (
			b          *Blockd
			dst_stride int = x.Dst.Y_stride
		)
		x.Block[0].Bmi = x.Mode_info_context.Bmi[0]
		x.Block[2].Bmi = x.Mode_info_context.Bmi[2]
		x.Block[8].Bmi = x.Mode_info_context.Bmi[8]
		x.Block[10].Bmi = x.Mode_info_context.Bmi[10]
		if x.Mode_info_context.Mbmi.Need_to_clamp_mvs != 0 {
			clamp_mv_to_umv_border(&x.Block[0].Bmi.Mv.As_mv, x)
			clamp_mv_to_umv_border(&x.Block[2].Bmi.Mv.As_mv, x)
			clamp_mv_to_umv_border(&x.Block[8].Bmi.Mv.As_mv, x)
			clamp_mv_to_umv_border(&x.Block[10].Bmi.Mv.As_mv, x)
		}
		b = &x.Block[0]
		build_inter_predictors4b(x, b, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), b.Offset)), dst_stride, base_pre, dst_stride)
		b = &x.Block[2]
		build_inter_predictors4b(x, b, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), b.Offset)), dst_stride, base_pre, dst_stride)
		b = &x.Block[8]
		build_inter_predictors4b(x, b, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), b.Offset)), dst_stride, base_pre, dst_stride)
		b = &x.Block[10]
		build_inter_predictors4b(x, b, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), b.Offset)), dst_stride, base_pre, dst_stride)
	} else {
		for i = 0; i < 16; i += 2 {
			var (
				d0         *Blockd = &x.Block[i]
				d1         *Blockd = &x.Block[i+1]
				dst_stride int     = x.Dst.Y_stride
			)
			x.Block[i+0].Bmi = x.Mode_info_context.Bmi[i+0]
			x.Block[i+1].Bmi = x.Mode_info_context.Bmi[i+1]
			if x.Mode_info_context.Mbmi.Need_to_clamp_mvs != 0 {
				clamp_mv_to_umv_border(&x.Block[i+0].Bmi.Mv.As_mv, x)
				clamp_mv_to_umv_border(&x.Block[i+1].Bmi.Mv.As_mv, x)
			}
			if d0.Bmi.Mv.As_int == d1.Bmi.Mv.As_int {
				build_inter_predictors2b(x, d0, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), d0.Offset)), dst_stride, base_pre, dst_stride)
			} else {
				build_inter_predictors_b(d0, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), d0.Offset)), dst_stride, base_pre, dst_stride, x.Subpixel_predict)
				build_inter_predictors_b(d1, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), d1.Offset)), dst_stride, base_pre, dst_stride, x.Subpixel_predict)
			}
		}
	}
	base_dst = (*uint8)(unsafe.Pointer(x.Dst.U_buffer))
	base_pre = (*uint8)(unsafe.Pointer(x.Pre.U_buffer))
	for i = 16; i < 20; i += 2 {
		var (
			d0         *Blockd = &x.Block[i]
			d1         *Blockd = &x.Block[i+1]
			dst_stride int     = x.Dst.Uv_stride
		)
		if d0.Bmi.Mv.As_int == d1.Bmi.Mv.As_int {
			build_inter_predictors2b(x, d0, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), d0.Offset)), dst_stride, base_pre, dst_stride)
		} else {
			build_inter_predictors_b(d0, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), d0.Offset)), dst_stride, base_pre, dst_stride, x.Subpixel_predict)
			build_inter_predictors_b(d1, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), d1.Offset)), dst_stride, base_pre, dst_stride, x.Subpixel_predict)
		}
	}
	base_dst = (*uint8)(unsafe.Pointer(x.Dst.V_buffer))
	base_pre = (*uint8)(unsafe.Pointer(x.Pre.V_buffer))
	for i = 20; i < 24; i += 2 {
		var (
			d0         *Blockd = &x.Block[i]
			d1         *Blockd = &x.Block[i+1]
			dst_stride int     = x.Dst.Uv_stride
		)
		if d0.Bmi.Mv.As_int == d1.Bmi.Mv.As_int {
			build_inter_predictors2b(x, d0, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), d0.Offset)), dst_stride, base_pre, dst_stride)
		} else {
			build_inter_predictors_b(d0, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), d0.Offset)), dst_stride, base_pre, dst_stride, x.Subpixel_predict)
			build_inter_predictors_b(d1, (*uint8)(unsafe.Add(unsafe.Pointer(base_dst), d1.Offset)), dst_stride, base_pre, dst_stride, x.Subpixel_predict)
		}
	}
}
func build_4x4uvmvs(x *MacroBlockd) {
	var (
		i int
		j int
	)
	for i = 0; i < 2; i++ {
		for j = 0; j < 2; j++ {
			var (
				yoffset int = i*8 + j*2
				uoffset int = i*2 + 16 + j
				voffset int = i*2 + 20 + j
				temp    int
			)
			temp = int(x.Mode_info_context.Bmi[yoffset+0].Mv.As_mv.Row) + int(x.Mode_info_context.Bmi[yoffset+1].Mv.As_mv.Row) + int(x.Mode_info_context.Bmi[yoffset+4].Mv.As_mv.Row) + int(x.Mode_info_context.Bmi[yoffset+5].Mv.As_mv.Row)
			temp += ((temp >> int(CHAR_BIT*unsafe.Sizeof(int(0))-1)) * 8) + 4
			x.Block[uoffset].Bmi.Mv.As_mv.Row = int16((temp / 8) & x.Fullpixel_mask)
			temp = int(x.Mode_info_context.Bmi[yoffset+0].Mv.As_mv.Col) + int(x.Mode_info_context.Bmi[yoffset+1].Mv.As_mv.Col) + int(x.Mode_info_context.Bmi[yoffset+4].Mv.As_mv.Col) + int(x.Mode_info_context.Bmi[yoffset+5].Mv.As_mv.Col)
			temp += ((temp >> int(CHAR_BIT*unsafe.Sizeof(int(0))-1)) * 8) + 4
			x.Block[uoffset].Bmi.Mv.As_mv.Col = int16((temp / 8) & x.Fullpixel_mask)
			if x.Mode_info_context.Mbmi.Need_to_clamp_mvs != 0 {
				clamp_uvmv_to_umv_border(&x.Block[uoffset].Bmi.Mv.As_mv, x)
			}
			x.Block[voffset].Bmi.Mv.As_int = x.Block[uoffset].Bmi.Mv.As_int
		}
	}
}
func vp8_build_inter_predictors_mb(xd *MacroBlockd) {
	if int(xd.Mode_info_context.Mbmi.Mode) != SPLITMV {
		vp8_build_inter16x16_predictors_mb(xd, (*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.U_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.V_buffer)), xd.Dst.Y_stride, xd.Dst.Uv_stride)
	} else {
		build_4x4uvmvs(xd)
		build_inter4x4_predictors_mb(xd)
	}
}
