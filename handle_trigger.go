package rpz

import (
	"context"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-rpz-plugin/triggers"
)

var TriggerTypeAliasMap = map[string]string{
	"name-regex":   "regex",
	"qname-regex":  "regex",
	"domain-regex": "regex",

	"qname":  "name",
	"domain": "name",

	"qtype": "type",

	"ip-address": "cidr",
	"ip-range":   "cidr",
	"client-ip":  "cidr",
}

func HandleTrigger(state request.Request, ctx context.Context, trigger RuleTrigger) (bool, error) {
	alias := trigger.GetAliasType()

	switch alias {
	case "type":
		return triggers.MatchQTypeTrigger(state, ctx, trigger.Value)

	case "cidr":
		return triggers.MatchCidrTrigger(state, ctx, trigger.Value)

	case "name":
		return triggers.MatchQNameTrigger(state, ctx, trigger.Value)

	case "time":
		return triggers.MatchTimeTrigger(state, ctx, trigger.Value)

	case "cron":
		return triggers.MatchCronTrigger(state, ctx, trigger.Value)

	case "regex":
		return triggers.MatchRegexTrigger(state, ctx, trigger.Value)
	}

	return true, nil // Return true, so any type that doesn't match will be "skipped"
}

func (trigger RuleTrigger) GetAliasType() string {
	t := strings.ReplaceAll(strings.ToLower(trigger.Type), " ", "-")
	if alias, exist := TriggerTypeAliasMap[t]; exist {
		return alias
	}

	return t
}
