package runtime

import (
	"context"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
)

func HandlePoliciesSequence(state request.Request, ctx context.Context, request *dns.Msg, p []policies.Policy) (*policies.Policy, *Response, error) {
	for _, policy := range p {
		start := time.Now()
		response, err := HandlePolicy(state, ctx, request, policy)
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
