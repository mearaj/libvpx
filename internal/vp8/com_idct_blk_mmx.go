package vp8

func vp8_dequantize_b_mmx(d *Blockd, DQC *int16) {
	var (
		sq *int16 = d.Qcoeff
		dq *int16 = d.Dqcoeff
	)
	vp8_dequantize_b_impl_mmx(sq, dq, DQC)
}
