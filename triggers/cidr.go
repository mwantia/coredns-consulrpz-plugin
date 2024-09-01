package triggers

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
)

type CidrData struct {
	Networks []net.IPNet
}

func ProcessCidrData(value json.RawMessage) (interface{}, error) {
	var cidrs []string
	if err := json.Unmarshal(value, &cidrs); err != nil {
		return nil, err
	}

	data := CidrData{}

	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}

		data.Networks = append(data.Networks, *network)
	}

	return data, nil
}

func MatchCidrTrigger(state request.Request, ctx context.Context, data CidrData) (bool, error) {
	ip := state.IP()
	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false, fmt.Errorf("unable to parse client IP '%s'", ip)
	}

	for _, network := range data.Networks {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
			logging.Log.Debugf("Checking cidr '%s' with client '%s'", network, clientIP)

			if network.Contains(clientIP) {
				return true, nil
			}
		}
	}

	return false, nil
}
