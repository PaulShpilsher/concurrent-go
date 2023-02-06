package concurrency

type Runner interface {
	Run(task func()) error
	WaitAndClose()
	GetNumberOfRunningTasks() int
}
