package internal

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/mearaj/libvpx/internal/ports"
	"unsafe"
)

var IVF_SIGNATURE *byte = libc.CString("DKIF")

func fix_framerate(num *int, den *int) {
	if *den > 0 && *den < 1000000000 && *num > 0 && *num < 1000 {
		if *num&1 != 0 {
			*den *= 2
		} else {
			*num /= 2
		}
	} else {
		*num = 30
		*den = 1
	}
}
func FileIsIvf(input_ctx *VpxInputContext) int {
	var (
		raw_hdr [32]byte
		is_ivf  int = 0
	)
	if int(input_ctx.File.ReadN(&raw_hdr[0], 1, 32)) == 32 {
		if libc.MemCmp(unsafe.Pointer(IVF_SIGNATURE), unsafe.Pointer(&raw_hdr[0]), 4) == 0 {
			is_ivf = 1
			if ports.GetLe16AsInt(unsafe.Pointer(&raw_hdr[4])) != 0 {
				stdio.Fprintf(stdio.Stderr(), "Error: Unrecognized IVF version! This file may not decode properly.")
			}
			input_ctx.Fourcc = uint32(ports.GetLe32AsInt(unsafe.Pointer(&raw_hdr[8])))
			input_ctx.Width = uint32(ports.GetLe16AsInt(unsafe.Pointer(&raw_hdr[12])))
			input_ctx.Height = uint32(ports.GetLe16AsInt(unsafe.Pointer(&raw_hdr[14])))
			input_ctx.Framerate.Numerator = int(ports.GetLe32AsInt(unsafe.Pointer(&raw_hdr[16])))
			input_ctx.Framerate.Denominator = int(ports.GetLe32AsInt(unsafe.Pointer(&raw_hdr[20])))
			fix_framerate(&input_ctx.Framerate.Numerator, &input_ctx.Framerate.Denominator)
		}
	}
	if is_ivf == 0 {
		Rewind(input_ctx.File)
		input_ctx.Detect.Buf_read = 0
	} else {
		input_ctx.Detect.Position = 4
	}
	return is_ivf
}
func IvfReadFrame(infile *stdio.File, buffer **uint8, bytes_read *uint64, buffer_size *uint64) int {
	var (
		raw_header [12]byte = [12]byte{}
		frame_size uint64   = 0
	)
	if int(infile.ReadN(&raw_header[0], 4+8, 1)) != 1 {
		if int(infile.IsEOF()) == 0 {
			Warn(libc.CString("Failed to read frame size"))
		}
	} else {
		frame_size = uint64(ports.GetLe32AsInt(unsafe.Pointer(&raw_header[0])))
		if frame_size > 256*1024*1024 {
			Warn(libc.CString("Read invalid frame size (%u)"), uint(frame_size))
			frame_size = 0
		}
		if frame_size > *buffer_size {
			var new_buffer *uint8 = (*uint8)(libc.Realloc(unsafe.Pointer(*buffer), int(frame_size*2)))
			if new_buffer != nil {
				*buffer = new_buffer
				*buffer_size = frame_size * 2
			} else {
				Warn(libc.CString("Failed to allocate compressed data buffer"))
				frame_size = 0
			}
		}
	}
	if int(infile.IsEOF()) == 0 {
		if uint64(infile.ReadN((*byte)(unsafe.Pointer(*buffer)), 1, int(frame_size))) != frame_size {
			Warn(libc.CString("Failed to read full frame"))
			return 1
		}
		*bytes_read = frame_size
		return 0
	}
	return 1
}
