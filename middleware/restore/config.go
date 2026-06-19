package restore

import (
	"time"

	"github.com/Compogo/compogo"
)

// DelayFieldName — имя поля для задержки между попытками.
const DelayFieldName = "runner.middleware.restore.delay"

// DelayDefault — задержка по умолчанию (100 мс).
var DelayDefault = 100 * time.Millisecond

// Config содержит конфигурацию middleware восстановления.
type Config struct {
	Delay time.Duration
}

// NewConfig создаёт новую конфигурацию.
func NewConfig() *Config {
	return &Config{}
}

// Configuration загружает конфигурацию из Configurator.
func Configuration(config *Config, configurator compogo.Configurator) *Config {
	if config.Delay == 0 || config.Delay == DelayDefault {
		configurator.SetDefault(DelayFieldName, DelayDefault)
		config.Delay = configurator.GetDuration(DelayFieldName)
	}

	return config
}
