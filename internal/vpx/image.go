package vpx

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/mem"
	"math"
	"unsafe"
)

func ImgAllocHelper(img *Image, fmt ImgFmt, d_w uint, d_h uint, buf_align uint, stride_align uint, img_data *uint8) *Image {
	var (
		h               uint
		w               uint
		s               uint
		xcs             uint
		ycs             uint
		bps             uint
		stride_in_bytes uint
		align           uint
	)
	if img != nil {
		*img = Image{}
	}
	if buf_align == 0 {
		buf_align = 1
	}
	if buf_align&(buf_align-1) != 0 {
		goto fail
	}
	if stride_align == 0 {
		stride_align = 1
	}
	if stride_align&(stride_align-1) != 0 {
		goto fail
	}
	switch fmt {
	case VPX_IMG_FMT_I420:
		fallthrough
	case VPX_IMG_FMT_YV12:
		fallthrough
	case VPX_IMG_FMT_NV12:
		bps = 12
	case VPX_IMG_FMT_I422:
		fallthrough
	case VPX_IMG_FMT_I440:
		bps = 16
	case VPX_IMG_FMT_I444:
		bps = 24
	case VPX_IMG_FMT_I42016:
		bps = 24
	case VPX_IMG_FMT_I42216:
		fallthrough
	case VPX_IMG_FMT_I44016:
		bps = 32
	case VPX_IMG_FMT_I44416:
		bps = 48
	default:
		bps = 16
	}
	switch fmt {
	case VPX_IMG_FMT_I420:
		fallthrough
	case VPX_IMG_FMT_YV12:
		fallthrough
	case VPX_IMG_FMT_I422:
		fallthrough
	case VPX_IMG_FMT_I42016:
		fallthrough
	case VPX_IMG_FMT_I42216:
		xcs = 1
	default:
		xcs = 0
	}
	switch fmt {
	case VPX_IMG_FMT_I420:
		fallthrough
	case VPX_IMG_FMT_NV12:
		fallthrough
	case VPX_IMG_FMT_I440:
		fallthrough
	case VPX_IMG_FMT_YV12:
		fallthrough
	case VPX_IMG_FMT_I42016:
		fallthrough
	case VPX_IMG_FMT_I44016:
		ycs = 1
	default:
		ycs = 0
	}
	w = d_w
	h = d_h
	if (fmt & VPX_IMG_FMT_PLANAR) != 0 {
		s = w
	} else {
		s = bps * w / 8
	}
	s = (s + stride_align - 1) & ^(stride_align - 1)
	if (fmt & VPX_IMG_FMT_HIGHBITDEPTH) != 0 {
		stride_in_bytes = s * 2
	} else {
		stride_in_bytes = s
	}
	if img == nil {
		img = &make([]Image, 1)[0]
		if img == nil {
			goto fail
		}
		img.Self_allocd = 1
	}
	img.Img_data = img_data
	if img_data == nil {
		var alloc_size uint64
		align = (1 << xcs) - 1
		w = (d_w + align) & ^align
		align = (1 << ycs) - 1
		h = (d_h + align) & ^align
		if (fmt & VPX_IMG_FMT_PLANAR) != 0 {
			s = w
		} else {
			s = bps * w / 8
		}
		s = (s + stride_align - 1) & ^(stride_align - 1)
		if (fmt & VPX_IMG_FMT_HIGHBITDEPTH) != 0 {
			stride_in_bytes = s * 2
		} else {
			stride_in_bytes = s
		}
		if (fmt & VPX_IMG_FMT_PLANAR) != 0 {
			alloc_size = uint64(h) * uint64(s) * uint64(bps) / 8
		} else {
			alloc_size = uint64(h) * uint64(s)
		}
		if alloc_size != alloc_size {
			goto fail
		}
		img.Img_data = (*uint8)(mem.VpxMemAlign(uint64(buf_align), alloc_size))
		img.Img_data_owner = 1
	}
	if img.Img_data == nil {
		goto fail
	}
	img.Fmt = fmt
	if (fmt & VPX_IMG_FMT_HIGHBITDEPTH) != 0 {
		img.Bit_depth = 16
	} else {
		img.Bit_depth = 8
	}
	img.W = w
	img.H = h
	img.X_chroma_shift = xcs
	img.Y_chroma_shift = ycs
	img.Bps = int(bps)
	img.Stride[VPX_PLANE_Y] = func() int {
		p := &img.Stride[VPX_PLANE_ALPHA]
		img.Stride[VPX_PLANE_ALPHA] = int(stride_in_bytes)
		return *p
	}()
	img.Stride[VPX_PLANE_U] = func() int {
		p := &img.Stride[VPX_PLANE_V]
		img.Stride[VPX_PLANE_V] = int(stride_in_bytes >> xcs)
		return *p
	}()
	if ImgSetRect(img, 0, 0, d_w, d_h) == 0 {
		return img
	}
fail:
	ImgFree(img)
	return nil
}
func ImgAlloc(img *Image, fmt ImgFmt, d_w uint, d_h uint, align uint) *Image {
	return ImgAllocHelper(img, fmt, d_w, d_h, align, align, nil)
}
func ImgWrap(img *Image, fmt ImgFmt, d_w uint, d_h uint, stride_align uint, img_data *uint8) *Image {
	return ImgAllocHelper(img, fmt, d_w, d_h, 1, stride_align, img_data)
}
func ImgSetRect(img *Image, x uint, y uint, w uint, h uint) int {
	if x <= math.MaxUint64-w && x+w <= img.W && y <= math.MaxUint64-h && y+h <= img.H {
		img.D_w = w
		img.D_h = h
		if (img.Fmt & VPX_IMG_FMT_PLANAR) == 0 {
			img.Planes[VPX_PLANE_PACKED] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img.Img_data), x*uint(img.Bps)/8))), y*uint(img.Stride[VPX_PLANE_PACKED])))
		} else {
			var bytes_per_sample int
			if (img.Fmt & VPX_IMG_FMT_HIGHBITDEPTH) != 0 {
				bytes_per_sample = 2
			} else {
				bytes_per_sample = 1
			}
			var data *uint8 = img.Img_data
			if img.Fmt&VPX_IMG_FMT_HAS_ALPHA != 0 {
				img.Planes[VPX_PLANE_ALPHA] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), x*uint(bytes_per_sample)))), y*uint(img.Stride[VPX_PLANE_ALPHA])))
				data = (*uint8)(unsafe.Add(unsafe.Pointer(data), img.H*uint(img.Stride[VPX_PLANE_ALPHA])))
			}
			img.Planes[VPX_PLANE_Y] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), x*uint(bytes_per_sample)))), y*uint(img.Stride[VPX_PLANE_Y])))
			data = (*uint8)(unsafe.Add(unsafe.Pointer(data), img.H*uint(img.Stride[VPX_PLANE_Y])))
			if img.Fmt == ImgFmt(VPX_IMG_FMT_NV12) {
				img.Planes[VPX_PLANE_U] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), x>>img.X_chroma_shift))), (y>>img.Y_chroma_shift)*uint(img.Stride[VPX_PLANE_U])))
				img.Planes[VPX_PLANE_V] = (*uint8)(unsafe.Add(unsafe.Pointer(img.Planes[VPX_PLANE_U]), 1))
			} else if (img.Fmt & VPX_IMG_FMT_UV_FLIP) == 0 {
				img.Planes[VPX_PLANE_U] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), (x>>img.X_chroma_shift)*uint(bytes_per_sample)))), (y>>img.Y_chroma_shift)*uint(img.Stride[VPX_PLANE_U])))
				data = (*uint8)(unsafe.Add(unsafe.Pointer(data), (img.H>>img.Y_chroma_shift)*uint(img.Stride[VPX_PLANE_U])))
				img.Planes[VPX_PLANE_V] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), (x>>img.X_chroma_shift)*uint(bytes_per_sample)))), (y>>img.Y_chroma_shift)*uint(img.Stride[VPX_PLANE_V])))
			} else {
				img.Planes[VPX_PLANE_V] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), (x>>img.X_chroma_shift)*uint(bytes_per_sample)))), (y>>img.Y_chroma_shift)*uint(img.Stride[VPX_PLANE_V])))
				data = (*uint8)(unsafe.Add(unsafe.Pointer(data), (img.H>>img.Y_chroma_shift)*uint(img.Stride[VPX_PLANE_V])))
				img.Planes[VPX_PLANE_U] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(data), (x>>img.X_chroma_shift)*uint(bytes_per_sample)))), (y>>img.Y_chroma_shift)*uint(img.Stride[VPX_PLANE_U])))
			}
		}
		return 0
	}
	return -1
}
func ImgFlip(img *Image) {
	img.Planes[VPX_PLANE_Y] = (*uint8)(unsafe.Add(unsafe.Pointer(img.Planes[VPX_PLANE_Y]), int(img.D_h-1)*img.Stride[VPX_PLANE_Y]))
	img.Stride[VPX_PLANE_Y] = -img.Stride[VPX_PLANE_Y]
	img.Planes[VPX_PLANE_U] = (*uint8)(unsafe.Add(unsafe.Pointer(img.Planes[VPX_PLANE_U]), int((img.D_h>>img.Y_chroma_shift)-1)*img.Stride[VPX_PLANE_U]))
	img.Stride[VPX_PLANE_U] = -img.Stride[VPX_PLANE_U]
	img.Planes[VPX_PLANE_V] = (*uint8)(unsafe.Add(unsafe.Pointer(img.Planes[VPX_PLANE_V]), int((img.D_h>>img.Y_chroma_shift)-1)*img.Stride[VPX_PLANE_V]))
	img.Stride[VPX_PLANE_V] = -img.Stride[VPX_PLANE_V]
	img.Planes[VPX_PLANE_ALPHA] = (*uint8)(unsafe.Add(unsafe.Pointer(img.Planes[VPX_PLANE_ALPHA]), int(img.D_h-1)*img.Stride[VPX_PLANE_ALPHA]))
	img.Stride[VPX_PLANE_ALPHA] = -img.Stride[VPX_PLANE_ALPHA]
}
func ImgFree(img *Image) {
	if img != nil {
		if img.Img_data != nil && img.Img_data_owner != 0 {
			mem.VpxFree(unsafe.Pointer(img.Img_data))
		}
		if img.Self_allocd != 0 {
			libc.Free(unsafe.Pointer(img))
		}
	}
}
