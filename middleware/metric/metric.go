package metric

import (
	"context"

	"github.com/Compogo/compogo"
	"github.com/Compogo/runner"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metric is a middleware that increments a Prometheus counter for each
// task execution. It helps track how many tasks are being run and can be
// used for monitoring task activity.
type Metric struct {
	counter prometheus.Counter
}

// NewMetric creates a new Metric middleware instance.
// The appConfig provides the application name which is added as a label
// to the counter, allowing per-application metrics in multi-service environments.
//
// The exported metric follows the naming convention:
//
//	compogo_runner_task{app="<app-name>"}
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

// Middleware wraps a task's process function and increments the counter
// each time the task is executed. The counter is incremented regardless
// of whether the task succeeds or fails.
func (m *Metric) Middleware(task *runner.Task, next runner.Process) runner.Process {
	return runner.ProcessFunc(func(ctx context.Context) error {
		m.counter.Inc()

		return next.Process(ctx)
	})
}
