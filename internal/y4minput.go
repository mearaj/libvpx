package internal

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"github.com/mearaj/libvpx/internal/vpx"
	"math"
	"unsafe"
)

type Y4mInput struct {
	Pic_w           int
	Pic_h           int
	Fps_n           int
	Fps_d           int
	Par_n           int
	Par_d           int
	Interlace       int8
	Src_c_dec_h     int
	Src_c_dec_v     int
	Dst_c_dec_h     int
	Dst_c_dec_v     int
	Chroma_type     [16]byte
	Dst_buf_sz      uint64
	Dst_buf_read_sz uint64
	Aux_buf_sz      uint64
	Aux_buf_read_sz uint64
	Convert         y4m_convert_func
	Dst_buf         *uint8
	Aux_buf         *uint8
	Vpx_fmt         vpx.ImgFmt
	Bps             int
	Bit_depth       uint
}
type y4m_convert_func func(_y4m *Y4mInput, _dst *uint8, _src *uint8)

func file_read(buf unsafe.Pointer, size uint64, file *stdio.File) int {
	var (
		kMaxRetries int = 5
		retry_count int = 0
		file_error  int
		len_        uint64 = 0
	)
	for {
		{
			var n uint64 = uint64(file.ReadN((*byte)(unsafe.Add(unsafe.Pointer((*uint8)(buf)), len_)), 1, int(size-len_)))
			len_ += n
			file_error = int(ferror(file))
			if file_error != 0 {
				if libc.Errno == 4 || libc.Errno == 11 {
					clearerr(file)
					continue
				} else {
					stdio.Fprintf(stdio.Stderr(), "Error reading file: %u of %u bytes read, %d: %s\n", uint32(len_), uint32(size), libc.Errno, libc.StrError(libc.Errno))
					return 0
				}
			}
		}
		if !(int(file.IsEOF()) == 0 && len_ < size && func() int {
			p := &retry_count
			*p++
			return *p
		}() < kMaxRetries) {
			break
		}
	}
	if int(file.IsEOF()) == 0 && len_ != size {
		stdio.Fprintf(stdio.Stderr(), "Error reading file: %u of %u bytes read, error: %d, retries: %d, %d: %s\n", uint32(len_), uint32(size), file_error, retry_count, libc.Errno, libc.StrError(libc.Errno))
	}
	return int(libc.BoolToInt(len_ == size))
}
func y4m_parse_tags(_y4m *Y4mInput, _tags *byte) int {
	var (
		p *byte
		q *byte
	)
	for p = _tags; ; p = q {
		for *p == byte(' ') {
			p = (*byte)(unsafe.Add(unsafe.Pointer(p), 1))
		}
		if *(*byte)(unsafe.Add(unsafe.Pointer(p), 0)) == byte('\x00') {
			break
		}
		for q = (*byte)(unsafe.Add(unsafe.Pointer(p), 1)); *q != byte('\x00') && *q != byte(' '); q = (*byte)(unsafe.Add(unsafe.Pointer(q), 1)) {
		}
		switch *(*byte)(unsafe.Add(unsafe.Pointer(p), 0)) {
		case 'W':
			if stdio.Sscanf((*byte)(unsafe.Add(unsafe.Pointer(p), 1)), "%d", &_y4m.Pic_w) != 1 {
				return -1
			}
		case 'H':
			if stdio.Sscanf((*byte)(unsafe.Add(unsafe.Pointer(p), 1)), "%d", &_y4m.Pic_h) != 1 {
				return -1
			}
		case 'F':
			if stdio.Sscanf((*byte)(unsafe.Add(unsafe.Pointer(p), 1)), "%d:%d", &_y4m.Fps_n, &_y4m.Fps_d) != 2 {
				return -1
			}
		case 'I':
			_y4m.Interlace = int8(*(*byte)(unsafe.Add(unsafe.Pointer(p), 1)))
		case 'A':
			if stdio.Sscanf((*byte)(unsafe.Add(unsafe.Pointer(p), 1)), "%d:%d", &_y4m.Par_n, &_y4m.Par_d) != 2 {
				return -1
			}
		case 'C':
			if int64(uintptr(unsafe.Pointer(q))-uintptr(unsafe.Pointer(p))) > 16 {
				return -1
			}
			libc.MemCpy(unsafe.Pointer(&_y4m.Chroma_type[0]), unsafe.Add(unsafe.Pointer(p), 1), int(int64(uintptr(unsafe.Pointer(q))-uintptr(unsafe.Pointer(p)))-1))
			_y4m.Chroma_type[int64(uintptr(unsafe.Pointer(q))-uintptr(unsafe.Pointer(p)))-1] = byte('\x00')
		}
	}
	return 0
}
func copy_tag(buf *byte, buf_len uint64, end_tag *byte, file *stdio.File) int {
	var i uint64
	libc.Assert(buf_len >= 1)
	for {
		if file_read(unsafe.Pointer(buf), 1, file) == 0 {
			return 0
		}
		if *(*byte)(unsafe.Add(unsafe.Pointer(buf), 0)) != byte(' ') {
			break
		}
	}
	if *(*byte)(unsafe.Add(unsafe.Pointer(buf), 0)) == byte('\n') {
		*(*byte)(unsafe.Add(unsafe.Pointer(buf), 0)) = byte('\x00')
		*end_tag = byte('\n')
		return 1
	}
	for i = 1; i < buf_len; i++ {
		if file_read(unsafe.Add(unsafe.Pointer(buf), i), 1, file) == 0 {
			return 0
		}
		if *(*byte)(unsafe.Add(unsafe.Pointer(buf), i)) == byte(' ') || *(*byte)(unsafe.Add(unsafe.Pointer(buf), i)) == byte('\n') {
			break
		}
	}
	if i == buf_len {
		stdio.Fprintf(stdio.Stderr(), "Error: Y4M header tags must be less than %lu characters\n", uint(i))
		return 0
	}
	*end_tag = *(*byte)(unsafe.Add(unsafe.Pointer(buf), i))
	*(*byte)(unsafe.Add(unsafe.Pointer(buf), i)) = byte('\x00')
	return 1
}
func parse_tags(y4m_ctx *Y4mInput, file *stdio.File) int {
	var (
		tag [256]byte
		end int8
	)
	y4m_ctx.Pic_w = -1
	y4m_ctx.Pic_h = -1
	y4m_ctx.Fps_n = -1
	y4m_ctx.Par_n = 0
	y4m_ctx.Par_d = 0
	y4m_ctx.Interlace = '?'
	stdio.Snprintf(&y4m_ctx.Chroma_type[0], int(unsafe.Sizeof([16]byte{})), "420")
	for {
		if copy_tag(&tag[0], uint64(unsafe.Sizeof([256]byte{})), (*byte)(unsafe.Pointer(&end)), file) == 0 {
			return 0
		}
		if y4m_parse_tags(y4m_ctx, &tag[0]) != 0 {
			return 0
		}
		if int(end) == int('\n') {
			break
		}
	}
	if y4m_ctx.Pic_w == -1 {
		stdio.Fprintf(stdio.Stderr(), "Width field missing\n")
		return 0
	}
	if y4m_ctx.Pic_h == -1 {
		stdio.Fprintf(stdio.Stderr(), "Height field missing\n")
		return 0
	}
	if y4m_ctx.Fps_n == -1 {
		stdio.Fprintf(stdio.Stderr(), "FPS field missing\n")
		return 0
	}
	return 1
}
func y4m_42xmpeg2_42xjpeg_helper(_dst *uint8, _src *uint8, _c_w int, _c_h int) {
	var (
		y int
		x int
	)
	for y = 0; y < _c_h; y++ {
		for x = 0; x < (func() int {
			if _c_w > 2 {
				return 2
			}
			return _c_w
		}()); x++ {
			if 0 < (func() int {
				if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), 0)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x - 1) < 0 {
						return 0
					}
					return x - 1
				}())))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 1) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 1
				}())))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 2) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 2
				}())))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 3) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 3
				}()))) + 64) >> 7) > math.MaxUint8 {
					return math.MaxUint8
				}
				return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), 0)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x - 1) < 0 {
						return 0
					}
					return x - 1
				}())))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 1) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 1
				}())))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 2) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 2
				}())))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 3) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 3
				}()))) + 64) >> 7
			}()) {
				if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), 0)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x - 1) < 0 {
						return 0
					}
					return x - 1
				}())))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 1) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 1
				}())))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 2) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 2
				}())))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 3) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 3
				}()))) + 64) >> 7) > math.MaxUint8 {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), x)) = math.MaxUint8
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), x)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), 0)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
						if (x - 1) < 0 {
							return 0
						}
						return x - 1
					}())))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
						if (x + 1) > (_c_w - 1) {
							return _c_w - 1
						}
						return x + 1
					}())))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
						if (x + 2) > (_c_w - 1) {
							return _c_w - 1
						}
						return x + 2
					}())))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
						if (x + 3) > (_c_w - 1) {
							return _c_w - 1
						}
						return x + 3
					}()))) + 64) >> 7))
				}
			} else {
				*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), x)) = 0
			}
		}
		for ; x < _c_w-3; x++ {
			if 0 < (func() int {
				if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-2)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-1)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+1)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+2)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+3))) + 64) >> 7) > math.MaxUint8 {
					return math.MaxUint8
				}
				return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-2)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-1)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+1)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+2)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+3))) + 64) >> 7
			}()) {
				if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-2)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-1)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+1)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+2)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+3))) + 64) >> 7) > math.MaxUint8 {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), x)) = math.MaxUint8
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), x)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-2)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-1)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+1)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+2)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x+3))) + 64) >> 7))
				}
			} else {
				*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), x)) = 0
			}
		}
		for ; x < _c_w; x++ {
			if 0 < (func() int {
				if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-2)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-1)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 1) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 1
				}())))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 2) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 2
				}())))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), _c_w-1))) + 64) >> 7) > math.MaxUint8 {
					return math.MaxUint8
				}
				return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-2)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-1)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 1) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 1
				}())))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 2) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 2
				}())))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), _c_w-1))) + 64) >> 7
			}()) {
				if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-2)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-1)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 1) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 1
				}())))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
					if (x + 2) > (_c_w - 1) {
						return _c_w - 1
					}
					return x + 2
				}())))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), _c_w-1))) + 64) >> 7) > math.MaxUint8 {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), x)) = math.MaxUint8
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), x)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-2)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x-1)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), x)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
						if (x + 1) > (_c_w - 1) {
							return _c_w - 1
						}
						return x + 1
					}())))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), func() int {
						if (x + 2) > (_c_w - 1) {
							return _c_w - 1
						}
						return x + 2
					}())))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), _c_w-1))) + 64) >> 7))
				}
			} else {
				*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), x)) = 0
			}
		}
		_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), _c_w))
		_src = (*uint8)(unsafe.Add(unsafe.Pointer(_src), _c_w))
	}
}
func y4m_convert_42xpaldv_42xjpeg(_y4m *Y4mInput, _dst *uint8, _aux *uint8) {
	var (
		tmp  *uint8
		c_w  int
		c_h  int
		c_sz int
		pli  int
		y    int
		x    int
	)
	_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), _y4m.Pic_w*_y4m.Pic_h))
	c_w = (_y4m.Pic_w + 1) / 2
	c_h = (_y4m.Pic_h + _y4m.Dst_c_dec_h - 1) / _y4m.Dst_c_dec_h
	c_sz = c_w * c_h
	tmp = (*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_sz*2))
	for pli = 1; pli < 3; pli++ {
		y4m_42xmpeg2_42xjpeg_helper(tmp, _aux, c_w, c_h)
		_aux = (*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_sz))
		switch pli {
		case 1:
			for x = 0; x < c_w; x++ {
				for y = 0; y < (func() int {
					if c_h > 3 {
						return 3
					}
					return c_h
				}()); y++ {
					if 0 < (func() int {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), 0))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y - 2) < 0 {
								return 0
							}
							return y - 2
						}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y - 1) < 0 {
								return 0
							}
							return y - 1
						}())*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 2) > (c_h - 1) {
								return c_h - 1
							}
							return y + 2
						}())*c_w)))*4 + 64) >> 7) > math.MaxUint8 {
							return math.MaxUint8
						}
						return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), 0))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y - 2) < 0 {
								return 0
							}
							return y - 2
						}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y - 1) < 0 {
								return 0
							}
							return y - 1
						}())*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 2) > (c_h - 1) {
								return c_h - 1
							}
							return y + 2
						}())*c_w)))*4 + 64) >> 7
					}()) {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), 0))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y - 2) < 0 {
								return 0
							}
							return y - 2
						}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y - 1) < 0 {
								return 0
							}
							return y - 1
						}())*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 2) > (c_h - 1) {
								return c_h - 1
							}
							return y + 2
						}())*c_w)))*4 + 64) >> 7) > math.MaxUint8 {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = math.MaxUint8
						} else {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), 0))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y - 2) < 0 {
									return 0
								}
								return y - 2
							}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y - 1) < 0 {
									return 0
								}
								return y - 1
							}())*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y + 1) > (c_h - 1) {
									return c_h - 1
								}
								return y + 1
							}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y + 2) > (c_h - 1) {
									return c_h - 1
								}
								return y + 2
							}())*c_w)))*4 + 64) >> 7))
						}
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = 0
					}
				}
				for ; y < c_h-2; y++ {
					if 0 < (func() int {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-3)*c_w))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+2)*c_w)))*4 + 64) >> 7) > math.MaxUint8 {
							return math.MaxUint8
						}
						return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-3)*c_w))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+2)*c_w)))*4 + 64) >> 7
					}()) {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-3)*c_w))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+2)*c_w)))*4 + 64) >> 7) > math.MaxUint8 {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = math.MaxUint8
						} else {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-3)*c_w))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+2)*c_w)))*4 + 64) >> 7))
						}
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = 0
					}
				}
				for ; y < c_h; y++ {
					if 0 < (func() int {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-3)*c_w))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (c_h-1)*c_w)))*4 + 64) >> 7) > math.MaxUint8 {
							return math.MaxUint8
						}
						return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-3)*c_w))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (c_h-1)*c_w)))*4 + 64) >> 7
					}()) {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-3)*c_w))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (c_h-1)*c_w)))*4 + 64) >> 7) > math.MaxUint8 {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = math.MaxUint8
						} else {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-3)*c_w))) - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*35 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y + 1) > (c_h - 1) {
									return c_h - 1
								}
								return y + 1
							}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (c_h-1)*c_w)))*4 + 64) >> 7))
						}
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = 0
					}
				}
				_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), 1))
				tmp = (*uint8)(unsafe.Add(unsafe.Pointer(tmp), 1))
			}
			_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), c_sz-c_w))
			tmp = (*uint8)(unsafe.Add(unsafe.Pointer(tmp), -c_w))
		case 2:
			for x = 0; x < c_w; x++ {
				for y = 0; y < (func() int {
					if c_h > 2 {
						return 2
					}
					return c_h
				}()); y++ {
					if 0 < (func() int {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), 0)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y - 1) < 0 {
								return 0
							}
							return y - 1
						}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 2) > (c_h - 1) {
								return c_h - 1
							}
							return y + 2
						}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 3) > (c_h - 1) {
								return c_h - 1
							}
							return y + 3
						}())*c_w))) + 64) >> 7) > math.MaxUint8 {
							return math.MaxUint8
						}
						return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), 0)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y - 1) < 0 {
								return 0
							}
							return y - 1
						}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 2) > (c_h - 1) {
								return c_h - 1
							}
							return y + 2
						}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 3) > (c_h - 1) {
								return c_h - 1
							}
							return y + 3
						}())*c_w))) + 64) >> 7
					}()) {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), 0)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y - 1) < 0 {
								return 0
							}
							return y - 1
						}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 2) > (c_h - 1) {
								return c_h - 1
							}
							return y + 2
						}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 3) > (c_h - 1) {
								return c_h - 1
							}
							return y + 3
						}())*c_w))) + 64) >> 7) > math.MaxUint8 {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = math.MaxUint8
						} else {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), 0)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y - 1) < 0 {
									return 0
								}
								return y - 1
							}())*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y + 1) > (c_h - 1) {
									return c_h - 1
								}
								return y + 1
							}())*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y + 2) > (c_h - 1) {
									return c_h - 1
								}
								return y + 2
							}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y + 3) > (c_h - 1) {
									return c_h - 1
								}
								return y + 3
							}())*c_w))) + 64) >> 7))
						}
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = 0
					}
				}
				for ; y < c_h-3; y++ {
					if 0 < (func() int {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+1)*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+3)*c_w))) + 64) >> 7) > math.MaxUint8 {
							return math.MaxUint8
						}
						return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+1)*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+3)*c_w))) + 64) >> 7
					}()) {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+1)*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+3)*c_w))) + 64) >> 7) > math.MaxUint8 {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = math.MaxUint8
						} else {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+1)*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+2)*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y+3)*c_w))) + 64) >> 7))
						}
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = 0
					}
				}
				for ; y < c_h; y++ {
					if 0 < (func() int {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 2) > (c_h - 1) {
								return c_h - 1
							}
							return y + 2
						}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (c_h-1)*c_w))) + 64) >> 7) > math.MaxUint8 {
							return math.MaxUint8
						}
						return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 2) > (c_h - 1) {
								return c_h - 1
							}
							return y + 2
						}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (c_h-1)*c_w))) + 64) >> 7
					}()) {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 1) > (c_h - 1) {
								return c_h - 1
							}
							return y + 1
						}())*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
							if (y + 2) > (c_h - 1) {
								return c_h - 1
							}
							return y + 2
						}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (c_h-1)*c_w))) + 64) >> 7) > math.MaxUint8 {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = math.MaxUint8
						} else {
							*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-2)*c_w)))*4 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (y-1)*c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), y*c_w)))*114 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y + 1) > (c_h - 1) {
									return c_h - 1
								}
								return y + 1
							}())*c_w)))*35 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (func() int {
								if (y + 2) > (c_h - 1) {
									return c_h - 1
								}
								return y + 2
							}())*c_w)))*9 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), (c_h-1)*c_w))) + 64) >> 7))
						}
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), y*c_w)) = 0
					}
				}
				_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), 1))
				tmp = (*uint8)(unsafe.Add(unsafe.Pointer(tmp), 1))
			}
		}
	}
}
func y4m_422jpeg_420jpeg_helper(_dst *uint8, _src *uint8, _c_w int, _c_h int) {
	var (
		y int
		x int
	)
	for x = 0; x < _c_w; x++ {
		for y = 0; y < (func() int {
			if _c_h > 2 {
				return 2
			}
			return _c_h
		}()); y += 2 {
			if 0 < (func() int {
				if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), 0)))*64 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if 1 > (_c_h - 1) {
						return _c_h - 1
					}
					return 1
				}())*_c_w)))*78 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if 2 > (_c_h - 1) {
						return _c_h - 1
					}
					return 2
				}())*_c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if 3 > (_c_h - 1) {
						return _c_h - 1
					}
					return 3
				}())*_c_w)))*3 + 64) >> 7) > math.MaxUint8 {
					return math.MaxUint8
				}
				return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), 0)))*64 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if 1 > (_c_h - 1) {
						return _c_h - 1
					}
					return 1
				}())*_c_w)))*78 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if 2 > (_c_h - 1) {
						return _c_h - 1
					}
					return 2
				}())*_c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if 3 > (_c_h - 1) {
						return _c_h - 1
					}
					return 3
				}())*_c_w)))*3 + 64) >> 7
			}()) {
				if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), 0)))*64 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if 1 > (_c_h - 1) {
						return _c_h - 1
					}
					return 1
				}())*_c_w)))*78 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if 2 > (_c_h - 1) {
						return _c_h - 1
					}
					return 2
				}())*_c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if 3 > (_c_h - 1) {
						return _c_h - 1
					}
					return 3
				}())*_c_w)))*3 + 64) >> 7) > math.MaxUint8 {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), (y>>1)*_c_w)) = math.MaxUint8
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), (y>>1)*_c_w)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), 0)))*64 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
						if 1 > (_c_h - 1) {
							return _c_h - 1
						}
						return 1
					}())*_c_w)))*78 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
						if 2 > (_c_h - 1) {
							return _c_h - 1
						}
						return 2
					}())*_c_w)))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
						if 3 > (_c_h - 1) {
							return _c_h - 1
						}
						return 3
					}())*_c_w)))*3 + 64) >> 7))
				}
			} else {
				*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), (y>>1)*_c_w)) = 0
			}
		}
		for ; y < _c_h-3; y += 2 {
			if 0 < (func() int {
				if (((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-2)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+3)*_c_w))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-1)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+2)*_c_w))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), y*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+1)*_c_w))))*78 + 64) >> 7) > math.MaxUint8 {
					return math.MaxUint8
				}
				return ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-2)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+3)*_c_w))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-1)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+2)*_c_w))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), y*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+1)*_c_w))))*78 + 64) >> 7
			}()) {
				if (((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-2)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+3)*_c_w))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-1)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+2)*_c_w))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), y*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+1)*_c_w))))*78 + 64) >> 7) > math.MaxUint8 {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), (y>>1)*_c_w)) = math.MaxUint8
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), (y>>1)*_c_w)) = uint8(int8(((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-2)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+3)*_c_w))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-1)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+2)*_c_w))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), y*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y+1)*_c_w))))*78 + 64) >> 7))
				}
			} else {
				*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), (y>>1)*_c_w)) = 0
			}
		}
		for ; y < _c_h; y += 2 {
			if 0 < (func() int {
				if (((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-2)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (_c_h-1)*_c_w))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-1)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if (y + 2) > (_c_h - 1) {
						return _c_h - 1
					}
					return y + 2
				}())*_c_w))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), y*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if (y + 1) > (_c_h - 1) {
						return _c_h - 1
					}
					return y + 1
				}())*_c_w))))*78 + 64) >> 7) > math.MaxUint8 {
					return math.MaxUint8
				}
				return ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-2)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (_c_h-1)*_c_w))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-1)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if (y + 2) > (_c_h - 1) {
						return _c_h - 1
					}
					return y + 2
				}())*_c_w))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), y*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if (y + 1) > (_c_h - 1) {
						return _c_h - 1
					}
					return y + 1
				}())*_c_w))))*78 + 64) >> 7
			}()) {
				if (((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-2)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (_c_h-1)*_c_w))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-1)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if (y + 2) > (_c_h - 1) {
						return _c_h - 1
					}
					return y + 2
				}())*_c_w))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), y*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
					if (y + 1) > (_c_h - 1) {
						return _c_h - 1
					}
					return y + 1
				}())*_c_w))))*78 + 64) >> 7) > math.MaxUint8 {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), (y>>1)*_c_w)) = math.MaxUint8
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), (y>>1)*_c_w)) = uint8(int8(((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-2)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (_c_h-1)*_c_w))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (y-1)*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
						if (y + 2) > (_c_h - 1) {
							return _c_h - 1
						}
						return y + 2
					}())*_c_w))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), y*_c_w)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_src), (func() int {
						if (y + 1) > (_c_h - 1) {
							return _c_h - 1
						}
						return y + 1
					}())*_c_w))))*78 + 64) >> 7))
				}
			} else {
				*(*uint8)(unsafe.Add(unsafe.Pointer(_dst), (y>>1)*_c_w)) = 0
			}
		}
		_src = (*uint8)(unsafe.Add(unsafe.Pointer(_src), 1))
		_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), 1))
	}
}
func y4m_convert_422jpeg_420jpeg(_y4m *Y4mInput, _dst *uint8, _aux *uint8) {
	var (
		c_w      int
		c_h      int
		c_sz     int
		dst_c_w  int
		dst_c_h  int
		dst_c_sz int
		pli      int
	)
	_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), _y4m.Pic_w*_y4m.Pic_h))
	c_w = (_y4m.Pic_w + _y4m.Src_c_dec_h - 1) / _y4m.Src_c_dec_h
	c_h = _y4m.Pic_h
	dst_c_w = (_y4m.Pic_w + _y4m.Dst_c_dec_h - 1) / _y4m.Dst_c_dec_h
	dst_c_h = (_y4m.Pic_h + _y4m.Dst_c_dec_v - 1) / _y4m.Dst_c_dec_v
	c_sz = c_w * c_h
	dst_c_sz = dst_c_w * dst_c_h
	for pli = 1; pli < 3; pli++ {
		y4m_422jpeg_420jpeg_helper(_dst, _aux, c_w, c_h)
		_aux = (*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_sz))
		_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), dst_c_sz))
	}
}
func y4m_convert_422_420jpeg(_y4m *Y4mInput, _dst *uint8, _aux *uint8) {
	var (
		tmp      *uint8
		c_w      int
		c_h      int
		c_sz     int
		dst_c_h  int
		dst_c_sz int
		pli      int
	)
	_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), _y4m.Pic_w*_y4m.Pic_h))
	c_w = (_y4m.Pic_w + _y4m.Src_c_dec_h - 1) / _y4m.Src_c_dec_h
	c_h = _y4m.Pic_h
	dst_c_h = (_y4m.Pic_h + _y4m.Dst_c_dec_v - 1) / _y4m.Dst_c_dec_v
	c_sz = c_w * c_h
	dst_c_sz = c_w * dst_c_h
	tmp = (*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_sz*2))
	for pli = 1; pli < 3; pli++ {
		y4m_42xmpeg2_42xjpeg_helper(tmp, _aux, c_w, c_h)
		y4m_422jpeg_420jpeg_helper(_dst, tmp, c_w, c_h)
		_aux = (*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_sz))
		_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), dst_c_sz))
	}
}
func y4m_convert_411_420jpeg(_y4m *Y4mInput, _dst *uint8, _aux *uint8) {
	var (
		tmp      *uint8
		c_w      int
		c_h      int
		c_sz     int
		dst_c_w  int
		dst_c_h  int
		dst_c_sz int
		tmp_sz   int
		pli      int
		y        int
		x        int
	)
	_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), _y4m.Pic_w*_y4m.Pic_h))
	c_w = (_y4m.Pic_w + _y4m.Src_c_dec_h - 1) / _y4m.Src_c_dec_h
	c_h = _y4m.Pic_h
	dst_c_w = (_y4m.Pic_w + _y4m.Dst_c_dec_h - 1) / _y4m.Dst_c_dec_h
	dst_c_h = (_y4m.Pic_h + _y4m.Dst_c_dec_v - 1) / _y4m.Dst_c_dec_v
	c_sz = c_w * c_h
	dst_c_sz = dst_c_w * dst_c_h
	tmp_sz = dst_c_w * c_h
	tmp = (*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_sz*2))
	for pli = 1; pli < 3; pli++ {
		for y = 0; y < c_h; y++ {
			for x = 0; x < (func() int {
				if c_w > 1 {
					return 1
				}
				return c_w
			}()); x++ {
				if 0 < (func() int {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*111 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 1 > (c_w - 1) {
							return c_w - 1
						}
						return 1
					}())))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 2 > (c_w - 1) {
							return c_w - 1
						}
						return 2
					}()))) + 64) >> 7) > math.MaxUint8 {
						return math.MaxUint8
					}
					return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*111 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 1 > (c_w - 1) {
							return c_w - 1
						}
						return 1
					}())))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 2 > (c_w - 1) {
							return c_w - 1
						}
						return 2
					}()))) + 64) >> 7
				}()) {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*111 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 1 > (c_w - 1) {
							return c_w - 1
						}
						return 1
					}())))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 2 > (c_w - 1) {
							return c_w - 1
						}
						return 2
					}()))) + 64) >> 7) > math.MaxUint8 {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1)) = math.MaxUint8
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*111 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if 1 > (c_w - 1) {
								return c_w - 1
							}
							return 1
						}())))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if 2 > (c_w - 1) {
								return c_w - 1
							}
							return 2
						}()))) + 64) >> 7))
					}
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1)) = 0
				}
				if 0 < (func() int {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*47 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 1 > (c_w - 1) {
							return c_w - 1
						}
						return 1
					}())))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 2 > (c_w - 1) {
							return c_w - 1
						}
						return 2
					}())))*5 + 64) >> 7) > math.MaxUint8 {
						return math.MaxUint8
					}
					return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*47 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 1 > (c_w - 1) {
							return c_w - 1
						}
						return 1
					}())))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 2 > (c_w - 1) {
							return c_w - 1
						}
						return 2
					}())))*5 + 64) >> 7
				}()) {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*47 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 1 > (c_w - 1) {
							return c_w - 1
						}
						return 1
					}())))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 2 > (c_w - 1) {
							return c_w - 1
						}
						return 2
					}())))*5 + 64) >> 7) > math.MaxUint8 {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1|1)) = math.MaxUint8
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1|1)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*47 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if 1 > (c_w - 1) {
								return c_w - 1
							}
							return 1
						}())))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if 2 > (c_w - 1) {
								return c_w - 1
							}
							return 2
						}())))*5 + 64) >> 7))
					}
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1|1)) = 0
				}
			}
			for ; x < c_w-2; x++ {
				if 0 < (func() int {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*110 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1)))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2))) + 64) >> 7) > math.MaxUint8 {
						return math.MaxUint8
					}
					return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*110 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1)))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2))) + 64) >> 7
				}()) {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*110 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1)))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2))) + 64) >> 7) > math.MaxUint8 {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1)) = math.MaxUint8
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*110 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1)))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2))) + 64) >> 7))
					}
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1)) = 0
				}
				if 0 < (func() int {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))*(-3) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*50 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1)))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2)))*5 + 64) >> 7) > math.MaxUint8 {
						return math.MaxUint8
					}
					return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))*(-3) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*50 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1)))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2)))*5 + 64) >> 7
				}()) {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))*(-3) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*50 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1)))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2)))*5 + 64) >> 7) > math.MaxUint8 {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1|1)) = math.MaxUint8
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1|1)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))*(-3) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*50 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1)))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2)))*5 + 64) >> 7))
					}
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1|1)) = 0
				}
			}
			for ; x < c_w; x++ {
				if 0 < (func() int {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*110 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if (x + 1) > (c_w - 1) {
							return c_w - 1
						}
						return x + 1
					}())))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1))) + 64) >> 7) > math.MaxUint8 {
						return math.MaxUint8
					}
					return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*110 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if (x + 1) > (c_w - 1) {
							return c_w - 1
						}
						return x + 1
					}())))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1))) + 64) >> 7
				}()) {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*110 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if (x + 1) > (c_w - 1) {
							return c_w - 1
						}
						return x + 1
					}())))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1))) + 64) >> 7) > math.MaxUint8 {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1)) = math.MaxUint8
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1))) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*110 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if (x + 1) > (c_w - 1) {
								return c_w - 1
							}
							return x + 1
						}())))*18 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1))) + 64) >> 7))
					}
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1)) = 0
				}
				if (x<<1 | 1) < dst_c_w {
					if 0 < (func() int {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))*(-3) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*50 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if (x + 1) > (c_w - 1) {
								return c_w - 1
							}
							return x + 1
						}())))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1)))*5 + 64) >> 7) > math.MaxUint8 {
							return math.MaxUint8
						}
						return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))*(-3) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*50 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if (x + 1) > (c_w - 1) {
								return c_w - 1
							}
							return x + 1
						}())))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1)))*5 + 64) >> 7
					}()) {
						if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))*(-3) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*50 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if (x + 1) > (c_w - 1) {
								return c_w - 1
							}
							return x + 1
						}())))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1)))*5 + 64) >> 7) > math.MaxUint8 {
							*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1|1)) = math.MaxUint8
						} else {
							*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1|1)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))*(-3) + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))*50 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
								if (x + 1) > (c_w - 1) {
									return c_w - 1
								}
								return x + 1
							}())))*86 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1)))*5 + 64) >> 7))
						}
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x<<1|1)) = 0
					}
				}
			}
			tmp = (*uint8)(unsafe.Add(unsafe.Pointer(tmp), dst_c_w))
			_aux = (*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w))
		}
		tmp = (*uint8)(unsafe.Add(unsafe.Pointer(tmp), -tmp_sz))
		y4m_422jpeg_420jpeg_helper(_dst, tmp, dst_c_w, c_h)
		_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), dst_c_sz))
	}
}
func y4m_convert_444_420jpeg(_y4m *Y4mInput, _dst *uint8, _aux *uint8) {
	var (
		tmp      *uint8
		c_w      int
		c_h      int
		c_sz     int
		dst_c_w  int
		dst_c_h  int
		dst_c_sz int
		tmp_sz   int
		pli      int
		y        int
		x        int
	)
	_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), _y4m.Pic_w*_y4m.Pic_h))
	c_w = (_y4m.Pic_w + _y4m.Src_c_dec_h - 1) / _y4m.Src_c_dec_h
	c_h = _y4m.Pic_h
	dst_c_w = (_y4m.Pic_w + _y4m.Dst_c_dec_h - 1) / _y4m.Dst_c_dec_h
	dst_c_h = (_y4m.Pic_h + _y4m.Dst_c_dec_v - 1) / _y4m.Dst_c_dec_v
	c_sz = c_w * c_h
	dst_c_sz = dst_c_w * dst_c_h
	tmp_sz = dst_c_w * c_h
	tmp = (*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_sz*2))
	for pli = 1; pli < 3; pli++ {
		for y = 0; y < c_h; y++ {
			for x = 0; x < (func() int {
				if c_w > 2 {
					return 2
				}
				return c_w
			}()); x += 2 {
				if 0 < (func() int {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*64 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 1 > (c_w - 1) {
							return c_w - 1
						}
						return 1
					}())))*78 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 2 > (c_w - 1) {
							return c_w - 1
						}
						return 2
					}())))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 3 > (c_w - 1) {
							return c_w - 1
						}
						return 3
					}())))*3 + 64) >> 7) > math.MaxUint8 {
						return math.MaxUint8
					}
					return (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*64 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 1 > (c_w - 1) {
							return c_w - 1
						}
						return 1
					}())))*78 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 2 > (c_w - 1) {
							return c_w - 1
						}
						return 2
					}())))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 3 > (c_w - 1) {
							return c_w - 1
						}
						return 3
					}())))*3 + 64) >> 7
				}()) {
					if ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*64 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 1 > (c_w - 1) {
							return c_w - 1
						}
						return 1
					}())))*78 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 2 > (c_w - 1) {
							return c_w - 1
						}
						return 2
					}())))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if 3 > (c_w - 1) {
							return c_w - 1
						}
						return 3
					}())))*3 + 64) >> 7) > math.MaxUint8 {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x>>1)) = math.MaxUint8
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x>>1)) = uint8(int8((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), 0)))*64 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if 1 > (c_w - 1) {
								return c_w - 1
							}
							return 1
						}())))*78 - int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if 2 > (c_w - 1) {
								return c_w - 1
							}
							return 2
						}())))*17 + int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if 3 > (c_w - 1) {
								return c_w - 1
							}
							return 3
						}())))*3 + 64) >> 7))
					}
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x>>1)) = 0
				}
			}
			for ; x < c_w-3; x += 2 {
				if 0 < (func() int {
					if (((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-2)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+3))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1))))*78 + 64) >> 7) > math.MaxUint8 {
						return math.MaxUint8
					}
					return ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-2)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+3))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1))))*78 + 64) >> 7
				}()) {
					if (((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-2)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+3))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1))))*78 + 64) >> 7) > math.MaxUint8 {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x>>1)) = math.MaxUint8
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x>>1)) = uint8(int8(((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-2)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+3))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+2))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x+1))))*78 + 64) >> 7))
					}
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x>>1)) = 0
				}
			}
			for ; x < c_w; x += 2 {
				if 0 < (func() int {
					if (((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-2)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if (x + 2) > (c_w - 1) {
							return c_w - 1
						}
						return x + 2
					}()))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if (x + 1) > (c_w - 1) {
							return c_w - 1
						}
						return x + 1
					}()))))*78 + 64) >> 7) > math.MaxUint8 {
						return math.MaxUint8
					}
					return ((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-2)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if (x + 2) > (c_w - 1) {
							return c_w - 1
						}
						return x + 2
					}()))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if (x + 1) > (c_w - 1) {
							return c_w - 1
						}
						return x + 1
					}()))))*78 + 64) >> 7
				}()) {
					if (((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-2)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if (x + 2) > (c_w - 1) {
							return c_w - 1
						}
						return x + 2
					}()))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
						if (x + 1) > (c_w - 1) {
							return c_w - 1
						}
						return x + 1
					}()))))*78 + 64) >> 7) > math.MaxUint8 {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x>>1)) = math.MaxUint8
					} else {
						*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x>>1)) = uint8(int8(((int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-2)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w-1))))*3 - (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x-1)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if (x + 2) > (c_w - 1) {
								return c_w - 1
							}
							return x + 2
						}()))))*17 + (int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), x)))+int(*(*uint8)(unsafe.Add(unsafe.Pointer(_aux), func() int {
							if (x + 1) > (c_w - 1) {
								return c_w - 1
							}
							return x + 1
						}()))))*78 + 64) >> 7))
					}
				} else {
					*(*uint8)(unsafe.Add(unsafe.Pointer(tmp), x>>1)) = 0
				}
			}
			tmp = (*uint8)(unsafe.Add(unsafe.Pointer(tmp), dst_c_w))
			_aux = (*uint8)(unsafe.Add(unsafe.Pointer(_aux), c_w))
		}
		tmp = (*uint8)(unsafe.Add(unsafe.Pointer(tmp), -tmp_sz))
		y4m_422jpeg_420jpeg_helper(_dst, tmp, dst_c_w, c_h)
		_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), dst_c_sz))
	}
}
func y4m_convert_mono_420jpeg(_y4m *Y4mInput, _dst *uint8, _aux *uint8) {
	var c_sz int
	_ = _aux
	_dst = (*uint8)(unsafe.Add(unsafe.Pointer(_dst), _y4m.Pic_w*_y4m.Pic_h))
	c_sz = ((_y4m.Pic_w + _y4m.Dst_c_dec_h - 1) / _y4m.Dst_c_dec_h) * ((_y4m.Pic_h + _y4m.Dst_c_dec_v - 1) / _y4m.Dst_c_dec_v)
	libc.MemSet(unsafe.Pointer(_dst), 128, c_sz*2)
}
func y4m_convert_null(_y4m *Y4mInput, _dst *uint8, _aux *uint8) {
	_ = _y4m
	_ = _dst
	_ = _aux
}

var TAG [10]byte = func() [10]byte {
	var t [10]byte
	copy(t[:], ([]byte)("YUV4MPEG2"))
	return t
}()

func y4m_input_open(y4m_ctx *Y4mInput, file *stdio.File, skip_buffer *byte, num_skip int, only_420 int) int {
	var tag_buffer [9]byte
	libc.Assert(num_skip >= 0 && num_skip <= 8)
	if num_skip > 0 {
		libc.MemCpy(unsafe.Pointer(&tag_buffer[0]), unsafe.Pointer(skip_buffer), num_skip)
	}
	if file_read(unsafe.Pointer(&tag_buffer[num_skip]), uint64(9-num_skip), file) == 0 {
		return -1
	}
	if libc.MemCmp(unsafe.Pointer(&TAG[0]), unsafe.Pointer(&tag_buffer[0]), 9) != 0 {
		stdio.Fprintf(stdio.Stderr(), "Error parsing header: must start with %s\n", &TAG[0])
		return -1
	}
	if file_read(unsafe.Pointer(&tag_buffer[0]), 1, file) == 0 || tag_buffer[0] != byte(' ') {
		stdio.Fprintf(stdio.Stderr(), "Error parsing header: space must follow %s\n", &TAG[0])
		return -1
	}
	if parse_tags(y4m_ctx, file) == 0 {
		stdio.Fprintf(stdio.Stderr(), "Error parsing %s header.\n", &TAG[0])
	}
	if int(y4m_ctx.Interlace) == int('?') {
		stdio.Fprintf(stdio.Stderr(), "Warning: Input video interlacing format unknown; assuming progressive scan.\n")
	} else if int(y4m_ctx.Interlace) != int('p') {
		stdio.Fprintf(stdio.Stderr(), "Input video is interlaced; Only progressive scan handled.\n")
		return -1
	}
	y4m_ctx.Vpx_fmt = vpx.ImgFmt(VPX_IMG_FMT_I420)
	y4m_ctx.Bps = 12
	y4m_ctx.Bit_depth = 8
	y4m_ctx.Aux_buf = nil
	y4m_ctx.Dst_buf = nil
	if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("420")) == 0 || libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("420jpeg")) == 0 || libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("420mpeg2")) == 0 {
		y4m_ctx.Src_c_dec_h = func() int {
			p := &y4m_ctx.Dst_c_dec_h
			y4m_ctx.Dst_c_dec_h = func() int {
				p := &y4m_ctx.Src_c_dec_v
				y4m_ctx.Src_c_dec_v = func() int {
					p := &y4m_ctx.Dst_c_dec_v
					y4m_ctx.Dst_c_dec_v = 2
					return *p
				}()
				return *p
			}()
			return *p
		}()
		y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w*y4m_ctx.Pic_h + ((y4m_ctx.Pic_w+1)/2)*2*((y4m_ctx.Pic_h+1)/2))
		y4m_ctx.Aux_buf_sz = func() uint64 {
			p := &y4m_ctx.Aux_buf_read_sz
			y4m_ctx.Aux_buf_read_sz = 0
			return *p
		}()
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_null(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("420p10")) == 0 {
		y4m_ctx.Src_c_dec_h = 2
		y4m_ctx.Dst_c_dec_h = 2
		y4m_ctx.Src_c_dec_v = 2
		y4m_ctx.Dst_c_dec_v = 2
		y4m_ctx.Dst_buf_read_sz = uint64((y4m_ctx.Pic_w*y4m_ctx.Pic_h + ((y4m_ctx.Pic_w+1)/2)*2*((y4m_ctx.Pic_h+1)/2)) * 2)
		y4m_ctx.Aux_buf_sz = func() uint64 {
			p := &y4m_ctx.Aux_buf_read_sz
			y4m_ctx.Aux_buf_read_sz = 0
			return *p
		}()
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_null(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
		y4m_ctx.Bit_depth = 10
		y4m_ctx.Bps = 15
		y4m_ctx.Vpx_fmt = vpx.ImgFmt(VPX_IMG_FMT_I42016)
		if only_420 != 0 {
			stdio.Fprintf(stdio.Stderr(), "Unsupported conversion from 420p10 to 420jpeg\n")
			return -1
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("420p12")) == 0 {
		y4m_ctx.Src_c_dec_h = 2
		y4m_ctx.Dst_c_dec_h = 2
		y4m_ctx.Src_c_dec_v = 2
		y4m_ctx.Dst_c_dec_v = 2
		y4m_ctx.Dst_buf_read_sz = uint64((y4m_ctx.Pic_w*y4m_ctx.Pic_h + ((y4m_ctx.Pic_w+1)/2)*2*((y4m_ctx.Pic_h+1)/2)) * 2)
		y4m_ctx.Aux_buf_sz = func() uint64 {
			p := &y4m_ctx.Aux_buf_read_sz
			y4m_ctx.Aux_buf_read_sz = 0
			return *p
		}()
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_null(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
		y4m_ctx.Bit_depth = 12
		y4m_ctx.Bps = 18
		y4m_ctx.Vpx_fmt = vpx.ImgFmt(VPX_IMG_FMT_I42016)
		if only_420 != 0 {
			stdio.Fprintf(stdio.Stderr(), "Unsupported conversion from 420p12 to 420jpeg\n")
			return -1
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("420paldv")) == 0 {
		y4m_ctx.Src_c_dec_h = func() int {
			p := &y4m_ctx.Dst_c_dec_h
			y4m_ctx.Dst_c_dec_h = func() int {
				p := &y4m_ctx.Src_c_dec_v
				y4m_ctx.Src_c_dec_v = func() int {
					p := &y4m_ctx.Dst_c_dec_v
					y4m_ctx.Dst_c_dec_v = 2
					return *p
				}()
				return *p
			}()
			return *p
		}()
		y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w * y4m_ctx.Pic_h)
		y4m_ctx.Aux_buf_sz = uint64(((y4m_ctx.Pic_w + 1) / 2) * 3 * ((y4m_ctx.Pic_h + 1) / 2))
		y4m_ctx.Aux_buf_read_sz = uint64(((y4m_ctx.Pic_w + 1) / 2) * 2 * ((y4m_ctx.Pic_h + 1) / 2))
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_42xpaldv_42xjpeg(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("422jpeg")) == 0 {
		y4m_ctx.Src_c_dec_h = func() int {
			p := &y4m_ctx.Dst_c_dec_h
			y4m_ctx.Dst_c_dec_h = 2
			return *p
		}()
		y4m_ctx.Src_c_dec_v = 1
		y4m_ctx.Dst_c_dec_v = 2
		y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w * y4m_ctx.Pic_h)
		y4m_ctx.Aux_buf_sz = func() uint64 {
			p := &y4m_ctx.Aux_buf_read_sz
			y4m_ctx.Aux_buf_read_sz = uint64(((y4m_ctx.Pic_w + 1) / 2) * 2 * y4m_ctx.Pic_h)
			return *p
		}()
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_422jpeg_420jpeg(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("422")) == 0 {
		y4m_ctx.Src_c_dec_h = 2
		y4m_ctx.Src_c_dec_v = 1
		if only_420 != 0 {
			y4m_ctx.Dst_c_dec_h = 2
			y4m_ctx.Dst_c_dec_v = 2
			y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w * y4m_ctx.Pic_h)
			y4m_ctx.Aux_buf_read_sz = uint64(((y4m_ctx.Pic_w + 1) / 2) * 2 * y4m_ctx.Pic_h)
			y4m_ctx.Aux_buf_sz = y4m_ctx.Aux_buf_read_sz + uint64(((y4m_ctx.Pic_w+1)/2)*y4m_ctx.Pic_h)
			y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
					y4m_convert_422_420jpeg(_y4m, _dst, _src)
				}(_y4m, _dst, _src)
			}
		} else {
			y4m_ctx.Vpx_fmt = vpx.ImgFmt(VPX_IMG_FMT_I422)
			y4m_ctx.Bps = 16
			y4m_ctx.Dst_c_dec_h = y4m_ctx.Src_c_dec_h
			y4m_ctx.Dst_c_dec_v = y4m_ctx.Src_c_dec_v
			y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w*y4m_ctx.Pic_h + ((y4m_ctx.Pic_w+1)/2)*2*y4m_ctx.Pic_h)
			y4m_ctx.Aux_buf_sz = func() uint64 {
				p := &y4m_ctx.Aux_buf_read_sz
				y4m_ctx.Aux_buf_read_sz = 0
				return *p
			}()
			y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
					y4m_convert_null(_y4m, _dst, _src)
				}(_y4m, _dst, _src)
			}
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("422p10")) == 0 {
		y4m_ctx.Src_c_dec_h = 2
		y4m_ctx.Src_c_dec_v = 1
		y4m_ctx.Vpx_fmt = vpx.ImgFmt(VPX_IMG_FMT_I42216)
		y4m_ctx.Bps = 20
		y4m_ctx.Bit_depth = 10
		y4m_ctx.Dst_c_dec_h = y4m_ctx.Src_c_dec_h
		y4m_ctx.Dst_c_dec_v = y4m_ctx.Src_c_dec_v
		y4m_ctx.Dst_buf_read_sz = uint64((y4m_ctx.Pic_w*y4m_ctx.Pic_h + ((y4m_ctx.Pic_w+1)/2)*2*y4m_ctx.Pic_h) * 2)
		y4m_ctx.Aux_buf_sz = func() uint64 {
			p := &y4m_ctx.Aux_buf_read_sz
			y4m_ctx.Aux_buf_read_sz = 0
			return *p
		}()
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_null(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
		if only_420 != 0 {
			stdio.Fprintf(stdio.Stderr(), "Unsupported conversion from 422p10 to 420jpeg\n")
			return -1
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("422p12")) == 0 {
		y4m_ctx.Src_c_dec_h = 2
		y4m_ctx.Src_c_dec_v = 1
		y4m_ctx.Vpx_fmt = vpx.ImgFmt(VPX_IMG_FMT_I42216)
		y4m_ctx.Bps = 24
		y4m_ctx.Bit_depth = 12
		y4m_ctx.Dst_c_dec_h = y4m_ctx.Src_c_dec_h
		y4m_ctx.Dst_c_dec_v = y4m_ctx.Src_c_dec_v
		y4m_ctx.Dst_buf_read_sz = uint64((y4m_ctx.Pic_w*y4m_ctx.Pic_h + ((y4m_ctx.Pic_w+1)/2)*2*y4m_ctx.Pic_h) * 2)
		y4m_ctx.Aux_buf_sz = func() uint64 {
			p := &y4m_ctx.Aux_buf_read_sz
			y4m_ctx.Aux_buf_read_sz = 0
			return *p
		}()
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_null(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
		if only_420 != 0 {
			stdio.Fprintf(stdio.Stderr(), "Unsupported conversion from 422p12 to 420jpeg\n")
			return -1
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("411")) == 0 {
		y4m_ctx.Src_c_dec_h = 4
		y4m_ctx.Dst_c_dec_h = 2
		y4m_ctx.Src_c_dec_v = 1
		y4m_ctx.Dst_c_dec_v = 2
		y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w * y4m_ctx.Pic_h)
		y4m_ctx.Aux_buf_read_sz = uint64(((y4m_ctx.Pic_w + 3) / 4) * 2 * y4m_ctx.Pic_h)
		y4m_ctx.Aux_buf_sz = y4m_ctx.Aux_buf_read_sz + uint64(((y4m_ctx.Pic_w+1)/2)*y4m_ctx.Pic_h)
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_411_420jpeg(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
		stdio.Fprintf(stdio.Stderr(), "Unsupported conversion from yuv 411\n")
		return -1
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("444")) == 0 {
		y4m_ctx.Src_c_dec_h = 1
		y4m_ctx.Src_c_dec_v = 1
		if only_420 != 0 {
			y4m_ctx.Dst_c_dec_h = 2
			y4m_ctx.Dst_c_dec_v = 2
			y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w * y4m_ctx.Pic_h)
			y4m_ctx.Aux_buf_read_sz = uint64(y4m_ctx.Pic_w * 2 * y4m_ctx.Pic_h)
			y4m_ctx.Aux_buf_sz = y4m_ctx.Aux_buf_read_sz + uint64(((y4m_ctx.Pic_w+1)/2)*y4m_ctx.Pic_h)
			y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
					y4m_convert_444_420jpeg(_y4m, _dst, _src)
				}(_y4m, _dst, _src)
			}
		} else {
			y4m_ctx.Vpx_fmt = vpx.ImgFmt(VPX_IMG_FMT_I444)
			y4m_ctx.Bps = 24
			y4m_ctx.Dst_c_dec_h = y4m_ctx.Src_c_dec_h
			y4m_ctx.Dst_c_dec_v = y4m_ctx.Src_c_dec_v
			y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w * 3 * y4m_ctx.Pic_h)
			y4m_ctx.Aux_buf_sz = func() uint64 {
				p := &y4m_ctx.Aux_buf_read_sz
				y4m_ctx.Aux_buf_read_sz = 0
				return *p
			}()
			y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
					y4m_convert_null(_y4m, _dst, _src)
				}(_y4m, _dst, _src)
			}
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("444p10")) == 0 {
		y4m_ctx.Src_c_dec_h = 1
		y4m_ctx.Src_c_dec_v = 1
		y4m_ctx.Vpx_fmt = vpx.ImgFmt(VPX_IMG_FMT_I44416)
		y4m_ctx.Bps = 30
		y4m_ctx.Bit_depth = 10
		y4m_ctx.Dst_c_dec_h = y4m_ctx.Src_c_dec_h
		y4m_ctx.Dst_c_dec_v = y4m_ctx.Src_c_dec_v
		y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w * (2 * 3) * y4m_ctx.Pic_h)
		y4m_ctx.Aux_buf_sz = func() uint64 {
			p := &y4m_ctx.Aux_buf_read_sz
			y4m_ctx.Aux_buf_read_sz = 0
			return *p
		}()
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_null(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
		if only_420 != 0 {
			stdio.Fprintf(stdio.Stderr(), "Unsupported conversion from 444p10 to 420jpeg\n")
			return -1
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("444p12")) == 0 {
		y4m_ctx.Src_c_dec_h = 1
		y4m_ctx.Src_c_dec_v = 1
		y4m_ctx.Vpx_fmt = vpx.ImgFmt(VPX_IMG_FMT_I44416)
		y4m_ctx.Bps = 36
		y4m_ctx.Bit_depth = 12
		y4m_ctx.Dst_c_dec_h = y4m_ctx.Src_c_dec_h
		y4m_ctx.Dst_c_dec_v = y4m_ctx.Src_c_dec_v
		y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w * (2 * 3) * y4m_ctx.Pic_h)
		y4m_ctx.Aux_buf_sz = func() uint64 {
			p := &y4m_ctx.Aux_buf_read_sz
			y4m_ctx.Aux_buf_read_sz = 0
			return *p
		}()
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_null(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
		if only_420 != 0 {
			stdio.Fprintf(stdio.Stderr(), "Unsupported conversion from 444p12 to 420jpeg\n")
			return -1
		}
	} else if libc.StrCmp(&y4m_ctx.Chroma_type[0], libc.CString("mono")) == 0 {
		y4m_ctx.Src_c_dec_h = func() int {
			p := &y4m_ctx.Src_c_dec_v
			y4m_ctx.Src_c_dec_v = 0
			return *p
		}()
		y4m_ctx.Dst_c_dec_h = func() int {
			p := &y4m_ctx.Dst_c_dec_v
			y4m_ctx.Dst_c_dec_v = 2
			return *p
		}()
		y4m_ctx.Dst_buf_read_sz = uint64(y4m_ctx.Pic_w * y4m_ctx.Pic_h)
		y4m_ctx.Aux_buf_sz = func() uint64 {
			p := &y4m_ctx.Aux_buf_read_sz
			y4m_ctx.Aux_buf_read_sz = 0
			return *p
		}()
		y4m_ctx.Convert = func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
			func(_y4m *Y4mInput, _dst *uint8, _src *uint8) {
				y4m_convert_mono_420jpeg(_y4m, _dst, _src)
			}(_y4m, _dst, _src)
		}
	} else {
		stdio.Fprintf(stdio.Stderr(), "Unknown chroma sampling type: %s\n", &y4m_ctx.Chroma_type[0])
		return -1
	}
	y4m_ctx.Dst_buf_sz = uint64(y4m_ctx.Pic_w*y4m_ctx.Pic_h + ((y4m_ctx.Pic_w+y4m_ctx.Dst_c_dec_h-1)/y4m_ctx.Dst_c_dec_h)*2*((y4m_ctx.Pic_h+y4m_ctx.Dst_c_dec_v-1)/y4m_ctx.Dst_c_dec_v))
	if y4m_ctx.Bit_depth == 8 {
		y4m_ctx.Dst_buf = (*uint8)(libc.Malloc(int(y4m_ctx.Dst_buf_sz)))
	} else {
		y4m_ctx.Dst_buf = (*uint8)(libc.Malloc(int(y4m_ctx.Dst_buf_sz * 2)))
	}
	if y4m_ctx.Aux_buf_sz > 0 {
		y4m_ctx.Aux_buf = (*uint8)(libc.Malloc(int(y4m_ctx.Aux_buf_sz)))
	}
	return 0
}
func y4m_input_close(_y4m *Y4mInput) {
	libc.Free(unsafe.Pointer(_y4m.Dst_buf))
	libc.Free(unsafe.Pointer(_y4m.Aux_buf))
}
func y4m_input_fetch_frame(_y4m *Y4mInput, _fin *stdio.File, _img *vpx.Image) int {
	var (
		frame            [6]byte
		pic_sz           int
		c_w              int
		c_h              int
		c_sz             int
		bytes_per_sample int
	)
	if _y4m.Bit_depth > 8 {
		bytes_per_sample = 2
	} else {
		bytes_per_sample = 1
	}
	if file_read(unsafe.Pointer(&frame[0]), 6, _fin) == 0 {
		return 0
	}
	if libc.MemCmp(unsafe.Pointer(&frame[0]), unsafe.Pointer(libc.CString("FRAME")), 5) != 0 {
		stdio.Fprintf(stdio.Stderr(), "Loss of framing in Y4M input data\n")
		return -1
	}
	if frame[5] != byte('\n') {
		var (
			c int8
			j int
		)
		for j = 0; j < 79 && file_read(unsafe.Pointer(&c), 1, _fin) != 0 && int(c) != int('\n'); j++ {
		}
		if j == 79 {
			stdio.Fprintf(stdio.Stderr(), "Error parsing Y4M frame header\n")
			return -1
		}
	}
	if file_read(unsafe.Pointer(_y4m.Dst_buf), _y4m.Dst_buf_read_sz, _fin) == 0 {
		stdio.Fprintf(stdio.Stderr(), "Error reading Y4M frame data.\n")
		return -1
	}
	if file_read(unsafe.Pointer(_y4m.Aux_buf), _y4m.Aux_buf_read_sz, _fin) == 0 {
		stdio.Fprintf(stdio.Stderr(), "Error reading Y4M frame data.\n")
		return -1
	}
	(_y4m.Convert)(_y4m, _y4m.Dst_buf, _y4m.Aux_buf)
	*_img = vpx.Image{}
	_img.Fmt = vpx.ImgFmt(_y4m.Vpx_fmt)
	_img.W = func() uint {
		p := &_img.D_w
		_img.D_w = uint(_y4m.Pic_w)
		return *p
	}()
	_img.H = func() uint {
		p := &_img.D_h
		_img.D_h = uint(_y4m.Pic_h)
		return *p
	}()
	_img.X_chroma_shift = uint(_y4m.Dst_c_dec_h >> 1)
	_img.Y_chroma_shift = uint(_y4m.Dst_c_dec_v >> 1)
	_img.Bps = _y4m.Bps
	pic_sz = _y4m.Pic_w * _y4m.Pic_h * bytes_per_sample
	c_w = (_y4m.Pic_w + _y4m.Dst_c_dec_h - 1) / _y4m.Dst_c_dec_h
	c_w *= bytes_per_sample
	c_h = (_y4m.Pic_h + _y4m.Dst_c_dec_v - 1) / _y4m.Dst_c_dec_v
	c_sz = c_w * c_h
	_img.Stride[VPX_PLANE_Y] = func() int {
		p := &_img.Stride[VPX_PLANE_ALPHA]
		_img.Stride[VPX_PLANE_ALPHA] = _y4m.Pic_w * bytes_per_sample
		return *p
	}()
	_img.Stride[VPX_PLANE_U] = func() int {
		p := &_img.Stride[VPX_PLANE_V]
		_img.Stride[VPX_PLANE_V] = c_w
		return *p
	}()
	_img.Planes[VPX_PLANE_Y] = _y4m.Dst_buf
	_img.Planes[VPX_PLANE_U] = (*uint8)(unsafe.Add(unsafe.Pointer(_y4m.Dst_buf), pic_sz))
	_img.Planes[VPX_PLANE_V] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_y4m.Dst_buf), pic_sz))), c_sz))
	_img.Planes[VPX_PLANE_ALPHA] = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer(_y4m.Dst_buf), pic_sz))), c_sz*2))
	return 1
}
