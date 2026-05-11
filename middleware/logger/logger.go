package logger

import (
	"context"
	"runtime/debug"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/runner"
)

type Logger struct {
	logger logger.Logger
}

func NewLogger(logger logger.Logger) *Logger {
	return &Logger{logger: logger.GetLogger("runner").GetLogger("middleware").GetLogger("logger")}
}

func (m *Logger) Middleware(process runner.Process, next runner.ProcessFunc) runner.ProcessFunc {
	return func(ctx context.Context) (err error) {
		m.logger.Infof("task '%s' running", process.Name())

		if err = next(ctx); err != nil {
			m.logger.Errorf("task '%s' error: %s\n%s", process.Name(), err, debug.Stack())
		}

		m.logger.Infof("task '%s' shutdown", process.Name())

		return err
	}
}
