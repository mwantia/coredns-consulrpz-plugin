package runtime

import (
	"context"
	"time"

	cmetrics "github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
	"github.com/mwantia/coredns-consulrpz-plugin/responses"
)

func HandlePoliciesSequence(state request.Request, ctx context.Context, policies []policies.Policy) (*policies.Policy, *responses.PolicyResponse, error) {
	server := cmetrics.WithServer(ctx)

	for _, policy := range policies {
		if policy.Disabled {
			logging.Log.Debugf("Policy '%s' is disabled and will be skipped", policy.Name)
			continue
		}

		start := time.Now()
		response, err := HandlePolicyResponse(state, nil, server, policy)
		duration := time.Since(start).Seconds()

		metrics.MetricPolicyExecutionTime(server, policy.Name, duration)

		if err != nil {
			return &policy, nil, err
		}

		if response != nil {
			return &policy, response, nil
		}
	}

	return nil, nil, nil
}
