package vp8

import (
	"github.com/mearaj/libvpx/internal/vpx"
	"math"
)

const MINQ = 0
const MAXQ = math.MaxInt8
const QINDEX_RANGE = 128
const NUM_YV12_BUFFERS = 4
const MAX_PARTITIONS = 9

type frame_contexts struct {
	Bmode_prob      [9]uint8
	Ymode_prob      [4]uint8
	Uv_mode_prob    [3]uint8
	Sub_mv_ref_prob [3]uint8
	Coef_probs      [4][8][3][11]uint8
	Mvc             [2]MV_CONTEXT
}
type FRAME_CONTEXT frame_contexts
type TOKEN_PARTITION int

const (
	ONE_PARTITION   TOKEN_PARTITION = 0
	TWO_PARTITION   TOKEN_PARTITION = 1
	FOUR_PARTITION  TOKEN_PARTITION = 2
	EIGHT_PARTITION TOKEN_PARTITION = 3
)

type CLAMP_TYPE int

const (
	RECON_CLAMP_REQUIRED    CLAMP_TYPE = 0
	RECON_CLAMP_NOTREQUIRED CLAMP_TYPE = 1
)

type VP8Common struct {
	Error                     vpx.InternalErrorInfo
	Y1dequant                 [128][2]int16
	Y2dequant                 [128][2]int16
	UVdequant                 [128][2]int16
	Width                     int
	Height                    int
	Horiz_scale               int
	Vert_scale                int
	Clamp_type                CLAMP_TYPE
	Frame_to_show             *scale.Yv12BufferConfig
	Yv12_fb                   [4]scale.Yv12BufferConfig
	Fb_idx_ref_cnt            [4]int
	New_fb_idx                int
	Lst_fb_idx                int
	Gld_fb_idx                int
	Alt_fb_idx                int
	Temp_scale_frame          scale.Yv12BufferConfig
	Post_proc_buffer          scale.Yv12BufferConfig
	Post_proc_buffer_int      scale.Yv12BufferConfig
	Post_proc_buffer_int_used int
	Pp_limits_buffer          *uint8
	Last_frame_type           int
	Frame_type                int
	Show_frame                int
	Frame_flags               int
	MBs                       int
	Mb_rows                   int
	Mb_cols                   int
	Mode_info_stride          int
	Mb_no_coeff_skip          int
	No_lpf                    int
	Use_bilinear_mc_filter    int
	Full_pixel                int
	Base_qindex               int
	Y1dc_delta_q              int
	Y2dc_delta_q              int
	Y2ac_delta_q              int
	Uvdc_delta_q              int
	Uvac_delta_q              int
	Mip                       *ModeInfo
	Mi                        *ModeInfo
	Show_frame_mi             *ModeInfo
	Filter_type               LOOPFILTERTYPE
	Lf_info                   loop_filter_info_n
	Filter_level              int
	Last_sharpness_level      int
	Sharpness_level           int
	Refresh_last_frame        int
	Refresh_golden_frame      int
	Refresh_alt_ref_frame     int
	Copy_buffer_to_gf         int
	Copy_buffer_to_arf        int
	Refresh_entropy_probs     int
	Ref_frame_sign_bias       [4]int
	Above_context             *ENTROPY_CONTEXT_PLANES
	Left_context              ENTROPY_CONTEXT_PLANES
	Lfc                       FRAME_CONTEXT
	Fc                        FRAME_CONTEXT
	Current_video_frame       uint
	Version                   int
	Multi_token_partition     TOKEN_PARTITION
	Processor_core_count      int
	Postproc_state            PostProcState
	Cpu_caps                  int
}
