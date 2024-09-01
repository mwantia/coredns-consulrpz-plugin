package policies

import (
	"bufio"
	"strings"
)

func ParsePolicyHostsRule(policy *Policy) error {
	reader, err := GetPolicyTargetReader(policy.Target)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
