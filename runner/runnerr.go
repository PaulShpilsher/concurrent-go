package runner

import (
	"errors"
	"fmt"
	"math"
	"sync/atomic"
)

const DefaultQuota = 16

type ConcurrencyRunner struct {
	quota          int
	freeSlots      chan struct{}
	executingCount int32
}

// Creates new ConcurrencyRunner with specified quota of maximum concurrently running functions
// Returnss new;y constructed &ConcurrencyRunner struct
// Error if "quota" argument is out of range [1..math.MaxInt32]
func NewConcurrencyRunner(quota int) (*ConcurrencyRunner, error) {
	if quota <= 0 || quota > math.MaxInt32 {
		return &ConcurrencyRunner{}, fmt.Errorf("quota must be in range [1..%d]", math.MaxInt32)
	}

	slots := make(chan struct{}, quota)
	for i := 0; i < quota; i++ {
		slots <- struct{}{}
	}

	return &ConcurrencyRunner{
		quota:          quota,
		freeSlots:      slots,
		executingCount: 0,
	}, nil
}

// Waits for all still running functions to complete,
// Then releases internally used resources.
// No more calls to Run() are possible.
func (t *ConcurrencyRunner) Close() {
	for i := 0; i < t.quota; i++ {
		<-t.freeSlots
	}
	close(t.freeSlots)
}

// Gets number of currently executing routines
func (t *ConcurrencyRunner) GetNumberOfRunningTasks() int {
	return int(atomic.LoadInt32(&t.executingCount))
}

// Concurrently executes a function wrapped in a goroutine.
// Hoever if the quota of currently running functions are reached
// this blocks until one of the rinning function completes and "makes room"
func (t *ConcurrencyRunner) Run(task func()) error {
	if task == nil {
		return errors.New("nil  argument")
	}

	if _, ok := <-t.freeSlots; !ok {
		return errors.New("channel closed")
	}

	atomic.AddInt32(&t.executingCount, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered panic in goroutine", r)
			}
			atomic.AddInt32(&t.executingCount, -1)
			t.freeSlots <- struct{}{}
		}()

		task()
	}()

	return nil
}
