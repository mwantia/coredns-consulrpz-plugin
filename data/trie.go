package data

import "strings"

type TrieNode struct {
	Children map[string]*TrieNode
	IsEnd    bool
}

type Trie struct {
	Root *TrieNode
}

func NewRootTrie() *Trie {
	return &Trie{
		Root: &TrieNode{
			Children: make(map[string]*TrieNode),
		},
	}
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
