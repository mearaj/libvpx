package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/mem"
	"github.com/mearaj/libvpx/internal/scale"
	"unsafe"
)

func vp8_de_alloc_frame_buffers(oci *VP8Common) {
	var i int
	for i = 0; i < NUM_YV12_BUFFERS; i++ {
		scale.Vp8Yv12DeAllocFrameBuffer(&oci.Yv12_fb[i])
	}
	scale.Vp8Yv12DeAllocFrameBuffer(&oci.Temp_scale_frame)
	scale.Vp8Yv12DeAllocFrameBuffer(&oci.Post_proc_buffer)
	if oci.Post_proc_buffer_int_used != 0 {
		scale.Vp8Yv12DeAllocFrameBuffer(&oci.Post_proc_buffer_int)
	}
	libc.Free(unsafe.Pointer(oci.Pp_limits_buffer))
	oci.Pp_limits_buffer = nil
	libc.Free(unsafe.Pointer(oci.Postproc_state.Generated_noise))
	oci.Postproc_state.Generated_noise = nil
	libc.Free(unsafe.Pointer(oci.Above_context))
	libc.Free(unsafe.Pointer(oci.Mip))
	oci.Above_context = nil
	oci.Mip = nil
}
func vp8_alloc_frame_buffers(oci *VP8Common, width int, height int) int {
	var i int
	vp8_de_alloc_frame_buffers(oci)
	if (width & 15) != 0 {
		width += 16 - (width & 15)
	}
	if (height & 15) != 0 {
		height += 16 - (height & 15)
	}
	for i = 0; i < NUM_YV12_BUFFERS; i++ {
		oci.Fb_idx_ref_cnt[i] = 0
		oci.Yv12_fb[i].Flags = 0
		if scale.Vp8Yv12AllocFrameBuffer(&oci.Yv12_fb[i], width, height, scale.VP8BORDERINPIXELS) < 0 {
			goto allocation_fail
		}
	}
	oci.New_fb_idx = 0
	oci.Lst_fb_idx = 1
	oci.Gld_fb_idx = 2
	oci.Alt_fb_idx = 3
	oci.Fb_idx_ref_cnt[0] = 1
	oci.Fb_idx_ref_cnt[1] = 1
	oci.Fb_idx_ref_cnt[2] = 1
	oci.Fb_idx_ref_cnt[3] = 1
	if scale.Vp8Yv12AllocFrameBuffer(&oci.Temp_scale_frame, width, 16, scale.VP8BORDERINPIXELS) < 0 {
		goto allocation_fail
	}
	oci.Mb_rows = height >> 4
	oci.Mb_cols = width >> 4
	oci.MBs = oci.Mb_rows * oci.Mb_cols
	oci.Mode_info_stride = oci.Mb_cols + 1
	oci.Mip = &make([]ModeInfo, (oci.Mb_cols+1)*(oci.Mb_rows+1))[0]
	if oci.Mip == nil {
		goto allocation_fail
	}
	oci.Mi = (*ModeInfo)(unsafe.Add(unsafe.Pointer((*ModeInfo)(unsafe.Add(unsafe.Pointer(oci.Mip), unsafe.Sizeof(ModeInfo{})*uintptr(oci.Mode_info_stride)))), unsafe.Sizeof(ModeInfo{})*1))
	oci.Above_context = (*ENTROPY_CONTEXT_PLANES)(libc.Calloc(oci.Mb_cols*int(unsafe.Sizeof(ENTROPY_CONTEXT_PLANES{})), 1))
	if oci.Above_context == nil {
		goto allocation_fail
	}
	if scale.Vp8Yv12AllocFrameBuffer(&oci.Post_proc_buffer, width, height, scale.VP8BORDERINPIXELS) < 0 {
		goto allocation_fail
	}
	oci.Post_proc_buffer_int_used = 0
	oci.Postproc_state = PostProcState{}
	libc.MemSet(unsafe.Pointer(oci.Post_proc_buffer.Buffer_alloc), 128, int(oci.Post_proc_buffer.Frame_size))
	oci.Pp_limits_buffer = (*uint8)(mem.VpxMemAlign(16, uint64(((oci.Mb_cols+1) & ^int(1))*24)))
	if oci.Pp_limits_buffer == nil {
		goto allocation_fail
	}
	return 0
allocation_fail:
	vp8_de_alloc_frame_buffers(oci)
	return 1
}
func vp8_setup_version(cm *VP8Common) {
	switch cm.Version {
	case 0:
		cm.No_lpf = 0
		cm.Filter_type = LOOPFILTERTYPE(NORMAL_LOOPFILTER)
		cm.Use_bilinear_mc_filter = 0
		cm.Full_pixel = 0
	case 1:
		cm.No_lpf = 0
		cm.Filter_type = LOOPFILTERTYPE(SIMPLE_LOOPFILTER)
		cm.Use_bilinear_mc_filter = 1
		cm.Full_pixel = 0
	case 2:
		cm.No_lpf = 1
		cm.Filter_type = LOOPFILTERTYPE(NORMAL_LOOPFILTER)
		cm.Use_bilinear_mc_filter = 1
		cm.Full_pixel = 0
	case 3:
		cm.No_lpf = 1
		cm.Filter_type = LOOPFILTERTYPE(SIMPLE_LOOPFILTER)
		cm.Use_bilinear_mc_filter = 1
		cm.Full_pixel = 1
	default:
		cm.No_lpf = 0
		cm.Filter_type = LOOPFILTERTYPE(NORMAL_LOOPFILTER)
		cm.Use_bilinear_mc_filter = 0
		cm.Full_pixel = 0
	}
}
func vp8_create_common(oci *VP8Common) {
	vp8_machine_specific_config(oci)
	vp8_init_mbmode_probs(oci)
	vp8_default_bmode_probs(oci.Fc.Bmode_prob)
	oci.Mb_no_coeff_skip = 1
	oci.No_lpf = 0
	oci.Filter_type = LOOPFILTERTYPE(NORMAL_LOOPFILTER)
	oci.Use_bilinear_mc_filter = 0
	oci.Full_pixel = 0
	oci.Multi_token_partition = TOKEN_PARTITION(ONE_PARTITION)
	oci.Clamp_type = CLAMP_TYPE(RECON_CLAMP_REQUIRED)
	*(*[4]int)(unsafe.Pointer(&oci.Ref_frame_sign_bias[0])) = [4]int{}
	oci.Copy_buffer_to_gf = 0
	oci.Copy_buffer_to_arf = 0
}
func vp8_remove_common(oci *VP8Common) {
	vp8_de_alloc_frame_buffers(oci)
}
