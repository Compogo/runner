package recover

import (
	"context"
	"runtime/debug"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/runner"
)

type Recover struct {
	logger logger.Logger
}

func NewRecover(logger logger.Logger) *Recover {
	return &Recover{logger: logger.GetLogger("runner").GetLogger("middleware").GetLogger("recover")}
}

func (m *Recover) Middleware(process runner.Process, next runner.ProcessFunc) runner.ProcessFunc {
	return func(ctx context.Context) error {
		defer func() {
			if err := recover(); err != nil {
				m.logger.Errorf("process '%s' panic: %s\n%s", process.Name(), err, debug.Stack())
			}
		}()

		return next(ctx)
	}
}
