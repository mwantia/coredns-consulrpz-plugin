package policies

import (
	"bufio"
	"encoding/json"
	"strings"
)

func ParsePolicyRpzRule(policy *Policy) error {
	reader, err := GetPolicyTargetReader(policy.Target)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)

	var domains []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) > 2 && parts[1] == "CNAME" && parts[2] == "." {
			domains = append(domains, parts[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if len(domains) > 0 {
		json, err := json.Marshal(domains)
		if err != nil {
			return err
		}

		policy.Rules = append(policy.Rules, PolicyRule{
			Triggers: []RuleTrigger{
				{
					Type:  "name",
					Value: json,
				},
			},
			Actions: []RuleAction{
				{
					Type: "deny",
				},
			},
		})
	}

	return nil
}
