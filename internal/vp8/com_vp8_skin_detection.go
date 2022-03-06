package vp8

import (
	"github.com/mearaj/libvpx/internal/dsp"
	"unsafe"
)

type SKIN_DETECTION_BLOCK_SIZE int

const (
	SKIN_8X8 = SKIN_DETECTION_BLOCK_SIZE(iota)
	SKIN_16X16
)

func avg_2x2(s *uint8, p int) int {
	var (
		i   int
		j   int
		sum int = 0
	)
	for i = 0; i < 2; func() *uint8 {
		i++
		return func() *uint8 {
			s = (*uint8)(unsafe.Add(unsafe.Pointer(s), uintptr(p)))
			return s
		}()
	}() {
		for j = 0; j < 2; j++ {
			sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(s), j)))
		}
	}
	return (sum + 2) >> 2
}
func vp8_compute_skin_block(y *uint8, u *uint8, v *uint8, stride int, strideuv int, bsize SKIN_DETECTION_BLOCK_SIZE, consec_zeromv int, curr_motion_magn int) int {
	if consec_zeromv > 60 && curr_motion_magn == 0 {
		return 0
	} else {
		var motion int = 1
		if consec_zeromv > 25 && curr_motion_magn == 0 {
			motion = 0
		}
		if bsize == SKIN_DETECTION_BLOCK_SIZE(SKIN_16X16) {
			var (
				ysource int = avg_2x2((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(y), stride*7))), 7)), stride)
				usource int = avg_2x2((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(u), strideuv*3))), 3)), strideuv)
				vsource int = avg_2x2((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(v), strideuv*3))), 3)), strideuv)
			)
			return dsp.VpxSkinPixel(ysource, usource, vsource, motion)
		} else {
			var (
				num_skin int = 0
				i        int
				j        int
			)
			for i = 0; i < 2; i++ {
				for j = 0; j < 2; j++ {
					var (
						ysource int = avg_2x2((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(y), stride*3))), 3)), stride)
						usource int = avg_2x2((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(u), strideuv))), 1)), strideuv)
						vsource int = avg_2x2((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(v), strideuv))), 1)), strideuv)
					)
					num_skin += dsp.VpxSkinPixel(ysource, usource, vsource, motion)
					if num_skin >= 2 {
						return 1
					}
					y = (*uint8)(unsafe.Add(unsafe.Pointer(y), 8))
					u = (*uint8)(unsafe.Add(unsafe.Pointer(u), 4))
					v = (*uint8)(unsafe.Add(unsafe.Pointer(v), 4))
				}
				y = (*uint8)(unsafe.Add(unsafe.Pointer(y), (stride<<3)-16))
				u = (*uint8)(unsafe.Add(unsafe.Pointer(u), (strideuv<<2)-8))
				v = (*uint8)(unsafe.Add(unsafe.Pointer(v), (strideuv<<2)-8))
			}
			return 0
		}
	}
}
