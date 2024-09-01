package policies

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/mwantia/coredns-consulrpz-plugin/logging"
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

	policy.Hash = CalculateHash(buffer)
	if len(policy.Type) > 0 && len(policy.Target) > 0 {
		logging.Log.Infof("Parsing policy target: '%s' as '%s'", policy.Target, policy.Type)
		switch strings.ToLower(policy.Type) {
		case "hosts":
			if err := ParsePolicyHostsRule(&policy); err != nil {
				return nil, err
			}
		case "rpz":
			if err := ParsePolicyRpzRule(&policy); err != nil {
				return nil, err
			}
		case "abp":
			if err := ParsePolicyAbpRule(&policy); err != nil {
				return nil, err
			}
		}
	}

	return &policy, nil
}

func (p *Policy) ValidatePolicy() (bool, error) {
	if p == nil || len(p.Name) <= 0 {
		return false, fmt.Errorf("policy name '%s' is undefined or invalid", p.Name)
	}

	if p.Version != CurrentPolicyVersion {
		return false, fmt.Errorf("policy version '%s' does not match the current version '%s'", p.Version, CurrentPolicyVersion)
	}

	rules := make([]PolicyRule, 0)
	for _, r := range p.Rules {
		if len(r.Triggers) > 0 && len(r.Actions) > 0 {
			rules = append(rules, r)
		}
	}

	if len(rules) <= 0 {
		return false, fmt.Errorf("policy has none or only empty rules")
	}

	p.Rules = rules
	return true, nil
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
