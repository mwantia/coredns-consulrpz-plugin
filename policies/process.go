package policies

import "github.com/mwantia/coredns-rpz-plugin/triggers"

func ProcessPolicyData(policies []Policy) error {
	for i := range policies {
		for j := range policies[i].Rules {
			for k := range policies[i].Rules[j].Triggers {
				trigger := &policies[i].Rules[j].Triggers[k]
				var err error

				if trigger.Data, err = trigger.ProcessData(); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (t RuleTrigger) ProcessData() (interface{}, error) {
	alias := t.GetAliasType()
	switch alias {
	case "type":
		return triggers.ProcessQTypeData(t.Value)

	case "cidr":
		return triggers.ProcessCidrData(t.Value)

	case "name":
		return triggers.ProcessQNameData(t.Value)

	case "time":
		return triggers.ProcessTimeData(t.Value)

	case "cron":
		return triggers.ProcessCronData(t.Value)

	case "regex":
		return triggers.ProcessRegexData(t.Value)
	}

	return nil, nil
}
