package runtime

import (
	"context"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
	"github.com/mwantia/coredns-consulrpz-plugin/responses"
)

func HandlePoliciesSequence(state request.Request, ctx context.Context, policies []policies.Policy) (*policies.Policy, *responses.PolicyResponse, error) {
	for _, policy := range policies {

		start := time.Now()
		response, err := HandlePolicyResponse(state, nil, policy)
		duration := time.Since(start).Seconds()

		metrics.MetricPolicyExecutionTime(policy.Name, duration)

		if err != nil {
			return &policy, nil, err
		}

		if response != nil {
			return &policy, response, nil
		}
	}

	return nil, nil, nil
}
