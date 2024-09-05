package responses

import (
	"context"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/matches"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
)

func HandleResponse(state request.Request, ctx context.Context, server string, result *matches.MatchResult, policy policies.Policy, rule policies.PolicyRule) (*PolicyResponse, error) {
	response := &PolicyResponse{}
	for _, rr := range rule.Responses {
		alias := rr.GetAliasType()

		switch alias {
		case "deny":
			response.Deny = true
			metrics.MetricsPolicyResponses(server, alias)

		case "fallthrough":
			response.Fallthrough = true
			metrics.MetricsPolicyResponses(server, alias)

		case "extra":
			if err := AppendExtraToResponse(state, rr.Value, response); err != nil {
				return response, err
			}
			metrics.MetricsPolicyResponses(server, alias)

		case "code":
			if err := AppendRcodeToResponse(state, rr.Value, response); err != nil {
				return response, err
			}
			metrics.MetricsPolicyResponses(server, alias)

		case "record":
			if err := AppendRecordToResponse(state, rr.Value, response); err != nil {
				return response, err
			}
			metrics.MetricsPolicyResponses(server, alias)

		case "log":
			if err := HandleLogResponse(state, rr.Value, result, policy, response); err != nil {
				return response, err
			}
			metrics.MetricsPolicyResponses(server, alias)
		}

		if strings.HasPrefix(alias, "inaddr_") {
			if err := AppendInAddrToResponse(state, alias, response); err != nil {
				return response, err
			}
			metrics.MetricsPolicyResponses(server, alias)
		}
	}

	return response, nil
}
