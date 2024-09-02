package matches

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type QTypeData struct {
	Types []string
}

func ProcessQTypeData(value json.RawMessage) (interface{}, error) {
	var types []string
	if err := json.Unmarshal(value, &types); err != nil {
		return nil, err
	}

	data := QTypeData{}

	for _, t := range types {
		data.Types = append(data.Types, strings.ToUpper(t))
	}

	return data, nil
}

func MatchQType(state request.Request, ctx context.Context, data QTypeData) (bool, error) {
	qtype := state.QType()

	for _, t := range data.Types {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
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
	}

	return false, nil
}
