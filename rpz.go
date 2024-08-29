package rpz

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-rpz-plugin/logging"
)

func (p RpzPlugin) Name() string { return "rpz" }

func (p RpzPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qtype := state.QType()
	qname := dns.Fqdn(state.Name())

	start := time.Now()
	policy, response, err := p.HandlePoliciesParallel(state, ctx, r)
	duration := time.Since(start).Seconds()

	if err != nil {
		logging.Log.Errorf("Unable to handle request for '%s': %s", qname, err)

		p.SetQueryStatus(ctx, qtype, QueryStatusError, duration, policy)
		return dns.RcodeServerFailure, err
	}

	if policy == nil || response == nil {
		p.SetQueryStatus(ctx, qtype, QueryStatusNoMatch, duration, policy)
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}

	if response.Fallthrough {
		p.SetQueryStatus(ctx, qtype, QueryStatusFallthrough, duration, policy)
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}

	if response.Deny {
		p.SetQueryStatus(ctx, qtype, QueryStatusDeny, duration, policy)
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

		p.SetQueryStatus(ctx, qtype, QueryStatusError, duration, policy)
		return dns.RcodeServerFailure, err
	}

	p.SetQueryStatus(ctx, qtype, QueryStatusSuccess, duration, policy)
	return msg.Rcode, nil
}

func (p RpzPlugin) SetQueryStatus(ctx context.Context, qtype uint16, status string, duration float64, policy *Policy) {
	name := ""
	if policy != nil {
		name = policy.Name
	}

	MetricRequestDurationSeconds(status, duration)
	MetricQueryRequestsTotal(status, name, qtype)

	p.SetMetadataQueryStatus(ctx, status)
}
