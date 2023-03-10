package runner

import (
	"errors"
	"log"
	"math"
	"runtime"
	"sync/atomic"

	"github.com/PaulShpilsher/concurrent-go/concurrency"
)

type channelRunner struct {
	quota          int
	freeSlots      chan struct{}
	executingCount int32
	closed         bool
}

// Creates new concurrency Runner with specified quota on
// the maximum number of concurrently running functions.
// Returns concurrency.Runner interface
// The "quota" argument is in range [1..math.MaxInt32]
// if specified value is outside of this range the quota falls back
// to the number of CPU cores, as runtime.GOMAXPROCS(0)
func New(quota int) concurrency.Runner {
	if quota <= 0 || quota > math.MaxInt32 {
		quota = runtime.GOMAXPROCS(0)
		log.Printf("Using default runner quota %d\n", quota)
	}

	slots := make(chan struct{}, quota)
	for i := 0; i < quota; i++ {
		slots <- struct{}{}
	}

	return &channelRunner{
		quota:          quota,
		freeSlots:      slots,
		executingCount: 0,
		closed:         false,
	}
}

// Waits for all running functions to complete,
// Then releases internally used resources.
// No more calls to Run() are possible.
func (r *channelRunner) WaitAndClose() {
	if !r.closed {
		r.closed = true
		for i := 0; i < r.quota; i++ {
			<-r.freeSlots
		}
		close(r.freeSlots)
	}
}

// Gets number of currently executing routines
func (r *channelRunner) GetNumberOfRunningTasks() int {
	return int(atomic.LoadInt32(&r.executingCount))
}

// Returns the quota limit
func (r *channelRunner) GetQuota() int {
	return r.quota
}

// Concurrently executes a function wrapped in a goroutine.
// If the quota of currently running functions is reached
// a call to this function will block until another running function finishes.
func (r *channelRunner) Run(task func()) error {
	if task == nil {
		return errors.New("nil  argument")
	}

	if r.closed {
		return errors.New("runned closed")
	}

	if _, ok := <-r.freeSlots; !ok {
		return errors.New("channel closed")
	}

	atomic.AddInt32(&r.executingCount, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered panic in goroutine", r)
			}
			atomic.AddInt32(&r.executingCount, -1)
			r.freeSlots <- struct{}{}
		}()

		task()
	}()

	return nil
}
