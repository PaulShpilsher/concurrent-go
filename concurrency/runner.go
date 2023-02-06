package concurrency

import (
	"log"
	"runtime"
)

type Runner interface {
	Run(task func()) error
	WaitAndClose()
}

func init() {
	maxWorkers := runtime.GOMAXPROCS(0)
	log.Println("runtime.GOMAXPROCS", maxWorkers)
}
