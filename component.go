package runner

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
)

var Component = &component.Component{
	Name: "runner",
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(
			newRunner,
			func(r *runner) Runner { return r },
		)
	}),
	Stop: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(d *runner) error {
			return d.Close()
		})
	}),
}
