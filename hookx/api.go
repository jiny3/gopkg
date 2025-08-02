package hookx

// add hooks to exitlist,
// when receive os.Interrupt or syscall.SIGTERM signal, run all hooks
// use ExitWait() to start a goroutine to wait for exit signal,
func Exit(opts ...func()) {
	exitListLock.Lock()
	defer exitListLock.Unlock()
	exitList = append(exitList, opts...)
}

// run hooks and ensure only run once,
// suggest to use in main.init(),
// input: a list of *func()
func Init(opts ...(*func())) {
	onceRun(opts)
}
