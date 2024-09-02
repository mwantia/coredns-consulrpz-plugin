package data

import (
	"github.com/miekg/dns"
)

type TrieNode struct {
	Children map[byte]*TrieNode
	IsEnd    bool
}

type Trie struct {
	Root *TrieNode
}

func NewRootTrie() *Trie {
	return &Trie{
		Root: &TrieNode{
			Children: make(map[byte]*TrieNode),
		},
	}
}

func (t *Trie) Insert(name string) {
	name = dns.Fqdn(name)
	node := t.Root
	for i := len(name) - 1; i >= 0; i-- {
		char := name[i]
		if char == '.' {
			continue
		}
		if _, exists := node.Children[char]; !exists {
			node.Children[char] = &TrieNode{
				Children: make(map[byte]*TrieNode),
			}
		}
		node = node.Children[char]
	}
	node.IsEnd = true
}

func (t *Trie) Search(name string) bool {
	name = dns.Fqdn(name)
	node := t.Root
	for i := len(name) - 1; i >= 0; i-- {
		char := name[i]
		if char == '.' {
			continue
		}
		if _, exists := node.Children[char]; !exists {
			return false
		}
		node = node.Children[char]
		if node.IsEnd {
			return true
		}
	}
	return false
}
