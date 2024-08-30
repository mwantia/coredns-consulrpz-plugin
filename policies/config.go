package policies

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

var TriggerTypeAliasMap = map[string]string{
	"name-regex":   "regex",
	"qname-regex":  "regex",
	"domain-regex": "regex",

	"qname":  "name",
	"domain": "name",

	"qtype": "type",

	"ip-address": "cidr",
	"ip-range":   "cidr",
	"client-ip":  "cidr",
}

var ActionTypeAliasMap = map[string]string{
	"rcode": "code",
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

func CalculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (trigger RuleTrigger) GetAliasType() string {
	t := strings.ReplaceAll(strings.ToLower(trigger.Type), " ", "-")
	if alias, exist := TriggerTypeAliasMap[t]; exist {
		return alias
	}

	return t
}

func (action RuleAction) GetAliasType() string {
	t := strings.ToLower(action.Type)
	if alias, exist := ActionTypeAliasMap[t]; exist {
		return alias
	}

	return t
}
