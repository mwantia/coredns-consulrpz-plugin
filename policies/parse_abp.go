package policies

import (
	"bufio"
	"strings"

	"github.com/mwantia/coredns-consulrpz-plugin/logging"
)

func ParsePolicyAbpRule(policy *Policy) error {
	reader, err := GetPolicyTargetReader(policy.Target)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		logging.Log.Debug(line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
