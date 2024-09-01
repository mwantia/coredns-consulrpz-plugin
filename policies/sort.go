package policies

import "sort"

func (p Policy) SortPolicy() {
	for _, rule := range p.Rules {
		sort.Slice(rule.Triggers, func(i, j int) bool {
			ipriority := rule.Triggers[i].GetPriority()
			jpriority := rule.Triggers[j].GetPriority()

			return ipriority < jpriority
		})
		sort.Slice(rule.Actions, func(i, j int) bool {
			ipriority := rule.Actions[i].GetPriority()
			jpriority := rule.Actions[j].GetPriority()

			return ipriority < jpriority
		})
	}

	sort.Slice(p.Rules, func(i, j int) bool {
		ipriority := p.Rules[i].GetPriority()
		jpriority := p.Rules[j].GetPriority()

		if ipriority != jpriority {
			return ipriority < jpriority
		}

		return len(p.Rules[i].Triggers) < len(p.Rules[j].Triggers)
	})
}

func SortPolicies(policies []Policy) {
	sort.Slice(policies, func(i, j int) bool {
		ipriority := policies[i].GetPriority()
		jpriority := policies[j].GetPriority()

		if ipriority != jpriority {
			return ipriority < jpriority
		}

		return len(policies[i].Rules) < len(policies[j].Rules)
	})
}

func (p *Policy) GetPriority() int {
	if p == nil || p.Priority == nil {
		return DefaultPolicyPriority
	}
	return *p.Priority
}

func (r *PolicyRule) GetPriority() int {
	if r == nil || r.Priority == nil {
		return DefaultPolicyPriority
	}
	return *r.Priority
}

func (t *RuleTrigger) GetPriority() int {
	if t != nil {
		alias := t.GetAliasType()
		switch alias {
		case "type":
			return 0
		case "name":
			return 2
		case "cidr":
			return 1
		case "time":
			return 3
		case "cron":
			return 4
		case "regex":
			return 5
		}
	}
	return 1000
}

func (a *RuleAction) GetPriority() int {
	if a != nil {
		alias := a.GetAliasType()
		switch alias {
		case "deny":
			return 0
		case "fallthrough":
			return 1
		case "code":
			return 2
		case "record":
			return 3
		}
	}
	return 1000
}
