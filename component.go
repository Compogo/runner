package runner

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
)

// Component is a ready-to-use Compogo component that provides the Runner.
// It automatically:
//   - Registers NewRunner in the DI container
//   - Calls Close() on the runner during the Stop phase
//
// Usage:
//
//	compogo.WithComponents(runner.Component)
var Component = &component.Component{
	// Init registers the Runner constructor in the DI container.
	// Any component that needs a runner can request it via dependency injection.
	Init: component.StepFunc(func(container container.Container) error {
		return container.Provide(NewRunner)
	}),

	// Stop ensures all running tasks are gracefully shut down.
	// This runs during application shutdown, after all Wait components have
	// completed and before the process exits.
	Stop: component.StepFunc(func(container container.Container) error {
		return container.Invoke(func(r Runner) error {
			return r.Close()
		})
	}),
}
