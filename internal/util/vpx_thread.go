package util

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/pthread"
	"sync"
	"unsafe"
)

const MAX_DECODE_THREADS = 8

type VPxWorkerStatus int

const (
	NOT_OK = VPxWorkerStatus(iota)
	OK
	WORK
)

type VPxWorkerHook func(unsafe.Pointer, unsafe.Pointer) int
type VPxWorkerImpl struct {
	Mutex_     pthread.Mutex
	Condition_ sync.Cond
	Thread_    *pthread.Thread
}
type VPxWorker struct {
	Impl_     *VPxWorkerImpl
	Status_   VPxWorkerStatus
	Hook      VPxWorkerHook
	Data1     unsafe.Pointer
	Data2     unsafe.Pointer
	Had_error int
}
type VPxWorkerInterface struct {
	Init    func(worker *VPxWorker)
	Reset   func(worker *VPxWorker) int
	Sync    func(worker *VPxWorker) int
	Launch  func(worker *VPxWorker)
	Execute func(worker *VPxWorker)
	End     func(worker *VPxWorker)
}

func thread_loop(ptr unsafe.Pointer) unsafe.Pointer {
	var (
		worker *VPxWorker = (*VPxWorker)(ptr)
		done   int        = 0
	)
	for done == 0 {
		worker.Impl_.Mutex_.CLock()
		for worker.Status_ == VPxWorkerStatus(OK) {
			worker.Impl_.Condition_.L = &worker.Impl_.Mutex_
			worker.Impl_.Condition_.Wait()
		}
		if worker.Status_ == VPxWorkerStatus(WORK) {
			execute(worker)
			worker.Status_ = VPxWorkerStatus(OK)
		} else if worker.Status_ == VPxWorkerStatus(NOT_OK) {
			done = 1
		}
		worker.Impl_.Condition_.Signal()
		worker.Impl_.Mutex_.CUnlock()
	}
	return nil
}
func change_state(worker *VPxWorker, new_status VPxWorkerStatus) {
	if worker.Impl_ == nil {
		return
	}
	worker.Impl_.Mutex_.CLock()
	if worker.Status_ >= VPxWorkerStatus(OK) {
		for worker.Status_ != VPxWorkerStatus(OK) {
			worker.Impl_.Condition_.L = &worker.Impl_.Mutex_
			worker.Impl_.Condition_.Wait()
		}
		if new_status != VPxWorkerStatus(OK) {
			worker.Status_ = new_status
			worker.Impl_.Condition_.Signal()
		}
	}
	worker.Impl_.Mutex_.CUnlock()
}
func vpxInit(worker *VPxWorker) {
	*worker = VPxWorker{}
	worker.Status_ = VPxWorkerStatus(NOT_OK)
}
func vpxSync(worker *VPxWorker) int {
	change_state(worker, VPxWorkerStatus(OK))
	libc.Assert(worker.Status_ <= VPxWorkerStatus(OK))
	return int(libc.BoolToInt(worker.Had_error == 0))
}
func reset(worker *VPxWorker) int {
	var ok int = 1
	worker.Had_error = 0
	if worker.Status_ < VPxWorkerStatus(OK) {
		worker.Impl_ = (*VPxWorkerImpl)(vpx_calloc(1, uint64(unsafe.Sizeof(VPxWorkerImpl{}))))
		if worker.Impl_ == nil {
			return 0
		}
		if int(worker.Impl_.Mutex_.Init(nil)) != 0 {
			goto Error
		}
		if int(pthread.CondInit(&worker.Impl_.Condition_, nil)) != 0 {
			worker.Impl_.Mutex_.Destroy()
			goto Error
		}
		worker.Impl_.Mutex_.CLock()
		ok = int(libc.BoolToInt(int(pthread.Create(&worker.Impl_.Thread_, nil, thread_loop, unsafe.Pointer(worker))) == 0))
		if ok != 0 {
			worker.Status_ = VPxWorkerStatus(OK)
		}
		worker.Impl_.Mutex_.CUnlock()
		if ok == 0 {
			worker.Impl_.Mutex_.Destroy()
			pthread.CondFree(&worker.Impl_.Condition_)
		Error:
			vpx_free(unsafe.Pointer(worker.Impl_))
			worker.Impl_ = nil
			return 0
		}
	} else if worker.Status_ > VPxWorkerStatus(OK) {
		ok = vpxSync(worker)
	}
	libc.Assert(ok == 0 || worker.Status_ == VPxWorkerStatus(OK))
	return ok
}
func execute(worker *VPxWorker) {
	if worker.Hook != nil {
		worker.Had_error |= int(libc.BoolToInt(worker.Hook(worker.Data1, worker.Data2) == 0))
	}
}
func launch(worker *VPxWorker) {
	change_state(worker, VPxWorkerStatus(WORK))
}
func end(worker *VPxWorker) {
	if worker.Impl_ != nil {
		change_state(worker, VPxWorkerStatus(NOT_OK))
		worker.Impl_.Thread_.Join(nil)
		worker.Impl_.Mutex_.Destroy()
		pthread.CondFree(&worker.Impl_.Condition_)
		vpx_free(unsafe.Pointer(worker.Impl_))
		worker.Impl_ = nil
	}
	libc.Assert(worker.Status_ == VPxWorkerStatus(NOT_OK))
}

var g_worker_interface VPxWorkerInterface = VPxWorkerInterface{Init: vpxInit, Reset: reset, Sync: vpxSync, Launch: launch, Execute: execute, End: end}

func vpx_set_worker_interface(winterface *VPxWorkerInterface) int {
	if winterface == nil || winterface.Init == nil || winterface.Reset == nil || winterface.Sync == nil || winterface.Launch == nil || winterface.Execute == nil || winterface.End == nil {
		return 0
	}
	g_worker_interface = *winterface
	return 1
}
func vpx_get_worker_interface() *VPxWorkerInterface {
	return &g_worker_interface
}
