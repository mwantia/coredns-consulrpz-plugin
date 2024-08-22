package rpz

import (
	"encoding/json"
	"fmt"
	"io"
)

type Policy struct {
	Name     string       `json:"name"`
	Version  string       `json:"version"`
	Priority int          `json:"priority,omitempty"`
	Rules    []PolicyRule `json:"rules"`
}

type PolicyRule struct {
	Priority int           `json:"priority,omitempty"`
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

	if p.Version != "1.0" {
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
		return nil, fmt.Errorf("unable to validate")
	}

	return &policy, nil
}
