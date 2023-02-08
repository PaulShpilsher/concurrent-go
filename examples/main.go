package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/PaulShpilsher/concurrent-go/concurrency/sync/runner"
)

func worker(id int) {
	fmt.Println("worker start", id)
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	fmt.Println("worker end  ", id)
}

func main() {
	var quota int
	var numbeOfTasks int

	flag.IntVar(&quota, "limit", 8, "maxiumum number of goroutines allowed to run concurrently")
	flag.IntVar(&numbeOfTasks, "tasks", 64, "the number of tasks we need to execute")
	flag.Parse()

	fmt.Printf("Executing %d tasks with concurrency limit %d\n", numbeOfTasks, quota)

	//r, _ := runner.NewConcurrencyRunner(quota)
	r := runner.New(quota)
	for i := 0; i < numbeOfTasks; i++ {
		id := i
		r.Run(func() { worker(id) })
	}

	r.WaitAndClose() // waits for all pending tasks to complete and closes runner
}
