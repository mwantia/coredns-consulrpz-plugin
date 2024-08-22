package triggers

import (
	"encoding/json"
	"regexp"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func MatchDomainTrigger(state request.Request, value json.RawMessage) (bool, error) {
	var domains []string
	if err := json.Unmarshal(value, &domains); err != nil {
		return false, err
	}

	qname := dns.Fqdn(state.Name())

	for _, d := range domains {
		if qname == dns.Fqdn(d) {
			return true, nil
		}

		match, err := regexp.MatchString(d, qname)
		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}

	return false, nil
}
