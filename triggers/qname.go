package triggers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type QNameData struct {
	Trie *Trie
}

type TrieNode struct {
	Children map[string]*TrieNode
	IsEnd    bool
}

type Trie struct {
	Root *TrieNode
}

func ProcessQNameData(value json.RawMessage) (interface{}, error) {
	var names []string
	if err := json.Unmarshal(value, &names); err != nil {
		return nil, err
	}

	data := QNameData{
		Trie: &Trie{
			Root: &TrieNode{
				Children: make(map[string]*TrieNode),
			},
		},
	}

	for _, name := range names {
		name = dns.Fqdn(name)
		data.Trie.Insert(name)
	}

	return data, nil
}

func MatchQNameTrigger(state request.Request, ctx context.Context, data QNameData) (bool, error) {
	name := dns.Fqdn(state.Name())
	ok := data.Trie.Search(name)
	return ok, nil
}

func (t *Trie) Insert(name string) {
	node := t.Root
	parts := strings.Split(strings.TrimSuffix(name, "."), ".")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if _, exists := node.Children[part]; !exists {
			node.Children[part] = &TrieNode{
				Children: make(map[string]*TrieNode),
			}
		}
		node = node.Children[part]
	}
	node.IsEnd = true
}

func (t *Trie) Search(name string) bool {
	node := t.Root
	parts := strings.Split(strings.TrimSuffix(name, "."), ".")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if _, exists := node.Children[part]; !exists {
			return false
		}
		node = node.Children[part]
		if node.IsEnd {
			return true
		}
	}
	return false
}
