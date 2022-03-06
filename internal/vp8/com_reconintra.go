package vp8

import (
	"github.com/mearaj/libvpx/internal/dsp"
	"github.com/mearaj/libvpx/internal/ports"
	"unsafe"
)

const (
	SIZE_16 = iota
	SIZE_8
	NUM_SIZES
)

type intra_pred_fn func(dst *uint8, stride int64, above *uint8, left *uint8)

var pred [4][2]intra_pred_fn
var dc_pred [2][2][2]intra_pred_fn

func vp8_init_intra_predictors_internal() {
	pred[V_PRED][SIZE_16] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxVPredictor16x16C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	pred[H_PRED][SIZE_16] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxHPredictor16x16C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	pred[TM_PRED][SIZE_16] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxTmPredictor16x16C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	dc_pred[0][0][SIZE_16] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxDc128Predictor16x16C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	dc_pred[0][1][SIZE_16] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxDcTopPredictor16x16C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	dc_pred[1][0][SIZE_16] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxDcLeftPredictor16x16C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	dc_pred[1][1][SIZE_16] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxDcPredictor16x16C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	pred[V_PRED][SIZE_8] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxVPredictor8x8C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	pred[H_PRED][SIZE_8] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxHPredictor8x8C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	pred[TM_PRED][SIZE_8] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxTmPredictor8x8C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	dc_pred[0][0][SIZE_8] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxDc128Predictor8x8C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	dc_pred[0][1][SIZE_8] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxDcTopPredictor8x8C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	dc_pred[1][0][SIZE_8] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxDcLeftPredictor8x8C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	dc_pred[1][1][SIZE_8] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxDcPredictor8x8C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	vp8_init_intra4x4_predictors_internal()
}
func vp8_build_intra_predictors_mby_s(x *MacroBlockd, yabove_row *uint8, yleft *uint8, left_stride int, ypred_ptr *uint8, y_stride int) {
	var (
		mode      int = int(x.Mode_info_context.Mbmi.Mode)
		yleft_col [16]uint8
		i         int
		fn        intra_pred_fn
	)
	for i = 0; i < 16; i++ {
		yleft_col[i] = *(*uint8)(unsafe.Add(unsafe.Pointer(yleft), i*left_stride))
	}
	if mode == int(DC_PRED) {
		fn = dc_pred[x.Left_available][x.Up_available][SIZE_16]
	} else {
		fn = pred[mode][SIZE_16]
	}
	fn(ypred_ptr, int64(y_stride), yabove_row, &yleft_col[0])
}
func vp8_build_intra_predictors_mbuv_s(x *MacroBlockd, uabove_row *uint8, vabove_row *uint8, uleft *uint8, vleft *uint8, left_stride int, upred_ptr *uint8, vpred_ptr *uint8, pred_stride int) {
	var (
		uvmode    int = int(x.Mode_info_context.Mbmi.Uv_mode)
		uleft_col [8]uint8
		vleft_col [8]uint8
		i         int
		fn        intra_pred_fn
	)
	for i = 0; i < 8; i++ {
		uleft_col[i] = *(*uint8)(unsafe.Add(unsafe.Pointer(uleft), i*left_stride))
		vleft_col[i] = *(*uint8)(unsafe.Add(unsafe.Pointer(vleft), i*left_stride))
	}
	if uvmode == int(DC_PRED) {
		fn = dc_pred[x.Left_available][x.Up_available][SIZE_8]
	} else {
		fn = pred[uvmode][SIZE_8]
	}
	fn(upred_ptr, int64(pred_stride), uabove_row, &uleft_col[0])
	fn(vpred_ptr, int64(pred_stride), vabove_row, &vleft_col[0])
}
func vp8_init_intra_predictors() {
	ports.Once(vp8_init_intra_predictors_internal)
}
