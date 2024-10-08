package metrics

import "github.com/prometheus/client_golang/prometheus"

func Register() error {
	prometheus.MustRegister(metricsRpzRequestDurationSeconds)
	prometheus.MustRegister(metricsQueryRequestsTotal)
	prometheus.MustRegister(metricsPolicyExecutionTime)
	prometheus.MustRegister(metricsPoliciesUpdatedTotal)
	prometheus.MustRegister(metricsPolicyResponsesTotal)
	return nil
}
