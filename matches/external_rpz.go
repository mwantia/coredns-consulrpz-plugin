package matches

import (
	"bufio"
	"io"
	"strings"

	"github.com/mwantia/coredns-consulrpz-plugin/data"
)

func ParseRpzMatches(reader io.Reader, trie *data.OldTrie) error {
	scanner := bufio.NewScanner(reader)

	var names []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) > 2 && parts[1] == "CNAME" && parts[2] == "." {
			names = append(names, parts[0])
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for _, name := range names {
		trie.Insert(name)
	}

	return nil
}
