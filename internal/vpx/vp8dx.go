package vpx

import "unsafe"

type Vp8DecControlId int

const (
	VP8D_GET_LAST_REF_UPDATES    Vp8DecControlId = Vp8DecControlId(VP8_DECODER_CTRL_ID_START)
	VP8D_GET_FRAME_CORRUPTED     Vp8DecControlId = 0
	VP8D_GET_LAST_REF_USED       Vp8DecControlId = 1
	VPXD_SET_DECRYPTOR           Vp8DecControlId = 2
	VP8D_SET_DECRYPTOR           Vp8DecControlId = Vp8DecControlId(VPXD_SET_DECRYPTOR)
	VP9D_GET_FRAME_SIZE          Vp8DecControlId = 3
	VP9D_GET_DISPLAY_SIZE        Vp8DecControlId = 4
	VP9D_GET_BIT_DEPTH           Vp8DecControlId = 5
	VP9_SET_BYTE_ALIGNMENT       Vp8DecControlId = 6
	VP9_INVERT_TILE_DECODE_ORDER Vp8DecControlId = 7
	VP9_SET_SKIP_LOOP_FILTER     Vp8DecControlId = 8
	VP9_DECODE_SVC_SPATIAL_LAYER Vp8DecControlId = 9
	VPXD_GET_LAST_QUANTIZER      Vp8DecControlId = 10
	VP9D_SET_ROW_MT              Vp8DecControlId = 11
	VP9D_SET_LOOP_FILTER_OPT     Vp8DecControlId = 12
	VP8_DECODER_CTRL_ID_MAX      Vp8DecControlId = 13
)

type DecryptCb func(DecryptState unsafe.Pointer, input *uint8, output *uint8, count int)
type DecryptInit struct {
	Decrypt_cb    DecryptCb
	Decrypt_state unsafe.Pointer
}

func Vp8dGetLastRefUpdates(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8dGetFrameCorrupted(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8dGetLastRefUsed(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func dSetDecryptor(ctx *CodecCtx, CtrlId int, data *DecryptInit) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8dSetDecryptor(ctx *CodecCtx, CtrlId int, data *DecryptInit) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9dGetFrameSize(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9dGetDisplaySize(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9dGetBitDepth(ctx *CodecCtx, CtrlId int, data *uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9SetByteAlignment(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9InvertTitleDecodeOrder(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9SetSkipLoopFilter(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9DecodeSvcSpatialLayer(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func VpxdGetLastQuantizer(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9dSetRowMt(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9dSetLoopFilterOpt(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
