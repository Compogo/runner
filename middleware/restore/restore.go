package restore

import (
	"context"
	"time"

	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
)

// Restore — middleware для автоматического восстановления процессов.
// При ошибке выполнения процесса перезапускает его с задержкой.
// Работает до тех пор, пока не будет отменён контекст.
type Restore struct {
	config *Config

	runner runner.Runner
	logger compogo.Logger
}

func NewRestore(config *Config, runner runner.Runner, logger compogo.Logger) *Restore {
	return &Restore{
		config: config,
		runner: runner,
		logger: logger.GetLogger("runner").GetLogger("middleware").GetLogger("restore"),
	}
}

func (restore *Restore) Middleware(_ runner.Process, next runner.ProcessFunc) runner.ProcessFunc {
	return func(ctx context.Context) (err error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for {
			err = next(ctx)
			if err == nil {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(restore.config.Delay):
				continue
			}
		}
	}
}
