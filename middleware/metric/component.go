package metric

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/runner"
)

// Component is a ready-to-use Compogo component that automatically adds
// Prometheus metrics middleware to the application's runner.
//
// It depends on runner.Component and registers itself during the PreRun phase.
// The exported metrics follow the naming convention:
//
//	compogo_runner_task{app="<app-name>"}
//
// Usage:
//
//	compogo.WithComponents(
//	    runner.Component,
//	    metric.Component,  // ← adds Prometheus metrics to all tasks
//	)
var Component = &component.Component{
	Dependencies: component.Components{
		runner.Component,
	},
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provide(NewMetric)
	}),
	PreRun: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(r runner.Runner, middleware *Metric) {
			r.Use(middleware)
		})
	}),
}
