package runner

// Middleware — интерфейс для цепочки обработчиков процессов.
// Позволяет оборачивать ProcessFunc для добавления логирования,
// метрик, восстановления после паник и т.д.
//
// Пример:
//
//	type LoggerMiddleware struct {
//	    logger Logger
//	}
//
//	func (m *LoggerMiddleware) Middleware(p Process, next ProcessFunc) ProcessFunc {
//	    return func(ctx context.Context) error {
//	        m.logger.Info("starting", "process", p.Name())
//	        err := next(ctx)
//	        m.logger.Info("finished", "process", p.Name())
//	        return err
//	    }
//	}
type Middleware interface {
	// Middleware оборачивает ProcessFunc.
	// Получает текущий Process и следующую функцию в цепочке.
	// Возвращает новую функцию, которая будет вызвана вместо следующей.
	Middleware(Process, ProcessFunc) ProcessFunc
}

// MiddlewareFunc — функциональный адаптер для Middleware.
// Позволяет использовать функции как middleware.
//
// Пример:
//
//	middleware := MiddlewareFunc(func(p Process, next ProcessFunc) ProcessFunc {
//	    return func(ctx context.Context) error {
//	        // before
//	        err := next(ctx)
//	        // after
//	        return err
//	    }
//	})
type MiddlewareFunc func(Process, ProcessFunc) ProcessFunc

// Middleware реализует интерфейс Middleware для MiddlewareFunc.
func (m MiddlewareFunc) Middleware(p Process, next ProcessFunc) ProcessFunc {
	return m(p, next)
}
