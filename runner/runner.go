package runner

import (
	"errors"
	"fmt"
	"sync/atomic"
)

const DefaultMaxTasks = 50

type TaskRunner struct {
	availableSlots chan struct{}
	runnungTasks   int32
}

func NewTaskRunner(limit int) *TaskRunner {
	if limit <= 0 {
		limit = DefaultMaxTasks
	}

	slots := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
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
}
