package hookx

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	exitListLock sync.Mutex
	exitList     = []func(){}
	onceMap      = sync.Map{}
)

func exitWait() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	exitRun()
	os.Exit(0)
}

func exitRun() {
	for _, exit := range exitList {
		exit()
	}
}

func onceRun(hooks []*func()) {
	for _, ptr := range hooks {
		_, loaded := onceMap.LoadOrStore(ptr, struct{}{})
		if !loaded {
			(*ptr)()
		}
	}
}
