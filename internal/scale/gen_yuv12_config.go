package scale

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/mearaj/libvpx/internal/mem"
	"github.com/mearaj/libvpx/internal/vpx"
	"math"
	"unsafe"
)

func Vp8Yv12DeAllocFrameBuffer(ybf *Yv12BufferConfig) int {
	if ybf != nil {
		if ybf.Buffer_alloc_sz > 0 {
			mem.VpxFree(unsafe.Pointer(ybf.Buffer_alloc))
		}
		*ybf = Yv12BufferConfig{}
	} else {
		return -1
	}
	return 0
}
func Vp8Yv12ReAllocFrameBuffer(ybf *Yv12BufferConfig, width int, height int, border int) int {
	if ybf != nil {
		var (
			aligned_width  int    = (width + 15) & ^int(15)
			aligned_height int    = (height + 15) & ^int(15)
			y_stride       int    = ((aligned_width + border*2) + 31) & ^int(31)
			yplane_size    int    = (aligned_height + border*2) * y_stride
			uv_width       int    = aligned_width >> 1
			uv_height      int    = aligned_height >> 1
			uv_stride      int    = y_stride >> 1
			uvplane_size   int    = (uv_height + border) * uv_stride
			frame_size     uint64 = uint64(yplane_size + uvplane_size*2)
		)
		if ybf.Buffer_alloc == nil {
			ybf.Buffer_alloc = (*uint8)(mem.VpxMemAlign(32, frame_size))
			if ybf.Buffer_alloc == nil {
				ybf.Buffer_alloc_sz = 0
				return -1
			}
			ybf.Buffer_alloc_sz = frame_size
		}
		if ybf.Buffer_alloc_sz < frame_size {
			return -1
		}
		if border&31 != 0 {
			return -3
		}
		ybf.Y_crop_width = width
		ybf.Y_crop_height = height
		ybf.Y_width = aligned_width
		ybf.Y_height = aligned_height
		ybf.Y_stride = y_stride
		ybf.Uv_crop_width = (width + 1) / 2
		ybf.Uv_crop_height = (height + 1) / 2
		ybf.Uv_width = uv_width
		ybf.Uv_height = uv_height
		ybf.Uv_stride = uv_stride
		ybf.Alpha_width = 0
		ybf.Alpha_height = 0
		ybf.Alpha_stride = 0
		ybf.Border = border
		ybf.Frame_size = frame_size
		ybf.Y_buffer = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ybf.Buffer_alloc), border*y_stride))), border))
		ybf.U_buffer = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ybf.Buffer_alloc), yplane_size))), border/2*uv_stride))), border/2))
		ybf.V_buffer = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(ybf.Buffer_alloc), yplane_size))), uvplane_size))), border/2*uv_stride))), border/2))
		ybf.Alpha_buffer = nil
		ybf.Corrupted = 0
		return 0
	}
	return -2
}
func Vp8Yv12AllocFrameBuffer(ybf *Yv12BufferConfig, width int, height int, border int) int {
	if ybf != nil {
		Vp8Yv12DeAllocFrameBuffer(ybf)
		return Vp8Yv12ReAllocFrameBuffer(ybf, width, height, border)
	}
	return -2
}
func VpxFreeFrameBuffer(ybf *Yv12BufferConfig) int {
	if ybf != nil {
		if ybf.Buffer_alloc_sz > 0 {
			mem.VpxFree(unsafe.Pointer(ybf.Buffer_alloc))
		}
		*ybf = Yv12BufferConfig{}
	} else {
		return -1
	}
	return 0
}
func VpxReAllocFrameBuffer(ybf *Yv12BufferConfig, width int, height int, ss_x int, ss_y int, border int, byte_alignment int, fb *vpx.CodecFrameBuffer, cb vpx.GetFrameBufferCbFn, cb_priv unsafe.Pointer) int {
	if border&31 != 0 {
		return -3
	}
	if ybf != nil {
		var vp9_byte_align int
		if byte_alignment == 0 {
			vp9_byte_align = 1
		} else {
			vp9_byte_align = byte_alignment
		}
		var aligned_width int = (width + 7) & ^int(7)
		var aligned_height int = (height + 7) & ^int(7)
		var y_stride int = ((aligned_width + border*2) + 31) & ^int(31)
		var yplane_size uint64 = uint64((aligned_height+border*2)*int(uint64(y_stride)) + byte_alignment)
		var uv_width int = aligned_width >> ss_x
		var uv_height int = aligned_height >> ss_y
		var uv_stride int = y_stride >> ss_x
		var uv_border_w int = border >> ss_x
		var uv_border_h int = border >> ss_y
		var uvplane_size uint64 = uint64((uv_height+uv_border_h*2)*int(uint64(uv_stride)) + byte_alignment)
		var frame_size uint64 = yplane_size + uvplane_size*2
		var buf *uint8 = nil
		if frame_size > math.MaxUint64 {
			return -1
		}
		if cb != nil {
			var (
				align_addr_extra_size int    = 31
				external_frame_size   uint64 = frame_size + uint64(align_addr_extra_size)
			)
			libc.Assert(fb != nil)
			if external_frame_size != external_frame_size {
				return -1
			}
			if cb(cb_priv, external_frame_size, fb) < 0 {
				return -1
			}
			if fb.Data == nil || fb.Size < external_frame_size {
				return -1
			}
			ybf.Buffer_alloc = (*uint8)(unsafe.Pointer(uintptr((uint64(uintptr(unsafe.Pointer(fb.Data))) + (32 - 1)) & 0xFFFFFFFFFFFFFFE0)))
		} else if frame_size > ybf.Buffer_alloc_sz {
			mem.VpxFree(unsafe.Pointer(ybf.Buffer_alloc))
			ybf.Buffer_alloc = nil
			ybf.Buffer_alloc_sz = 0
			ybf.Buffer_alloc = (*uint8)(mem.VpxMemAlign(32, frame_size))
			if ybf.Buffer_alloc == nil {
				return -1
			}
			ybf.Buffer_alloc_sz = frame_size
			libc.MemSet(unsafe.Pointer(ybf.Buffer_alloc), 0, int(ybf.Buffer_alloc_sz))
		}
		ybf.Y_crop_width = width
		ybf.Y_crop_height = height
		ybf.Y_width = aligned_width
		ybf.Y_height = aligned_height
		ybf.Y_stride = y_stride
		ybf.Uv_crop_width = (width + ss_x) >> ss_x
		ybf.Uv_crop_height = (height + ss_y) >> ss_y
		ybf.Uv_width = uv_width
		ybf.Uv_height = uv_height
		ybf.Uv_stride = uv_stride
		ybf.Border = border
		ybf.Frame_size = frame_size
		ybf.Subsampling_x = ss_x
		ybf.Subsampling_y = ss_y
		buf = ybf.Buffer_alloc
		ybf.Y_buffer = (*uint8)(unsafe.Pointer(uintptr((uint64(uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(buf), border*y_stride))), border))))) + uint64(vp9_byte_align-1)) & uint64(-vp9_byte_align))))
		ybf.U_buffer = (*uint8)(unsafe.Pointer(uintptr((uint64(uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(buf), yplane_size))), uv_border_h*uv_stride))), uv_border_w))))) + uint64(vp9_byte_align-1)) & uint64(-vp9_byte_align))))
		ybf.V_buffer = (*uint8)(unsafe.Pointer(uintptr((uint64(uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(buf), yplane_size))), uvplane_size))), uv_border_h*uv_stride))), uv_border_w))))) + uint64(vp9_byte_align-1)) & uint64(-vp9_byte_align))))
		ybf.Corrupted = 0
		return 0
	}
	return -2
}
func VpxAllocFrameBuffer(ybf *Yv12BufferConfig, width int, height int, ss_x int, ss_y int, border int, byte_alignment int) int {
	if ybf != nil {
		VpxFreeFrameBuffer(ybf)
		return VpxReAllocFrameBuffer(ybf, width, height, ss_x, ss_y, border, byte_alignment, nil, nil, nil)
	}
	return -2
}
