package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/dsp"
	"github.com/mearaj/libvpx/internal/mem"
	"github.com/mearaj/libvpx/internal/ports"
	"github.com/mearaj/libvpx/internal/scale"
	"github.com/mearaj/libvpx/internal/util"
	"github.com/mearaj/libvpx/internal/vpx"
	"unsafe"
)

const VP8_CAP_POSTPROC = 0x40000
const VP8_CAP_ERROR_CONCEALMENT = 0

type MemSegId int

const (
	VP8_SEG_ALG_PRIV = MemSegId(iota + 256)
	VP8_SEG_MAX
)

type CodecAlgPvt struct {
	Base               vpx.CodecPriv
	Cfg                vpx.CodecDecCfg
	Si                 vpx.CodecStreamInfo
	Decoder_init       int
	Restart_threads    int
	Postproc_cfg_set   int
	Postproc_cfg       vpx.Vp8PostProcCfg
	Decrypt_cb         vpx.DecryptCb
	Decrypt_state      unsafe.Pointer
	Img                vpx.Image
	Img_setup          int
	Yv12_frame_buffers frame_buffers
	User_priv          unsafe.Pointer
	Fragments          FRAGMENT_DATA
}

func vp8_init_ctx(ctx *vpx.CodecCtx) int {
	var priv *CodecAlgPvt = (*CodecAlgPvt)(mem.VpxCalloc(1, uint64(unsafe.Sizeof(CodecAlgPvt{}))))
	if priv == nil {
		return 1
	}
	ctx.Priv = &priv.Base
	ctx.Priv.Init_flags = ctx.Init_flags
	priv.Si.Sz = uint(unsafe.Sizeof(vpx.CodecStreamInfo{}))
	priv.Decrypt_cb = nil
	priv.Decrypt_state = nil
	if ctx.Config.Dec != nil {
		priv.Cfg = vpx.CodecDecCfg(*ctx.Config.Dec)
		ctx.Config.Dec = (*vpx.CodecDecCfg)(unsafe.Pointer(&priv.Cfg))
	}
	return 0
}
func vp8_init(ctx *vpx.CodecCtx, data *vpx.CodecPvtEncMrCfg) vpx.CodecErr {
	var res vpx.CodecErr = vpx.CodecErr(VPX_CODEC_OK)
	_ = data
	Vp8Rtcd()
	dsp.DspRtcd()
	scale.ScaleRtcd()
	if ctx.Priv == nil {
		var priv *CodecAlgPvt
		if vp8_init_ctx(ctx) != 0 {
			return vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR)
		}
		priv = (*CodecAlgPvt)(unsafe.Pointer(ctx.Priv))
		priv.Fragments.Count = 0
		priv.Fragments.Enabled = int(priv.Base.Init_flags & VPX_CODEC_USE_INPUT_FRAGMENTS)
	}
	return res
}
func vp8_destroy(ctx *CodecAlgPvt) vpx.CodecErr {
	vp8_remove_decoder_instances(&ctx.Yv12_frame_buffers)
	mem.VpxFree(unsafe.Pointer(ctx))
	return vpx.CodecErr(VPX_CODEC_OK)
}
func vp8_peek_si_internal(data *uint8, data_sz uint, si *vpx.CodecStreamInfo, decrypt_cb vpx.DecryptCb, decrypt_state unsafe.Pointer) vpx.CodecErr {
	var res vpx.CodecErr = vpx.CodecErr(VPX_CODEC_OK)
	libc.Assert(data != nil)
	if uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), data_sz)))) <= uintptr(unsafe.Pointer(data)) {
		res = vpx.CodecErr(vpx.VPX_CODEC_INVALID_PARAM)
	} else {
		var (
			clear_buffer [10]uint8
			clear        *uint8 = data
		)
		if decrypt_cb != nil {
			var n int = int(func() uintptr {
				if (unsafe.Sizeof([10]uint8{})) < uintptr(data_sz) {
					return unsafe.Sizeof([10]uint8{})
				}
				return uintptr(data_sz)
			}())
			decrypt_cb(decrypt_state, data, &clear_buffer[0], n)
			clear = &clear_buffer[0]
		}
		si.Is_kf = 0
		if data_sz >= 10 && (int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 0)))&1) == 0 {
			si.Is_kf = 1
			if int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 3))) != 157 || int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 4))) != 1 || int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 5))) != 42 {
				return vpx.CodecErr(vpx.VPX_CODEC_UNSUP_BITSTREAM)
			}
			si.W = uint((int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 6))) | int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 7)))<<8) & 0x3FFF)
			si.H = uint((int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 8))) | int(*(*uint8)(unsafe.Add(unsafe.Pointer(clear), 9)))<<8) & 0x3FFF)
			if !(si.H != 0 && si.W != 0) {
				res = vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME)
			}
		} else {
			res = vpx.CodecErr(vpx.VPX_CODEC_UNSUP_BITSTREAM)
		}
	}
	return res
}
func vp8_peek_si(data *uint8, data_sz uint, si *vpx.CodecStreamInfo) vpx.CodecErr {
	return vp8_peek_si_internal(data, data_sz, si, nil, nil)
}
func vp8_get_si(ctx *CodecAlgPvt, si *vpx.CodecStreamInfo) vpx.CodecErr {
	var sz uint
	if si.Sz >= uint(unsafe.Sizeof(vpx.CodecStreamInfo{})) {
		sz = uint(unsafe.Sizeof(vpx.CodecStreamInfo{}))
	} else {
		sz = uint(unsafe.Sizeof(vpx.CodecStreamInfo{}))
	}
	libc.MemCpy(unsafe.Pointer(si), unsafe.Pointer(&ctx.Si), int(sz))
	si.Sz = sz
	return vpx.CodecErr(VPX_CODEC_OK)
}
func yuvconfig2image(img *vpx.Image, yv12 *scale.Yv12BufferConfig, user_priv unsafe.Pointer) {
	img.Fmt = vpx.ImgFmt(vpx.VPX_IMG_FMT_I420)
	img.W = uint(yv12.Y_stride)
	img.H = uint((yv12.Y_height + int(scale.VP8BORDERINPIXELS*2) + 15) & ^int(15))
	img.D_w = func() uint {
		p := &img.R_w
		img.R_w = uint(yv12.Y_width)
		return *p
	}()
	img.D_h = func() uint {
		p := &img.R_h
		img.R_h = uint(yv12.Y_height)
		return *p
	}()
	img.X_chroma_shift = 1
	img.Y_chroma_shift = 1
	img.Planes[vpx.VPX_PLANE_Y] = yv12.Y_buffer
	img.Planes[vpx.VPX_PLANE_U] = yv12.U_buffer
	img.Planes[vpx.VPX_PLANE_V] = yv12.V_buffer
	img.Planes[vpx.VPX_PLANE_ALPHA] = nil
	img.Stride[vpx.VPX_PLANE_Y] = yv12.Y_stride
	img.Stride[vpx.VPX_PLANE_U] = yv12.Uv_stride
	img.Stride[vpx.VPX_PLANE_V] = yv12.Uv_stride
	img.Stride[vpx.VPX_PLANE_ALPHA] = yv12.Y_stride
	img.Bit_depth = 8
	img.Bps = 12
	img.User_priv = user_priv
	img.Img_data = yv12.Buffer_alloc
	img.Img_data_owner = 0
	img.Self_allocd = 0
}
func update_fragments(ctx *CodecAlgPvt, data *uint8, data_sz uint, res *vpx.CodecErr) int {
	*res = vpx.CodecErr(VPX_CODEC_OK)
	if ctx.Fragments.Count == 0 {
		*(*[9]*uint8)(unsafe.Pointer(&ctx.Fragments.Ptrs[0])) = [9]*uint8{}
		*(*[9]uint)(unsafe.Pointer(&ctx.Fragments.Sizes[0])) = [9]uint{}
	}
	if ctx.Fragments.Enabled != 0 && !(data == nil && data_sz == 0) {
		ctx.Fragments.Ptrs[ctx.Fragments.Count] = data
		ctx.Fragments.Sizes[ctx.Fragments.Count] = data_sz
		ctx.Fragments.Count++
		if ctx.Fragments.Count > uint((1<<EIGHT_PARTITION)+1) {
			ctx.Fragments.Count = 0
			*res = vpx.CodecErr(vpx.VPX_CODEC_INVALID_PARAM)
			return -1
		}
		return 0
	}
	if ctx.Fragments.Enabled == 0 && (data == nil && data_sz == 0) {
		return 0
	}
	if ctx.Fragments.Enabled == 0 {
		ctx.Fragments.Ptrs[0] = data
		ctx.Fragments.Sizes[0] = data_sz
		ctx.Fragments.Count = 1
	}
	return 1
}
func vp8_decode(ctx *CodecAlgPvt, data *uint8, data_sz uint, user_priv unsafe.Pointer, deadline int) vpx.CodecErr {
	var (
		res               vpx.CodecErr
		resolution_change uint = 0
		w                 uint
		h                 uint
	)
	if ctx.Fragments.Enabled == 0 && (data == nil && data_sz == 0) {
		return 0
	}
	if update_fragments(ctx, data, data_sz, &res) <= 0 {
		return res
	}
	w = ctx.Si.W
	h = ctx.Si.H
	res = vp8_peek_si_internal(ctx.Fragments.Ptrs[0], ctx.Fragments.Sizes[0], &ctx.Si, ctx.Decrypt_cb, ctx.Decrypt_state)
	if res == vpx.CodecErr(vpx.VPX_CODEC_UNSUP_BITSTREAM) && ctx.Si.Is_kf == 0 {
		res = vpx.CodecErr(VPX_CODEC_OK)
	}
	if ctx.Decoder_init == 0 && ctx.Si.Is_kf == 0 {
		res = vpx.CodecErr(vpx.VPX_CODEC_UNSUP_BITSTREAM)
	}
	if ctx.Si.H != h || ctx.Si.W != w {
		resolution_change = 1
	}
	if res == 0 && ctx.Restart_threads != 0 {
		var (
			fb  *frame_buffers = &ctx.Yv12_frame_buffers
			pbi *VP8D_COMP     = ctx.Yv12_frame_buffers.Pbi[0]
			pc  *VP8Common     = &pbi.Common
		)
		if pbi.Common.Error.Jmp.SetJump() != 0 {
			vp8_remove_decoder_instances(fb)
			*(*[32]*VP8D_COMP)(unsafe.Pointer(&fb.Pbi[0])) = [32]*VP8D_COMP{}
			ports.ClearSystemState()
			return vpx.CodecErr(VPX_CODEC_ERROR)
		}
		pbi.Common.Error.Setjmp = 1
		pbi.Max_threads = int(ctx.Cfg.Threads)
		DecoderCreateThreads(pbi)
		if util.AtomicLoadAcquire(&pbi.B_multithreaded_rd) != 0 {
			vp8mt_alloc_temp_buffers(pbi, pc.Width, pc.Mb_rows)
		}
		ctx.Restart_threads = 0
		pbi.Common.Error.Setjmp = 0
	}
	if res == 0 && ctx.Decoder_init == 0 {
		var oxcf VP8D_CONFIG
		oxcf.Width = int(ctx.Si.W)
		oxcf.Height = int(ctx.Si.H)
		oxcf.Version = 9
		oxcf.Postprocess = 0
		oxcf.Max_threads = int(ctx.Cfg.Threads)
		oxcf.Error_concealment = int(ctx.Base.Init_flags & VPX_CODEC_USE_ERROR_CONCEALMENT)
		if ctx.Postproc_cfg_set == 0 && (ctx.Base.Init_flags&vpx.VPX_CODEC_USE_POSTPROC) != 0 {
			ctx.Postproc_cfg.Post_proc_flag = VP8_DEBLOCK | VP8_DEMACROBLOCK | VP8_MFQE
			ctx.Postproc_cfg.Deblocking_level = 4
			ctx.Postproc_cfg.Noise_level = 0
		}
		res = vpx.CodecErr(vp8_create_decoder_instances(&ctx.Yv12_frame_buffers, &oxcf))
		if res == vpx.CodecErr(VPX_CODEC_OK) {
			ctx.Decoder_init = 1
		}
	}
	if ctx.Decoder_init != 0 {
		ctx.Yv12_frame_buffers.Pbi[0].Decrypt_cb = ctx.Decrypt_cb
		ctx.Yv12_frame_buffers.Pbi[0].Decrypt_state = ctx.Decrypt_state
	}
	if res == 0 {
		var (
			pbi *VP8D_COMP = ctx.Yv12_frame_buffers.Pbi[0]
			pc  *VP8Common = &pbi.Common
		)
		if resolution_change != 0 {
			var (
				xd *MacroBlockd = &pbi.Mb
				i  int
			)
			pc.Width = int(ctx.Si.W)
			pc.Height = int(ctx.Si.H)
			{
				var prev_mb_rows int = pc.Mb_rows
				if pbi.Common.Error.Jmp.SetJump() != 0 {
					pbi.Common.Error.Setjmp = 0
					ctx.Si.W = 0
					ctx.Si.H = 0
					ports.ClearSystemState()
					return -1
				}
				pbi.Common.Error.Setjmp = 1
				if pc.Width <= 0 {
					pc.Width = int(w)
					vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Invalid frame width"))
				}
				if pc.Height <= 0 {
					pc.Height = int(h)
					vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME), libc.CString("Invalid frame height"))
				}
				if vp8_alloc_frame_buffers(pc, pc.Width, pc.Height) != 0 {
					vpx.InternalError(&pc.Error, vpx.CodecErr(vpx.VPX_CODEC_MEM_ERROR), libc.CString("Failed to allocate frame buffers"))
				}
				xd.Pre = pc.Yv12_fb[pc.Lst_fb_idx]
				xd.Dst = pc.Yv12_fb[pc.New_fb_idx]
				for i = 0; i < pbi.Allocated_decoding_thread_count; i++ {
					(*(*MB_ROW_DEC)(unsafe.Add(unsafe.Pointer(pbi.Mb_row_di), unsafe.Sizeof(MB_ROW_DEC{})*uintptr(i)))).Mbd.Dst = pc.Yv12_fb[pc.New_fb_idx]
					vp8_build_block_doffsets(&(*(*MB_ROW_DEC)(unsafe.Add(unsafe.Pointer(pbi.Mb_row_di), unsafe.Sizeof(MB_ROW_DEC{})*uintptr(i)))).Mbd)
				}
				vp8_build_block_doffsets(&pbi.Mb)
				if util.AtomicLoadAcquire(&pbi.B_multithreaded_rd) != 0 {
					vp8mt_alloc_temp_buffers(pbi, pc.Width, prev_mb_rows)
				}
			}
			pbi.Common.Error.Setjmp = 0
			pbi.Common.Fb_idx_ref_cnt[0] = 0
		}
		if pbi.Common.Error.Jmp.SetJump() != 0 {
			ports.ClearSystemState()
			pc.Yv12_fb[pc.Lst_fb_idx].Corrupted = 1
			if pc.Fb_idx_ref_cnt[pc.New_fb_idx] > 0 {
				pc.Fb_idx_ref_cnt[pc.New_fb_idx]--
			}
			pc.Error.Setjmp = 0
			if pbi.Restart_threads != 0 {
				ctx.Si.W = 0
				ctx.Si.H = 0
				ctx.Restart_threads = 1
			}
			res = UpdateErrorState(ctx, &pbi.Common.Error)
			return res
		}
		pbi.Common.Error.Setjmp = 1
		pbi.Fragments = ctx.Fragments
		pbi.Restart_threads = 0
		ctx.User_priv = user_priv
		if vp8dx_receive_compressed_data(pbi, int64(deadline)) != 0 {
			res = UpdateErrorState(ctx, &pbi.Common.Error)
		}
		ctx.Fragments.Count = 0
	}
	return res
}
func vp8_get_frame(ctx *CodecAlgPvt, iter *vpx.CodecIter) *vpx.Image {
	var img *vpx.Image = nil
	if (*iter) == nil && ctx.Yv12_frame_buffers.Pbi[0] != nil {
		var (
			sd             scale.Yv12BufferConfig
			time_stamp     int64 = 0
			time_end_stamp int64 = 0
			flags          Vp8PpFlags
		)
		flags = Vp8PpFlags{}
		if ctx.Base.Init_flags&vpx.VPX_CODEC_USE_POSTPROC != 0 {
			flags.Post_proc_flag = ctx.Postproc_cfg.Post_proc_flag
			flags.Deblocking_level = ctx.Postproc_cfg.Deblocking_level
			flags.Noise_level = ctx.Postproc_cfg.Noise_level
		}
		if vp8dx_get_raw_frame(ctx.Yv12_frame_buffers.Pbi[0], &sd, &time_stamp, &time_end_stamp, &flags) == 0 {
			yuvconfig2image(&ctx.Img, &sd, ctx.User_priv)
			img = &ctx.Img
			*iter = vpx.CodecIter(img)
		}
	}
	return img
}
func image2yuvconfig(img *vpx.Image, yv12 *scale.Yv12BufferConfig) vpx.CodecErr {
	var (
		y_w  int          = int(img.D_w)
		y_h  int          = int(img.D_h)
		uv_w int          = int((img.D_w + 1) / 2)
		uv_h int          = int((img.D_h + 1) / 2)
		res  vpx.CodecErr = vpx.CodecErr(VPX_CODEC_OK)
	)
	yv12.Y_buffer = img.Planes[vpx.VPX_PLANE_Y]
	yv12.U_buffer = img.Planes[vpx.VPX_PLANE_U]
	yv12.V_buffer = img.Planes[vpx.VPX_PLANE_V]
	yv12.Y_crop_width = y_w
	yv12.Y_crop_height = y_h
	yv12.Y_width = y_w
	yv12.Y_height = y_h
	yv12.Uv_crop_width = uv_w
	yv12.Uv_crop_height = uv_h
	yv12.Uv_width = uv_w
	yv12.Uv_height = uv_h
	yv12.Y_stride = img.Stride[vpx.VPX_PLANE_Y]
	yv12.Uv_stride = img.Stride[vpx.VPX_PLANE_U]
	yv12.Border = (img.Stride[vpx.VPX_PLANE_Y] - int(img.D_w)) / 2
	return res
}
func vp8_set_reference(ctx *CodecAlgPvt, args libc.ArgList) vpx.CodecErr {
	var data *vpx.RefFrame = args.Arg().(*vpx.RefFrame)
	if data != nil {
		var (
			frame *vpx.RefFrame = data
			sd    scale.Yv12BufferConfig
		)
		image2yuvconfig(&frame.Img, &sd)
		return vp8dx_set_reference(ctx.Yv12_frame_buffers.Pbi[0], vpx.VpxRefFrameType(frame.Frame_type), &sd)
	} else {
		return vpx.CodecErr(vpx.VPX_CODEC_INVALID_PARAM)
	}
}
func vp8_get_reference(ctx *CodecAlgPvt, args libc.ArgList) vpx.CodecErr {
	var data *vpx.RefFrame = args.Arg().(*vpx.RefFrame)
	if data != nil {
		var (
			frame *vpx.RefFrame = data
			sd    scale.Yv12BufferConfig
		)
		image2yuvconfig(&frame.Img, &sd)
		return vp8dx_get_reference(ctx.Yv12_frame_buffers.Pbi[0], vpx.VpxRefFrameType(frame.Frame_type), &sd)
	} else {
		return vpx.CodecErr(vpx.VPX_CODEC_INVALID_PARAM)
	}
}
func vp8_get_quantizer(ctx *CodecAlgPvt, args libc.ArgList) vpx.CodecErr {
	var (
		arg *int       = args.Arg().(*int)
		pbi *VP8D_COMP = ctx.Yv12_frame_buffers.Pbi[0]
	)
	if arg == nil {
		return vpx.CodecErr(vpx.VPX_CODEC_INVALID_PARAM)
	}
	if pbi == nil {
		return vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME)
	}
	*arg = vp8dx_get_quantizer(pbi)
	return vpx.CodecErr(VPX_CODEC_OK)
}
func vp8_set_postproc(ctx *CodecAlgPvt, args libc.ArgList) vpx.CodecErr {
	var data *vpx.Vp8PostProcCfg = args.Arg().(*vpx.Vp8PostProcCfg)
	if data != nil {
		ctx.Postproc_cfg_set = 1
		ctx.Postproc_cfg = *data
		return vpx.CodecErr(VPX_CODEC_OK)
	} else {
		return vpx.CodecErr(vpx.VPX_CODEC_INVALID_PARAM)
	}
}
func vp8_get_last_ref_updates(ctx *CodecAlgPvt, args libc.ArgList) vpx.CodecErr {
	var update_info *int = args.Arg().(*int)
	if update_info != nil {
		var pbi *VP8D_COMP = ctx.Yv12_frame_buffers.Pbi[0]
		if pbi == nil {
			return vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME)
		}
		*update_info = pbi.Common.Refresh_alt_ref_frame*VP8_ALTR_FRAME + pbi.Common.Refresh_golden_frame*VP8_GOLD_FRAME + pbi.Common.Refresh_last_frame*VP8_LAST_FRAME
		return vpx.CodecErr(VPX_CODEC_OK)
	} else {
		return vpx.CodecErr(vpx.VPX_CODEC_INVALID_PARAM)
	}
}
func vp8_get_last_ref_frame(ctx *CodecAlgPvt, args libc.ArgList) vpx.CodecErr {
	var ref_info *int = args.Arg().(*int)
	if ref_info != nil {
		var pbi *VP8D_COMP = ctx.Yv12_frame_buffers.Pbi[0]
		if pbi != nil {
			var oci *VP8Common = &pbi.Common
			*ref_info = (func() int {
				if vp8dx_references_buffer(oci, int(ALTREF_FRAME)) != 0 {
					return VP8_ALTR_FRAME
				}
				return 0
			}()) | (func() int {
				if vp8dx_references_buffer(oci, int(GOLDEN_FRAME)) != 0 {
					return VP8_GOLD_FRAME
				}
				return 0
			}()) | (func() int {
				if vp8dx_references_buffer(oci, int(LAST_FRAME)) != 0 {
					return VP8_LAST_FRAME
				}
				return 0
			}())
			return vpx.CodecErr(VPX_CODEC_OK)
		} else {
			return vpx.CodecErr(vpx.VPX_CODEC_CORRUPT_FRAME)
		}
	} else {
		return vpx.CodecErr(vpx.VPX_CODEC_INVALID_PARAM)
	}
}
func vp8_get_frame_corrupted(ctx *CodecAlgPvt, args libc.ArgList) vpx.CodecErr {
	var (
		corrupted *int       = args.Arg().(*int)
		pbi       *VP8D_COMP = ctx.Yv12_frame_buffers.Pbi[0]
	)
	if corrupted != nil && pbi != nil {
		var frame *scale.Yv12BufferConfig = pbi.Common.Frame_to_show
		if frame == nil {
			return vpx.CodecErr(VPX_CODEC_ERROR)
		}
		*corrupted = frame.Corrupted
		return vpx.CodecErr(VPX_CODEC_OK)
	} else {
		return vpx.CodecErr(vpx.VPX_CODEC_INVALID_PARAM)
	}
}
func vp8_set_decryptor(ctx *CodecAlgPvt, args libc.ArgList) vpx.CodecErr {
	var init *vpx.DecryptInit = args.Arg().(*vpx.DecryptInit)
	if init != nil {
		ctx.Decrypt_cb = init.Decrypt_cb
		ctx.Decrypt_state = init.Decrypt_state
	} else {
		ctx.Decrypt_cb = nil
		ctx.Decrypt_state = nil
	}
	return vpx.CodecErr(VPX_CODEC_OK)
}
func CodecDxFn() *vpx.CodecIFace {
	return &vp8DxAlgo
}
