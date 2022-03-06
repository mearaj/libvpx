package vp8

const VP8_HEADER_SIZE = 3

type VP8_HEADER struct {
	Type                            uint
	Version                         uint
	Show_frame                      uint
	First_partition_length_in_bytes uint
}
