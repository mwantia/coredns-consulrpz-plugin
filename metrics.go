package rpz

import (
	"github.com/coredns/coredns/plugin"
	"github.com/prometheus/client_golang/prometheus"
)

var metricsSubsystem = "rpz"

var metricsRpzRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: plugin.Namespace,
	Subsystem: metricsSubsystem,
	Name:      "request_duration_seconds",
	Help:      "Histogram of the time (in seconds) each request to Consul took.",
	Buckets:   []float64{.001, .002, .005, .01, .02, .05, .1, .2, .5, 1},
}, []string{"status"})

func IncrementMetricsRpzRequestDurationSeconds(status string, duration float64) {
	metricsRpzRequestDurationSeconds.WithLabelValues(status).Observe(duration)
}
