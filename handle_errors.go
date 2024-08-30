package rpz

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-rpz-plugin/policies"
)

func HandleError(state request.Request, rcode int, e error) (int, error) {
	msg := PrepareResponseRcode(state.Req, rcode, true)
	if err := state.W.WriteMsg(msg); err != nil {
		return dns.RcodeServerFailure, err
	}

	return rcode, e
}

func HandleDenyPolicy(state request.Request, policy policies.Policy) (int, error) {
	msg := PrepareResponseRcode(state.Req, dns.RcodeRefused, false)
	WriteExtraPolicyHandle(msg, state, policy)
	if err := state.W.WriteMsg(msg); err != nil {
		return dns.RcodeServerFailure, err
	}

	return dns.RcodeRefused, nil
}

func HandleNXDomain(state request.Request) (int, error) {
	msg := PrepareResponseRcode(state.Req, dns.RcodeNameError, false)
	if err := state.W.WriteMsg(msg); err != nil {
		return dns.RcodeServerFailure, err
	}

	return dns.RcodeNameError, nil
}
