package util

import "sync/atomic"

type VpxAtomicInt struct {
	Value int64
}

func VpxAtomicInit(p *VpxAtomicInt, value int) {
	p.Value = int64(value)
}
func VpxAtomicStoreRelease(p *VpxAtomicInt, value int) {
	atomic.StoreInt64(&p.Value, int64(value))
}
func AtomicLoadAcquire(p *VpxAtomicInt) int {
	return int(atomic.LoadInt64(&p.Value))
}
