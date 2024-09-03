package responses

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func AppendRcodeToResponse(state request.Request, value json.RawMessage, response *PolicyResponse) error {
	var s string
	if err := json.Unmarshal(value, &s); err != nil {
		return err
	}

	rcode, err := StringToRcode(s)
	if err != nil {
		return err
	}

	response.Rcode = &rcode
	return nil
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
