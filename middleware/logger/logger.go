package logger

import (
	"context"
	"runtime/debug"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/runner"
)

// Logger is a middleware that logs task lifecycle events:
// - When a task starts running
// - When a task finishes (with error details if any)
type Logger struct {
	logger logger.Logger
}

// NewLogger creates a new Logger middleware instance.
// The provided logger is used to output task lifecycle messages.
func NewLogger(logger logger.Logger) *Logger {
	return &Logger{logger: logger}
}

// Middleware wraps a task's process function with lifecycle logging.
// It logs "running" before the task starts, and "shutdown" (with error)
// after the task completes. If the task returns an error, it is logged
// with a stack trace for debugging.
func (m *Logger) Middleware(task *runner.Task, next runner.Process) runner.Process {
	return runner.ProcessFunc(func(ctx context.Context) error {
		m.logger.Infof("task '%s' running", task.Name())

		err := next.Process(ctx)
		if err != nil {
			m.logger.Errorf("task '%s' error: %s\n%s", task.Name(), err, debug.Stack())
		}

		m.logger.Infof("task '%s' shutdown", task.Name())

		return err
	})
}
