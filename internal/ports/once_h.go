package ports

func Once(func_ func()) {
	var lock pthread_once_t = PTHREAD_ONCE_INIT
	pthread_once(&lock, func_)
}
