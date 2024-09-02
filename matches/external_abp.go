package matches

import (
	"bufio"
	"io"
	"strings"

	"github.com/mwantia/coredns-consulrpz-plugin/data"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
)

func ParseAbpMatches(reader io.Reader, trie *data.OldTrie) error {
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
