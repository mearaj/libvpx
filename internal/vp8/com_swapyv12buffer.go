package vp8

func vp8_swap_yv12_buffer(new_frame *scale.Yv12BufferConfig, last_frame *scale.Yv12BufferConfig) {
	var temp *uint8
	temp = last_frame.Buffer_alloc
	last_frame.Buffer_alloc = new_frame.Buffer_alloc
	new_frame.Buffer_alloc = temp
	temp = last_frame.Y_buffer
	last_frame.Y_buffer = new_frame.Y_buffer
	new_frame.Y_buffer = temp
	temp = last_frame.U_buffer
	last_frame.U_buffer = new_frame.U_buffer
	new_frame.U_buffer = temp
	temp = last_frame.V_buffer
	last_frame.V_buffer = new_frame.V_buffer
	new_frame.V_buffer = temp
}
