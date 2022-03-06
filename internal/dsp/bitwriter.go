package dsp

import (
	"math"
	"unsafe"
)

type vpx_writer struct {
	Lowvalue uint
	Range    uint
	Count    int
	Pos      uint
	Buffer   *uint8
}

func vpx_write(br *vpx_writer, bit int, probability int) {
	var (
		split    uint
		count    int  = br.Count
		range_   uint = br.Range
		lowvalue uint = br.Lowvalue
		shift    int
	)
	split = (((range_ - 1) * uint(probability)) >> 8) + 1
	range_ = split
	if bit != 0 {
		lowvalue += split
		range_ = br.Range - split
	}
	shift = int(vpx_norm[range_])
	range_ <<= uint(shift)
	count += shift
	if count >= 0 {
		var offset int = shift - count
		if (lowvalue<<uint(offset-1))&0x80000000 != 0 {
			var x int = int(br.Pos - 1)
			for x >= 0 && *(*uint8)(unsafe.Add(unsafe.Pointer(br.Buffer), x)) == math.MaxUint8 {
				*(*uint8)(unsafe.Add(unsafe.Pointer(br.Buffer), x)) = 0
				x--
			}
			*(*uint8)(unsafe.Add(unsafe.Pointer(br.Buffer), x)) += 1
		}
		*(*uint8)(unsafe.Add(unsafe.Pointer(br.Buffer), func() uint {
			p := &br.Pos
			x := *p
			*p++
			return x
		}())) = uint8((lowvalue >> uint(24-offset)) & math.MaxUint8)
		lowvalue <<= uint(offset)
		shift = count
		lowvalue &= 0xFFFFFF
		count -= 8
	}
	lowvalue <<= uint(shift)
	br.Count = count
	br.Lowvalue = lowvalue
	br.Range = range_
}
func vpx_write_bit(w *vpx_writer, bit int) {
	vpx_write(w, bit, 128)
}
func vpx_write_literal(w *vpx_writer, data int, bits int) {
	var bit int
	for bit = bits - 1; bit >= 0; bit-- {
		vpx_write_bit(w, (data>>bit)&1)
	}
}
func vpx_start_encode(br *vpx_writer, source *uint8) {
	br.Lowvalue = 0
	br.Range = math.MaxUint8
	br.Count = -24
	br.Buffer = source
	br.Pos = 0
	vpx_write_bit(br, 0)
}
func vpx_stop_encode(br *vpx_writer) {
	var i int
	for i = 0; i < 32; i++ {
		vpx_write_bit(br, 0)
	}
	if (*(*uint8)(unsafe.Add(unsafe.Pointer(br.Buffer), br.Pos-1)) & 224) == 192 {
		*(*uint8)(unsafe.Add(unsafe.Pointer(br.Buffer), func() uint {
			p := &br.Pos
			x := *p
			*p++
			return x
		}())) = 0
	}
}
