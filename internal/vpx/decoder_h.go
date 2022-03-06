package vpx

import "unsafe"

const VPX_DECODER_ABI_VERSION = 12
const VPX_CODEC_CAP_PUT_SLICE = 0x10000
const VPX_CODEC_CAP_PUT_FRAME = 0x20000
const VPX_CODEC_CAP_POSTPROC = 0x40000
const VPX_CODEC_CAP_ERROR_CONCEALMENT = 0x80000
const VPX_CODEC_CAP_INPUT_FRAGMENTS = 0x100000
const VPX_CODEC_CAP_FRAME_THREADING = 0x200000
const VPX_CODEC_CAP_EXTERNAL_FRAME_BUFFER = 0x400000
const VPX_CODEC_USE_POSTPROC = 0x10000
const VPX_CODEC_USE_ERROR_CONCEALMENT = 0x20000
const VPX_CODEC_USE_INPUT_FRAGMENTS = 0x40000
const VPX_CODEC_USE_FRAME_THREADING = 0x80000

type CodecStreamInfo struct {
	Sz    uint
	W     uint
	H     uint
	Is_kf uint
}
type CodecDecCfg struct {
	Threads uint
	W       uint
	H       uint
}
type vpx_codec_put_frame_cb_fn_t func(user_priv unsafe.Pointer, img *Image)
type vpx_codec_put_slice_cb_fn_t func(user_priv unsafe.Pointer, img *Image, valid *ImageRect, update *ImageRect)
