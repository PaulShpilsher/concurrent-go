package runner

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const DefaultConcurrencyLimit = 50

type ConcurrentRunner struct {
	concurrencyLimit int
	availableSlots   chan struct{}
	runnungTasks     int32
}

// Creates new ConcurrentRunner with DefaultConcurrencyLimit concurrency limit
func DefaultConcurrentRunner() *ConcurrentRunner {
	return NewConcurrentRunner(DefaultConcurrencyLimit)
}

// Creates new ConcurrentRunner with specified concurrency limit
func NewConcurrentRunner(concurrencyLimit int) *ConcurrentRunner {
	if concurrencyLimit <= 0 {
		panic("concurrency limit must be > 0")
	}

	slots := make(chan struct{}, concurrencyLimit)
	for i := 0; i < concurrencyLimit; i++ {
		slots <- struct{}{}
	}

	return &ConcurrentRunner{
		concurrencyLimit: concurrencyLimit,
		availableSlots:   slots,
		runnungTasks:     0,
	}
}

// Closes ConcurrencyLimiter
// But first it waits for all pending tasks to complete
func (t *ConcurrentRunner) Close() {
	for i := 0; i < t.concurrencyLimit; i++ {
		<-t.availableSlots
	}
	close(t.availableSlots)
}

// Gets number of currently executing tasks
func (t *ConcurrentRunner) GetNumberOfRunningTasks() int {
	return int(atomic.LoadInt32(&t.runnungTasks))
}

// Executes a task concurrently
// if there are no available slots it blocks until one becomes available
func (t *ConcurrentRunner) Run(task func()) error {
	if task == nil {
		return errors.New("nil task argument")
	}

	if _, ok := <-t.availableSlots; !ok {
		return errors.New("channel closed")
	}

	atomic.AddInt32(&t.runnungTasks, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered panic in goroutine", r)
			}
			atomic.AddInt32(&t.runnungTasks, -1)
			t.availableSlots <- struct{}{}
		}()

		task()
	}()

	return nil
}
