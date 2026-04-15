package runner

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/Compogo/compogo/closer"
	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/types/mapper"
)

// Runner defines the interface for managing concurrent tasks.
// It provides methods to start, stop, and manage the lifecycle of multiple tasks.
type Runner interface {
	// Closer stops all running tasks and waits for their completion.
	// It implements io.Closer for integration with Compogo's lifecycle.
	io.Closer

	// RunTask starts a single task. Returns an error if the task is already running.
	RunTask(task *Task) error

	// RunTasks starts multiple tasks sequentially. If any task fails to start,
	// the error is returned immediately and subsequent tasks are not started.
	RunTasks(tasks ...*Task) error

	// StopTask stops a specific task by canceling its context.
	// The task is removed from the runner after it completes.
	StopTask(task *Task) error

	// StopTaskByName stops a task identified by its name.
	// Returns TaskUndefinedError if no task with the given name exists.
	StopTaskByName(name string) error

	// HasTaskByName checks if a task with the given name is currently running.
	// Returns true if the task exists, false otherwise.
	HasTaskByName(name string) bool

	// HasTask checks if the specific task instance is currently running.
	// This is useful for checking against a known task pointer.
	HasTask(task *Task) bool

	// Use registers one or more middlewares that will wrap all tasks
	// executed by this runner. Middlewares are applied in the order they are added,
	// with the last added middleware being the outermost wrapper.
	Use(middlewares ...Middleware)
}

// runner is the internal implementation of the Runner interface.
type runner struct {
	wg sync.WaitGroup

	// tasks stores running tasks indexed by both name and value.
	// Using Mapper allows efficient lookup by name and deduplication.
	tasks *mapper.Mapper[*Task]

	closer closer.Closer
	logger logger.Logger

	// middleware holds the chain of middlewares that wrap all tasks.
	// They are applied in the order they were added via Use().
	middleware []Middleware

	rwMutex sync.RWMutex
}

// NewRunner creates a new Runner instance.
// The provided closer is used to derive contexts for tasks, enabling
// graceful shutdown when the application stops.
// The logger is used for task lifecycle events and error reporting
func NewRunner(closer closer.Closer, logger logger.Logger) Runner {
	return &runner{
		tasks:  mapper.NewMapper[*Task](),
		closer: closer,
		logger: logger.GetLogger("runner"),
	}
}

// Use registers one or more middlewares that will wrap all tasks
// executed by this runner. Middlewares are applied in the order they are added,
// meaning the first middleware added will be the innermost wrapper,
// and the last middleware added will be the outermost wrapper.
//
// Example:
//
//	r.Use(loggerMiddleware, metricsMiddleware)
//	// Result: metricsMiddleware(loggerMiddleware(task))
func (r *runner) Use(middlewares ...Middleware) {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	r.middleware = append(r.middleware, middlewares...)
}

// Close stops all running tasks and waits for their completion.
// It first cancels all task contexts, then waits for the WaitGroup.
// Returns the first error encountered during task shutdown.
func (r *runner) Close() (err error) {
	r.rwMutex.RLock()
	tasks := r.tasks
	r.rwMutex.RUnlock()

	for _, task := range tasks.All() {
		if err = r.StopTask(task); err != nil {
			return err
		}
	}

	r.wg.Wait()

	return nil
}

// RunTasks starts multiple tasks sequentially.
// If any task fails to start (e.g., due to duplicate name), the error is
// returned immediately and no further tasks are started.
func (r *runner) RunTasks(tasks ...*Task) (err error) {
	for _, task := range tasks {
		if err = r.RunTask(task); err != nil {
			return err
		}
	}

	return nil
}

// RunTask starts a single task.
// The task is registered in the runner, its context is derived from the
// closer's context, and it's launched in a goroutine.
// Returns TaskAlreadyExistsError if a task with the same pointer is already running.
func (r *runner) RunTask(task *Task) error {
	if r.HasTask(task) {
		return fmt.Errorf("[runner] task '%s': %w", task.name, TaskAlreadyExistsError)
	}

	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	r.tasks.Add(task)

	task.ctx, task.cancelFunc = context.WithCancel(r.closer.GetContext())

	process := task.process
	for _, middleware := range r.middleware {
		process = middleware.Middleware(task, process)
	}

	r.wg.Add(1)
	go func(task *Task, process Process) {
		defer r.wg.Done()
		defer func() {
			if err := r.StopTaskByName(task.Name()); err != nil {
				r.logger.Errorf("task '%s' remove failed: %w", task.Name(), err)
			}
		}()

		_ = process.Process(task.ctx)
	}(task, process)

	return nil
}

// StopTaskByName stops a task identified by its name.
// It looks up the task in the mapper and delegates to StopTask.
// Returns TaskUndefinedError if no task with the given name exists.
func (r *runner) StopTaskByName(name string) error {
	if !r.HasTaskByName(name) {
		return fmt.Errorf("[runner] task '%s': %w", name, TaskUndefinedError)
	}

	r.rwMutex.RLock()
	task, _ := r.tasks.Get(name)
	r.rwMutex.RUnlock()

	return r.StopTask(task)
}

// StopTask stops a specific task by canceling its context.
// The task is removed from the runner immediately (the goroutine will
// clean up after the Process function returns).
func (r *runner) StopTask(task *Task) error {
	task.cancelFunc()

	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	r.tasks.RemoveByValue(task)

	return nil
}

// HasTaskByName checks if a task with the given name is currently running.
// It acquires a read lock to safely check the internal task map.
// Returns true if a task with that name exists, otherwise false.
func (r *runner) HasTaskByName(name string) bool {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	return r.tasks.HasByKey(name)
}

// HasTask checks if the specific task instance is currently running.
// This is useful for verifying the state of a task you have a reference to,
// especially after operations like RunTask or StopTask.
// Returns true if the exact task is found, otherwise false.
func (r *runner) HasTask(task *Task) bool {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	return r.tasks.HasByValue(task)
}
