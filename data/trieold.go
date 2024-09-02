package data

import "strings"

type OldTrieNode struct {
	Children map[string]*OldTrieNode
	IsEnd    bool
}

type OldTrie struct {
	Root *OldTrieNode
}

func OldRootTrie() *OldTrie {
	return &OldTrie{
		Root: &OldTrieNode{
			Children: make(map[string]*OldTrieNode),
		},
	}
}

func (t *OldTrie) Insert(name string) {
	node := t.Root
	parts := strings.Split(strings.TrimSuffix(name, "."), ".")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if _, exists := node.Children[part]; !exists {
			node.Children[part] = &OldTrieNode{
				Children: make(map[string]*OldTrieNode),
			}
		}
		node = node.Children[part]
	}
	node.IsEnd = true
}

func (t *OldTrie) Search(name string) bool {
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
