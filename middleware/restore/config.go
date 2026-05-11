package restore

import (
	"time"

	"github.com/Compogo/compogo/configurator"
)

const (
	DelayFieldName = "runner.middleware.restore.delay"

	DelayDefault = 100 * time.Millisecond
)

type Config struct {
	Delay time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func Configuration(config *Config, configurator configurator.Configurator) *Config {
	if config.Delay == 0 || config.Delay == DelayDefault {
		configurator.SetDefault(DelayFieldName, DelayDefault)
		config.Delay = configurator.GetDuration(DelayFieldName)
	}

	return config
}
