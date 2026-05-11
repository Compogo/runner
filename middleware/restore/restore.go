package restore

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/runner"
)

type Restore struct {
	config *Config

	runner runner.Runner
	logger logger.Logger
}

func NewRestore(config *Config, runner runner.Runner, logger logger.Logger) *Restore {
	return &Restore{
		config: config,
		runner: runner,
		logger: logger.GetLogger("runner").GetLogger("middleware").GetLogger("restore"),
	}
}

func (restore *Restore) Middleware(process runner.Process, next runner.ProcessFunc) runner.ProcessFunc {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				if err := restore.process(ctx, process, next); err != nil {
					restore.logger.Errorf("process '%s' failed: %s", process.Name(), err)
					<-time.After(restore.config.Delay)
					continue
				}
			}

			return nil
		}
	}
}

func (restore *Restore) process(ctx context.Context, process runner.Process, next runner.ProcessFunc) (err error) {
	defer func() {
		if errr := recover(); errr != nil {
			err = fmt.Errorf("process '%s' panic: %s\n%s", process.Name(), errr, debug.Stack())
		}
	}()

	return next(ctx)
}
