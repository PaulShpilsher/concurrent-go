package concurrency

// Defines concurrency runner functionality.
type Runner interface {
	// Concurrently executes a function wrapped in a goroutine.
	Run(task func()) error

	// Waits for all running functions to complete and frees resources.
	WaitAndClose() error

	// Returns the number of currently executing functions.
	GetNumberOfRunningTasks() int

	// Returns the quota limit
	GetQuota() int
}
