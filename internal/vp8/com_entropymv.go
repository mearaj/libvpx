package vp8

import "math"

const (
	mv_max       = 1023
	MVvals       = (mv_max * 2) + 1
	mvfp_max     = math.MaxUint8
	MVfpvals     = (mvfp_max * 2) + 1
	mvlong_width = 10
	mvnum_short  = 8
	mvpis_short  = 0
	MVPsign      = 1
	MVPshort     = 2
	MVPbits      = MVPshort + mvnum_short - 1
	MVPcount     = MVPbits + mvlong_width
)

type mv_context struct {
	Prob [19]uint8
}
type MV_CONTEXT mv_context

var vp8_mv_update_probs [2]MV_CONTEXT = [2]MV_CONTEXT{{Prob: [19]uint8{237, 246, 253, 253, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 250, 250, 252, 254, 254}}, {Prob: [19]uint8{231, 243, 245, 253, 254, 254, 254, 254, 254, 254, 254, 254, 254, 254, 251, 251, 254, 254, 254}}}
var vp8_default_mv_context [2]MV_CONTEXT = [2]MV_CONTEXT{{Prob: [19]uint8{162, 128, 225, 146, 172, 147, 214, 39, 156, 128, 129, 132, 75, 145, 178, 206, 239, 254, 254}}, {Prob: [19]uint8{164, 128, 204, 170, 119, 235, 140, 230, 228, 128, 130, 130, 74, 148, 180, 203, 236, 254, 254}}}
