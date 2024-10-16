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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func HandlePoliciesSequence(state request.Request, ctx context.Context, policies []policies.Policy) (*policies.Policy, *responses.PolicyResponse, error) {
	tracer := otel.Tracer("coredns.otel")
	ctx, span := tracer.Start(ctx, "HandlePoliciesSequence",
		trace.WithAttributes(
			attribute.Int("policy.count", len(policies)),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	server := cmetrics.WithServer(ctx)

	for _, policy := range policies {
		if policy.Disabled {
			logging.Log.Debugf("Policy '%s' is disabled and will be skipped", policy.Name)
			continue
		}

		start := time.Now()
		response, err := HandlePolicyResponse(ctx, state, server, policy)
		duration := time.Since(start).Seconds()

		metrics.MetricPolicyExecutionTime(server, policy.Name, duration)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return &policy, nil, err
		}

		if response != nil {
			return &policy, response, nil
		}
	}

	return nil, nil, nil
}
