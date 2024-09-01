package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/mwantia/coredns-consulrpz-plugin/policies"
	"github.com/mwantia/coredns-consulrpz-plugin/triggers"
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

func HandleResponse(state request.Request, ctx context.Context, r *dns.Msg, rule policies.PolicyRule) (*Response, error) {
	response := &Response{}
	for _, action := range rule.Actions {
		if handled, err := HandleActionResponse(state, action, response); handled || err != nil {
			return response, err
		}
	}

	return response, nil
}

func HandleActionResponse(state request.Request, action policies.RuleAction, response *Response) (bool, error) {
	alias := action.GetAliasType()

	switch alias {
	case "deny":
		response.Deny = true
		return true, nil

	case "fallthrough":
		response.Fallthrough = true
		return true, nil

	case "code":
		if err := response.AppendRcode(state, action); err != nil {
			return false, err
		}

	case "record":
		if err := response.AppendRecord(state, action); err != nil {
			return false, err
		}
	}

	return false, nil
}

func HandleTrigger(state request.Request, ctx context.Context, trigger policies.RuleTrigger) (bool, error) {
	alias := trigger.GetAliasType()

	switch alias {
	case "type":
		if data, ok := trigger.Data.(triggers.QTypeData); ok {
			return triggers.MatchQTypeTrigger(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "cidr":
		if data, ok := trigger.Data.(triggers.CidrData); ok {
			return triggers.MatchCidrTrigger(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "name":
		if data, ok := trigger.Data.(triggers.QNameData); ok {
			return triggers.MatchQNameTrigger(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "time":
		if data, ok := trigger.Data.(triggers.TimeData); ok {
			return triggers.MatchTimeTrigger(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "cron":
		if data, ok := trigger.Data.(triggers.CronData); ok {
			return triggers.MatchCronTrigger(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "regex":
		if data, ok := trigger.Data.(triggers.RegexData); ok {
			return triggers.MatchRegexTrigger(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)
	}

	return true, nil // Return true, so any type that doesn't match will be "skipped"
}

func (r *Response) AppendRcode(state request.Request, action policies.RuleAction) error {
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

func (r *Response) AppendRecord(state request.Request, action policies.RuleAction) error {
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

func StringToRcode(s string) (uint16, error) {
	switch strings.ToUpper(s) {
	case "NOERROR":
		return dns.RcodeSuccess, nil
	case "FORMERR":
		return dns.RcodeFormatError, nil
	case "SERVFAIL":
		return dns.RcodeServerFailure, nil
	case "NXDOMAIN":
		return dns.RcodeNameError, nil
	case "NOTIMP":
		return dns.RcodeNotImplemented, nil
	case "REFUSED":
		return dns.RcodeRefused, nil
	case "YXDOMAIN":
		return dns.RcodeYXDomain, nil
	case "YXRRSET":
		return dns.RcodeYXRrset, nil
	case "NXRRSET":
		return dns.RcodeNXRrset, nil
	case "NOTAUTH":
		return dns.RcodeNotAuth, nil
	case "NOTZONE":
		return dns.RcodeNotZone, nil
	case "BADSIG", "BADVERS":
		return dns.RcodeBadVers, nil
	case "BADKEY":
		return dns.RcodeBadKey, nil
	case "BADTIME":
		return dns.RcodeBadTime, nil
	case "BADMODE":
		return dns.RcodeBadMode, nil
	case "BADNAME":
		return dns.RcodeBadName, nil
	case "BADALG":
		return dns.RcodeBadAlg, nil
	case "BADTRUNC":
		return dns.RcodeBadTrunc, nil
	default:
		return 0, fmt.Errorf("unknown rcode: %s", s)
	}
}
