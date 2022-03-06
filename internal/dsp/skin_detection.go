package dsp

import "github.com/gotranspile/cxgo/runtime/libc"

const MODEL_MODE = 1

var skin_mean [5][2]int = [5][2]int{{7463, 9614}, {6400, 0x2800}, {7040, 0x2800}, {8320, 9280}, {6800, 9614}}
var skin_inv_cov [4]int = [4]int{4107, 1663, 1663, 2157}
var skin_threshold [6]int = [6]int{0x17F74C, 1400000, 800000, 800000, 800000, 800000}
var y_low int = 40
var y_high int = 220

func vpx_evaluate_skin_color_difference(cb int, cr int, idx int) int {
	var (
		cb_q6         int = cb << 6
		cr_q6         int = cr << 6
		cb_diff_q12   int = (cb_q6 - skin_mean[idx][0]) * (cb_q6 - skin_mean[idx][0])
		cbcr_diff_q12 int = (cb_q6 - skin_mean[idx][0]) * (cr_q6 - skin_mean[idx][1])
		cr_diff_q12   int = (cr_q6 - skin_mean[idx][1]) * (cr_q6 - skin_mean[idx][1])
		cb_diff_q2    int = (cb_diff_q12 + (1 << 9)) >> 10
		cbcr_diff_q2  int = (cbcr_diff_q12 + (1 << 9)) >> 10
		cr_diff_q2    int = (cr_diff_q12 + (1 << 9)) >> 10
		skin_diff     int = skin_inv_cov[0]*cb_diff_q2 + skin_inv_cov[1]*cbcr_diff_q2 + skin_inv_cov[2]*cbcr_diff_q2 + skin_inv_cov[3]*cr_diff_q2
	)
	return skin_diff
}
func VpxSkinPixel(y int, cb int, cr int, motion int) int {
	if y < y_low || y > y_high {
		return 0
	} else if MODEL_MODE == 0 {
		return int(libc.BoolToInt(vpx_evaluate_skin_color_difference(cb, cr, 0) < skin_threshold[0]))
	} else {
		var i int = 0
		if cb == 128 && cr == 128 {
			return 0
		}
		if cb > 150 && cr < 110 {
			return 0
		}
		for ; i < 5; i++ {
			var skin_color_diff int = vpx_evaluate_skin_color_difference(cb, cr, i)
			if skin_color_diff < skin_threshold[i+1] {
				if y < 60 && skin_color_diff > (skin_threshold[i+1]>>2)*3 {
					return 0
				} else if motion == 0 && skin_color_diff > (skin_threshold[i+1]>>1) {
					return 0
				} else {
					return 1
				}
			}
			if skin_color_diff > (skin_threshold[i+1] << 3) {
				return 0
			}
		}
		return 0
	}
}
