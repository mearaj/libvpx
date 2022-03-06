package vpx

import "unsafe"

const VP8_EFLAG_NO_REF_LAST = 0x10000
const VP8_EFLAG_NO_REF_GF = 0x20000
const VP8_EFLAG_NO_REF_ARF = 0x200000
const VP8_EFLAG_NO_UPD_LAST = 0x40000
const VP8_EFLAG_NO_UPD_GF = 0x400000
const VP8_EFLAG_NO_UPD_ARF = 0x800000
const VP8_EFLAG_FORCE_GF = 0x80000
const VP8_EFLAG_FORCE_ARF = 0x1000000
const VP8_EFLAG_NO_UPD_ENTROPY = 0x100000

type Vp8eEncControlId int

const (
	VP8E_SET_ROI_MAP                    Vp8eEncControlId = 8
	VP8E_SET_ACTIVEMAP                  Vp8eEncControlId = 9
	VP8E_SET_SCALEMODE                  Vp8eEncControlId = 11
	VP8E_SET_CPUUSED                    Vp8eEncControlId = 13
	VP8E_SET_ENABLEAUTOALTREF           Vp8eEncControlId = 14
	VP8E_SET_NOISE_SENSITIVITY          Vp8eEncControlId = 15
	VP8E_SET_SHARPNESS                  Vp8eEncControlId = 16
	VP8E_SET_STATIC_THRESHOLD           Vp8eEncControlId = 17
	VP8E_SET_TOKEN_PARTITIONS           Vp8eEncControlId = 18
	VP8E_GET_LAST_QUANTIZER             Vp8eEncControlId = 19
	VP8E_GET_LAST_QUANTIZER_64          Vp8eEncControlId = 20
	VP8E_SET_ARNR_MAXFRAMES             Vp8eEncControlId = 21
	VP8E_SET_ARNR_STRENGTH              Vp8eEncControlId = 22
	VP8E_SET_ARNR_TYPE                  Vp8eEncControlId = 23
	VP8E_SET_TUNING                     Vp8eEncControlId = 24
	VP8E_SET_CQ_LEVEL                   Vp8eEncControlId = 25
	VP8E_SET_MAX_INTRA_BITRATE_PCT      Vp8eEncControlId = 26
	VP8E_SET_FRAME_FLAGS                Vp8eEncControlId = 27
	VP9E_SET_MAX_INTER_BITRATE_PCT      Vp8eEncControlId = 28
	VP9E_SET_GF_CBR_BOOST_PCT           Vp8eEncControlId = 29
	VP8E_SET_TEMPORAL_LAYER_ID          Vp8eEncControlId = 30
	VP8E_SET_SCREEN_CONTENT_MODE        Vp8eEncControlId = 31
	VP9E_SET_LOSSLESS                   Vp8eEncControlId = 32
	VP9E_SET_TILE_COLUMNS               Vp8eEncControlId = 33
	VP9E_SET_TILE_ROWS                  Vp8eEncControlId = 34
	VP9E_SET_FRAME_PARALLEL_DECODING    Vp8eEncControlId = 35
	VP9E_SET_AQ_MODE                    Vp8eEncControlId = 36
	VP9E_SET_FRAME_PERIODIC_BOOST       Vp8eEncControlId = 37
	VP9E_SET_NOISE_SENSITIVITY          Vp8eEncControlId = 38
	VP9E_SET_SVC                        Vp8eEncControlId = 39
	VP9E_SET_ROI_MAP                    Vp8eEncControlId = 40
	VP9E_SET_SVC_PARAMETERS             Vp8eEncControlId = 41
	VP9E_SET_SVC_LAYER_ID               Vp8eEncControlId = 42
	VP9E_SET_TUNE_CONTENT               Vp8eEncControlId = 43
	VP9E_GET_SVC_LAYER_ID               Vp8eEncControlId = 44
	VP9E_REGISTER_CX_CALLBACK           Vp8eEncControlId = 45
	VP9E_SET_COLOR_SPACE                Vp8eEncControlId = 46
	VP9E_SET_TEMPORAL_LAYERING_MODE     Vp8eEncControlId = 47
	VP9E_SET_MIN_GF_INTERVAL            Vp8eEncControlId = 48
	VP9E_SET_MAX_GF_INTERVAL            Vp8eEncControlId = 49
	VP9E_GET_ACTIVEMAP                  Vp8eEncControlId = 50
	VP9E_SET_COLOR_RANGE                Vp8eEncControlId = 51
	VP9E_SET_SVC_REF_FRAME_CONFIG       Vp8eEncControlId = 52
	VP9E_SET_RENDER_SIZE                Vp8eEncControlId = 53
	VP9E_SET_TARGET_LEVEL               Vp8eEncControlId = 54
	VP9E_SET_ROW_MT                     Vp8eEncControlId = 55
	VP9E_GET_LEVEL                      Vp8eEncControlId = 56
	VP9E_SET_ALT_REF_AQ                 Vp8eEncControlId = 57
	VP8E_SET_GF_CBR_BOOST_PCT           Vp8eEncControlId = 58
	VP9E_ENABLE_MOTION_VECTOR_UNIT_TEST Vp8eEncControlId = 59
	VP9E_SET_SVC_INTER_LAYER_PRED       Vp8eEncControlId = 60
	VP9E_SET_SVC_FRAME_DROP_LAYER       Vp8eEncControlId = 61
	VP9E_GET_SVC_REF_FRAME_CONFIG       Vp8eEncControlId = 62
	VP9E_SET_SVC_GF_TEMPORAL_REF        Vp8eEncControlId = 63
	VP9E_SET_SVC_SPATIAL_LAYER_SYNC     Vp8eEncControlId = 64
	VP9E_SET_TPL                        Vp8eEncControlId = 65
	VP9E_SET_POSTENCODE_DROP            Vp8eEncControlId = 66
	VP9E_SET_DELTA_Q_UV                 Vp8eEncControlId = 67
	VP9E_SET_DISABLE_OVERSHOOT_MAXQ_CBR Vp8eEncControlId = 68
	VP9E_SET_DISABLE_LOOPFILTER         Vp8eEncControlId = 69
	VP9E_SET_EXTERNAL_RATE_CONTROL      Vp8eEncControlId = 70
	VP9E_SET_RTC_EXTERNAL_RATECTRL      Vp8eEncControlId = 71
	VP9E_GET_LOOPFILTER_LEVEL           Vp8eEncControlId = 72
	VP9E_GET_LAST_QUANTIZER_SVC_LAYERS  Vp8eEncControlId = 73
	VP8E_SET_RTC_EXTERNAL_RATECTRL      Vp8eEncControlId = 74
)

type VPX_SCALING_MODE int

const (
	VP8E_NORMAL    VPX_SCALING_MODE = 0
	VP8E_FOURFIVE  VPX_SCALING_MODE = 1
	VP8E_THREEFIVE VPX_SCALING_MODE = 2
	VP8E_ONETWO    VPX_SCALING_MODE = 3
)

type VP9E_TEMPORAL_LAYERING_MODE int

const (
	VP9E_TEMPORAL_LAYERING_MODE_NOLAYERING VP9E_TEMPORAL_LAYERING_MODE = 0
	VP9E_TEMPORAL_LAYERING_MODE_BYPASS     VP9E_TEMPORAL_LAYERING_MODE = 1
	VP9E_TEMPORAL_LAYERING_MODE_0101       VP9E_TEMPORAL_LAYERING_MODE = 2
	VP9E_TEMPORAL_LAYERING_MODE_0212       VP9E_TEMPORAL_LAYERING_MODE = 3
)

type RoiMap struct {
	Enabled          uint8
	Roi_map          *uint8
	Rows             uint
	Cols             uint
	Delta_q          [8]int
	Delta_lf         [8]int
	Skip             [8]int
	Ref_frame        [8]int
	Static_threshold [4]uint
}
type vpx_active_map struct {
	Active_map *uint8
	Rows       uint
	Cols       uint
}
type ScalingMode struct {
	H_scaling_mode VPX_SCALING_MODE
	V_scaling_mode VPX_SCALING_MODE
}
type Vp8eTokenPartitions int

const (
	VP8_ONE_TOKENPARTITION   Vp8eTokenPartitions = 0
	VP8_TWO_TOKENPARTITION   Vp8eTokenPartitions = 1
	VP8_FOUR_TOKENPARTITION  Vp8eTokenPartitions = 2
	VP8_EIGHT_TOKENPARTITION Vp8eTokenPartitions = 3
)

type Vp9eTuneContent int

const (
	VP9E_CONTENT_DEFAULT = Vp9eTuneContent(iota)
	VP9E_CONTENT_SCREEN
	VP9E_CONTENT_FILM
	VP9E_CONTENT_INVALID
)

type Vp8eTuning int

const (
	VP8_TUNE_PSNR = Vp8eTuning(iota)
	VP8_TUNE_SSIM
)

type SvcLayerId struct {
	Spatial_layer_id              int
	Temporal_layer_id             int
	Temporal_layer_id_per_spatial [5]int
}
type SvcRefFrameConfig struct {
	Lst_fb_idx         [5]int
	Gld_fb_idx         [5]int
	Alt_fb_idx         [5]int
	Update_buffer_slot [5]int
	Update_last        [5]int
	Update_golden      [5]int
	Update_alt_ref     [5]int
	Reference_last     [5]int
	Reference_golden   [5]int
	Reference_alt_ref  [5]int
	Duration           [5]int64
}
type SVC_LAYER_DROP_MODE int

const (
	CONSTRAINED_LAYER_DROP = SVC_LAYER_DROP_MODE(iota)
	LAYER_DROP
	FULL_SUPERFRAME_DROP
	CONSTRAINED_FROM_ABOVE_DROP
)

type SvcFrameDrop struct {
	Framedrop_thresh [5]int
	Framedrop_mode   SVC_LAYER_DROP_MODE
	Max_consec_drop  int
}
type SvcSpatialLayerSync struct {
	Spatial_layer_sync    [5]int
	Base_layer_intra_only int
}

func Vp8eSetRoiMap(ctx *CodecCtx, CtrlId int, data *RoiMap) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetActiveMap(ctx *CodecCtx, CtrlId int, data *vpx_active_map) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetScaleMode(ctx *CodecCtx, CtrlId int, data *ScalingMode) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetCpuUsed(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetEnableAutoAltRef(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetNoiseSensitivity(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetSharpness(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetStaticThreshold(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetTokenPartitions(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eGetLastQuantizer(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eGetLastQuantizer64(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetArnrMaxFrames(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetArnrStrength(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetArnrType(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetTuning(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetCqLevel(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetMaxIntraBitratePct(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetFrameFlags(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetMaxInterBitratePct(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetGfCbrBoostPct(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetTemporalLayerId(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetScreenContentMode(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetLossless(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetTileColumns(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetTileRows(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetFrameParallelDecoding(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetAqMode(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetFramePeriodicBoost(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetNoiseSensitivity(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetSvc(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetRoiMap(ctx *CodecCtx, CtrlId int, data *RoiMap) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetSvcParameters(ctx *CodecCtx, CtrlId int, data unsafe.Pointer) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetSvcLayerId(ctx *CodecCtx, CtrlId int, data *SvcLayerId) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetTuneContent(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eGetSvcLayerId(ctx *CodecCtx, CtrlId int, data *SvcLayerId) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eRegisterCxCallback(ctx *CodecCtx, CtrlId int, data unsafe.Pointer) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetColorSpace(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetTemporalLayeringMode(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetMinGfInterval(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetMaxGfInterval(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eGetActiveMap(ctx *CodecCtx, CtrlId int, data *vpx_active_map) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetColorRange(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetSvcRefFrameConfig(ctx *CodecCtx, CtrlId int, data *SvcRefFrameConfig) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetRenderSize(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetTargetLevel(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetRowMt(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eGetLevel(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetAltRefAq(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetGfCbrBoostPct(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eEnableMotionVectorUnitTest(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetSvcInterLayerPred(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetSvcFrameDropLayer(ctx *CodecCtx, CtrlId int, data *SvcFrameDrop) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eGetSvcRefFrameConfig(ctx *CodecCtx, CtrlId int, data *SvcRefFrameConfig) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetSvcGfTemporalRef(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetSvcSpatialLayerSync(ctx *CodecCtx, CtrlId int, data *SvcSpatialLayerSync) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetTpl(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetPostEncodeDrop(ctx *CodecCtx, CtrlId int, data uint) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetDeltaQUv(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetDisableOverShootMaxqCbr(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetDisableLoopFilter(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetExternalRateControl(ctx *CodecCtx, CtrlId int, data *RcFuncs) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eSetRtcExternalRateCtrl(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eGetLoopFilterLevel(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp9eGetLastQuantizerSvcLayers(ctx *CodecCtx, CtrlId int, data *int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
func Vp8eSetRtcExternalRateCtrl(ctx *CodecCtx, CtrlId int, data int) CodecErr {
	return CodecControl(ctx, CtrlId, data)
}
