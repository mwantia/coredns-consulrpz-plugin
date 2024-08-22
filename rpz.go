package rpz

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-rpz-plugin/logging"
	"github.com/mwantia/coredns-rpz-plugin/triggers"
)

func (p RpzPlugin) Name() string { return "rpz" }

func (p RpzPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	start := time.Now()

	for _, policy := range p.Config.Policies {
		response, err := p.HandlePolicy(state, ctx, r, policy)
		if err != nil {
			logging.Log.Errorf("Unable to handle request for '%s': %s", dns.Fqdn(state.Name()), err)
		}

		if response != nil {
			if response.Fallthrough {
				duration := time.Since(start).Seconds()
				IncrementMetricsRpzRequestDurationSeconds("FALLTHROUGH", duration)

				return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
			}
			if response.Deny {
				duration := time.Since(start).Seconds()
				IncrementMetricsRpzRequestDurationSeconds("DENY", duration)

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
				return dns.RcodeServerFailure, err
			}

			duration := time.Since(start).Seconds()
			IncrementMetricsRpzRequestDurationSeconds("SUCCESS", duration)
			return msg.Rcode, nil
		}
	}
	duration := time.Since(start).Seconds()
	IncrementMetricsRpzRequestDurationSeconds("NOMATCH", duration)
	// Pass onto the next plugin (fallthrough) by default
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}

func (p RpzPlugin) HandlePolicy(state request.Request, ctx context.Context, r *dns.Msg, policy Policy) (*Response, error) {
	logging.Log.Debugf("Handling policy named '%s'", policy.Name)

	for _, rule := range policy.Rules {
		if response, err := p.HandlePolicyRule(state, ctx, r, rule); response != nil || err != nil {
			return response, err
		}
	}

	return nil, nil
}

func (p RpzPlugin) HandlePolicyRule(state request.Request, ctx context.Context, r *dns.Msg, rule PolicyRule) (*Response, error) {
	logging.Log.Debugf("Handling policy rule with '%v' triggers and '%v' actions", len(rule.Triggers), len(rule.Actions))

	for _, trigger := range rule.Triggers {
		globalmatch := false
		switch trigger.Type {
		case "domain":
			match, err := triggers.MatchDomainTrigger(state, trigger.Value)
			globalmatch = match

			if err != nil {
				return nil, err
			}
		}

		if globalmatch {
			if response, err := HandleResponse(state, ctx, r, rule); response != nil || err != nil {
				return response, err
			}
		}
	}

	return nil, nil
}
