# Simple Golang concurrent task runner

Task runner limiting how many concurrent goroutines can be executing at the same time.

It uses channels to enforce maximum number of goroutines execiting simultaneously.

## Example
```go
    runner := concurrent.NewConcurrenRunner(limit)
	for i := 0; i < tasks; i++ {
		runner.Run(func() {
            // task's code to exectute
         })
	}

	runner.Close()
```
