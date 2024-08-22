package rpz

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func HandleError(state request.Request, rcode int, e error) (int, error) {
	msg := PrepareResponseRcode(state.Req, rcode, true)
	if err := state.W.WriteMsg(msg); err != nil {
		return dns.RcodeServerFailure, err
	}

	return rcode, e
}

func HandleDenyAll(state request.Request) (int, error) {
	msg := PrepareResponseRcode(state.Req, dns.RcodeRefused, false)
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
