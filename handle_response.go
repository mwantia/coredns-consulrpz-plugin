package rpz

import (
	"context"
	"encoding/json"
	"net"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Response struct {
	Deny        bool
	Fallthrough bool
	Rcode       *uint16
	Answers     []dns.RR
}

type ResponseRecord struct {
	TTL     *int `json:"ttl"`
	Records []struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	} `json:"records"`
}

func HandleResponse(state request.Request, ctx context.Context, r *dns.Msg, rule PolicyRule) (*Response, error) {
	response := &Response{}

	for _, action := range rule.Actions {
		switch action.Type {
		case "deny":
			return &Response{
				Deny: true,
			}, nil

		case "fallthrough":
			return &Response{
				Fallthrough: true,
			}, nil

		case "rcode":
			if err := response.AppendRcode(state, action); err != nil {
				return nil, err
			}

		case "record":
			if err := response.AppendRecord(state, action); err != nil {
				return nil, err
			}
		}
	}

	return response, nil
}

func (r *Response) AppendRcode(state request.Request, action RuleAction) error {
	var s string
	if err := json.Unmarshal(action.Value, &s); err != nil {
		return err
	}

	rcode, err := StringToRcode(s)
	if err != nil {
		return err
	}

	r.Rcode = &rcode
	return nil
}

func (r *Response) AppendRecord(state request.Request, action RuleAction) error {
	var record ResponseRecord
	if err := json.Unmarshal(action.Value, &record); err != nil {
		return err
	}

	ttl := record.GetRecordTTL()

	for _, rec := range record.Records {
		switch rec.Type {
		case "A":
			if err := r.AppendARecords(state, ttl, rec.Value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Response) AppendARecords(state request.Request, ttl int, value json.RawMessage) error {
	var addresses []string
	if err := json.Unmarshal(value, &addresses); err != nil {
		return nil
	}

	qname := dns.Fqdn(state.Name())

	for _, address := range addresses {
		rr := &dns.A{
			Hdr: dns.RR_Header{
				Name:   qname,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(ttl),
			},
			A: net.ParseIP(address),
		}
		r.Answers = append(r.Answers, rr)
	}

	return nil
}

func (r ResponseRecord) GetRecordTTL() int {
	if r.TTL != nil {
		return *r.TTL
	}

	return 3600 // Default TTL
}
