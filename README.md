# Golang concurrent task runner with quota

Goroutine runner limiting how many goroutines can be executing at the same time.

It uses channels to enforce maximum number of goroutines execiting simultaneously.  When the concurency quota is reached scheduling new executions will be blocked until other tasks finish.

Panic inside a goroutine is handled and is not going to kill the program.  The panic message however is sent to the console.

## Example

```go
import "github.com/PaulShpilsher/concurrent-go/runner"
```

```go
 runner := runner.NewConcurrentRunner(12)
 
 for i := 0; i < 1000; i++ {
  runner.Run(func() {
   // task's code to exectute
  })
 }

 runner.Close()
```
