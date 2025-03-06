package hookx

func init() {
	go exitWait()
}

// add hooks to exitlist,
// when receive os.Interrupt or syscall.SIGTERM signal, run all hooks
func Exit(hooks ...func()) {
	exitListLock.Lock()
	defer exitListLock.Unlock()
	exitList = append(exitList, hooks...)
}

// run hooks and ensure only run once,
// suggest to use in main.init(),
// input: a list of *func()
func Init(hooks ...(*func())) {
	onceRun(hooks)
}
