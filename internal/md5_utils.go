package internal

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

type MD5Context struct {
	Buf   [4]uint
	Bytes [2]uint
	In    [16]uint
}

func byteSwap(buf *uint, words uint) {
	var (
		p *uint8
		i int = 1
	)
	if *(*byte)(unsafe.Pointer(&i)) == 1 {
		return
	}
	p = (*uint8)(unsafe.Pointer(buf))
	for {
		*func() *uint {
			p := &buf
			x := *p
			*p = (*uint)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(uint(0))*1))
			return x
		}() = (uint(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 3)))<<8|uint(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 2))))<<16 | (uint(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 1)))<<8 | uint(*(*uint8)(unsafe.Add(unsafe.Pointer(p), 0))))
		p = (*uint8)(unsafe.Add(unsafe.Pointer(p), 4))
		if func() uint {
			p := &words
			*p--
			return *p
		}() == 0 {
			break
		}
	}
}
func MD5Init(ctx *MD5Context) {
	ctx.Buf[0] = 0x67452301
	ctx.Buf[1] = 0xEFCDAB89
	ctx.Buf[2] = 0x98BADCFE
	ctx.Buf[3] = 0x10325476
	ctx.Bytes[0] = 0
	ctx.Bytes[1] = 0
}
func MD5Update(ctx *MD5Context, buf *uint8, len_ uint) {
	var t uint
	t = ctx.Bytes[0]
	if (func() uint {
		p := &ctx.Bytes[0]
		ctx.Bytes[0] = t + len_
		return *p
	}()) < t {
		ctx.Bytes[1]++
	}
	t = 64 - (t & 63)
	if t > len_ {
		libc.MemCpy(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Pointer(&ctx.In[0]))), 64))), -int(t)), unsafe.Pointer(buf), int(len_))
		return
	}
	libc.MemCpy(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Pointer(&ctx.In[0]))), 64))), -int(t)), unsafe.Pointer(buf), int(t))
	byteSwap(&ctx.In[0], 16)
	MD5Transform(ctx.Buf, ctx.In)
	buf = (*uint8)(unsafe.Add(unsafe.Pointer(buf), t))
	len_ -= t
	for len_ >= 64 {
		libc.MemCpy(unsafe.Pointer(&ctx.In[0]), unsafe.Pointer(buf), 64)
		byteSwap(&ctx.In[0], 16)
		MD5Transform(ctx.Buf, ctx.In)
		buf = (*uint8)(unsafe.Add(unsafe.Pointer(buf), 64))
		len_ -= 64
	}
	libc.MemCpy(unsafe.Pointer(&ctx.In[0]), unsafe.Pointer(buf), int(len_))
}
func MD5Final(digest [16]uint8, ctx *MD5Context) {
	var (
		count int    = int(ctx.Bytes[0] & 63)
		p     *uint8 = (*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(unsafe.Pointer(&ctx.In[0]))), count))
	)
	*func() *uint8 {
		p := &p
		x := *p
		*p = (*uint8)(unsafe.Add(unsafe.Pointer(*p), 1))
		return x
	}() = 128
	count = 56 - 1 - count
	if count < 0 {
		libc.MemSet(unsafe.Pointer(p), 0, count+8)
		byteSwap(&ctx.In[0], 16)
		MD5Transform(ctx.Buf, ctx.In)
		p = (*uint8)(unsafe.Pointer(&ctx.In[0]))
		count = 56
	}
	libc.MemSet(unsafe.Pointer(p), 0, count)
	byteSwap(&ctx.In[0], 14)
	ctx.In[14] = ctx.Bytes[0] << 3
	ctx.In[15] = ctx.Bytes[1]<<3 | ctx.Bytes[0]>>29
	MD5Transform(ctx.Buf, ctx.In)
	byteSwap(&ctx.Buf[0], 4)
	libc.MemCpy(unsafe.Pointer(&digest[0]), unsafe.Pointer(&ctx.Buf[0]), 16)
	*ctx = MD5Context{}
}
func MD5Transform(buf [4]uint, in [16]uint) {
	var (
		a uint
		b uint
		c uint
		d uint
	)
	a = buf[0]
	b = buf[1]
	c = buf[2]
	d = buf[3]
	a += (d ^ (b & (c ^ d))) + in[0] + 0xD76AA478
	a = (a<<7 | a>>(32-7)) + b
	d += (c ^ (a & (b ^ c))) + in[1] + 0xE8C7B756
	d = (d<<12 | d>>(32-12)) + a
	c += (b ^ (d & (a ^ b))) + in[2] + 0x242070DB
	c = (c<<17 | c>>(32-17)) + d
	b += (a ^ (c & (d ^ a))) + in[3] + 0xC1BDCEEE
	b = (b<<22 | b>>(32-22)) + c
	a += (d ^ (b & (c ^ d))) + in[4] + 0xF57C0FAF
	a = (a<<7 | a>>(32-7)) + b
	d += (c ^ (a & (b ^ c))) + in[5] + 0x4787C62A
	d = (d<<12 | d>>(32-12)) + a
	c += (b ^ (d & (a ^ b))) + in[6] + 0xA8304613
	c = (c<<17 | c>>(32-17)) + d
	b += (a ^ (c & (d ^ a))) + in[7] + 0xFD469501
	b = (b<<22 | b>>(32-22)) + c
	a += (d ^ (b & (c ^ d))) + in[8] + 0x698098D8
	a = (a<<7 | a>>(32-7)) + b
	d += (c ^ (a & (b ^ c))) + in[9] + 0x8B44F7AF
	d = (d<<12 | d>>(32-12)) + a
	c += (b ^ (d & (a ^ b))) + in[10] + 0xFFFF5BB1
	c = (c<<17 | c>>(32-17)) + d
	b += (a ^ (c & (d ^ a))) + in[11] + 0x895CD7BE
	b = (b<<22 | b>>(32-22)) + c
	a += (d ^ (b & (c ^ d))) + in[12] + 0x6B901122
	a = (a<<7 | a>>(32-7)) + b
	d += (c ^ (a & (b ^ c))) + in[13] + 0xFD987193
	d = (d<<12 | d>>(32-12)) + a
	c += (b ^ (d & (a ^ b))) + in[14] + 0xA679438E
	c = (c<<17 | c>>(32-17)) + d
	b += (a ^ (c & (d ^ a))) + in[15] + 0x49B40821
	b = (b<<22 | b>>(32-22)) + c
	a += (c ^ (d & (b ^ c))) + in[1] + 0xF61E2562
	a = (a<<5 | a>>(32-5)) + b
	d += (b ^ (c & (a ^ b))) + in[6] + 0xC040B340
	d = (d<<9 | d>>(32-9)) + a
	c += (a ^ (b & (d ^ a))) + in[11] + 0x265E5A51
	c = (c<<14 | c>>(32-14)) + d
	b += (d ^ (a & (c ^ d))) + in[0] + 0xE9B6C7AA
	b = (b<<20 | b>>(32-20)) + c
	a += (c ^ (d & (b ^ c))) + in[5] + 0xD62F105D
	a = (a<<5 | a>>(32-5)) + b
	d += (b ^ (c & (a ^ b))) + in[10] + 0x2441453
	d = (d<<9 | d>>(32-9)) + a
	c += (a ^ (b & (d ^ a))) + in[15] + 0xD8A1E681
	c = (c<<14 | c>>(32-14)) + d
	b += (d ^ (a & (c ^ d))) + in[4] + 0xE7D3FBC8
	b = (b<<20 | b>>(32-20)) + c
	a += (c ^ (d & (b ^ c))) + in[9] + 0x21E1CDE6
	a = (a<<5 | a>>(32-5)) + b
	d += (b ^ (c & (a ^ b))) + in[14] + 0xC33707D6
	d = (d<<9 | d>>(32-9)) + a
	c += (a ^ (b & (d ^ a))) + in[3] + 0xF4D50D87
	c = (c<<14 | c>>(32-14)) + d
	b += (d ^ (a & (c ^ d))) + in[8] + 0x455A14ED
	b = (b<<20 | b>>(32-20)) + c
	a += (c ^ (d & (b ^ c))) + in[13] + 0xA9E3E905
	a = (a<<5 | a>>(32-5)) + b
	d += (b ^ (c & (a ^ b))) + in[2] + 0xFCEFA3F8
	d = (d<<9 | d>>(32-9)) + a
	c += (a ^ (b & (d ^ a))) + in[7] + 0x676F02D9
	c = (c<<14 | c>>(32-14)) + d
	b += (d ^ (a & (c ^ d))) + in[12] + 0x8D2A4C8A
	b = (b<<20 | b>>(32-20)) + c
	a += (b ^ c ^ d) + in[5] + 0xFFFA3942
	a = (a<<4 | a>>(32-4)) + b
	d += (a ^ b ^ c) + in[8] + 0x8771F681
	d = (d<<11 | d>>(32-11)) + a
	c += (d ^ a ^ b) + in[11] + 0x6D9D6122
	c = (c<<16 | c>>(32-16)) + d
	b += (c ^ d ^ a) + in[14] + 0xFDE5380C
	b = (b<<23 | b>>(32-23)) + c
	a += (b ^ c ^ d) + in[1] + 0xA4BEEA44
	a = (a<<4 | a>>(32-4)) + b
	d += (a ^ b ^ c) + in[4] + 0x4BDECFA9
	d = (d<<11 | d>>(32-11)) + a
	c += (d ^ a ^ b) + in[7] + 0xF6BB4B60
	c = (c<<16 | c>>(32-16)) + d
	b += (c ^ d ^ a) + in[10] + 0xBEBFBC70
	b = (b<<23 | b>>(32-23)) + c
	a += (b ^ c ^ d) + in[13] + 0x289B7EC6
	a = (a<<4 | a>>(32-4)) + b
	d += (a ^ b ^ c) + in[0] + 0xEAA127FA
	d = (d<<11 | d>>(32-11)) + a
	c += (d ^ a ^ b) + in[3] + 0xD4EF3085
	c = (c<<16 | c>>(32-16)) + d
	b += (c ^ d ^ a) + in[6] + 0x4881D05
	b = (b<<23 | b>>(32-23)) + c
	a += (b ^ c ^ d) + in[9] + 0xD9D4D039
	a = (a<<4 | a>>(32-4)) + b
	d += (a ^ b ^ c) + in[12] + 0xE6DB99E5
	d = (d<<11 | d>>(32-11)) + a
	c += (d ^ a ^ b) + in[15] + 530742520
	c = (c<<16 | c>>(32-16)) + d
	b += (c ^ d ^ a) + in[2] + 0xC4AC5665
	b = (b<<23 | b>>(32-23)) + c
	a += (c ^ (b | ^d)) + in[0] + 0xF4292244
	a = (a<<6 | a>>(32-6)) + b
	d += (b ^ (a | ^c)) + in[7] + 0x432AFF97
	d = (d<<10 | d>>(32-10)) + a
	c += (a ^ (d | ^b)) + in[14] + 0xAB9423A7
	c = (c<<15 | c>>(32-15)) + d
	b += (d ^ (c | ^a)) + in[5] + 0xFC93A039
	b = (b<<21 | b>>(32-21)) + c
	a += (c ^ (b | ^d)) + in[12] + 0x655B59C3
	a = (a<<6 | a>>(32-6)) + b
	d += (b ^ (a | ^c)) + in[3] + 0x8F0CCC92
	d = (d<<10 | d>>(32-10)) + a
	c += (a ^ (d | ^b)) + in[10] + 0xFFEFF47D
	c = (c<<15 | c>>(32-15)) + d
	b += (d ^ (c | ^a)) + in[1] + 0x85845DD1
	b = (b<<21 | b>>(32-21)) + c
	a += (c ^ (b | ^d)) + in[8] + 0x6FA87E4F
	a = (a<<6 | a>>(32-6)) + b
	d += (b ^ (a | ^c)) + in[15] + 0xFE2CE6E0
	d = (d<<10 | d>>(32-10)) + a
	c += (a ^ (d | ^b)) + in[6] + 0xA3014314
	c = (c<<15 | c>>(32-15)) + d
	b += (d ^ (c | ^a)) + in[13] + 0x4E0811A1
	b = (b<<21 | b>>(32-21)) + c
	a += (c ^ (b | ^d)) + in[4] + 0xF7537E82
	a = (a<<6 | a>>(32-6)) + b
	d += (b ^ (a | ^c)) + in[11] + 0xBD3AF235
	d = (d<<10 | d>>(32-10)) + a
	c += (a ^ (d | ^b)) + in[2] + 0x2AD7D2BB
	c = (c<<15 | c>>(32-15)) + d
	b += (d ^ (c | ^a)) + in[9] + 0xEB86D391
	b = (b<<21 | b>>(32-21)) + c
	buf[0] += a
	buf[1] += b
	buf[2] += c
	buf[3] += d
}
