package vp8

const MAX_LOOP_FILTER = 63
const PARTIAL_FRAME_FRACTION = 8
const SIMD_WIDTH = 16

type LOOPFILTERTYPE int

const (
	NORMAL_LOOPFILTER LOOPFILTERTYPE = 0
	SIMPLE_LOOPFILTER LOOPFILTERTYPE = 1
)

type loop_filter_info_n struct {
	Mblim       [64][16]uint8
	Blim        [64][16]uint8
	Lim         [64][16]uint8
	Hev_thr     [4][16]uint8
	Lvl         [4][4][4]uint8
	Hev_thr_lut [2][64]uint8
	Mode_lf_lut [10]uint8
}
type loop_filter_info struct {
	Mblim   *uint8
	Blim    *uint8
	Lim     *uint8
	Hev_thr *uint8
}
type loop_filter_uvfunction func(u *uint8, p int, blimit *uint8, limit *uint8, thresh *uint8, v *uint8)
