package runner

import "context"

type Task struct {
	name        string
	processFunc ProcessFunc

	ctx        context.Context
	cancelFunc context.CancelFunc

	middlewares []Middleware
}

func NewTask(name string, processFunc ProcessFunc, middlewares ...Middleware) *Task {
	return &Task{
		name:        name,
		processFunc: processFunc,
		middlewares: middlewares,
	}
}

func (task *Task) Process(ctx context.Context) error {
	task.ctx, task.cancelFunc = context.WithCancel(ctx)
	defer task.cancelFunc()

	processFunc := task.processFunc
	for i := len(task.middlewares) - 1; i >= 0; i-- {
		processFunc = task.middlewares[i].Middleware(task, processFunc)
	}

	return processFunc(task.ctx)
}

func (task *Task) Name() string {
	return task.name
}

func (task *Task) Close() error {
	if task.cancelFunc != nil {
		task.cancelFunc()
	}

	return nil
}
