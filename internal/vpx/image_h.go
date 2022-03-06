package vpx

import "unsafe"

const VPX_IMAGE_ABI_VERSION = 5
const VPX_IMG_FMT_PLANAR = 256
const VPX_IMG_FMT_UV_FLIP = 512
const VPX_IMG_FMT_HAS_ALPHA = 1024
const VPX_IMG_FMT_HIGHBITDEPTH = 2048
const VPX_PLANE_PACKED = 0
const VPX_PLANE_Y = 0
const VPX_PLANE_U = 1
const VPX_PLANE_V = 2
const VPX_PLANE_ALPHA = 3

type ImgFmt int

const (
	VPX_IMG_FMT_NONE   ImgFmt = 0
	VPX_IMG_FMT_YV12   ImgFmt = ImgFmt(int(VPX_IMG_FMT_PLANAR|VPX_IMG_FMT_UV_FLIP) | 1)
	VPX_IMG_FMT_I420   ImgFmt = ImgFmt(int(VPX_IMG_FMT_PLANAR | 2))
	VPX_IMG_FMT_I422   ImgFmt = ImgFmt(int(VPX_IMG_FMT_PLANAR | 5))
	VPX_IMG_FMT_I444   ImgFmt = ImgFmt(int(VPX_IMG_FMT_PLANAR | 6))
	VPX_IMG_FMT_I440   ImgFmt = ImgFmt(int(VPX_IMG_FMT_PLANAR | 7))
	VPX_IMG_FMT_NV12   ImgFmt = ImgFmt(int(VPX_IMG_FMT_PLANAR | 9))
	VPX_IMG_FMT_I42016 ImgFmt = ImgFmt(VPX_IMG_FMT_I420 | VPX_IMG_FMT_HIGHBITDEPTH)
	VPX_IMG_FMT_I42216 ImgFmt = ImgFmt(VPX_IMG_FMT_I422 | VPX_IMG_FMT_HIGHBITDEPTH)
	VPX_IMG_FMT_I44416 ImgFmt = ImgFmt(VPX_IMG_FMT_I444 | VPX_IMG_FMT_HIGHBITDEPTH)
	VPX_IMG_FMT_I44016 ImgFmt = ImgFmt(VPX_IMG_FMT_I440 | VPX_IMG_FMT_HIGHBITDEPTH)
)

type ColorSpace int

const (
	VPX_CS_UNKNOWN   ColorSpace = 0
	VPX_CS_BT_601    ColorSpace = 1
	VPX_CS_BT_709    ColorSpace = 2
	VPX_CS_SMPTE_170 ColorSpace = 3
	VPX_CS_SMPTE_240 ColorSpace = 4
	VPX_CS_BT_2020   ColorSpace = 5
	VPX_CS_RESERVED  ColorSpace = 6
	VPX_CS_SRGB      ColorSpace = 7
)

type ColorRange int

const (
	VPX_CR_STUDIO_RANGE ColorRange = 0
	VPX_CR_FULL_RANGE   ColorRange = 1
)

type Image struct {
	Fmt            ImgFmt
	Cs             ColorSpace
	Range          ColorRange
	W              uint
	H              uint
	Bit_depth      uint
	D_w            uint
	D_h            uint
	R_w            uint
	R_h            uint
	X_chroma_shift uint
	Y_chroma_shift uint
	Planes         [4]*uint8
	Stride         [4]int
	Bps            int
	User_priv      unsafe.Pointer
	Img_data       *uint8
	Img_data_owner int
	Self_allocd    int
	Fb_priv        unsafe.Pointer
}
type vpx_image_rect struct {
	X uint
	Y uint
	W uint
	H uint
}
type ImageRect vpx_image_rect
