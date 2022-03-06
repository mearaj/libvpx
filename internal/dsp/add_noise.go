package dsp

import (
	"math"
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
)

func VpxPlaneAddNoiseC(start *uint8, noise *int8, blackclamp int, whiteclamp int, width int, height int, pitch int) {
	var (
		i         int
		j         int
		bothclamp int = blackclamp + whiteclamp
	)
	for i = 0; i < height; i++ {
		var (
			pos *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer(start), i*pitch))
			ref *int8  = ((*int8)(unsafe.Add(unsafe.Pointer(noise), int(libc.Rand())&math.MaxUint8)))
		)
		for j = 0; j < width; j++ {
			var v int = int(*(*uint8)(unsafe.Add(unsafe.Pointer(pos), j)))
			v = clamp(v-blackclamp, 0, math.MaxUint8)
			v = clamp(v+bothclamp, 0, math.MaxUint8)
			v = clamp(v-whiteclamp, 0, math.MaxUint8)
			*(*uint8)(unsafe.Add(unsafe.Pointer(pos), j)) = uint8(int8(v + int(*(*int8)(unsafe.Add(unsafe.Pointer(ref), j)))))
		}
	}
}
func gaussian(sigma float64, mu float64, x float64) float64 {
	return 1 / (sigma * math.Sqrt(2.0*3.14159265)) * math.Exp(-(x-mu)*(x-mu)/(sigma*2*sigma))
}
func VpxSetupNoise(sigma float64, noise *int8, size int) int {
	var (
		char_dist [256]int8
		next      int = 0
		i         int
		j         int
	)
	for i = -32; i < 32; i++ {
		var a_i int = int(gaussian(sigma, 0, float64(i))*256 + 0.5)
		if a_i != 0 {
			for j = 0; j < a_i; j++ {
				if next+j >= 256 {
					goto set_noise
				}
				char_dist[next+j] = int8(i)
			}
			next = next + j
		}
	}
	for ; next < 256; next++ {
		char_dist[next] = 0
	}
set_noise:
	for i = 0; i < size; i++ {
		*(*int8)(unsafe.Add(unsafe.Pointer(noise), i)) = char_dist[int(libc.Rand())&math.MaxUint8]
	}
	return int(-char_dist[0])
}
