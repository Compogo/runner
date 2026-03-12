package recover

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/runner"
)

// Component is a ready-to-use Compogo component that automatically adds
// panic recovery middleware to the application's runner.
//
// It depends on runner.Component and registers itself during the PreRun phase.
//
// Usage:
//
//	compogo.WithComponents(
//	    runner.Component,
//	    recover.Component,  // ← protects all tasks from panics
//	)
var Component = &component.Component{
	Name: "runner.recover",
	Dependencies: component.Components{
		runner.Component,
	},
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provide(NewRecover)
	}),
	PreExecute: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(r runner.Runner, middleware *Recover) {
			r.Use(middleware)
		})
	}),
}
