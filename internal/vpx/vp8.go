package vpx

type Vp8ComControlId int

const (
	VP8_SET_REFERENCE         Vp8ComControlId = 1
	VP8_COPY_REFERENCE        Vp8ComControlId = 2
	VP8_SET_POSTPROC          Vp8ComControlId = 3
	VP9_GET_REFERENCE         Vp8ComControlId = 128
	VP8_COMMON_CTRL_ID_MAX    Vp8ComControlId = 129
	VP8_DECODER_CTRL_ID_START Vp8ComControlId = 256
)

type Vp8PostProcLevel int

const (
	VP8_NOFILTERING  Vp8PostProcLevel = 0
	VP8_DEBLOCK      Vp8PostProcLevel = 1 << 0
	VP8_DEMACROBLOCK Vp8PostProcLevel = 1 << 1
	VP8_ADDNOISE     Vp8PostProcLevel = 1 << 2
	VP8_MFQE         Vp8PostProcLevel = 1 << 3
)

type Vp8PostProcCfg struct {
	Post_proc_flag   int
	Deblocking_level int
	Noise_level      int
}
type VpxRefFrameType int

const (
	VP8_LAST_FRAME VpxRefFrameType = 1
	VP8_GOLD_FRAME VpxRefFrameType = 2
	VP8_ALTR_FRAME VpxRefFrameType = 4
)

type RefFrame struct {
	Frame_type VpxRefFrameType
	Img        Image
}
type Vp9RefFrame struct {
	Idx int
	Img Image
}

func Vp8SetReference(ctx *CodecCtx, CtrlId int, data *RefFrame) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8CopyReference(ctx *CodecCtx, CtrlId int, data *RefFrame) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8SetPostProc(ctx *CodecCtx, CtrlId int, data *Vp8PostProcCfg) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9GetReference(ctx *CodecCtx, CtrlId int, data *Vp9RefFrame) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
