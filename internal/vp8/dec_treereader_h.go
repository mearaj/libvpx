package vp8

import (
	"github.com/mearaj/libvpx/internal/vpx"
	"unsafe"
)

type BOOL_DECODER struct {
	User_buffer_end *uint8
	User_buffer     *uint8
	Value           uint64
	Count           int
	Range           uint
	Decrypt_cb      vpx.DecryptCb
	Decrypt_state   unsafe.Pointer
}
type vp8_reader BOOL_DECODER

func vp8_treed_read(r *vp8_reader, t []int8, p *uint8) int {
	var i int8 = 0
	for int(func() int8 {
		i = t[int(i)+vp8dx_decode_bool((*BOOL_DECODER)(unsafe.Pointer(r)), int(*(*uint8)(unsafe.Add(unsafe.Pointer(p), int(i)>>1))))]
		return i
	}()) > 0 {
	}
	return int(-i)
}
