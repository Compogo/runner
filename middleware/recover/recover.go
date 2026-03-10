package recover

import (
	"context"
	"runtime/debug"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/runner"
)

// Recover is a middleware that catches panics in task execution,
// logs them with stack traces, and prevents them from crashing the application.
type Recover struct {
	logger logger.Logger
}

// NewRecover creates a new Recover middleware instance.
// The provided logger is used to log panic details including stack traces.
func NewRecover(logger logger.Logger) *Recover {
	return &Recover{logger: logger}
}

// Middleware wraps a task's process function with panic recovery.
// If the wrapped process panics, it is recovered, logged, and the panic is
// not propagated further. The task will continue running (the panic only
// affects the current execution, not the task itself).
func (m *Recover) Middleware(task *runner.Task, next runner.Process) runner.Process {
	return runner.ProcessFunc(func(ctx context.Context) error {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Errorf("task '%s' panic: %s\n%s", task.Name(), err, debug.Stack())
			}
		}()

		return next.Process(ctx)
	})
}
