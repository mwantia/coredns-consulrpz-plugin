package responses

import (
	"encoding/json"

	"github.com/miekg/dns"
)

type PolicyResponse struct {
	Deny        bool     `json:"deny"`
	Fallthrough bool     `json:"fallthrough"`
	Rcode       *uint16  `json:"rcode"`
	Records     []dns.RR `json:"records"`
	Extras      []string `json:"extras"`
}

type PolicyRecord struct {
	TTL     *int `json:"ttl"`
	Records []struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	} `json:"records"`
}
