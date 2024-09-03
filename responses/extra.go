package responses

import (
	"encoding/json"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func AppendExtraToResponse(state request.Request, value json.RawMessage, response *PolicyResponse) error {
	var extra []string
	if err := json.Unmarshal(value, &extra); err != nil {
		return err
	}

	response.Extra = append(response.Extra, extra...)
	return nil
}

func WriteExtraHandle(msg *dns.Msg, state request.Request, extra []string) {
	if len(extra) > 0 {
		qname := dns.Fqdn(state.Name())
		msg.Extra = append(msg.Extra, &dns.TXT{
			Hdr: dns.RR_Header{Name: qname, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0},
			Txt: extra,
		})
	}
}
