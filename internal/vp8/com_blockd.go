package vp8

import (
	"github.com/mearaj/libvpx/internal/vpx"
	"unsafe"
)

const DCPREDSIMTHRESH = 0
const DCPREDCNTTHRESH = 3
const MB_FEATURE_TREE_PROBS = 3
const MAX_MB_SEGMENTS = 4
const MAX_REF_LF_DELTAS = 4
const MAX_MODE_LF_DELTAS = 4
const SEGMENT_DELTADATA = 0
const SEGMENT_ABSDATA = 1
const PLANE_TYPE_Y_NO_DC = 0
const PLANE_TYPE_Y2 = 1
const PLANE_TYPE_UV = 2
const PLANE_TYPE_Y_WITH_DC = 3
const SEGMENT_ALTQ = 1
const SEGMENT_ALT_LF = 2

type POS struct {
	R int
	C int
}
type ENTROPY_CONTEXT_PLANES struct {
	Y1 [4]byte
	U  [2]byte
	V  [2]byte
	Y2 int8
}

const (
	KEY_FRAME   int = 0
	INTER_FRAME int = 1
)
const (
	DC_PRED = int(iota)
	V_PRED
	H_PRED
	TM_PRED
	B_PRED
	NEARESTMV
	NEARMV
	ZEROMV
	NEWMV
	SPLITMV
	MB_MODE_COUNT
)
const (
	MB_LVL_ALT_Q  int = 0
	MB_LVL_ALT_LF int = 1
	MB_LVL_MAX    int = 2
)

type B_PREDICTION_MODE int

const (
	B_DC_PRED = B_PREDICTION_MODE(iota)
	B_TM_PRED
	B_VE_PRED
	B_HE_PRED
	B_LD_PRED
	B_RD_PRED
	B_VR_PRED
	B_VL_PRED
	B_HD_PRED
	B_HU_PRED
	LEFT4X4
	ABOVE4X4
	ZERO4X4
	NEW4X4
	B_MODE_COUNT
)

type b_mode_info struct {
	// union
	As_mode B_PREDICTION_MODE
	Mv      int_mv
}

const (
	INTRA_FRAME    int = 0
	LAST_FRAME     int = 1
	GOLDEN_FRAME   int = 2
	ALTREF_FRAME   int = 3
	MAX_REF_FRAMES int = 4
)

type MB_MODE_INFO struct {
	Mode              uint8
	Uv_mode           uint8
	Ref_frame         uint8
	Is_4x4            uint8
	Mv                int_mv
	Partitioning      uint8
	Mb_skip_coeff     uint8
	Need_to_clamp_mvs uint8
	Segment_id        uint8
}
type ModeInfo struct {
	Mbmi MB_MODE_INFO
	Bmi  [16]b_mode_info
}
type Blockd struct {
	Qcoeff    *int16
	Dqcoeff   *int16
	Predictor *uint8
	Dequant   *int16
	Offset    int
	Eob       *byte
	Bmi       b_mode_info
}
type vp8_subpix_fn_t func(src_ptr *uint8, src_pixels_per_line int, xoffset int, yoffset int, dst_ptr *uint8, dst_pitch int)
type MacroBlockd struct {
	Predictor                   [384]uint8
	Qcoeff                      [400]int16
	Dqcoeff                     [400]int16
	Eobs                        [25]byte
	Dequant_y1                  [16]int16
	Dequant_y1_dc               [16]int16
	Dequant_y2                  [16]int16
	Dequant_uv                  [16]int16
	Block                       [25]Blockd
	Fullpixel_mask              int
	Pre                         scale.Yv12BufferConfig
	Dst                         scale.Yv12BufferConfig
	Mode_info_context           *ModeInfo
	Mode_info_stride            int
	Frame_type                  int
	Up_available                int
	Left_available              int
	Recon_above                 [3]*uint8
	Recon_left                  [3]*uint8
	Recon_left_stride           [2]int
	Above_context               *ENTROPY_CONTEXT_PLANES
	Left_context                *ENTROPY_CONTEXT_PLANES
	Segmentation_enabled        uint8
	Update_mb_segmentation_map  uint8
	Update_mb_segmentation_data uint8
	Mb_segement_abs_delta       uint8
	Mb_segment_tree_probs       [3]uint8
	Segment_feature_data        [2][4]int8
	Mode_ref_lf_delta_enabled   uint8
	Mode_ref_lf_delta_update    uint8
	Last_ref_lf_deltas          [4]int8
	Ref_lf_deltas               [4]int8
	Last_mode_lf_deltas         [4]int8
	Mode_lf_deltas              [4]int8
	Mb_to_left_edge             int
	Mb_to_right_edge            int
	Mb_to_top_edge              int
	Mb_to_bottom_edge           int
	Subpixel_predict            vp8_subpix_fn_t
	Subpixel_predict8x4         vp8_subpix_fn_t
	Subpixel_predict8x8         vp8_subpix_fn_t
	Subpixel_predict16x16       vp8_subpix_fn_t
	Current_bc                  unsafe.Pointer
	Corrupted                   int
	Error_info                  vpx.InternalErrorInfo
	Y_buf                       [704]uint8
}

var vp8_block2left [25]uint8 = [25]uint8{0, 0, 0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 5, 5, 6, 6, 7, 7, 8}
var vp8_block2above [25]uint8 = [25]uint8{0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 4, 5, 4, 5, 6, 7, 6, 7, 8}
