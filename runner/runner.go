package runner

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const DefaultConcurrencyLimit = 50

type TaskRunner struct {
	availableSlots chan struct{}
	runnungTasks   int32
}

func NewTaskRunner(concurrencyLimit int) *TaskRunner {
	if concurrencyLimit <= 0 {
		concurrencyLimit = DefaultConcurrencyLimit
	}

	slots := make(chan struct{}, concurrencyLimit)
	for i := 0; i < concurrencyLimit; i++ {
		slots <- struct{}{}
	}

	return &TaskRunner{
		availableSlots: slots,
		runnungTasks:   0,
	}
}

func (t *TaskRunner) ExecuteTask(task func()) error {
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
				fmt.Println("Recovered panic in task", r)
			}
			atomic.AddInt32(&t.runnungTasks, -1)
			t.availableSlots <- struct{}{}
		}()

		task()
	}()

	return nil
}
