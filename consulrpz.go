package consulrpz

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
	"github.com/mwantia/coredns-consulrpz-plugin/metrics"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
	"github.com/mwantia/coredns-consulrpz-plugin/responses"
	"github.com/mwantia/coredns-consulrpz-plugin/runtime"
)

func (plug ConsulRpzPlugin) Name() string { return "consulrpz" }

func (plug ConsulRpzPlugin) ServeDNS(ctx context.Context, writer dns.ResponseWriter, msg *dns.Msg) (int, error) {
	state := request.Request{W: writer, Req: msg}
	qtype := state.QType()
	qname := dns.Fqdn(state.Name())

	var policy *policies.Policy
	var response *responses.PolicyResponse
	var err error

	start := time.Now()
	execution := strings.ToLower(plug.Cfg.Execution)
	switch execution {
	case "parallel":
		policy, response, err = runtime.HandlePoliciesParallel(state, ctx, plug.Policies)
	case "sequence":
		policy, response, err = runtime.HandlePoliciesSequence(state, ctx, plug.Policies)
	}
	duration := time.Since(start).Seconds()

	if err != nil && !errors.Is(err, context.Canceled) {
		logging.Log.Errorf("Unable to handle request for '%s': %s", qname, err)

		plug.SetQueryStatus(ctx, qtype, metrics.QueryStatusError, duration, policy)
		return dns.RcodeServerFailure, err
	}

	if policy == nil || response == nil {
		plug.SetQueryStatus(ctx, qtype, metrics.QueryStatusNoMatch, duration, policy)
		return plugin.NextOrFailure("consulrpz", plug.Next, ctx, writer, msg)
	}

	if response.Fallthrough {
		plug.SetQueryStatus(ctx, qtype, metrics.QueryStatusFallthrough, duration, policy)
		return plugin.NextOrFailure("consulrpz", plug.Next, ctx, writer, msg)
	}

	if response.Deny {
		plug.SetQueryStatus(ctx, qtype, metrics.QueryStatusDeny, duration, policy)
		return HandleDenyPolicy(state, *policy)
	}

	responsemsg := PrepareResponseReply(state.Req, true)
	if response.Rcode != nil {
		responsemsg.Rcode = int(*response.Rcode)
	}
	responsemsg.SetReply(msg)
	responsemsg.Answer = response.Records
	WriteExtraPolicyHandle(responsemsg, state, *policy)

	if response.Rcode != nil {
		responsemsg.Rcode = int(*response.Rcode)
	} else {
		if len(responsemsg.Answer) > 0 {
			responsemsg.Rcode = dns.RcodeSuccess
		} else {
			responsemsg.Rcode = dns.RcodeNameError
		}
	}

	if err := writer.WriteMsg(responsemsg); err != nil {
		logging.Log.Errorf("Unable to send response for '%s': %s", qname, err)

		plug.SetQueryStatus(ctx, qtype, metrics.QueryStatusError, duration, policy)
		return dns.RcodeServerFailure, err
	}

	plug.SetQueryStatus(ctx, qtype, metrics.QueryStatusSuccess, duration, policy)
	return responsemsg.Rcode, nil
}

func (plug ConsulRpzPlugin) SetQueryStatus(ctx context.Context, qtype uint16, status string, duration float64, policy *policies.Policy) {
	name := ""
	if policy != nil {
		name = policy.Name
	}

	metrics.MetricRequestDurationSeconds(status, duration)
	metrics.MetricQueryRequestsTotal(status, name, qtype)

	plug.SetMetadataQueryStatus(ctx, status)
}
