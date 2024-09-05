package policies

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

var MatchTypeAliasMap = map[string]string{
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

var ResponseTypeAliasMap = map[string]string{
	"rcode":   "code",
	"allow":   "fallthrough",
	"logging": "log",

	"any":              "inaddr_any",
	"loopback":         "inaddr_loopback",
	"broadcast":        "inaddr_broadcast",
	"inaddr6_any":      "inaddr_any",
	"inaddr6_loopback": "inaddr_loopback",
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

	return &policy, nil
}

func (policy *Policy) ValidatePolicy() (bool, error) {
	if policy == nil || len(policy.Name) <= 0 {
		return false, fmt.Errorf("policy name '%s' is undefined or invalid", policy.Name)
	}

	if policy.Version != CurrentPolicyVersion {
		return false, fmt.Errorf("policy version '%s' does not match the current version '%s'", policy.Version, CurrentPolicyVersion)
	}

	rules := make([]PolicyRule, 0)
	for _, rule := range policy.Rules {
		if len(rule.Matches) > 0 && len(rule.Responses) > 0 {
			rules = append(rules, rule)
		}
	}

	if len(rules) <= 0 {
		return false, fmt.Errorf("policy has none or only empty rules")
	}

	policy.Rules = rules
	return true, nil
}

func CalculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (match RuleMatch) GetAliasType() string {
	t := strings.ReplaceAll(strings.ToLower(match.Type), " ", "-")
	if alias, exist := MatchTypeAliasMap[t]; exist {
		return alias
	}

	return t
}

func (response RuleResponse) GetAliasType() string {
	t := strings.ToLower(response.Type)
	if alias, exist := ResponseTypeAliasMap[t]; exist {
		return alias
	}

	return t
}
