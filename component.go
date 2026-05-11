package runner

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
)

var Component = &component.Component{
	Name: "runner",
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provides(
			NewDefault,
			func(d *Default) Runner { return d },
		)
	}),
	Stop: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(d *Default) error {
			return d.Close()
		})
	}),
}
