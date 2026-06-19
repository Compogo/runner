package metric

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
)

// Component — компонент middleware метрик для runner.Runner.
var Component = compogo.Component{
	Name: "runner.middleware.metric",
	Dependencies: compogo.Components{
		&runner.Component,
	},
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provide(NewMetric)
	}),
	PreExecute: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(r runner.Runner, middleware *Metric) {
			r.Use(middleware)
		})
	}),
}
