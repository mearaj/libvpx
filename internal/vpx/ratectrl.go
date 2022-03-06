package vpx

import "unsafe"

const VPX_EXT_RATECTRL_ABI_VERSION = 1

type vpx_rc_model_t unsafe.Pointer
type RcEncodeFrameDecision struct {
	Q_index        int
	Max_frame_size int
}
type RcEncodeFrameInfo struct {
	Frame_type               int
	Show_index               int
	Coding_index             int
	Gop_index                int
	Ref_frame_coding_indexes [3]int
	Ref_frame_valid_list     [3]int
}
type RcEncodeFrameResult struct {
	Sse                    int64
	Bit_count              int64
	Pixel_count            int64
	Actual_encoding_qindex int
}
type RcStatus int

const (
	VPX_RC_OK    RcStatus = 0
	VPX_RC_ERROR RcStatus = 1
)

type RcFrameStats struct {
	Frame              float64
	Weight             float64
	Intra_error        float64
	Coded_error        float64
	Sr_coded_error     float64
	Frame_noise_energy float64
	Pcnt_inter         float64
	Pcnt_motion        float64
	Pcnt_second_ref    float64
	Pcnt_neutral       float64
	Pcnt_intra_low     float64
	Pcnt_intra_high    float64
	Intra_skip_pct     float64
	Intra_smooth_pct   float64
	Inactive_zone_rows float64
	Inactive_zone_cols float64
	MVr                float64
	Mvr_abs            float64
	MVc                float64
	Mvc_abs            float64
	MVrv               float64
	MVcv               float64
	Mv_in_out_count    float64
	Duration           float64
	Count              float64
}
type RcFirstPassStats struct {
	Frame_stats *RcFrameStats
	Num_frames  int
}
type RcConfig struct {
	Frame_width         int
	Frame_height        int
	Show_frame_count    int
	Target_bitrate_kbps int
	Frame_rate_num      int
	Frame_rate_den      int
}
type RcCreateModelCbFn func(priv unsafe.Pointer, ratectrl_config *RcConfig, rate_ctrl_model_pt *vpx_rc_model_t) RcStatus
type RcSendFirstpassStatsCbFn func(rate_ctrl_model vpx_rc_model_t, first_pass_stats *RcFirstPassStats) RcStatus
type RcGetEncodeframeDecisionCbFn func(rate_ctrl_model vpx_rc_model_t, encode_frame_info *RcEncodeFrameInfo, frame_decision *RcEncodeFrameDecision) RcStatus
type RcUpdateEncodeframeResultCbFn func(rate_ctrl_model vpx_rc_model_t, encode_frame_result *RcEncodeFrameResult) RcStatus
type RcDeleteModelCbFn func(rate_ctrl_model vpx_rc_model_t) RcStatus
type RcFuncs struct {
	Create_model              RcCreateModelCbFn
	Send_firstpass_stats      RcSendFirstpassStatsCbFn
	Get_encodeframe_decision  RcGetEncodeframeDecisionCbFn
	Update_encodeframe_result RcUpdateEncodeframeResultCbFn
	Delete_model              RcDeleteModelCbFn
	Priv                      unsafe.Pointer
}
