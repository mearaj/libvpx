package vp8

import (
	"github.com/gotranspile/cxgo/runtime/cmath"
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/dsp"
	"github.com/mearaj/libvpx/internal/scale"
	"unsafe"
)

func FilterByWeight(src *uint8, src_stride int, dst *uint8, dst_stride int, block_size int, src_weight int) {
	var (
		dst_weight       = (int(1 << MFQE_PRECISION)) - src_weight
		rounding_bit int = 1 << (int(MFQE_PRECISION - 1))
		r            int
		c            int
	)
	for r = 0; r < block_size; r++ {
		for c = 0; c < block_size; c++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), c)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(src), c)))*src_weight + int(*(*uint8)(unsafe.Add(unsafe.Pointer(dst), c)))*dst_weight + rounding_bit) >> MFQE_PRECISION))
		}
		src = (*uint8)(unsafe.Add(unsafe.Pointer(src), src_stride))
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), dst_stride))
	}
}
func Vp8FilterByWeight16x16C(src *uint8, src_stride int, dst *uint8, dst_stride int, src_weight int) {
	FilterByWeight(src, src_stride, dst, dst_stride, 16, src_weight)
}
func Vp8FilterByWeight8x8C(src *uint8, src_stride int, dst *uint8, dst_stride int, src_weight int) {
	FilterByWeight(src, src_stride, dst, dst_stride, 8, src_weight)
}
func Vp8FilterByWeight4x4C(src *uint8, src_stride int, dst *uint8, dst_stride int, src_weight int) {
	FilterByWeight(src, src_stride, dst, dst_stride, 4, src_weight)
}
func apply_ifactor(y_src *uint8, y_src_stride int, y_dst *uint8, y_dst_stride int, u_src *uint8, v_src *uint8, uv_src_stride int, u_dst *uint8, v_dst *uint8, uv_dst_stride int, block_size int, src_weight int) {
	if block_size == 16 {
		Vp8FilterByWeight16x16C(y_src, y_src_stride, y_dst, y_dst_stride, src_weight)
		Vp8FilterByWeight8x8C(u_src, uv_src_stride, u_dst, uv_dst_stride, src_weight)
		Vp8FilterByWeight8x8C(v_src, uv_src_stride, v_dst, uv_dst_stride, src_weight)
	} else {
		Vp8FilterByWeight8x8C(y_src, y_src_stride, y_dst, y_dst_stride, src_weight)
		Vp8FilterByWeight4x4C(u_src, uv_src_stride, u_dst, uv_dst_stride, src_weight)
		Vp8FilterByWeight4x4C(v_src, uv_src_stride, v_dst, uv_dst_stride, src_weight)
	}
}
func int_sqrt(x uint) uint {
	var (
		y     = x
		guess uint
		p     = 1
	)
	for func() uint {
		y >>= 1
		return y
	}() != 0 {
		p++
	}
	p >>= 1
	guess = 0
	for p >= 0 {
		guess |= uint(1 << p)
		if x < guess*guess {
			guess -= uint(1 << p)
		}
		p--
	}
	return guess + uint(libc.BoolToInt(guess*guess+guess+1 <= x))
}
func multiframe_quality_enhance_block(blksize int, qcurr int, qprev int, y *uint8, u *uint8, v *uint8, y_stride int, uv_stride int, yd *uint8, ud *uint8, vd *uint8, yd_stride int, uvd_stride int) {
	var (
		VP8_ZEROS = [16]uint8{}
		uvblksize = blksize >> 1
		qdiff               = qcurr - qprev
		i         int
		up        *uint8
		udp       *uint8
		vp        *uint8
		vdp       *uint8
		act       uint32
		actd      uint32
		sad       uint32
		usad      uint32
		vsad      uint32
		sse       uint32
		thr       uint32
		thrsq     uint32
		actrisk   uint32
	)
	if blksize == 16 {
		actd = (dsp.VpxVariance16x16C((*uint8)(unsafe.Pointer(yd)), yd_stride, (*uint8)(unsafe.Pointer(&VP8_ZEROS[0])), 0, &sse) + 128) >> 8
		act = (dsp.VpxVariance16x16C((*uint8)(unsafe.Pointer(y)), y_stride, (*uint8)(unsafe.Pointer(&VP8_ZEROS[0])), 0, &sse) + 128) >> 8
		dsp.VpxVariance16x16C((*uint8)(unsafe.Pointer(y)), y_stride, (*uint8)(unsafe.Pointer(yd)), yd_stride, &sse)
		sad = (sse + 128) >> 8
		dsp.VpxVariance8x8C((*uint8)(unsafe.Pointer(u)), uv_stride, (*uint8)(unsafe.Pointer(ud)), uvd_stride, &sse)
		usad = (sse + 32) >> 6
		dsp.VpxVariance8x8C((*uint8)(unsafe.Pointer(v)), uv_stride, (*uint8)(unsafe.Pointer(vd)), uvd_stride, &sse)
		vsad = (sse + 32) >> 6
	} else {
		actd = (dsp.VpxVariance8x8C((*uint8)(unsafe.Pointer(yd)), yd_stride, (*uint8)(unsafe.Pointer(&VP8_ZEROS[0])), 0, &sse) + 32) >> 6
		act = (dsp.VpxVariance8x8C((*uint8)(unsafe.Pointer(y)), y_stride, (*uint8)(unsafe.Pointer(&VP8_ZEROS[0])), 0, &sse) + 32) >> 6
		dsp.VpxVariance8x8C((*uint8)(unsafe.Pointer(y)), y_stride, (*uint8)(unsafe.Pointer(yd)), yd_stride, &sse)
		sad = (sse + 32) >> 6
		dsp.VpxVariance4x4C((*uint8)(unsafe.Pointer(u)), uv_stride, (*uint8)(unsafe.Pointer(ud)), uvd_stride, &sse)
		usad = (sse + 8) >> 4
		dsp.VpxVariance4x4C((*uint8)(unsafe.Pointer(v)), uv_stride, (*uint8)(unsafe.Pointer(vd)), uvd_stride, &sse)
		vsad = (sse + 8) >> 4
	}
	actrisk = uint32(libc.BoolToInt(actd > act*5))
	thr = uint32(qdiff >> 4)
	for func() uint {
		actd >>= 1
		return uint(actd)
	}() != 0 {
		thr++
	}
	for func() int {
		qprev >>= 2
		return qprev
	}() != 0 {
		thr++
	}
	thrsq = thr * thr
	if sad < thrsq && usad*4 < thrsq && vsad*4 < thrsq && actrisk == 0 {
		var ifactor int
		sad = uint32(int_sqrt(uint(sad)))
		ifactor = int((sad << MFQE_PRECISION) / thr)
		ifactor >>= qdiff >> 5
		if ifactor != 0 {
			apply_ifactor(y, y_stride, yd, yd_stride, u, v, uv_stride, ud, vd, uvd_stride, blksize, ifactor)
		}
	} else {
		if blksize == 16 {
			Vp8CopyMem16x16C(y, y_stride, yd, yd_stride)
			Vp8CopyMem8x8C(u, uv_stride, ud, uvd_stride)
			Vp8CopyMem8x8C(v, uv_stride, vd, uvd_stride)
		} else {
			Vp8CopyMem8x8C(y, y_stride, yd, yd_stride)
			for func() int {
				up = u
				udp = ud
				return func() int {
					i = 0
					return i
				}()
			}(); i < uvblksize; func() *uint8 {
				i++
				up = (*uint8)(unsafe.Add(unsafe.Pointer(up), uv_stride))
				return func() *uint8 {
					udp = (*uint8)(unsafe.Add(unsafe.Pointer(udp), uintptr(uvd_stride)))
					return udp
				}()
			}() {
				libc.MemCpy(unsafe.Pointer(udp), unsafe.Pointer(up), uvblksize)
			}
			for func() int {
				vp = v
				vdp = vd
				return func() int {
					i = 0
					return i
				}()
			}(); i < uvblksize; func() *uint8 {
				i++
				vp = (*uint8)(unsafe.Add(unsafe.Pointer(vp), uv_stride))
				return func() *uint8 {
					vdp = (*uint8)(unsafe.Add(unsafe.Pointer(udp), uintptr(uvd_stride)))
					return vdp
				}()
			}() {
				libc.MemCpy(unsafe.Pointer(vdp), unsafe.Pointer(vp), uvblksize)
			}
		}
	}
}
func qualify_inter_mb(mode_info_context *ModeInfo, map_ *int) int {
	if mode_info_context.Mbmi.Mb_skip_coeff != 0 {
		*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*0)) = func() int {
			p := (*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*1))
			*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*1)) = func() int {
				p := (*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*2))
				*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*2)) = func() int {
					p := (*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*3))
					*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*3)) = 1
					return *p
				}()
				return *p
			}()
			return *p
		}()
	} else if int(mode_info_context.Mbmi.Mode) == SPLITMV {
		var (
			ndx = [4][4]int{{0, 1, 4, 5}, {2, 3, 6, 7}, {8, 9, 12, 13}, {10, 11, 14, 15}}
			i   int
			j   int
		)
		*map_ = 0
		for i = 0; i < 4; i++ {
			*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*uintptr(i))) = 1
			for j = 0; j < 4 && *(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*uintptr(j))) != 0; j++ {
				*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*uintptr(i))) &= int(libc.BoolToInt(int(mode_info_context.Bmi[ndx[i][j]].Mv.As_mv.Row) <= 2 && int(mode_info_context.Bmi[ndx[i][j]].Mv.As_mv.Col) <= 2))
			}
		}
	} else {
		*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*0)) = func() int {
			p := (*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*1))
			*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*1)) = func() int {
				p := (*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*2))
				*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*2)) = func() int {
					p := (*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*3))
					*(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*3)) = int(libc.BoolToInt(int(mode_info_context.Mbmi.Mode) > B_PRED && cmath.Abs(int64(mode_info_context.Mbmi.Mv.As_mv.Row)) <= 2 && cmath.Abs(int64(mode_info_context.Mbmi.Mv.As_mv.Col)) <= 2))
					return *p
				}()
				return *p
			}()
			return *p
		}()
	}
	return *(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*0)) + *(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*1)) + *(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*2)) + *(*int)(unsafe.Add(unsafe.Pointer(map_), unsafe.Sizeof(int(0))*3))
}
func vp8_multiframe_quality_enhance(cm *VP8Common) {
	var (
		show = cm.Frame_to_show
		dest = &cm.Post_proc_buffer
		frame_type                         = cm.Frame_type
		mode_info_context                         = cm.Mi
		mb_row            int
		mb_col            int
		totmap            int
		map_  [4]int
		qcurr = cm.Base_qindex
		qprev = cm.Postproc_state.Last_base_qindex
		y_ptr *uint8
		u_ptr             *uint8
		v_ptr             *uint8
		yd_ptr            *uint8
		ud_ptr            *uint8
		vd_ptr            *uint8
	)
	y_ptr = (*uint8)(unsafe.Pointer(show.Y_buffer))
	u_ptr = (*uint8)(unsafe.Pointer(show.U_buffer))
	v_ptr = (*uint8)(unsafe.Pointer(show.V_buffer))
	yd_ptr = (*uint8)(unsafe.Pointer(dest.Y_buffer))
	ud_ptr = (*uint8)(unsafe.Pointer(dest.U_buffer))
	vd_ptr = (*uint8)(unsafe.Pointer(dest.V_buffer))
	for mb_row = 0; mb_row < cm.Mb_rows; mb_row++ {
		for mb_col = 0; mb_col < cm.Mb_cols; mb_col++ {
			if frame_type == int(INTER_FRAME) {
				totmap = qualify_inter_mb(mode_info_context, &map_[0])
			} else {
				if frame_type == int(KEY_FRAME) {
					totmap = 4
				} else {
					totmap = 0
				}
			}
			if totmap != 0 {
				if totmap < 4 {
					var (
						i int
						j int
					)
					for i = 0; i < 2; i++ {
						for j = 0; j < 2; j++ {
							if map_[i*2+j] != 0 {
								multiframe_quality_enhance_block(8, qcurr, qprev, (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), (i*show.Y_stride+j)*8)), (*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), (i*show.Uv_stride+j)*4)), (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), (i*show.Uv_stride+j)*4)), show.Y_stride, show.Uv_stride, (*uint8)(unsafe.Add(unsafe.Pointer(yd_ptr), (i*dest.Y_stride+j)*8)), (*uint8)(unsafe.Add(unsafe.Pointer(ud_ptr), (i*dest.Uv_stride+j)*4)), (*uint8)(unsafe.Add(unsafe.Pointer(vd_ptr), (i*dest.Uv_stride+j)*4)), dest.Y_stride, dest.Uv_stride)
							} else {
								var (
									k   int
									up  = (*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), (i*show.Uv_stride+j)*4))
									udp = (*uint8)(unsafe.Add(unsafe.Pointer(ud_ptr), (i*dest.Uv_stride+j)*4))
									vp         = (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), (i*show.Uv_stride+j)*4))
									vdp        = (*uint8)(unsafe.Add(unsafe.Pointer(vd_ptr), (i*dest.Uv_stride+j)*4))
								)
								Vp8CopyMem8x8C((*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), (i*show.Y_stride+j)*8)), show.Y_stride, (*uint8)(unsafe.Add(unsafe.Pointer(yd_ptr), (i*dest.Y_stride+j)*8)), dest.Y_stride)
								for k = 0; k < 4; func() *uint8 {
									k++
									up = (*uint8)(unsafe.Add(unsafe.Pointer(up), show.Uv_stride))
									udp = (*uint8)(unsafe.Add(unsafe.Pointer(udp), dest.Uv_stride))
									vp = (*uint8)(unsafe.Add(unsafe.Pointer(vp), show.Uv_stride))
									return func() *uint8 {
										vdp = (*uint8)(unsafe.Add(unsafe.Pointer(vdp), uintptr(dest.Uv_stride)))
										return vdp
									}()
								}() {
									libc.MemCpy(unsafe.Pointer(udp), unsafe.Pointer(up), 4)
									libc.MemCpy(unsafe.Pointer(vdp), unsafe.Pointer(vp), 4)
								}
							}
						}
					}
				} else {
					multiframe_quality_enhance_block(16, qcurr, qprev, y_ptr, u_ptr, v_ptr, show.Y_stride, show.Uv_stride, yd_ptr, ud_ptr, vd_ptr, dest.Y_stride, dest.Uv_stride)
				}
			} else {
				Vp8CopyMem16x16C(y_ptr, show.Y_stride, yd_ptr, dest.Y_stride)
				Vp8CopyMem8x8C(u_ptr, show.Uv_stride, ud_ptr, dest.Uv_stride)
				Vp8CopyMem8x8C(v_ptr, show.Uv_stride, vd_ptr, dest.Uv_stride)
			}
			y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), 16))
			u_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), 8))
			v_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), 8))
			yd_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(yd_ptr), 16))
			ud_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(ud_ptr), 8))
			vd_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(vd_ptr), 8))
			mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
		}
		y_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(y_ptr), show.Y_stride*16-cm.Mb_cols*16))
		u_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(u_ptr), show.Uv_stride*8-cm.Mb_cols*8))
		v_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(v_ptr), show.Uv_stride*8-cm.Mb_cols*8))
		yd_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(yd_ptr), dest.Y_stride*16-cm.Mb_cols*16))
		ud_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(ud_ptr), dest.Uv_stride*8-cm.Mb_cols*8))
		vd_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(vd_ptr), dest.Uv_stride*8-cm.Mb_cols*8))
		mode_info_context = (*ModeInfo)(unsafe.Add(unsafe.Pointer(mode_info_context), unsafe.Sizeof(ModeInfo{})*1))
	}
}
