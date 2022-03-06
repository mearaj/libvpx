package vpx

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"unsafe"
)

func CodecVersion() int {
	return (int(VERSION_MAJOR << 16)) | (int(VERSION_MINOR << 8)) | VERSION_PATCH
}
func CodecVersionStr() *byte {
	return libc.CString(VERSION_STRING_NOSP)
}
func CodecVersionExtraStr() *byte {
	return libc.CString(VERSION_EXTRA)
}
func CodecIFaceName(iface *CodecIFace) *byte {
	if iface != nil {
		return iface.Name
	}
	return libc.CString("<invalid interface>")
}
func CodecErrToString(err CodecErr) *byte {
	switch err {
	case VPX_CODEC_OK:
		return libc.CString("Success")
	case VPX_CODEC_ERROR:
		return libc.CString("Unspecified internal error")
	case VPX_CODEC_MEM_ERROR:
		return libc.CString("Memory allocation error")
	case VPX_CODEC_ABI_MISMATCH:
		return libc.CString("ABI version mismatch")
	case VPX_CODEC_INCAPABLE:
		return libc.CString("Codec does not implement requested capability")
	case VPX_CODEC_UNSUP_BITSTREAM:
		return libc.CString("Bitstream not supported by this decoder")
	case VPX_CODEC_UNSUP_FEATURE:
		return libc.CString("Bitstream required feature not supported by this decoder")
	case VPX_CODEC_CORRUPT_FRAME:
		return libc.CString("Corrupt frame detected")
	case VPX_CODEC_INVALID_PARAM:
		return libc.CString("Invalid parameter")
	case VPX_CODEC_LIST_END:
		return libc.CString("End of iterated list")
	}
	return libc.CString("Unrecognized error code")
}
func CodecError(ctx *CodecCtx) *byte {
	if ctx != nil {
		return CodecErrToString(ctx.Err)
	}
	return CodecErrToString(CodecErr(VPX_CODEC_INVALID_PARAM))
}
func CodecErrorDetail(ctx *CodecCtx) *byte {
	if ctx != nil && ctx.Err != 0 {
		if ctx.Priv != nil {
			return ctx.Priv.Err_detail
		}
		return ctx.Err_detail
	}
	return nil
}
func CodecDestroy(ctx *CodecCtx) CodecErr {
	var res CodecErr
	if ctx == nil {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if ctx.Iface == nil || ctx.Priv == nil {
		res = CodecErr(VPX_CODEC_ERROR)
	} else {
		ctx.Iface.Destroy((*CodecAlgPvt)(unsafe.Pointer(ctx.Priv)))
		ctx.Iface = nil
		ctx.Name = nil
		ctx.Priv = nil
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
func CodecGetCaps(iface *CodecIFace) CodecCaps {
	if iface != nil {
		return iface.Caps
	}
	return 0
}
func CodecControl(ctx *CodecCtx, CtrlId int, _rest ...interface{}) CodecErr {
	var res CodecErr
	if ctx == nil || CtrlId == 0 {
		res = CodecErr(VPX_CODEC_INVALID_PARAM)
	} else if ctx.Iface == nil || ctx.Priv == nil || ctx.Iface.Ctrl_maps == nil {
		res = CodecErr(VPX_CODEC_ERROR)
	} else {
		var entry *FnMap
		res = CodecErr(VPX_CODEC_INCAPABLE)
		for entry = ctx.Iface.Ctrl_maps; entry.Fn != nil; entry = (*FnMap)(unsafe.Add(unsafe.Pointer(entry), unsafe.Sizeof(FnMap{})*1)) {
			if entry.Ctrl_id == 0 || entry.Ctrl_id == CtrlId {
				var ap libc.ArgList
				ap.Start(CtrlId, _rest)
				res = entry.Fn((*CodecAlgPvt)(unsafe.Pointer(ctx.Priv)), ap)
				ap.End()
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
func InternalError(info *InternalErrorInfo, error CodecErr, fmt *byte, _rest ...interface{}) {
	var ap libc.ArgList
	info.Error_code = error
	info.Has_detail = 0
	if fmt != nil {
		var sz uint64 = uint64(unsafe.Sizeof([80]byte{}))
		info.Has_detail = 1
		ap.Start(fmt, _rest)
		stdio.Vsnprintf(&info.Detail[0], int(sz-1), libc.GoString(fmt), ap)
		ap.End()
		info.Detail[sz-1] = byte('\x00')
	}
	if info.Setjmp != 0 {
		info.Jmp.LongJump(int(info.Error_code))
	}
}
