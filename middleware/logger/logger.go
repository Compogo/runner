package logger

import (
	"context"

	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
)

// Logger — middleware для логирования выполнения процессов.
// Логирует:
//   - Начало выполнения процесса
//   - Успешное завершение
//   - Ошибки
type Logger struct {
	logger compogo.Logger
}

func NewLogger(logger compogo.Logger) *Logger {
	return &Logger{logger: logger.GetLogger("runner").GetLogger("middleware").GetLogger("logger")}
}

func (m *Logger) Middleware(process runner.Process, next runner.ProcessFunc) runner.ProcessFunc {
	return func(ctx context.Context) (err error) {
		m.logger.Infof("process '%s' running", process.Name())

		if err = next(ctx); err != nil {
			m.logger.Errorf("process '%s' error: %s", process.Name(), err)
		}

		m.logger.Infof("process '%s' shutdown", process.Name())

		return err
	}
}
