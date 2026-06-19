package recover

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
)

// Component — компонент middleware восстановления после паник для runner.Runner.
var Component = &compogo.Component{
	Name: "runner.middleware.recover",
	Dependencies: compogo.Components{
		&runner.Component,
	},
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provide(NewRecover)
	}),
	PreExecute: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(r runner.Runner, middleware *Recover) {
			r.Use(middleware)
		})
	}),
}
