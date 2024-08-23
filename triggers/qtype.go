package triggers

import (
	"encoding/json"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func MatchQTypeTrigger(state request.Request, value json.RawMessage) (bool, error) {
	var types []string
	if err := json.Unmarshal(value, &types); err != nil {
		return false, err
	}

	qtype := state.QType()

	for _, t := range types {
		t = strings.ToUpper(t)
		switch qtype {
		case dns.TypeA:
			return t == "A", nil

		case dns.TypeAAAA:
			return t == "AAAA", nil

		case dns.TypeCNAME:
			return t == "CNAME", nil

		case dns.TypeHTTPS:
			return t == "HTTPS", nil

		case dns.TypeTXT:
			return t == "TXT", nil

		case dns.TypeSOA:
			return t == "SOA", nil

		case dns.TypeNS:
			return t == "NS", nil
		}
	}

	return false, nil
}
