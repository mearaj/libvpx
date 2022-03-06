package vp8

const (
	VP8D_NOFILTERING  = 0
	VP8D_DEBLOCK      = 1 << 0
	VP8D_DEMACROBLOCK = 1 << 1
	VP8D_ADDNOISE     = 1 << 2
	VP8D_MFQE         = 1 << 3
)

type Vp8PpFlags struct {
	Post_proc_flag         int
	Deblocking_level       int
	Noise_level            int
	Display_ref_frame_flag int
	Display_mb_modes_flag  int
	Display_b_modes_flag   int
	Display_mv_flag        int
}
