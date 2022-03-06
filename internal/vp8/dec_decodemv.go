package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func read_bmode(bc *vp8_reader, p *uint8) B_PREDICTION_MODE {
	var i int = vp8_treed_read(bc, []int8(vp8_bmode_tree), p)
	return B_PREDICTION_MODE(i)
}
func read_ymode(bc *vp8_reader, p *uint8) int {
	var i int = vp8_treed_read(bc, []int8(vp8_ymode_tree), p)
	return int(i)
}
func read_kf_ymode(bc *vp8_reader, p *uint8) int {
	var i int = vp8_treed_read(bc, []int8(vp8_kf_ymode_tree), p)
	return int(i)
}
func read_uv_mode(bc *vp8_reader, p *uint8) int {
	var i int = vp8_treed_read(bc, []int8(vp8_uv_mode_tree), p)
	return int(i)
}
func read_kf_modes(pbi *VP8D_COMP, mi *ModeInfo) {
	var (
		bc  *vp8_reader = &pbi.Mbc[8]
		mis int         = pbi.Common.Mode_info_stride
	)
	mi.Mbmi.Ref_frame = uint8(int8(INTRA_FRAME))
	mi.Mbmi.Mode = uint8(int8(read_kf_ymode(bc, &vp8_kf_ymode_prob[0])))
	if int(mi.Mbmi.Mode) == B_PRED {
		var i int = 0
		mi.Mbmi.Is_4x4 = 1
		for {
			{
				var (
					A B_PREDICTION_MODE = above_block_mode(mi, i, mis)
					L B_PREDICTION_MODE = left_block_mode(mi, i)
				)
				mi.Bmi[i].As_mode = read_bmode(bc, &vp8_kf_bmode_prob[A][L][0])
			}
			if func() int {
				p := &i
				*p++
				return *p
			}() >= 16 {
				break
			}
		}
	}
	mi.Mbmi.Uv_mode = uint8(int8(read_uv_mode(bc, &vp8_kf_uv_mode_prob[0])))
}
func read_mvcomponent(r *vp8_reader, mvc *MV_CONTEXT) int {
	var (
		p *uint8 = (*uint8)(unsafe.Pointer(mvc))
		x int    = 0
	)
	if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(r)), int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), mvpis_short)))) != 0 {
		var i int = 0
		for {
			x += vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(r)), int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), MVPbits+i)))) << i
			if func() int {
				p := &i
				*p++
				return *p
			}() >= 3 {
				break
			}
		}
		i = mvlong_width - 1
		for {
			x += vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(r)), int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), MVPbits+i)))) << i
			if func() int {
				p := &i
				*p--
				return *p
			}() <= 3 {
				break
			}
		}
		if (x&0xFFF0) == 0 || vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(r)), int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), MVPbits+3)))) != 0 {
			x += 8
		}
	} else {
		x = vp8_treed_read(r, []int8(vp8_small_mvtree), (*uint8)(unsafe.Add(unsafe.Pointer(p), MVPshort)))
	}
	if x != 0 && vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(r)), int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), MVPsign)))) != 0 {
		x = -x
	}
	return x
}
func read_mv(r *vp8_reader, mv *MV, mvc *MV_CONTEXT) {
	mv.Row = int16(read_mvcomponent(r, mvc) * 2)
	mv.Col = int16(read_mvcomponent(r, func() *MV_CONTEXT {
		p := &mvc
		*p = (*MV_CONTEXT)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(MV_CONTEXT{})*1))
		return *p
	}()) * 2)
}
func read_mvcontexts(bc *vp8_reader, mvc *MV_CONTEXT) {
	var i int = 0
	for {
		{
			var (
				up    *uint8 = &vp8_mv_update_probs[i].Prob[0]
				p     *uint8 = (*uint8)(unsafe.Pointer((*MV_CONTEXT)(unsafe.Add(unsafe.Pointer(mvc), unsafe.Sizeof(MV_CONTEXT{})*uintptr(i)))))
				pstop *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(p), MVPcount))
			)
			for {
				if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), int(*func() *uint8 {
					p := &up
					x := *p
					*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}())) != 0 {
					var x uint8 = uint8(int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 7)))
					if x != 0 {
						*p = x << 1
					} else {
						*p = 1
					}
				}
				if uintptr(unsafe.Pointer(func() *uint8 {
					p := &p
					*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
					return *p
				}())) >= uintptr(unsafe.Pointer(pstop)) {
					break
				}
			}
		}
		if func() int {
			p := &i
			*p++
			return *p
		}() >= 2 {
			break
		}
	}
}

var mbsplit_fill_count [4]uint8 = [4]uint8{8, 8, 4, 1}
var mbsplit_fill_offset [4][16]uint8 = [4][16]uint8{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, {0, 1, 4, 5, 8, 9, 12, 13, 2, 3, 6, 7, 10, 11, 14, 15}, {0, 1, 4, 5, 2, 3, 6, 7, 8, 9, 12, 13, 10, 11, 14, 15}, {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}}

func mb_mode_mv_init(pbi *VP8D_COMP) {
	var (
		bc  *vp8_reader = &pbi.Mbc[8]
		mvc *MV_CONTEXT = &pbi.Common.Fc.Mvc[0]
	)
	pbi.Common.Mb_no_coeff_skip = vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128)
	pbi.Prob_skip_false = 0
	if pbi.Common.Mb_no_coeff_skip != 0 {
		pbi.Prob_skip_false = uint8(int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 8)))
	}
	if pbi.Common.Frame_type != int(KEY_FRAME) {
		pbi.Prob_intra = uint8(int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 8)))
		pbi.Prob_last = uint8(int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 8)))
		pbi.Prob_gf = uint8(int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 8)))
		if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
			var i int = 0
			for {
				pbi.Common.Fc.Ymode_prob[i] = uint8(int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 8)))
				if func() int {
					p := &i
					*p++
					return *p
				}() >= 4 {
					break
				}
			}
		}
		if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 128) != 0 {
			var i int = 0
			for {
				pbi.Common.Fc.Uv_mode_prob[i] = uint8(int8(vp8_decode_value((*BOOL_DECODER)(unsafe.Pointer(bc)), 8)))
				if func() int {
					p := &i
					*p++
					return *p
				}() >= 3 {
					break
				}
			}
		}
		read_mvcontexts(bc, mvc)
	}
}

var vp8_sub_mv_ref_prob3 [8][3]uint8 = [8][3]uint8{{147, 136, 18}, {223, 1, 34}, {106, 145, 1}, {208, 1, 1}, {179, 121, 1}, {223, 1, 34}, {179, 121, 1}, {208, 1, 1}}

func get_sub_mv_ref_prob(left int, above int) *uint8 {
	var (
		lez  int = int(libc.BoolToInt(left == 0))
		aez  int = int(libc.BoolToInt(above == 0))
		lea  int = int(libc.BoolToInt(left == above))
		prob *uint8
	)
	prob = &vp8_sub_mv_ref_prob3[(aez<<2)|lez<<1|lea][0]
	return prob
}
func decode_split_mv(bc *vp8_reader, mi *ModeInfo, left_mb *ModeInfo, above_mb *ModeInfo, mbmi *MB_MODE_INFO, best_mv int_mv, mvc *MV_CONTEXT, mb_to_left_edge int, mb_to_right_edge int, mb_to_top_edge int, mb_to_bottom_edge int) {
	var (
		s     int
		num_p int
		j     int = 0
	)
	s = 3
	num_p = 16
	if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 110) != 0 {
		s = 2
		num_p = 4
		if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 111) != 0 {
			s = vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), 150)
			num_p = 2
		}
	}
	for {
		{
			var (
				leftmv  int_mv
				abovemv int_mv
				blockmv int_mv
				k       int
				prob    *uint8
			)
			k = int(vp8_mbsplit_offset[s][j])
			if (k & 3) == 0 {
				if int(left_mb.Mbmi.Mode) != SPLITMV {
					leftmv.As_int = left_mb.Mbmi.Mv.As_int
				} else {
					leftmv.As_int = ((*b_mode_info)(unsafe.Add(unsafe.Pointer(&left_mb.Bmi[k+4]), -int(unsafe.Sizeof(b_mode_info{})*1)))).Mv.As_int
				}
			} else {
				leftmv.As_int = ((*b_mode_info)(unsafe.Add(unsafe.Pointer(&mi.Bmi[k]), -int(unsafe.Sizeof(b_mode_info{})*1)))).Mv.As_int
			}
			if (k >> 2) == 0 {
				if int(above_mb.Mbmi.Mode) != SPLITMV {
					abovemv.As_int = above_mb.Mbmi.Mv.As_int
				} else {
					abovemv.As_int = ((*b_mode_info)(unsafe.Add(unsafe.Pointer(&above_mb.Bmi[k+16]), -int(unsafe.Sizeof(b_mode_info{})*4)))).Mv.As_int
				}
			} else {
				abovemv.As_int = ((*b_mode_info)(unsafe.Add(unsafe.Pointer(&mi.Bmi[k]), -int(unsafe.Sizeof(b_mode_info{})*4)))).Mv.As_int
			}
			prob = get_sub_mv_ref_prob(int(leftmv.As_int), int(abovemv.As_int))
			if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), int(*(*uint8)(unsafe.Add(unsafe.Pointer(prob), 0)))) != 0 {
				if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), int(*(*uint8)(unsafe.Add(unsafe.Pointer(prob), 1)))) != 0 {
					blockmv.As_int = 0
					if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), int(*(*uint8)(unsafe.Add(unsafe.Pointer(prob), 2)))) != 0 {
						blockmv.As_mv.Row = int16(read_mvcomponent(bc, (*MV_CONTEXT)(unsafe.Add(unsafe.Pointer(mvc), unsafe.Sizeof(MV_CONTEXT{})*0))) * 2)
						blockmv.As_mv.Row += best_mv.As_mv.Row
						blockmv.As_mv.Col = int16(read_mvcomponent(bc, (*MV_CONTEXT)(unsafe.Add(unsafe.Pointer(mvc), unsafe.Sizeof(MV_CONTEXT{})*1))) * 2)
						blockmv.As_mv.Col += best_mv.As_mv.Col
					}
				} else {
					blockmv.As_int = abovemv.As_int
				}
			} else {
				blockmv.As_int = leftmv.As_int
			}
			mbmi.Need_to_clamp_mvs |= uint8(vp8_check_mv_bounds(&blockmv, mb_to_left_edge, mb_to_right_edge, mb_to_top_edge, mb_to_bottom_edge))
			{
				var (
					fill_offset *uint8
					fill_count  uint = uint(mbsplit_fill_count[s])
				)
				fill_offset = &mbsplit_fill_offset[s][int(uint8(int8(j)))*int(mbsplit_fill_count[s])]
				for {
					mi.Bmi[*fill_offset].Mv.As_int = blockmv.As_int
					fill_offset = (*uint8)(unsafe.Add(unsafe.Pointer(fill_offset), 1))
					if func() uint {
						p := &fill_count
						*p--
						return *p
					}() == 0 {
						break
					}
				}
			}
		}
		if func() int {
			p := &j
			*p++
			return *p
		}() >= num_p {
			break
		}
	}
	mbmi.Partitioning = uint8(int8(s))
}
func read_mb_modes_mv(pbi *VP8D_COMP, mi *ModeInfo, mbmi *MB_MODE_INFO) {
	var bc *vp8_reader = &pbi.Mbc[8]
	mbmi.Ref_frame = uint8(int8(int(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), int(pbi.Prob_intra)))))
	if mbmi.Ref_frame != 0 {
		const (
			CNT_INTRA = iota
			CNT_NEAREST
			CNT_NEAR
			CNT_SPLITMV
		)
		var cnt [4]int
		var cntx *int = &cnt[0]
		var near_mvs [4]int_mv
		var nmv *int_mv = &near_mvs[0]
		var mis int = pbi.Mb.Mode_info_stride
		var above *ModeInfo = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mi), -int(unsafe.Sizeof(ModeInfo{})*uintptr(mis))))
		var left *ModeInfo = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mi), -int(unsafe.Sizeof(ModeInfo{})*1)))
		var aboveleft *ModeInfo = (*ModeInfo)(unsafe.Add(unsafe.Pointer(above), -int(unsafe.Sizeof(ModeInfo{})*1)))
		var ref_frame_sign_bias *int = &pbi.Common.Ref_frame_sign_bias[0]
		mbmi.Need_to_clamp_mvs = 0
		if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), int(pbi.Prob_last)) != 0 {
			mbmi.Ref_frame = uint8(int8(int(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), int(pbi.Prob_gf)) + 2)))
		}
		(*(*int_mv)(unsafe.Add(unsafe.Pointer(nmv), unsafe.Sizeof(int_mv{})*0))).As_int = func() uint32 {
			p := &(*(*int_mv)(unsafe.Add(unsafe.Pointer(nmv), unsafe.Sizeof(int_mv{})*1))).As_int
			(*(*int_mv)(unsafe.Add(unsafe.Pointer(nmv), unsafe.Sizeof(int_mv{})*1))).As_int = func() uint32 {
				p := &(*(*int_mv)(unsafe.Add(unsafe.Pointer(nmv), unsafe.Sizeof(int_mv{})*2))).As_int
				(*(*int_mv)(unsafe.Add(unsafe.Pointer(nmv), unsafe.Sizeof(int_mv{})*2))).As_int = 0
				return *p
			}()
			return *p
		}()
		cnt[0] = func() int {
			p := &cnt[1]
			cnt[1] = func() int {
				p := &cnt[2]
				cnt[2] = func() int {
					p := &cnt[3]
					cnt[3] = 0
					return *p
				}()
				return *p
			}()
			return *p
		}()
		if int(above.Mbmi.Ref_frame) != INTRA_FRAME {
			if above.Mbmi.Mv.As_int != 0 {
				(func() *int_mv {
					p := &nmv
					*p = (*int_mv)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int_mv{})*1))
					return *p
				}()).As_int = above.Mbmi.Mv.As_int
				mv_bias(*(*int)(unsafe.Add(unsafe.Pointer(ref_frame_sign_bias), unsafe.Sizeof(int(0))*uintptr(above.Mbmi.Ref_frame))), int(mbmi.Ref_frame), nmv, ref_frame_sign_bias)
				cntx = (*int)(unsafe.Add(unsafe.Pointer(cntx), unsafe.Sizeof(int(0))*1))
			}
			*cntx += 2
		}
		if int(left.Mbmi.Ref_frame) != INTRA_FRAME {
			if left.Mbmi.Mv.As_int != 0 {
				var this_mv int_mv
				this_mv.As_int = left.Mbmi.Mv.As_int
				mv_bias(*(*int)(unsafe.Add(unsafe.Pointer(ref_frame_sign_bias), unsafe.Sizeof(int(0))*uintptr(left.Mbmi.Ref_frame))), int(mbmi.Ref_frame), &this_mv, ref_frame_sign_bias)
				if this_mv.As_int != nmv.As_int {
					(func() *int_mv {
						p := &nmv
						*p = (*int_mv)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int_mv{})*1))
						return *p
					}()).As_int = this_mv.As_int
					cntx = (*int)(unsafe.Add(unsafe.Pointer(cntx), unsafe.Sizeof(int(0))*1))
				}
				*cntx += 2
			} else {
				cnt[CNT_INTRA] += 2
			}
		}
		if int(aboveleft.Mbmi.Ref_frame) != INTRA_FRAME {
			if aboveleft.Mbmi.Mv.As_int != 0 {
				var this_mv int_mv
				this_mv.As_int = aboveleft.Mbmi.Mv.As_int
				mv_bias(*(*int)(unsafe.Add(unsafe.Pointer(ref_frame_sign_bias), unsafe.Sizeof(int(0))*uintptr(aboveleft.Mbmi.Ref_frame))), int(mbmi.Ref_frame), &this_mv, ref_frame_sign_bias)
				if this_mv.As_int != nmv.As_int {
					(func() *int_mv {
						p := &nmv
						*p = (*int_mv)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int_mv{})*1))
						return *p
					}()).As_int = this_mv.As_int
					cntx = (*int)(unsafe.Add(unsafe.Pointer(cntx), unsafe.Sizeof(int(0))*1))
				}
				*cntx += 1
			} else {
				cnt[CNT_INTRA] += 1
			}
		}
		if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), vp8_mode_contexts[cnt[CNT_INTRA]][0]) != 0 {
			cnt[CNT_NEAREST] += int(libc.BoolToInt(cnt[CNT_SPLITMV] > 0) & libc.BoolToInt(nmv.As_int == near_mvs[CNT_NEAREST].As_int))
			if cnt[CNT_NEAR] > cnt[CNT_NEAREST] {
				var tmp int
				tmp = cnt[CNT_NEAREST]
				cnt[CNT_NEAREST] = cnt[CNT_NEAR]
				cnt[CNT_NEAR] = tmp
				tmp = int(near_mvs[CNT_NEAREST].As_int)
				near_mvs[CNT_NEAREST].As_int = near_mvs[CNT_NEAR].As_int
				near_mvs[CNT_NEAR].As_int = uint32(tmp)
			}
			if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), vp8_mode_contexts[cnt[CNT_NEAREST]][1]) != 0 {
				if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), vp8_mode_contexts[cnt[CNT_NEAR]][2]) != 0 {
					var (
						mb_to_top_edge    int
						mb_to_bottom_edge int
						mb_to_left_edge   int
						mb_to_right_edge  int
						mvc               *MV_CONTEXT = &pbi.Common.Fc.Mvc[0]
						near_index        int
					)
					mb_to_top_edge = pbi.Mb.Mb_to_top_edge
					mb_to_bottom_edge = pbi.Mb.Mb_to_bottom_edge
					mb_to_top_edge -= 16 << 3
					mb_to_bottom_edge += 16 << 3
					mb_to_right_edge = pbi.Mb.Mb_to_right_edge
					mb_to_right_edge += 16 << 3
					mb_to_left_edge = pbi.Mb.Mb_to_left_edge
					mb_to_left_edge -= 16 << 3
					near_index = CNT_INTRA + int(libc.BoolToInt(cnt[CNT_NEAREST] >= cnt[CNT_INTRA]))
					vp8_clamp_mv2(&near_mvs[near_index], &pbi.Mb)
					cnt[CNT_SPLITMV] = int(libc.BoolToInt(int(above.Mbmi.Mode) == SPLITMV)+libc.BoolToInt(int(left.Mbmi.Mode) == SPLITMV))*2 + int(libc.BoolToInt(int(aboveleft.Mbmi.Mode) == SPLITMV))
					if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(bc)), vp8_mode_contexts[cnt[CNT_SPLITMV]][3]) != 0 {
						decode_split_mv(bc, mi, left, above, mbmi, near_mvs[near_index], mvc, mb_to_left_edge, mb_to_right_edge, mb_to_top_edge, mb_to_bottom_edge)
						mbmi.Mv.As_int = mi.Bmi[15].Mv.As_int
						mbmi.Mode = uint8(int8(SPLITMV))
						mbmi.Is_4x4 = 1
					} else {
						var mbmi_mv *int_mv = &mbmi.Mv
						read_mv(bc, &mbmi_mv.As_mv, mvc)
						mbmi_mv.As_mv.Row += near_mvs[near_index].As_mv.Row
						mbmi_mv.As_mv.Col += near_mvs[near_index].As_mv.Col
						mbmi.Need_to_clamp_mvs = uint8(vp8_check_mv_bounds(mbmi_mv, mb_to_left_edge, mb_to_right_edge, mb_to_top_edge, mb_to_bottom_edge))
						mbmi.Mode = uint8(int8(NEWMV))
					}
				} else {
					mbmi.Mode = uint8(int8(NEARMV))
					mbmi.Mv.As_int = near_mvs[CNT_NEAR].As_int
					vp8_clamp_mv2(&mbmi.Mv, &pbi.Mb)
				}
			} else {
				mbmi.Mode = uint8(int8(NEARESTMV))
				mbmi.Mv.As_int = near_mvs[CNT_NEAREST].As_int
				vp8_clamp_mv2(&mbmi.Mv, &pbi.Mb)
			}
		} else {
			mbmi.Mode = uint8(int8(ZEROMV))
			mbmi.Mv.As_int = 0
		}
	} else {
		mbmi.Mv.As_int = 0
		if int(func() uint8 {
			p := &mbmi.Mode
			mbmi.Mode = uint8(int8(read_ymode(bc, &pbi.Common.Fc.Ymode_prob[0])))
			return *p
		}()) == B_PRED {
			var j int = 0
			mbmi.Is_4x4 = 1
			for {
				mi.Bmi[j].As_mode = read_bmode(bc, &pbi.Common.Fc.Bmode_prob[0])
				if func() int {
					p := &j
					*p++
					return *p
				}() >= 16 {
					break
				}
			}
		}
		mbmi.Uv_mode = uint8(int8(read_uv_mode(bc, &pbi.Common.Fc.Uv_mode_prob[0])))
	}
}
func read_mb_features(r *vp8_reader, mi *MB_MODE_INFO, x *MacroBlockd) {
	if int(x.Segmentation_enabled) != 0 && int(x.Update_mb_segmentation_map) != 0 {
		if vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(r)), int(x.Mb_segment_tree_probs[0])) != 0 {
			mi.Segment_id = uint8(uint8(int8(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(r)), int(x.Mb_segment_tree_probs[2])) + 2)))
		} else {
			mi.Segment_id = uint8(uint8(int8(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(r)), int(x.Mb_segment_tree_probs[1])))))
		}
	}
}
func decode_mb_mode_mvs(pbi *VP8D_COMP, mi *ModeInfo) {
	if int(pbi.Mb.Update_mb_segmentation_map) != 0 {
		read_mb_features(&pbi.Mbc[8], &mi.Mbmi, &pbi.Mb)
	} else if pbi.Common.Frame_type == int(KEY_FRAME) {
		mi.Mbmi.Segment_id = 0
	}
	if pbi.Common.Mb_no_coeff_skip != 0 {
		mi.Mbmi.Mb_skip_coeff = uint8(int8(vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(&pbi.Mbc[8])), int(pbi.Prob_skip_false))))
	} else {
		mi.Mbmi.Mb_skip_coeff = 0
	}
	mi.Mbmi.Is_4x4 = 0
	if pbi.Common.Frame_type == int(KEY_FRAME) {
		read_kf_modes(pbi, mi)
	} else {
		read_mb_modes_mv(pbi, mi, &mi.Mbmi)
	}
}
func vp8_decode_mode_mvs(pbi *VP8D_COMP) {
	var (
		mi                     *ModeInfo = pbi.Common.Mi
		mb_row                 int       = -1
		mb_to_right_edge_start int
	)
	mb_mode_mv_init(pbi)
	pbi.Mb.Mb_to_top_edge = 0
	pbi.Mb.Mb_to_bottom_edge = ((pbi.Common.Mb_rows - 1) * 16) << 3
	mb_to_right_edge_start = ((pbi.Common.Mb_cols - 1) * 16) << 3
	for func() int {
		p := &mb_row
		*p++
		return *p
	}() < pbi.Common.Mb_rows {
		var mb_col int = -1
		pbi.Mb.Mb_to_left_edge = 0
		pbi.Mb.Mb_to_right_edge = mb_to_right_edge_start
		for func() int {
			p := &mb_col
			*p++
			return *p
		}() < pbi.Common.Mb_cols {
			decode_mb_mode_mvs(pbi, mi)
			pbi.Mb.Mb_to_left_edge -= 16 << 3
			pbi.Mb.Mb_to_right_edge -= 16 << 3
			mi = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mi), unsafe.Sizeof(ModeInfo{})*1))
		}
		pbi.Mb.Mb_to_top_edge -= 16 << 3
		pbi.Mb.Mb_to_bottom_edge -= 16 << 3
		mi = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mi), unsafe.Sizeof(ModeInfo{})*1))
	}
}
