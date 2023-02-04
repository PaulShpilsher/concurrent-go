# Golang concurrent function runner with quota

A concurrent function runner with quota on how many functions can be executing at the same time.

## Features

It uses channels to enforce maximum number of goroutines executing simultaneously. It also maintains an atomic counter of how many functions are executing at any point of time.

When the concurrency quota is reached scheduling new function executions is blocked until some of the running functions finish.

A panic inside a goroutine is handled by logging it to console, and will not stop program execution.

## Usage

```bash
go get github.com/PaulShpilsher/concurrent-go
```

```go
import "github.com/PaulShpilsher/concurrent-go/runner"
```

```go
 concurrentRunn, err := runner.NewConcurrentRunner(12)
 if err != nil {
    panic(err.Error())
 }
 
 for i := 0; i < 1000; i++ {
  concurrentRunn.Run(func() {
   // task's code to exectute
  })
 }

 concurrentRunn.Close()
```
