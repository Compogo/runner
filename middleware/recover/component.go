package recover

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/runner"
)

var Component = &component.Component{
	Name: "runner.middleware.recover",
	Dependencies: component.Components{
		runner.Component,
	},
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provide(NewRecover)
	}),
	Configuration: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(r runner.Runner, middleware *Recover) {
			r.Use(middleware)
		})
	}),
}
