package vpx

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func GetAlgPvt(ctx *CodecCtx) *CodecAlgPvt {
	return (*CodecAlgPvt)(unsafe.Pointer(ctx.Priv))
}
func VpxCodecEncInitVer(ctx *CodecCtx, iface *CodecIFace, cfg *CodecEncCfg, flags CodecFlags, ver int) CodecErr {
	var res CodecErr
	if ver != (15 + (4 + 5) + 1) {
		res = CodecErr(VPX_CODEC_ABI_MISMATCH)
	} else if ctx == nil || iface == nil || cfg == nil {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if iface.Abi_version != 5 {
		res = CodecErr(VPX_CODEC_ABI_MISMATCH)
	} else if (iface.Caps & VPX_CODEC_CAP_ENCODER) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else if (flags&VPX_CODEC_USE_PSNR) != 0 && (iface.Caps&VPX_CODEC_CAP_PSNR) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else if (flags&VPX_CODEC_USE_OUTPUT_PARTITION) != 0 && (iface.Caps&VPX_CODEC_CAP_OUTPUT_PARTITION) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else {
		ctx.Iface = iface
		ctx.Name = iface.Name
		ctx.Priv = nil
		ctx.Init_flags = flags
		ctx.Config.Enc = cfg
		res = ctx.Iface.Init(ctx, nil)
		if res != 0 {
			if ctx.Priv != nil {
				ctx.Err_detail = ctx.Priv.Err_detail
			} else {
				ctx.Err_detail = nil
			}
			CodecDestroy(ctx)
		}
	}
	if ctx != nil {
		return func() CodecErr {
			p := &ctx.Err
			ctx.Err = res
			return *p
		}()
	}
	return res
}
func vpx_codec_enc_init_multi_ver(ctx *CodecCtx, iface *CodecIFace, cfg *CodecEncCfg, num_enc int, flags CodecFlags, dsf *Rational, ver int) CodecErr {
	var res CodecErr = CodecErr(VPX_CODEC_OK)
	if ver != (15 + (4 + 5) + 1) {
		res = CodecErr(VPX_CODEC_ABI_MISMATCH)
	} else if ctx == nil || iface == nil || cfg == nil || (num_enc > 16 || num_enc < 1) {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if iface.Abi_version != 5 {
		res = CodecErr(VPX_CODEC_ABI_MISMATCH)
	} else if (iface.Caps & VPX_CODEC_CAP_ENCODER) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else if (flags&VPX_CODEC_USE_PSNR) != 0 && (iface.Caps&VPX_CODEC_CAP_PSNR) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else if (flags&VPX_CODEC_USE_OUTPUT_PARTITION) != 0 && (iface.Caps&VPX_CODEC_CAP_OUTPUT_PARTITION) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else {
		var (
			i       int
			mem_loc unsafe.Pointer = nil
		)
		if iface.Enc.MrGetMemLoc == nil {
			return CodecErr(VPX_CODEC_INCAPABLE)
		}
		if (func() CodecErr {
			res = iface.Enc.MrGetMemLoc(cfg, &mem_loc)
			return res
		}()) == 0 {
			for i = 0; i < num_enc; i++ {
				var mr_cfg CodecPvtEncMrCfg
				if dsf.Num < 1 || dsf.Num > 4096 || dsf.Den < 1 || dsf.Den > dsf.Num {
					res = CodecErr(VPX_CODEC_INVALID_PARAM)
				} else {
					mr_cfg.Mr_low_res_mode_info = mem_loc
					mr_cfg.Mr_total_resolutions = uint(num_enc)
					mr_cfg.Mr_encoder_id = uint(num_enc - 1 - i)
					mr_cfg.Mr_down_sampling_factor.Num = dsf.Num
					mr_cfg.Mr_down_sampling_factor.Den = dsf.Den
					ctx.Iface = iface
					ctx.Name = iface.Name
					ctx.Priv = nil
					ctx.Init_flags = flags
					ctx.Config.Enc = cfg
					res = ctx.Iface.Init(ctx, &mr_cfg)
				}
				if res != 0 {
					var error_detail *byte
					if ctx.Priv != nil {
						error_detail = ctx.Priv.Err_detail
					} else {
						error_detail = nil
					}
					ctx.Err_detail = error_detail
					CodecDestroy(ctx)
					for i != 0 {
						ctx = (*CodecCtx)(unsafe.Add(unsafe.Pointer(ctx), -int(unsafe.Sizeof(CodecCtx{})*1)))
						ctx.Err_detail = error_detail
						CodecDestroy(ctx)
						i--
					}
					if ctx != nil {
						return func() CodecErr {
							p := &ctx.Err
							ctx.Err = res
							return *p
						}()
					}
					return res
				}
				ctx = (*CodecCtx)(unsafe.Add(unsafe.Pointer(ctx), unsafe.Sizeof(CodecCtx{})*1))
				cfg = (*CodecEncCfg)(unsafe.Add(unsafe.Pointer(cfg), unsafe.Sizeof(CodecEncCfg{})*1))
				dsf = (*Rational)(unsafe.Add(unsafe.Pointer(dsf), unsafe.Sizeof(Rational{})*1))
			}
			ctx = (*CodecCtx)(unsafe.Add(unsafe.Pointer(ctx), -int(unsafe.Sizeof(CodecCtx{})*1)))
		}
	}
	if ctx != nil {
		return func() CodecErr {
			p := &ctx.Err
			ctx.Err = res
			return *p
		}()
	}
	return res
}
func vpx_codec_enc_config_default(iface *CodecIFace, cfg *CodecEncCfg, usage uint) CodecErr {
	var res CodecErr
	if iface == nil || cfg == nil || usage != 0 {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if (iface.Caps & VPX_CODEC_CAP_ENCODER) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else {
		libc.Assert(iface.Enc.CfgMapCount == 1)
		*cfg = iface.Enc.CfgMaps.Cfg
		res = CodecErr(VPX_CODEC_OK)
	}
	return res
}
func vpx_codec_encode(ctx *CodecCtx, img *Image, pts CodecPts, duration uint, flags vpx_enc_frame_flags_t, deadline uint) CodecErr {
	var res CodecErr = CodecErr(VPX_CODEC_OK)
	if ctx == nil || img != nil && duration == 0 {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if ctx.Iface == nil || ctx.Priv == nil {
		res = CodecErr(VPX_CODEC_ERROR)
	} else if (ctx.Iface.Caps & VPX_CODEC_CAP_ENCODER) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else {
		var num_enc uint = ctx.Priv.Enc.Total_encoders
		for {
			{
				var x87_orig_mode uint16 = uint16(x87_set_double_precision())
				if num_enc == 1 {
					res = ctx.Iface.Enc.Encode(GetAlgPvt(ctx), img, pts, duration, flags, deadline)
				} else {
					var i int
					ctx = (*CodecCtx)(unsafe.Add(unsafe.Pointer(ctx), unsafe.Sizeof(CodecCtx{})*uintptr(num_enc-1)))
					if img != nil {
						img = (*Image)(unsafe.Add(unsafe.Pointer(img), unsafe.Sizeof(Image{})*uintptr(num_enc-1)))
					}
					for i = int(num_enc - 1); i >= 0; i-- {
						if (func() CodecErr {
							res = ctx.Iface.Enc.Encode(GetAlgPvt(ctx), img, pts, duration, flags, deadline)
							return res
						}()) != 0 {
							break
						}
						ctx = (*CodecCtx)(unsafe.Add(unsafe.Pointer(ctx), -int(unsafe.Sizeof(CodecCtx{})*1)))
						if img != nil {
							img = (*Image)(unsafe.Add(unsafe.Pointer(img), -int(unsafe.Sizeof(Image{})*1)))
						}
					}
					ctx = (*CodecCtx)(unsafe.Add(unsafe.Pointer(ctx), unsafe.Sizeof(CodecCtx{})*1))
				}
				vpx_winx64_fldcw(x87_orig_mode)
			}
			if true {
				break
			}
		}
	}
	if ctx != nil {
		return func() CodecErr {
			p := &ctx.Err
			ctx.Err = res
			return *p
		}()
	}
	return res
}
func vpx_codec_get_cx_data(ctx *CodecCtx, iter *CodecIter) *CodecCxPkt {
	var pkt *CodecCxPkt = nil
	if ctx != nil {
		if iter == nil {
			ctx.Err = CodecErr(VPX_CODEC_INVALID_PARAM)
		} else if ctx.Iface == nil || ctx.Priv == nil {
			ctx.Err = CodecErr(VPX_CODEC_ERROR)
		} else if (ctx.Iface.Caps & VPX_CODEC_CAP_ENCODER) == 0 {
			ctx.Err = CodecErr(VPX_CODEC_INCAPABLE)
		} else {
			pkt = ctx.Iface.Enc.GetCxData(GetAlgPvt(ctx), iter)
		}
	}
	if pkt != nil && pkt.Kind == CodecCxPktKind(VPX_CODEC_CX_FRAME_PKT) {
		var (
			priv    *CodecPriv = ctx.Priv
			dst_buf *byte      = (*byte)(priv.Enc.Cx_data_dst_buf.Buf)
		)
		if dst_buf != nil && pkt.Data.Raw.Buf != unsafe.Pointer(dst_buf) && pkt.Data.Raw.Sz+uint64(priv.Enc.Cx_data_pad_before)+uint64(priv.Enc.Cx_data_pad_after) <= priv.Enc.Cx_data_dst_buf.Sz {
			var modified_pkt *CodecCxPkt = &priv.Enc.Cx_data_pkt
			libc.MemCpy(unsafe.Add(unsafe.Pointer(dst_buf), priv.Enc.Cx_data_pad_before), pkt.Data.Raw.Buf, int(pkt.Data.Raw.Sz))
			*modified_pkt = *pkt
			modified_pkt.Data.Raw.Buf = unsafe.Pointer(dst_buf)
			modified_pkt.Data.Raw.Sz += uint64(priv.Enc.Cx_data_pad_before + priv.Enc.Cx_data_pad_after)
			pkt = modified_pkt
		}
		if unsafe.Pointer(dst_buf) == pkt.Data.Raw.Buf {
			priv.Enc.Cx_data_dst_buf.Buf = unsafe.Add(unsafe.Pointer(dst_buf), pkt.Data.Raw.Sz)
			priv.Enc.Cx_data_dst_buf.Sz -= pkt.Data.Raw.Sz
		}
	}
	return pkt
}
func vpx_codec_set_cx_data_buf(ctx *CodecCtx, buf *FixedBuf, pad_before uint, pad_after uint) CodecErr {
	if ctx == nil || ctx.Priv == nil {
		return CodecErr(VPX_CODEC_INVALID_PARAM)
	}
	if buf != nil {
		ctx.Priv.Enc.Cx_data_dst_buf = *buf
		ctx.Priv.Enc.Cx_data_pad_before = pad_before
		ctx.Priv.Enc.Cx_data_pad_after = pad_after
	} else {
		ctx.Priv.Enc.Cx_data_dst_buf.Buf = nil
		ctx.Priv.Enc.Cx_data_dst_buf.Sz = 0
		ctx.Priv.Enc.Cx_data_pad_before = 0
		ctx.Priv.Enc.Cx_data_pad_after = 0
	}
	return CodecErr(VPX_CODEC_OK)
}
func vpx_codec_get_preview_frame(ctx *CodecCtx) *Image {
	var img *Image = nil
	if ctx != nil {
		if ctx.Iface == nil || ctx.Priv == nil {
			ctx.Err = CodecErr(VPX_CODEC_ERROR)
		} else if (ctx.Iface.Caps & VPX_CODEC_CAP_ENCODER) == 0 {
			ctx.Err = CodecErr(VPX_CODEC_INCAPABLE)
		} else if ctx.Iface.Enc.GetPreview == nil {
			ctx.Err = CodecErr(VPX_CODEC_INCAPABLE)
		} else {
			img = ctx.Iface.Enc.GetPreview(GetAlgPvt(ctx))
		}
	}
	return img
}
func vpx_codec_get_global_headers(ctx *CodecCtx) *FixedBuf {
	var buf *FixedBuf = nil
	if ctx != nil {
		if ctx.Iface == nil || ctx.Priv == nil {
			ctx.Err = CodecErr(VPX_CODEC_ERROR)
		} else if (ctx.Iface.Caps & VPX_CODEC_CAP_ENCODER) == 0 {
			ctx.Err = CodecErr(VPX_CODEC_INCAPABLE)
		} else if ctx.Iface.Enc.GetGlobHdrs == nil {
			ctx.Err = CodecErr(VPX_CODEC_INCAPABLE)
		} else {
			buf = ctx.Iface.Enc.GetGlobHdrs(GetAlgPvt(ctx))
		}
	}
	return buf
}
func vpx_codec_enc_config_set(ctx *CodecCtx, cfg *CodecEncCfg) CodecErr {
	var res CodecErr
	if ctx == nil || ctx.Iface == nil || ctx.Priv == nil || cfg == nil {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if (ctx.Iface.Caps & VPX_CODEC_CAP_ENCODER) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else {
		res = ctx.Iface.Enc.CfgSet(GetAlgPvt(ctx), cfg)
	}
	if ctx != nil {
		return func() CodecErr {
			p := &ctx.Err
			ctx.Err = res
			return *p
		}()
	}
	return res
}
func vpx_codec_pkt_list_add(list *CodecPktList, pkt *CodecCxPkt) int {
	if list.Cnt < list.Max {
		list.Pkts[func() uint {
			p := &list.Cnt
			x := *p
			*p++
			return x
		}()] = *pkt
		return 0
	}
	return 1
}
