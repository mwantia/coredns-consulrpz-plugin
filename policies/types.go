package policies

import "encoding/json"

const CurrentPolicyVersion = "1.0"
const DefaultPolicyPriority = 1000

type Policy struct {
	Name     string       `json:"name"`
	Disabled bool         `json:"disabled"`
	Version  string       `json:"version"`
	Priority *int         `json:"priority,omitempty"`
	Rules    []PolicyRule `json:"rules"`
	Hash     string       `json:"-"`
}

type PolicyRule struct {
	Priority  *int           `json:"priority,omitempty"`
	Matches   []RuleMatch    `json:"matches"`
	Responses []RuleResponse `json:"responses"`
}

type RuleMatch struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
	Data  interface{}     `json:"-"`
}

type RuleResponse struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
}
