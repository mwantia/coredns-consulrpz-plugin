package policies

import "sort"

func (policy Policy) SortPolicy() {
	for _, rule := range policy.Rules {
		sort.Slice(rule.Matches, func(i, j int) bool {
			ipriority := rule.Matches[i].GetPriority()
			jpriority := rule.Matches[j].GetPriority()

			return ipriority < jpriority
		})
		sort.Slice(rule.Responses, func(i, j int) bool {
			ipriority := rule.Responses[i].GetPriority()
			jpriority := rule.Responses[j].GetPriority()

			return ipriority < jpriority
		})
	}

	sort.Slice(policy.Rules, func(i, j int) bool {
		ipriority := policy.Rules[i].GetPriority()
		jpriority := policy.Rules[j].GetPriority()

		if ipriority != jpriority {
			return ipriority < jpriority
		}

		return len(policy.Rules[i].Matches) < len(policy.Rules[j].Matches)
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

func (policy *Policy) GetPriority() int {
	if policy == nil || policy.Priority == nil {
		return DefaultPolicyPriority
	}
	return *policy.Priority
}

func (rule *PolicyRule) GetPriority() int {
	if rule == nil || rule.Priority == nil {
		return DefaultPolicyPriority
	}
	return *rule.Priority
}

func (match *RuleMatch) GetPriority() int {
	if match != nil {
		alias := match.GetAliasType()
		switch alias {
		case "type":
			return 0
		case "cidr":
			return 1
		case "name":
			return 2
		case "external":
			return 3
		case "time":
			return 4
		case "cron":
			return 5
		case "regex":
			return 6
		}
	}
	return 1000
}

func (response *RuleResponse) GetPriority() int {
	if response != nil {
		alias := response.GetAliasType()
		switch alias {
		case "deny":
			return 0
		case "fallthrough":
			return 1
		case "code":
			return 2
		case "inaddr_any":
			return 3
		case "inaddr_loopback":
			return 4
		case "inaddr_broadcast":
			return 5
		case "record":
			return 6
		case "log":
			return 7
		}
	}
	return 1000
}
