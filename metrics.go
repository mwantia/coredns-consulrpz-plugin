package rpz

import (
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsSubsystem = "rpz"

	StatusError       = "ERROR"
	StatusDeny        = "DENY"
	StatusFallthrough = "FALLTHROUGH"
	StatusSuccess     = "SUCCESS"
	StatusNoMatch     = "NOMATCH"
)

var metricsRpzRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "request_duration_seconds",
	Help:      "Histogram of the time (in seconds) each request to Consul took.",
	Buckets:   []float64{.001, .002, .005, .01, .02, .05, .1, .2, .5, 1},
}, []string{"status"})

func MetricRequestDurationSeconds(status string, duration float64) {
	s := strings.ToUpper(status)
	metricsRpzRequestDurationSeconds.WithLabelValues(s).Observe(duration)
}

var metricsQueryRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "query_requests_total",
	Help:      "Count the amount of queries received as request by the plugin.",
}, []string{"status", "policy", "type"})

func MetricQueryRequestsTotal(status, policy string, qtype uint16) {
	t := dns.TypeToString[qtype]
	s := strings.ToUpper(status)
	p := strings.ReplaceAll(strings.ToLower(policy), " ", "_")
	metricsQueryRequestsTotal.WithLabelValues(s, p, t).Inc()
}
