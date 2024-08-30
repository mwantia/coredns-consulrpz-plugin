package rpz

import (
	"context"
	"fmt"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-rpz-plugin/policies"
	"github.com/mwantia/coredns-rpz-plugin/triggers"
)

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
