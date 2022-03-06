package dsp

type filter8_1dfunction func(src_ptr *uint8, src_pitch int64, output_ptr *uint8, out_pitch int64, output_height uint32, filter *int16)
