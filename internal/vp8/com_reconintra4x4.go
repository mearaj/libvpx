package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/dsp"
	"unsafe"
)

func intra_prediction_down_copy(xd *MacroBlockd, above_right_src *uint8) {
	var (
		dst_stride      int    = xd.Dst.Y_stride
		above_right_dst *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(xd.Dst.Y_buffer), -dst_stride))), 16))
		src_ptr         *uint  = (*uint)(unsafe.Pointer(above_right_src))
		dst_ptr0        *uint  = (*uint)(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(above_right_dst), dst_stride*4))))
	)
	_ = dst_ptr0
	var dst_ptr1 *uint = (*uint)(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(above_right_dst), dst_stride*8))))
	_ = dst_ptr1
	var dst_ptr2 *uint = (*uint)(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(above_right_dst), dst_stride*12))))
	_ = dst_ptr2
	*dst_ptr0 = *src_ptr
	*dst_ptr1 = *src_ptr
	*dst_ptr2 = *src_ptr
}

var predalt [10]intra_pred_fn

func vp8_init_intra4x4_predictors_internal() {
	predalt[B_DC_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxDcPredictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	predalt[B_TM_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxTmPredictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	predalt[B_VE_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxVePredictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	predalt[B_HE_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxHePredictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	predalt[B_LD_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxD45ePredictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	predalt[B_RD_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxD135Predictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	predalt[B_VR_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxD117Predictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	predalt[B_VL_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxD63ePredictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	predalt[B_HD_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxD153Predictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
	predalt[B_HU_PRED] = func(dst *uint8, stride int64, above *uint8, left *uint8) {
		func(dst *uint8, stride int64, above *uint8, left *uint8) {
			dsp.VpxD207Predictor4x4C(dst, stride, above, left)
		}(dst, stride, above, left)
	}
}
func vp8_intra4x4_predict(above *uint8, yleft *uint8, left_stride int, b_mode B_PREDICTION_MODE, dst *uint8, dst_stride int, top_left uint8) {
	var (
		Aboveb [12]uint8
		Above  *uint8 = &Aboveb[4]
		Left   [4]uint8
	)
	Left[0] = *(*uint8)(unsafe.Add(unsafe.Pointer(yleft), 0))
	Left[1] = *(*uint8)(unsafe.Add(unsafe.Pointer(yleft), left_stride))
	Left[2] = *(*uint8)(unsafe.Add(unsafe.Pointer(yleft), left_stride*2))
	Left[3] = *(*uint8)(unsafe.Add(unsafe.Pointer(yleft), left_stride*3))
	libc.MemCpy(unsafe.Pointer(Above), unsafe.Pointer(above), 8)
	*(*uint8)(unsafe.Add(unsafe.Pointer(Above), -1)) = top_left
	predalt[b_mode]((*uint8)(unsafe.Pointer(dst)), int64(dst_stride), (*uint8)(unsafe.Pointer(Above)), (*uint8)(unsafe.Pointer(&Left[0])))
}
