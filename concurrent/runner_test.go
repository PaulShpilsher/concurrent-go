package concurrent_test

import (
	"psp/concurrent-worker-pool-go/concurrent"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func smallDelay() {
	time.Sleep(time.Millisecond)
}

func TestRunnerDefault(t *testing.T) {
	concurrentLimiter := concurrent.DefaultConcurrentRunner()
	concurrentLimiter.Close()
}

func TestRunningBeforeClose(t *testing.T) {
	runner := concurrent.DefaultConcurrentRunner()

	err := runner.Run(func() {})
	if err != nil {
		t.Error("Execution failed")
	}

	runner.Close()
}

func TestNoRunningAfterClose(t *testing.T) {
	runner := concurrent.DefaultConcurrentRunner()
	runner.Close()

	err := runner.Run(func() {})
	if err == nil {
		t.Error("Executing after close should not be allowed")
	}
}

func TestExectuteHappyPath(t *testing.T) {
	runner := concurrent.NewConcurrenRunner(25)

	const numTasks = 1000
	cnt := int32(0)

	for i := 0; i < numTasks; i++ {
		runner.Run(func() {
			atomic.AddInt32(&cnt, 1)
			smallDelay()
		})
	}

	runner.Close()

	if cnt != numTasks {
		t.Error("Not all tasks had executed")
	}
}

func TestConcurrencyLimit(t *testing.T) {
	const limit = 13
	const numTasks = 10000

	runner := concurrent.NewConcurrenRunner(limit)

	mutex := &sync.Mutex{}
	failed := false

	for i := 0; i < numTasks; i++ {
		runner.Run(func() {
			cnt := runner.GetNumberOfRunningTasks()
			if cnt > limit {
				mutex.Lock()
				failed = true
				mutex.Unlock()
			}
		})
	}

	runner.Close()

	if failed {
		t.Error("Concurrency limit broken.")
	}
}
