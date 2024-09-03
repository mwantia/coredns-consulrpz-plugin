package consulrpz

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

func PrepareResponseRcode(request *dns.Msg, rcode int, recursionAvailable bool) *dns.Msg {
	m := new(dns.Msg)
	m.SetRcode(request, rcode)
	m.Authoritative = false
	m.RecursionAvailable = recursionAvailable

	return m
}

func PrepareResponseReply(request *dns.Msg, recursionAvailable bool) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(request)
	m.Authoritative = false
	m.RecursionAvailable = recursionAvailable
	return m
}

func GetFileNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

func CalculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
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
