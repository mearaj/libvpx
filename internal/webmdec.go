package internal

import "unsafe"

type WebmInputContext struct {
	Reader            unsafe.Pointer
	Segment           unsafe.Pointer
	Buffer            *uint8
	Cluster           unsafe.Pointer
	Block_entry       unsafe.Pointer
	Block             unsafe.Pointer
	Block_frame_index int
	Video_track_index int
	Timestamp_ns      uint64
	Is_key_frame      int
	Reached_eos       int
}
