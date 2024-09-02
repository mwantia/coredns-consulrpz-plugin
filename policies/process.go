package policies

import "github.com/mwantia/coredns-consulrpz-plugin/matches"

func (policy Policy) ProcessPolicyData() error {
	for i := range policy.Rules {
		for j := range policy.Rules[i].Matches {
			trigger := &policy.Rules[i].Matches[j]
			var err error

			if trigger.Data, err = trigger.ProcessData(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (match RuleMatch) ProcessData() (interface{}, error) {
	alias := match.GetAliasType()
	switch alias {
	case "type":
		return matches.ProcessQTypeData(match.Value)

	case "cidr":
		return matches.ProcessCidrData(match.Value)

	case "name":
		return matches.ProcessQNameData(match.Value)

	case "external":
		return matches.ProcessExternalData(match.Value)

	case "time":
		return matches.ProcessTimeData(match.Value)

	case "cron":
		return matches.ProcessCronData(match.Value)

	case "regex":
		return matches.ProcessRegexData(match.Value)
	}

	return nil, nil
}
