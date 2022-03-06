package vpx

import "unsafe"

const VPX_TS_MAX_PERIODICITY = 16
const VPX_TS_MAX_LAYERS = 5
const VPX_MAX_LAYERS = 12
const VPX_SS_MAX_LAYERS = 5
const VPX_SS_DEFAULT_LAYERS = 1
const VPX_ENCODER_ABI_VERSION = 25
const VPX_CODEC_CAP_PSNR = 0x10000
const VPX_CODEC_CAP_OUTPUT_PARTITION = 0x20000
const VPX_CODEC_USE_PSNR = 0x10000
const VPX_CODEC_USE_OUTPUT_PARTITION = 0x20000
const VPX_CODEC_USE_HIGHBITDEPTH = 0x40000
const VPX_FRAME_IS_KEY = 1
const VPX_FRAME_IS_DROPPABLE = 2
const VPX_FRAME_IS_INVISIBLE = 4
const VPX_FRAME_IS_FRAGMENT = 8
const VPX_ERROR_RESILIENT_DEFAULT = 1
const VPX_ERROR_RESILIENT_PARTITIONS = 2
const VPX_EFLAG_FORCE_KF = 1
const VPX_DL_REALTIME = 1
const VPX_DL_GOOD_QUALITY = 1000000
const VPX_DL_BEST_QUALITY = 0

type vpx_fixed_buf struct {
	Buf unsafe.Pointer
	Sz  uint64
}
type FixedBuf vpx_fixed_buf
type CodecPts int64
type CodecFrameFlags uint32
type CodecErFlags uint32
type CodecCxPktKind int

const (
	VPX_CODEC_CX_FRAME_PKT   CodecCxPktKind = 0
	VPX_CODEC_STATS_PKT      CodecCxPktKind = 1
	VPX_CODEC_FPMB_STATS_PKT CodecCxPktKind = 2
	VPX_CODEC_PSNR_PKT       CodecCxPktKind = 3
	VPX_CODEC_CUSTOM_PKT     CodecCxPktKind = 256
)

type CodecCxPkt struct {
	Kind CodecCxPktKind
	Data struct {
		// union
		Frame struct {
			Buf                   unsafe.Pointer
			Sz                    uint64
			Pts                   CodecPts
			Duration              uint
			Flags                 CodecFrameFlags
			Partition_id          int
			Width                 [5]uint
			Height                [5]uint
			Spatial_layer_encoded [5]uint8
		}
		Twopass_stats      FixedBuf
		Firstpass_mb_stats FixedBuf
		Psnr               vpx_psnr_pkt
		Raw                FixedBuf
		Pad                [120]byte
	}
}
type CodecEncOutputCxPktCbFn func(pkt *CodecCxPkt, user_data unsafe.Pointer)
type CodecEncOutputCxCbPair struct {
	Output_cx_pkt CodecEncOutputCxPktCbFn
	User_priv     unsafe.Pointer
}
type CodecPvtOutputCxPktCbPair CodecEncOutputCxCbPair
type vpx_rational struct {
	Num int
	Den int
}
type Rational vpx_rational
type vpx_enc_pass int

const (
	VPX_RC_ONE_PASS = vpx_enc_pass(iota)
	VPX_RC_FIRST_PASS
	VPX_RC_LAST_PASS
)

type vpx_rc_mode int

const (
	VPX_VBR = vpx_rc_mode(iota)
	VPX_CBR
	VPX_CQ
	VPX_Q
)

type vpx_kf_mode int

const (
	VPX_KF_FIXED    vpx_kf_mode = 0
	VPX_KF_AUTO     vpx_kf_mode = 1
	VPX_KF_DISABLED vpx_kf_mode = 0
)

type vpx_enc_frame_flags_t int
type CodecEncCfg struct {
	G_usage                         uint
	G_threads                       uint
	G_profile                       uint
	G_w                             uint
	G_h                             uint
	G_bit_depth                     BitDepth
	G_input_bit_depth               uint
	G_timebase                      vpx_rational
	G_error_resilient               CodecErFlags
	G_pass                          vpx_enc_pass
	G_lag_in_frames                 uint
	Rc_dropframe_thresh             uint
	Rc_resize_allowed               uint
	Rc_scaled_width                 uint
	Rc_scaled_height                uint
	Rc_resize_up_thresh             uint
	Rc_resize_down_thresh           uint
	Rc_end_usage                    vpx_rc_mode
	Rc_twopass_stats_in             FixedBuf
	Rc_firstpass_mb_stats_in        FixedBuf
	Rc_target_bitrate               uint
	Rc_min_quantizer                uint
	Rc_max_quantizer                uint
	Rc_undershoot_pct               uint
	Rc_overshoot_pct                uint
	Rc_buf_sz                       uint
	Rc_buf_initial_sz               uint
	Rc_buf_optimal_sz               uint
	Rc_2pass_vbr_bias_pct           uint
	Rc_2pass_vbr_minsection_pct     uint
	Rc_2pass_vbr_maxsection_pct     uint
	Rc_2pass_vbr_corpus_complexity  uint
	Kf_mode                         vpx_kf_mode
	Kf_min_dist                     uint
	Kf_max_dist                     uint
	Ss_number_layers                uint
	Ss_enable_auto_alt_ref          [5]int
	Ss_target_bitrate               [5]uint
	Ts_number_layers                uint
	Ts_target_bitrate               [5]uint
	Ts_rate_decimator               [5]uint
	Ts_periodicity                  uint
	Ts_layer_id                     [16]uint
	Layer_target_bitrate            [12]uint
	Temporal_layering_mode          int
	Use_vizier_rc_params            int
	Active_wq_factor                Rational
	Err_per_mb_factor               Rational
	Sr_default_decay_limit          Rational
	Sr_diff_factor                  Rational
	Kf_err_per_mb_factor            Rational
	Kf_frame_min_boost_factor       Rational
	Kf_frame_max_boost_first_factor Rational
	Kf_frame_max_boost_subs_factor  Rational
	Kf_max_total_boost_factor       Rational
	Gf_max_total_boost_factor       Rational
	Gf_frame_max_boost_factor       Rational
	Zm_factor                       Rational
	Rd_mult_inter_qp_fac            Rational
	Rd_mult_arf_qp_fac              Rational
	Rd_mult_key_qp_fac              Rational
}
type SvcParameters struct {
	Max_quantizers         [12]int
	Min_quantizers         [12]int
	Scaling_factor_num     [12]int
	Scaling_factor_den     [12]int
	Speed_per_layer        [12]int
	Temporal_layering_mode int
	Loopfilter_ctrl        [12]int
}
type SvcExtraCfg SvcParameters
