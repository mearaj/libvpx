package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const LEFT_TOP_MARGIN = 128
const RIGHT_BOTTOM_MARGIN = 128

func mv_bias(refmb_ref_frame_sign_bias int, refframe int, mvp *int_mv, ref_frame_sign_bias *int) {
	if refmb_ref_frame_sign_bias != *(*int)(unsafe.Add(unsafe.Pointer(ref_frame_sign_bias), unsafe.Sizeof(int(0))*uintptr(refframe))) {
		mvp.As_mv.Row *= -1
		mvp.As_mv.Col *= -1
	}
}
func vp8_clamp_mv2(mv *int_mv, xd *MacroBlockd) {
	if int(mv.As_mv.Col) < (xd.Mb_to_left_edge - (16 << 3)) {
		mv.As_mv.Col = int16(xd.Mb_to_left_edge - (16 << 3))
	} else if int(mv.As_mv.Col) > xd.Mb_to_right_edge+(16<<3) {
		mv.As_mv.Col = int16(xd.Mb_to_right_edge + (16 << 3))
	}
	if int(mv.As_mv.Row) < (xd.Mb_to_top_edge - (16 << 3)) {
		mv.As_mv.Row = int16(xd.Mb_to_top_edge - (16 << 3))
	} else if int(mv.As_mv.Row) > xd.Mb_to_bottom_edge+(16<<3) {
		mv.As_mv.Row = int16(xd.Mb_to_bottom_edge + (16 << 3))
	}
}
func vp8_clamp_mv(mv *int_mv, mb_to_left_edge int, mb_to_right_edge int, mb_to_top_edge int, mb_to_bottom_edge int) {
	if int(mv.As_mv.Col) < mb_to_left_edge {
		mv.As_mv.Col = int16(mb_to_left_edge)
	} else {
		mv.As_mv.Col = mv.As_mv.Col
	}
	if int(mv.As_mv.Col) > mb_to_right_edge {
		mv.As_mv.Col = int16(mb_to_right_edge)
	} else {
		mv.As_mv.Col = mv.As_mv.Col
	}
	if int(mv.As_mv.Row) < mb_to_top_edge {
		mv.As_mv.Row = int16(mb_to_top_edge)
	} else {
		mv.As_mv.Row = mv.As_mv.Row
	}
	if int(mv.As_mv.Row) > mb_to_bottom_edge {
		mv.As_mv.Row = int16(mb_to_bottom_edge)
	} else {
		mv.As_mv.Row = mv.As_mv.Row
	}
}
func vp8_check_mv_bounds(mv *int_mv, mb_to_left_edge int, mb_to_right_edge int, mb_to_top_edge int, mb_to_bottom_edge int) uint {
	var need_to_clamp uint
	need_to_clamp = uint(libc.BoolToInt(int(mv.As_mv.Col) < mb_to_left_edge))
	need_to_clamp |= uint(libc.BoolToInt(int(mv.As_mv.Col) > mb_to_right_edge))
	need_to_clamp |= uint(libc.BoolToInt(int(mv.As_mv.Row) < mb_to_top_edge))
	need_to_clamp |= uint(libc.BoolToInt(int(mv.As_mv.Row) > mb_to_bottom_edge))
	return need_to_clamp
}
func left_block_mv(cur_mb *ModeInfo, b int) uint32 {
	if (b & 3) == 0 {
		cur_mb = (*ModeInfo)(unsafe.Add(unsafe.Pointer(cur_mb), -int(unsafe.Sizeof(ModeInfo{})*1)))
		if int(cur_mb.Mbmi.Mode) != SPLITMV {
			return cur_mb.Mbmi.Mv.As_int
		}
		b += 4
	}
	return ((*b_mode_info)(unsafe.Add(unsafe.Pointer(&cur_mb.Bmi[b]), -int(unsafe.Sizeof(b_mode_info{})*1)))).Mv.As_int
}
func above_block_mv(cur_mb *ModeInfo, b int, mi_stride int) uint32 {
	if (b >> 2) == 0 {
		cur_mb = (*ModeInfo)(unsafe.Add(unsafe.Pointer(cur_mb), -int(unsafe.Sizeof(ModeInfo{})*uintptr(mi_stride))))
		if int(cur_mb.Mbmi.Mode) != SPLITMV {
			return cur_mb.Mbmi.Mv.As_int
		}
		b += 16
	}
	return (&cur_mb.Bmi[b-4]).Mv.As_int
}
func left_block_mode(cur_mb *ModeInfo, b int) B_PREDICTION_MODE {
	if (b & 3) == 0 {
		cur_mb = (*ModeInfo)(unsafe.Add(unsafe.Pointer(cur_mb), -int(unsafe.Sizeof(ModeInfo{})*1)))
		switch int(cur_mb.Mbmi.Mode) {
		case B_PRED:
			return (&cur_mb.Bmi[b+3]).As_mode
		case DC_PRED:
			return B_PREDICTION_MODE(B_DC_PRED)
		case V_PRED:
			return B_PREDICTION_MODE(B_VE_PRED)
		case H_PRED:
			return B_PREDICTION_MODE(B_HE_PRED)
		case TM_PRED:
			return B_PREDICTION_MODE(B_TM_PRED)
		default:
			return B_PREDICTION_MODE(B_DC_PRED)
		}
	}
	return ((*b_mode_info)(unsafe.Add(unsafe.Pointer(&cur_mb.Bmi[b]), -int(unsafe.Sizeof(b_mode_info{})*1)))).As_mode
}
func above_block_mode(cur_mb *ModeInfo, b int, mi_stride int) B_PREDICTION_MODE {
	if (b >> 2) == 0 {
		cur_mb = (*ModeInfo)(unsafe.Add(unsafe.Pointer(cur_mb), -int(unsafe.Sizeof(ModeInfo{})*uintptr(mi_stride))))
		switch int(cur_mb.Mbmi.Mode) {
		case B_PRED:
			return (&cur_mb.Bmi[b+12]).As_mode
		case DC_PRED:
			return B_PREDICTION_MODE(B_DC_PRED)
		case V_PRED:
			return B_PREDICTION_MODE(B_VE_PRED)
		case H_PRED:
			return B_PREDICTION_MODE(B_HE_PRED)
		case TM_PRED:
			return B_PREDICTION_MODE(B_TM_PRED)
		default:
			return B_PREDICTION_MODE(B_DC_PRED)
		}
	}
	return ((*b_mode_info)(unsafe.Add(unsafe.Pointer(&cur_mb.Bmi[b]), -int(unsafe.Sizeof(b_mode_info{})*4)))).As_mode
}

var vp8_mbsplit_offset [4][16]uint8 = [4][16]uint8{{0, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 2, 8, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}}

func vp8_find_near_mvs(xd *MacroBlockd, here *ModeInfo, nearest *int_mv, nearby *int_mv, best_mv *int_mv, near_mv_ref_cnts [4]int, refframe int, ref_frame_sign_bias *int) {
	var (
		above     *ModeInfo = (*ModeInfo)(unsafe.Add(unsafe.Pointer(here), -int(unsafe.Sizeof(ModeInfo{})*uintptr(xd.Mode_info_stride))))
		left      *ModeInfo = (*ModeInfo)(unsafe.Add(unsafe.Pointer(here), -int(unsafe.Sizeof(ModeInfo{})*1)))
		aboveleft *ModeInfo = (*ModeInfo)(unsafe.Add(unsafe.Pointer(above), -int(unsafe.Sizeof(ModeInfo{})*1)))
		near_mvs  [4]int_mv
		mv        *int_mv = &near_mvs[0]
		cntx      *int    = &near_mv_ref_cnts[0]
	)
	const (
		CNT_INTRA = iota
		CNT_NEAREST
		CNT_NEAR
		CNT_SPLITMV
	)
	(*(*int_mv)(unsafe.Add(unsafe.Pointer(mv), unsafe.Sizeof(int_mv{})*0))).As_int = func() uint32 {
		p := &(*(*int_mv)(unsafe.Add(unsafe.Pointer(mv), unsafe.Sizeof(int_mv{})*1))).As_int
		(*(*int_mv)(unsafe.Add(unsafe.Pointer(mv), unsafe.Sizeof(int_mv{})*1))).As_int = func() uint32 {
			p := &(*(*int_mv)(unsafe.Add(unsafe.Pointer(mv), unsafe.Sizeof(int_mv{})*2))).As_int
			(*(*int_mv)(unsafe.Add(unsafe.Pointer(mv), unsafe.Sizeof(int_mv{})*2))).As_int = 0
			return *p
		}()
		return *p
	}()
	near_mv_ref_cnts[0] = func() int {
		p := &near_mv_ref_cnts[1]
		near_mv_ref_cnts[1] = func() int {
			p := &near_mv_ref_cnts[2]
			near_mv_ref_cnts[2] = func() int {
				p := &near_mv_ref_cnts[3]
				near_mv_ref_cnts[3] = 0
				return *p
			}()
			return *p
		}()
		return *p
	}()
	if int(above.Mbmi.Ref_frame) != INTRA_FRAME {
		if above.Mbmi.Mv.As_int != 0 {
			(func() *int_mv {
				p := &mv
				*p = (*int_mv)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int_mv{})*1))
				return *p
			}()).As_int = above.Mbmi.Mv.As_int
			mv_bias(*(*int)(unsafe.Add(unsafe.Pointer(ref_frame_sign_bias), unsafe.Sizeof(int(0))*uintptr(above.Mbmi.Ref_frame))), refframe, mv, ref_frame_sign_bias)
			cntx = (*int)(unsafe.Add(unsafe.Pointer(cntx), unsafe.Sizeof(int(0))*1))
		}
		*cntx += 2
	}
	if int(left.Mbmi.Ref_frame) != INTRA_FRAME {
		if left.Mbmi.Mv.As_int != 0 {
			var this_mv int_mv
			this_mv.As_int = left.Mbmi.Mv.As_int
			mv_bias(*(*int)(unsafe.Add(unsafe.Pointer(ref_frame_sign_bias), unsafe.Sizeof(int(0))*uintptr(left.Mbmi.Ref_frame))), refframe, &this_mv, ref_frame_sign_bias)
			if this_mv.As_int != mv.As_int {
				(func() *int_mv {
					p := &mv
					*p = (*int_mv)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int_mv{})*1))
					return *p
				}()).As_int = this_mv.As_int
				cntx = (*int)(unsafe.Add(unsafe.Pointer(cntx), unsafe.Sizeof(int(0))*1))
			}
			*cntx += 2
		} else {
			near_mv_ref_cnts[CNT_INTRA] += 2
		}
	}
	if int(aboveleft.Mbmi.Ref_frame) != INTRA_FRAME {
		if aboveleft.Mbmi.Mv.As_int != 0 {
			var this_mv int_mv
			this_mv.As_int = aboveleft.Mbmi.Mv.As_int
			mv_bias(*(*int)(unsafe.Add(unsafe.Pointer(ref_frame_sign_bias), unsafe.Sizeof(int(0))*uintptr(aboveleft.Mbmi.Ref_frame))), refframe, &this_mv, ref_frame_sign_bias)
			if this_mv.As_int != mv.As_int {
				(func() *int_mv {
					p := &mv
					*p = (*int_mv)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int_mv{})*1))
					return *p
				}()).As_int = this_mv.As_int
				cntx = (*int)(unsafe.Add(unsafe.Pointer(cntx), unsafe.Sizeof(int(0))*1))
			}
			*cntx += 1
		} else {
			near_mv_ref_cnts[CNT_INTRA] += 1
		}
	}
	if near_mv_ref_cnts[CNT_SPLITMV] != 0 {
		if mv.As_int == near_mvs[CNT_NEAREST].As_int {
			near_mv_ref_cnts[CNT_NEAREST] += 1
		}
	}
	near_mv_ref_cnts[CNT_SPLITMV] = int(libc.BoolToInt(int(above.Mbmi.Mode) == int(SPLITMV))+libc.BoolToInt(int(left.Mbmi.Mode) == SPLITMV))*2 + int(libc.BoolToInt(int(aboveleft.Mbmi.Mode) == SPLITMV))
	if near_mv_ref_cnts[CNT_NEAR] > near_mv_ref_cnts[CNT_NEAREST] {
		var tmp int
		tmp = near_mv_ref_cnts[CNT_NEAREST]
		near_mv_ref_cnts[CNT_NEAREST] = near_mv_ref_cnts[CNT_NEAR]
		near_mv_ref_cnts[CNT_NEAR] = tmp
		tmp = int(near_mvs[CNT_NEAREST].As_int)
		near_mvs[CNT_NEAREST].As_int = near_mvs[CNT_NEAR].As_int
		near_mvs[CNT_NEAR].As_int = uint32(tmp)
	}
	if near_mv_ref_cnts[CNT_NEAREST] >= near_mv_ref_cnts[CNT_INTRA] {
		near_mvs[CNT_INTRA] = near_mvs[CNT_NEAREST]
	}
	best_mv.As_int = near_mvs[0].As_int
	nearest.As_int = near_mvs[CNT_NEAREST].As_int
	nearby.As_int = near_mvs[CNT_NEAR].As_int
}
func invert_and_clamp_mvs(inv *int_mv, src *int_mv, xd *MacroBlockd) {
	inv.As_mv.Row = int16(int(src.As_mv.Row) * (-1))
	inv.As_mv.Col = int16(int(src.As_mv.Col) * (-1))
	vp8_clamp_mv2(inv, xd)
	vp8_clamp_mv2(src, xd)
}
func vp8_find_near_mvs_bias(xd *MacroBlockd, here *ModeInfo, mode_mv_sb [2][10]int_mv, best_mv_sb [2]int_mv, cnt [4]int, refframe int, ref_frame_sign_bias *int) int {
	var sign_bias int = *(*int)(unsafe.Add(unsafe.Pointer(ref_frame_sign_bias), unsafe.Sizeof(int(0))*uintptr(refframe)))
	vp8_find_near_mvs(xd, here, &mode_mv_sb[sign_bias][NEARESTMV], &mode_mv_sb[sign_bias][NEARMV], &best_mv_sb[sign_bias], cnt, refframe, ref_frame_sign_bias)
	invert_and_clamp_mvs(&mode_mv_sb[libc.BoolToInt(sign_bias == 0)][NEARESTMV], &mode_mv_sb[sign_bias][NEARESTMV], xd)
	invert_and_clamp_mvs(&mode_mv_sb[libc.BoolToInt(sign_bias == 0)][NEARMV], &mode_mv_sb[sign_bias][NEARMV], xd)
	invert_and_clamp_mvs(&best_mv_sb[libc.BoolToInt(sign_bias == 0)], &best_mv_sb[sign_bias], xd)
	return sign_bias
}
func vp8_mv_ref_probs(p [4]uint8, near_mv_ref_ct [4]int) *uint8 {
	p[0] = uint8(int8(vp8_mode_contexts[near_mv_ref_ct[0]][0]))
	p[1] = uint8(int8(vp8_mode_contexts[near_mv_ref_ct[1]][1]))
	p[2] = uint8(int8(vp8_mode_contexts[near_mv_ref_ct[2]][2]))
	p[3] = uint8(int8(vp8_mode_contexts[near_mv_ref_ct[3]][3]))
	return &p[0]
}
