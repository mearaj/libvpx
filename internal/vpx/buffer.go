package vpx

import "unsafe"

const VPX_MAXIMUM_WORK_BUFFERS = 8
const VP9_MAXIMUM_REF_BUFFERS = 8

type CodecFrameBuffer struct {
	Data *uint8
	Size uint64
	Priv unsafe.Pointer
}
type GetFrameBufferCbFn func(priv unsafe.Pointer, minSize uint64, fb *CodecFrameBuffer) int
type ReleaseFrameBufferCbFn func(priv unsafe.Pointer, fb *CodecFrameBuffer) int
