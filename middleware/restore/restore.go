package restore

import (
	"context"
	"time"

	"github.com/Compogo/runner"
)

// Restore — middleware для автоматического восстановления процессов.
// При ошибке выполнения процесса перезапускает его с задержкой.
// Работает до тех пор, пока не будет отменён контекст.
type Restore struct {
	config *Config

	runner runner.Runner
}

func NewRestore(config *Config, runner runner.Runner) *Restore {
	return &Restore{
		config: config,
		runner: runner,
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
