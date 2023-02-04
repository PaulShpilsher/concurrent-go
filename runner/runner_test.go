package runner_test

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/PaulShpilsher/concurrent-go/runner"
)

func smallDelay() {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
}

func TestRunnerDefault(t *testing.T) {
	r, err := runner.NewConcurrencyRunner(runner.DefaultQuota)
	if err != nil {
		t.Error("NewConcurrencyRunner failed", err)
	}
	r.Close()

}

func TestRunningBeforeClose(t *testing.T) {
	r, _ := runner.NewConcurrencyRunner(runner.DefaultQuota)

	err := r.Run(func() {})
	if err != nil {
		t.Error("Execution failed")
	}

	r.Close()
}

func TestNoRunningAfterClose(t *testing.T) {
	r, _ := runner.NewConcurrencyRunner(runner.DefaultQuota)
	r.Close()

	err := r.Run(func() {})
	if err == nil {
		t.Error("Executing after close should not be allowed")
	}
}

func TestExectuteHappyPath(t *testing.T) {
	r, _ := runner.NewConcurrencyRunner(25)

	const numTasks = 1000
	cnt := int32(0)

	for i := 0; i < numTasks; i++ {
		r.Run(func() {
			atomic.AddInt32(&cnt, 1)
			smallDelay()
		})
	}

	r.Close()

	if cnt != numTasks {
		t.Error("Not all tasks had executed")
	}
}

func TestConcurrencyLimit(t *testing.T) {
	const quota = 13
	const numTasks = 10000

	r, _ := runner.NewConcurrencyRunner(quota)

	mutex := &sync.Mutex{}
	failed := false

	for i := 0; i < numTasks; i++ {
		r.Run(func() {
			cnt := r.GetNumberOfRunningTasks()
			if cnt > quota {
				mutex.Lock()
				failed = true
				mutex.Unlock()
			}
		})
	}

	r.Close()

	if failed {
		t.Error("Concurrency limit broken.")
	}
}
