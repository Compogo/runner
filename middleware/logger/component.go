package logger

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
)

// Component — компонент middleware логирования для runner.Runner.
var Component = compogo.Component{
	Name: "runner.middleware.logger",
	Dependencies: compogo.Components{
		&runner.Component,
	},
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provide(NewLogger)
	}),
	PreExecute: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(r runner.Runner, middleware *Logger) {
			r.Use(middleware)
		})
	}),
}
