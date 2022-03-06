package scale

import "github.com/mearaj/libvpx/internal/vpx"

const VP8BORDERINPIXELS = 32
const VP9INNERBORDERINPIXELS = 96
const VP9_INTERP_EXTEND = 4
const VP9_ENC_BORDER_IN_PIXELS = 160
const VP9_DEC_BORDER_IN_PIXELS = 32
const YV12_FLAG_HIGHBITDEPTH = 8

type Yv12BufferConfig struct {
	Y_width         int
	Y_height        int
	Y_crop_width    int
	Y_crop_height   int
	Y_stride        int
	Uv_width        int
	Uv_height       int
	Uv_crop_width   int
	Uv_crop_height  int
	Uv_stride       int
	Alpha_width     int
	Alpha_height    int
	Alpha_stride    int
	Y_buffer        *uint8
	U_buffer        *uint8
	V_buffer        *uint8
	Alpha_buffer    *uint8
	Buffer_alloc    *uint8
	Buffer_alloc_sz uint64
	Border          int
	Frame_size      uint64
	Subsampling_x   int
	Subsampling_y   int
	Bit_depth       uint
	Color_space     vpx.ColorSpace
	Color_range     vpx_color_range_t
	Render_width    int
	Render_height   int
	Corrupted       int
	Flags           int
}
