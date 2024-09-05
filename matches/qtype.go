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

func MatchQType(state request.Request, ctx context.Context, data QTypeData) (*MatchResult, error) {
	qtype := state.QType()

	for _, t := range data.Types {
		switch qtype {
		case dns.TypeA:
			return &MatchResult{
				Handled: t == "A",
				Data:    t,
			}, nil

		case dns.TypeAAAA:
			return &MatchResult{
				Handled: t == "AAAA",
				Data:    t,
			}, nil

		case dns.TypeCNAME:
			return &MatchResult{
				Handled: t == "CNAME",
				Data:    t,
			}, nil

		case dns.TypeHTTPS:
			return &MatchResult{
				Handled: t == "HTTPS",
				Data:    t,
			}, nil

		case dns.TypeTXT:
			return &MatchResult{
				Handled: t == "TXT",
				Data:    t,
			}, nil

		case dns.TypeSOA:
			return &MatchResult{
				Handled: t == "SOA",
				Data:    t,
			}, nil

		case dns.TypeNS:
			return &MatchResult{
				Handled: t == "NS",
				Data:    t,
			}, nil
		}
	}

	return nil, nil
}
