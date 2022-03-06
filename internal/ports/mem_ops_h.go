package ports

import (
	"math"
	"unsafe"
)

func GetBe16AsInt(vmem unsafe.Pointer) uint {
	var (
		val uint
		mem *uint8 = (*uint8)(vmem)
	)
	val = uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0))) << 8)
	val |= uint(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1)))
	return val
}
func GetBe24AsInt(vmem unsafe.Pointer) uint {
	var (
		val uint
		mem *uint8 = (*uint8)(vmem)
	)
	val = uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0))) << 16)
	val |= uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1))) << 8)
	val |= uint(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 2)))
	return val
}
func GetBe32AsInt(vmem unsafe.Pointer) uint {
	var (
		val uint
		mem *uint8 = (*uint8)(vmem)
	)
	val = (uint(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)))) << 24
	val |= uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1))) << 16)
	val |= uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 2))) << 8)
	val |= uint(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 3)))
	return val
}
func GetLe16AsInt(vmem unsafe.Pointer) uint {
	var (
		val uint
		mem *uint8 = (*uint8)(vmem)
	)
	val = uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1))) << 8)
	val |= uint(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)))
	return val
}
func GetLe24AsInt(vmem unsafe.Pointer) uint {
	var (
		val uint
		mem *uint8 = (*uint8)(vmem)
	)
	val = uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 2))) << 16)
	val |= uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1))) << 8)
	val |= uint(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)))
	return val
}
func GetLe32AsInt(vmem unsafe.Pointer) uint {
	var (
		val uint
		mem *uint8 = (*uint8)(vmem)
	)
	val = (uint(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 3)))) << 24
	val |= uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 2))) << 16)
	val |= uint(int(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1))) << 8)
	val |= uint(*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)))
	return val
}
func GetSbe16AsInt(vmem unsafe.Pointer) int {
	var (
		mem *uint8 = (*uint8)(vmem)
		val int    = int(GetBe16AsInt(unsafe.Pointer(mem)))
	)
	return (val << int((unsafe.Sizeof(int(0))<<3)-16)) >> int((unsafe.Sizeof(int(0))<<3)-16)
}
func GetSbe24AsInt(vmem unsafe.Pointer) int {
	var (
		mem *uint8 = (*uint8)(vmem)
		val int    = int(GetBe24AsInt(unsafe.Pointer(mem)))
	)
	return (val << int((unsafe.Sizeof(int(0))<<3)-24)) >> int((unsafe.Sizeof(int(0))<<3)-24)
}
func GetSbe32AsInt(vmem unsafe.Pointer) int {
	var (
		mem *uint8 = (*uint8)(vmem)
		val int    = int(GetBe32AsInt(unsafe.Pointer(mem)))
	)
	return (val << int((unsafe.Sizeof(int(0))<<3)-32)) >> int((unsafe.Sizeof(int(0))<<3)-32)
}
func GetSle16AsInt(vmem unsafe.Pointer) int {
	var (
		mem *uint8 = (*uint8)(vmem)
		val int    = int(GetLe16AsInt(unsafe.Pointer(mem)))
	)
	return (val << int((unsafe.Sizeof(int(0))<<3)-16)) >> int((unsafe.Sizeof(int(0))<<3)-16)
}
func GetSle24AsInt(vmem unsafe.Pointer) int {
	var (
		mem *uint8 = (*uint8)(vmem)
		val int    = int(GetLe24AsInt(unsafe.Pointer(mem)))
	)
	return (val << int((unsafe.Sizeof(int(0))<<3)-24)) >> int((unsafe.Sizeof(int(0))<<3)-24)
}
func GetSle32AsInt(vmem unsafe.Pointer) int {
	var (
		mem *uint8 = (*uint8)(vmem)
		val int    = int(GetLe32AsInt(unsafe.Pointer(mem)))
	)
	return (val << int((unsafe.Sizeof(int(0))<<3)-32)) >> int((unsafe.Sizeof(int(0))<<3)-32)
}
func PutBe16AsInt(vmem unsafe.Pointer, val int) {
	var mem *uint8 = (*uint8)(vmem)
	_ = mem
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)) = uint8(int8((val >> 8) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1)) = uint8(int8((val >> 0) & math.MaxUint8))
}
func PutBe24AsInt(vmem unsafe.Pointer, val int) {
	var mem *uint8 = (*uint8)(vmem)
	_ = mem
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)) = uint8(int8((val >> 16) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1)) = uint8(int8((val >> 8) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 2)) = uint8(int8((val >> 0) & math.MaxUint8))
}
func PutBe32AsInt(vmem unsafe.Pointer, val int) {
	var mem *uint8 = (*uint8)(vmem)
	_ = mem
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)) = uint8(int8((val >> 24) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1)) = uint8(int8((val >> 16) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 2)) = uint8(int8((val >> 8) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 3)) = uint8(int8((val >> 0) & math.MaxUint8))
}
func PutLe16AsInt(vmem unsafe.Pointer, val int) {
	var mem *uint8 = (*uint8)(vmem)
	_ = mem
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)) = uint8(int8((val >> 0) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1)) = uint8(int8((val >> 8) & math.MaxUint8))
}
func PutLe24AsInt(vmem unsafe.Pointer, val int) {
	var mem *uint8 = (*uint8)(vmem)
	_ = mem
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)) = uint8(int8((val >> 0) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1)) = uint8(int8((val >> 8) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 2)) = uint8(int8((val >> 16) & math.MaxUint8))
}
func PutLe32AsInt(vmem unsafe.Pointer, val int) {
	var mem *uint8 = (*uint8)(vmem)
	_ = mem
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 0)) = uint8(int8((val >> 0) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 1)) = uint8(int8((val >> 8) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 2)) = uint8(int8((val >> 16) & math.MaxUint8))
	*(*uint8)(unsafe.Add(unsafe.Pointer(mem), 3)) = uint8(int8((val >> 24) & math.MaxUint8))
}
