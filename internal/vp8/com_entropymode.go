package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const VP8_NUMMBSPLITS = 4
const SUBMVREF_COUNT = 5
const (
	SUBMVREF_NORMAL = int(iota)
	SUBMVREF_LEFT_ZED
	SUBMVREF_ABOVE_ZED
	SUBMVREF_LEFT_ABOVE_SAME
	SUBMVREF_LEFT_ABOVE_ZED
)

type vp8_mbsplit [16]int

func vp8_mv_cont(l *int_mv, a *int_mv) int {
	var (
		lez int = int(libc.BoolToInt(l.As_int == 0))
		aez int = int(libc.BoolToInt(a.As_int == 0))
		lea int = int(libc.BoolToInt(l.As_int == a.As_int))
	)
	if lea != 0 && lez != 0 {
		return SUBMVREF_LEFT_ABOVE_ZED
	}
	if lea != 0 {
		return SUBMVREF_LEFT_ABOVE_SAME
	}
	if aez != 0 {
		return SUBMVREF_ABOVE_ZED
	}
	if lez != 0 {
		return SUBMVREF_LEFT_ZED
	}
	return SUBMVREF_NORMAL
}
func vp8_init_mbmode_probs(x *VP8Common) {
	libc.MemCpy(unsafe.Pointer(&x.Fc.Ymode_prob[0]), unsafe.Pointer(&vp8_ymode_prob[0]), int(unsafe.Sizeof([4]uint8{})))
	libc.MemCpy(unsafe.Pointer(&x.Fc.Uv_mode_prob[0]), unsafe.Pointer(&vp8_uv_mode_prob[0]), int(unsafe.Sizeof([3]uint8{})))
	libc.MemCpy(unsafe.Pointer(&x.Fc.Sub_mv_ref_prob[0]), unsafe.Pointer(&sub_mv_ref_prob[0]), int(unsafe.Sizeof([3]uint8{})))
}
func vp8_default_bmode_probs(dest [9]uint8) {
	libc.MemCpy(unsafe.Pointer(&dest[0]), unsafe.Pointer(&vp8_bmode_prob[0]), int(unsafe.Sizeof([9]uint8{})))
}
