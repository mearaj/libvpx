package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func D207Predictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		r int
		c int
	)
	_ = above
	for r = 0; r < bs-1; r++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), r*int(stride))) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r+1))) + 1) >> 1))
	}
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), (bs-1)*int(stride))) = *(*uint8)(unsafe.Add(unsafe.Pointer(left), bs-1))
	dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), 1))
	for r = 0; r < bs-2; r++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), r*int(stride))) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r+1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r+2))) + 2) >> 2))
	}
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), (bs-2)*int(stride))) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), bs-2))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), bs-1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), bs-1))) + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), (bs-1)*int(stride))) = *(*uint8)(unsafe.Add(unsafe.Pointer(left), bs-1))
	dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), 1))
	for c = 0; c < bs-2; c++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), (bs-1)*int(stride)+c)) = *(*uint8)(unsafe.Add(unsafe.Pointer(left), bs-1))
	}
	for r = bs - 2; r >= 0; r-- {
		for c = 0; c < bs-2; c++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), r*int(stride)+c)) = *(*uint8)(unsafe.Add(unsafe.Pointer(dst), (r+1)*int(stride)+c-2))
		}
	}
}
func D63Predictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		r    int
		c    int
		size int
	)
	_ = left
	for c = 0; c < bs; c++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), c)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c+1))) + 1) >> 1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride+int64(c))) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c+1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c+2))) + 2) >> 2))
	}
	for func() int {
		r = 2
		return func() int {
			size = bs - 2
			return size
		}()
	}(); r < bs; func() int {
		r += 2
		return func() int {
			p := &size
			*p--
			return *p
		}()
	}() {
		libc.MemCpy(unsafe.Add(unsafe.Pointer(dst), (r+0)*int(stride)), unsafe.Add(unsafe.Pointer(dst), r>>1), size)
		libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(dst), (r+0)*int(stride)))), size), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(above), bs-1))), bs-size)
		libc.MemCpy(unsafe.Add(unsafe.Pointer(dst), (r+1)*int(stride)), unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))), r>>1), size)
		libc.MemSet(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(dst), (r+1)*int(stride)))), size), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(above), bs-1))), bs-size)
	}
}
func D45Predictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		above_right uint8  = *(*uint8)(unsafe.Add(unsafe.Pointer(above), bs-1))
		dst_row0    *uint8 = dst
		x           int
		size        int
	)
	_ = left
	for x = 0; x < bs-1; x++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), x)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), x))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), x+1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), x+2))) + 2) >> 2))
	}
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), bs-1)) = above_right
	dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	for func() int {
		x = 1
		return func() int {
			size = bs - 2
			return size
		}()
	}(); x < bs; func() int {
		x++
		return func() int {
			p := &size
			*p--
			return *p
		}()
	}() {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Add(unsafe.Pointer(dst_row0), x), size)
		libc.MemSet(unsafe.Add(unsafe.Pointer(dst), size), byte(above_right), x+1)
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func D117Predictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		r int
		c int
	)
	for c = 0; c < bs; c++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), c)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c))) + 1) >> 1))
	}
	dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 0)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0))) + 2) >> 2))
	for c = 1; c < bs; c++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), c)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c-2))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c-1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c))) + 2) >> 2))
	}
	dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 0)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 1))) + 2) >> 2))
	for r = 3; r < bs; r++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), (r-2)*int(stride))) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r-3))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r-2)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r-1))) + 2) >> 2))
	}
	for r = 2; r < bs; r++ {
		for c = 1; c < bs; c++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), c)) = *(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*(-2)+int64(c)-1))
		}
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func D135Predictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		i      int
		border [63]uint8
	)
	for i = 0; i < bs-2; i++ {
		border[i] = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), bs-3-i))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), bs-2-i)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), bs-1-i))) + 2) >> 2))
	}
	border[bs-2] = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 1))) + 2) >> 2))
	border[bs-1] = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0))) + 2) >> 2))
	border[bs-0] = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 1))) + 2) >> 2))
	for i = 0; i < bs-2; i++ {
		border[bs+1+i] = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), i))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), i+1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), i+2))) + 2) >> 2))
	}
	for i = 0; i < bs; i++ {
		libc.MemCpy(unsafe.Add(unsafe.Pointer(dst), i*int(stride)), unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(&border[bs]), -1))), -i), bs)
	}
}
func D153Predictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		r int
		c int
	)
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 0)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0))) + 1) >> 1))
	for r = 1; r < bs; r++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), r*int(stride))) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r))) + 1) >> 1))
	}
	dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), 1))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 0)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0))) + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 1))) + 2) >> 2))
	for r = 2; r < bs; r++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), r*int(stride))) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r-2))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r-1)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r))) + 2) >> 2))
	}
	dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), 1))
	for c = 0; c < bs-2; c++ {
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), c)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c)))*2 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c+1))) + 2) >> 2))
	}
	dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	for r = 1; r < bs; r++ {
		for c = 0; c < bs-2; c++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), c)) = *(*uint8)(unsafe.Add(unsafe.Pointer(dst), -stride+int64(c)-2))
		}
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func VPredictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var r int
	_ = left
	for r = 0; r < bs; r++ {
		libc.MemCpy(unsafe.Pointer(dst), unsafe.Pointer(above), bs)
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func HPredictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var r int
	_ = above
	for r = 0; r < bs; r++ {
		libc.MemSet(unsafe.Pointer(dst), byte(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r))), bs)
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func TmPredictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		r         int
		c         int
		ytop_left int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1)))
	)
	for r = 0; r < bs; r++ {
		for c = 0; c < bs; c++ {
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), c)) = clip_pixel(int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), r))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), c))) - ytop_left)
		}
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func Dc128Predictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var r int
	_ = above
	_ = left
	for r = 0; r < bs; r++ {
		libc.MemSet(unsafe.Pointer(dst), 128, bs)
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func DcLeftPredictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		i           int
		r           int
		expected_dc int
		sum         int = 0
	)
	_ = above
	for i = 0; i < bs; i++ {
		sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), i)))
	}
	expected_dc = (sum + (bs >> 1)) / bs
	for r = 0; r < bs; r++ {
		libc.MemSet(unsafe.Pointer(dst), byte(int8(expected_dc)), bs)
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func DcTopPredictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		i           int
		r           int
		expected_dc int
		sum         int = 0
	)
	_ = left
	for i = 0; i < bs; i++ {
		sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), i)))
	}
	expected_dc = (sum + (bs >> 1)) / bs
	for r = 0; r < bs; r++ {
		libc.MemSet(unsafe.Pointer(dst), byte(int8(expected_dc)), bs)
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func DcPredictor(dst *uint8, stride int64, bs int, above *uint8, left *uint8) {
	var (
		i           int
		r           int
		expected_dc int
		sum         int = 0
		count       int = bs * 2
	)
	for i = 0; i < bs; i++ {
		sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), i)))
		sum += int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), i)))
	}
	expected_dc = (sum + (count >> 1)) / count
	for r = 0; r < bs; r++ {
		libc.MemSet(unsafe.Pointer(dst), byte(int8(expected_dc)), bs)
		dst = (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride))
	}
}
func VpxHePredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		H int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1)))
		I int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0)))
		J int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 1)))
		K int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 2)))
		L int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 3)))
	)
	libc.MemSet(unsafe.Add(unsafe.Pointer(dst), stride*0), byte(int8((H+I*2+J+2)>>2)), 4)
	libc.MemSet(unsafe.Add(unsafe.Pointer(dst), stride*1), byte(int8((I+J*2+K+2)>>2)), 4)
	libc.MemSet(unsafe.Add(unsafe.Pointer(dst), stride*2), byte(int8((J+K*2+L+2)>>2)), 4)
	libc.MemSet(unsafe.Add(unsafe.Pointer(dst), stride*3), byte(int8((K+L*2+L+2)>>2)), 4)
}
func VpxVePredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		H int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1)))
		I int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0)))
		J int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 1)))
		K int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 2)))
		L int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 3)))
		M int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 4)))
	)
	_ = left
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 0)) = uint8(int8((H + I*2 + J + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 1)) = uint8(int8((I + J*2 + K + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 2)) = uint8(int8((J + K*2 + L + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), 3)) = uint8(int8((K + L*2 + M + 2) >> 2))
	libc.MemCpy(unsafe.Add(unsafe.Pointer(dst), stride*1), unsafe.Pointer(dst), 4)
	libc.MemCpy(unsafe.Add(unsafe.Pointer(dst), stride*2), unsafe.Pointer(dst), 4)
	libc.MemCpy(unsafe.Add(unsafe.Pointer(dst), stride*3), unsafe.Pointer(dst), 4)
}
func VpxD207Predictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		I int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0)))
		J int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 1)))
		K int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 2)))
		L int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 3)))
	)
	_ = above
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+0)) = uint8(int8((I + J + 1) >> 1))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0)) = uint8(int8((J + K + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0)) = uint8(int8((K + L + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+1)) = uint8(int8((I + J*2 + K + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1)) = uint8(int8((J + K*2 + L + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1)) = uint8(int8((K + L*2 + L + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0)) = func() uint8 {
				p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1))
				*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1)) = func() uint8 {
					p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2))
					*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2)) = func() uint8 {
						p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3))
						*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3)) = uint8(int8(L))
						return *p
					}()
					return *p
				}()
				return *p
			}()
			return *p
		}()
		return *p
	}()
}
func VpxD63Predictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		A int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0)))
		B int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 1)))
		C int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 2)))
		D int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 3)))
		E int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 4)))
		F int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 5)))
		G int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 6)))
	)
	_ = left
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+0)) = uint8(int8((A + B + 1) >> 1))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0)) = uint8(int8((B + C + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1)) = uint8(int8((C + D + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2)) = uint8(int8((D + E + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3)) = uint8(int8((E + F + 1) >> 1))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0)) = uint8(int8((A + B*2 + C + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0)) = uint8(int8((B + C*2 + D + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1)) = uint8(int8((C + D*2 + E + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2)) = uint8(int8((D + E*2 + F + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3)) = uint8(int8((E + F*2 + G + 2) >> 2))
}
func VpxD63ePredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		A int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0)))
		B int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 1)))
		C int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 2)))
		D int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 3)))
		E int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 4)))
		F int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 5)))
		G int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 6)))
		H int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 7)))
	)
	_ = left
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+0)) = uint8(int8((A + B + 1) >> 1))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0)) = uint8(int8((B + C + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1)) = uint8(int8((C + D + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2)) = uint8(int8((D + E + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3)) = uint8(int8((E + F*2 + G + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0)) = uint8(int8((A + B*2 + C + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0)) = uint8(int8((B + C*2 + D + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1)) = uint8(int8((C + D*2 + E + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2)) = uint8(int8((D + E*2 + F + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3)) = uint8(int8((F + G*2 + H + 2) >> 2))
}
func VpxD45Predictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		A int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0)))
		B int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 1)))
		C int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 2)))
		D int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 3)))
		E int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 4)))
		F int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 5)))
		G int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 6)))
		H int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 7)))
	)
	_ = stride
	_ = left
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+0)) = uint8(int8((A + B*2 + C + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0)) = uint8(int8((B + C*2 + D + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0)) = uint8(int8((C + D*2 + E + 2) >> 2))
			return *p
		}()
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1)) = func() uint8 {
				p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0))
				*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0)) = uint8(int8((D + E*2 + F + 2) >> 2))
				return *p
			}()
			return *p
		}()
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1)) = uint8(int8((E + F*2 + G + 2) >> 2))
			return *p
		}()
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2)) = uint8(int8((F + G*2 + H + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3)) = uint8(int8(H))
}
func VpxD45ePredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		A int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0)))
		B int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 1)))
		C int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 2)))
		D int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 3)))
		E int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 4)))
		F int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 5)))
		G int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 6)))
		H int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 7)))
	)
	_ = stride
	_ = left
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+0)) = uint8(int8((A + B*2 + C + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0)) = uint8(int8((B + C*2 + D + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0)) = uint8(int8((C + D*2 + E + 2) >> 2))
			return *p
		}()
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1)) = func() uint8 {
				p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0))
				*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0)) = uint8(int8((D + E*2 + F + 2) >> 2))
				return *p
			}()
			return *p
		}()
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1)) = uint8(int8((E + F*2 + G + 2) >> 2))
			return *p
		}()
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2)) = uint8(int8((F + G*2 + H + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3)) = uint8(int8((G + H*2 + H + 2) >> 2))
}
func VpxD117Predictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		I int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0)))
		J int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 1)))
		K int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 2)))
		X int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1)))
		A int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0)))
		B int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 1)))
		C int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 2)))
		D int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 3)))
	)
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+0)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1)) = uint8(int8((X + A + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2)) = uint8(int8((A + B + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3)) = uint8(int8((B + C + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+3)) = uint8(int8((C + D + 1) >> 1))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0)) = uint8(int8((K + J*2 + I + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0)) = uint8(int8((J + I*2 + X + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1)) = uint8(int8((I + X*2 + A + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2)) = uint8(int8((X + A*2 + B + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3)) = uint8(int8((A + B*2 + C + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+3)) = uint8(int8((B + C*2 + D + 2) >> 2))
}
func VpxD135Predictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		I int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0)))
		J int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 1)))
		K int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 2)))
		L int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 3)))
		X int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1)))
		A int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0)))
		B int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 1)))
		C int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 2)))
		D int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 3)))
	)
	_ = stride
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0)) = uint8(int8((J + K*2 + L + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0)) = uint8(int8((I + J*2 + K + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0)) = uint8(int8((X + I*2 + J + 2) >> 2))
			return *p
		}()
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1)) = func() uint8 {
				p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+0))
				*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+0)) = uint8(int8((A + X*2 + I + 2) >> 2))
				return *p
			}()
			return *p
		}()
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2)) = func() uint8 {
			p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+1))
			*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+1)) = uint8(int8((B + A*2 + X + 2) >> 2))
			return *p
		}()
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+3)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+2)) = uint8(int8((C + B*2 + A + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+3)) = uint8(int8((D + C*2 + B + 2) >> 2))
}
func VpxD153Predictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	var (
		I int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 0)))
		J int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 1)))
		K int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 2)))
		L int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(left), 3)))
		X int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), -1)))
		A int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 0)))
		B int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 1)))
		C int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(above), 2)))
	)
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+0)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+2)) = uint8(int8((I + X + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+0)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+2)) = uint8(int8((J + I + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+0)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+2)) = uint8(int8((K + J + 1) >> 1))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+0)) = uint8(int8((L + K + 1) >> 1))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+3)) = uint8(int8((A + B*2 + C + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+2)) = uint8(int8((X + A*2 + B + 2) >> 2))
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*0+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+3))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+3)) = uint8(int8((I + X*2 + A + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*1+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+3)) = uint8(int8((J + I*2 + X + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*2+1)) = func() uint8 {
		p := (*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3))
		*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+3)) = uint8(int8((K + J*2 + I + 2) >> 2))
		return *p
	}()
	*(*uint8)(unsafe.Add(unsafe.Pointer(dst), stride*3+1)) = uint8(int8((L + K*2 + J + 2) >> 2))
}
func VpxD207Predictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D207Predictor(dst, stride, 8, above, left)
}
func VpxD207Predictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D207Predictor(dst, stride, 16, above, left)
}
func VpxD207Predictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D207Predictor(dst, stride, 32, above, left)
}
func VpxD63Predictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D63Predictor(dst, stride, 8, above, left)
}
func VpxD63Predictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D63Predictor(dst, stride, 16, above, left)
}
func VpxD63Predictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D63Predictor(dst, stride, 32, above, left)
}
func VpxD45Predictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D45Predictor(dst, stride, 8, above, left)
}
func VpxD45Predictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D45Predictor(dst, stride, 16, above, left)
}
func VpxD45Predictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D45Predictor(dst, stride, 32, above, left)
}
func VpxD117Predictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D117Predictor(dst, stride, 8, above, left)
}
func VpxD117Predictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D117Predictor(dst, stride, 16, above, left)
}
func VpxD117Predictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D117Predictor(dst, stride, 32, above, left)
}
func VpxD135Predictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D135Predictor(dst, stride, 8, above, left)
}
func VpxD135Predictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D135Predictor(dst, stride, 16, above, left)
}
func VpxD135Predictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D135Predictor(dst, stride, 32, above, left)
}
func Vpxd153Predictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D153Predictor(dst, stride, 8, above, left)
}
func Vpxd153Predictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	D153Predictor(dst, stride, 16, above, left)
}
func Vpxd153Predictor32x32_c(dst *uint8, stride int64, above *uint8, left *uint8) {
	D153Predictor(dst, stride, 32, above, left)
}
func VpxVPredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	VPredictor(dst, stride, 4, above, left)
}
func VpxVPredictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	VPredictor(dst, stride, 8, above, left)
}
func VpxVPredictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	VPredictor(dst, stride, 16, above, left)
}
func VpxVPredictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	VPredictor(dst, stride, 32, above, left)
}
func VpxHPredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	HPredictor(dst, stride, 4, above, left)
}
func VpxHPredictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	HPredictor(dst, stride, 8, above, left)
}
func VpxHPredictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	HPredictor(dst, stride, 16, above, left)
}
func VpxHPredictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	HPredictor(dst, stride, 32, above, left)
}
func VpxTmPredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	TmPredictor(dst, stride, 4, above, left)
}
func VpxTmPredictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	TmPredictor(dst, stride, 8, above, left)
}
func VpxTmPredictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	TmPredictor(dst, stride, 16, above, left)
}
func VpxTmPredictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	TmPredictor(dst, stride, 32, above, left)
}
func VpxDc128Predictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	Dc128Predictor(dst, stride, 4, above, left)
}
func VpxDc128Predictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	Dc128Predictor(dst, stride, 8, above, left)
}
func VpxDc128Predictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	Dc128Predictor(dst, stride, 16, above, left)
}
func VpxDc128Predictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	Dc128Predictor(dst, stride, 32, above, left)
}
func VpxDcLeftPredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcLeftPredictor(dst, stride, 4, above, left)
}
func VpxDcLeftPredictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcLeftPredictor(dst, stride, 8, above, left)
}
func VpxDcLeftPredictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcLeftPredictor(dst, stride, 16, above, left)
}
func VpxDcLeftPredictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcLeftPredictor(dst, stride, 32, above, left)
}
func VpxDcTopPredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcTopPredictor(dst, stride, 4, above, left)
}
func VpxDcTopPredictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcTopPredictor(dst, stride, 8, above, left)
}
func VpxDcTopPredictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcTopPredictor(dst, stride, 16, above, left)
}
func VpxDcTopPredictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcTopPredictor(dst, stride, 32, above, left)
}
func VpxDcPredictor4x4C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcPredictor(dst, stride, 4, above, left)
}
func VpxDcPredictor8x8C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcPredictor(dst, stride, 8, above, left)
}
func VpxDcPredictor16x16C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcPredictor(dst, stride, 16, above, left)
}
func VpxDcPredictor32x32C(dst *uint8, stride int64, above *uint8, left *uint8) {
	DcPredictor(dst, stride, 32, above, left)
}
