package util

import (
	"github.com/gotranspile/cxgo/runtime/libc"
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
	Mutex_     pthread_mutex_t
	Condition_ pthread_cond_t
	Thread_    pthread_t
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
		pthread_mutex_lock(&worker.Impl_.Mutex_)
		for worker.Status_ == VPxWorkerStatus(OK) {
			pthread_cond_wait(&worker.Impl_.Condition_, &worker.Impl_.Mutex_)
		}
		if worker.Status_ == VPxWorkerStatus(WORK) {
			execute(worker)
			worker.Status_ = VPxWorkerStatus(OK)
		} else if worker.Status_ == VPxWorkerStatus(NOT_OK) {
			done = 1
		}
		pthread_cond_signal(&worker.Impl_.Condition_)
		pthread_mutex_unlock(&worker.Impl_.Mutex_)
	}
	return nil
}
func change_state(worker *VPxWorker, new_status VPxWorkerStatus) {
	if worker.Impl_ == nil {
		return
	}
	pthread_mutex_lock(&worker.Impl_.Mutex_)
	if worker.Status_ >= VPxWorkerStatus(OK) {
		for worker.Status_ != VPxWorkerStatus(OK) {
			pthread_cond_wait(&worker.Impl_.Condition_, &worker.Impl_.Mutex_)
		}
		if new_status != VPxWorkerStatus(OK) {
			worker.Status_ = new_status
			pthread_cond_signal(&worker.Impl_.Condition_)
		}
	}
	pthread_mutex_unlock(&worker.Impl_.Mutex_)
}
func init(worker *VPxWorker) {
	*worker = VPxWorker{}
	worker.Status_ = VPxWorkerStatus(NOT_OK)
}
func sync(worker *VPxWorker) int {
	change_state(worker, VPxWorkerStatus(OK))
	if worker.Status_ <= VPxWorkerStatus(OK) {
	} else {
		__assert_fail(libc.CString("worker->status_ <= OK"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
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
		if pthread_mutex_init(&worker.Impl_.Mutex_, nil) != 0 {
			goto Error
		}
		if pthread_cond_init(&worker.Impl_.Condition_, nil) != 0 {
			pthread_mutex_destroy(&worker.Impl_.Mutex_)
			goto Error
		}
		pthread_mutex_lock(&worker.Impl_.Mutex_)
		ok = int(libc.BoolToInt(pthread_create(&worker.Impl_.Thread_, nil, thread_loop, unsafe.Pointer(worker)) == 0))
		if ok != 0 {
			worker.Status_ = VPxWorkerStatus(OK)
		}
		pthread_mutex_unlock(&worker.Impl_.Mutex_)
		if ok == 0 {
			pthread_mutex_destroy(&worker.Impl_.Mutex_)
			pthread_cond_destroy(&worker.Impl_.Condition_)
		Error:
			vpx_free(unsafe.Pointer(worker.Impl_))
			worker.Impl_ = nil
			return 0
		}
	} else if worker.Status_ > VPxWorkerStatus(OK) {
		ok = sync(worker)
	}
	if ok == 0 || worker.Status_ == VPxWorkerStatus(OK) {
	} else {
		__assert_fail(libc.CString("!ok || (worker->status_ == OK)"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
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
		pthread_join(worker.Impl_.Thread_, nil)
		pthread_mutex_destroy(&worker.Impl_.Mutex_)
		pthread_cond_destroy(&worker.Impl_.Condition_)
		vpx_free(unsafe.Pointer(worker.Impl_))
		worker.Impl_ = nil
	}
	if worker.Status_ == VPxWorkerStatus(NOT_OK) {
	} else {
		__assert_fail(libc.CString("worker->status_ == NOT_OK"), libc.CString(__FILE__), __LINE__, (*byte)(nil))
	}
}

var g_worker_interface VPxWorkerInterface = VPxWorkerInterface{Init: init, Reset: reset, Sync: sync, Launch: launch, Execute: execute, End: end}

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
