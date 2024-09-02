package matches

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/data"
	"github.com/mwantia/coredns-consulrpz-plugin/logging"
)

type ExternalData struct {
	Trie *data.OldTrie
}

func ProcessExternalData(value json.RawMessage) (interface{}, error) {
	var targets []struct {
		Target  string `json:"target"`
		Type    string `json:"type"`
		Refresh string `json:"refresh"`
	}
	if err := json.Unmarshal(value, &targets); err != nil {
		return nil, err
	}

	trie := data.OldRootTrie()

	for _, target := range targets {
		// We skip targets that can't be resolved
		if len(target.Target) <= 0 {
			logging.Log.Warning("Defined external target is empty or undefined")
			continue
		}

		reader, err := GetTargetReader(target.Target)
		if err != nil {
			logging.Log.Warningf("Unable to load target: %v", err)
			continue
		}

		switch strings.ToLower(target.Type) {
		case "hosts":
			if err := ParseHostsMatches(reader, trie); err != nil {
				logging.Log.Warningf("Unable to parse to '%s': %v", target.Type, err)
				continue
			}
		case "rpz":
			if err := ParseRpzMatches(reader, trie); err != nil {
				logging.Log.Warningf("Unable to parse to '%s': %v", target.Type, err)
				continue
			}
		case "abp":
			if err := ParseAbpMatches(reader, trie); err != nil {
				logging.Log.Warningf("Unable to parse to '%s': %v", target.Type, err)
				continue
			}
		}
	}

	return ExternalData{
		Trie: trie,
	}, nil
}

func MatchExternal(state request.Request, ctx context.Context, data ExternalData) (bool, error) {
	name := state.Name()
	ok := data.Trie.Search(name)
	return ok, nil
}

func GetTargetReader(target string) (io.Reader, error) {
	if strings.HasPrefix(target, "https://") || strings.HasPrefix(target, "http://") {
		response, err := http.Get(target)
		if err != nil {
			return nil, err
		}

		return response.Body, nil
	}

	return os.Open(target)
}
