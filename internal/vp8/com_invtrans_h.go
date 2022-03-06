package vp8

import "unsafe"

func eob_adjust(eobs *byte, diff *int16) {
	var js int
	for js = 0; js < 16; js++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(eobs), js)) == 0 && int(*(*int16)(unsafe.Add(unsafe.Pointer(diff), unsafe.Sizeof(int16(0))*0))) != 0 {
			*(*byte)(unsafe.Add(unsafe.Pointer(eobs), js))++
		}
		diff = (*int16)(unsafe.Add(unsafe.Pointer(diff), unsafe.Sizeof(int16(0))*16))
	}
}
func vp8_inverse_transform_mby(xd *MacroBlockd) {
	var DQC *int16 = &xd.Dequant_y1[0]
	if int(xd.Mode_info_context.Mbmi.Mode) != int(SPLITMV) {
		if xd.Eobs[24] > 1 {
			Vp8ShortInvWalsh4x4C((*int16)(unsafe.Add(unsafe.Pointer(xd.Block[24].Dqcoeff), unsafe.Sizeof(int16(0))*0)), &xd.Qcoeff[0])
		} else {
			vp8_short_inv_walsh4x4_1_c((*int16)(unsafe.Add(unsafe.Pointer(xd.Block[24].Dqcoeff), unsafe.Sizeof(int16(0))*0)), &xd.Qcoeff[0])
		}
		eob_adjust(&xd.Eobs[0], &xd.Qcoeff[0])
		DQC = &xd.Dequant_y1_dc[0]
	}
	vp8_dequant_idct_add_y_block_sse2(&xd.Qcoeff[0], DQC, (*uint8)(unsafe.Pointer(xd.Dst.Y_buffer)), xd.Dst.Y_stride, &xd.Eobs[0])
}
