package matches

import (
	"context"
	"encoding/json"

	"github.com/coredns/coredns/request"
	"github.com/mwantia/coredns-consulrpz-plugin/data"
)

type QNameData struct {
	Trie *data.Trie
}

func ProcessQNameData(value json.RawMessage) (interface{}, error) {
	var names []string
	if err := json.Unmarshal(value, &names); err != nil {
		return nil, err
	}

	trie := data.NewRootTrie()

	for _, name := range names {
		trie.Insert(name)
	}

	return QNameData{
		Trie: trie,
	}, nil
}

func MatchQName(state request.Request, ctx context.Context, data QNameData) (bool, error) {
	name := state.Name()
	ok := data.Trie.Search(name)
	return ok, nil
}
