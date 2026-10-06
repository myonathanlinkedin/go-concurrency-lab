package main

import (
	"sort"
)

// TrieNode represents a node in the trie.
type TrieNode struct {
	children map[byte]*TrieNode
	isWord   bool
	freq     int
	word     string
}

// Trie represents the prefix tree.
type Trie struct {
	root *TrieNode
}

// NewTrie creates a new empty trie.
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{children: make(map[byte]*TrieNode)},
	}
}

// Insert adds a word with its frequency to the trie.
func (t *Trie) Insert(word string, freq int) {
	node := t.root
	for i := 0; i < len(word); i++ {
		c := word[i]
		if node.children[c] == nil {
			node.children[c] = &TrieNode{children: make(map[byte]*TrieNode)}
		}
		node = node.children[c]
	}
	node.isWord = true
	node.freq = freq
	node.word = word
}

// Search checks if a word exists in the trie.
func (t *Trie) Search(word string) bool {
	node := t.root
	for i := 0; i < len(word); i++ {
		c := word[i]
		if node.children[c] == nil {
			return false
		}
		node = node.children[c]
	}
	return node.isWord
}

// AutoComplete returns up to limit words that start with the given prefix,
// sorted by descending frequency and then lexicographically.
func (t *Trie) AutoComplete(prefix string, limit int) []string {
	node := t.root
	for i := 0; i < len(prefix); i++ {
		c := prefix[i]
		if node.children[c] == nil {
			return []string{}
		}
		node = node.children[c]
	}
	type entry struct {
		word string
		freq int
	}
	var collect func(*TrieNode)
	results := []entry{}
	collect = func(n *TrieNode) {
		if n.isWord {
			results = append(results, entry{word: n.word, freq: n.freq})
		}
		for _, child := range n.children {
			collect(child)
		}
	}
	collect(node)
	sort.Slice(results, func(i, j int) bool {
		if results[i].freq != results[j].freq {
			return results[i].freq > results[j].freq
		}
		return results[i].word < results[j].word
	})
	if limit > len(results) {
		limit = len(results)
	}
	out := make([]string, limit)
	for i := 0; i < limit; i++ {
		out[i] = results[i].word
	}
	return out
}
