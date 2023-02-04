# Simple Golang concurrent task runner

Task runner limiting how many concurrent goroutines can be executing at the same time.

It uses channels to enforce maximum number of goroutines execiting simultaneously.  When the concurency quota is reached scheduling new execution will be blocked until at least one of the executing tasks finishes.

Panic inside a goroutine is handled and is not going to kill the program.  The panic message however is sent to the console.

## Example

```go
 runner := concurrent.NewConcurrenRunner(12)
 
 for i := 0; i < 1000; i++ {
  runner.Run(func() {
   // task's code to exectute
  })
 }

 runner.Close()
```
