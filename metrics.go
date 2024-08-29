package rpz

import (
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsSubsystem = "rpz"

	QueryStatusError       = "ERROR"
	QueryStatusDeny        = "DENY"
	QueryStatusFallthrough = "FALLTHROUGH"
	QueryStatusSuccess     = "SUCCESS"
	QueryStatusNoMatch     = "NOMATCH"
)

var metricsRpzRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "request_duration_seconds",
	Help:      "Histogram of the time (in seconds) each RPZ request took.",
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

var metricsPolicyExecutionTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "policy_execution_time_seconds",
	Help:      "Histogram of the time (in seconds) each policy execution took.",
	Buckets:   []float64{.0001, .0005, .001, .005, .01, .05, .1, .5, 1},
}, []string{"policy"})

func MetricPolicyExecutionTime(policy string, duration float64) {
	p := strings.ReplaceAll(strings.ToLower(policy), " ", "_")
	metricsPolicyExecutionTime.WithLabelValues(p).Observe(duration)
}

var metricsTriggerMatchCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "trigger_match_total",
	Help:      "Count of trigger matches per policy.",
}, []string{"policy", "trigger"})

func MetricTriggerMatchCount(policy, trigger string) {
	p := strings.ReplaceAll(strings.ToLower(policy), " ", "_")
	metricsTriggerMatchCount.WithLabelValues(p, trigger).Inc()
}
