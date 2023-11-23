package runner_test

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/paulshpilsher/concurrent-go/concurrency/chan/runner"
)

func smallDelay() {
	time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
}

func TestRunnerNewDefault(t *testing.T) {
	r := runner.New(0)
	if r == nil {
		t.Error("Failed to construct new runner with defaults")
	}
	r.WaitAndClose()
}
func TestRunnerNew(t *testing.T) {
	r := runner.New(5)
	if r == nil {
		t.Error("Failed to construct new runner")
	}
	r.WaitAndClose()
}

func TestRunningBeforeClose(t *testing.T) {
	r := runner.New(5)

	err := r.Run(func() {})
	if err != nil {
		t.Error("Execution failed")
	}

	r.WaitAndClose()
}

func TestNoRunningAfterClose(t *testing.T) {
	r := runner.New(5)
	r.WaitAndClose()

	err := r.Run(func() {})
	if err == nil {
		t.Error("Executing after close should not be allowed")
	}
}

func TestRunningNilTask(t *testing.T) {
	r := runner.New(5)
	err := r.Run(nil)
	if err == nil {
		t.Error("Failed to error on nil argument")
	}
	r.WaitAndClose()
}

func TestExectuteHappyPath(t *testing.T) {
	r := runner.New(25)

	const numTasks = 1000
	cnt := int32(0)

	for i := 0; i < numTasks; i++ {
		r.Run(func() {
			atomic.AddInt32(&cnt, 1)
			smallDelay()
		})
	}

	r.WaitAndClose()

	if cnt != numTasks {
		t.Error("Not all tasks had executed")
	}
}

func TestConcurrencyLimit(t *testing.T) {
	const quota = 13
	const numTasks = 10000

	r := runner.New(quota)

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

	r.WaitAndClose()

	if failed {
		t.Error("Concurrency limit broken.")
	}
}

func TestQuota(t *testing.T) {
	r := runner.New(10)
	actual := r.GetQuota()
	if actual != 10 {
		t.Errorf("Exoected quota 10, actual %d", actual)
	}
}
func BenchmarkChannelRunner(b *testing.B) {
	r := runner.New(2)
	for i := 0; i < b.N; i++ {
		r.Run(func() {})
	}
	r.WaitAndClose()
}
