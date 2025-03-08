package hookx

func init() {
	// init signal
	go exitWait()
}

// add hooks to exitlist,
// when receive os.Interrupt or syscall.SIGTERM signal, run all hooks
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
