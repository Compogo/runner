package metric

import (
	"context"

	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metric — middleware для сбора метрик Prometheus.
// Считает количество запусков процессов.
type Metric struct {
	counter prometheus.Counter
}

func NewMetric(appConfig *compogo.Config) *Metric {
	return &Metric{
		counter: promauto.NewCounter(prometheus.CounterOpts{
			Name: compogo.MetricNamePrefix + "runner_task",
			Help: "number of running tasks",
			ConstLabels: map[string]string{
				compogo.MetricAppNameFieldName: appConfig.Name,
			},
		}),
	}
}

func (m *Metric) Middleware(_ runner.Process, next runner.ProcessFunc) runner.ProcessFunc {
	return func(ctx context.Context) error {
		m.counter.Inc()

		return next(ctx)
	}
}
