package dsp

import "github.com/mearaj/libvpx/internal/ports"

func DspRtcd() {
	ports.Once(setup_rtcd_internal)
}
