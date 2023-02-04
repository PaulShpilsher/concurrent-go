package runner

const DefaultMaxTasks = 50

type TaskRunner struct {
	concurrent   chan struct{}
	runnungTasks uint
}

func NewTaskRunner(limit uint) *TaskRunner {
	if limit == 0 {
		limit = DefaultMaxTasks
	}
	return &TaskRunner{
		concurrent:   make(chan struct{}, limit),
		runnungTasks: 0,
	}
}
