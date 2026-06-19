package recover

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
)

// Recover — middleware для восстановления после паник.
// Перехватывает panic в процессе, логирует стек-трейс и преобразует panic в ошибку.
type Recover struct {
	logger compogo.Logger
}

func NewRecover(logger compogo.Logger) *Recover {
	return &Recover{logger: logger.GetLogger("runner").GetLogger("middleware").GetLogger("recover")}
}

func (m *Recover) Middleware(process runner.Process, next runner.ProcessFunc) runner.ProcessFunc {
	return func(ctx context.Context) (err error) {
		defer func() {
			if rerr := recover(); rerr != nil {
				m.logger.Errorf("process '%s' panic: %s\n%s", process.Name(), rerr, debug.Stack())

				var panicErr error
				if e, ok := rerr.(error); ok {
					panicErr = e
				} else {
					panicErr = fmt.Errorf("%v", rerr)
				}

				if err == nil {
					err = fmt.Errorf("panic: %w", panicErr)
				} else {
					err = fmt.Errorf("%w; panic: %w", err, panicErr)
				}
			}
		}()

		err = next(ctx)
		return err
	}
}
