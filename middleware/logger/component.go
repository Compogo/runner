package logger

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/runner"
)

// Component is a ready-to-use Compogo component that automatically adds
// task lifecycle logging middleware to the application's runner.
//
// It depends on runner.Component and registers itself during the PreRun phase.
// Logs will show when tasks start, complete, or encounter errors.
//
// Usage:
//
//	compogo.WithComponents(
//	    runner.Component,
//	    logger.Component,  // ← adds lifecycle logging to all tasks
//	)
var Component = &component.Component{
	Dependencies: component.Components{
		runner.Component,
	},
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provide(NewLogger)
	}),
	PreRun: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(r runner.Runner, middleware *Logger) {
			r.Use(middleware)
		})
	}),
}
