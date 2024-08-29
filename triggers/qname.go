package triggers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func MatchQNameTrigger(state request.Request, ctx context.Context, value json.RawMessage) (bool, error) {
	var names []string
	if err := json.Unmarshal(value, &names); err != nil {
		return false, err
	}

	qname := dns.Fqdn(state.Name())

	for _, name := range names {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			suffix := dns.Fqdn(name)
			if strings.HasSuffix(qname, suffix) {
				return true, nil
			}
		}
	}

	return false, nil
}
