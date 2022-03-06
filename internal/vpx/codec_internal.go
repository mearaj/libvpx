package vpx

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const VPX_CODEC_INTERNAL_ABI_VERSION = 5

type CodecPvtEncMrCfg struct {
	Mr_total_resolutions    uint
	Mr_encoder_id           uint
	Mr_down_sampling_factor vpx_rational
	Mr_low_res_mode_info    unsafe.Pointer
}
type CodecInitFn func(ctx *CodecCtx, data *CodecPvtEncMrCfg) CodecErr
type CodecDestroyFn func(ctx *CodecAlgPvt) CodecErr
type CodecPeekSiFn func(data *uint8, data_sz uint, si *CodecStreamInfo) CodecErr
type CodecGetSiFn func(ctx *CodecAlgPvt, si *CodecStreamInfo) CodecErr
type CodecControlFn func(ctx *CodecAlgPvt, ap libc.ArgList) CodecErr
type FnMap struct {
	Ctrl_id int
	Fn      CodecControlFn
}
type CodecDecodeFn func(ctx *CodecAlgPvt, data *uint8, data_sz uint, user_priv unsafe.Pointer, deadline int) CodecErr
type CodecGetFrameFn func(ctx *CodecAlgPvt, iter *CodecIter) *Image
type CodecSetFbFn func(ctx *CodecAlgPvt, cb_get GetFrameBufferCbFn, cb_release ReleaseFrameBufferCbFn, cb_priv unsafe.Pointer) CodecErr
type CodecEncodeFn func(ctx *CodecAlgPvt, img *Image, pts CodecPts, duration uint, flags vpx_enc_frame_flags_t, deadline uint) CodecErr
type CodecGetCxDataFn func(ctx *CodecAlgPvt, iter *CodecIter) *CodecCxPkt
type CodecEncConfigSetFn func(ctx *CodecAlgPvt, cfg *CodecEncCfg) CodecErr
type CodecGetGlobalHeadersFn func(ctx *CodecAlgPvt) *FixedBuf
type CodecGetPreviewFrameFn func(ctx *CodecAlgPvt) *Image
type CodecEncMrGetMemLocFn func(cfg *CodecEncCfg, mem_loc *unsafe.Pointer) CodecErr
type CodecEncCfgMap struct {
	Usage int
	Cfg   CodecEncCfg
}
type CodecIFace struct {
	Name        *byte
	Abi_version int
	Caps        CodecCaps
	Init        CodecInitFn
	Destroy     CodecDestroyFn
	Ctrl_maps   *FnMap
	Dec         CodecDecIFace
	Enc         CodecEncIFace
}
type CodecPvtCbPair struct {
	U struct {
		// union
		Put_frame vpx_codec_put_frame_cb_fn_t
		Put_slice vpx_codec_put_slice_cb_fn_t
	}
	User_priv unsafe.Pointer
}
type CodecPvt struct {
	Err_detail *byte
	Init_flags CodecFlags
	Dec        struct {
		Put_frame_cb CodecPvtCbPair
		Put_slice_cb CodecPvtCbPair
	}
	Enc struct {
		Cx_data_dst_buf    FixedBuf
		Cx_data_pad_before uint
		Cx_data_pad_after  uint
		Cx_data_pkt        CodecCxPkt
		Total_encoders     uint
	}
}
type CodecPktList struct {
	Cnt  uint
	Max  uint
	Pkts [1]CodecCxPkt
}
type InternalErrorInfo struct {
	Error_code CodecErr
	Has_detail int
	Detail     [80]byte
	Setjmp     int
	Jmp        libc.JumpBuf
}
