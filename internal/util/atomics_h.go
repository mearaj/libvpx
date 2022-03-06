package util

type VpxAtomicInt struct {
	Value int
}

func VpxAtomicInit(atomic *VpxAtomicInt, value int) {
	atomic.Value = value
}
func VpxAtomicStoreRelease(atomic *VpxAtomicInt, value int) {
	asm
	atomic.Value = value
}
func AtomicLoadAcquire(atomic *VpxAtomicInt) int {
	var v int = atomic.Value
	asm
	return v
}
