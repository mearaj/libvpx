package vp8

import "github.com/mearaj/libvpx/internal/ports"

func Vp8Rtcd() {
	ports.Once(SetupRtcdInternal)
}
