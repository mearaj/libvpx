package dsp

import "math"

func clip_pixel(val int) uint8 {
	if val > math.MaxUint8 {
		return math.MaxUint8
	}
	if val < 0 {
		return 0
	}
	return uint8(int8(val))
}
func clamp(value int, low int, high int) int {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
func fclamp(value float64, low float64, high float64) float64 {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
func lclamp(value int64, low int64, high int64) int64 {
	if value < low {
		return low
	}
	if value > high {
		return high
	}
	return value
}
func clip_pixel_highbd(val int, bd int) uint16 {
	switch bd {
	case 8:
		fallthrough
	default:
		return uint16(int16(clamp(val, 0, math.MaxUint8)))
	case 10:
		return uint16(int16(clamp(val, 0, 1023)))
	case 12:
		return uint16(int16(clamp(val, 0, 4095)))
	}
}
