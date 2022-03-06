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

func initialize_dec() {
	var init_done int = 0
	if init_done == 0 {
		dsp.DspRtcd()
		vp8_init_intra_predictors()
		init_done = 1
	}
}
func remove_decompressor(pbi *VP8D_COMP) {
	vp8_remove_common(&pbi.Common)
	mem.VpxFree(unsafe.Pointer(pbi))
}
func create_decompressor(oxcf *VP8D_CONFIG) *VP8D_COMP {
	var pbi *VP8D_COMP = (*VP8D_COMP)(mem.VpxMemAlign(32, uint64(unsafe.Sizeof(VP8D_COMP{}))))
	if pbi == nil {
		return nil
	}
	*pbi = VP8D_COMP{}
	if pbi.Common.Error.Jmp.SetJump() != 0 {
		pbi.Common.Error.Setjmp = 0
		remove_decompressor(pbi)
		return nil
	}
	pbi.Common.Error.Setjmp = 1
	vp8_create_common(&pbi.Common)
	pbi.Common.Current_video_frame = 0
	pbi.Ready_for_new_data = 1
	vp8cx_init_de_quantizer(pbi)
	vp8_loop_filter_init(&pbi.Common)
	pbi.Common.Error.Setjmp = 0
	_ = oxcf
	pbi.Ec_enabled = 0
	pbi.Ec_active = 0
	pbi.Decoded_key_frame = 0
	pbi.Independent_partitions = 0
	vp8_setup_block_dptrs(&pbi.Mb)
	ports.Once(initialize_dec)
	return pbi
}
func vp8dx_get_reference(pbi *VP8D_COMP, ref_frame_flag vpx.VpxRefFrameType, sd *scale.Yv12BufferConfig) vpx.CodecErr {
	var (
		cm         *VP8Common = &pbi.Common
		ref_fb_idx int
	)
	if ref_frame_flag == vpx.VpxRefFrameType(VP8_LAST_FRAME) {
		ref_fb_idx = cm.Lst_fb_idx
	} else if ref_frame_flag == vpx.VpxRefFrameType(VP8_GOLD_FRAME) {
		ref_fb_idx = cm.Gld_fb_idx
	} else if ref_frame_flag == vpx.VpxRefFrameType(VP8_ALTR_FRAME) {
		ref_fb_idx = cm.Alt_fb_idx
	} else {
		vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(VPX_CODEC_ERROR), libc.CString("Invalid reference frame"))
		return pbi.Common.Error.Error_code
	}
	if cm.Yv12_fb[ref_fb_idx].Y_height != sd.Y_height || cm.Yv12_fb[ref_fb_idx].Y_width != sd.Y_width || cm.Yv12_fb[ref_fb_idx].Uv_height != sd.Uv_height || cm.Yv12_fb[ref_fb_idx].Uv_width != sd.Uv_width {
		vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(VPX_CODEC_ERROR), libc.CString("Incorrect buffer dimensions"))
	} else {
		scale.Vp8Yv12CopyFrameC((*scale.Yv12BufferConfig)(unsafe.Pointer(&cm.Yv12_fb[ref_fb_idx])), (*scale.Yv12BufferConfig)(unsafe.Pointer(sd)))
	}
	return pbi.Common.Error.Error_code
}
func vp8dx_set_reference(pbi *VP8D_COMP, ref_frame_flag vpx.VpxRefFrameType, sd *scale.Yv12BufferConfig) vpx.CodecErr {
	var (
		cm         *VP8Common = &pbi.Common
		ref_fb_ptr *int       = nil
		free_fb    int
	)
	if ref_frame_flag == vpx.VpxRefFrameType(VP8_LAST_FRAME) {
		ref_fb_ptr = &cm.Lst_fb_idx
	} else if ref_frame_flag == vpx.VpxRefFrameType(VP8_GOLD_FRAME) {
		ref_fb_ptr = &cm.Gld_fb_idx
	} else if ref_frame_flag == vpx.VpxRefFrameType(VP8_ALTR_FRAME) {
		ref_fb_ptr = &cm.Alt_fb_idx
	} else {
		vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(VPX_CODEC_ERROR), libc.CString("Invalid reference frame"))
		return pbi.Common.Error.Error_code
	}
	if cm.Yv12_fb[*ref_fb_ptr].Y_height != sd.Y_height || cm.Yv12_fb[*ref_fb_ptr].Y_width != sd.Y_width || cm.Yv12_fb[*ref_fb_ptr].Uv_height != sd.Uv_height || cm.Yv12_fb[*ref_fb_ptr].Uv_width != sd.Uv_width {
		vpx.InternalError(&pbi.Common.Error, vpx.CodecErr(VPX_CODEC_ERROR), libc.CString("Incorrect buffer dimensions"))
	} else {
		free_fb = get_free_fb(cm)
		cm.Fb_idx_ref_cnt[free_fb]--
		ref_cnt_fb(&cm.Fb_idx_ref_cnt[0], ref_fb_ptr, free_fb)
		scale.Vp8Yv12CopyFrameC((*scale.Yv12BufferConfig)(unsafe.Pointer(sd)), (*scale.Yv12BufferConfig)(unsafe.Pointer(&cm.Yv12_fb[*ref_fb_ptr])))
	}
	return pbi.Common.Error.Error_code
}
func get_free_fb(cm *VP8Common) int {
	var i int
	for i = 0; i < NUM_YV12_BUFFERS; i++ {
		if cm.Fb_idx_ref_cnt[i] == 0 {
			break
		}
	}
	libc.Assert(i < NUM_YV12_BUFFERS)
	cm.Fb_idx_ref_cnt[i] = 1
	return i
}
func ref_cnt_fb(buf *int, idx *int, new_idx int) {
	if *(*int)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int(0))*uintptr(*idx))) > 0 {
		*(*int)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int(0))*uintptr(*idx)))--
	}
	*idx = new_idx
	*(*int)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int(0))*uintptr(new_idx)))++
}
func swap_frame_buffers(cm *VP8Common) int {
	var err int = 0
	if cm.Copy_buffer_to_arf != 0 {
		var new_fb int = 0
		if cm.Copy_buffer_to_arf == 1 {
			new_fb = cm.Lst_fb_idx
		} else if cm.Copy_buffer_to_arf == 2 {
			new_fb = cm.Gld_fb_idx
		} else {
			err = -1
		}
		ref_cnt_fb(&cm.Fb_idx_ref_cnt[0], &cm.Alt_fb_idx, new_fb)
	}
	if cm.Copy_buffer_to_gf != 0 {
		var new_fb int = 0
		if cm.Copy_buffer_to_gf == 1 {
			new_fb = cm.Lst_fb_idx
		} else if cm.Copy_buffer_to_gf == 2 {
			new_fb = cm.Alt_fb_idx
		} else {
			err = -1
		}
		ref_cnt_fb(&cm.Fb_idx_ref_cnt[0], &cm.Gld_fb_idx, new_fb)
	}
	if cm.Refresh_golden_frame != 0 {
		ref_cnt_fb(&cm.Fb_idx_ref_cnt[0], &cm.Gld_fb_idx, cm.New_fb_idx)
	}
	if cm.Refresh_alt_ref_frame != 0 {
		ref_cnt_fb(&cm.Fb_idx_ref_cnt[0], &cm.Alt_fb_idx, cm.New_fb_idx)
	}
	if cm.Refresh_last_frame != 0 {
		ref_cnt_fb(&cm.Fb_idx_ref_cnt[0], &cm.Lst_fb_idx, cm.New_fb_idx)
		cm.Frame_to_show = &cm.Yv12_fb[cm.Lst_fb_idx]
	} else {
		cm.Frame_to_show = &cm.Yv12_fb[cm.New_fb_idx]
	}
	cm.Fb_idx_ref_cnt[cm.New_fb_idx]--
	return err
}
func check_fragments_for_errors(pbi *VP8D_COMP) int {
	if pbi.Ec_active == 0 && pbi.Fragments.Count <= 1 && pbi.Fragments.Sizes[0] == 0 {
		var cm *VP8Common = &pbi.Common
		if cm.Fb_idx_ref_cnt[cm.Lst_fb_idx] > 1 {
			var prev_idx int = cm.Lst_fb_idx
			cm.Fb_idx_ref_cnt[prev_idx]--
			cm.Lst_fb_idx = get_free_fb(cm)
			scale.Vp8Yv12CopyFrameC((*scale.Yv12BufferConfig)(unsafe.Pointer(&cm.Yv12_fb[prev_idx])), (*scale.Yv12BufferConfig)(unsafe.Pointer(&cm.Yv12_fb[cm.Lst_fb_idx])))
		}
		cm.Yv12_fb[cm.Lst_fb_idx].Corrupted = 1
		cm.Show_frame = 0
		return 0
	}
	return 1
}
func vp8dx_receive_compressed_data(pbi *VP8D_COMP, time_stamp int64) int {
	var (
		cm      *VP8Common = &pbi.Common
		retcode int        = -1
	)
	pbi.Common.Error.Error_code = vpx.CodecErr(VPX_CODEC_OK)
	retcode = check_fragments_for_errors(pbi)
	if retcode <= 0 {
		return retcode
	}
	cm.New_fb_idx = get_free_fb(cm)
	pbi.Dec_fb_ref[INTRA_FRAME] = &cm.Yv12_fb[cm.New_fb_idx]
	pbi.Dec_fb_ref[LAST_FRAME] = &cm.Yv12_fb[cm.Lst_fb_idx]
	pbi.Dec_fb_ref[GOLDEN_FRAME] = &cm.Yv12_fb[cm.Gld_fb_idx]
	pbi.Dec_fb_ref[ALTREF_FRAME] = &cm.Yv12_fb[cm.Alt_fb_idx]
	retcode = vp8_decode_frame(pbi)
	if retcode < 0 {
		if cm.Fb_idx_ref_cnt[cm.New_fb_idx] > 0 {
			cm.Fb_idx_ref_cnt[cm.New_fb_idx]--
		}
		pbi.Common.Error.Error_code = vpx.CodecErr(VPX_CODEC_ERROR)
		if pbi.Mb.Error_info.Error_code != 0 {
			pbi.Common.Error.Error_code = pbi.Mb.Error_info.Error_code
			libc.MemCpy(unsafe.Pointer(&pbi.Common.Error.Detail[0]), unsafe.Pointer(&pbi.Mb.Error_info.Detail[0]), int(unsafe.Sizeof([80]byte{})))
		}
		goto decode_exit
	}
	if swap_frame_buffers(cm) != 0 {
		pbi.Common.Error.Error_code = vpx.CodecErr(VPX_CODEC_ERROR)
		goto decode_exit
	}
	ports.ClearSystemState()
	if cm.Show_frame != 0 {
		cm.Current_video_frame++
		cm.Show_frame_mi = cm.Mi
	}
	pbi.Ready_for_new_data = 0
	pbi.Last_time_stamp = time_stamp
decode_exit:
	ports.ClearSystemState()
	return retcode
}
func vp8dx_get_raw_frame(pbi *VP8D_COMP, sd *scale.Yv12BufferConfig, time_stamp *int64, time_end_stamp *int64, flags *Vp8PpFlags) int {
	var ret int = -1
	if pbi.Ready_for_new_data == 1 {
		return ret
	}
	if pbi.Common.Show_frame == 0 {
		return ret
	}
	pbi.Ready_for_new_data = 1
	*time_stamp = pbi.Last_time_stamp
	*time_end_stamp = 0
	ret = vp8_post_proc_frame(&pbi.Common, sd, flags)
	ports.ClearSystemState()
	return ret
}
func vp8dx_references_buffer(oci *VP8Common, ref_frame int) int {
	var (
		mi     *ModeInfo = oci.Mi
		mb_row int
		mb_col int
	)
	for mb_row = 0; mb_row < oci.Mb_rows; mb_row++ {
		for mb_col = 0; mb_col < oci.Mb_cols; func() *ModeInfo {
			mb_col++
			return func() *ModeInfo {
				p := &mi
				x := *p
				*p = (*ModeInfo)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(ModeInfo{})*1))
				return x
			}()
		}() {
			if int(mi.Mbmi.Ref_frame) == ref_frame {
				return 1
			}
		}
		mi = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mi), unsafe.Sizeof(ModeInfo{})*1))
	}
	return 0
}
func vp8_create_decoder_instances(fb *frame_buffers, oxcf *VP8D_CONFIG) int {
	fb.Pbi[0] = create_decompressor(oxcf)
	if fb.Pbi[0] == nil {
		return VPX_CODEC_ERROR
	}
	if fb.Pbi[0].Common.Error.Jmp.SetJump() != 0 {
		vp8_remove_decoder_instances(fb)
		*(*[32]*VP8D_COMP)(unsafe.Pointer(&fb.Pbi[0])) = [32]*VP8D_COMP{}
		ports.ClearSystemState()
		return VPX_CODEC_ERROR
	}
	fb.Pbi[0].Common.Error.Setjmp = 1
	fb.Pbi[0].Max_threads = oxcf.Max_threads
	DecoderCreateThreads(fb.Pbi[0])
	fb.Pbi[0].Common.Error.Setjmp = 0
	return VPX_CODEC_OK
}
func vp8_remove_decoder_instances(fb *frame_buffers) int {
	var pbi *VP8D_COMP = fb.Pbi[0]
	if pbi == nil {
		return VPX_CODEC_ERROR
	}
	vp8_decoder_remove_threads(pbi)
	remove_decompressor(pbi)
	return VPX_CODEC_OK
}
func vp8dx_get_quantizer(pbi *VP8D_COMP) int {
	return pbi.Common.Base_qindex
}
