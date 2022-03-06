package scale

import "github.com/mearaj/libvpx/internal/ports"

func ScaleRtcd() {
	ports.Once(setup_rtcd_internal)
}
