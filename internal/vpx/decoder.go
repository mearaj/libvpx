package vpx

import "unsafe"

func CodecDecInitVer(ctx *CodecCtx, iface *CodecIFace, cfg *CodecDecCfg, flags CodecFlags, ver int) CodecErr {
	var res CodecErr
	if ver != (3 + (4 + 5)) {
		res = CodecErr(VPX_CODEC_ABI_MISMATCH)
	} else if ctx == nil || iface == nil {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if iface.Abi_version != 5 {
		res = CodecErr(VPX_CODEC_ABI_MISMATCH)
	} else if (flags&VPX_CODEC_USE_POSTPROC) != 0 && (iface.Caps&VPX_CODEC_CAP_POSTPROC) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else if (flags&VPX_CODEC_USE_ERROR_CONCEALMENT) != 0 && (iface.Caps&VPX_CODEC_CAP_ERROR_CONCEALMENT) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else if (flags&VPX_CODEC_USE_INPUT_FRAGMENTS) != 0 && (iface.Caps&VPX_CODEC_CAP_INPUT_FRAGMENTS) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else if (iface.Caps & VPX_CODEC_CAP_DECODER) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else {
		*ctx = CodecCtx{}
		ctx.Iface = iface
		ctx.Name = iface.Name
		ctx.Priv = nil
		ctx.Init_flags = flags
		ctx.Config.Dec = cfg
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
func CodecPeekStreamInfo(iface *CodecIFace, data *uint8, data_sz uint, si *CodecStreamInfo) CodecErr {
	var res CodecErr
	if iface == nil || data == nil || data_sz == 0 || si == nil || si.Sz < uint(unsafe.Sizeof(CodecStreamInfo{})) {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else {
		si.W = 0
		si.H = 0
		res = iface.Dec.PeekSi(data, data_sz, si)
	}
	return res
}
func CodecGetStreamInfo(ctx *CodecCtx, si *CodecStreamInfo) CodecErr {
	var res CodecErr
	if ctx == nil || si == nil || si.Sz < uint(unsafe.Sizeof(CodecStreamInfo{})) {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if ctx.Iface == nil || ctx.Priv == nil {
		res = CodecErr(VPX_CODEC_ERROR)
	} else {
		si.W = 0
		si.H = 0
		res = ctx.Iface.Dec.GetSi(GetAlgPvt(ctx), si)
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
func CodecDecode(ctx *CodecCtx, data *uint8, data_sz uint, user_priv unsafe.Pointer, deadline int) CodecErr {
	var res CodecErr
	if ctx == nil || data == nil && data_sz != 0 || data != nil && data_sz == 0 {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if ctx.Iface == nil || ctx.Priv == nil {
		res = CodecErr(VPX_CODEC_ERROR)
	} else {
		res = ctx.Iface.Dec.Decode(GetAlgPvt(ctx), data, data_sz, user_priv, deadline)
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
func CodecGetFrame(ctx *CodecCtx, iter *CodecIter) *Image {
	var img *Image
	if ctx == nil || iter == nil || ctx.Iface == nil || ctx.Priv == nil {
		img = nil
	} else {
		img = ctx.Iface.Dec.GetFrame(GetAlgPvt(ctx), iter)
	}
	return img
}
func CodecRegisterPutFrameCb(ctx *CodecCtx, cb vpx_codec_put_frame_cb_fn_t, user_priv unsafe.Pointer) CodecErr {
	var res CodecErr
	if ctx == nil || cb == nil {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if ctx.Iface == nil || ctx.Priv == nil {
		res = CodecErr(VPX_CODEC_ERROR)
	} else if (ctx.Iface.Caps & VPX_CODEC_CAP_PUT_FRAME) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else {
		ctx.Priv.Dec.Put_frame_cb.U.Put_frame = cb
		ctx.Priv.Dec.Put_frame_cb.User_priv = user_priv
		res = CodecErr(VPX_CODEC_OK)
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
func CodecRegisterPutSliceCb(ctx *CodecCtx, cb vpx_codec_put_slice_cb_fn_t, user_priv unsafe.Pointer) CodecErr {
	var res CodecErr
	if ctx == nil || cb == nil {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if ctx.Iface == nil || ctx.Priv == nil {
		res = CodecErr(VPX_CODEC_ERROR)
	} else if (ctx.Iface.Caps & VPX_CODEC_CAP_PUT_SLICE) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else {
		ctx.Priv.Dec.Put_slice_cb.U.Put_slice = cb
		ctx.Priv.Dec.Put_slice_cb.User_priv = user_priv
		res = CodecErr(VPX_CODEC_OK)
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
func CodecSetFrameBufferFunctions(ctx *CodecCtx, cb_get GetFrameBufferCbFn, cb_release ReleaseFrameBufferCbFn, cb_priv unsafe.Pointer) CodecErr {
	var res CodecErr
	if ctx == nil || cb_get == nil || cb_release == nil {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if ctx.Iface == nil || ctx.Priv == nil {
		res = CodecErr(VPX_CODEC_ERROR)
	} else if (ctx.Iface.Caps & VPX_CODEC_CAP_EXTERNAL_FRAME_BUFFER) == 0 {
		res = CodecErr(VPX_CODEC_INCAPABLE)
	} else {
		res = ctx.Iface.Dec.SetFbFn(GetAlgPvt(ctx), cb_get, cb_release, cb_priv)
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
