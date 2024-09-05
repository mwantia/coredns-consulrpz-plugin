package runtime

import (
	"context"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/matches"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
	"github.com/mwantia/coredns-consulrpz-plugin/responses"
)

func HandlePolicyResponse(state request.Request, ctx context.Context, server string, policy policies.Policy) (*responses.PolicyResponse, error) {
	for _, rule := range policy.Rules {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				if response, err := HandlePolicyResponseRule(state, ctx, server, policy, rule); response != nil || err != nil {
					return response, err
				}
			}
		} else if response, err := HandlePolicyResponseRule(state, ctx, server, policy, rule); response != nil || err != nil {
			return response, err
		}
	}

	return nil, nil
}

func HandlePolicyResponseRule(state request.Request, ctx context.Context, server string, policy policies.Policy, rule policies.PolicyRule) (*responses.PolicyResponse, error) {
	var err error
	var result *matches.MatchResult

	for _, match := range rule.Matches {
		alias := match.GetAliasType()

		if ctx != nil {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				result, err = matches.HandleMatches(state, ctx, alias, match.Data)
				if err != nil {
					return nil, err
				}
				if result != nil && !result.Handled {
					return nil, nil
				}
			}
		} else {
			result, err = matches.HandleMatches(state, ctx, alias, match.Data)
			if err != nil {
				return nil, err
			}
			if result != nil && !result.Handled {
				return nil, nil
			}
		}
	}

	if ctx != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if response, err := responses.HandleResponse(state, ctx, server, result, policy, rule); response != nil || err != nil {
				return response, err
			}
		}
	} else if response, err := responses.HandleResponse(state, ctx, server, result, policy, rule); response != nil || err != nil {
		return response, err
	}

	return nil, nil
}
