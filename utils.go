package rpz

import (
	"fmt"
	"sort"
	"strings"

	"github.com/miekg/dns"
)

func PrepareResponseRcode(r *dns.Msg, rcode int, recursionAvailable bool) *dns.Msg {
	m := new(dns.Msg)
	m.SetRcode(r, rcode)
	m.Authoritative = true
	m.RecursionAvailable = recursionAvailable

	return m
}

func PrepareResponseReply(r *dns.Msg, recursionAvailable bool) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.RecursionAvailable = recursionAvailable

	return m
}

func GetPriority(p *int) int {
	if p == nil {
		return DefaultPolicyPriority
	}
	return *p
}

func SortPolicies(policies []Policy) {
	sort.Slice(policies, func(i, j int) bool {
		ipriority := GetPriority(policies[i].Priority)
		jpriority := GetPriority(policies[j].Priority)

		if ipriority != jpriority {
			return ipriority < jpriority
		}

		return len(policies[i].Rules) < len(policies[j].Rules)
	})
	for _, policy := range policies {
		SortPolicyRules(policy.Rules)
	}
}

func SortPolicyRules(rules []PolicyRule) {
	sort.Slice(rules, func(i, j int) bool {
		ipriority := GetPriority(rules[i].Priority)
		jpriority := GetPriority(rules[j].Priority)

		if ipriority != jpriority {
			return ipriority < jpriority
		}

		return len(rules[i].Triggers) < len(rules[j].Triggers)
	})
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
