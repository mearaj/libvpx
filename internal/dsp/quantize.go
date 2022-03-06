package dsp

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func vpx_quantize_dc(coeff_ptr *int16, n_coeffs int, skip_block int, round_ptr *int16, quant int16, qcoeff_ptr *int16, dqcoeff_ptr *int16, dequant int16, eob_ptr *uint16) {
	var (
		rc         int = 0
		coeff      int = int(*(*int16)(unsafe.Add(unsafe.Pointer(coeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))))
		coeff_sign int = (coeff >> 31)
		abs_coeff  int = (coeff ^ coeff_sign) - coeff_sign
		tmp        int
		eob        int = -1
	)
	libc.MemSet(unsafe.Pointer(qcoeff_ptr), 0, n_coeffs*int(unsafe.Sizeof(int16(0))))
	libc.MemSet(unsafe.Pointer(dqcoeff_ptr), 0, n_coeffs*int(unsafe.Sizeof(int16(0))))
	if skip_block == 0 {
		tmp = clamp(abs_coeff+int(*(*int16)(unsafe.Add(unsafe.Pointer(round_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0)))), math.MinInt16, math.MaxInt16)
		tmp = (tmp * int(quant)) >> 16
		*(*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))) = int16((tmp ^ coeff_sign) - coeff_sign)
		*(*int16)(unsafe.Add(unsafe.Pointer(dqcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))) = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc)))) * int(dequant))
		if tmp != 0 {
			eob = 0
		}
	}
	*eob_ptr = uint16(int16(eob + 1))
}
func vpx_quantize_dc_32x32(coeff_ptr *int16, skip_block int, round_ptr *int16, quant int16, qcoeff_ptr *int16, dqcoeff_ptr *int16, dequant int16, eob_ptr *uint16) {
	var (
		n_coeffs   int = 1024
		rc         int = 0
		coeff      int = int(*(*int16)(unsafe.Add(unsafe.Pointer(coeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))))
		coeff_sign int = (coeff >> 31)
		abs_coeff  int = (coeff ^ coeff_sign) - coeff_sign
		tmp        int
		eob        int = -1
	)
	libc.MemSet(unsafe.Pointer(qcoeff_ptr), 0, n_coeffs*int(unsafe.Sizeof(int16(0))))
	libc.MemSet(unsafe.Pointer(dqcoeff_ptr), 0, n_coeffs*int(unsafe.Sizeof(int16(0))))
	if skip_block == 0 {
		tmp = clamp(abs_coeff+((int(*(*int16)(unsafe.Add(unsafe.Pointer(round_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0))))+(1<<(1-1)))>>1), math.MinInt16, math.MaxInt16)
		tmp = (tmp * int(quant)) >> 15
		*(*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))) = int16((tmp ^ coeff_sign) - coeff_sign)
		*(*int16)(unsafe.Add(unsafe.Pointer(dqcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))) = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc)))) * int(dequant) / 2)
		if tmp != 0 {
			eob = 0
		}
	}
	*eob_ptr = uint16(int16(eob + 1))
}
func vpx_quantize_b_c(coeff_ptr *int16, n_coeffs int64, skip_block int, zbin_ptr *int16, round_ptr *int16, quant_ptr *int16, quant_shift_ptr *int16, qcoeff_ptr *int16, dqcoeff_ptr *int16, dequant_ptr *int16, eob_ptr *uint16, scan *int16, iscan *int16) {
	var (
		i              int
		non_zero_count int    = int(n_coeffs)
		eob            int    = -1
		zbins          [2]int = [2]int{int(*(*int16)(unsafe.Add(unsafe.Pointer(zbin_ptr), unsafe.Sizeof(int16(0))*0))), int(*(*int16)(unsafe.Add(unsafe.Pointer(zbin_ptr), unsafe.Sizeof(int16(0))*1)))}
		nzbins         [2]int = [2]int{zbins[0] * (-1), zbins[1] * (-1)}
	)
	_ = iscan
	_ = skip_block
	libc.Assert(skip_block == 0)
	libc.MemSet(unsafe.Pointer(qcoeff_ptr), 0, int(n_coeffs*int64(unsafe.Sizeof(int16(0)))))
	libc.MemSet(unsafe.Pointer(dqcoeff_ptr), 0, int(n_coeffs*int64(unsafe.Sizeof(int16(0)))))
	for i = int(n_coeffs) - 1; i >= 0; i-- {
		var (
			rc    int = int(*(*int16)(unsafe.Add(unsafe.Pointer(scan), unsafe.Sizeof(int16(0))*uintptr(i))))
			coeff int = int(*(*int16)(unsafe.Add(unsafe.Pointer(coeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))))
		)
		if coeff < zbins[rc != 0] && coeff > nzbins[rc != 0] {
			non_zero_count--
		} else {
			break
		}
	}
	for i = 0; i < non_zero_count; i++ {
		var (
			rc         int = int(*(*int16)(unsafe.Add(unsafe.Pointer(scan), unsafe.Sizeof(int16(0))*uintptr(i))))
			coeff      int = int(*(*int16)(unsafe.Add(unsafe.Pointer(coeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))))
			coeff_sign int = (coeff >> 31)
			abs_coeff  int = (coeff ^ coeff_sign) - coeff_sign
		)
		if abs_coeff >= zbins[rc != 0] {
			var tmp int = clamp(abs_coeff+int(*(*int16)(unsafe.Add(unsafe.Pointer(round_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0)))), math.MinInt16, math.MaxInt16)
			tmp = ((((tmp * int(*(*int16)(unsafe.Add(unsafe.Pointer(quant_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0))))) >> 16) + tmp) * int(*(*int16)(unsafe.Add(unsafe.Pointer(quant_shift_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0))))) >> 16
			*(*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))) = int16((tmp ^ coeff_sign) - coeff_sign)
			*(*int16)(unsafe.Add(unsafe.Pointer(dqcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))) = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc)))) * int(*(*int16)(unsafe.Add(unsafe.Pointer(dequant_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0)))))
			if tmp != 0 {
				eob = i
			}
		}
	}
	*eob_ptr = uint16(int16(eob + 1))
}
func vpx_quantize_b_32x32_c(coeff_ptr *int16, n_coeffs int64, skip_block int, zbin_ptr *int16, round_ptr *int16, quant_ptr *int16, quant_shift_ptr *int16, qcoeff_ptr *int16, dqcoeff_ptr *int16, dequant_ptr *int16, eob_ptr *uint16, scan *int16, iscan *int16) {
	var (
		zbins   [2]int = [2]int{((int(*(*int16)(unsafe.Add(unsafe.Pointer(zbin_ptr), unsafe.Sizeof(int16(0))*0))) + (1 << (1 - 1))) >> 1), ((int(*(*int16)(unsafe.Add(unsafe.Pointer(zbin_ptr), unsafe.Sizeof(int16(0))*1))) + (1 << (1 - 1))) >> 1)}
		nzbins  [2]int = [2]int{zbins[0] * (-1), zbins[1] * (-1)}
		idx     int    = 0
		idx_arr [1024]int
		i       int
		eob     int = -1
	)
	_ = iscan
	_ = skip_block
	libc.Assert(skip_block == 0)
	libc.MemSet(unsafe.Pointer(qcoeff_ptr), 0, int(n_coeffs*int64(unsafe.Sizeof(int16(0)))))
	libc.MemSet(unsafe.Pointer(dqcoeff_ptr), 0, int(n_coeffs*int64(unsafe.Sizeof(int16(0)))))
	for i = 0; i < int(n_coeffs); i++ {
		var (
			rc    int = int(*(*int16)(unsafe.Add(unsafe.Pointer(scan), unsafe.Sizeof(int16(0))*uintptr(i))))
			coeff int = int(*(*int16)(unsafe.Add(unsafe.Pointer(coeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))))
		)
		if coeff >= zbins[rc != 0] || coeff <= nzbins[rc != 0] {
			idx_arr[func() int {
				p := &idx
				x := *p
				*p++
				return x
			}()] = i
		}
	}
	for i = 0; i < idx; i++ {
		var (
			rc         int = int(*(*int16)(unsafe.Add(unsafe.Pointer(scan), unsafe.Sizeof(int16(0))*uintptr(idx_arr[i]))))
			coeff      int = int(*(*int16)(unsafe.Add(unsafe.Pointer(coeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))))
			coeff_sign int = (coeff >> 31)
			tmp        int
			abs_coeff  int = (coeff ^ coeff_sign) - coeff_sign
		)
		abs_coeff += (int(*(*int16)(unsafe.Add(unsafe.Pointer(round_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0)))) + (1 << (1 - 1))) >> 1
		abs_coeff = clamp(abs_coeff, math.MinInt16, math.MaxInt16)
		tmp = ((((abs_coeff * int(*(*int16)(unsafe.Add(unsafe.Pointer(quant_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0))))) >> 16) + abs_coeff) * int(*(*int16)(unsafe.Add(unsafe.Pointer(quant_shift_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0))))) >> 15
		*(*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))) = int16((tmp ^ coeff_sign) - coeff_sign)
		*(*int16)(unsafe.Add(unsafe.Pointer(dqcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))) = int16(clamp(int(*(*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*uintptr(rc))))*int(*(*int16)(unsafe.Add(unsafe.Pointer(dequant_ptr), unsafe.Sizeof(int16(0))*uintptr(rc != 0))))/2, math.MinInt16, math.MaxInt16))
		if tmp != 0 {
			eob = idx_arr[i]
		}
	}
	*eob_ptr = uint16(int16(eob + 1))
}
