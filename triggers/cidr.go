package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-rpz-plugin/logging"
)

func MatchCidrTrigger(state request.Request, ctx context.Context, value json.RawMessage) (bool, error) {
	var cidrs []string
	if err := json.Unmarshal(value, &cidrs); err != nil {
		return false, err
	}

	ip := state.IP()
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false, fmt.Errorf("unable to parse client IP '%s'", ip)
	}

	for _, cidr := range cidrs {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			logging.Log.Debugf("Checking cidr '%s' with client '%s'", cidr, clientIP)
			// Simplest check and should always be tried first
			if ip == cidr {
				return true, nil
			}

			_, ipnet, err := net.ParseCIDR(cidr)
			if err != nil {
				return false, nil // Ignore parse errors
			}

			if ipnet.Contains(clientIP) {
				return true, nil
			}
		}
	}

	return false, nil
}
