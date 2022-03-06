package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const SSIM_C1 = 0
const SSIM_C2 = 0
const FS_NLEVELS = 4

type fs_level struct {
	Im1  *uint32
	Im2  *uint32
	Ssim *float64
	W    int
	H    int
}
type fs_ctx struct {
	Level   *fs_level
	Nlevels int
	Col_buf *uint
}

func fs_ctx_init(_ctx *fs_ctx, _w int, _h int, _nlevels int) {
	var (
		data      *uint8
		data_size uint64
		lw        int
		lh        int
		l         int
	)
	lw = (_w + 1) >> 1
	lh = (_h + 1) >> 1
	data_size = uint64(_nlevels*int(unsafe.Sizeof(fs_level{})) + (lw+8)*2*8*int(unsafe.Sizeof(uint(0))))
	for l = 0; l < _nlevels; l++ {
		var (
			im_size    uint64
			level_size uint64
		)
		im_size = uint64(lw * int(uint64(lh)))
		level_size = im_size * 2 * uint64(unsafe.Sizeof(uint32(0)))
		level_size += uint64(unsafe.Sizeof(float64(0)) - 1)
		level_size /= uint64(unsafe.Sizeof(float64(0)))
		level_size += im_size
		level_size *= uint64(unsafe.Sizeof(float64(0)))
		data_size += level_size
		lw = (lw + 1) >> 1
		lh = (lh + 1) >> 1
	}
	data = (*uint8)(libc.Malloc(int(data_size)))
	_ctx.Level = (*fs_level)(unsafe.Pointer(data))
	_ctx.Nlevels = _nlevels
	data = (*uint8)(unsafe.Add(unsafe.Pointer(data), _nlevels*int(unsafe.Sizeof(fs_level{}))))
	lw = (_w + 1) >> 1
	lh = (_h + 1) >> 1
	for l = 0; l < _nlevels; l++ {
		var (
			im_size    uint64
			level_size uint64
		)
		(*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(l)))).W = lw
		(*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(l)))).H = lh
		im_size = uint64(lw * int(uint64(lh)))
		level_size = im_size * 2 * uint64(unsafe.Sizeof(uint32(0)))
		level_size += uint64(unsafe.Sizeof(float64(0)) - 1)
		level_size /= uint64(unsafe.Sizeof(float64(0)))
		level_size *= uint64(unsafe.Sizeof(float64(0)))
		(*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(l)))).Im1 = (*uint32)(unsafe.Pointer(data))
		(*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(l)))).Im2 = (*uint32)(unsafe.Add(unsafe.Pointer((*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(l)))).Im1), unsafe.Sizeof(uint32(0))*uintptr(im_size)))
		data = (*uint8)(unsafe.Add(unsafe.Pointer(data), level_size))
		(*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(l)))).Ssim = (*float64)(unsafe.Pointer(data))
		data = (*uint8)(unsafe.Add(unsafe.Pointer(data), im_size*uint64(unsafe.Sizeof(float64(0)))))
		lw = (lw + 1) >> 1
		lh = (lh + 1) >> 1
	}
	_ctx.Col_buf = (*uint)(unsafe.Pointer(data))
}
func fs_ctx_clear(_ctx *fs_ctx) {
	libc.Free(unsafe.Pointer(_ctx.Level))
}
func fs_downsample_level(_ctx *fs_ctx, _l int) {
	var (
		src1 *uint32
		src2 *uint32
		dst1 *uint32
		dst2 *uint32
		w2   int
		h2   int
		w    int
		h    int
		i    int
		j    int
	)
	w = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).W
	h = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).H
	dst1 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).Im1
	dst2 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).Im2
	w2 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l-1)))).W
	h2 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l-1)))).H
	src1 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l-1)))).Im1
	src2 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l-1)))).Im2
	for j = 0; j < h; j++ {
		var (
			j0offs int
			j1offs int
		)
		j0offs = j * 2 * w2
		j1offs = (func() int {
			if (j*2 + 1) < h2 {
				return j*2 + 1
			}
			return h2
		}()) * w2
		for i = 0; i < w; i++ {
			var (
				i0 int
				i1 int
			)
			i0 = i * 2
			if (i0 + 1) < w2 {
				i1 = i0 + 1
			} else {
				i1 = w2
			}
			*(*uint32)(unsafe.Add(unsafe.Pointer(dst1), unsafe.Sizeof(uint32(0))*uintptr(j*w+i))) = uint32(int64(*(*uint32)(unsafe.Add(unsafe.Pointer(src1), unsafe.Sizeof(uint32(0))*uintptr(j0offs+i0)))) + int64(*(*uint32)(unsafe.Add(unsafe.Pointer(src1), unsafe.Sizeof(uint32(0))*uintptr(j0offs+i1)))) + int64(*(*uint32)(unsafe.Add(unsafe.Pointer(src1), unsafe.Sizeof(uint32(0))*uintptr(j1offs+i0)))) + int64(*(*uint32)(unsafe.Add(unsafe.Pointer(src1), unsafe.Sizeof(uint32(0))*uintptr(j1offs+i1)))))
			*(*uint32)(unsafe.Add(unsafe.Pointer(dst2), unsafe.Sizeof(uint32(0))*uintptr(j*w+i))) = uint32(int64(*(*uint32)(unsafe.Add(unsafe.Pointer(src2), unsafe.Sizeof(uint32(0))*uintptr(j0offs+i0)))) + int64(*(*uint32)(unsafe.Add(unsafe.Pointer(src2), unsafe.Sizeof(uint32(0))*uintptr(j0offs+i1)))) + int64(*(*uint32)(unsafe.Add(unsafe.Pointer(src2), unsafe.Sizeof(uint32(0))*uintptr(j1offs+i0)))) + int64(*(*uint32)(unsafe.Add(unsafe.Pointer(src2), unsafe.Sizeof(uint32(0))*uintptr(j1offs+i1)))))
		}
	}
}
func fs_downsample_level0(_ctx *fs_ctx, _src1 *uint8, _s1ystride int, _src2 *uint8, _s2ystride int, _w int, _h int, bd uint32, shift uint32) {
	var (
		dst1 *uint32
		dst2 *uint32
		w    int
		h    int
		i    int
		j    int
	)
	w = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*0))).W
	h = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*0))).H
	dst1 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*0))).Im1
	dst2 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*0))).Im2
	for j = 0; j < h; j++ {
		var (
			j0 int
			j1 int
		)
		j0 = j * 2
		if (j0 + 1) < _h {
			j1 = j0 + 1
		} else {
			j1 = _h
		}
		for i = 0; i < w; i++ {
			var (
				i0 int
				i1 int
			)
			i0 = i * 2
			if (i0 + 1) < _w {
				i1 = i0 + 1
			} else {
				i1 = _w
			}
			if bd == 8 && shift == 0 {
				*(*uint32)(unsafe.Add(unsafe.Pointer(dst1), unsafe.Sizeof(uint32(0))*uintptr(j*w+i))) = uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(_src1), j0*_s1ystride+i0)) + *(*uint8)(unsafe.Add(unsafe.Pointer(_src1), j0*_s1ystride+i1)) + *(*uint8)(unsafe.Add(unsafe.Pointer(_src1), j1*_s1ystride+i0)) + *(*uint8)(unsafe.Add(unsafe.Pointer(_src1), j1*_s1ystride+i1)))
				*(*uint32)(unsafe.Add(unsafe.Pointer(dst2), unsafe.Sizeof(uint32(0))*uintptr(j*w+i))) = uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(_src2), j0*_s2ystride+i0)) + *(*uint8)(unsafe.Add(unsafe.Pointer(_src2), j0*_s2ystride+i1)) + *(*uint8)(unsafe.Add(unsafe.Pointer(_src2), j1*_s2ystride+i0)) + *(*uint8)(unsafe.Add(unsafe.Pointer(_src2), j1*_s2ystride+i1)))
			} else {
				var (
					src1s *uint16 = ((*uint16)(unsafe.Pointer(uintptr((uint64(uintptr(unsafe.Pointer(_src1)))) << 1))))
					src2s *uint16 = ((*uint16)(unsafe.Pointer(uintptr((uint64(uintptr(unsafe.Pointer(_src2)))) << 1))))
				)
				*(*uint32)(unsafe.Add(unsafe.Pointer(dst1), unsafe.Sizeof(uint32(0))*uintptr(j*w+i))) = (uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(src1s), unsafe.Sizeof(uint16(0))*uintptr(j0*_s1ystride+i0)))) >> shift) + (uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(src1s), unsafe.Sizeof(uint16(0))*uintptr(j0*_s1ystride+i1)))) >> shift) + (uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(src1s), unsafe.Sizeof(uint16(0))*uintptr(j1*_s1ystride+i0)))) >> shift) + (uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(src1s), unsafe.Sizeof(uint16(0))*uintptr(j1*_s1ystride+i1)))) >> shift)
				*(*uint32)(unsafe.Add(unsafe.Pointer(dst2), unsafe.Sizeof(uint32(0))*uintptr(j*w+i))) = (uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(src2s), unsafe.Sizeof(uint16(0))*uintptr(j0*_s2ystride+i0)))) >> shift) + (uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(src2s), unsafe.Sizeof(uint16(0))*uintptr(j0*_s2ystride+i1)))) >> shift) + (uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(src2s), unsafe.Sizeof(uint16(0))*uintptr(j1*_s2ystride+i0)))) >> shift) + (uint32(*(*uint16)(unsafe.Add(unsafe.Pointer(src2s), unsafe.Sizeof(uint16(0))*uintptr(j1*_s2ystride+i1)))) >> shift)
			}
		}
	}
}
func fs_apply_luminance(_ctx *fs_ctx, _l int, bit_depth int) {
	var (
		col_sums_x *uint
		col_sums_y *uint
		im1        *uint32
		im2        *uint32
		ssim       *float64
		c1         float64
		w          int
		h          int
		j0offs     int
		j1offs     int
		i          int
		j          int
		ssim_c1    float64 = (math.MaxUint8 * math.MaxUint8 * 0.01 * 0.01)
	)
	if bit_depth == 8 {
	} else {
		__assert_fail(libc.CString("bit_depth == 8"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	_ = bit_depth
	w = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).W
	h = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).H
	col_sums_x = _ctx.Col_buf
	col_sums_y = (*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(w)))
	im1 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).Im1
	im2 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).Im2
	for i = 0; i < w; i++ {
		*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i))) = uint(*(*uint32)(unsafe.Add(unsafe.Pointer(im1), unsafe.Sizeof(uint32(0))*uintptr(i))) * 5)
	}
	for i = 0; i < w; i++ {
		*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_y), unsafe.Sizeof(uint(0))*uintptr(i))) = uint(*(*uint32)(unsafe.Add(unsafe.Pointer(im2), unsafe.Sizeof(uint32(0))*uintptr(i))) * 5)
	}
	for j = 1; j < 4; j++ {
		j1offs = (func() int {
			if j < (h - 1) {
				return j
			}
			return h - 1
		}()) * w
		for i = 0; i < w; i++ {
			*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i))) += uint(*(*uint32)(unsafe.Add(unsafe.Pointer(im1), unsafe.Sizeof(uint32(0))*uintptr(j1offs+i))))
		}
		for i = 0; i < w; i++ {
			*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_y), unsafe.Sizeof(uint(0))*uintptr(i))) += uint(*(*uint32)(unsafe.Add(unsafe.Pointer(im2), unsafe.Sizeof(uint32(0))*uintptr(j1offs+i))))
		}
	}
	ssim = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).Ssim
	c1 = ssim_c1 * 4096 * float64(1<<(_l*4))
	for j = 0; j < h; j++ {
		var (
			mux int64
			muy int64
			i0  int
			i1  int
		)
		mux = int64(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*0)) * 5)
		muy = int64(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_y), unsafe.Sizeof(uint(0))*0)) * 5)
		for i = 1; i < 4; i++ {
			if i < (w - 1) {
				i1 = i
			} else {
				i1 = w - 1
			}
			mux += int64(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i1))))
			muy += int64(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_y), unsafe.Sizeof(uint(0))*uintptr(i1))))
		}
		for i = 0; i < w; i++ {
			*(*float64)(unsafe.Add(unsafe.Pointer(ssim), unsafe.Sizeof(float64(0))*uintptr(j*w+i))) *= (float64(mux*2)*float64(muy) + c1) / (float64(mux)*float64(mux) + float64(muy)*float64(muy) + c1)
			if i+1 < w {
				if 0 > (i - 4) {
					i0 = 0
				} else {
					i0 = i - 4
				}
				if (i + 4) < (w - 1) {
					i1 = i + 4
				} else {
					i1 = w - 1
				}
				mux += int64(int(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i1)))) - int(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i0)))))
				muy += int64(int(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i1)))) - int(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i0)))))
			}
		}
		if j+1 < h {
			j0offs = (func() int {
				if 0 > (j - 4) {
					return 0
				}
				return j - 4
			}()) * w
			for i = 0; i < w; i++ {
				*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i))) -= uint(*(*uint32)(unsafe.Add(unsafe.Pointer(im1), unsafe.Sizeof(uint32(0))*uintptr(j0offs+i))))
			}
			for i = 0; i < w; i++ {
				*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_y), unsafe.Sizeof(uint(0))*uintptr(i))) -= uint(*(*uint32)(unsafe.Add(unsafe.Pointer(im2), unsafe.Sizeof(uint32(0))*uintptr(j0offs+i))))
			}
			j1offs = (func() int {
				if (j + 4) < (h - 1) {
					return j + 4
				}
				return h - 1
			}()) * w
			for i = 0; i < w; i++ {
				*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i))) = uint(uint32(int64(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_x), unsafe.Sizeof(uint(0))*uintptr(i)))) + int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im1), unsafe.Sizeof(uint32(0))*uintptr(j1offs+i))))))
			}
			for i = 0; i < w; i++ {
				*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_y), unsafe.Sizeof(uint(0))*uintptr(i))) = uint(uint32(int64(*(*uint)(unsafe.Add(unsafe.Pointer(col_sums_y), unsafe.Sizeof(uint(0))*uintptr(i)))) + int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im2), unsafe.Sizeof(uint32(0))*uintptr(j1offs+i))))))
			}
		}
	}
}
func fs_calc_structure(_ctx *fs_ctx, _l int, bit_depth int) {
	var (
		im1           *uint32
		im2           *uint32
		gx_buf        *uint
		gy_buf        *uint
		ssim          *float64
		col_sums_gx2  [8]float64
		col_sums_gy2  [8]float64
		col_sums_gxgy [8]float64
		c2            float64
		stride        int
		w             int
		h             int
		i             int
		j             int
		ssim_c2       float64 = (math.MaxUint8 * math.MaxUint8 * 0.03 * 0.03)
	)
	if bit_depth == 8 {
	} else {
		__assert_fail(libc.CString("bit_depth == 8"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	_ = bit_depth
	w = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).W
	h = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).H
	im1 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).Im1
	im2 = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).Im2
	ssim = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).Ssim
	gx_buf = _ctx.Col_buf
	stride = w + 8
	gy_buf = (*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(stride*8)))
	libc.MemSet(unsafe.Pointer(gx_buf), 0, stride*(2*8)*int(unsafe.Sizeof(uint(0))))
	c2 = ssim_c2 * float64(1<<(_l*4)) * 16 * 104
	for j = 0; j < h+4; j++ {
		if j < h-1 {
			for i = 0; i < w-1; i++ {
				var (
					g1 int64
					g2 int64
					gx int64
					gy int64
				)
				g1 = int64(labs(int(int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im1), unsafe.Sizeof(uint32(0))*uintptr((j+1)*w+i+1)))) - int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im1), unsafe.Sizeof(uint32(0))*uintptr(j*w+i)))))))
				g2 = int64(labs(int(int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im1), unsafe.Sizeof(uint32(0))*uintptr((j+1)*w+i)))) - int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im1), unsafe.Sizeof(uint32(0))*uintptr(j*w+i+1)))))))
				gx = (func() int64 {
					if g1 > g2 {
						return g1
					}
					return g2
				}())*4 + (func() int64 {
					if g1 < g2 {
						return g1
					}
					return g2
				}())
				g1 = int64(labs(int(int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im2), unsafe.Sizeof(uint32(0))*uintptr((j+1)*w+i+1)))) - int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im2), unsafe.Sizeof(uint32(0))*uintptr(j*w+i)))))))
				g2 = int64(labs(int(int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im2), unsafe.Sizeof(uint32(0))*uintptr((j+1)*w+i)))) - int64(*(*uint32)(unsafe.Add(unsafe.Pointer(im2), unsafe.Sizeof(uint32(0))*uintptr(j*w+i+1)))))))
				gy = (func() int64 {
					if g1 > g2 {
						return g1
					}
					return g2
				}())*4 + (func() int64 {
					if g1 < g2 {
						return g1
					}
					return g2
				}())
				*(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr((j&7)*stride+i+4))) = uint(uint32(gx))
				*(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr((j&7)*stride+i+4))) = uint(uint32(gy))
			}
		} else {
			libc.MemSet(unsafe.Pointer((*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr((j&7)*stride)))), 0, stride*int(unsafe.Sizeof(uint(0))))
			libc.MemSet(unsafe.Pointer((*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr((j&7)*stride)))), 0, stride*int(unsafe.Sizeof(uint(0))))
		}
		if j >= 4 {
			var k int
			col_sums_gx2[3] = func() float64 {
				p := &col_sums_gx2[2]
				col_sums_gx2[2] = func() float64 {
					p := &col_sums_gx2[1]
					col_sums_gx2[1] = func() float64 {
						p := &col_sums_gx2[0]
						col_sums_gx2[0] = 0
						return *p
					}()
					return *p
				}()
				return *p
			}()
			col_sums_gy2[3] = func() float64 {
				p := &col_sums_gy2[2]
				col_sums_gy2[2] = func() float64 {
					p := &col_sums_gy2[1]
					col_sums_gy2[1] = func() float64 {
						p := &col_sums_gy2[0]
						col_sums_gy2[0] = 0
						return *p
					}()
					return *p
				}()
				return *p
			}()
			col_sums_gxgy[3] = func() float64 {
				p := &col_sums_gxgy[2]
				col_sums_gxgy[2] = func() float64 {
					p := &col_sums_gxgy[1]
					col_sums_gxgy[1] = func() float64 {
						p := &col_sums_gxgy[0]
						col_sums_gxgy[0] = 0
						return *p
					}()
					return *p
				}()
				return *p
			}()
			for i = 4; i < 8; i++ {
				for {
					{
						var (
							gx uint
							gy uint
						)
						gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-1))&7)*stride+i+0)))
						gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-1))&7)*stride+i+0)))
						col_sums_gx2[i] = float64(gx) * float64(gx)
						col_sums_gy2[i] = float64(gy) * float64(gy)
						col_sums_gxgy[i] = float64(gx) * float64(gy)
					}
					if true {
						break
					}
				}
				for {
					{
						var (
							gx uint
							gy uint
						)
						gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+0)&7)*stride+i+0)))
						gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+0)&7)*stride+i+0)))
						col_sums_gx2[i] += float64(gx) * float64(gx)
						col_sums_gy2[i] += float64(gy) * float64(gy)
						col_sums_gxgy[i] += float64(gx) * float64(gy)
					}
					if true {
						break
					}
				}
				for k = 1; k < 8-i; k++ {
					for {
						col_sums_gx2[i] = col_sums_gx2[i] * 2
						col_sums_gy2[i] = col_sums_gy2[i] * 2
						col_sums_gxgy[i] = col_sums_gxgy[i] * 2
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+(-k-1))&7)*stride+i+0)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+(-k-1))&7)*stride+i+0)))
							col_sums_gx2[i] += float64(gx) * float64(gx)
							col_sums_gy2[i] += float64(gy) * float64(gy)
							col_sums_gxgy[i] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+k)&7)*stride+i+0)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+k)&7)*stride+i+0)))
							col_sums_gx2[i] += float64(gx) * float64(gx)
							col_sums_gy2[i] += float64(gy) * float64(gy)
							col_sums_gxgy[i] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
				}
			}
			for i = 0; i < w; i++ {
				var (
					mugx2  float64
					mugy2  float64
					mugxgy float64
				)
				mugx2 = col_sums_gx2[0]
				for k = 1; k < 8; k++ {
					mugx2 += col_sums_gx2[k]
				}
				mugy2 = col_sums_gy2[0]
				for k = 1; k < 8; k++ {
					mugy2 += col_sums_gy2[k]
				}
				mugxgy = col_sums_gxgy[0]
				for k = 1; k < 8; k++ {
					mugxgy += col_sums_gxgy[k]
				}
				*(*float64)(unsafe.Add(unsafe.Pointer(ssim), unsafe.Sizeof(float64(0))*uintptr((j-4)*w+i))) = (mugxgy*2 + c2) / (mugx2 + mugy2 + c2)
				if i+1 < w {
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-1))&7)*stride+i+1)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-1))&7)*stride+i+1)))
							col_sums_gx2[0] = float64(gx) * float64(gx)
							col_sums_gy2[0] = float64(gy) * float64(gy)
							col_sums_gxgy[0] = float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+0)&7)*stride+i+1)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+0)&7)*stride+i+1)))
							col_sums_gx2[0] += float64(gx) * float64(gx)
							col_sums_gy2[0] += float64(gy) * float64(gy)
							col_sums_gxgy[0] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-3))&7)*stride+i+2)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-3))&7)*stride+i+2)))
							col_sums_gx2[2] -= float64(gx) * float64(gx)
							col_sums_gy2[2] -= float64(gy) * float64(gy)
							col_sums_gxgy[2] -= float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+2)&7)*stride+i+2)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+2)&7)*stride+i+2)))
							col_sums_gx2[2] -= float64(gx) * float64(gx)
							col_sums_gy2[2] -= float64(gy) * float64(gy)
							col_sums_gxgy[2] -= float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						col_sums_gx2[1] = col_sums_gx2[2] * 0.5
						col_sums_gy2[1] = col_sums_gy2[2] * 0.5
						col_sums_gxgy[1] = col_sums_gxgy[2] * 0.5
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-4))&7)*stride+i+3)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-4))&7)*stride+i+3)))
							col_sums_gx2[3] -= float64(gx) * float64(gx)
							col_sums_gy2[3] -= float64(gy) * float64(gy)
							col_sums_gxgy[3] -= float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+3)&7)*stride+i+3)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+3)&7)*stride+i+3)))
							col_sums_gx2[3] -= float64(gx) * float64(gx)
							col_sums_gy2[3] -= float64(gy) * float64(gy)
							col_sums_gxgy[3] -= float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						col_sums_gx2[2] = col_sums_gx2[3] * 0.5
						col_sums_gy2[2] = col_sums_gy2[3] * 0.5
						col_sums_gxgy[2] = col_sums_gxgy[3] * 0.5
						if true {
							break
						}
					}
					for {
						col_sums_gx2[3] = col_sums_gx2[4]
						col_sums_gy2[3] = col_sums_gy2[4]
						col_sums_gxgy[3] = col_sums_gxgy[4]
						if true {
							break
						}
					}
					for {
						col_sums_gx2[4] = col_sums_gx2[5] * 2
						col_sums_gy2[4] = col_sums_gy2[5] * 2
						col_sums_gxgy[4] = col_sums_gxgy[5] * 2
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-4))&7)*stride+i+5)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-4))&7)*stride+i+5)))
							col_sums_gx2[4] += float64(gx) * float64(gx)
							col_sums_gy2[4] += float64(gy) * float64(gy)
							col_sums_gxgy[4] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+3)&7)*stride+i+5)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+3)&7)*stride+i+5)))
							col_sums_gx2[4] += float64(gx) * float64(gx)
							col_sums_gy2[4] += float64(gy) * float64(gy)
							col_sums_gxgy[4] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						col_sums_gx2[5] = col_sums_gx2[6] * 2
						col_sums_gy2[5] = col_sums_gy2[6] * 2
						col_sums_gxgy[5] = col_sums_gxgy[6] * 2
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-3))&7)*stride+i+6)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-3))&7)*stride+i+6)))
							col_sums_gx2[5] += float64(gx) * float64(gx)
							col_sums_gy2[5] += float64(gy) * float64(gy)
							col_sums_gxgy[5] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+2)&7)*stride+i+6)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+2)&7)*stride+i+6)))
							col_sums_gx2[5] += float64(gx) * float64(gx)
							col_sums_gy2[5] += float64(gy) * float64(gy)
							col_sums_gxgy[5] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						col_sums_gx2[6] = col_sums_gx2[7] * 2
						col_sums_gy2[6] = col_sums_gy2[7] * 2
						col_sums_gxgy[6] = col_sums_gxgy[7] * 2
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-2))&7)*stride+i+7)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-2))&7)*stride+i+7)))
							col_sums_gx2[6] += float64(gx) * float64(gx)
							col_sums_gy2[6] += float64(gy) * float64(gy)
							col_sums_gxgy[6] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+1)&7)*stride+i+7)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+1)&7)*stride+i+7)))
							col_sums_gx2[6] += float64(gx) * float64(gx)
							col_sums_gy2[6] += float64(gy) * float64(gy)
							col_sums_gxgy[6] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-1))&7)*stride+i+8)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+int(-1))&7)*stride+i+8)))
							col_sums_gx2[7] = float64(gx) * float64(gx)
							col_sums_gy2[7] = float64(gy) * float64(gy)
							col_sums_gxgy[7] = float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
					for {
						{
							var (
								gx uint
								gy uint
							)
							gx = *(*uint)(unsafe.Add(unsafe.Pointer(gx_buf), unsafe.Sizeof(uint(0))*uintptr(((j+0)&7)*stride+i+8)))
							gy = *(*uint)(unsafe.Add(unsafe.Pointer(gy_buf), unsafe.Sizeof(uint(0))*uintptr(((j+0)&7)*stride+i+8)))
							col_sums_gx2[7] += float64(gx) * float64(gx)
							col_sums_gy2[7] += float64(gy) * float64(gy)
							col_sums_gxgy[7] += float64(gx) * float64(gy)
						}
						if true {
							break
						}
					}
				}
			}
		}
	}
}

var FS_WEIGHTS [4]float64 = [4]float64{0.2989654541015625, 0.3141326904296875, 0.2473602294921875, 0.1395416259765625}

func fs_average(_ctx *fs_ctx, _l int) float64 {
	var (
		ssim *float64
		ret  float64
		w    int
		h    int
		i    int
		j    int
	)
	w = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).W
	h = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).H
	ssim = (*(*fs_level)(unsafe.Add(unsafe.Pointer(_ctx.Level), unsafe.Sizeof(fs_level{})*uintptr(_l)))).Ssim
	ret = 0
	for j = 0; j < h; j++ {
		for i = 0; i < w; i++ {
			ret += *(*float64)(unsafe.Add(unsafe.Pointer(ssim), unsafe.Sizeof(float64(0))*uintptr(j*w+i)))
		}
	}
	return pow(ret/float64(w*h), FS_WEIGHTS[_l])
}
func convert_ssim_db(_ssim float64, _weight float64) float64 {
	if _weight >= _ssim {
	} else {
		__assert_fail(libc.CString("_weight >= _ssim"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	if (_weight - _ssim) < 1e-10 {
		return 100.0
	}
	return (log10(_weight) - log10(_weight-_ssim)) * 10
}
func calc_ssim(_src *uint8, _systride int, _dst *uint8, _dystride int, _w int, _h int, _bd uint32, _shift uint32) float64 {
	var (
		ctx fs_ctx
		ret float64
		l   int
	)
	ret = 1
	fs_ctx_init(&ctx, _w, _h, 4)
	fs_downsample_level0(&ctx, _src, _systride, _dst, _dystride, _w, _h, _bd, _shift)
	for l = 0; l < 4-1; l++ {
		fs_calc_structure(&ctx, l, int(_bd))
		ret *= fs_average(&ctx, l)
		fs_downsample_level(&ctx, l+1)
	}
	fs_calc_structure(&ctx, l, int(_bd))
	fs_apply_luminance(&ctx, l, int(_bd))
	ret *= fs_average(&ctx, l)
	fs_ctx_clear(&ctx)
	return ret
}
func vpx_calc_fastssim(source *YV12_BUFFER_CONFIG, dest *YV12_BUFFER_CONFIG, ssim_y *float64, ssim_u *float64, ssim_v *float64, bd uint32, in_bd uint32) float64 {
	var (
		ssimv    float64
		bd_shift uint32 = 0
	)
	vpx_clear_system_state()
	if bd >= in_bd {
	} else {
		__assert_fail(libc.CString("bd >= in_bd"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
	bd_shift = bd - in_bd
	*ssim_y = calc_ssim(source.Y_buffer, source.Y_stride, dest.Y_buffer, dest.Y_stride, source.Y_crop_width, source.Y_crop_height, in_bd, bd_shift)
	*ssim_u = calc_ssim(source.U_buffer, source.Uv_stride, dest.U_buffer, dest.Uv_stride, source.Uv_crop_width, source.Uv_crop_height, in_bd, bd_shift)
	*ssim_v = calc_ssim(source.V_buffer, source.Uv_stride, dest.V_buffer, dest.Uv_stride, source.Uv_crop_width, source.Uv_crop_height, in_bd, bd_shift)
	ssimv = (*ssim_y)*0.8 + ((*ssim_u)+(*ssim_v))*0.1
	return convert_ssim_db(ssimv, 1.0)
}
