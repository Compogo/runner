package runner

import (
	"context"
	"sync"
)

// Task — базовая реализация интерфейса Process.
//
// Пример:
//
//	task := NewTask("worker", func(ctx context.Context) error {
//	    // работа
//	    return nil
//	})
//
//	runner.RunProcess(task)
type Task struct {
	rwm sync.RWMutex

	name        string
	processFunc ProcessFunc

	cancelFunc  context.CancelFunc
	middlewares []Middleware
}

// NewTask создаёт новый Task.
// Принимает имя, функцию выполнения и опциональные middleware.
func NewTask(name string, processFunc ProcessFunc, middlewares ...Middleware) *Task {
	return &Task{
		name:        name,
		processFunc: processFunc,
		middlewares: middlewares,
	}
}

// Process запускает выполнение задачи.
// Создаёт отменяемый контекст и применяет middleware.
// Реализует интерфейс Process.
func (task *Task) Process(ctx context.Context) error {
	task.rwm.Lock()
	ctx, task.cancelFunc = context.WithCancel(ctx)
	defer task.cancelFunc()
	task.rwm.Unlock()

	processFunc := task.processFunc
	for i := len(task.middlewares) - 1; i >= 0; i-- {
		processFunc = task.middlewares[i].Middleware(task, processFunc)
	}

	return processFunc(ctx)
}

// Name возвращает имя задачи.
// Реализует интерфейс Process.
func (task *Task) Name() string {
	return task.name
}

// Close останавливает задачу, отменяя её контекст.
// Реализует интерфейс io.Closer.
func (task *Task) Close() error {
	task.rwm.RLock()
	defer task.rwm.RUnlock()

	if task.cancelFunc != nil {
		task.cancelFunc()
	}

	return nil
}
