# Golang concurrent function runner with quota
[![Go Report Card](https://goreportcard.com/badge/github.com/PaulShpilsher/concurrent-go)](https://goreportcard.com/report/github.com/PaulShpilsher/concurrent-go)

A concurrent function runner with quota on how many functions can be executing at the same time.

## Features

A concurrent runner purpose is to enforce maximum number of goroutines that can execute simultaneously (quota).  When the quota is reached scheduling new function executions is blocked until some of the running functions are finished.

It also maintains an atomic counter of how many functions are executing at any point of time.

## Implementation

There are two flavors of concurrent runners are implemented.  One that uses semaphore synchronization primitive and the other uses channels.

Both have common functionality described by the interface:

```go
type Runner interface {
   // Concurrently executes a function wrapped in a goroutine.
   Run(task func()) error

   // Waits for all running functions to complete and frees resources.
   WaitAndClose()

   // Returns the number of currently executing functions.
   GetNumberOfRunningTasks() int

   // Returns the quota limit
   GetQuota() int
}
```

## Usage

Get the package

```bash
go get github.com/PaulShpilsher/concurrent-go
```

In code import runner.

To use channel-based runner:

```go
import "github.com/PaulShpilsher/concurrent-go/concurrency/chan/runner"
```

To use sync-based runner:

```go
import "github.com/PaulShpilsher/concurrent-go/concurrency/sync/runner"
```

Use the runner

```go
 theRunner := runner.New(quota)
 if err != nil {
    panic(err.Error())
 }
 
 for i := 0; i < 1000; i++ {
  theRunner.Run(func() {
   // put some code to be exectuted
  })
 }

 theRunner.WaitAndClose()
```

## Examples

The exapmples are in the ./examples/ directory.

Running examples using make utility:

```shell
make run-example-channel
```

or

```shell
make run-example-sync
```

## Testing

Running unit tests using make utility:

```shell
make test-channel-runner
```

or

```shell
make test-sync-runner
```

## Acknowledgements

- Inspiration ["Simple Made Easy" - Rich Hickey (2011)](https://www.youtube.com/watch?v=SxdOUGdseq4)
- References
