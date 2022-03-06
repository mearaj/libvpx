package vpx

import "unsafe"

const VPX_CODEC_ABI_VERSION = 9
const VPX_CODEC_CAP_DECODER = 1
const VPX_CODEC_CAP_ENCODER = 2
const VPX_CODEC_CAP_HIGHBITDEPTH = 4

type CodecErr int

const (
	VPX_CODEC_OK = CodecErr(iota)
	VPX_CODEC_ERROR
	VPX_CODEC_MEM_ERROR
	VPX_CODEC_ABI_MISMATCH
	VPX_CODEC_INCAPABLE
	VPX_CODEC_UNSUP_BITSTREAM
	VPX_CODEC_UNSUP_FEATURE
	VPX_CODEC_CORRUPT_FRAME
	VPX_CODEC_INVALID_PARAM
	VPX_CODEC_LIST_END
)

type CodecCaps int
type CodecFlags int
type CodecPriv CodecPvt
type CodecIter unsafe.Pointer
type CodecCtx struct {
	Name       *byte
	Iface      *CodecIFace
	Err        CodecErr
	Err_detail *byte
	Init_flags CodecFlags
	Config     struct {
		// union
		Dec *CodecDecCfg
		Enc *CodecEncCfg
		Raw unsafe.Pointer
	}
	Priv *CodecPriv
}
type BitDepth int

const (
	VPX_BITS_8  BitDepth = 8
	VPX_BITS_10 BitDepth = 10
	VPX_BITS_12 BitDepth = 12
)
