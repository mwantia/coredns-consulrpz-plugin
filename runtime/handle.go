package runtime

import (
	"context"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/matches"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
	"github.com/mwantia/coredns-consulrpz-plugin/responses"
)

func HandlePolicyResponse(state request.Request, ctx context.Context, policy policies.Policy) (*responses.PolicyResponse, error) {
	for _, rule := range policy.Rules {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if response, err := HandlePolicyResponseRule(state, ctx, policy, rule); response != nil || err != nil {
					return response, err
				}
			}
		} else if response, err := HandlePolicyResponseRule(state, ctx, policy, rule); response != nil || err != nil {
			return response, err
		}
	}

	return nil, nil
}

func HandlePolicyResponseRule(state request.Request, ctx context.Context, policy policies.Policy,
	rule policies.PolicyRule) (*responses.PolicyResponse, error) {
	for _, match := range rule.Matches {
		alias := match.GetAliasType()

		if ctx != nil {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if handled, err := matches.HandleMatches(state, ctx, alias, match.Data); !handled || err != nil {
					return nil, err
				}

				metrics.MetricTriggerMatchCount(policy.Name, match.GetAliasType())
			}
		} else {
			if handled, err := matches.HandleMatches(state, ctx, alias, match.Data); !handled || err != nil {
				return nil, err
			}

			metrics.MetricTriggerMatchCount(policy.Name, alias)
		}
	}

	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if response, err := responses.HandleResponse(state, ctx, rule); response != nil || err != nil {
				return response, err
			}
		}
	} else if response, err := responses.HandleResponse(state, ctx, rule); response != nil || err != nil {
		return response, err
	}

	return nil, nil
}
