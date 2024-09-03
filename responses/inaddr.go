package responses

import (
	"fmt"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func AppendInAddrToResponse(state request.Request, aliastype string, response *PolicyResponse) error {
	qname := dns.Fqdn(state.Name())
	qtype := state.QType()
	ttl := uint32(3600)

	var rr dns.RR

	switch aliastype {
	case "inaddr_any":
		if qtype == dns.TypeAAAA {
			rr = CreateDnsRecord(qname, dns.TypeAAAA, ttl, "::")
		} else {
			rr = CreateDnsRecord(qname, dns.TypeA, ttl, "0.0.0.0")
		}

	case "inaddr_loopback":
		if qtype == dns.TypeAAAA {
			rr = CreateDnsRecord(qname, dns.TypeAAAA, ttl, "::1")
		} else {
			rr = CreateDnsRecord(qname, dns.TypeA, ttl, "127.0.0.1")
		}

	case "inaddr_broadcast":
		if qtype == dns.TypeAAAA {
			return nil // Lets just ignore AAAA requests
		} else {
			rr = CreateDnsRecord(qname, dns.TypeA, ttl, "255.255.255.255")
		}
	}

	if rr != nil {
		response.Records = append(response.Records, rr)
		return nil
	}

	return fmt.Errorf("no matching inaddr with the type '%s' found", aliastype)
}
