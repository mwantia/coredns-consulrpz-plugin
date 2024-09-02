package matches

import (
	"context"
	"fmt"

	"github.com/coredns/coredns/request"
)

func HandleMatches(state request.Request, ctx context.Context, alias string, data interface{}) (bool, error) {
	switch alias {
	case "type":
		if data, ok := data.(QTypeData); ok {
			return MatchQType(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "cidr":
		if data, ok := data.(CidrData); ok {
			return MatchCidr(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "name":
		if data, ok := data.(QNameData); ok {
			return MatchQName(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "external":
		if data, ok := data.(ExternalData); ok {
			return MatchExternal(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "time":
		if data, ok := data.(TimeData); ok {
			return MatchTime(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "cron":
		if data, ok := data.(CronData); ok {
			return MatchCron(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)

	case "regex":
		if data, ok := data.(RegexData); ok {
			return MatchRegex(state, ctx, data)
		}
		return false, fmt.Errorf("unable to process trigger data as '%s'", alias)
	}

	return true, nil // Return true, so any type that doesn't match will be "skipped"
}
