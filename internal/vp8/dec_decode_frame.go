package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/mearaj/libvpx/internal/scale"
	"github.com/mearaj/libvpx/internal/util"
	"github.com/mearaj/libvpx/internal/vpx"
	"log"
	"math"
	"unsafe"
)

func vp8cx_init_de_quantizer(pbi *VP8D_COMP) {
	var (
		Q  int
		pc = &pbi.Common
	)
	for Q = 0; Q < (int(MAXQ + 1)); Q++ {
		pc.Y1dequant[Q][0] = int16(vp8_dc_quant(Q, pc.Y1dc_delta_q))
		pc.Y2dequant[Q][0] = int16(vp8_dc2quant(Q, pc.Y2dc_delta_q))
		pc.UVdequant[Q][0] = int16(vp8_dc_uv_quant(Q, pc.Uvdc_delta_q))
		pc.Y1dequant[Q][1] = int16(vp8_ac_yquant(Q))
		pc.Y2dequant[Q][1] = int16(vp8_ac2quant(Q, pc.Y2ac_delta_q))
		pc.UVdequant[Q][1] = int16(vp8_ac_uv_quant(Q, pc.Uvac_delta_q))
	}
}
func vp8_mb_init_dequantizer(pbi *VP8D_COMP, xd *MacroBlockd) {
	var (
		i      int
		QIndex int
		mbmi   = &xd.Mode_info_context.Mbmi
		pc     = &pbi.Common
	)
	if int(xd.Segmentation_enabled) != 0 {
		if int(xd.Mb_segement_abs_delta) == SEGMENT_ABSDATA {
			QIndex = int(xd.Segment_feature_data[MB_LVL_ALT_Q][mbmi.Segment_id])
		} else {
			QIndex = pc.Base_qindex + int(xd.Segment_feature_data[MB_LVL_ALT_Q][mbmi.Segment_id])
		}
		if QIndex >= 0 {
			if QIndex <= MAXQ {
				QIndex = QIndex
			} else {
				QIndex = MAXQ
			}
		} else {
			QIndex = 0
		}
	} else {
		QIndex = pc.Base_qindex
	}
	xd.Dequant_y1_dc[0] = 1
	xd.Dequant_y1[0] = pc.Y1dequant[QIndex][0]
	xd.Dequant_y2[0] = pc.Y2dequant[QIndex][0]
	xd.Dequant_uv[0] = pc.UVdequant[QIndex][0]
	for i = 1; i < 16; i++ {
		xd.Dequant_y1_dc[i] = func() int16 {
			p := &xd.Dequant_y1[i]
			xd.Dequant_y1[i] = pc.Y1dequant[QIndex][1]
			return *p
		}()
		xd.Dequant_y2[i] = pc.Y2dequant[QIndex][1]
		xd.Dequant_uv[i] = pc.UVdequant[QIndex][1]
	}
}
func decode_macroblock(pbi *VP8D_COMP, xd *MacroBlockd, mb_idx uint) {
	var (
		mode int
		i    int
	)
	_ = mb_idx
	if xd.Mode_info_context.Mbmi.Mb_skip_coeff != 0 {
		vp8_reset_mb_tokens_context(xd)
	} else if vp8dx_bool_error((*BOOL_DECODER)(xd.Current_bc)) == 0 {
		var eobtotal int
		eobtotal = vp8_decode_mb_tokens(pbi, xd)
		xd.Mode_info_context.Mbmi.Mb_skip_coeff = uint8(int8(libc.BoolToInt(eobtotal == 0)))
	}
	mode = int(xd.Mode_info_context.Mbmi.Mode)
	if int(xd.Segmentation_enabled) != 0 {
		vp8_mb_init_dequantizer(pbi, xd)
	}
	if int(xd.Mode_info_context.Mbmi.Ref_frame) == INTRA_FRAME {
		vp8_build_intra_predictors_mbuv_s(xd, xd.Recon_above[1], xd.Recon_above[2], xd.Recon_left[1], xd.Recon_left[2], xd.Recon_left_stride[1], (*uint8)(unsafe.Pointer(xd.Dst.U_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.V_buffer)), xd.Dst.Uv_stride)
		if mode != int(B_PRED) {
			vp8_build_intra_predictors_mby_s(xd, xd.Recon_above[0], xd.Recon_left[0], xd.Recon_left_stride[0], (*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), xd.Dst.Y_stride)
		} else {
			var (
				DQC        = &xd.Dequant_y1[0]
				dst_stride = xd.Dst.Y_stride
			)
			if xd.Mode_info_context.Mbmi.Mb_skip_coeff != 0 {
				libc.MemSet(unsafe.Pointer(&xd.Eobs[0]), 0, 25)
			}
			intra_prediction_down_copy(xd, (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[0]), 16)))
			for i = 0; i < 16; i++ {
				var (
					b   = &xd.Block[i]
					dst = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.Y_buffer), b.Offset))
					b_mode         = xd.Mode_info_context.Bmi[i].As_mode
					Above         = (*uint8)(unsafe.Add(unsafe.Pointer(dst), -dst_stride))
					yleft                    = (*uint8)(unsafe.Add(unsafe.Pointer(dst), -1))
					left_stride        = dst_stride
					top_left           = *(*uint8)(unsafe.Add(unsafe.Pointer(Above), -1))
				)
				vp8_intra4x4_predict(Above, yleft, left_stride, b_mode, dst, dst_stride, top_left)
				if xd.Eobs[i] != 0 {
					if xd.Eobs[i] > 1 {
						Vp8DequantIdctAddC(b.Qcoeff, DQC, dst, dst_stride)
					} else {
						Vp8DcOnlyIdctAddC(int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(b.Qcoeff), unsafe.Sizeof(int16(0))*0)))*int(*(*int16)(unsafe.Add(unsafe.Pointer(DQC), unsafe.Sizeof(int16(0))*0)))), dst, dst_stride, dst, dst_stride)
						libc.MemSet(unsafe.Pointer(b.Qcoeff), 0, int(2*unsafe.Sizeof(int16(0))))
					}
				}
			}
		}
	} else {
		vp8_build_inter_predictors_mb(xd)
	}
	if xd.Mode_info_context.Mbmi.Mb_skip_coeff == 0 {
		if mode != int(B_PRED) {
			var DQC = &xd.Dequant_y1[0]
			if mode != int(SPLITMV) {
				var b = &xd.Block[24]
				if xd.Eobs[24] > 1 {
					vp8_dequantize_b_mmx(b, &xd.Dequant_y2[0])
					Vp8ShortInvWalsh4x4C((*int16)(unsafe.Add(unsafe.Pointer(b.Dqcoeff), unsafe.Sizeof(int16(0))*0)), &xd.Qcoeff[0])
					libc.MemSet(unsafe.Pointer(b.Qcoeff), 0, int(16*unsafe.Sizeof(int16(0))))
				} else {
					*(*int16)(unsafe.Add(unsafe.Pointer(b.Dqcoeff), unsafe.Sizeof(int16(0))*0)) = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(b.Qcoeff), unsafe.Sizeof(int16(0))*0))) * int(xd.Dequant_y2[0]))
					vp8_short_inv_walsh4x4_1_c((*int16)(unsafe.Add(unsafe.Pointer(b.Dqcoeff), unsafe.Sizeof(int16(0))*0)), &xd.Qcoeff[0])
					libc.MemSet(unsafe.Pointer(b.Qcoeff), 0, int(2*unsafe.Sizeof(int16(0))))
				}
				DQC = &xd.Dequant_y1_dc[0]
			}
			vp8_dequant_idct_add_y_block_sse2(&xd.Qcoeff[0], DQC, (*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), xd.Dst.Y_stride, &xd.Eobs[0])
		}
		vp8_dequant_idct_add_uv_block_sse2(&xd.Qcoeff[16*16], &xd.Dequant_uv[0], (*uint8)(unsafe.Pointer(xd.Dst.U_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.V_buffer)), xd.Dst.Uv_stride, &xd.Eobs[16])
	}
}
func get_delta_q(bc *vp8_reader, prev int, q_update *int) int {
	var ret_val = 0
	if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
		ret_val = vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 4)
		if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
			ret_val = -ret_val
		}
	}
	if ret_val != prev {
		*q_update = 1
	}
	return ret_val
}
func yv12_extend_frame_top_c(ybf *scale.Yv12BufferConfig) {
	var (
		i            int
		src_ptr1     *uint8
		dest_ptr1    *uint8
		Border       uint
		plane_stride int
	)
	Border = uint(ybf.Border)
	plane_stride = ybf.Y_stride
	src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(ybf.Y_buffer), -int(Border)))
	dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), -int(Border*uint(plane_stride))))
	for i = 0; i < int(Border); i++ {
		libc.MemCpy(unsafe.Pointer(dest_ptr1), unsafe.Pointer(src_ptr1), plane_stride)
		dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr1), plane_stride))
	}
	plane_stride = ybf.Uv_stride
	Border /= 2
	src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(ybf.U_buffer), -int(Border)))
	dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), -int(Border*uint(plane_stride))))
	for i = 0; i < int(Border); i++ {
		libc.MemCpy(unsafe.Pointer(dest_ptr1), unsafe.Pointer(src_ptr1), plane_stride)
		dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr1), plane_stride))
	}
	src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(ybf.V_buffer), -int(Border)))
	dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), -int(Border*uint(plane_stride))))
	for i = 0; i < int(Border); i++ {
		libc.MemCpy(unsafe.Pointer(dest_ptr1), unsafe.Pointer(src_ptr1), plane_stride)
		dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr1), plane_stride))
	}
}
func yv12_extend_frame_bottom_c(ybf *scale.Yv12BufferConfig) {
	var (
		i            int
		src_ptr1     *uint8
		src_ptr2     *uint8
		dest_ptr2    *uint8
		Border       uint
		plane_stride int
		plane_height int
	)
	Border = uint(ybf.Border)
	plane_stride = ybf.Y_stride
	plane_height = ybf.Y_height
	src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(ybf.Y_buffer), -int(Border)))
	src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), plane_height*plane_stride))), -plane_stride))
	dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), plane_stride))
	for i = 0; i < int(Border); i++ {
		libc.MemCpy(unsafe.Pointer(dest_ptr2), unsafe.Pointer(src_ptr2), plane_stride)
		dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr2), plane_stride))
	}
	plane_stride = ybf.Uv_stride
	plane_height = ybf.Uv_height
	Border /= 2
	src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(ybf.U_buffer), -int(Border)))
	src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), plane_height*plane_stride))), -plane_stride))
	dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), plane_stride))
	for i = 0; i < int(Border); i++ {
		libc.MemCpy(unsafe.Pointer(dest_ptr2), unsafe.Pointer(src_ptr2), plane_stride)
		dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr2), plane_stride))
	}
	src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(ybf.V_buffer), -int(Border)))
	src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), plane_height*plane_stride))), -plane_stride))
	dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), plane_stride))
	for i = 0; i < int(Border); i++ {
		libc.MemCpy(unsafe.Pointer(dest_ptr2), unsafe.Pointer(src_ptr2), plane_stride)
		dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr2), plane_stride))
	}
}
func yv12_extend_frame_left_right_c(ybf *scale.Yv12BufferConfig, y_src *uint8, u_src *uint8, v_src *uint8) {
	var (
		i            int
		src_ptr1     *uint8
		src_ptr2     *uint8
		dest_ptr1    *uint8
		dest_ptr2    *uint8
		Border       uint
		plane_stride int
		plane_height int
		plane_width  int
	)
	Border = uint(ybf.Border)
	plane_stride = ybf.Y_stride
	plane_height = 16
	plane_width = ybf.Y_width
	src_ptr1 = y_src
	src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), plane_width))), -1))
	dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), -int(Border)))
	dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), 1))
	for i = 0; i < plane_height; i++ {
		libc.MemSet(unsafe.Pointer(dest_ptr1), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), 0))), int(Border))
		libc.MemSet(unsafe.Pointer(dest_ptr2), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), 0))), int(Border))
		src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), plane_stride))
		src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), plane_stride))
		dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr1), plane_stride))
		dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr2), plane_stride))
	}
	plane_stride = ybf.Uv_stride
	plane_height = 8
	plane_width = ybf.Uv_width
	Border /= 2
	src_ptr1 = u_src
	src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), plane_width))), -1))
	dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), -int(Border)))
	dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), 1))
	for i = 0; i < plane_height; i++ {
		libc.MemSet(unsafe.Pointer(dest_ptr1), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), 0))), int(Border))
		libc.MemSet(unsafe.Pointer(dest_ptr2), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), 0))), int(Border))
		src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), plane_stride))
		src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), plane_stride))
		dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr1), plane_stride))
		dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr2), plane_stride))
	}
	src_ptr1 = v_src
	src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), plane_width))), -1))
	dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), -int(Border)))
	dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), 1))
	for i = 0; i < plane_height; i++ {
		libc.MemSet(unsafe.Pointer(dest_ptr1), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), 0))), int(Border))
		libc.MemSet(unsafe.Pointer(dest_ptr2), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), 0))), int(Border))
		src_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr1), plane_stride))
		src_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(src_ptr2), plane_stride))
		dest_ptr1 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr1), plane_stride))
		dest_ptr2 = (*uint8)(unsafe.Add(unsafe.Pointer(dest_ptr2), plane_stride))
	}
}
func decode_mb_rows(pbi *VP8D_COMP) {
	var (
		pc = &pbi.Common
		xd = &pbi.Mb
		lf_mic            = xd.Mode_info_context
		ibc                 = 0
		num_part           = int(1 << pc.Multi_token_partition)
		recon_yoffset int
		recon_uvoffset   int
		mb_row           int
		mb_col      int
		mb_idx      = 0
		yv12_fb_new = pbi.Dec_fb_ref[INTRA_FRAME]
		recon_y_stride     = yv12_fb_new.Y_stride
		recon_uv_stride                         = yv12_fb_new.Uv_stride
		ref_buffer      [4][3]*uint8
		dst_buffer       [3]*uint8
		lf_dst           [3]*uint8
		eb_dst           [3]*uint8
		i                int
		ref_fb_corrupted [4]int
	)
	ref_fb_corrupted[INTRA_FRAME] = 0
	for i = 1; i < MAX_REF_FRAMES; i++ {
		var this_fb = pbi.Dec_fb_ref[i]
		ref_buffer[i][0] = (*uint8)(unsafe.Pointer(this_fb.Y_buffer))
		ref_buffer[i][1] = (*uint8)(unsafe.Pointer(this_fb.U_buffer))
		ref_buffer[i][2] = (*uint8)(unsafe.Pointer(this_fb.V_buffer))
		ref_fb_corrupted[i] = this_fb.Corrupted
	}
	eb_dst[0] = func() *uint8 {
		p := &lf_dst[0]
		lf_dst[0] = func() *uint8 {
			p := &dst_buffer[0]
			dst_buffer[0] = (*uint8)(unsafe.Pointer(yv12_fb_new.Y_buffer))
			return *p
		}()
		return *p
	}()
	eb_dst[1] = func() *uint8 {
		p := &lf_dst[1]
		lf_dst[1] = func() *uint8 {
			p := &dst_buffer[1]
			dst_buffer[1] = (*uint8)(unsafe.Pointer(yv12_fb_new.U_buffer))
			return *p
		}()
		return *p
	}()
	eb_dst[2] = func() *uint8 {
		p := &lf_dst[2]
		lf_dst[2] = func() *uint8 {
			p := &dst_buffer[2]
			dst_buffer[2] = (*uint8)(unsafe.Pointer(yv12_fb_new.V_buffer))
			return *p
		}()
		return *p
	}()
	xd.Up_available = 0
	if pc.Filter_level != 0 {
		vp8_loop_filter_frame_init(pc, xd, pc.Filter_level)
	}
	vp8_setup_intra_recon_top_line(yv12_fb_new)
	for mb_row = 0; mb_row < pc.Mb_rows; mb_row++ {
		if num_part > 1 {
			xd.Current_bc = unsafe.Pointer(&pbi.Mbc[ibc])
			ibc++
			if ibc == num_part {
				ibc = 0
			}
		}
		recon_yoffset = mb_row * recon_y_stride * 16
		recon_uvoffset = mb_row * recon_uv_stride * 8
		xd.Above_context = pc.Above_context
		*xd.Left_context = ENTROPY_CONTEXT_PLANES{}
		xd.Left_available = 0
		xd.Mb_to_top_edge = -((mb_row * 16) << 3)
		xd.Mb_to_bottom_edge = ((pc.Mb_rows - 1 - mb_row) * 16) << 3
		xd.Recon_above[0] = (*uint8)(unsafe.Add(unsafe.Pointer(dst_buffer[0]), recon_yoffset))
		xd.Recon_above[1] = (*uint8)(unsafe.Add(unsafe.Pointer(dst_buffer[1]), recon_uvoffset))
		xd.Recon_above[2] = (*uint8)(unsafe.Add(unsafe.Pointer(dst_buffer[2]), recon_uvoffset))
		xd.Recon_left[0] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[0]), -1))
		xd.Recon_left[1] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[1]), -1))
		xd.Recon_left[2] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[2]), -1))
		xd.Recon_above[0] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[0]), -xd.Dst.Y_stride))
		xd.Recon_above[1] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[1]), -xd.Dst.Uv_stride))
		xd.Recon_above[2] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[2]), -xd.Dst.Uv_stride))
		xd.Recon_left_stride[0] = xd.Dst.Y_stride
		xd.Recon_left_stride[1] = xd.Dst.Uv_stride
		setup_intra_recon_left(xd.Recon_left[0], xd.Recon_left[1], xd.Recon_left[2], xd.Dst.Y_stride, xd.Dst.Uv_stride)
		for mb_col = 0; mb_col < pc.Mb_cols; mb_col++ {
			xd.Mb_to_left_edge = -((mb_col * 16) << 3)
			xd.Mb_to_right_edge = ((pc.Mb_cols - 1 - mb_col) * 16) << 3
			xd.Dst.Y_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(dst_buffer[0]), recon_yoffset))
			xd.Dst.U_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(dst_buffer[1]), recon_uvoffset))
			xd.Dst.V_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(dst_buffer[2]), recon_uvoffset))
			if int(xd.Mode_info_context.Mbmi.Ref_frame) >= LAST_FRAME {
				var ref = int(xd.Mode_info_context.Mbmi.Ref_frame)
				xd.Pre.Y_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(ref_buffer[ref][0]), recon_yoffset))
				xd.Pre.U_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(ref_buffer[ref][1]), recon_uvoffset))
				xd.Pre.V_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(ref_buffer[ref][2]), recon_uvoffset))
			} else {
				xd.Pre.Y_buffer = nil
				xd.Pre.U_buffer = nil
				xd.Pre.V_buffer = nil
			}
			xd.Corrupted |= ref_fb_corrupted[xd.Mode_info_context.Mbmi.Ref_frame]
			decode_macroblock(pbi, xd, uint(mb_idx))
			mb_idx++
			xd.Left_available = 1
			xd.Corrupted |= vp8dx_bool_error((*BOOL_DECODER)(xd.Current_bc))
			xd.Recon_above[0] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[0]), 16))
			xd.Recon_above[1] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[1]), 8))
			xd.Recon_above[2] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[2]), 8))
			xd.Recon_left[0] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_left[0]), 16))
			xd.Recon_left[1] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_left[1]), 8))
			xd.Recon_left[2] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_left[2]), 8))
			recon_yoffset += 16
			recon_uvoffset += 8
			xd.Mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(xd.Mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
			xd.Above_context = (*ENTROPY_CONTEXT_PLANES)(unsafe.Add(unsafe.Pointer(xd.Above_context), unsafe.Sizeof(ENTROPY_CONTEXT_PLANES{})*1))
		}
		vp8_extend_mb_row(yv12_fb_new, (*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.Y_buffer), 16)), (*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.U_buffer), 8)), (*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.V_buffer), 8)))
		xd.Mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(xd.Mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
		xd.Up_available = 1
		if pc.Filter_level != 0 {
			if mb_row > 0 {
				if pc.Filter_type == LOOPFILTERTYPE(NORMAL_LOOPFILTER) {
					vp8_loop_filter_row_normal(pc, lf_mic, mb_row-1, recon_y_stride, recon_uv_stride, lf_dst[0], lf_dst[1], lf_dst[2])
				} else {
					vp8_loop_filter_row_simple(pc, lf_mic, mb_row-1, recon_y_stride, lf_dst[0])
				}
				if mb_row > 1 {
					yv12_extend_frame_left_right_c(yv12_fb_new, eb_dst[0], eb_dst[1], eb_dst[2])
					eb_dst[0] = (*uint8)(unsafe.Add(unsafe.Pointer(eb_dst[0]), recon_y_stride*16))
					eb_dst[1] = (*uint8)(unsafe.Add(unsafe.Pointer(eb_dst[1]), recon_uv_stride*8))
					eb_dst[2] = (*uint8)(unsafe.Add(unsafe.Pointer(eb_dst[2]), recon_uv_stride*8))
				}
				lf_dst[0] = (*uint8)(unsafe.Add(unsafe.Pointer(lf_dst[0]), recon_y_stride*16))
				lf_dst[1] = (*uint8)(unsafe.Add(unsafe.Pointer(lf_dst[1]), recon_uv_stride*8))
				lf_dst[2] = (*uint8)(unsafe.Add(unsafe.Pointer(lf_dst[2]), recon_uv_stride*8))
				lf_mic = (*ModeInfo)(unsafe.Add(unsafe.Pointer(lf_mic), unsafe.Sizeof(ModeInfo{})*uintptr(pc.Mb_cols)))
				lf_mic = (*ModeInfo)(unsafe.Add(unsafe.Pointer(lf_mic), unsafe.Sizeof(ModeInfo{})*1))
			}
		} else {
			if mb_row > 0 {
				yv12_extend_frame_left_right_c(yv12_fb_new, eb_dst[0], eb_dst[1], eb_dst[2])
				eb_dst[0] = (*uint8)(unsafe.Add(unsafe.Pointer(eb_dst[0]), recon_y_stride*16))
				eb_dst[1] = (*uint8)(unsafe.Add(unsafe.Pointer(eb_dst[1]), recon_uv_stride*8))
				eb_dst[2] = (*uint8)(unsafe.Add(unsafe.Pointer(eb_dst[2]), recon_uv_stride*8))
			}
		}
	}
	if pc.Filter_level != 0 {
		if pc.Filter_type == LOOPFILTERTYPE(NORMAL_LOOPFILTER) {
			vp8_loop_filter_row_normal(pc, lf_mic, mb_row-1, recon_y_stride, recon_uv_stride, lf_dst[0], lf_dst[1], lf_dst[2])
		} else {
			vp8_loop_filter_row_simple(pc, lf_mic, mb_row-1, recon_y_stride, lf_dst[0])
		}
		yv12_extend_frame_left_right_c(yv12_fb_new, eb_dst[0], eb_dst[1], eb_dst[2])
		eb_dst[0] = (*uint8)(unsafe.Add(unsafe.Pointer(eb_dst[0]), recon_y_stride*16))
		eb_dst[1] = (*uint8)(unsafe.Add(unsafe.Pointer(eb_dst[1]), recon_uv_stride*8))
		eb_dst[2] = (*uint8)(unsafe.Add(unsafe.Pointer(eb_dst[2]), recon_uv_stride*8))
	}
	yv12_extend_frame_left_right_c(yv12_fb_new, eb_dst[0], eb_dst[1], eb_dst[2])
	yv12_extend_frame_top_c(yv12_fb_new)
	yv12_extend_frame_bottom_c(yv12_fb_new)
}
func read_partition_size(pbi *VP8D_COMP, cx_size *uint8) uint {
	var temp [3]uint8
	if pbi.Decrypt_cb != nil {
		pbi.Decrypt_cb(pbi.Decrypt_state, cx_size, &temp[0], 3)
		cx_size = &temp[0]
	}
	return uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(cx_size), 0))) + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(cx_size), 1))) << 8) + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(cx_size), 2))) << 16))
}
func read_is_valid(start *uint8, len_ uint64, end *uint8) int {
	return int(libc.BoolToInt(len_ != 0 && uintptr(unsafe.Pointer(end)) > uintptr(unsafe.Pointer(start)) && len_ <= uint64(int64(uintptr(unsafe.Pointer(end))-uintptr(unsafe.Pointer(start))))))
}
func read_available_partition_size(pbi *VP8D_COMP, token_part_sizes *uint8, fragment_start *uint8, first_fragment_end *uint8, fragment_end *uint8, i int, num_part int) uint {
	var (
		pc                 = &pbi.Common
		partition_size_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(token_part_sizes), i*3))
		partition_size     uint       = 0
		bytes_left          = int64(uintptr(unsafe.Pointer(fragment_end)) - uintptr(unsafe.Pointer(fragment_start)))
	)
	if bytes_left < 0 {
		vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Truncated packet or corrupt partition. No bytes left %d."), int(bytes_left))
	}
	if i < num_part-1 {
		if read_is_valid(partition_size_ptr, 3, first_fragment_end) != 0 {
			partition_size = read_partition_size(pbi, partition_size_ptr)
		} else if pbi.Ec_active != 0 {
			partition_size = uint(bytes_left)
		} else {
			vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Truncated partition size data"))
		}
	} else {
		partition_size = uint(bytes_left)
	}
	if read_is_valid(fragment_start, uint64(partition_size), fragment_end) == 0 {
		if pbi.Ec_active != 0 {
			partition_size = uint(bytes_left)
		} else {
			vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Truncated packet or corrupt partition %d length"), i+1)
		}
	}
	return partition_size
}
func setup_token_decoder(pbi *VP8D_COMP, token_part_sizes *uint8) {
	var (
		bool_decoder  = &pbi.Mbc[0]
		partition_idx uint
		fragment_idx          uint
		num_token_partitions  uint
		first_fragment_end    = (*uint8)(unsafe.Add(unsafe.Pointer(pbi.Fragments.Ptrs[0]), pbi.Fragments.Sizes[0]))
		multi_token_partition = TOKEN_PARTITION(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(&pbi.Mbc[8])), 2))
	)
	if vp8dx_bool_error((*BOOL_DECODER)(unsafe.Pointer(&pbi.Mbc[8]))) == 0 {
		pbi.Common.Multi_token_partition = multi_token_partition
	}
	num_token_partitions = uint(1 << pbi.Common.Multi_token_partition)
	for fragment_idx = 0; fragment_idx < pbi.Fragments.Count; fragment_idx++ {
		var (
			fragment_size = pbi.Fragments.Sizes[fragment_idx]
			fragment_end  = (*uint8)(unsafe.Add(unsafe.Pointer(pbi.Fragments.Ptrs[fragment_idx]), fragment_size))
		)
		if fragment_idx == 0 {
			var ext_first_part_size = int64(uintptr(unsafe.Pointer(token_part_sizes))-uintptr(unsafe.Pointer(pbi.Fragments.Ptrs[0]))) + int64((num_token_partitions-1)*3)
			if fragment_size < uint(ext_first_part_size) {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Corrupted fragment size %d"), fragment_size)
			}
			fragment_size -= uint(ext_first_part_size)
			if fragment_size > 0 {
				pbi.Fragments.Sizes[0] = uint(ext_first_part_size)
				fragment_idx++
				pbi.Fragments.Ptrs[fragment_idx] = (*uint8)(unsafe.Add(unsafe.Pointer(pbi.Fragments.Ptrs[0]), pbi.Fragments.Sizes[0]))
			}
		}
		for fragment_size > 0 {
			var partition_size = int64(read_available_partition_size(pbi, token_part_sizes, pbi.Fragments.Ptrs[fragment_idx], first_fragment_end, fragment_end, int(fragment_idx-1), int(num_token_partitions)))
			pbi.Fragments.Sizes[fragment_idx] = uint(partition_size)
			if fragment_size < uint(partition_size) {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Corrupted fragment size %d"), fragment_size)
			}
			fragment_size -= uint(partition_size)
			if fragment_idx <= num_token_partitions {
			} else {
				// Todo:
				log.Fatal("error")

			}
			if fragment_size > 0 {
				fragment_idx++
				pbi.Fragments.Ptrs[fragment_idx] = (*uint8)(unsafe.Add(unsafe.Pointer(pbi.Fragments.Ptrs[fragment_idx-1]), partition_size))
			}
		}
	}
	pbi.Fragments.Count = num_token_partitions + 1
	for partition_idx = 1; partition_idx < pbi.Fragments.Count; partition_idx++ {
		if vp8dx_start_decode((*BOOL_DECODER)(unsafe.Pointer(bool_decoder)), pbi.Fragments.Ptrs[partition_idx], pbi.Fragments.Sizes[partition_idx], pbi.Decrypt_cb, pbi.Decrypt_state) != 0 {
			vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate bool decoder %d"), partition_idx)
		}
		bool_decoder = (*vp8_reader)(unsafe.Add(unsafe.Pointer(bool_decoder), unsafe.Sizeof(vp8_reader{})*1))
	}
	if pbi.Decoding_thread_count > num_token_partitions-1 {
		pbi.Decoding_thread_count = num_token_partitions - 1
	}
	if int(pbi.Decoding_thread_count) > pbi.Common.Mb_rows-1 {
		if pbi.Common.Mb_rows > 0 {
		} else {
			// Todo:
			log.Fatal("error")

		}
		pbi.Decoding_thread_count = uint(pbi.Common.Mb_rows - 1)
	}
}
func init_frame(pbi *VP8D_COMP) {
	var (
		pc = &pbi.Common
		xd = &pbi.Mb
	)
	if pc.Frame_type == int(KEY_FRAME) {
		libc.MemCpy(unsafe.Pointer(&pc.Fc.Mvc[0]), unsafe.Pointer(&vp8_default_mv_context[0]), int(unsafe.Sizeof([2]MV_CONTEXT{})))
		vp8_init_mbmode_probs(pc)
		vp8_default_coef_probs(pc)
		*(*[2][4]int8)(unsafe.Pointer(&xd.Segment_feature_data[0][0])) = [2][4]int8{}
		xd.Mb_segement_abs_delta = SEGMENT_DELTADATA
		*(*[4]int8)(unsafe.Pointer(&xd.Ref_lf_deltas[0])) = [4]int8{}
		*(*[4]int8)(unsafe.Pointer(&xd.Mode_lf_deltas[0])) = [4]int8{}
		pc.Refresh_golden_frame = 1
		pc.Refresh_alt_ref_frame = 1
		pc.Copy_buffer_to_gf = 0
		pc.Copy_buffer_to_arf = 0
		pc.Ref_frame_sign_bias[GOLDEN_FRAME] = 0
		pc.Ref_frame_sign_bias[ALTREF_FRAME] = 0
	} else {
		if pc.Use_bilinear_mc_filter == 0 {
			xd.Subpixel_predict = func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
				func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
					Vp8SixtapPredict4x4C(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
				}(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
			}
			xd.Subpixel_predict8x4 = func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
				func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
					Vp8SixtapPredict8x4C(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
				}(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
			}
			xd.Subpixel_predict8x8 = func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
				func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
					Vp8SixtapPredict8x8C(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
				}(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
			}
			xd.Subpixel_predict16x16 = func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
				func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
					Vp8SixtapPredict16x16C(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
				}(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
			}
		} else {
			xd.Subpixel_predict = func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
				func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
					Vp8BilinearPredict4x4C(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
				}(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
			}
			xd.Subpixel_predict8x4 = func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
				func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
					Vp8BilinearPredict8x4C(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
				}(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
			}
			xd.Subpixel_predict8x8 = func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
				func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
					Vp8BilinearPredict8x8C(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
				}(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
			}
			xd.Subpixel_predict16x16 = func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
				func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int) {
					Vp8BilinearPredict16x16C(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
				}(src_ptr, src_pixels_per_line, xoffset, yoffset, dst_ptr, dst_pitch)
			}
		}
		if pbi.Decoded_key_frame != 0 && pbi.Ec_enabled != 0 && pbi.Ec_active == 0 {
			pbi.Ec_active = 1
		}
	}
	xd.Left_context = &pc.Left_context
	xd.Mode_info_context = pc.Mi
	xd.Frame_type = pc.Frame_type
	xd.Mode_info_context.Mbmi.Mode = uint8(int8(DC_PRED))
	xd.Mode_info_stride = pc.Mode_info_stride
	xd.Corrupted = 0
	xd.Fullpixel_mask = ^int(0)
	if pc.Full_pixel != 0 {
		xd.Fullpixel_mask = ^int(7)
	}
}
func vp8_decode_frame(pbi *VP8D_COMP) int {
	var (
		bc = &pbi.Mbc[8]
		pc = &pbi.Common
		xd             = &pbi.Mb
		data            = pbi.Fragments.Ptrs[0]
		data_sz              = pbi.Fragments.Sizes[0]
		data_end        = (*uint8)(unsafe.Add(unsafe.Pointer(data), data_sz))
		first_partition_length_in_bytes int64
		i                               int
		j                               int
		k                               int
		l                    int
		mb_feature_data_bits = &vp8_mb_feature_data_bits[0]
		corrupt_tokens       = 0
		prev_independent_partitions      = pbi.Independent_partitions
		yv12_fb_new                     = pbi.Dec_fb_ref[INTRA_FRAME]
	)
	xd.Corrupted = 0
	yv12_fb_new.Corrupted = 0
	if int64(uintptr(unsafe.Pointer(data_end))-uintptr(unsafe.Pointer(data))) < 3 {
		if pbi.Ec_active == 0 {
			vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Truncated packet"))
		}
		pc.Frame_type = int(INTER_FRAME)
		pc.Version = 0
		pc.Show_frame = 1
		first_partition_length_in_bytes = 0
	} else {
		var (
			clear_buffer [10]uint8
			clear        = data
		)
		if pbi.Decrypt_cb != nil {
			var n = int(func() uintptr {
				if (unsafe.Sizeof([10]uint8{})) < uintptr(data_sz) {
					return unsafe.Sizeof([10]uint8{})
				}
				return uintptr(data_sz)
			}())
			pbi.Decrypt_cb(pbi.Decrypt_state, data, &clear_buffer[0], n)
			clear = &clear_buffer[0]
		}
		pc.Frame_type = int(int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 0))) & 1)
		pc.Version = (int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 0))) >> 1) & 7
		pc.Show_frame = (int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 0))) >> 4) & 1
		first_partition_length_in_bytes = int64((int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 0))) | int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 1)))<<8 | int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 2)))<<16) >> 5)
		if pbi.Ec_active == 0 && (uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), first_partition_length_in_bytes)))) > uintptr(unsafe.Pointer(data_end)) || uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), first_partition_length_in_bytes)))) < uintptr(unsafe.Pointer(data))) {
			vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Truncated packet or corrupt partition 0 length"))
		}
		data = (*uint8)(unsafe.Add(unsafe.Pointer(data), 3))
		clear = (*uint8)(unsafe.Add(unsafe.Pointer(clear), 3))
		vp8_setup_version(pc)
		if pc.Frame_type == int(KEY_FRAME) {
			if uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), 3)))) < uintptr(unsafe.Pointer(data_end)) {
				if int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 0))) != 157 || int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 1))) != 1 || int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 2))) != 42 {
					vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_UNSUP_BITSTREAM), libc.CString("Invalid frame sync code"))
				}
			}
			if uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), 6)))) < uintptr(unsafe.Pointer(data_end)) {
				pc.Width = (int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 3))) | int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 4)))<<8) & 0x3FFF
				pc.Horiz_scale = int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 4))) >> 6
				pc.Height = (int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 5))) | int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 6)))<<8) & 0x3FFF
				pc.Vert_scale = int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 6))) >> 6
				data = (*uint8)(unsafe.Add(unsafe.Pointer(data), 7))
			} else if pbi.Ec_active == 0 {
				vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Truncated key frame header"))
			} else {
				data = data_end
			}
		} else {
			libc.MemCpy(unsafe.Pointer(&xd.Pre), unsafe.Pointer(yv12_fb_new), int(unsafe.Sizeof(scale.Yv12BufferConfig{})))
			libc.MemCpy(unsafe.Pointer(&xd.Dst), unsafe.Pointer(yv12_fb_new), int(unsafe.Sizeof(scale.Yv12BufferConfig{})))
		}
	}
	if pbi.Decoded_key_frame == 0 && pc.Frame_type != int(KEY_FRAME) {
		return -1
	}
	init_frame(pbi)
	if vp8dx_start_decode((*BOOL_DECODER)(unsafe.Pointer(bc)), data, uint(int64(uintptr(unsafe.Pointer(data_end))-uintptr(unsafe.Pointer(data)))), pbi.Decrypt_cb, pbi.Decrypt_state) != 0 {
		vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate bool decoder 0"))
	}
	if pc.Frame_type == int(KEY_FRAME) {
		vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)
		pc.Clamp_type = CLAMP_TYPE(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128))
	}
	xd.Segmentation_enabled = uint8(int8(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)))
	if int(xd.Segmentation_enabled) != 0 {
		xd.Update_mb_segmentation_map = uint8(int8(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)))
		xd.Update_mb_segmentation_data = uint8(int8(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)))
		if int(xd.Update_mb_segmentation_data) != 0 {
			xd.Mb_segement_abs_delta = uint8(int8(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)))
			*(*[2][4]int8)(unsafe.Pointer(&xd.Segment_feature_data[0][0])) = [2][4]int8{}
			for i = 0; i < MB_LVL_MAX; i++ {
				for j = 0; j < MAX_MB_SEGMENTS; j++ {
					if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
						xd.Segment_feature_data[i][j] = int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), *(*int)(unsafe.Add(unsafe.Pointer(mb_feature_data_bits), unsafe.Sizeof(int(0))*uintptr(i)))))
						if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
							xd.Segment_feature_data[i][j] = int8(int(-xd.Segment_feature_data[i][j]))
						}
					} else {
						xd.Segment_feature_data[i][j] = 0
					}
				}
			}
		}
		if int(xd.Update_mb_segmentation_map) != 0 {
			libc.MemSet(unsafe.Pointer(&xd.Mb_segment_tree_probs[0]), math.MaxUint8, int(unsafe.Sizeof([3]uint8{})))
			for i = 0; i < MB_FEATURE_TREE_PROBS; i++ {
				if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
					xd.Mb_segment_tree_probs[i] = uint8(int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 8)))
				}
			}
		}
	} else {
		xd.Update_mb_segmentation_map = 0
		xd.Update_mb_segmentation_data = 0
	}
	pc.Filter_type = LOOPFILTERTYPE(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128))
	pc.Filter_level = vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 6)
	pc.Sharpness_level = vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 3)
	xd.Mode_ref_lf_delta_update = 0
	xd.Mode_ref_lf_delta_enabled = uint8(int8(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)))
	if int(xd.Mode_ref_lf_delta_enabled) != 0 {
		xd.Mode_ref_lf_delta_update = uint8(int8(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)))
		if int(xd.Mode_ref_lf_delta_update) != 0 {
			for i = 0; i < MAX_REF_LF_DELTAS; i++ {
				if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
					xd.Ref_lf_deltas[i] = int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 6))
					if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
						xd.Ref_lf_deltas[i] = int8(int(xd.Ref_lf_deltas[i]) * (-1))
					}
				}
			}
			for i = 0; i < MAX_MODE_LF_DELTAS; i++ {
				if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
					xd.Mode_lf_deltas[i] = int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 6))
					if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
						xd.Mode_lf_deltas[i] = int8(int(xd.Mode_lf_deltas[i]) * (-1))
					}
				}
			}
		}
	}
	setup_token_decoder(pbi, (*uint8)(unsafe.Add(unsafe.Pointer(data), first_partition_length_in_bytes)))
	xd.Current_bc = unsafe.Pointer(&pbi.Mbc[0])
	{
		var (
			Q        int
			q_update int
		)
		Q = vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 7)
		pc.Base_qindex = Q
		q_update = 0
		pc.Y1dc_delta_q = get_delta_q(bc, pc.Y1dc_delta_q, &q_update)
		pc.Y2dc_delta_q = get_delta_q(bc, pc.Y2dc_delta_q, &q_update)
		pc.Y2ac_delta_q = get_delta_q(bc, pc.Y2ac_delta_q, &q_update)
		pc.Uvdc_delta_q = get_delta_q(bc, pc.Uvdc_delta_q, &q_update)
		pc.Uvac_delta_q = get_delta_q(bc, pc.Uvac_delta_q, &q_update)
		if q_update != 0 {
			vp8cx_init_de_quantizer(pbi)
		}
		vp8_mb_init_dequantizer(pbi, &pbi.Mb)
	}
	if pc.Frame_type != int(KEY_FRAME) {
		pc.Refresh_golden_frame = vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)
		pc.Refresh_alt_ref_frame = vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)
		pc.Copy_buffer_to_gf = 0
		if pc.Refresh_golden_frame == 0 {
			pc.Copy_buffer_to_gf = vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 2)
		}
		pc.Copy_buffer_to_arf = 0
		if pc.Refresh_alt_ref_frame == 0 {
			pc.Copy_buffer_to_arf = vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 2)
		}
		pc.Ref_frame_sign_bias[GOLDEN_FRAME] = vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)
		pc.Ref_frame_sign_bias[ALTREF_FRAME] = vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)
	}
	pc.Refresh_entropy_probs = vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)
	if pc.Refresh_entropy_probs == 0 {
		libc.MemCpy(unsafe.Pointer(&pc.Lfc), unsafe.Pointer(&pc.Fc), int(unsafe.Sizeof(FRAME_CONTEXT{})))
	}
	pc.Refresh_last_frame = int(libc.BoolToInt(pc.Frame_type == int(KEY_FRAME) || vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0))
	if false {
		var z = (*stdio.File)(unsafe.Pointer(stdio.FOpen("decodestats.stt", "a")))
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(z)), "%6d F:%d,G:%d,A:%d,L:%d,Q:%d\n", pc.Current_video_frame, pc.Frame_type, pc.Refresh_golden_frame, pc.Refresh_alt_ref_frame, pc.Refresh_last_frame, pc.Base_qindex)
		fclose(z)
	}
	{
		pbi.Independent_partitions = 1
		for i = 0; i < BLOCK_TYPES; i++ {
			for j = 0; j < COEF_BANDS; j++ {
				for k = 0; k < PREV_COEF_CONTEXTS; k++ {
					for l = 0; l < ENTROPY_NODES; l++ {
						var p = &pc.Fc.Coef_probs[i][j][k][l]
						if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), int(vp8_coef_update_probs[i][j][k][l])) != 0 {
							*p = uint8(int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 8)))
						}
						if k > 0 && *p != pc.Fc.Coef_probs[i][j][k-1][l] {
							pbi.Independent_partitions = 0
						}
					}
				}
			}
		}
	}
	*(*[400]int16)(unsafe.Pointer(&xd.Qcoeff[0])) = [400]int16{}
	vp8_decode_mode_mvs(pbi)
	libc.MemSet(unsafe.Pointer(pc.Above_context), 0, pc.Mb_cols*int(unsafe.Sizeof(ENTROPY_CONTEXT_PLANES{})))
	pbi.Frame_corrupt_residual = 0
	if util.AtomicLoadAcquire(&pbi.B_multithreaded_rd) != 0 && pc.Multi_token_partition != TOKEN_PARTITION(ONE_PARTITION) {
		var thread uint
		if vp8mt_decode_mb_rows(pbi, xd) != 0 {
			vp8_decoder_remove_threads(pbi)
			pbi.Restart_threads = 1
			vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), nil)
		}
		scale.Vp8Yv12ExtendFrameBorders((*scale.Yv12BufferConfig)(unsafe.Pointer(yv12_fb_new)))
		for thread = 0; thread < pbi.Decoding_thread_count; thread++ {
			corrupt_tokens |= (*(*MB_ROW_DEC)(unsafe.Add(unsafe.Pointer(pbi.Mb_row_di), unsafe.Sizeof(MB_ROW_DEC{})*uintptr(thread)))).Mbd.Corrupted
		}
	} else {
		decode_mb_rows(pbi)
		corrupt_tokens |= xd.Corrupted
	}
	yv12_fb_new.Corrupted = vp8dx_bool_error((*BOOL_DECODER)(unsafe.Pointer(bc)))
	yv12_fb_new.Corrupted |= corrupt_tokens
	if pbi.Decoded_key_frame == 0 {
		if pc.Frame_type == int(KEY_FRAME) && yv12_fb_new.Corrupted == 0 {
			pbi.Decoded_key_frame = 1
		} else {
			vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("A stream must start with a complete key frame"))
		}
	}
	if pc.Refresh_entropy_probs == 0 {
		libc.MemCpy(unsafe.Pointer(&pc.Fc), unsafe.Pointer(&pc.Lfc), int(unsafe.Sizeof(FRAME_CONTEXT{})))
		pbi.Independent_partitions = prev_independent_partitions
	}
	return 0
}
