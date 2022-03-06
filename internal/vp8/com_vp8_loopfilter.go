package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/scale"
	"unsafe"
)

func lf_init_lut(lfi *loop_filter_info_n) {
	var filt_lvl int
	for filt_lvl = 0; filt_lvl <= MAX_LOOP_FILTER; filt_lvl++ {
		if filt_lvl >= 40 {
			lfi.Hev_thr_lut[KEY_FRAME][filt_lvl] = 2
			lfi.Hev_thr_lut[INTER_FRAME][filt_lvl] = 3
		} else if filt_lvl >= 20 {
			lfi.Hev_thr_lut[KEY_FRAME][filt_lvl] = 1
			lfi.Hev_thr_lut[INTER_FRAME][filt_lvl] = 2
		} else if filt_lvl >= 15 {
			lfi.Hev_thr_lut[KEY_FRAME][filt_lvl] = 1
			lfi.Hev_thr_lut[INTER_FRAME][filt_lvl] = 1
		} else {
			lfi.Hev_thr_lut[KEY_FRAME][filt_lvl] = 0
			lfi.Hev_thr_lut[INTER_FRAME][filt_lvl] = 0
		}
	}
	lfi.Mode_lf_lut[DC_PRED] = 1
	lfi.Mode_lf_lut[V_PRED] = 1
	lfi.Mode_lf_lut[H_PRED] = 1
	lfi.Mode_lf_lut[TM_PRED] = 1
	lfi.Mode_lf_lut[B_PRED] = 0
	lfi.Mode_lf_lut[ZEROMV] = 1
	lfi.Mode_lf_lut[NEARESTMV] = 2
	lfi.Mode_lf_lut[NEARMV] = 2
	lfi.Mode_lf_lut[NEWMV] = 2
	lfi.Mode_lf_lut[SPLITMV] = 3
}
func vp8_loop_filter_update_sharpness(lfi *loop_filter_info_n, sharpness_lvl int) {
	var i int
	for i = 0; i <= MAX_LOOP_FILTER; i++ {
		var (
			filt_lvl           = i
			block_inside_limit = 0
		)
		block_inside_limit = filt_lvl >> int(libc.BoolToInt(sharpness_lvl > 0))
		block_inside_limit = block_inside_limit >> int(libc.BoolToInt(sharpness_lvl > 4))
		if sharpness_lvl > 0 {
			if block_inside_limit > (9 - sharpness_lvl) {
				block_inside_limit = 9 - sharpness_lvl
			}
		}
		if block_inside_limit < 1 {
			block_inside_limit = 1
		}
		libc.MemSet(unsafe.Pointer(&lfi.Lim[i][0]), byte(int8(block_inside_limit)), SIMD_WIDTH)
		libc.MemSet(unsafe.Pointer(&lfi.Blim[i][0]), byte(int8(filt_lvl*2+block_inside_limit)), SIMD_WIDTH)
		libc.MemSet(unsafe.Pointer(&lfi.Mblim[i][0]), byte(int8((filt_lvl+2)*2+block_inside_limit)), SIMD_WIDTH)
	}
}
func vp8_loop_filter_init(cm *VP8Common) {
	var (
		lfi = &cm.Lf_info
		i   int
	)
	vp8_loop_filter_update_sharpness(lfi, cm.Sharpness_level)
	cm.Last_sharpness_level = cm.Sharpness_level
	lf_init_lut(lfi)
	for i = 0; i < 4; i++ {
		libc.MemSet(unsafe.Pointer(&lfi.Hev_thr[i][0]), byte(int8(i)), SIMD_WIDTH)
	}
}
func vp8_loop_filter_frame_init(cm *VP8Common, mbd *MacroBlockd, default_filt_lvl int) {
	var (
		seg  int
		ref  int
		mode int
		lfi  = &cm.Lf_info
	)
	if cm.Last_sharpness_level != cm.Sharpness_level {
		vp8_loop_filter_update_sharpness(lfi, cm.Sharpness_level)
		cm.Last_sharpness_level = cm.Sharpness_level
	}
	for seg = 0; seg < MAX_MB_SEGMENTS; seg++ {
		var (
			lvl_seg = default_filt_lvl
			lvl_ref int
			lvl_mode int
		)
		if int(mbd.Segmentation_enabled) != 0 {
			if int(mbd.Mb_segement_abs_delta) == SEGMENT_ABSDATA {
				lvl_seg = int(mbd.Segment_feature_data[MB_LVL_ALT_LF][seg])
			} else {
				lvl_seg += int(mbd.Segment_feature_data[MB_LVL_ALT_LF][seg])
			}
			if lvl_seg > 0 {
				if lvl_seg > 63 {
					lvl_seg = 63
				} else {
					lvl_seg = lvl_seg
				}
			} else {
				lvl_seg = 0
			}
		}
		if int(mbd.Mode_ref_lf_delta_enabled) == 0 {
			libc.MemSet(unsafe.Pointer(&lfi.Lvl[seg][0][0]), byte(int8(lvl_seg)), 4*4)
			continue
		}
		ref = INTRA_FRAME
		lvl_ref = lvl_seg + int(mbd.Ref_lf_deltas[ref])
		mode = 0
		lvl_mode = lvl_ref + int(mbd.Mode_lf_deltas[mode])
		if lvl_mode > 0 {
			if lvl_mode > 63 {
				lvl_mode = 63
			} else {
				lvl_mode = lvl_mode
			}
		} else {
			lvl_mode = 0
		}
		lfi.Lvl[seg][ref][mode] = uint8(int8(lvl_mode))
		mode = 1
		if lvl_ref > 0 {
			if lvl_ref > 63 {
				lvl_mode = 63
			} else {
				lvl_mode = lvl_ref
			}
		} else {
			lvl_mode = 0
		}
		lfi.Lvl[seg][ref][mode] = uint8(int8(lvl_mode))
		for ref = 1; ref < MAX_REF_FRAMES; ref++ {
			lvl_ref = lvl_seg + int(mbd.Ref_lf_deltas[ref])
			for mode = 1; mode < 4; mode++ {
				lvl_mode = lvl_ref + int(mbd.Mode_lf_deltas[mode])
				if lvl_mode > 0 {
					if lvl_mode > 63 {
						lvl_mode = 63
					} else {
						lvl_mode = lvl_mode
					}
				} else {
					lvl_mode = 0
				}
				lfi.Lvl[seg][ref][mode] = uint8(int8(lvl_mode))
			}
		}
	}
}
func vp8_loop_filter_row_normal(cm *VP8Common, mode_info_context *ModeInfo, mb_row int, post_ystride int, post_uvstride int, y_ptr *uint8, u_ptr *uint8, v_ptr *uint8) {
	var (
		mb_col       int
		filter_level int
		lfi_n        = &cm.Lf_info
		lfi          loop_filter_info
		frame_type = cm.Frame_type
	)
	for mb_col = 0; mb_col < cm.Mb_cols; mb_col++ {
		var (
			skip_lf    = int(libc.BoolToInt(int(mode_info_context.Mbmi.Mode) != B_PRED && int(mode_info_context.Mbmi.Mode) != SPLITMV && mode_info_context.Mbmi.Mb_skip_coeff != 0))
			mode_index = int(lfi_n.Mode_lf_lut[mode_info_context.Mbmi.Mode])
			seg            = int(mode_info_context.Mbmi.Segment_id)
			ref_frame      = int(mode_info_context.Mbmi.Ref_frame)
		)
		filter_level = int(lfi_n.Lvl[seg][ref_frame][mode_index])
		if filter_level != 0 {
			var hev_index = int(lfi_n.Hev_thr_lut[frame_type][filter_level])
			lfi.Mblim = &lfi_n.Mblim[filter_level][0]
			lfi.Blim = &lfi_n.Blim[filter_level][0]
			lfi.Lim = &lfi_n.Lim[filter_level][0]
			lfi.Hev_thr = &lfi_n.Hev_thr[hev_index][0]
			if mb_col > 0 {
				vp8_loop_filter_mbv_sse2(y_ptr, u_ptr, v_ptr, post_ystride, post_uvstride, &lfi)
			}
			if skip_lf == 0 {
				vp8_loop_filter_bv_sse2(y_ptr, u_ptr, v_ptr, post_ystride, post_uvstride, &lfi)
			}
			if mb_row > 0 {
				vp8_loop_filter_mbh_sse2(y_ptr, u_ptr, v_ptr, post_ystride, post_uvstride, &lfi)
			}
			if skip_lf == 0 {
				vp8_loop_filter_bh_sse2(y_ptr, u_ptr, v_ptr, post_ystride, post_uvstride, &lfi)
			}
		}
		y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 16))
		u_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), 8))
		v_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), 8))
		mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
	}
}
func vp8_loop_filter_row_simple(cm *VP8Common, mode_info_context *ModeInfo, mb_row int, post_ystride int, y_ptr *uint8) {
	var (
		mb_col       int
		filter_level int
		lfi_n        = &cm.Lf_info
	)
	for mb_col = 0; mb_col < cm.Mb_cols; mb_col++ {
		var (
			skip_lf    = int(libc.BoolToInt(int(mode_info_context.Mbmi.Mode) != B_PRED && int(mode_info_context.Mbmi.Mode) != SPLITMV && mode_info_context.Mbmi.Mb_skip_coeff != 0))
			mode_index = int(lfi_n.Mode_lf_lut[mode_info_context.Mbmi.Mode])
			seg            = int(mode_info_context.Mbmi.Segment_id)
			ref_frame      = int(mode_info_context.Mbmi.Ref_frame)
		)
		filter_level = int(lfi_n.Lvl[seg][ref_frame][mode_index])
		if filter_level != 0 {
			if mb_col > 0 {
				Vp8LoopFilterSimpleVerticalEdgeC(y_ptr, post_ystride, &lfi_n.Mblim[filter_level][0])
			}
			if skip_lf == 0 {
				vp8_loop_filter_bvs_sse2(y_ptr, post_ystride, &lfi_n.Blim[filter_level][0])
			}
			if mb_row > 0 {
				Vp8LoopFilterSimpleHorizontalEdgeC(y_ptr, post_ystride, &lfi_n.Mblim[filter_level][0])
			}
			if skip_lf == 0 {
				vp8_loop_filter_bhs_sse2(y_ptr, post_ystride, &lfi_n.Blim[filter_level][0])
			}
		}
		y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 16))
		mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
	}
}
func vp8_loop_filter_frame(cm *VP8Common, mbd *MacroBlockd, frame_type int) {
	var (
		post  = cm.Frame_to_show
		lfi_n = &cm.Lf_info
		lfi   loop_filter_info
		mb_row            int
		mb_col  int
		mb_rows = cm.Mb_rows
		mb_cols = cm.Mb_cols
		filter_level int
		y_ptr             *uint8
		u_ptr             *uint8
		v_ptr             *uint8
		mode_info_context = cm.Mi
		post_y_stride     = post.Y_stride
		post_uv_stride              = post.Uv_stride
	)
	vp8_loop_filter_frame_init(cm, mbd, cm.Filter_level)
	y_ptr = (*uint8)(unsafe.Pointer(post.Y_buffer))
	u_ptr = (*uint8)(unsafe.Pointer(post.U_buffer))
	v_ptr = (*uint8)(unsafe.Pointer(post.V_buffer))
	if cm.Filter_type == LOOPFILTERTYPE(NORMAL_LOOPFILTER) {
		for mb_row = 0; mb_row < mb_rows; mb_row++ {
			for mb_col = 0; mb_col < mb_cols; mb_col++ {
				var (
					skip_lf    = int(libc.BoolToInt(int(mode_info_context.Mbmi.Mode) != B_PRED && int(mode_info_context.Mbmi.Mode) != SPLITMV && mode_info_context.Mbmi.Mb_skip_coeff != 0))
					mode_index = int(lfi_n.Mode_lf_lut[mode_info_context.Mbmi.Mode])
					seg            = int(mode_info_context.Mbmi.Segment_id)
					ref_frame      = int(mode_info_context.Mbmi.Ref_frame)
				)
				filter_level = int(lfi_n.Lvl[seg][ref_frame][mode_index])
				if filter_level != 0 {
					var hev_index = int(lfi_n.Hev_thr_lut[frame_type][filter_level])
					lfi.Mblim = &lfi_n.Mblim[filter_level][0]
					lfi.Blim = &lfi_n.Blim[filter_level][0]
					lfi.Lim = &lfi_n.Lim[filter_level][0]
					lfi.Hev_thr = &lfi_n.Hev_thr[hev_index][0]
					if mb_col > 0 {
						vp8_loop_filter_mbv_sse2(y_ptr, u_ptr, v_ptr, post_y_stride, post_uv_stride, &lfi)
					}
					if skip_lf == 0 {
						vp8_loop_filter_bv_sse2(y_ptr, u_ptr, v_ptr, post_y_stride, post_uv_stride, &lfi)
					}
					if mb_row > 0 {
						vp8_loop_filter_mbh_sse2(y_ptr, u_ptr, v_ptr, post_y_stride, post_uv_stride, &lfi)
					}
					if skip_lf == 0 {
						vp8_loop_filter_bh_sse2(y_ptr, u_ptr, v_ptr, post_y_stride, post_uv_stride, &lfi)
					}
				}
				y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 16))
				u_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), 8))
				v_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), 8))
				mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
			}
			y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), post_y_stride*16-post.Y_width))
			u_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), post_uv_stride*8-post.Uv_width))
			v_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), post_uv_stride*8-post.Uv_width))
			mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
		}
	} else {
		for mb_row = 0; mb_row < mb_rows; mb_row++ {
			for mb_col = 0; mb_col < mb_cols; mb_col++ {
				var (
					skip_lf    = int(libc.BoolToInt(int(mode_info_context.Mbmi.Mode) != B_PRED && int(mode_info_context.Mbmi.Mode) != SPLITMV && mode_info_context.Mbmi.Mb_skip_coeff != 0))
					mode_index = int(lfi_n.Mode_lf_lut[mode_info_context.Mbmi.Mode])
					seg            = int(mode_info_context.Mbmi.Segment_id)
					ref_frame      = int(mode_info_context.Mbmi.Ref_frame)
				)
				filter_level = int(lfi_n.Lvl[seg][ref_frame][mode_index])
				if filter_level != 0 {
					var (
						mblim = &lfi_n.Mblim[filter_level][0]
						blim  = &lfi_n.Blim[filter_level][0]
					)
					if mb_col > 0 {
						Vp8LoopFilterSimpleVerticalEdgeC(y_ptr, post_y_stride, mblim)
					}
					if skip_lf == 0 {
						vp8_loop_filter_bvs_sse2(y_ptr, post_y_stride, blim)
					}
					if mb_row > 0 {
						Vp8LoopFilterSimpleHorizontalEdgeC(y_ptr, post_y_stride, mblim)
					}
					if skip_lf == 0 {
						vp8_loop_filter_bhs_sse2(y_ptr, post_y_stride, blim)
					}
				}
				y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 16))
				u_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), 8))
				v_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), 8))
				mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
			}
			y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), post_y_stride*16-post.Y_width))
			u_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), post_uv_stride*8-post.Uv_width))
			v_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), post_uv_stride*8-post.Uv_width))
			mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
		}
	}
}
func vp8_loop_filter_frame_yonly(cm *VP8Common, mbd *MacroBlockd, default_filt_lvl int) {
	var (
		post  = cm.Frame_to_show
		y_ptr *uint8
		mb_row            int
		mb_col int
		lfi_n  = &cm.Lf_info
		lfi    loop_filter_info
		filter_level      int
		frame_type        = cm.Frame_type
		mode_info_context = cm.Mi
	)
	vp8_loop_filter_frame_init(cm, mbd, default_filt_lvl)
	y_ptr = (*uint8)(unsafe.Pointer(post.Y_buffer))
	for mb_row = 0; mb_row < cm.Mb_rows; mb_row++ {
		for mb_col = 0; mb_col < cm.Mb_cols; mb_col++ {
			var (
				skip_lf    = int(libc.BoolToInt(int(mode_info_context.Mbmi.Mode) != B_PRED && int(mode_info_context.Mbmi.Mode) != SPLITMV && mode_info_context.Mbmi.Mb_skip_coeff != 0))
				mode_index = int(lfi_n.Mode_lf_lut[mode_info_context.Mbmi.Mode])
				seg            = int(mode_info_context.Mbmi.Segment_id)
				ref_frame      = int(mode_info_context.Mbmi.Ref_frame)
			)
			filter_level = int(lfi_n.Lvl[seg][ref_frame][mode_index])
			if filter_level != 0 {
				if cm.Filter_type == LOOPFILTERTYPE(NORMAL_LOOPFILTER) {
					var hev_index = int(lfi_n.Hev_thr_lut[frame_type][filter_level])
					lfi.Mblim = &lfi_n.Mblim[filter_level][0]
					lfi.Blim = &lfi_n.Blim[filter_level][0]
					lfi.Lim = &lfi_n.Lim[filter_level][0]
					lfi.Hev_thr = &lfi_n.Hev_thr[hev_index][0]
					if mb_col > 0 {
						vp8_loop_filter_mbv_sse2(y_ptr, nil, nil, post.Y_stride, 0, &lfi)
					}
					if skip_lf == 0 {
						vp8_loop_filter_bv_sse2(y_ptr, nil, nil, post.Y_stride, 0, &lfi)
					}
					if mb_row > 0 {
						vp8_loop_filter_mbh_sse2(y_ptr, nil, nil, post.Y_stride, 0, &lfi)
					}
					if skip_lf == 0 {
						vp8_loop_filter_bh_sse2(y_ptr, nil, nil, post.Y_stride, 0, &lfi)
					}
				} else {
					if mb_col > 0 {
						Vp8LoopFilterSimpleVerticalEdgeC(y_ptr, post.Y_stride, &lfi_n.Mblim[filter_level][0])
					}
					if skip_lf == 0 {
						vp8_loop_filter_bvs_sse2(y_ptr, post.Y_stride, &lfi_n.Blim[filter_level][0])
					}
					if mb_row > 0 {
						Vp8LoopFilterSimpleHorizontalEdgeC(y_ptr, post.Y_stride, &lfi_n.Mblim[filter_level][0])
					}
					if skip_lf == 0 {
						vp8_loop_filter_bhs_sse2(y_ptr, post.Y_stride, &lfi_n.Blim[filter_level][0])
					}
				}
			}
			y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 16))
			mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
		}
		y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), post.Y_stride*16-post.Y_width))
		mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
	}
}
func vp8_loop_filter_partial_frame(cm *VP8Common, mbd *MacroBlockd, default_filt_lvl int) {
	var (
		post  = cm.Frame_to_show
		y_ptr *uint8
		mb_row            int
		mb_col  int
		mb_cols = post.Y_width >> 4
		mb_rows = post.Y_height >> 4
		linestocopy int
		lfi_n       = &cm.Lf_info
		lfi         loop_filter_info
		filter_level      int
		frame_type        = cm.Frame_type
		mode_info_context *ModeInfo
	)
	vp8_loop_filter_frame_init(cm, mbd, default_filt_lvl)
	linestocopy = mb_rows / PARTIAL_FRAME_FRACTION
	if linestocopy != 0 {
		linestocopy = linestocopy << 4
	} else {
		linestocopy = 16
	}
	y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(post.Y_buffer), ((post.Y_height>>5)*16)*post.Y_stride))
	mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(cm.Mi), unsafe.Sizeof(ModeInfo{})*uintptr((post.Y_height>>5)*(mb_cols+1))))
	for mb_row = 0; mb_row < (linestocopy >> 4); mb_row++ {
		for mb_col = 0; mb_col < mb_cols; mb_col++ {
			var (
				skip_lf    = int(libc.BoolToInt(int(mode_info_context.Mbmi.Mode) != B_PRED && int(mode_info_context.Mbmi.Mode) != SPLITMV && mode_info_context.Mbmi.Mb_skip_coeff != 0))
				mode_index = int(lfi_n.Mode_lf_lut[mode_info_context.Mbmi.Mode])
				seg            = int(mode_info_context.Mbmi.Segment_id)
				ref_frame      = int(mode_info_context.Mbmi.Ref_frame)
			)
			filter_level = int(lfi_n.Lvl[seg][ref_frame][mode_index])
			if filter_level != 0 {
				if cm.Filter_type == LOOPFILTERTYPE(NORMAL_LOOPFILTER) {
					var hev_index = int(lfi_n.Hev_thr_lut[frame_type][filter_level])
					lfi.Mblim = &lfi_n.Mblim[filter_level][0]
					lfi.Blim = &lfi_n.Blim[filter_level][0]
					lfi.Lim = &lfi_n.Lim[filter_level][0]
					lfi.Hev_thr = &lfi_n.Hev_thr[hev_index][0]
					if mb_col > 0 {
						vp8_loop_filter_mbv_sse2(y_ptr, nil, nil, post.Y_stride, 0, &lfi)
					}
					if skip_lf == 0 {
						vp8_loop_filter_bv_sse2(y_ptr, nil, nil, post.Y_stride, 0, &lfi)
					}
					vp8_loop_filter_mbh_sse2(y_ptr, nil, nil, post.Y_stride, 0, &lfi)
					if skip_lf == 0 {
						vp8_loop_filter_bh_sse2(y_ptr, nil, nil, post.Y_stride, 0, &lfi)
					}
				} else {
					if mb_col > 0 {
						Vp8LoopFilterSimpleVerticalEdgeC(y_ptr, post.Y_stride, &lfi_n.Mblim[filter_level][0])
					}
					if skip_lf == 0 {
						vp8_loop_filter_bvs_sse2(y_ptr, post.Y_stride, &lfi_n.Blim[filter_level][0])
					}
					Vp8LoopFilterSimpleHorizontalEdgeC(y_ptr, post.Y_stride, &lfi_n.Mblim[filter_level][0])
					if skip_lf == 0 {
						vp8_loop_filter_bhs_sse2(y_ptr, post.Y_stride, &lfi_n.Blim[filter_level][0])
					}
				}
			}
			y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 16))
			mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
		}
		y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), post.Y_stride*16-post.Y_width))
		mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
	}
}
