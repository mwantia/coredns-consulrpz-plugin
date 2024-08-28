package rpz

import (
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-rpz-plugin/triggers"
)

var TriggerTypeAliasMap = map[string]string{
	"domain": "name",
	"qname":  "name",

	"qtype": "type",

	"ip-address": "cidr",
	"ip-range":   "cidr",
	"client-ip":  "cidr",
}

func HandleTrigger(state request.Request, trigger RuleTrigger) (bool, error) {
	alias := trigger.GetAliasType()

	switch alias {
	case "type":
		return triggers.MatchQTypeTrigger(state, trigger.Value)

	case "cidr":
		return triggers.MatchCidrTrigger(state, trigger.Value)

	case "name":
		return triggers.MatchQNameTrigger(state, trigger.Value)

	case "time":
		return triggers.MatchTimeTrigger(state, trigger.Value)

	case "cron":
		return triggers.MatchCronTrigger(state, trigger.Value)
	}

	return true, nil // Return true, so any type that doesn't match will be "skipped"
}

func (trigger RuleTrigger) GetAliasType() string {
	t := strings.ToLower(trigger.Type)
	if alias, exist := TriggerTypeAliasMap[t]; exist {
		return alias
	}

	return t
}
