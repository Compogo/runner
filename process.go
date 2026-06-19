package runner

import (
	"context"
	"io"
)

// Process — интерфейс выполняемого процесса.
// Сочетает в себе:
//   - io.Closer — для остановки процесса
//   - Process(ctx) — для выполнения работы
//   - Name() — для идентификации
//
// Процессы запускаются через Runner и выполняются в отдельных горутинах.
type Process interface {
	io.Closer
	Process(ctx context.Context) error
	Name() string
}

// ProcessFunc — функция, выполняемая процессом.
// Принимает контекст и возвращает ошибку.
type ProcessFunc func(ctx context.Context) error
