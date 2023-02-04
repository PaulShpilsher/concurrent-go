package main

import (
	"flag"
	"fmt"
	"psp/concurrent-worker-pool-go/concurrent"
	"time"
)

func worker(id int) {
	fmt.Println("worker start", id)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("worker end  ", id)
}

func main() {
	var limit int
	var tasks int

	flag.IntVar(&limit, "limit", 2, "maxiumum number of goroutines allowed to run concurrently")
	flag.IntVar(&tasks, "tasks", 5, "the number of tasks we need to execute")
	flag.Parse()

	fmt.Printf("Executing %d tasks with concurrency limit %d\n", tasks, limit)

	runner := concurrent.NewConcurrentLimiter(limit)
	for i := 0; i < tasks; i++ {
		id := i
		runner.Run(func() { worker(id) })
	}

	runner.Close()

}
