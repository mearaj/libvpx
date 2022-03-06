package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/mem"
	"github.com/mearaj/libvpx/internal/scale"
	"github.com/mearaj/libvpx/internal/util"
	"github.com/mearaj/libvpx/internal/vpx"
	"math"
	"unsafe"
)

func setup_decoding_thread_data(pbi *VP8D_COMP, xd *MacroBlockd, mbrd *MB_ROW_DEC, count int) {
	var (
		pc = &pbi.Common
		i  int
	)
	for i = 0; i < count; i++ {
		var mbd = &(*(*MB_ROW_DEC)(unsafe.Add(unsafe.Pointer(mbrd), unsafe.Sizeof(MB_ROW_DEC{})*uintptr(i)))).Mbd
		mbd.Subpixel_predict = xd.Subpixel_predict
		mbd.Subpixel_predict8x4 = xd.Subpixel_predict8x4
		mbd.Subpixel_predict8x8 = xd.Subpixel_predict8x8
		mbd.Subpixel_predict16x16 = xd.Subpixel_predict16x16
		mbd.Frame_type = pc.Frame_type
		mbd.Pre = xd.Pre
		mbd.Dst = xd.Dst
		mbd.Segmentation_enabled = xd.Segmentation_enabled
		mbd.Mb_segement_abs_delta = xd.Mb_segement_abs_delta
		libc.MemCpy(unsafe.Pointer(&mbd.Segment_feature_data[0][0]), unsafe.Pointer(&xd.Segment_feature_data[0][0]), int(unsafe.Sizeof([2][4]int8{})))
		libc.MemCpy(unsafe.Pointer(&mbd.Ref_lf_deltas[0]), unsafe.Pointer(&xd.Ref_lf_deltas[0]), int(unsafe.Sizeof([4]int8{})))
		libc.MemCpy(unsafe.Pointer(&mbd.Mode_lf_deltas[0]), unsafe.Pointer(&xd.Mode_lf_deltas[0]), int(unsafe.Sizeof([4]int8{})))
		mbd.Mode_ref_lf_delta_enabled = xd.Mode_ref_lf_delta_enabled
		mbd.Mode_ref_lf_delta_update = xd.Mode_ref_lf_delta_update
		mbd.Current_bc = unsafe.Pointer(&pbi.Mbc[0])
		libc.MemCpy(unsafe.Pointer(&mbd.Dequant_y1_dc[0]), unsafe.Pointer(&xd.Dequant_y1_dc[0]), int(unsafe.Sizeof([16]int16{})))
		libc.MemCpy(unsafe.Pointer(&mbd.Dequant_y1[0]), unsafe.Pointer(&xd.Dequant_y1[0]), int(unsafe.Sizeof([16]int16{})))
		libc.MemCpy(unsafe.Pointer(&mbd.Dequant_y2[0]), unsafe.Pointer(&xd.Dequant_y2[0]), int(unsafe.Sizeof([16]int16{})))
		libc.MemCpy(unsafe.Pointer(&mbd.Dequant_uv[0]), unsafe.Pointer(&xd.Dequant_uv[0]), int(unsafe.Sizeof([16]int16{})))
		mbd.Fullpixel_mask = math.MaxUint32
		if pc.Full_pixel != 0 {
			mbd.Fullpixel_mask = 0xFFFFFFF8
		}
	}
	for i = 0; i < pc.Mb_rows; i++ {
		util.VpxAtomicStoreRelease((*util.VpxAtomicInt)(unsafe.Add(unsafe.Pointer(pbi.Mt_current_mb_col), unsafe.Sizeof(util.VpxAtomicInt{})*uintptr(i))), -1)
	}
}
func mt_decode_macroblock(pbi *VP8D_COMP, xd *MacroBlockd, mb_idx uint) {
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
					Above  *uint8
					yleft       *uint8
					left_stride int
					top_left    uint8
				)
				if i < 4 && pbi.Common.Filter_level != 0 {
					Above = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[0]), b.Offset))
				} else {
					Above = (*uint8)(unsafe.Add(unsafe.Pointer(dst), -dst_stride))
				}
				if i%4 == 0 && pbi.Common.Filter_level != 0 {
					yleft = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_left[0]), i))
					left_stride = 1
				} else {
					yleft = (*uint8)(unsafe.Add(unsafe.Pointer(dst), -1))
					left_stride = dst_stride
				}
				if (i == 4 || i == 8 || i == 12) && pbi.Common.Filter_level != 0 {
					top_left = *((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_left[0]), i))), -1)))
				} else {
					top_left = *(*uint8)(unsafe.Add(unsafe.Pointer(Above), -1))
				}
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
func mt_decode_mb_rows(pbi *VP8D_COMP, xd *MacroBlockd, start_mb_row int) {
	var (
		last_row_current_mb_col *util.VpxAtomicInt
		current_mb_col          *util.VpxAtomicInt
		mb_row int
		pc     = &pbi.Common
		nsync  = pbi.Sync_range
		first_row_no_sync_above            = util.VpxAtomicInt{Value: pc.Mb_cols + nsync}
		num_part                    = int(1 << pbi.Common.Multi_token_partition)
		last_mb_row                               = start_mb_row
		yv12_fb_new     = pbi.Dec_fb_ref[INTRA_FRAME]
		yv12_fb_lst     = pbi.Dec_fb_ref[LAST_FRAME]
		recon_y_stride                         = yv12_fb_new.Y_stride
		recon_uv_stride                         = yv12_fb_new.Uv_stride
		ref_buffer      [4][3]*uint8
		dst_buffer              [3]*uint8
		i                       int
		ref_fb_corrupted        [4]int
	)
	ref_fb_corrupted[INTRA_FRAME] = 0
	for i = 1; i < MAX_REF_FRAMES; i++ {
		var this_fb = pbi.Dec_fb_ref[i]
		ref_buffer[i][0] = (*uint8)(unsafe.Pointer(this_fb.Y_buffer))
		ref_buffer[i][1] = (*uint8)(unsafe.Pointer(this_fb.U_buffer))
		ref_buffer[i][2] = (*uint8)(unsafe.Pointer(this_fb.V_buffer))
		ref_fb_corrupted[i] = this_fb.Corrupted
	}
	dst_buffer[0] = (*uint8)(unsafe.Pointer(yv12_fb_new.Y_buffer))
	dst_buffer[1] = (*uint8)(unsafe.Pointer(yv12_fb_new.U_buffer))
	dst_buffer[2] = (*uint8)(unsafe.Pointer(yv12_fb_new.V_buffer))
	xd.Up_available = int(libc.BoolToInt(start_mb_row != 0))
	xd.Mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(pc.Mi), unsafe.Sizeof(ModeInfo{})*uintptr(pc.Mode_info_stride*start_mb_row)))
	xd.Mode_info_stride = pc.Mode_info_stride
	for mb_row = start_mb_row; mb_row < pc.Mb_rows; mb_row += int(pbi.Decoding_thread_count + 1) {
		var (
			recon_yoffset  int
			recon_uvoffset int
			mb_col         int
			filter_level int
			lfi_n        = &pc.Lf_info
		)
		last_mb_row = mb_row
		xd.Current_bc = unsafe.Pointer(&pbi.Mbc[mb_row%num_part])
		if mb_row > 0 {
			last_row_current_mb_col = (*util.VpxAtomicInt)(unsafe.Add(unsafe.Pointer(pbi.Mt_current_mb_col), unsafe.Sizeof(util.VpxAtomicInt{})*uintptr(mb_row-1)))
		} else {
			last_row_current_mb_col = &first_row_no_sync_above
		}
		current_mb_col = (*util.VpxAtomicInt)(unsafe.Add(unsafe.Pointer(pbi.Mt_current_mb_col), unsafe.Sizeof(util.VpxAtomicInt{})*uintptr(mb_row)))
		recon_yoffset = mb_row * recon_y_stride * 16
		recon_uvoffset = mb_row * recon_uv_stride * 8
		xd.Above_context = pc.Above_context
		*xd.Left_context = ENTROPY_CONTEXT_PLANES{}
		xd.Left_available = 0
		xd.Mb_to_top_edge = -((mb_row * 16) << 3)
		xd.Mb_to_bottom_edge = ((pc.Mb_rows - 1 - mb_row) * 16) << 3
		if pbi.Common.Filter_level != 0 {
			xd.Recon_above[0] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row)))), 0*16))), 32))
			xd.Recon_above[1] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row)))), 0*8))), 16))
			xd.Recon_above[2] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row)))), 0*8))), 16))
			xd.Recon_left[0] = *(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row)))
			xd.Recon_left[1] = *(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row)))
			xd.Recon_left[2] = *(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row)))
			xd.Recon_left_stride[0] = 1
			xd.Recon_left_stride[1] = 1
		} else {
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
		}
		for mb_col = 0; mb_col < pc.Mb_cols; mb_col++ {
			if ((mb_col - 1) % nsync) == 0 {
				util.VpxAtomicStoreRelease(current_mb_col, mb_col-1)
			}
			if mb_row != 0 && (mb_col&(nsync-1)) == 0 {
				vp8_atomic_spin_wait(mb_col, last_row_current_mb_col, nsync)
			}
			xd.Mb_to_left_edge = -((mb_col * 16) << 3)
			xd.Mb_to_right_edge = ((pc.Mb_cols - 1 - mb_col) * 16) << 3
			xd.Dst.Y_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(dst_buffer[0]), recon_yoffset))
			xd.Dst.U_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(dst_buffer[1]), recon_uvoffset))
			xd.Dst.V_buffer = (*uint8)(unsafe.Add(unsafe.Pointer(dst_buffer[2]), recon_uvoffset))
			xd.Corrupted |= ref_fb_corrupted[xd.Mode_info_context.Mbmi.Ref_frame]
			if xd.Corrupted != 0 {
				for ; mb_row < pc.Mb_rows; mb_row += int(pbi.Decoding_thread_count + 1) {
					current_mb_col = (*util.VpxAtomicInt)(unsafe.Add(unsafe.Pointer(pbi.Mt_current_mb_col), unsafe.Sizeof(util.VpxAtomicInt{})*uintptr(mb_row)))
					util.VpxAtomicStoreRelease(current_mb_col, pc.Mb_cols+nsync)
				}
				vpx.InternalError(&xd.Error_info, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Corrupted reference frame"))
			}
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
			mt_decode_macroblock(pbi, xd, 0)
			xd.Left_available = 1
			xd.Corrupted |= vp8dx_bool_error((*BOOL_DECODER)(xd.Current_bc))
			xd.Recon_above[0] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[0]), 16))
			xd.Recon_above[1] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[1]), 8))
			xd.Recon_above[2] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_above[2]), 8))
			if pbi.Common.Filter_level == 0 {
				xd.Recon_left[0] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_left[0]), 16))
				xd.Recon_left[1] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_left[1]), 8))
				xd.Recon_left[2] = (*uint8)(unsafe.Add(unsafe.Pointer(xd.Recon_left[2]), 8))
			}
			if pbi.Common.Filter_level != 0 {
				var (
					skip_lf    = int(libc.BoolToInt(int(xd.Mode_info_context.Mbmi.Mode) != B_PRED && int(xd.Mode_info_context.Mbmi.Mode) != SPLITMV && xd.Mode_info_context.Mbmi.Mb_skip_coeff != 0))
					mode_index = int(lfi_n.Mode_lf_lut[xd.Mode_info_context.Mbmi.Mode])
					seg            = int(xd.Mode_info_context.Mbmi.Segment_id)
					ref_frame      = int(xd.Mode_info_context.Mbmi.Ref_frame)
				)
				filter_level = int(lfi_n.Lvl[seg][ref_frame][mode_index])
				if mb_row != pc.Mb_rows-1 {
					libc.MemCpy(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row+1)))), 32))), mb_col*16))), unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.Y_buffer), recon_y_stride*15))), 16)
					libc.MemCpy(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row+1)))), 16))), mb_col*8))), unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.U_buffer), recon_uv_stride*7))), 8)
					libc.MemCpy(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row+1)))), 16))), mb_col*8))), unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.V_buffer), recon_uv_stride*7))), 8)
				}
				if mb_col != pc.Mb_cols-1 {
					var next = (*ModeInfo)(unsafe.Add(unsafe.Pointer(xd.Mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
					if int(next.Mbmi.Ref_frame) == INTRA_FRAME {
						for i = 0; i < 16; i++ {
							*(*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row)))), i)) = uint8(*(*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.Y_buffer), i*recon_y_stride+15)))
						}
						for i = 0; i < 8; i++ {
							*(*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row)))), i)) = uint8(*(*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.U_buffer), i*recon_uv_stride+7)))
							*(*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row)))), i)) = uint8(*(*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.V_buffer), i*recon_uv_stride+7)))
						}
					}
				}
				if filter_level != 0 {
					if pc.Filter_type == LOOPFILTERTYPE(NORMAL_LOOPFILTER) {
						var (
							lfi        loop_filter_info
							frame_type = pc.Frame_type
							hev_index  = int(lfi_n.Hev_thr_lut[frame_type][filter_level])
						)
						lfi.Mblim = &lfi_n.Mblim[filter_level][0]
						lfi.Blim = &lfi_n.Blim[filter_level][0]
						lfi.Lim = &lfi_n.Lim[filter_level][0]
						lfi.Hev_thr = &lfi_n.Hev_thr[hev_index][0]
						if mb_col > 0 {
							vp8_loop_filter_mbv_sse2((*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.U_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.V_buffer)), recon_y_stride, recon_uv_stride, &lfi)
						}
						if skip_lf == 0 {
							vp8_loop_filter_bv_sse2((*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.U_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.V_buffer)), recon_y_stride, recon_uv_stride, &lfi)
						}
						if mb_row > 0 {
							vp8_loop_filter_mbh_sse2((*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.U_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.V_buffer)), recon_y_stride, recon_uv_stride, &lfi)
						}
						if skip_lf == 0 {
							vp8_loop_filter_bh_sse2((*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.U_buffer)), (*uint8)(unsafe.Pointer(xd.Dst.V_buffer)), recon_y_stride, recon_uv_stride, &lfi)
						}
					} else {
						if mb_col > 0 {
							Vp8LoopFilterSimpleVerticalEdgeC((*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), recon_y_stride, &lfi_n.Mblim[filter_level][0])
						}
						if skip_lf == 0 {
							vp8_loop_filter_bvs_sse2((*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), recon_y_stride, &lfi_n.Blim[filter_level][0])
						}
						if mb_row > 0 {
							Vp8LoopFilterSimpleHorizontalEdgeC((*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), recon_y_stride, &lfi_n.Mblim[filter_level][0])
						}
						if skip_lf == 0 {
							vp8_loop_filter_bhs_sse2((*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), recon_y_stride, &lfi_n.Blim[filter_level][0])
						}
					}
				}
			}
			recon_yoffset += 16
			recon_uvoffset += 8
			xd.Mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(xd.Mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
			xd.Above_context = (*ENTROPY_CONTEXT_PLANES)(unsafe.Add(unsafe.Pointer(xd.Above_context), unsafe.Sizeof(ENTROPY_CONTEXT_PLANES{})*1))
		}
		if pbi.Common.Filter_level != 0 {
			if mb_row != pc.Mb_rows-1 {
				var (
					lasty  = yv12_fb_lst.Y_width + scale.VP8BORDERINPIXELS
					lastuv = (yv12_fb_lst.Y_width >> 1) + (int(scale.VP8BORDERINPIXELS >> 1))
				)
				for i = 0; i < 4; i++ {
					*(*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row+1)))), lasty+i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row+1)))), lasty-1))
					*(*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row+1)))), lastuv+i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row+1)))), lastuv-1))
					*(*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row+1)))), lastuv+i)) = *(*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(mb_row+1)))), lastuv-1))
				}
			}
		} else {
			vp8_extend_mb_row(yv12_fb_new, (*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.Y_buffer), 16)), (*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.U_buffer), 8)), (*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.V_buffer), 8)))
		}
		util.VpxAtomicStoreRelease(current_mb_col, mb_col+nsync)
		xd.Mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(xd.Mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
		xd.Up_available = 1
		xd.Mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(xd.Mode_info_context), unsafe.Sizeof(ModeInfo{})*uintptr(xd.Mode_info_stride*int(pbi.Decoding_thread_count))))
	}
	if last_mb_row+int(pbi.Decoding_thread_count)+1 >= pc.Mb_rows {
		sem_post(&pbi.H_event_end_decoding)
	}
}
func thread_decoding_proc(p_data unsafe.Pointer) unsafe.Pointer {
	var (
		ithread = ((*DECODETHREAD_DATA)(p_data)).Ithread
		pbi     = (*VP8D_COMP)(((*DECODETHREAD_DATA)(p_data)).Ptr1)
		mbrd        = (*MB_ROW_DEC)(((*DECODETHREAD_DATA)(p_data)).Ptr2)
		mb_row_left_context ENTROPY_CONTEXT_PLANES
	)
	for {
		if util.AtomicLoadAcquire(&pbi.B_multithreaded_rd) == 0 {
			break
		}
		if sem_wait((*sem_t)(unsafe.Add(unsafe.Pointer(pbi.H_event_start_decoding), unsafe.Sizeof(sem_t{})*uintptr(ithread)))) == 0 {
			if util.AtomicLoadAcquire(&pbi.B_multithreaded_rd) == 0 {
				break
			} else {
				var xd = &mbrd.Mbd
				xd.Left_context = &mb_row_left_context
				if _setjmp(([1]__jmp_buf_tag)(xd.Error_info.Jmp)) != 0 {
					xd.Error_info.Setjmp = 0
					sem_post(&pbi.H_event_end_decoding)
					continue
				}
				xd.Error_info.Setjmp = 1
				mt_decode_mb_rows(pbi, xd, ithread+1)
			}
		}
	}
	return nil
}
func DecoderCreateThreads(pbi *VP8D_COMP) {
	var (
		core_count = 0
		ithread    uint
	)
	vpx_atomic_init(&pbi.B_multithreaded_rd, 0)
	pbi.Allocated_decoding_thread_count = 0
	if pbi.Max_threads > 8 {
		core_count = 8
	} else {
		core_count = pbi.Max_threads
	}
	if core_count > pbi.Common.Processor_core_count {
		core_count = pbi.Common.Processor_core_count
	}
	if core_count > 1 {
		vpx_atomic_init(&pbi.B_multithreaded_rd, 1)
		pbi.Decoding_thread_count = uint(core_count - 1)
		for {
			pbi.H_decoding_thread = (*pthread_t)(mem.VpxCalloc(uint64(unsafe.Sizeof(pthread_t(0))), uint64(pbi.Decoding_thread_count)))
			if pbi.H_decoding_thread == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->h_decoding_thread)"))
			}
			if true {
				break
			}
		}
		for {
			pbi.H_event_start_decoding = (*sem_t)(mem.VpxCalloc(uint64(unsafe.Sizeof(sem_t{})), uint64(pbi.Decoding_thread_count)))
			if pbi.H_event_start_decoding == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->h_event_start_decoding)"))
			}
			if true {
				break
			}
		}
		for {
			for {
				pbi.Mb_row_di = (*MB_ROW_DEC)(mem.VpxMemAlign(32, uint64(pbi.Decoding_thread_count*uint(unsafe.Sizeof(MB_ROW_DEC{})))))
				if pbi.Mb_row_di == nil {
					vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->mb_row_di)"))
				}
				if true {
					break
				}
			}
			libc.MemSet(unsafe.Pointer(pbi.Mb_row_di), 0, int(pbi.Decoding_thread_count*uint(unsafe.Sizeof(MB_ROW_DEC{}))))
			if true {
				break
			}
		}
		for {
			pbi.De_thread_data = (*DECODETHREAD_DATA)(mem.VpxCalloc(uint64(unsafe.Sizeof(DECODETHREAD_DATA{})), uint64(pbi.Decoding_thread_count)))
			if pbi.De_thread_data == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->de_thread_data)"))
			}
			if true {
				break
			}
		}
		if sem_init(&pbi.H_event_end_decoding, 0, 0) != 0 {
			vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to initialize semaphore"))
		}
		for ithread = 0; ithread < pbi.Decoding_thread_count; ithread++ {
			if sem_init((*sem_t)(unsafe.Add(unsafe.Pointer(pbi.H_event_start_decoding), unsafe.Sizeof(sem_t{})*uintptr(ithread))), 0, 0) != 0 {
				break
			}
			vp8_setup_block_dptrs(&(*(*MB_ROW_DEC)(unsafe.Add(unsafe.Pointer(pbi.Mb_row_di), unsafe.Sizeof(MB_ROW_DEC{})*uintptr(ithread)))).Mbd)
			(*(*DECODETHREAD_DATA)(unsafe.Add(unsafe.Pointer(pbi.De_thread_data), unsafe.Sizeof(DECODETHREAD_DATA{})*uintptr(ithread)))).Ithread = int(ithread)
			(*(*DECODETHREAD_DATA)(unsafe.Add(unsafe.Pointer(pbi.De_thread_data), unsafe.Sizeof(DECODETHREAD_DATA{})*uintptr(ithread)))).Ptr1 = unsafe.Pointer(pbi)
			(*(*DECODETHREAD_DATA)(unsafe.Add(unsafe.Pointer(pbi.De_thread_data), unsafe.Sizeof(DECODETHREAD_DATA{})*uintptr(ithread)))).Ptr2 = unsafe.Pointer((*MB_ROW_DEC)(unsafe.Add(unsafe.Pointer(pbi.Mb_row_di), unsafe.Sizeof(MB_ROW_DEC{})*uintptr(ithread))))
			if pthread_create((*pthread_t)(unsafe.Add(unsafe.Pointer(pbi.H_decoding_thread), unsafe.Sizeof(pthread_t(0))*uintptr(ithread))), nil, thread_decoding_proc, unsafe.Pointer((*DECODETHREAD_DATA)(unsafe.Add(unsafe.Pointer(pbi.De_thread_data), unsafe.Sizeof(DECODETHREAD_DATA{})*uintptr(ithread))))) != 0 {
				sem_destroy((*sem_t)(unsafe.Add(unsafe.Pointer(pbi.H_event_start_decoding), unsafe.Sizeof(sem_t{})*uintptr(ithread))))
				break
			}
		}
		pbi.Allocated_decoding_thread_count = int(ithread)
		if pbi.Allocated_decoding_thread_count != int(pbi.Decoding_thread_count) {
			if pbi.Allocated_decoding_thread_count == 0 {
				sem_destroy(&pbi.H_event_end_decoding)
			}
			vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to create threads"))
		}
	}
}
func vp8mt_de_alloc_temp_buffers(pbi *VP8D_COMP, mb_rows int) {
	var i int
	mem.VpxFree(unsafe.Pointer(pbi.Mt_current_mb_col))
	pbi.Mt_current_mb_col = nil
	if pbi.Mt_yabove_row != nil {
		for i = 0; i < mb_rows; i++ {
			mem.VpxFree(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))))
			*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = nil
		}
		mem.VpxFree(unsafe.Pointer(pbi.Mt_yabove_row))
		pbi.Mt_yabove_row = nil
	}
	if pbi.Mt_uabove_row != nil {
		for i = 0; i < mb_rows; i++ {
			mem.VpxFree(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))))
			*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = nil
		}
		mem.VpxFree(unsafe.Pointer(pbi.Mt_uabove_row))
		pbi.Mt_uabove_row = nil
	}
	if pbi.Mt_vabove_row != nil {
		for i = 0; i < mb_rows; i++ {
			mem.VpxFree(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))))
			*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = nil
		}
		mem.VpxFree(unsafe.Pointer(pbi.Mt_vabove_row))
		pbi.Mt_vabove_row = nil
	}
	if pbi.Mt_yleft_col != nil {
		for i = 0; i < mb_rows; i++ {
			mem.VpxFree(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))))
			*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = nil
		}
		mem.VpxFree(unsafe.Pointer(pbi.Mt_yleft_col))
		pbi.Mt_yleft_col = nil
	}
	if pbi.Mt_uleft_col != nil {
		for i = 0; i < mb_rows; i++ {
			mem.VpxFree(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))))
			*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = nil
		}
		mem.VpxFree(unsafe.Pointer(pbi.Mt_uleft_col))
		pbi.Mt_uleft_col = nil
	}
	if pbi.Mt_vleft_col != nil {
		for i = 0; i < mb_rows; i++ {
			mem.VpxFree(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))))
			*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = nil
		}
		mem.VpxFree(unsafe.Pointer(pbi.Mt_vleft_col))
		pbi.Mt_vleft_col = nil
	}
}
func vp8mt_alloc_temp_buffers(pbi *VP8D_COMP, width int, prev_mb_rows int) {
	var (
		pc = &pbi.Common
		i  int
		uv_width int
	)
	if util.AtomicLoadAcquire(&pbi.B_multithreaded_rd) != 0 {
		vp8mt_de_alloc_temp_buffers(pbi, prev_mb_rows)
		if (width & 15) != 0 {
			width += 16 - (width & 15)
		}
		if width < 640 {
			pbi.Sync_range = 1
		} else if width <= 1280 {
			pbi.Sync_range = 8
		} else if width <= 2560 {
			pbi.Sync_range = 16
		} else {
			pbi.Sync_range = 32
		}
		uv_width = width >> 1
		for {
			pbi.Mt_current_mb_col = (*util.VpxAtomicInt)(mem.VpxMalloc(uint64(pc.Mb_rows * int(unsafe.Sizeof(util.VpxAtomicInt{})))))
			if pbi.Mt_current_mb_col == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate pbi->mt_current_mb_col"))
			}
			if true {
				break
			}
		}
		for i = 0; i < pc.Mb_rows; i++ {
			vpx_atomic_init((*util.VpxAtomicInt)(unsafe.Add(unsafe.Pointer(pbi.Mt_current_mb_col), unsafe.Sizeof(util.VpxAtomicInt{})*uintptr(i))), 0)
		}
		for {
			pbi.Mt_yabove_row = (**uint8)(mem.VpxCalloc(uint64(unsafe.Sizeof((*uint8)(nil))), uint64(pc.Mb_rows)))
			if pbi.Mt_yabove_row == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->mt_yabove_row)"))
			}
			if true {
				break
			}
		}
		for i = 0; i < pc.Mb_rows; i++ {
			for {
				*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = (*uint8)(mem.VpxMemAlign(16, uint64((width+(int(scale.VP8BORDERINPIXELS<<1)))*int(unsafe.Sizeof(uint8(0))))))
				if (*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))) == nil {
					vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate pbi->mt_yabove_row[i]"))
				}
				if true {
					break
				}
			}
			libc.MemSet(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))), 0, (width+(int(scale.VP8BORDERINPIXELS<<1)))*int(unsafe.Sizeof(uint8(0))))
		}
		for {
			pbi.Mt_uabove_row = (**uint8)(mem.VpxCalloc(uint64(unsafe.Sizeof((*uint8)(nil))), uint64(pc.Mb_rows)))
			if pbi.Mt_uabove_row == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->mt_uabove_row)"))
			}
			if true {
				break
			}
		}
		for i = 0; i < pc.Mb_rows; i++ {
			for {
				*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = (*uint8)(mem.VpxMemAlign(16, uint64((uv_width+scale.VP8BORDERINPIXELS)*int(unsafe.Sizeof(uint8(0))))))
				if (*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))) == nil {
					vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate pbi->mt_uabove_row[i]"))
				}
				if true {
					break
				}
			}
			libc.MemSet(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))), 0, (uv_width+scale.VP8BORDERINPIXELS)*int(unsafe.Sizeof(uint8(0))))
		}
		for {
			pbi.Mt_vabove_row = (**uint8)(mem.VpxCalloc(uint64(unsafe.Sizeof((*uint8)(nil))), uint64(pc.Mb_rows)))
			if pbi.Mt_vabove_row == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->mt_vabove_row)"))
			}
			if true {
				break
			}
		}
		for i = 0; i < pc.Mb_rows; i++ {
			for {
				*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = (*uint8)(mem.VpxMemAlign(16, uint64((uv_width+scale.VP8BORDERINPIXELS)*int(unsafe.Sizeof(uint8(0))))))
				if (*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))) == nil {
					vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate pbi->mt_vabove_row[i]"))
				}
				if true {
					break
				}
			}
			libc.MemSet(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))), 0, (uv_width+scale.VP8BORDERINPIXELS)*int(unsafe.Sizeof(uint8(0))))
		}
		for {
			pbi.Mt_yleft_col = (**uint8)(mem.VpxCalloc(uint64(unsafe.Sizeof((*uint8)(nil))), uint64(pc.Mb_rows)))
			if pbi.Mt_yleft_col == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->mt_yleft_col)"))
			}
			if true {
				break
			}
		}
		for i = 0; i < pc.Mb_rows; i++ {
			for {
				*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = (*uint8)(mem.VpxCalloc(uint64(unsafe.Sizeof(uint8(0))*16), 1))
				if (*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))) == nil {
					vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate pbi->mt_yleft_col[i]"))
				}
				if true {
					break
				}
			}
		}
		for {
			pbi.Mt_uleft_col = (**uint8)(mem.VpxCalloc(uint64(unsafe.Sizeof((*uint8)(nil))), uint64(pc.Mb_rows)))
			if pbi.Mt_uleft_col == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->mt_uleft_col)"))
			}
			if true {
				break
			}
		}
		for i = 0; i < pc.Mb_rows; i++ {
			for {
				*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = (*uint8)(mem.VpxCalloc(uint64(unsafe.Sizeof(uint8(0))*8), 1))
				if (*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))) == nil {
					vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate pbi->mt_uleft_col[i]"))
				}
				if true {
					break
				}
			}
		}
		for {
			pbi.Mt_vleft_col = (**uint8)(mem.VpxCalloc(uint64(unsafe.Sizeof((*uint8)(nil))), uint64(pc.Mb_rows)))
			if pbi.Mt_vleft_col == nil {
				vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate (pbi->mt_vleft_col)"))
			}
			if true {
				break
			}
		}
		for i = 0; i < pc.Mb_rows; i++ {
			for {
				*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i))) = (*uint8)(mem.VpxCalloc(uint64(unsafe.Sizeof(uint8(0))*8), 1))
				if (*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(i)))) == nil {
					vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate pbi->mt_vleft_col[i]"))
				}
				if true {
					break
				}
			}
		}
	}
}
func vp8_decoder_remove_threads(pbi *VP8D_COMP) {
	if util.AtomicLoadAcquire(&pbi.B_multithreaded_rd) != 0 {
		var i int
		util.VpxAtomicStoreRelease(&pbi.B_multithreaded_rd, 0)
		for i = 0; i < pbi.Allocated_decoding_thread_count; i++ {
			sem_post((*sem_t)(unsafe.Add(unsafe.Pointer(pbi.H_event_start_decoding), unsafe.Sizeof(sem_t{})*uintptr(i))))
			pthread_join(*(*pthread_t)(unsafe.Add(unsafe.Pointer(pbi.H_decoding_thread), unsafe.Sizeof(pthread_t(0))*uintptr(i))), nil)
		}
		for i = 0; i < pbi.Allocated_decoding_thread_count; i++ {
			sem_destroy((*sem_t)(unsafe.Add(unsafe.Pointer(pbi.H_event_start_decoding), unsafe.Sizeof(sem_t{})*uintptr(i))))
		}
		if pbi.Allocated_decoding_thread_count != 0 {
			sem_destroy(&pbi.H_event_end_decoding)
		}
		mem.VpxFree(unsafe.Pointer(pbi.H_decoding_thread))
		pbi.H_decoding_thread = nil
		mem.VpxFree(unsafe.Pointer(pbi.H_event_start_decoding))
		pbi.H_event_start_decoding = nil
		mem.VpxFree(unsafe.Pointer(pbi.Mb_row_di))
		pbi.Mb_row_di = nil
		mem.VpxFree(unsafe.Pointer(pbi.De_thread_data))
		pbi.De_thread_data = nil
		vp8mt_de_alloc_temp_buffers(pbi, pbi.Common.Mb_rows)
	}
}
func vp8mt_decode_mb_rows(pbi *VP8D_COMP, xd *MacroBlockd) int {
	var (
		pc = &pbi.Common
		i  uint
		j            int
		filter_level = pc.Filter_level
		yv12_fb_new  = pbi.Dec_fb_ref[INTRA_FRAME]
	)
	if filter_level != 0 {
		libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*0))), scale.VP8BORDERINPIXELS))), -1), math.MaxInt8, yv12_fb_new.Y_width+5)
		libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*0))), int(scale.VP8BORDERINPIXELS>>1)))), -1), math.MaxInt8, (yv12_fb_new.Y_width>>1)+5)
		libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*0))), int(scale.VP8BORDERINPIXELS>>1)))), -1), math.MaxInt8, (yv12_fb_new.Y_width>>1)+5)
		for j = 1; j < pc.Mb_rows; j++ {
			libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(j)))), scale.VP8BORDERINPIXELS))), -1), 129, 1)
			libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(j)))), int(scale.VP8BORDERINPIXELS>>1)))), -1), 129, 1)
			libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vabove_row), unsafe.Sizeof((*uint8)(nil))*uintptr(j)))), int(scale.VP8BORDERINPIXELS>>1)))), -1), 129, 1)
		}
		for j = 0; j < pc.Mb_rows; j++ {
			libc.MemSet(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_yleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(j)))), 129, 16)
			libc.MemSet(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_uleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(j)))), 129, 8)
			libc.MemSet(unsafe.Pointer(*(**uint8)(unsafe.Add(unsafe.Pointer(pbi.Mt_vleft_col), unsafe.Sizeof((*uint8)(nil))*uintptr(j)))), 129, 8)
		}
		vp8_loop_filter_frame_init(pc, &pbi.Mb, filter_level)
	} else {
		vp8_setup_intra_recon_top_line(yv12_fb_new)
	}
	setup_decoding_thread_data(pbi, xd, pbi.Mb_row_di, int(pbi.Decoding_thread_count))
	for i = 0; i < pbi.Decoding_thread_count; i++ {
		sem_post((*sem_t)(unsafe.Add(unsafe.Pointer(pbi.H_event_start_decoding), unsafe.Sizeof(sem_t{})*uintptr(i))))
	}
	if _setjmp(([1]__jmp_buf_tag)(xd.Error_info.Jmp)) != 0 {
		xd.Error_info.Setjmp = 0
		xd.Corrupted = 1
		for i = 0; i < pbi.Decoding_thread_count; i++ {
			sem_wait(&pbi.H_event_end_decoding)
		}
		return -1
	}
	xd.Error_info.Setjmp = 1
	mt_decode_mb_rows(pbi, xd, 0)
	for i = 0; i < pbi.Decoding_thread_count+1; i++ {
		sem_wait(&pbi.H_event_end_decoding)
	}
	return 0
}
