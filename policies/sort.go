package policies

import "sort"

func SortPolicies(policies []Policy) {
	sort.Slice(policies, func(i, j int) bool {
		ipriority := policies[i].GetPriority()
		jpriority := policies[j].GetPriority()

		if ipriority != jpriority {
			return ipriority < jpriority
		}

		return len(policies[i].Rules) < len(policies[j].Rules)
	})

	for _, policy := range policies {
		sort.Slice(policy.Rules, func(i, j int) bool {
			ipriority := policy.Rules[i].GetPriority()
			jpriority := policy.Rules[j].GetPriority()

			if ipriority != jpriority {
				return ipriority < jpriority
			}

			return len(policy.Rules[i].Triggers) < len(policy.Rules[j].Triggers)
		})

		for _, rule := range policy.Rules {
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
	}
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
