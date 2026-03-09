package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime/debug"
	"sync"

	"github.com/Compogo/compogo/closer"
	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/compogo/types"
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
}

// runner is the internal implementation of the Runner interface.
type runner struct {
	wg sync.WaitGroup

	// tasks stores running tasks indexed by both name and value.
	// Using Mapper allows efficient lookup by name and deduplication.
	tasks   types.Mapper[*Task]
	rwMutex sync.RWMutex

	closer closer.Closer
	logger logger.Logger
}

// NewRunner creates a new Runner instance.
// The provided closer is used to derive contexts for tasks, enabling
// graceful shutdown when the application stops.
// The logger is used for task lifecycle events and error reporting
func NewRunner(closer closer.Closer, logger logger.Logger) Runner {
	return &runner{
		closer: closer,
		logger: logger,
	}
}

// Close stops all running tasks and waits for their completion.
// It first cancels all task contexts, then waits for the WaitGroup.
// Returns the first error encountered during task shutdown.
func (r *runner) Close() (err error) {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	for _, task := range r.tasks.All() {
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
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	if r.tasks.HasByValue(task) {
		return fmt.Errorf("[runner] task '%s': %w", task.name, TaskAlreadyExistsError)
	}

	r.tasks.Add(task)

	task.ctx, task.cancelFunc = context.WithCancel(r.closer.GetContext())

	r.wg.Add(1)
	go func(task *Task) {
		defer r.wg.Done()
		defer func() {
			if err := r.StopTaskByName(task.Name()); err != nil {
				r.logger.Errorf("task '%s' remove failed: %w", task.Name(), err)
			}
		}()
		defer func() {
			if err := recover(); err != nil {
				r.logger.Errorf("task '%s' panic: %s\n%s", task.Name(), err, debug.Stack())
			}
		}()

		r.logger.Infof("task '%s' running", task.Name())

		if err := task.process.Process(task.ctx); err != nil {
			r.logger.Errorf("task '%s' error: %s\n%s", task.Name(), err, debug.Stack())
		}

		r.logger.Infof("task '%s' shutdown", task.Name())
	}(task)

	return nil
}

// StopTaskByName stops a task identified by its name.
// It looks up the task in the mapper and delegates to StopTask.
// Returns TaskUndefinedError if no task with the given name exists.
func (r *runner) StopTaskByName(name string) error {
	r.rwMutex.RLock()

	task, err := r.tasks.Get(name)
	if errors.Is(err, types.DoesNotExistError) {
		r.rwMutex.RUnlock()
		return fmt.Errorf("[runner] task '%s': %w", name, TaskUndefinedError)
	}

	if err != nil {
		r.rwMutex.RUnlock()
		return fmt.Errorf("[runner] task '%s': %w", name, err)
	}

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
