package metrics

import (
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
)

var metricsRpzRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "request_duration_seconds",
	Help:      "Histogram of the time (in seconds) each RPZ request took.",
	Buckets:   []float64{.001, .002, .005, .01, .02, .05, .1, .2, .5, 1},
}, []string{"server", "status"})

func MetricsRequestDurationSeconds(server, status string, duration float64) {
	s := strings.ToUpper(status)
	metricsRpzRequestDurationSeconds.WithLabelValues(server, s).Observe(duration)
}

var metricsQueryRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "query_requests_total",
	Help:      "Count the amount of queries received as request by the plugin.",
}, []string{"server", "status", "policy", "type"})

func MetricsQueryRequests(server, status, policy string, qtype uint16) {
	t := dns.TypeToString[qtype]
	s := strings.ToUpper(status)
	p := strings.ReplaceAll(strings.ToLower(policy), " ", "_")
	metricsQueryRequestsTotal.WithLabelValues(server, s, p, t).Inc()
}

var metricsPolicyExecutionTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "policy_execution_time_seconds",
	Help:      "Histogram of the time (in seconds) each policy execution took.",
	Buckets:   []float64{.0001, .0005, .001, .005, .01, .05, .1, .5, 1},
}, []string{"server", "policy"})

func MetricPolicyExecutionTime(server, policy string, duration float64) {
	p := strings.ReplaceAll(strings.ToLower(policy), " ", "_")
	metricsPolicyExecutionTime.WithLabelValues(server, p).Observe(duration)
}

var metricsPoliciesUpdatedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "policies_updated_total",
	Help:      "Count the amount of updates performed per policy.",
}, []string{"policy"})

func MetricsPoliciesUpdated(policy string) {
	p := strings.ReplaceAll(strings.ToLower(policy), " ", "_")
	metricsPoliciesUpdatedTotal.WithLabelValues(p).Inc()
}

var metricsPolicyResponsesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: MetricsSubsystem,
	Name:      "policy_responses_total",
	Help:      "Count the amount of policy responses by type.",
}, []string{"server", "type"})

func MetricsPolicyResponses(server, rtype string) {
	metricsPolicyResponsesTotal.WithLabelValues(server, rtype).Inc()
}
