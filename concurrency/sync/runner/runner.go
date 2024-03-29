package runner

import (
	"context"
	"log"
	"math"
	"runtime"
	"sync/atomic"

	"github.com/paulshpilsher/concurrent-go/concurrency"
	"golang.org/x/sync/semaphore"
)

type semaphoreRunner struct {
	quota          int
	sem            *semaphore.Weighted
	ctx            context.Context
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
	}

	return &semaphoreRunner{
		quota:          quota,
		sem:            semaphore.NewWeighted(int64(quota)),
		ctx:            context.TODO(),
		executingCount: 0,
		closed:         false,
	}
}

// Concurrently executes a function wrapped in a goroutine.
// If the quota of currently running functions is reached
// a call to this function will block until another running
// function finishes.
func (r *semaphoreRunner) Run(task func()) error {
	if task == nil {
		return concurrency.ErrNilArgument
	}

	if r.closed {
		return concurrency.ErrRunnerClosed
	}

	if err := r.sem.Acquire(r.ctx, 1); err != nil {
		return err
	}

	atomic.AddInt32(&r.executingCount, 1)
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic recovered in goroutine %v", recovered)
			}
			atomic.AddInt32(&r.executingCount, -1)
			r.sem.Release(1)
		}()

		task()
	}()

	return nil
}

// Waits for all running functions to complete,
// Then releases internally used resources.
// No more calls to Run() are possible.
func (r *semaphoreRunner) WaitAndClose() error {
	if !r.closed {
		r.closed = true
		if err := r.sem.Acquire(r.ctx, int64(r.quota)); err != nil {
			return err
		}
		r.sem = nil
	}
	return nil
}

// Returns the number of currently executing functions
func (r *semaphoreRunner) GetNumberOfRunningTasks() int {
	return int(atomic.LoadInt32(&r.executingCount))
}

// Returns the quota limit
func (r *semaphoreRunner) GetQuota() int {
	return r.quota
}
