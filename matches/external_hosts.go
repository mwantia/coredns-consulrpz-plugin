package matches

import (
	"bufio"
	"io"
	"strings"

	"github.com/mwantia/coredns-consulrpz-plugin/data"
)

func ParseHostsMatches(reader io.Reader, trie *data.OldTrie) error {
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
