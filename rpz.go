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
	start := time.Now()

	qtype := state.QType()

	policy, response, err := p.HandlePoliciesParallel(state, ctx, r)

	if err != nil {
		logging.Log.Errorf("Unable to handle request for '%s': %s", dns.Fqdn(state.Name()), err)
	}

	if policy != nil && response != nil {
		if response.Fallthrough {
			duration := time.Since(start).Seconds()

			MetricRequestDurationSeconds(StatusFallthrough, duration)
			MetricQueryRequestsTotal(StatusFallthrough, policy.Name, qtype)

			return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
		}
		if response.Deny {
			duration := time.Since(start).Seconds()

			MetricRequestDurationSeconds(StatusDeny, duration)
			MetricQueryRequestsTotal(StatusDeny, policy.Name, qtype)

			return HandleDenyAll(state)
		}

		msg := PrepareResponseReply(state.Req, true)
		if response.Rcode != nil {
			msg.Rcode = int(*response.Rcode)
		}
		msg.SetReply(r)
		msg.Answer = response.Answers

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
			logging.Log.Errorf("Unable to send response: %s", err)

			duration := time.Since(start).Seconds()

			MetricRequestDurationSeconds(StatusError, duration)
			MetricQueryRequestsTotal(StatusError, policy.Name, qtype)

			return dns.RcodeServerFailure, err
		}

		duration := time.Since(start).Seconds()

		MetricRequestDurationSeconds(StatusSuccess, duration)
		MetricQueryRequestsTotal(StatusSuccess, policy.Name, qtype)

		return msg.Rcode, nil
	}
	duration := time.Since(start).Seconds()

	MetricRequestDurationSeconds(StatusNoMatch, duration)
	MetricQueryRequestsTotal(StatusNoMatch, "", qtype)

	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}
