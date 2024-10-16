package consulrpz

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/coredns/coredns/plugin"
	cmetrics "github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
	"github.com/mwantia/coredns-consulrpz-plugin/responses"
	"github.com/mwantia/coredns-consulrpz-plugin/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (plug ConsulRpzPlugin) Name() string { return "consulrpz" }

func (p ConsulRpzPlugin) ServeDNS(ctx context.Context, writer dns.ResponseWriter, msg *dns.Msg) (int, error) {
	state := request.Request{W: writer, Req: msg.Copy()}

	tracer := otel.Tracer("coredns.otel")
	ctx, span := tracer.Start(ctx, p.Name(),
		trace.WithAttributes(
			attribute.String("plugin.name", p.Name()),
		),
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	status, err := p.ServeDnsRequest(ctx, state)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return status, err
}

func (p ConsulRpzPlugin) ServeDnsRequest(ctx context.Context, state request.Request) (int, error) {
	var policy *policies.Policy
	var response *responses.PolicyResponse
	var err error

	start := time.Now()
	execution := strings.ToLower(p.Cfg.Execution)
	switch execution {
	case "parallel":
		// policy, response, err = runtime.HandlePoliciesParallel(state, ctx, p.Policies)
	case "sequence":
		policy, response, err = runtime.HandlePoliciesSequence(state, ctx, p.Policies)
	}
	duration := time.Since(start).Seconds()

	if err != nil && !errors.Is(err, context.Canceled) {
		logging.Log.Errorf("Unable to handle request for '%s': %s", dns.Fqdn(state.Name()), err)

		p.SetQueryStatus(ctx, state.QType(), metrics.QueryStatusError, duration, policy)
		return dns.RcodeServerFailure, err
	}

	if policy == nil || response == nil {
		p.SetQueryStatus(ctx, state.QType(), metrics.QueryStatusNoMatch, duration, policy)
		return p.HandleNextOrFailure(ctx, state)
	}

	if response.Fallthrough {
		p.SetQueryStatus(ctx, state.QType(), metrics.QueryStatusFallthrough, duration, policy)
		return p.HandleNextOrFailure(ctx, state)
	}

	if response.Deny {
		p.SetQueryStatus(ctx, state.QType(), metrics.QueryStatusDeny, duration, policy)
		return HandleDenyPolicy(state, *policy)
	}

	responsemsg := PrepareResponseReply(state.Req, true)
	if response.Rcode != nil {
		responsemsg.Rcode = int(*response.Rcode)
	}
	responsemsg.SetReply(state.Req)
	responsemsg.Answer = response.Records
	responses.WriteExtraHandle(responsemsg, state, response.Extra)

	if response.Rcode != nil {
		responsemsg.Rcode = int(*response.Rcode)
	} else {
		if len(responsemsg.Answer) > 0 {
			responsemsg.Rcode = dns.RcodeSuccess
		} else {
			responsemsg.Rcode = dns.RcodeNameError
		}
	}

	if err := state.W.WriteMsg(responsemsg); err != nil {
		logging.Log.Errorf("Unable to send response for '%s': %s", dns.Fqdn(state.Name()), err)

		p.SetQueryStatus(ctx, state.QType(), metrics.QueryStatusError, duration, policy)
		return dns.RcodeServerFailure, err
	}

	p.SetQueryStatus(ctx, state.QType(), metrics.QueryStatusSuccess, duration, policy)
	return responsemsg.Rcode, nil
}

func (p ConsulRpzPlugin) HandleNextOrFailure(ctx context.Context, state request.Request) (int, error) {
	tracer := otel.Tracer("coredns.otel")
	ctx, span := tracer.Start(ctx, "NextOrFailure",
		trace.WithAttributes(
			attribute.String("plugin.name", p.Name()),
			attribute.String("plugin.next", p.Next.Name()),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	status, err := plugin.NextOrFailure(p.Name(), p.Next, ctx, state.W, state.Req)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return status, err
}

func (plug ConsulRpzPlugin) SetQueryStatus(ctx context.Context, qtype uint16, status string, duration float64, policy *policies.Policy) {
	name := ""
	if policy != nil {
		name = policy.Name
	}

	server := cmetrics.WithServer(ctx)
	metrics.MetricsRequestDurationSeconds(server, status, duration)
	metrics.MetricsQueryRequests(server, status, name, qtype)

	plug.SetMetadataQueryStatus(ctx, status)
}
