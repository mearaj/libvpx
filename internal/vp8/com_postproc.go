package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/dsp"
	"github.com/mearaj/libvpx/internal/mem"
	"github.com/mearaj/libvpx/internal/ports"
	"github.com/mearaj/libvpx/internal/scale"
	"github.com/mearaj/libvpx/internal/vpx"
	"unsafe"
)

const MFQE_PRECISION = 4

type PostProcState struct {
	Last_q           int
	Last_noise       int
	Last_base_qindex int
	Last_frame_valid int
	Clamp            int
	Generated_noise  *int8
}

func q2mbl(x int) int {
	if x < 20 {
		x = 20
	}
	x = (x-50)*10/8 + 50
	return x * x / 3
}
func vp8_de_mblock(post *scale.Yv12BufferConfig, q int) {
	dsp.VpxMbPostProcAcrossIpC((*uint8)(unsafe.Pointer(post.Y_buffer)), post.Y_stride, post.Y_height, post.Y_width, q2mbl(q))
	dsp.VpxMbPostProcDownC((*uint8)(unsafe.Pointer(post.Y_buffer)), post.Y_stride, post.Y_height, post.Y_width, q2mbl(q))
}
func vp8_deblock(cm *VP8Common, source *scale.Yv12BufferConfig, post *scale.Yv12BufferConfig, q int) {
	var (
		level             float64   = float64(q)*6e-05*float64(q)*float64(q) - float64(q)*0.0067*float64(q) + float64(q)*0.306 + 0.0065
		ppl               int       = int(level + 0.5)
		mode_info_context *ModeInfo = cm.Mi
		mbr               int
		mbc               int
		ylimits           *uint8 = cm.Pp_limits_buffer
		uvlimits          *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(cm.Pp_limits_buffer), cm.Mb_cols*16))
	)
	if ppl > 0 {
		for mbr = 0; mbr < cm.Mb_rows; mbr++ {
			var (
				ylptr  *uint8 = ylimits
				uvlptr *uint8 = uvlimits
			)
			for mbc = 0; mbc < cm.Mb_cols; mbc++ {
				var mb_ppl uint8
				if mode_info_context.Mbmi.Mb_skip_coeff != 0 {
					mb_ppl = uint8(int8(int(uint8(int8(ppl))) >> 1))
				} else {
					mb_ppl = uint8(int8(ppl))
				}
				libc.MemSet(unsafe.Pointer(ylptr), byte(mb_ppl), 16)
				libc.MemSet(unsafe.Pointer(uvlptr), byte(mb_ppl), 8)
				ylptr = (*uint8)(unsafe.Add(unsafe.Pointer(ylptr), 16))
				uvlptr = (*uint8)(unsafe.Add(unsafe.Pointer(uvlptr), 8))
				mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
			}
			mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
			dsp.VpxPostProcDownAndAcrossMbRowC((*uint8)(unsafe.Add(unsafe.Pointer(source.Y_buffer), mbr*16*source.Y_stride)), (*uint8)(unsafe.Add(unsafe.Pointer(post.Y_buffer), mbr*16*post.Y_stride)), source.Y_stride, post.Y_stride, source.Y_width, ylimits, 16)
			dsp.VpxPostProcDownAndAcrossMbRowC((*uint8)(unsafe.Add(unsafe.Pointer(source.U_buffer), mbr*8*source.Uv_stride)), (*uint8)(unsafe.Add(unsafe.Pointer(post.U_buffer), mbr*8*post.Uv_stride)), source.Uv_stride, post.Uv_stride, source.Uv_width, uvlimits, 8)
			dsp.VpxPostProcDownAndAcrossMbRowC((*uint8)(unsafe.Add(unsafe.Pointer(source.V_buffer), mbr*8*source.Uv_stride)), (*uint8)(unsafe.Add(unsafe.Pointer(post.V_buffer), mbr*8*post.Uv_stride)), source.Uv_stride, post.Uv_stride, source.Uv_width, uvlimits, 8)
		}
	} else {
		scale.Vp8Yv12CopyFrameC((*scale.Yv12BufferConfig)(unsafe.Pointer(source)), (*scale.Yv12BufferConfig)(unsafe.Pointer(post)))
	}
}
func vp8_de_noise(cm *VP8Common, source *scale.Yv12BufferConfig, q int, uvfilter int) {
	var (
		mbr     int
		level   float64 = float64(q)*6e-05*float64(q)*float64(q) - float64(q)*0.0067*float64(q) + float64(q)*0.306 + 0.0065
		ppl     int     = int(level + 0.5)
		mb_rows int     = cm.Mb_rows
		mb_cols int     = cm.Mb_cols
		limits  *uint8  = cm.Pp_limits_buffer
	)
	libc.MemSet(unsafe.Pointer(limits), byte(uint8(int8(ppl))), mb_cols*16)
	for mbr = 0; mbr < mb_rows; mbr++ {
		dsp.VpxPostProcDownAndAcrossMbRowC((*uint8)(unsafe.Add(unsafe.Pointer(source.Y_buffer), mbr*16*source.Y_stride)), (*uint8)(unsafe.Add(unsafe.Pointer(source.Y_buffer), mbr*16*source.Y_stride)), source.Y_stride, source.Y_stride, source.Y_width, limits, 16)
		if uvfilter == 1 {
			dsp.VpxPostProcDownAndAcrossMbRowC((*uint8)(unsafe.Add(unsafe.Pointer(source.U_buffer), mbr*8*source.Uv_stride)), (*uint8)(unsafe.Add(unsafe.Pointer(source.U_buffer), mbr*8*source.Uv_stride)), source.Uv_stride, source.Uv_stride, source.Uv_width, limits, 8)
			dsp.VpxPostProcDownAndAcrossMbRowC((*uint8)(unsafe.Add(unsafe.Pointer(source.V_buffer), mbr*8*source.Uv_stride)), (*uint8)(unsafe.Add(unsafe.Pointer(source.V_buffer), mbr*8*source.Uv_stride)), source.Uv_stride, source.Uv_stride, source.Uv_width, limits, 8)
		}
	}
}
func vp8_post_proc_frame(oci *VP8Common, dest *scale.Yv12BufferConfig, ppflags *Vp8PpFlags) int {
	var (
		q             int = oci.Filter_level * 10 / 6
		flags         int = ppflags.Post_proc_flag
		deblock_level int = ppflags.Deblocking_level
		noise_level   int = ppflags.Noise_level
	)
	if oci.Frame_to_show == nil {
		return -1
	}
	if q > 63 {
		q = 63
	}
	if flags == 0 {
		*dest = *oci.Frame_to_show
		dest.Y_width = oci.Width
		dest.Y_height = oci.Height
		dest.Uv_height = dest.Y_height / 2
		oci.Postproc_state.Last_base_qindex = oci.Base_qindex
		oci.Postproc_state.Last_frame_valid = 1
		return 0
	}
	if flags&VP8D_ADDNOISE != 0 {
		if oci.Postproc_state.Generated_noise == nil {
			oci.Postproc_state.Generated_noise = (*int8)(mem.VpxCalloc(uint64(oci.Width+256), uint64(unsafe.Sizeof(int8(0)))))
			if oci.Postproc_state.Generated_noise == nil {
				return 1
			}
		}
	}
	if (flags&VP8D_MFQE) != 0 && oci.Post_proc_buffer_int_used == 0 {
		if (flags&VP8D_DEBLOCK) != 0 || (flags&VP8D_DEMACROBLOCK) != 0 {
			var (
				width  int = (oci.Width + 15) & ^int(15)
				height int = (oci.Height + 15) & ^int(15)
			)
			if scale.Vp8Yv12AllocFrameBuffer(&oci.Post_proc_buffer_int, width, height, scale.VP8BORDERINPIXELS) != 0 {
				vpx.InternalError(&oci.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate MFQE framebuffer"))
			}
			oci.Post_proc_buffer_int_used = 1
			libc.MemSet(unsafe.Pointer((&oci.Post_proc_buffer_int).Buffer_alloc), 128, int((&oci.Post_proc_buffer).Frame_size))
		}
	}
	ports.ClearSystemState()
	if (flags&VP8D_MFQE) != 0 && oci.Postproc_state.Last_frame_valid != 0 && oci.Current_video_frame > 10 && oci.Postproc_state.Last_base_qindex < 60 && oci.Base_qindex-oci.Postproc_state.Last_base_qindex >= 20 {
		vp8_multiframe_quality_enhance(oci)
		if ((flags&VP8D_DEBLOCK) != 0 || (flags&VP8D_DEMACROBLOCK) != 0) && oci.Post_proc_buffer_int_used != 0 {
			scale.Vp8Yv12CopyFrameC((*scale.Yv12BufferConfig)(unsafe.Pointer(&oci.Post_proc_buffer)), (*scale.Yv12BufferConfig)(unsafe.Pointer(&oci.Post_proc_buffer_int)))
			if flags&VP8D_DEMACROBLOCK != 0 {
				vp8_deblock(oci, &oci.Post_proc_buffer_int, &oci.Post_proc_buffer, q+(deblock_level-5)*10)
				vp8_de_mblock(&oci.Post_proc_buffer, q+(deblock_level-5)*10)
			} else if flags&VP8D_DEBLOCK != 0 {
				vp8_deblock(oci, &oci.Post_proc_buffer_int, &oci.Post_proc_buffer, q)
			}
		}
		oci.Postproc_state.Last_base_qindex = (oci.Postproc_state.Last_base_qindex*3 + oci.Base_qindex) >> 2
	} else if flags&VP8D_DEMACROBLOCK != 0 {
		vp8_deblock(oci, oci.Frame_to_show, &oci.Post_proc_buffer, q+(deblock_level-5)*10)
		vp8_de_mblock(&oci.Post_proc_buffer, q+(deblock_level-5)*10)
		oci.Postproc_state.Last_base_qindex = oci.Base_qindex
	} else if flags&VP8D_DEBLOCK != 0 {
		vp8_deblock(oci, oci.Frame_to_show, &oci.Post_proc_buffer, q)
		oci.Postproc_state.Last_base_qindex = oci.Base_qindex
	} else {
		scale.Vp8Yv12CopyFrameC((*scale.Yv12BufferConfig)(unsafe.Pointer(oci.Frame_to_show)), (*scale.Yv12BufferConfig)(unsafe.Pointer(&oci.Post_proc_buffer)))
		oci.Postproc_state.Last_base_qindex = oci.Base_qindex
	}
	oci.Postproc_state.Last_frame_valid = 1
	if flags&VP8D_ADDNOISE != 0 {
		if oci.Postproc_state.Last_q != q || oci.Postproc_state.Last_noise != noise_level {
			var (
				sigma   float64
				ppstate *PostProcState = &oci.Postproc_state
			)
			ports.ClearSystemState()
			sigma = float64(noise_level) + 0.5 + float64(q)*0.6/63.0
			ppstate.Clamp = dsp.VpxSetupNoise(sigma, ppstate.Generated_noise, oci.Width+256)
			ppstate.Last_q = q
			ppstate.Last_noise = noise_level
		}
		dsp.VpxPlaneAddNoiseC(oci.Post_proc_buffer.Y_buffer, oci.Postproc_state.Generated_noise, oci.Postproc_state.Clamp, oci.Postproc_state.Clamp, oci.Post_proc_buffer.Y_width, oci.Post_proc_buffer.Y_height, oci.Post_proc_buffer.Y_stride)
	}
	*dest = oci.Post_proc_buffer
	dest.Y_width = oci.Width
	dest.Y_height = oci.Height
	dest.Uv_height = dest.Y_height / 2
	return 0
}
