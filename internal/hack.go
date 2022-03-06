package internal

import (
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/mearaj/libvpx/internal/vpx"
	"os"
)

type FilterMode = FilterModeEnum

// Todo: uses webmdec.cc file, equivalent should be available in go
func WebmGuessFrameRate(ctx *WebmInputContext, ctx2 *VpxInputContext) int {
	return 0
}

// Todo: uses webmdec.cc file, equivalent should be available in go
func WebmFree(ctx *WebmInputContext) {

}

// Todo: uses webmdec.cc file, equivalent should be available in go
func WebmGuessFramerate(ctx *WebmInputContext, ctx2 *VpxInputContext) int {
	return 0
}

// Todo: uses webmdec.cc file, equivalent should be available in go
func FileIsWebm(ctx *WebmInputContext, ctx2 *VpxInputContext) int {
	return 0
}

// Todo: uses webmdec.cc file, equivalent should be available in go
func WebmReadFrame(ctx *WebmInputContext, buf **uint8, buffer *uint64) int {
	return 0
}

// Todo: uses third_party/libyuv/include/libyuv/libyuv_scale.h
func I420Scale(u *uint8, i int, u2 *uint8, i2 int, u3 *uint8, i3 int, i4 int, i5 int, u4 *uint8, i6 int, u5 *uint8, i7 int, u6 *uint8, i8 int, i9 int, i10 int, mode FilterMode) int {
	return 0
}

// Todo: uses stdio.h
func clearerr(file *stdio.File) {

}

// Todo: uses stdio.h
func ferror(file *stdio.File) int {
	return 0
}

// Todo: uses stdio.h
func Rewind(file *stdio.File) {

}

// Todo: uses unistd.h
func IsAtty(i int) int {
	return 0
}

// Todo: uses stdio.h
func ftello(file *stdio.File) int64 {
	return 0
}

// Todo: uses stdio.h
func fseeko(file *stdio.File, i int, i2 int64) int {
	return 0
}
func vpx_codec_vp9_dx() *vpx.CodecIFace {
	return nil
}

func vpx_codec_vp8_cx() *vpx.CodecIFace {
	return nil
}

func vpx_codec_vp9_cx() *vpx.CodecIFace {
	return nil
}

const VPX_IMG_FMT_HIGHBITDEPTH = vpx.VPX_IMG_FMT_HIGHBITDEPTH
const VPX_IMG_FMT_NV12 = vpx.VPX_IMG_FMT_NV12
const VPX_IMG_FMT_YV12 = vpx.VPX_IMG_FMT_YV12
const VPX_PLANE_V = vpx.VPX_PLANE_V
const VPX_PLANE_U = vpx.VPX_PLANE_U
const VPX_PLANE_Y = vpx.VPX_PLANE_Y
const VPX_PLANE_ALPHA = vpx.VPX_PLANE_ALPHA
const VPX_IMG_FMT_I44416 = vpx.VPX_IMG_FMT_I44416
const VPX_IMG_FMT_I42216 = vpx.VPX_IMG_FMT_I42216
const VPX_IMG_FMT_I444 = vpx.VPX_IMG_FMT_I444
const VPX_IMG_FMT_I44016 = vpx.VPX_IMG_FMT_I44016
const VPX_IMG_FMT_I420 = vpx.VPX_IMG_FMT_I420
const VPX_IMG_FMT_I422 = vpx.VPX_IMG_FMT_I422
const VPX_IMG_FMT_I42016 = vpx.VPX_IMG_FMT_I42016
const KFilterNone = kFilterNone
const KFilterLinear = kFilterLinear
const KFilterBilinear = kFilterBilinear
const KFilterBox = kFilterBox

func strtol(val *byte, i **byte, i2 int) int {
	return 0
}
func strtoul(val *byte, i **byte, i2 int) uint32 {
	return 0
}

// Todo: Implement Go Equivalent
func UsageExit() {
	// Todo: show_help as ShowHelp in golang
	// ShowHelp(stdio.Stderr(), 1)
	os.Exit(0)
}
