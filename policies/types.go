package policies

import "encoding/json"

const CurrentPolicyVersion = "1.0"
const DefaultPolicyPriority = 1000

type Policy struct {
	Name     string       `json:"name"`
	Version  string       `json:"version"`
	Priority *int         `json:"priority,omitempty"`
	Rules    []PolicyRule `json:"rules"`
	Target   string       `json:"target,omitempty"`
	Type     string       `json:"type,omitempty"`
	Hash     string       `json:"-"`
}

type PolicyRule struct {
	Priority *int          `json:"priority,omitempty"`
	Triggers []RuleTrigger `json:"triggers"`
	Actions  []RuleAction  `json:"actions"`
}

type RuleTrigger struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
	Data  interface{}     `json:"-"`
}

type RuleAction struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
}
