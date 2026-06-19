package runner

import (
	"github.com/Compogo/compogo"
)

// Component — компонент Runner для Compogo.
// Регистрирует Runner в DI-контейнере и останавливает его при завершении приложения.
var Component = compogo.Component{
	Name: "runner",
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provides(
			newRunner,
			func(r *runner) Runner { return r },
		)
	}),
	Stop: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(d *runner) error {
			return d.Close()
		})
	}),
}
