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
		if handled, err := triggers.MatchQTypeTrigger(state, trigger.Value); handled || err != nil {
			return handled, err
		}

	case "cidr":
		if handled, err := triggers.MatchCidrTrigger(state, trigger.Value); handled || err != nil {
			return handled, err
		}

	case "name":
		if handled, err := triggers.MatchQNameTrigger(state, trigger.Value); handled || err != nil {
			return handled, err
		}
	}

	return false, nil
}

func (trigger RuleTrigger) GetAliasType() string {
	t := strings.ToLower(trigger.Type)
	if alias, exist := TriggerTypeAliasMap[t]; exist {
		return alias
	}

	return t
}
