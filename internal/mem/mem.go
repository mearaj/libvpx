package mem

import (
	"unsafe"

	"github.com/gotranspile/cxgo/runtime/libc"
)

const VPX_MAX_ALLOCABLE_MEMORY = 0x10000000000

func CheckSizeArgOverflow(nmemb uint64, size uint64) int {
	var total_size uint64 = nmemb * size
	if nmemb == 0 {
		return 1
	}
	if size > (1<<40)/nmemb {
		return 0
	}
	if total_size != total_size {
		return 0
	}
	return 1
}
func GetMallocAddrLocation(mem unsafe.Pointer) *uint64 {
	return (*uint64)(unsafe.Add(unsafe.Pointer((*uint64)(mem)), -int(unsafe.Sizeof(uint64(0))*1)))
}
func GetAlignedMallocSize(size uint64, align uint64) uint64 {
	return size + align - 1 + uint64(unsafe.Sizeof(uint64(0)))
}
func SetActualMallocAddr(mem unsafe.Pointer, malloc_addr unsafe.Pointer) {
	var malloc_addr_location *uint64 = GetMallocAddrLocation(mem)
	_ = malloc_addr_location
	*malloc_addr_location = uint64(uintptr(malloc_addr))
}
func GetActualMallocAddress(mem unsafe.Pointer) unsafe.Pointer {
	var malloc_addr_location *uint64 = GetMallocAddrLocation(mem)
	return unsafe.Pointer(uintptr(*malloc_addr_location))
}
func VpxMemAlign(align uint64, size uint64) unsafe.Pointer {
	var (
		x            unsafe.Pointer = nil
		addr         unsafe.Pointer
		aligned_size uint64 = GetAlignedMallocSize(size, align)
	)
	if CheckSizeArgOverflow(1, aligned_size) == 0 {
		return nil
	}
	addr = libc.Malloc(int(aligned_size))
	if addr != nil {
		x = unsafe.Pointer(uintptr((uint64(uintptr(unsafe.Pointer((*uint8)(unsafe.Add(unsafe.Pointer((*uint8)(addr)), unsafe.Sizeof(uint64(0))))))) + (align - 1)) & ^(align - 1)))
		SetActualMallocAddr(x, addr)
	}
	return x
}
