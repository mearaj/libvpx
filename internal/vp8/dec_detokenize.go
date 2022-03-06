package vp8

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const NUM_PROBAS = 11
const NUM_CTX = 3

func vp8_reset_mb_tokens_context(x *MacroBlockd) {
	var (
		a_ctx *byte = ((*byte)(unsafe.Pointer(x.Above_context)))
		l_ctx *byte = ((*byte)(unsafe.Pointer(x.Left_context)))
	)
	libc.MemSet(unsafe.Pointer(a_ctx), 0, int(unsafe.Sizeof(ENTROPY_CONTEXT_PLANES{})-1))
	libc.MemSet(unsafe.Pointer(l_ctx), 0, int(unsafe.Sizeof(ENTROPY_CONTEXT_PLANES{})-1))
	if int(x.Mode_info_context.Mbmi.Is_4x4) == 0 {
		*(*byte)(unsafe.Add(unsafe.Pointer(a_ctx), 8)) = func() byte {
			p := (*byte)(unsafe.Add(unsafe.Pointer(l_ctx), 8))
			*(*byte)(unsafe.Add(unsafe.Pointer(l_ctx), 8)) = 0
			return *p
		}()
	}
}

var kBands [17]uint8 = [17]uint8{0, 1, 2, 3, 6, 4, 5, 6, 6, 6, 6, 6, 6, 6, 6, 7, 0}
var kCat3 [4]uint8 = [4]uint8{173, 148, 140, 0}
var kCat4 [5]uint8 = [5]uint8{176, 155, 140, 135, 0}
var kCat5 [6]uint8 = [6]uint8{180, 157, 141, 134, 130, 0}
var kCat6 [12]uint8 = [12]uint8{254, 254, 243, 230, 196, 177, 153, 140, 133, 130, 129, 0}
var kCat3456 [4]*uint8 = [4]*uint8{&kCat3[0], &kCat4[0], &kCat5[0], &kCat6[0]}
var kZigzag [16]uint8 = [16]uint8{0, 1, 4, 8, 5, 2, 3, 6, 9, 12, 13, 10, 7, 11, 14, 15}

type ProbaArray *[3][11]uint8

func GetSigned(br *BOOL_DECODER, value_to_sign int) int {
	var (
		split    int    = int((br.Range + 1) >> 1)
		bigsplit uint64 = uint64(split) << uint64((CHAR_BIT*int(unsafe.Sizeof(uint64(0))))-8)
		v        int
	)
	if br.Count < 0 {
		vp8dx_bool_decoder_fill(br)
	}
	if br.Value < bigsplit {
		br.Range = uint(split)
		v = value_to_sign
	} else {
		br.Range = br.Range - uint(split)
		br.Value = br.Value - bigsplit
		v = -value_to_sign
	}
	br.Range += br.Range
	br.Value += br.Value
	br.Count--
	return v
}
func GetCoeffs(br *BOOL_DECODER, prob ProbaArray, ctx int, n int, out *int16) int {
	var p *uint8 = &(*(*[3][11]uint8)(unsafe.Add(unsafe.Pointer(prob), unsafe.Sizeof([3][11]uint8{})*uintptr(n))))[ctx][0]
	if vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 0)))) == 0 {
		return 0
	}
	for {
		n++
		if vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 1)))) == 0 {
			p = &(*(*[3][11]uint8)(unsafe.Add(unsafe.Pointer(prob), unsafe.Sizeof([3][11]uint8{})*uintptr(kBands[n]))))[0][0]
		} else {
			var (
				v int
				j int
			)
			if vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 2)))) == 0 {
				p = &(*(*[3][11]uint8)(unsafe.Add(unsafe.Pointer(prob), unsafe.Sizeof([3][11]uint8{})*uintptr(kBands[n]))))[1][0]
				v = 1
			} else {
				if vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 3)))) == 0 {
					if vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 4)))) == 0 {
						v = 2
					} else {
						v = vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 5)))) + 3
					}
				} else {
					if vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 6)))) == 0 {
						if vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 7)))) == 0 {
							v = vp8dx_decode_bool(br, 159) + 5
						} else {
							v = vp8dx_decode_bool(br, 165)*2 + 7
							v += vp8dx_decode_bool(br, 145)
						}
					} else {
						var (
							tab  *uint8
							bit1 int = vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 8))))
							bit0 int = vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), bit1+9))))
							cat  int = bit1*2 + bit0
						)
						v = 0
						for tab = kCat3456[cat]; int(*tab) != 0; tab = (*uint8)(unsafe.Add(unsafe.Pointer(tab), 1)) {
							v += v + vp8dx_decode_bool(br, int(*tab))
						}
						v += (8 << cat) + 3
					}
				}
				p = &(*(*[3][11]uint8)(unsafe.Add(unsafe.Pointer(prob), unsafe.Sizeof([3][11]uint8{})*uintptr(kBands[n]))))[2][0]
			}
			j = int(kZigzag[n-1])
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(j))) = int16(GetSigned(br, v))
			if n == 16 || vp8dx_decode_bool(br, int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 0)))) == 0 {
				return n
			}
		}
		if n == 16 {
			return 16
		}
	}
}
func vp8_decode_mb_tokens(dx *VP8D_COMP, x *MacroBlockd) int {
	var (
		bc         *BOOL_DECODER  = (*BOOL_DECODER)(x.Current_bc)
		fc         *FRAME_CONTEXT = &dx.Common.Fc
		eobs       *byte          = &x.Eobs[0]
		i          int
		nonzeros   int
		eobtotal   int = 0
		qcoeff_ptr *int16
		coef_probs ProbaArray
		a_ctx      *byte = ((*byte)(unsafe.Pointer(x.Above_context)))
		l_ctx      *byte = ((*byte)(unsafe.Pointer(x.Left_context)))
		a          *byte
		l          *byte
		skip_dc    int = 0
	)
	qcoeff_ptr = &x.Qcoeff[0]
	if int(x.Mode_info_context.Mbmi.Is_4x4) == 0 {
		a = (*byte)(unsafe.Add(unsafe.Pointer(a_ctx), 8))
		l = (*byte)(unsafe.Add(unsafe.Pointer(l_ctx), 8))
		coef_probs = ProbaArray(unsafe.Pointer(&fc.Coef_probs[1][0][0][0]))
		nonzeros = GetCoeffs(bc, coef_probs, int(*a+*l), 0, (*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*(24*16))))
		*a = func() byte {
			p := l
			*l = byte(int8(libc.BoolToInt(nonzeros > 0)))
			return *p
		}()
		*(*byte)(unsafe.Add(unsafe.Pointer(eobs), 24)) = byte(int8(nonzeros))
		eobtotal += nonzeros - 16
		coef_probs = ProbaArray(unsafe.Pointer(&fc.Coef_probs[0][0][0][0]))
		skip_dc = 1
	} else {
		coef_probs = ProbaArray(unsafe.Pointer(&fc.Coef_probs[3][0][0][0]))
		skip_dc = 0
	}
	for i = 0; i < 16; i++ {
		a = (*byte)(unsafe.Add(unsafe.Pointer(a_ctx), i&3))
		l = (*byte)(unsafe.Add(unsafe.Pointer(l_ctx), (i&12)>>2))
		nonzeros = GetCoeffs(bc, coef_probs, int(*a+*l), skip_dc, qcoeff_ptr)
		*a = func() byte {
			p := l
			*l = byte(int8(libc.BoolToInt(nonzeros > 0)))
			return *p
		}()
		nonzeros += skip_dc
		*(*byte)(unsafe.Add(unsafe.Pointer(eobs), i)) = byte(int8(nonzeros))
		eobtotal += nonzeros
		qcoeff_ptr = (*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*16))
	}
	coef_probs = ProbaArray(unsafe.Pointer(&fc.Coef_probs[2][0][0][0]))
	a_ctx = (*byte)(unsafe.Add(unsafe.Pointer(a_ctx), 4))
	l_ctx = (*byte)(unsafe.Add(unsafe.Pointer(l_ctx), 4))
	for i = 16; i < 24; i++ {
		a = (*byte)(unsafe.Add(unsafe.Pointer((*byte)(unsafe.Add(unsafe.Pointer(a_ctx), int(libc.BoolToInt(i > 19))<<1))), i&1))
		l = (*byte)(unsafe.Add(unsafe.Pointer(l_ctx), int(libc.BoolToInt(i > 19)<<1+libc.BoolToInt((i&3) > 1))))
		nonzeros = GetCoeffs(bc, coef_probs, int(*a+*l), 0, qcoeff_ptr)
		*a = func() byte {
			p := l
			*l = byte(int8(libc.BoolToInt(nonzeros > 0)))
			return *p
		}()
		*(*byte)(unsafe.Add(unsafe.Pointer(eobs), i)) = byte(int8(nonzeros))
		eobtotal += nonzeros
		qcoeff_ptr = (*int16)(unsafe.Add(unsafe.Pointer(qcoeff_ptr), unsafe.Sizeof(int16(0))*16))
	}
	return eobtotal
}
