package policies

import "github.com/mwantia/coredns-consulrpz-plugin/triggers"

func (p Policy) ProcessPolicyData() error {
	for i := range p.Rules {
		for j := range p.Rules[i].Triggers {
			trigger := &p.Rules[i].Triggers[j]
			var err error

			if trigger.Data, err = trigger.ProcessData(); err != nil {
				return err
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
