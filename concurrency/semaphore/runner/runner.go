package runner

import (
	"context"
	"log"
	"math"
	"runtime"
	"sync/atomic"

	"github.com/PaulShpilsher/concurrent-go/concurrency"
	"golang.org/x/sync/semaphore"
)

type semaphoreRunner struct {
	quota          int
	sem            *semaphore.Weighted
	ctx            context.Context
	executingCount int32
}

func NewRunner(quota int) concurrency.Runner {
	if quota <= 0 || quota > math.MaxInt32 {
		quota = runtime.GOMAXPROCS(0)
		log.Printf("Using default runner quota %d", quota)
	}

	return &semaphoreRunner{
		quota:          quota,
		sem:            semaphore.NewWeighted(int64(quota)),
		ctx:            context.TODO(),
		executingCount: 0,
	}
}

func (r *semaphoreRunner) Run(task func()) error {
	if err := r.sem.Acquire(r.ctx, 1); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
		return err
	}

	atomic.AddInt32(&r.executingCount, 1)
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("Recovered panic in goroutine %v", recovered)
			}
			atomic.AddInt32(&r.executingCount, -1)
			r.sem.Release(1)
		}()

		task()
	}()

	return nil
}

func (r *semaphoreRunner) WaitAndClose() {
	if err := r.sem.Acquire(r.ctx, int64(r.quota)); err != nil {
		log.Printf("Failed to acquire semaphore: %v", err)
	}
}
