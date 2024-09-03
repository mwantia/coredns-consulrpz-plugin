package responses

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func AppendRecordToResponse(state request.Request, value json.RawMessage, response *PolicyResponse) error {
	var record PolicyRecord
	if err := json.Unmarshal(value, &record); err != nil {
		return err
	}

	ttl := record.GetRecordTTL()

	for _, rec := range record.Records {
		switch strings.ToUpper(rec.Type) {
		case "A":
			if err := AppendARecordToResponse(state, response, ttl, rec.Value); err != nil {
				return err
			}
		case "AAAA":
			if err := AppendAAAARecordToResponse(state, response, ttl, rec.Value); err != nil {
				return err
			}
		}
	}

	return nil
}

func AppendARecordToResponse(state request.Request, response *PolicyResponse, ttl int, value json.RawMessage) error {
	var addresses []string
	if err := json.Unmarshal(value, &addresses); err != nil {
		return err
	}

	qname := dns.Fqdn(state.Name())
	for _, address := range addresses {
		rr := CreateDnsRecord(qname, dns.TypeA, uint32(ttl), address)
		response.Records = append(response.Records, rr)
	}

	return nil
}

func AppendAAAARecordToResponse(state request.Request, response *PolicyResponse, ttl int, value json.RawMessage) error {
	var addresses []string
	if err := json.Unmarshal(value, &addresses); err != nil {
		return err
	}

	qname := dns.Fqdn(state.Name())
	for _, address := range addresses {
		rr := CreateDnsRecord(qname, dns.TypeAAAA, uint32(ttl), address)
		response.Records = append(response.Records, rr)
	}

	return nil
}

func CreateDnsRecord(qname string, qtype uint16, ttl uint32, address string) *dns.A {
	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   qname,
			Rrtype: qtype,
			Class:  dns.ClassINET,
			Ttl:    ttl,
		},
		A: net.ParseIP(address),
	}
}

func (record PolicyRecord) GetRecordTTL() int {
	if record.TTL != nil {
		return *record.TTL
	}

	return 3600 // Default TTL
}
