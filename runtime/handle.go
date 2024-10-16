package runtime

import (
	"context"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/matches"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
	"github.com/mwantia/coredns-consulrpz-plugin/responses"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func HandlePolicyResponse(ctx context.Context, state request.Request, server string, policy policies.Policy) (*responses.PolicyResponse, error) {
	tracer := otel.Tracer("coredns.otel")
	ctx, span := tracer.Start(ctx, "HandlePolicyResponse",
		trace.WithAttributes(
			attribute.String("policy.name", policy.Name),
			attribute.Int("policy.priority", policy.GetPriority()),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	for _, rule := range policy.Rules {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			response, err := HandlePolicyResponseRule(state, ctx, server, policy, rule)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				return nil, err
			}

			if response != nil {
				span.SetAttributes(
					attribute.Bool("policy.response.deny", response.Deny),
					attribute.Bool("dns.response.fallthrough", response.Fallthrough),
				)

				if response.Rcode != nil {
					span.SetAttributes(
						attribute.Int("dns.response.rcode", int(*response.Rcode)),
					)
				}

				return response, nil
			}
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
