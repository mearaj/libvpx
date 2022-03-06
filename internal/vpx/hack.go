package vpx

type vpx_psnr_pkt interface{}
type VpxUsecTimer interface{}
type CodecAlgPvt interface{}

// Todo: uses vpx_ports/vpx_timer.h which in turn uses time.h from <sys/time.h>
func VpxUsecTimerStart(v *VpxUsecTimer) {

}

// Todo: uses vpx_ports/vpx_timer.h which in turn uses time.h from <sys/time.h>
func VpxUsecTimerMark(v *VpxUsecTimer) {

}

// Todo: uses vpx_ports/vpx_timer.h which in turn uses time.h from <sys/time.h>
func VpxUsecTimerElapsed(v *VpxUsecTimer) uint64 {
	return 0
}

type CodecDecIFace struct {
	PeekSi   CodecPeekSiFn
	GetSi    CodecGetSiFn
	Decode   CodecDecodeFn
	GetFrame CodecGetFrameFn
	SetFbFn  CodecSetFbFn
}

type CodecEncIFace struct {
	CfgMapCount int
	CfgMaps     *CodecEncCfgMap
	Encode      CodecEncodeFn
	GetCxData   CodecGetCxDataFn
	CfgSet      CodecEncConfigSetFn
	GetGlobHdrs CodecGetGlobalHeadersFn
	GetPreview  CodecGetPreviewFrameFn
	MrGetMemLoc CodecEncMrGetMemLocFn
}

func vpx_winx64_fldcw(mode uint16) {

}

func x87_set_double_precision() uint16 {
	return 0
}
