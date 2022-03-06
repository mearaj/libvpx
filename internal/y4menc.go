package internal

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/mearaj/libvpx/internal/vpx"
)

const Y4M_BUFFER_SIZE = 128

func Y4mWriteFileHeader(buf *byte, len_ uint64, width int, height int, framerate *VpxRational, fmt vpx.ImgFmt, bit_depth uint) int {
	var color *byte
	switch bit_depth {
	case 8:
		color = libc.CString(func() string {
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I444) {
				return "C444\n"
			}
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I422) {
				return "C422\n"
			}
			return "C420jpeg\n"
		}())
	case 9:
		color = libc.CString(func() string {
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I44416) {
				return "C444p9 XYSCSS=444P9\n"
			}
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I42216) {
				return "C422p9 XYSCSS=422P9\n"
			}
			return "C420p9 XYSCSS=420P9\n"
		}())
	case 10:
		color = libc.CString(func() string {
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I44416) {
				return "C444p10 XYSCSS=444P10\n"
			}
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I42216) {
				return "C422p10 XYSCSS=422P10\n"
			}
			return "C420p10 XYSCSS=420P10\n"
		}())
	case 12:
		color = libc.CString(func() string {
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I44416) {
				return "C444p12 XYSCSS=444P12\n"
			}
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I42216) {
				return "C422p12 XYSCSS=422P12\n"
			}
			return "C420p12 XYSCSS=420P12\n"
		}())
	case 14:
		color = libc.CString(func() string {
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I44416) {
				return "C444p14 XYSCSS=444P14\n"
			}
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I42216) {
				return "C422p14 XYSCSS=422P14\n"
			}
			return "C420p14 XYSCSS=420P14\n"
		}())
	case 16:
		color = libc.CString(func() string {
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I44416) {
				return "C444p16 XYSCSS=444P16\n"
			}
			if fmt == vpx.ImgFmt(VPX_IMG_FMT_I42216) {
				return "C422p16 XYSCSS=422P16\n"
			}
			return "C420p16 XYSCSS=420P16\n"
		}())
	default:
		color = nil
		libc.Assert(false)
	}
	return stdio.Snprintf(buf, int(len_), "YUV4MPEG2 W%u H%u F%u:%u I%c %s", width, height, framerate.Numerator, framerate.Denominator, 'p', color)
}
func Y4mWriteFrameHeader(buf *byte, len_ uint64) int {
	return stdio.Snprintf(buf, int(len_), "FRAME\n")
}
