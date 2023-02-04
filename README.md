# Golang concurrent task runner

A runner with limit of how many concurrent goroutines can be execiting at the same time.

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
