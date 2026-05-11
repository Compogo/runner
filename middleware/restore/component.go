package restore

import (
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
	"github.com/Compogo/runner"
)

var (
	ComponentNoRegistration = &component.Component{
		Name: "runner.middleware.restore.no_registration",
		Dependencies: component.Components{
			runner.Component,
		},
		Init: component.StepFunc(func(container container.Container) error {
			return container.Provides(NewConfig, NewRestore)
		}),
		BindFlags: component.BindFlags(func(flagSet flag.FlagSet, container container.Container) error {
			return container.Invoke(func(config *Config) {
				flagSet.DurationVar(&config.Delay, DelayFieldName, DelayDefault, "")
			})
		}),
		Configuration: component.StepFunc(func(container container.Container) error {
			return container.Invoke(Configuration)
		}),
	}

	Component = &component.Component{
		Name: "runner.middleware.restore",
		Dependencies: component.Components{
			runner.Component,
		},
		Init: component.StepFunc(func(container container.Container) error {
			return container.Provides(NewConfig, NewRestore)
		}),
		BindFlags: component.BindFlags(func(flagSet flag.FlagSet, container container.Container) error {
			return container.Invoke(func(config *Config) {
				flagSet.DurationVar(&config.Delay, DelayFieldName, DelayDefault, "")
			})
		}),
		Configuration: component.StepFunc(func(container container.Container) error {
			if err := container.Invoke(Configuration); err != nil {
				return err
			}

			return container.Invoke(func(r runner.Runner, middleware *Restore) {
				r.Use(middleware)
			})
		}),
	}
)
