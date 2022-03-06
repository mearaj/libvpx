package ports

import "sync"

func Once(func_ func()) {
	var lock sync.Once = sync.Once{}
	lock.Do(func_)
}
