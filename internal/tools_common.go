package internal

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/mearaj/libvpx/internal/vp8"
	"github.com/mearaj/libvpx/internal/vpx"
	"math"
	"os"
	"unsafe"
)

const PATH_MAX = 512
const IVF_FRAME_HDR_SZ = 12
const IVF_FILE_HDR_SZ = 32
const VP8_FOURCC = 808996950
const VP9_FOURCC = 0x30395056

type FileOffset uint64
type VideoFileType int

const (
	FILE_TYPE_RAW = VideoFileType(iota)
	FILE_TYPE_IVF
	FILE_TYPE_Y4M
	FILE_TYPE_WEBM
)

type FileTypeDetectionBuffer struct {
	Buf      [4]byte
	Buf_read uint64
	Position uint64
}
type VpxRational struct {
	Numerator   int
	Denominator int
}
type VpxInputContext struct {
	Filename           *byte
	File               *stdio.File
	Length             int64
	Detect             FileTypeDetectionBuffer
	File_type          VideoFileType
	Width              uint32
	Height             uint32
	Pixel_aspect_ratio VpxRational
	Fmt                vpx.ImgFmt
	Bit_depth          vpx.BitDepth
	Only_i420          int
	Fourcc             uint32
	Framerate          VpxRational
	Y4m                Y4mInput
}
type VpxInterface struct {
	Name            *byte
	Fourcc          uint32
	Codec_interface func() *vpx.CodecIFace
}

func WrapFRead(ptr unsafe.Pointer, size uint64, nmemb uint64, stream *stdio.File) uint64 {
	return uint64(stream.ReadN((*byte)(ptr), int(size), int(nmemb)))
}
func SetBinaryMode(stream *stdio.File) *stdio.File {
	_ = stream
	return stream
}
func Die(fmt *byte, _rest ...interface{}) {
	for {
		{
			var (
				l  *byte = nil
				ap libc.ArgList
			)
			ap.Start(fmt, _rest)
			if l != nil {
				stdio.Fprintf(stdio.Stderr(), "%s: ", l)
			}
			stdio.Vfprintf(stdio.Stderr(), libc.GoString(fmt), ap)
			stdio.Fprintf(stdio.Stderr(), "\n")
			ap.End()
		}
		if true {
			break
		}
	}
	UsageExit()
}
func Fatal(fmt *byte, _rest ...interface{}) {
	for {
		{
			var (
				l  *byte = libc.CString("Fatal")
				ap libc.ArgList
			)
			ap.Start(fmt, _rest)
			if l != nil {
				stdio.Fprintf(stdio.Stderr(), "%s: ", l)
			}
			stdio.Vfprintf(stdio.Stderr(), libc.GoString(fmt), ap)
			stdio.Fprintf(stdio.Stderr(), "\n")
			ap.End()
		}
		if true {
			break
		}
	}
	os.Exit(0)
}
func Warn(fmt *byte, _rest ...interface{}) {
	for {
		{
			var (
				l  *byte = libc.CString("Warning")
				ap libc.ArgList
			)
			ap.Start(fmt, _rest)
			if l != nil {
				stdio.Fprintf(stdio.Stderr(), "%s: ", l)
			}
			stdio.Vfprintf(stdio.Stderr(), libc.GoString(fmt), ap)
			stdio.Fprintf(stdio.Stderr(), "\n")
			ap.End()
		}
		if true {
			break
		}
	}
}
func DieCodec(ctx *vpx.CodecCtx, s *byte) {
	var detail *byte = vpx.CodecErrorDetail(ctx)
	stdio.Printf("%s: %s\n", s, vpx.CodecError(ctx))
	if detail != nil {
		stdio.Printf("    %s\n", detail)
	}
	os.Exit(0)
}
func ReadYuvFrame(input_ctx *VpxInputContext, yuv_frame *vpx.Image) int {
	var (
		f         *stdio.File              = input_ctx.File
		detect    *FileTypeDetectionBuffer = &input_ctx.Detect
		plane     int                      = 0
		shortread int                      = 0
		bytespp   int
	)
	if (yuv_frame.Fmt & VPX_IMG_FMT_HIGHBITDEPTH) != 0 {
		bytespp = 2
	} else {
		bytespp = 1
	}
	for plane = 0; plane < 3; plane++ {
		var (
			ptr *uint8
			w   int = ImgPlaneWidth(yuv_frame, plane)
			h   int = ImgPlaneHeight(yuv_frame, plane)
			r   int
		)
		if yuv_frame.Fmt == vpx.ImgFmt(VPX_IMG_FMT_NV12) && plane > 1 {
			break
		}
		if yuv_frame.Fmt == vpx.ImgFmt(VPX_IMG_FMT_NV12) && plane == 1 {
			w = (w + 1) & ^int(1)
		}
		switch plane {
		case 1:
			ptr = yuv_frame.Planes[func() int {
				if yuv_frame.Fmt == vpx.ImgFmt(VPX_IMG_FMT_YV12) {
					return VPX_PLANE_V
				}
				return VPX_PLANE_U
			}()]
		case 2:
			ptr = yuv_frame.Planes[func() int {
				if yuv_frame.Fmt == vpx.ImgFmt(VPX_IMG_FMT_YV12) {
					return VPX_PLANE_U
				}
				return VPX_PLANE_V
			}()]
		default:
			ptr = yuv_frame.Planes[plane]
		}
		for r = 0; r < h; r++ {
			var (
				needed       uint64 = uint64(w * bytespp)
				buf_position uint64 = 0
				left         uint64 = detect.Buf_read - detect.Position
			)
			if left > 0 {
				var more uint64
				if left < needed {
					more = left
				} else {
					more = needed
				}
				libc.MemCpy(unsafe.Pointer(ptr), unsafe.Pointer(&detect.Buf[detect.Position]), int(more))
				buf_position = more
				needed -= more
				detect.Position += more
			}
			if needed > 0 {
				shortread |= int(libc.BoolToInt(uint64(f.ReadN((*byte)(unsafe.Add(unsafe.Pointer(ptr), buf_position)), 1, int(needed))) < needed))
			}
			ptr = (*uint8)(unsafe.Add(unsafe.Pointer(ptr), yuv_frame.Stride[plane]))
		}
	}
	return shortread
}

var VpxEncoders [2]VpxInterface = [2]VpxInterface{{Name: libc.CString("vp8"), Fourcc: VP8_FOURCC, Codec_interface: vpx_codec_vp8_cx}, {Name: libc.CString("vp9"), Fourcc: VP9_FOURCC, Codec_interface: vpx_codec_vp9_cx}}

func GetVpxEncoderCount() int {
	return int(unsafe.Sizeof([2]VpxInterface{}) / unsafe.Sizeof(VpxInterface{}))
}
func GetVpxEncoderByIndex(i int) *VpxInterface {
	return &VpxEncoders[i]
}
func GetVpxEncoderByName(name *byte) *VpxInterface {
	var i int
	for i = 0; i < GetVpxEncoderCount(); i++ {
		var encoder *VpxInterface = GetVpxEncoderByIndex(i)
		if libc.StrCmp(encoder.Name, name) == 0 {
			return encoder
		}
	}
	return nil
}

var VpxDecoders [2]VpxInterface = [2]VpxInterface{{Name: libc.CString("vp8"), Fourcc: VP8_FOURCC, Codec_interface: vp8.CodecDxFn}, {Name: libc.CString("vp9"), Fourcc: VP9_FOURCC, Codec_interface: vpx_codec_vp9_dx}}

func GetVpxDecoderCount() int {
	return int(unsafe.Sizeof([2]VpxInterface{}) / unsafe.Sizeof(VpxInterface{}))
}
func GetVpxDecoderByIndex(i int) *VpxInterface {
	return &VpxDecoders[i]
}
func GetVpxDecoderByName(name *byte) *VpxInterface {
	var i int
	for i = 0; i < GetVpxDecoderCount(); i++ {
		var decoder *VpxInterface = GetVpxDecoderByIndex(i)
		if libc.StrCmp(decoder.Name, name) == 0 {
			return decoder
		}
	}
	return nil
}
func GetVpxDecoderByFourcc(fourcc uint32) *VpxInterface {
	var i int
	for i = 0; i < GetVpxDecoderCount(); i++ {
		var decoder *VpxInterface = GetVpxDecoderByIndex(i)
		if int(decoder.Fourcc) == int(fourcc) {
			return decoder
		}
	}
	return nil
}
func ImgPlaneWidth(img *vpx.Image, plane int) int {
	if plane > 0 && img.X_chroma_shift > 0 {
		return int((img.D_w + 1) >> img.X_chroma_shift)
	} else {
		return int(img.D_w)
	}
}
func ImgPlaneHeight(img *vpx.Image, plane int) int {
	if plane > 0 && img.Y_chroma_shift > 0 {
		return int((img.D_h + 1) >> img.Y_chroma_shift)
	} else {
		return int(img.D_h)
	}
}
func VpxImgWrite(img *vpx.Image, file *stdio.File) {
	var plane int
	for plane = 0; plane < 3; plane++ {
		var (
			buf    *uint8 = img.Planes[plane]
			stride int    = img.Stride[plane]
			w      int    = ImgPlaneWidth(img, plane) * (func() int {
				if (img.Fmt & VPX_IMG_FMT_HIGHBITDEPTH) != 0 {
					return 2
				}
				return 1
			}())
			h int = ImgPlaneHeight(img, plane)
			y int
		)
		for y = 0; y < h; y++ {
			file.WriteN((*byte)(unsafe.Pointer(buf)), 1, w)
			buf = (*uint8)(unsafe.Add(unsafe.Pointer(buf), stride))
		}
	}
}
func VpxImgRead(img *vpx.Image, file *stdio.File) int {
	var plane int
	for plane = 0; plane < 3; plane++ {
		var (
			buf    *uint8 = img.Planes[plane]
			stride int    = img.Stride[plane]
			w      int    = ImgPlaneWidth(img, plane) * (func() int {
				if (img.Fmt & VPX_IMG_FMT_HIGHBITDEPTH) != 0 {
					return 2
				}
				return 1
			}())
			h int = ImgPlaneHeight(img, plane)
			y int
		)
		for y = 0; y < h; y++ {
			if uint64(file.ReadN((*byte)(unsafe.Pointer(buf)), 1, w)) != uint64(w) {
				return 0
			}
			buf = (*uint8)(unsafe.Add(unsafe.Pointer(buf), stride))
		}
	}
	return 1
}
func SseToPsnr(samples float64, peak float64, sse float64) float64 {
	var kMaxPSNR float64 = 100.0
	if sse > 0.0 {
		var psnr float64 = math.Log10(samples*peak*peak/sse) * 10.0
		if psnr > kMaxPSNR {
			return kMaxPSNR
		}
		return psnr
	} else {
		return kMaxPSNR
	}
}
func ReadFrame(input_ctx *VpxInputContext, img *vpx.Image) int {
	var (
		f         *stdio.File = input_ctx.File
		y4m       *Y4mInput   = &input_ctx.Y4m
		shortread int         = 0
	)
	if input_ctx.File_type == VideoFileType(FILE_TYPE_Y4M) {
		if y4m_input_fetch_frame(y4m, f, img) < 1 {
			return 0
		}
	} else {
		shortread = ReadYuvFrame(input_ctx, img)
	}
	return int(libc.BoolToInt(shortread == 0))
}
func FileIsY4m(detect [4]byte) int {
	if libc.MemCmp(unsafe.Pointer(&detect[0]), unsafe.Pointer(libc.CString("YUV4")), 4) == 0 {
		return 1
	}
	return 0
}
func FourccIsIvf(detect [4]byte) int {
	if libc.MemCmp(unsafe.Pointer(&detect[0]), unsafe.Pointer(libc.CString("DKIF")), 4) == 0 {
		return 1
	}
	return 0
}
func OpenInputFile(input *VpxInputContext) {
	if libc.StrCmp(input.Filename, libc.CString("-")) != 0 {
		input.File = stdio.FOpen(libc.GoString(input.Filename), "rb")
	} else {
		input.File = SetBinaryMode(stdio.Stdin())
	}
	if input.File == nil {
		Fatal(libc.CString("Failed to open input file"))
	}
	if fseeko(input.File, 0, int64(stdio.SEEK_END)) == 0 {
		input.Length = int64(ftello(input.File))
		Rewind(input.File)
	}
	input.Pixel_aspect_ratio.Numerator = 1
	input.Pixel_aspect_ratio.Denominator = 1
	input.Detect.Buf_read = uint64(input.File.ReadN(&input.Detect.Buf[0], 1, 4))
	input.Detect.Position = 0
	if input.Detect.Buf_read == 4 && FileIsY4m(input.Detect.Buf) != 0 {
		if y4m_input_open(&input.Y4m, input.File, &input.Detect.Buf[0], 4, input.Only_i420) >= 0 {
			input.File_type = VideoFileType(FILE_TYPE_Y4M)
			input.Width = uint32(int32(input.Y4m.Pic_w))
			input.Height = uint32(int32(input.Y4m.Pic_h))
			input.Pixel_aspect_ratio.Numerator = input.Y4m.Par_n
			input.Pixel_aspect_ratio.Denominator = input.Y4m.Par_d
			input.Framerate.Numerator = input.Y4m.Fps_n
			input.Framerate.Denominator = input.Y4m.Fps_d
			input.Fmt = vpx.ImgFmt(input.Y4m.Vpx_fmt)
			input.Bit_depth = vpx.BitDepth(input.Y4m.Bit_depth)
		} else {
			Fatal(libc.CString("Unsupported Y4M stream."))
		}
	} else if input.Detect.Buf_read == 4 && FourccIsIvf(input.Detect.Buf) != 0 {
		Fatal(libc.CString("IVF is not supported as input."))
	} else {
		input.File_type = VideoFileType(FILE_TYPE_RAW)
	}
}
func CloseInputFile(input *VpxInputContext) {
	input.File.Close()
	if input.File_type == VideoFileType(FILE_TYPE_Y4M) {
		y4m_input_close(&input.Y4m)
	}
}
func CompareImg(img1 *vpx.Image, img2 *vpx.Image) int {
	var (
		l_w   uint32 = uint32(img1.D_w)
		c_w   uint32 = uint32((img1.D_w + img1.X_chroma_shift) >> img1.X_chroma_shift)
		c_h   uint32 = uint32((img1.D_h + img1.Y_chroma_shift) >> img1.Y_chroma_shift)
		i     uint32
		match int = 1
	)
	match &= int(libc.BoolToInt(img1.Fmt == img2.Fmt))
	match &= int(libc.BoolToInt(img1.D_w == img2.D_w))
	match &= int(libc.BoolToInt(img1.D_h == img2.D_h))
	for i = 0; uint(i) < img1.D_h; i++ {
		match &= int(libc.BoolToInt(libc.MemCmp(unsafe.Add(unsafe.Pointer(img1.Planes[VPX_PLANE_Y]), int(i)*img1.Stride[VPX_PLANE_Y]), unsafe.Add(unsafe.Pointer(img2.Planes[VPX_PLANE_Y]), int(i)*img2.Stride[VPX_PLANE_Y]), int(l_w)) == 0))
	}
	for i = 0; int(i) < int(c_h); i++ {
		match &= int(libc.BoolToInt(libc.MemCmp(unsafe.Add(unsafe.Pointer(img1.Planes[VPX_PLANE_U]), int(i)*img1.Stride[VPX_PLANE_U]), unsafe.Add(unsafe.Pointer(img2.Planes[VPX_PLANE_U]), int(i)*img2.Stride[VPX_PLANE_U]), int(c_w)) == 0))
	}
	for i = 0; int(i) < int(c_h); i++ {
		match &= int(libc.BoolToInt(libc.MemCmp(unsafe.Add(unsafe.Pointer(img1.Planes[VPX_PLANE_V]), int(i)*img1.Stride[VPX_PLANE_V]), unsafe.Add(unsafe.Pointer(img2.Planes[VPX_PLANE_V]), int(i)*img2.Stride[VPX_PLANE_V]), int(c_w)) == 0))
	}
	return match
}
func FindMismatch(img1 *vpx.Image, img2 *vpx.Image, yloc [4]int, uloc [4]int, vloc [4]int) {
	var (
		bsize  uint32 = 64
		bsizey uint32 = uint32(uint(bsize) >> img1.Y_chroma_shift)
		bsizex uint32 = uint32(uint(bsize) >> img1.X_chroma_shift)
		c_w    uint32 = uint32((img1.D_w + img1.X_chroma_shift) >> img1.X_chroma_shift)
		c_h    uint32 = uint32((img1.D_h + img1.Y_chroma_shift) >> img1.Y_chroma_shift)
		match  int    = 1
		i      uint32
		j      uint32
	)
	yloc[0] = func() int {
		p := &yloc[1]
		yloc[1] = func() int {
			p := &yloc[2]
			yloc[2] = func() int {
				p := &yloc[3]
				yloc[3] = -1
				return *p
			}()
			return *p
		}()
		return *p
	}()
	for func() int {
		i = 0
		return func() int {
			match = 1
			return match
		}()
	}(); match != 0 && uint(i) < img1.D_h; i += bsize {
		for j = 0; match != 0 && uint(j) < img1.D_w; j += bsize {
			var (
				k  int
				l  int
				si int = (func() int {
					if (int(i) + int(bsize)) < int(img1.D_h) {
						return int(i) + int(bsize)
					}
					return int(img1.D_h)
				}()) - int(i)
				sj int = (func() int {
					if (int(j) + int(bsize)) < int(img1.D_w) {
						return int(j) + int(bsize)
					}
					return int(img1.D_w)
				}()) - int(j)
			)
			for k = 0; match != 0 && k < si; k++ {
				for l = 0; match != 0 && l < sj; l++ {
					if int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img1.Planes[VPX_PLANE_Y]), (int(i)+k)*img1.Stride[VPX_PLANE_Y]))), j))), l)))) != int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img2.Planes[VPX_PLANE_Y]), (int(i)+k)*img2.Stride[VPX_PLANE_Y]))), j))), l)))) {
						yloc[0] = int(i) + k
						yloc[1] = int(j) + l
						yloc[2] = int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img1.Planes[VPX_PLANE_Y]), (int(i)+k)*img1.Stride[VPX_PLANE_Y]))), j))), l))))
						yloc[3] = int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img2.Planes[VPX_PLANE_Y]), (int(i)+k)*img2.Stride[VPX_PLANE_Y]))), j))), l))))
						match = 0
						break
					}
				}
			}
		}
	}
	uloc[0] = func() int {
		p := &uloc[1]
		uloc[1] = func() int {
			p := &uloc[2]
			uloc[2] = func() int {
				p := &uloc[3]
				uloc[3] = -1
				return *p
			}()
			return *p
		}()
		return *p
	}()
	for func() int {
		i = 0
		return func() int {
			match = 1
			return match
		}()
	}(); match != 0 && int(i) < int(c_h); i += bsizey {
		for j = 0; match != 0 && int(j) < int(c_w); j += bsizex {
			var (
				k  int
				l  int
				si int = (func() int {
					if (int(i) + int(bsizey)) < (int(c_h) - int(i)) {
						return int(i) + int(bsizey)
					}
					return int(c_h) - int(i)
				}())
				sj int = (func() int {
					if (int(j) + int(bsizex)) < (int(c_w) - int(j)) {
						return int(j) + int(bsizex)
					}
					return int(c_w) - int(j)
				}())
			)
			for k = 0; match != 0 && k < si; k++ {
				for l = 0; match != 0 && l < sj; l++ {
					if int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img1.Planes[VPX_PLANE_U]), (int(i)+k)*img1.Stride[VPX_PLANE_U]))), j))), l)))) != int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img2.Planes[VPX_PLANE_U]), (int(i)+k)*img2.Stride[VPX_PLANE_U]))), j))), l)))) {
						uloc[0] = int(i) + k
						uloc[1] = int(j) + l
						uloc[2] = int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img1.Planes[VPX_PLANE_U]), (int(i)+k)*img1.Stride[VPX_PLANE_U]))), j))), l))))
						uloc[3] = int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img2.Planes[VPX_PLANE_U]), (int(i)+k)*img2.Stride[VPX_PLANE_U]))), j))), l))))
						match = 0
						break
					}
				}
			}
		}
	}
	vloc[0] = func() int {
		p := &vloc[1]
		vloc[1] = func() int {
			p := &vloc[2]
			vloc[2] = func() int {
				p := &vloc[3]
				vloc[3] = -1
				return *p
			}()
			return *p
		}()
		return *p
	}()
	for func() int {
		i = 0
		return func() int {
			match = 1
			return match
		}()
	}(); match != 0 && int(i) < int(c_h); i += bsizey {
		for j = 0; match != 0 && int(j) < int(c_w); j += bsizex {
			var (
				k  int
				l  int
				si int = (func() int {
					if (int(i) + int(bsizey)) < (int(c_h) - int(i)) {
						return int(i) + int(bsizey)
					}
					return int(c_h) - int(i)
				}())
				sj int = (func() int {
					if (int(j) + int(bsizex)) < (int(c_w) - int(j)) {
						return int(j) + int(bsizex)
					}
					return int(c_w) - int(j)
				}())
			)
			for k = 0; match != 0 && k < si; k++ {
				for l = 0; match != 0 && l < sj; l++ {
					if int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img1.Planes[VPX_PLANE_V]), (int(i)+k)*img1.Stride[VPX_PLANE_V]))), j))), l)))) != int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img2.Planes[VPX_PLANE_V]), (int(i)+k)*img2.Stride[VPX_PLANE_V]))), j))), l)))) {
						vloc[0] = int(i) + k
						vloc[1] = int(j) + l
						vloc[2] = int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img1.Planes[VPX_PLANE_V]), (int(i)+k)*img1.Stride[VPX_PLANE_V]))), j))), l))))
						vloc[3] = int(*((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(img2.Planes[VPX_PLANE_V]), (int(i)+k)*img2.Stride[VPX_PLANE_V]))), j))), l))))
						match = 0
						break
					}
				}
			}
		}
	}
}
