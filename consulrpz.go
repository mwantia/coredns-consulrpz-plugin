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
	"github.com/mwantia/coredns-consulrpz-plugin/runtime"
)

func (p ConsulRpzPlugin) Name() string { return "consulrpz" }

func (p ConsulRpzPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qtype := state.QType()
	qname := dns.Fqdn(state.Name())

	var policy *policies.Policy
	var response *runtime.Response
	var err error

	start := time.Now()
	switch strings.ToLower(p.Cfg.Execution) {
	case "parallel":
		policy, response, err = runtime.HandlePoliciesParallel(state, ctx, r, p.Policies)
	case "sequence":
		policy, response, err = runtime.HandlePoliciesSequence(state, ctx, r, p.Policies)
	}
	duration := time.Since(start).Seconds()

	if err != nil && !errors.Is(err, context.Canceled) {
		logging.Log.Errorf("Unable to handle request for '%s': %s", qname, err)

		p.SetQueryStatus(ctx, qtype, metrics.QueryStatusError, duration, policy)
		return dns.RcodeServerFailure, err
	}

	if policy == nil || response == nil {
		p.SetQueryStatus(ctx, qtype, metrics.QueryStatusNoMatch, duration, policy)
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}

	if response.Fallthrough {
		p.SetQueryStatus(ctx, qtype, metrics.QueryStatusFallthrough, duration, policy)
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}

	if response.Deny {
		p.SetQueryStatus(ctx, qtype, metrics.QueryStatusDeny, duration, policy)
		return HandleDenyPolicy(state, *policy)
	}

	msg := PrepareResponseReply(state.Req, true)
	if response.Rcode != nil {
		msg.Rcode = int(*response.Rcode)
	}
	msg.SetReply(r)
	msg.Answer = response.Answers
	WriteExtraPolicyHandle(msg, state, *policy)

	if response.Rcode != nil {
		msg.Rcode = int(*response.Rcode)
	} else {
		if len(msg.Answer) > 0 {
			msg.Rcode = dns.RcodeSuccess
		} else {
			msg.Rcode = dns.RcodeNameError
		}
	}

	if err := w.WriteMsg(msg); err != nil {
		logging.Log.Errorf("Unable to send response for '%s': %s", qname, err)

		p.SetQueryStatus(ctx, qtype, metrics.QueryStatusError, duration, policy)
		return dns.RcodeServerFailure, err
	}

	p.SetQueryStatus(ctx, qtype, metrics.QueryStatusSuccess, duration, policy)
	return msg.Rcode, nil
}

func (p ConsulRpzPlugin) SetQueryStatus(ctx context.Context, qtype uint16, status string, duration float64, policy *policies.Policy) {
	name := ""
	if policy != nil {
		name = policy.Name
	}

	metrics.MetricRequestDurationSeconds(status, duration)
	metrics.MetricQueryRequestsTotal(status, name, qtype)

	p.SetMetadataQueryStatus(ctx, status)
}
