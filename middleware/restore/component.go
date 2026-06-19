package restore

import (
	"github.com/Compogo/compogo"
	"github.com/Compogo/compogo/flag"
	"github.com/Compogo/runner"
)

// Component — компонент middleware восстановления для runner.Runner.
var Component = compogo.Component{
	Name: "runner.middleware.restore",
	Dependencies: compogo.Components{
		&runner.Component,
	},
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provides(NewConfig, NewRestore)
	}),
	BindFlags: compogo.BindFlags(func(flagSet flag.FlagSet, container compogo.Container) error {
		return container.Invoke(func(config *Config) {
			flagSet.DurationVar(&config.Delay, DelayFieldName, DelayDefault, "")
		})
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(Configuration)
	}),
	PreExecute: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(func(r runner.Runner, middleware *Restore) {
			r.Use(middleware)
		})
	}),
}
