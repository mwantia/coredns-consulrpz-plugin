package rpz

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

const CurrentPolicyVersion = "1.0"
const DefaultPolicyPriority = 1000

type Policy struct {
	Name             string       `json:"name"`
	Version          string       `json:"version"`
	Priority         *int         `json:"priority,omitempty"`
	AdaptivePriority *int         `json:"-"`
	Rules            []PolicyRule `json:"rules"`
	Hash             string       `json:"-"`
}

type PolicyRule struct {
	Priority *int          `json:"priority,omitempty"`
	Triggers []RuleTrigger `json:"triggers"`
	Actions  []RuleAction  `json:"actions"`
}

type RuleTrigger struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
}

type RuleAction struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
}

func (p *Policy) ValidatePolicy() bool {
	if p == nil || len(p.Name) <= 0 {
		return false
	}

	if p.Version != CurrentPolicyVersion {
		return false
	}

	rules := make([]PolicyRule, 0)
	for _, r := range p.Rules {
		if len(r.Triggers) > 0 && len(r.Actions) > 0 {
			rules = append(rules, r)
		}
	}

	if len(rules) <= 0 {
		return false
	}

	p.Rules = rules
	return true
}

func ParsePolicyFile(reader io.Reader) (*Policy, error) {
	var policy Policy

	buffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buffer, &policy)
	if err != nil {
		return nil, err
	}

	if !policy.ValidatePolicy() {
		return nil, fmt.Errorf("unable to validate policy")
	}

	policy.Hash = CalculateHash(buffer)
	return &policy, nil
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
	for _, policy := range policies {
		SortPolicyRules(policy.Rules)
	}
}

func SortPolicyRules(rules []PolicyRule) {
	sort.Slice(rules, func(i, j int) bool {
		ipriority := rules[i].GetPriority()
		jpriority := rules[j].GetPriority()

		if ipriority != jpriority {
			return ipriority < jpriority
		}

		return len(rules[i].Triggers) < len(rules[j].Triggers)
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
