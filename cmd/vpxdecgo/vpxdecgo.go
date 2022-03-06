package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/mearaj/libvpx/internal"
	"github.com/mearaj/libvpx/internal/ports"
	"github.com/mearaj/libvpx/internal/vpx"
	"os"
	"unsafe"
)

var exec_name *byte

type VpxDecInputContext struct {
	Vpx_input_ctx *internal.VpxInputContext
	Webm_ctx      *internal.WebmInputContext
}

var help internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("help"), Has_val: 0, Desc: libc.CString("Show usage options and exit"), Enums: nil}
var looparg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("loops"), Has_val: 1, Desc: libc.CString("Number of times to decode the file"), Enums: nil}
var codecarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("codec"), Has_val: 1, Desc: libc.CString("Codec to use"), Enums: nil}
var use_yv12 internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("yv12"), Has_val: 0, Desc: libc.CString("Output raw YV12 frames"), Enums: nil}
var use_i420 internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("i420"), Has_val: 0, Desc: libc.CString("Output raw I420 frames"), Enums: nil}
var flipuvarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("flipuv"), Has_val: 0, Desc: libc.CString("Flip the chroma planes in the output"), Enums: nil}
var rawvideo internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("rawvideo"), Has_val: 0, Desc: libc.CString("Output raw YUV frames"), Enums: nil}
var noblitarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("noblit"), Has_val: 0, Desc: libc.CString("Don't process the decoded frames"), Enums: nil}
var progressarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("progress"), Has_val: 0, Desc: libc.CString("Show progress after each frame decodes"), Enums: nil}
var limitarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("limit"), Has_val: 1, Desc: libc.CString("Stop decoding after n frames"), Enums: nil}
var skiparg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("skip"), Has_val: 1, Desc: libc.CString("Skip the first n input frames"), Enums: nil}
var postprocarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("postproc"), Has_val: 0, Desc: libc.CString("Postprocess decoded frames"), Enums: nil}
var summaryarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("summary"), Has_val: 0, Desc: libc.CString("Show timing summary"), Enums: nil}
var outputfile internal.ArgDefT = internal.ArgDefT{Short_name: libc.CString("o"), Long_name: libc.CString("output"), Has_val: 1, Desc: libc.CString("Output file name pattern (see below)"), Enums: nil}
var threadsarg internal.ArgDefT = internal.ArgDefT{Short_name: libc.CString("t"), Long_name: libc.CString("threads"), Has_val: 1, Desc: libc.CString("Max threads to use"), Enums: nil}
var frameparallelarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("frame-parallel"), Has_val: 0, Desc: libc.CString("Frame parallel decode (ignored)"), Enums: nil}
var verbosearg internal.ArgDefT = internal.ArgDefT{Short_name: libc.CString("v"), Long_name: libc.CString("verbose"), Has_val: 0, Desc: libc.CString("Show version string"), Enums: nil}
var error_concealment internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("error-concealment"), Has_val: 0, Desc: libc.CString("Enable decoder error-concealment"), Enums: nil}
var scalearg internal.ArgDefT = internal.ArgDefT{Short_name: libc.CString("S"), Long_name: libc.CString("scale"), Has_val: 0, Desc: libc.CString("Scale output frames uniformly"), Enums: nil}
var continuearg internal.ArgDefT = internal.ArgDefT{Short_name: libc.CString("k"), Long_name: libc.CString("keep-going"), Has_val: 0, Desc: libc.CString("(debug) Continue decoding after error"), Enums: nil}
var fb_arg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("frame-buffers"), Has_val: 1, Desc: libc.CString("Number of frame buffers to use"), Enums: nil}
var md5arg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("md5"), Has_val: 0, Desc: libc.CString("Compute the MD5 sum of the decoded frame"), Enums: nil}
var svcdecodingarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("svc-decode-layer"), Has_val: 1, Desc: libc.CString("Decode SVC stream up to given spatial layer"), Enums: nil}
var framestatsarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("framestats"), Has_val: 1, Desc: libc.CString("Output per-frame stats (.csv format)"), Enums: nil}
var rowmtarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("row-mt"), Has_val: 1, Desc: libc.CString("Enable multi-threading to run row-wise in VP9"), Enums: nil}
var lpfoptarg internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("lpf-opt"), Has_val: 1, Desc: libc.CString("Do loopfilter without waiting for all threads to sync."), Enums: nil}
var all_args [26]*internal.ArgDefT = [26]*internal.ArgDefT{&help, &codecarg, &use_yv12, &use_i420, &flipuvarg, &rawvideo, &noblitarg, &progressarg, &limitarg, &skiparg, &postprocarg, &summaryarg, &outputfile, &threadsarg, &frameparallelarg, &verbosearg, &scalearg, &fb_arg, &md5arg, &error_concealment, &continuearg, &svcdecodingarg, &framestatsarg, &rowmtarg, &lpfoptarg, nil}
var addnoise_level internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("noise-level"), Has_val: 1, Desc: libc.CString("Enable VP8 postproc add noise"), Enums: nil}
var deblock internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("deblock"), Has_val: 0, Desc: libc.CString("Enable VP8 deblocking"), Enums: nil}
var demacroblock_level internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("demacroblock-level"), Has_val: 1, Desc: libc.CString("Enable VP8 demacroblocking, w/ level"), Enums: nil}
var mfqe internal.ArgDefT = internal.ArgDefT{Short_name: nil, Long_name: libc.CString("mfqe"), Has_val: 0, Desc: libc.CString("Enable multiframe quality enhancement"), Enums: nil}
var vp8_pp_args [5]*internal.ArgDefT = [5]*internal.ArgDefT{&addnoise_level, &deblock, &demacroblock_level, &mfqe, nil}

func libyuv_scale(src *vpx.Image, dst *vpx.Image, mode internal.FilterModeEnum) int {
	libc.Assert(src.Fmt == vpx.ImgFmt(vpx.VPX_IMG_FMT_I420))
	libc.Assert(dst.Fmt == vpx.ImgFmt(vpx.VPX_IMG_FMT_I420))
	return internal.I420Scale(src.Planes[vpx.VPX_PLANE_Y], src.Stride[vpx.VPX_PLANE_Y], src.Planes[vpx.VPX_PLANE_U], src.Stride[vpx.VPX_PLANE_U], src.Planes[vpx.VPX_PLANE_V], src.Stride[vpx.VPX_PLANE_V], int(src.D_w), int(src.D_h), dst.Planes[vpx.VPX_PLANE_Y], dst.Stride[vpx.VPX_PLANE_Y], dst.Planes[vpx.VPX_PLANE_U], dst.Stride[vpx.VPX_PLANE_U], dst.Planes[vpx.VPX_PLANE_V], dst.Stride[vpx.VPX_PLANE_V], int(dst.D_w), int(dst.D_h), internal.FilterMode(mode))
}
func show_help(fout *stdio.File, shorthelp int) {
	var i int
	stdio.Fprintf(fout, "Usage: %s <options> filename\n\n", exec_name)
	if shorthelp != 0 {
		stdio.Fprintf(fout, "Use --help to see the full list of options.\n")
		return
	}
	stdio.Fprintf(fout, "Options:\n")
	internal.ArgShowUsage(fout, (**internal.ArgDef)(unsafe.Pointer(&all_args[0])))
	stdio.Fprintf(fout, "\nVP8 Postprocessing Options:\n")
	internal.ArgShowUsage(fout, (**internal.ArgDef)(unsafe.Pointer(&vp8_pp_args[0])))
	stdio.Fprintf(fout, "\nOutput File Patterns:\n\n  The -o argument specifies the name of the file(s) to write to. If the\n  argument does not include any escape characters, the output will be\n  written to a single file. Otherwise, the filename will be calculated by\n  expanding the following escape characters:\n")
	stdio.Fprintf(fout, "\n\t%%w   - Frame width\n\t%%h   - Frame height\n\t%%<n> - Frame number, zero padded to <n> places (1..9)\n\n  Pattern arguments are only supported in conjunction with the --yv12 and\n  --i420 options. If the -o option is not specified, the output will be\n  directed to stdout.\n")
	stdio.Fprintf(fout, "\nIncluded decoders:\n\n")
	for i = 0; i < internal.GetVpxDecoderCount(); i++ {
		var decoder *internal.VpxInterface = internal.GetVpxDecoderByIndex(i)
		stdio.Fprintf(fout, "    %-6s - %s\n", decoder.Name, vpx.CodecIFaceName(decoder.Codec_interface()))
	}
}
func usage_exit() {
	show_help(stdio.Stderr(), 1)
	os.Exit(0)
}
func raw_read_frame(infile *stdio.File, buffer **uint8, bytes_read *uint64, buffer_size *uint64) int {
	var (
		raw_hdr    [4]byte
		frame_size uint64 = 0
	)
	if int(infile.ReadN(&raw_hdr[0], int(unsafe.Sizeof(uint32(0))), 1)) != 1 {
		if int(infile.IsEOF()) == 0 {
			internal.Warn(libc.CString("Failed to read RAW frame size\n"))
		}
	} else {
		var (
			kCorruptFrameThreshold  uint64 = 256 * 1024 * 1024
			kFrameTooSmallThreshold uint64 = 256 * 1024
		)
		frame_size = uint64(ports.GetLe32AsInt(unsafe.Pointer(&raw_hdr[0])))
		if frame_size > kCorruptFrameThreshold {
			internal.Warn(libc.CString("Read invalid frame size (%u)\n"), uint(frame_size))
			frame_size = 0
		}
		if frame_size < kFrameTooSmallThreshold {
			internal.Warn(libc.CString("Warning: Read invalid frame size (%u) - not a raw file?\n"), uint(frame_size))
		}
		if frame_size > *buffer_size {
			var new_buf *uint8 = (*uint8)(libc.Realloc(unsafe.Pointer(*buffer), int(frame_size*2)))
			if new_buf != nil {
				*buffer = new_buf
				*buffer_size = frame_size * 2
			} else {
				internal.Warn(libc.CString("Failed to allocate compressed data buffer\n"))
				frame_size = 0
			}
		}
	}
	if int(infile.IsEOF()) == 0 {
		if uint64(infile.ReadN((*byte)(unsafe.Pointer(*buffer)), 1, int(frame_size))) != frame_size {
			internal.Warn(libc.CString("Failed to read full frame\n"))
			return 1
		}
		*bytes_read = frame_size
		return 0
	}
	return 1
}
func dec_read_frame(input *VpxDecInputContext, buf **uint8, bytes_in_buffer *uint64, buffer_size *uint64) int {
	switch input.Vpx_input_ctx.File_type {
	case internal.FILE_TYPE_WEBM:
		return internal.WebmReadFrame(input.Webm_ctx, buf, bytes_in_buffer)
	case internal.FILE_TYPE_RAW:
		return raw_read_frame(input.Vpx_input_ctx.File, buf, bytes_in_buffer, buffer_size)
	case internal.FILE_TYPE_IVF:
		return internal.IvfReadFrame(input.Vpx_input_ctx.File, buf, bytes_in_buffer, buffer_size)
	default:
		return 1
	}
}
func update_image_md5(img *vpx.Image, planes [3]int, md5 *internal.MD5Context) {
	var (
		i int
		y int
	)
	for i = 0; i < 3; i++ {
		var (
			plane  int    = planes[i]
			buf    *uint8 = img.Planes[plane]
			stride int    = img.Stride[plane]
			w      int    = internal.ImgPlaneWidth(img, plane) * (func() int {
				if (img.Fmt & vpx.VPX_IMG_FMT_HIGHBITDEPTH) != 0 {
					return 2
				}
				return 1
			}())
			h int = internal.ImgPlaneHeight(img, plane)
		)
		for y = 0; y < h; y++ {
			internal.MD5Update(md5, buf, uint(w))
			buf = (*uint8)(unsafe.Add(unsafe.Pointer(buf), stride))
		}
	}
}
func write_image_file(img *vpx.Image, planes [3]int, file *stdio.File) {
	var (
		i                int
		y                int
		bytes_per_sample int = 1
	)
	for i = 0; i < 3; i++ {
		var (
			plane  int    = planes[i]
			buf    *uint8 = img.Planes[plane]
			stride int    = img.Stride[plane]
			w      int    = internal.ImgPlaneWidth(img, plane)
			h      int    = internal.ImgPlaneHeight(img, plane)
		)
		for y = 0; y < h; y++ {
			file.WriteN((*byte)(unsafe.Pointer(buf)), bytes_per_sample, w)
			buf = (*uint8)(unsafe.Add(unsafe.Pointer(buf), stride))
		}
	}
}
func file_is_raw(input *internal.VpxInputContext) int {
	var (
		buf    [32]uint8
		is_raw int = 0
		si     vpx.CodecStreamInfo
	)
	si.Sz = uint(unsafe.Sizeof(vpx.CodecStreamInfo{}))
	if int(input.File.ReadN((*byte)(unsafe.Pointer(&buf[0])), 1, 32)) == 32 {
		var i int
		if ports.GetLe32AsInt(unsafe.Pointer(&buf[0])) < 256*1024*1024 {
			for i = 0; i < internal.GetVpxDecoderCount(); i++ {
				var decoder *internal.VpxInterface = internal.GetVpxDecoderByIndex(i)
				if vpx.CodecPeekStreamInfo(decoder.Codec_interface(), &buf[4], 32-4, &si) == 0 {
					is_raw = 1
					input.Fourcc = decoder.Fourcc
					input.Width = uint32(si.W)
					input.Height = uint32(si.H)
					input.Framerate.Numerator = 30
					input.Framerate.Denominator = 1
					break
				}
			}
		}
	}
	internal.Rewind(input.File)
	return is_raw
}
func show_progress(frame_in int, frame_out int, dx_time uint64) {
	stdio.Fprintf(stdio.Stderr(), "%d decoded frames/%d showed frames in %lld us (%.2f fps)\r", frame_in, frame_out, dx_time, float64(frame_out)*1e+06/float64(dx_time))
}

type ExternalFrameBuffer struct {
	Data   *uint8
	Size   uint64
	In_use int
}
type ExternalFrameBufferList struct {
	Num_external_frame_buffers int
	Ext_fb                     *ExternalFrameBuffer
}

func get_vp9_frame_buffer(cb_priv unsafe.Pointer, min_size uint64, fb *vpx.CodecFrameBuffer) int {
	var (
		i           int
		ext_fb_list *ExternalFrameBufferList = (*ExternalFrameBufferList)(cb_priv)
	)
	if ext_fb_list == nil {
		return -1
	}
	for i = 0; i < ext_fb_list.Num_external_frame_buffers; i++ {
		if (*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).In_use == 0 {
			break
		}
	}
	if i == ext_fb_list.Num_external_frame_buffers {
		return -1
	}
	if (*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).Size < min_size {
		libc.Free(unsafe.Pointer((*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).Data))
		(*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).Data = &make([]uint8, int(min_size))[0]
		if (*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).Data == nil {
			return -1
		}
		(*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).Size = min_size
	}
	fb.Data = (*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).Data
	fb.Size = (*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).Size
	(*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).In_use = 1
	fb.Priv = unsafe.Pointer((*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i))))
	return 0
}
func release_vp9_frame_buffer(cb_priv unsafe.Pointer, fb *vpx.CodecFrameBuffer) int {
	var ext_fb *ExternalFrameBuffer = (*ExternalFrameBuffer)(fb.Priv)
	_ = cb_priv
	ext_fb.In_use = 0
	return 0
}
func generate_filename(pattern *byte, out *byte, q_len uint64, d_w uint, d_h uint, frame_in uint) {
	var (
		p *byte = pattern
		q *byte = out
	)
	for {
		{
			var next_pat *byte = libc.StrChr(p, byte('%'))
			if p == next_pat {
				var pat_len uint64
				*(*byte)(unsafe.Add(unsafe.Pointer(q), q_len-1)) = byte('\x00')
				switch *(*byte)(unsafe.Add(unsafe.Pointer(p), 1)) {
				case 'w':
					stdio.Snprintf(q, int(q_len-1), "%d", d_w)
				case 'h':
					stdio.Snprintf(q, int(q_len-1), "%d", d_h)
				case '1':
					stdio.Snprintf(q, int(q_len-1), "%d", frame_in)
				case '2':
					stdio.Snprintf(q, int(q_len-1), "%02d", frame_in)
				case '3':
					stdio.Snprintf(q, int(q_len-1), "%03d", frame_in)
				case '4':
					stdio.Snprintf(q, int(q_len-1), "%04d", frame_in)
				case '5':
					stdio.Snprintf(q, int(q_len-1), "%05d", frame_in)
				case '6':
					stdio.Snprintf(q, int(q_len-1), "%06d", frame_in)
				case '7':
					stdio.Snprintf(q, int(q_len-1), "%07d", frame_in)
				case '8':
					stdio.Snprintf(q, int(q_len-1), "%08d", frame_in)
				case '9':
					stdio.Snprintf(q, int(q_len-1), "%09d", frame_in)
				default:
					internal.Die(libc.CString("Unrecognized pattern %%%c\n"), *(*byte)(unsafe.Add(unsafe.Pointer(p), 1)))
				}
				pat_len = uint64(libc.StrLen(q))
				if pat_len >= q_len-1 {
					internal.Die(libc.CString("Output filename too long.\n"))
				}
				q = (*byte)(unsafe.Add(unsafe.Pointer(q), pat_len))
				p = (*byte)(unsafe.Add(unsafe.Pointer(p), 2))
				q_len -= pat_len
			} else {
				var copy_len uint64
				if next_pat == nil {
					copy_len = uint64(libc.StrLen(p))
				} else {
					copy_len = uint64(int64(uintptr(unsafe.Pointer(next_pat)) - uintptr(unsafe.Pointer(p))))
				}
				if copy_len >= q_len-1 {
					internal.Die(libc.CString("Output filename too long.\n"))
				}
				libc.MemCpy(unsafe.Pointer(q), unsafe.Pointer(p), int(copy_len))
				*(*byte)(unsafe.Add(unsafe.Pointer(q), copy_len)) = byte('\x00')
				q = (*byte)(unsafe.Add(unsafe.Pointer(q), copy_len))
				p = (*byte)(unsafe.Add(unsafe.Pointer(p), copy_len))
				q_len -= copy_len
			}
		}
		if *p == 0 {
			break
		}
	}
}
func is_single_file(outfile_pattern *byte) int {
	var p *byte = outfile_pattern
	for {
		p = libc.StrChr(p, byte('%'))
		if p != nil && *(*byte)(unsafe.Add(unsafe.Pointer(p), 1)) >= byte('1') && *(*byte)(unsafe.Add(unsafe.Pointer(p), 1)) <= byte('9') {
			return 0
		}
		if p != nil {
			p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
		}
		if p == nil {
			break
		}
	}
	return 1
}
func print_md5(digest [16]uint8, filename *byte) {
	var i int
	for i = 0; i < 16; i++ {
		stdio.Printf("%02x", digest[i])
	}
	stdio.Printf("  %s\n", filename)
}
func open_outfile(name *byte) *stdio.File {
	if libc.StrCmp(libc.CString("-"), name) == 0 {
		internal.SetBinaryMode(stdio.Stdout())
		return stdio.Stdout()
	} else {
		var file *stdio.File = stdio.FOpen(libc.GoString(name), "wb")
		if file == nil {
			internal.Fatal(libc.CString("Failed to open output file '%s'"), name)
		}
		return file
	}
}
func main_loop(argc int, argv_ **byte) int {
	var (
		decoder                    vpx.CodecCtx
		fn                         *byte = nil
		i                          int
		ret                        int    = 0
		buf                        *uint8 = nil
		bytes_in_buffer            uint64 = 0
		buffer_size                uint64 = 0
		infile                     *stdio.File
		frame_in                   int                    = 0
		frame_out                  int                    = 0
		flipuv                     int                    = 0
		noblit                     int                    = 0
		do_md5                     int                    = 0
		progress                   int                    = 0
		stop_after                 int                    = 0
		postproc                   int                    = 0
		summary                    int                    = 0
		quiet                      int                    = 1
		arg_skip                   int                    = 0
		ec_enabled                 int                    = 0
		keep_going                 int                    = 0
		enable_row_mt              int                    = 0
		enable_lpf_opt             int                    = 0
		interface_                 *internal.VpxInterface = nil
		fourcc_interface           *internal.VpxInterface = nil
		dx_time                    uint64                 = 0
		arg                        internal.Arg
		argv                       **byte
		argi                       **byte
		argj                       **byte
		single_file                int
		use_y4m                    int                = 1
		opt_yv12                   int                = 0
		opt_i420                   int                = 0
		cfg                        vpx.CodecDecCfg    = vpx.CodecDecCfg{}
		svc_decoding               int                = 0
		svc_spatial_layer          int                = 0
		vp8_pp_cfg                 vpx.Vp8PostProcCfg = vpx.Vp8PostProcCfg{}
		frames_corrupted           int                = 0
		dec_flags                  int                = 0
		do_scale                   int                = 0
		scaled_img                 *vpx.Image         = nil
		frame_avail                int
		got_data                   int
		flush_decoder              int                     = 0
		num_external_frame_buffers int                     = 0
		ext_fb_list                ExternalFrameBufferList = ExternalFrameBufferList{}
		outfile_pattern            *byte                   = nil
		outfile_name               [512]byte               = [512]byte{}
		outfile                    *stdio.File             = nil
		framestats_file            *stdio.File             = nil
		md5_ctx                    internal.MD5Context
		md5_digest                 [16]uint8
		input                      VpxDecInputContext = VpxDecInputContext{}
		vpx_input_ctx              internal.VpxInputContext
		webm_ctx                   internal.WebmInputContext
	)
	webm_ctx = internal.WebmInputContext{}
	input.Webm_ctx = &webm_ctx
	input.Vpx_input_ctx = &vpx_input_ctx
	exec_name = *(**byte)(unsafe.Add(unsafe.Pointer(argv_), unsafe.Sizeof((*byte)(nil))*0))
	argv = internal.ArgvDup(argc-1, (**byte)(unsafe.Add(unsafe.Pointer(argv_), unsafe.Sizeof((*byte)(nil))*1)))
	for argi = func() **byte {
		argj = argv
		return argj
	}(); (func() *byte {
		p := argj
		*argj = *argi
		return *p
	}()) != nil; argi = (**byte)(unsafe.Add(unsafe.Pointer(argi), unsafe.Sizeof((*byte)(nil))*uintptr(arg.Argv_step))) {
		arg = internal.Arg{}
		arg.Argv_step = 1
		if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&help)), argi) != 0 {
			show_help(stdio.Stdout(), 0)
			os.Exit(1)
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&codecarg)), argi) != 0 {
			interface_ = internal.GetVpxDecoderByName(arg.Val)
			if interface_ == nil {
				internal.Die(libc.CString("Error: Unrecognized argument (%s) to --codec\n"), arg.Val)
			}
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&looparg)), argi) != 0 {
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&outputfile)), argi) != 0 {
			outfile_pattern = arg.Val
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&use_yv12)), argi) != 0 {
			use_y4m = 0
			flipuv = 1
			opt_yv12 = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&use_i420)), argi) != 0 {
			use_y4m = 0
			flipuv = 0
			opt_i420 = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&rawvideo)), argi) != 0 {
			use_y4m = 0
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&flipuvarg)), argi) != 0 {
			flipuv = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&noblitarg)), argi) != 0 {
			noblit = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&progressarg)), argi) != 0 {
			progress = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&limitarg)), argi) != 0 {
			stop_after = int(internal.ArgParseUint(&arg))
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&skiparg)), argi) != 0 {
			arg_skip = int(internal.ArgParseUint(&arg))
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&postprocarg)), argi) != 0 {
			postproc = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&md5arg)), argi) != 0 {
			do_md5 = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&summaryarg)), argi) != 0 {
			summary = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&threadsarg)), argi) != 0 {
			cfg.Threads = internal.ArgParseUint(&arg)
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&frameparallelarg)), argi) != 0 {
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&verbosearg)), argi) != 0 {
			quiet = 0
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&scalearg)), argi) != 0 {
			do_scale = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&fb_arg)), argi) != 0 {
			num_external_frame_buffers = int(internal.ArgParseUint(&arg))
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&continuearg)), argi) != 0 {
			keep_going = 1
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&svcdecodingarg)), argi) != 0 {
			svc_decoding = 1
			svc_spatial_layer = int(internal.ArgParseUint(&arg))
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&framestatsarg)), argi) != 0 {
			framestats_file = stdio.FOpen(libc.GoString(arg.Val), "w")
			if framestats_file == nil {
				internal.Die(libc.CString("Error: Could not open --framestats file (%s) for writing.\n"), arg.Val)
			}
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&rowmtarg)), argi) != 0 {
			enable_row_mt = int(internal.ArgParseUint(&arg))
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&lpfoptarg)), argi) != 0 {
			enable_lpf_opt = int(internal.ArgParseUint(&arg))
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&addnoise_level)), argi) != 0 {
			postproc = 1
			vp8_pp_cfg.Post_proc_flag |= int(vpx.VP8_ADDNOISE)
			vp8_pp_cfg.Noise_level = int(internal.ArgParseUint(&arg))
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&demacroblock_level)), argi) != 0 {
			postproc = 1
			vp8_pp_cfg.Post_proc_flag |= int(vpx.VP8_DEMACROBLOCK)
			vp8_pp_cfg.Deblocking_level = int(internal.ArgParseUint(&arg))
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&deblock)), argi) != 0 {
			postproc = 1
			vp8_pp_cfg.Post_proc_flag |= int(vpx.VP8_DEBLOCK)
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&mfqe)), argi) != 0 {
			postproc = 1
			vp8_pp_cfg.Post_proc_flag |= int(vpx.VP8_MFQE)
		} else if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&error_concealment)), argi) != 0 {
			ec_enabled = 1
		} else {
			argj = (**byte)(unsafe.Add(unsafe.Pointer(argj), unsafe.Sizeof((*byte)(nil))*1))
		}
	}
	for argi = argv; *argi != nil; argi = (**byte)(unsafe.Add(unsafe.Pointer(argi), unsafe.Sizeof((*byte)(nil))*1)) {
		if *(*byte)(unsafe.Add(unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(argi), unsafe.Sizeof((*byte)(nil))*0))), 0)) == byte('-') && libc.StrLen(*(**byte)(unsafe.Add(unsafe.Pointer(argi), unsafe.Sizeof((*byte)(nil))*0))) > 1 {
			internal.Die(libc.CString("Error: Unrecognized option %s\n"), *argi)
		}
	}
	fn = *(**byte)(unsafe.Add(unsafe.Pointer(argv), unsafe.Sizeof((*byte)(nil))*0))
	if fn == nil {
		libc.Free(unsafe.Pointer(argv))
		stdio.Fprintf(stdio.Stderr(), "No input file specified!\n")
		usage_exit()
	}
	if libc.StrCmp(fn, libc.CString("-")) != 0 {
		infile = stdio.FOpen(libc.GoString(fn), "rb")
	} else {
		infile = internal.SetBinaryMode(stdio.Stdin())
	}
	if infile == nil {
		internal.Fatal(libc.CString("Failed to open input file '%s'"), func() *byte {
			if libc.StrCmp(fn, libc.CString("-")) != 0 {
				return fn
			}
			return libc.CString("stdin")
		}())
	}
	if outfile_pattern == nil && internal.IsAtty(int(stdio.Stdout().FileNo())) != 0 && do_md5 == 0 && noblit == 0 {
		stdio.Fprintf(stdio.Stderr(), "Not dumping raw video to your terminal. Use '-o -' to override.\n")
		return 0
	}
	input.Vpx_input_ctx.File = infile
	if internal.FileIsIvf(input.Vpx_input_ctx) != 0 {
		input.Vpx_input_ctx.File_type = internal.VideoFileType(internal.FILE_TYPE_IVF)
	} else if internal.FileIsWebm(input.Webm_ctx, input.Vpx_input_ctx) != 0 {
		input.Vpx_input_ctx.File_type = internal.VideoFileType(internal.FILE_TYPE_WEBM)
	} else if file_is_raw(input.Vpx_input_ctx) != 0 {
		input.Vpx_input_ctx.File_type = internal.VideoFileType(internal.FILE_TYPE_RAW)
	} else {
		stdio.Fprintf(stdio.Stderr(), "Unrecognized input file type.\n")
		libc.Free(unsafe.Pointer(argv))
		return 0
	}
	if outfile_pattern != nil {
		outfile_pattern = outfile_pattern
	} else {
		outfile_pattern = libc.CString("-")
	}
	single_file = is_single_file(outfile_pattern)
	if noblit == 0 && single_file != 0 {
		generate_filename(outfile_pattern, &outfile_name[0], internal.PATH_MAX, uint(vpx_input_ctx.Width), uint(vpx_input_ctx.Height), 0)
		if do_md5 != 0 {
			internal.MD5Init(&md5_ctx)
		} else {
			outfile = open_outfile(&outfile_name[0])
		}
	}
	if use_y4m != 0 && noblit == 0 {
		if single_file == 0 {
			stdio.Fprintf(stdio.Stderr(), "YUV4MPEG2 not supported with output patterns, try --i420 or --yv12 or --rawvideo.\n")
			return 0
		}
		if vpx_input_ctx.File_type == internal.VideoFileType(internal.FILE_TYPE_WEBM) {
			if internal.WebmGuessFrameRate(input.Webm_ctx, input.Vpx_input_ctx) != 0 {
				stdio.Fprintf(stdio.Stderr(), "Failed to guess framerate -- error parsing webm file?\n")
				return 0
			}
		}
	}
	fourcc_interface = internal.GetVpxDecoderByFourcc(vpx_input_ctx.Fourcc)
	if interface_ != nil && fourcc_interface != nil && interface_ != fourcc_interface {
		internal.Warn(libc.CString("Header indicates codec: %s\n"), fourcc_interface.Name)
	} else {
		interface_ = fourcc_interface
	}
	if interface_ == nil {
		interface_ = internal.GetVpxDecoderByIndex(0)
	}
	dec_flags = (func() int {
		if postproc != 0 {
			return vpx.VPX_CODEC_USE_POSTPROC
		}
		return 0
	}()) | (func() int {
		if ec_enabled != 0 {
			return vpx.VPX_CODEC_USE_ERROR_CONCEALMENT
		}
		return 0
	}())
	if vpx.CodecDecInitVer(&decoder, interface_.Codec_interface(), &cfg, vpx.CodecFlags(dec_flags), 3+(4+5)) != 0 {
		stdio.Fprintf(stdio.Stderr(), "Failed to initialize decoder: %s\n", vpx.CodecError(&decoder))
		goto fail2
	}
	if svc_decoding != 0 {
		if vpx.Vp9DecodeSvcSpatialLayer(&decoder, int(vpx.VP9_DECODE_SVC_SPATIAL_LAYER), svc_spatial_layer) != 0 {
			stdio.Fprintf(stdio.Stderr(), "Failed to set spatial layer for svc decode: %s\n", vpx.CodecError(&decoder))
			goto fail
		}
	}
	if int(interface_.Fourcc) == internal.VP9_FOURCC && vpx.Vp9dSetRowMt(&decoder, int(vpx.VP9D_SET_ROW_MT), enable_row_mt) != 0 {
		stdio.Fprintf(stdio.Stderr(), "Failed to set decoder in row multi-thread mode: %s\n", vpx.CodecError(&decoder))
		goto fail
	}
	if int(interface_.Fourcc) == internal.VP9_FOURCC && vpx.Vp9dSetLoopFilterOpt(&decoder, int(vpx.VP9D_SET_LOOP_FILTER_OPT), enable_lpf_opt) != 0 {
		stdio.Fprintf(stdio.Stderr(), "Failed to set decoder in optimized loopfilter mode: %s\n", vpx.CodecError(&decoder))
		goto fail
	}
	if quiet == 0 {
		stdio.Fprintf(stdio.Stderr(), "%s\n", decoder.Name)
	}
	if vp8_pp_cfg.Post_proc_flag != 0 && vpx.Vp8SetPostProc(&decoder, int(vpx.VP8_SET_POSTPROC), &vp8_pp_cfg) != 0 {
		stdio.Fprintf(stdio.Stderr(), "Failed to configure postproc: %s\n", vpx.CodecError(&decoder))
		goto fail
	}
	if arg_skip != 0 {
		stdio.Fprintf(stdio.Stderr(), "Skipping first %d frames.\n", arg_skip)
	}
	for arg_skip != 0 {
		if dec_read_frame(&input, &buf, &bytes_in_buffer, &buffer_size) != 0 {
			break
		}
		arg_skip--
	}
	if num_external_frame_buffers > 0 {
		ext_fb_list.Num_external_frame_buffers = num_external_frame_buffers
		ext_fb_list.Ext_fb = &make([]ExternalFrameBuffer, num_external_frame_buffers)[0]
		if vpx.CodecSetFrameBufferFunctions(&decoder, func(priv unsafe.Pointer, min_size uint64, fb *vpx.CodecFrameBuffer) int {
			return get_vp9_frame_buffer(priv, min_size, fb)
		}, func(priv unsafe.Pointer, fb *vpx.CodecFrameBuffer) int {
			return release_vp9_frame_buffer(priv, fb)
		}, unsafe.Pointer(&ext_fb_list)) != 0 {
			stdio.Fprintf(stdio.Stderr(), "Failed to configure external frame buffers: %s\n", vpx.CodecError(&decoder))
			goto fail
		}
	}
	frame_avail = 1
	got_data = 0
	if framestats_file != nil {
		stdio.Fprintf(framestats_file, "bytes,qp\n")
	}
	for frame_avail != 0 || got_data != 0 {
		var (
			iter      vpx.CodecIter = nil
			img       *vpx.Image
			timer     vpx.VpxUsecTimer
			corrupted int = 0
		)
		frame_avail = 0
		if stop_after == 0 || frame_in < stop_after {
			if dec_read_frame(&input, &buf, &bytes_in_buffer, &buffer_size) == 0 {
				frame_avail = 1
				frame_in++
				vpx.VpxUsecTimerStart(&timer)
				if vpx.CodecDecode(&decoder, buf, uint(bytes_in_buffer), nil, 0) != 0 {
					var detail *byte = vpx.CodecErrorDetail(&decoder)
					internal.Warn(libc.CString("Failed to decode frame %d: %s"), frame_in, vpx.CodecError(&decoder))
					if detail != nil {
						internal.Warn(libc.CString("Additional information: %s"), detail)
					}
					corrupted = 1
					if keep_going == 0 {
						goto fail
					}
				}
				if framestats_file != nil {
					var qp int
					if vpx.VpxdGetLastQuantizer(&decoder, int(vpx.VPXD_GET_LAST_QUANTIZER), &qp) != 0 {
						internal.Warn(libc.CString("Failed int(vpx.VPXD_GET_LAST_QUANTIZER): %s"), vpx.CodecError(&decoder))
						if keep_going == 0 {
							goto fail
						}
					}
					stdio.Fprintf(framestats_file, "%d,%d\n", int(bytes_in_buffer), qp)
				}
				vpx.VpxUsecTimerMark(&timer)
				dx_time += uint64(vpx.VpxUsecTimerElapsed(&timer))
			} else {
				flush_decoder = 1
			}
		} else {
			flush_decoder = 1
		}
		vpx.VpxUsecTimerStart(&timer)
		if flush_decoder != 0 {
			if vpx.CodecDecode(&decoder, nil, 0, nil, 0) != 0 {
				internal.Warn(libc.CString("Failed to flush decoder: %s"), vpx.CodecError(&decoder))
				corrupted = 1
				if keep_going == 0 {
					goto fail
				}
			}
		}
		got_data = 0
		if (func() *vpx.Image {
			img = vpx.CodecGetFrame(&decoder, &iter)
			return img
		}()) != nil {
			frame_out++
			got_data = 1
		}
		vpx.VpxUsecTimerMark(&timer)
		dx_time += uint64(uint(vpx.VpxUsecTimerElapsed(&timer)))
		if corrupted == 0 && vpx.Vp8dGetFrameCorrupted(&decoder, int(vpx.VP8D_GET_FRAME_CORRUPTED), &corrupted) != 0 {
			internal.Warn(libc.CString("Failed VP8_GET_FRAME_CORRUPTED: %s"), vpx.CodecError(&decoder))
			if keep_going == 0 {
				goto fail
			}
		}
		frames_corrupted += corrupted
		if progress != 0 {
			show_progress(frame_in, frame_out, dx_time)
		}
		if noblit == 0 && img != nil {
			var (
				PLANES_YUV [3]int = [3]int{vpx.VPX_PLANE_Y, vpx.VPX_PLANE_U, vpx.VPX_PLANE_V}
				PLANES_YVU [3]int = [3]int{vpx.VPX_PLANE_Y, vpx.VPX_PLANE_V, vpx.VPX_PLANE_U}
				planes     [3]int
			)
			if flipuv != 0 {
				planes = PLANES_YVU
			} else {
				planes = PLANES_YUV
			}
			if do_scale != 0 {
				if frame_out == 1 {
					var (
						render_width  int = int(vpx_input_ctx.Width)
						render_height int = int(vpx_input_ctx.Height)
					)
					if render_width == 0 || render_height == 0 {
						var render_size [2]int
						if vpx.Vp9dGetDisplaySize(&decoder, int(vpx.VP9D_GET_DISPLAY_SIZE), &render_size[0]) != 0 {
							render_width = int(img.D_w)
							render_height = int(img.D_h)
						} else {
							render_width = render_size[0]
							render_height = render_size[1]
						}
					}
					scaled_img = vpx.ImgAlloc(nil, img.Fmt, uint(render_width), uint(render_height), 16)
					scaled_img.Bit_depth = img.Bit_depth
				}
				if img.D_w != scaled_img.D_w || img.D_h != scaled_img.D_h {
					libyuv_scale(img, scaled_img, internal.FilterModeEnum(internal.KFilterBox))
					img = scaled_img
				}
			}
			if single_file != 0 {
				if use_y4m != 0 {
					var (
						buf  [128]byte = [128]byte{}
						len_ uint64    = 0
					)
					if img.Fmt == vpx.ImgFmt(vpx.VPX_IMG_FMT_I440) || img.Fmt == vpx.ImgFmt(vpx.VPX_IMG_FMT_I44016) {
						stdio.Fprintf(stdio.Stderr(), "Cannot produce y4m output for 440 sampling.\n")
						goto fail
					}
					if frame_out == 1 {
						len_ = uint64(internal.Y4mWriteFileHeader(&buf[0], uint64(unsafe.Sizeof([128]byte{})), int(vpx_input_ctx.Width), int(vpx_input_ctx.Height), &vpx_input_ctx.Framerate, img.Fmt, img.Bit_depth))
						if do_md5 != 0 {
							internal.MD5Update(&md5_ctx, (*uint8)(unsafe.Pointer(&buf[0])), uint(len_))
						} else {
							outfile.PutS(&buf[0])
						}
					}
					len_ = uint64(internal.Y4mWriteFrameHeader(&buf[0], uint64(unsafe.Sizeof([128]byte{}))))
					if do_md5 != 0 {
						internal.MD5Update(&md5_ctx, (*uint8)(unsafe.Pointer(&buf[0])), uint(len_))
					} else {
						outfile.PutS(&buf[0])
					}
				} else {
					if frame_out == 1 {
						if opt_i420 != 0 {
							if img.Fmt != vpx.ImgFmt(vpx.VPX_IMG_FMT_I420) && img.Fmt != vpx.ImgFmt(vpx.VPX_IMG_FMT_I42016) {
								stdio.Fprintf(stdio.Stderr(), "Cannot produce i420 output for bit-stream.\n")
								goto fail
							}
						}
						if opt_yv12 != 0 {
							if img.Fmt != vpx.ImgFmt(vpx.VPX_IMG_FMT_I420) && img.Fmt != vpx.ImgFmt(vpx.VPX_IMG_FMT_YV12) || img.Bit_depth != 8 {
								stdio.Fprintf(stdio.Stderr(), "Cannot produce yv12 output for bit-stream.\n")
								goto fail
							}
						}
					}
				}
				if do_md5 != 0 {
					update_image_md5(img, ([3]int)(planes), &md5_ctx)
				} else {
					if corrupted == 0 {
						write_image_file(img, ([3]int)(planes), outfile)
					}
				}
			} else {
				generate_filename(outfile_pattern, &outfile_name[0], internal.PATH_MAX, img.D_w, img.D_h, uint(frame_in))
				if do_md5 != 0 {
					internal.MD5Init(&md5_ctx)
					update_image_md5(img, ([3]int)(planes), &md5_ctx)
					internal.MD5Final(md5_digest, &md5_ctx)
					print_md5(md5_digest, &outfile_name[0])
				} else {
					outfile = open_outfile(&outfile_name[0])
					write_image_file(img, ([3]int)(planes), outfile)
					outfile.Close()
				}
			}
		}
	}
	if summary != 0 || progress != 0 {
		show_progress(frame_in, frame_out, dx_time)
		stdio.Fprintf(stdio.Stderr(), "\n")
	}
	if frames_corrupted != 0 {
		stdio.Fprintf(stdio.Stderr(), "WARNING: %d frames corrupted.\n", frames_corrupted)
	} else {
		ret = 1
	}
fail:
	if vpx.CodecDestroy(&decoder) != 0 {
		stdio.Fprintf(stdio.Stderr(), "Failed to destroy decoder: %s\n", vpx.CodecError(&decoder))
	}
fail2:
	if noblit == 0 && single_file != 0 {
		if do_md5 != 0 {
			internal.MD5Final(md5_digest, &md5_ctx)
			print_md5(md5_digest, &outfile_name[0])
		} else {
			outfile.Close()
		}
	}
	if input.Vpx_input_ctx.File_type == internal.VideoFileType(internal.FILE_TYPE_WEBM) {
		internal.WebmFree(input.Webm_ctx)
	}
	if input.Vpx_input_ctx.File_type != internal.VideoFileType(internal.FILE_TYPE_WEBM) {
		libc.Free(unsafe.Pointer(buf))
	}
	if scaled_img != nil {
		vpx.ImgFree(scaled_img)
	}
	for i = 0; i < ext_fb_list.Num_external_frame_buffers; i++ {
		libc.Free(unsafe.Pointer((*(*ExternalFrameBuffer)(unsafe.Add(unsafe.Pointer(ext_fb_list.Ext_fb), unsafe.Sizeof(ExternalFrameBuffer{})*uintptr(i)))).Data))
	}
	libc.Free(unsafe.Pointer(ext_fb_list.Ext_fb))
	infile.Close()
	if framestats_file != nil {
		framestats_file.Close()
	}
	libc.Free(unsafe.Pointer(argv))
	return ret
}
func main() {
	var (
		argc  int    = len(os.Args)
		argv_ **byte = libc.CStringSlice(os.Args)
		loops uint   = 1
		i     uint
		argv  **byte
		argi  **byte
		argj  **byte
		arg   internal.Arg
		error int = 0
	)
	argv = internal.ArgvDup(argc-1, (**byte)(unsafe.Add(unsafe.Pointer(argv_), unsafe.Sizeof((*byte)(nil))*1)))
	for argi = func() **byte {
		argj = argv
		return argj
	}(); (func() *byte {
		p := argj
		*argj = *argi
		return *p
	}()) != nil; argi = (**byte)(unsafe.Add(unsafe.Pointer(argi), unsafe.Sizeof((*byte)(nil))*uintptr(arg.Argv_step))) {
		arg = internal.Arg{}
		arg.Argv_step = 1
		if internal.ArgMatch(&arg, (*internal.ArgDef)(unsafe.Pointer(&looparg)), argi) != 0 {
			loops = internal.ArgParseUint(&arg)
			break
		}
	}
	libc.Free(unsafe.Pointer(argv))
	for i = 0; error == 0 && i < loops; i++ {
		error = main_loop(argc, argv_)
	}
	os.Exit(error)
}
