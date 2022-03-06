package vp8

func vp8_setup_block_dptrs(x *MacroBlockd) {
	var (
		r int
		c int
	)
	for r = 0; r < 4; r++ {
		for c = 0; c < 4; c++ {
			x.Block[r*4+c].Predictor = &x.Predictor[r*4*16+c*4]
		}
	}
	for r = 0; r < 2; r++ {
		for c = 0; c < 2; c++ {
			x.Block[r*2+16+c].Predictor = &x.Predictor[r*4*8+256+c*4]
		}
	}
	for r = 0; r < 2; r++ {
		for c = 0; c < 2; c++ {
			x.Block[r*2+20+c].Predictor = &x.Predictor[r*4*8+320+c*4]
		}
	}
	for r = 0; r < 25; r++ {
		x.Block[r].Qcoeff = &x.Qcoeff[r*16]
		x.Block[r].Dqcoeff = &x.Dqcoeff[r*16]
		x.Block[r].Eob = &x.Eobs[r]
	}
}
func vp8_build_block_doffsets(x *MacroBlockd) {
	var block int
	for block = 0; block < 16; block++ {
		x.Block[block].Offset = (block>>2)*4*x.Dst.Y_stride + (block&3)*4
	}
	for block = 16; block < 20; block++ {
		x.Block[block+4].Offset = func() int {
			p := &x.Block[block].Offset
			x.Block[block].Offset = ((block-16)>>1)*4*x.Dst.Uv_stride + (block&1)*4
			return *p
		}()
	}
}
