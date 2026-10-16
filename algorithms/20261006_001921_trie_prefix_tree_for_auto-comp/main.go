package main

import (
	"fmt"
	"os"
	"time"
)

type testResult struct {
	name   string
	passed bool
	err    string
}

func runTest(name string, fn func() error) testResult {
	err := fn()
	if err != nil {
		return testResult{name, false, err.Error()}
	}
	return testResult{name, true, ""}
}

func testInsertAndSearch() error {
	trie := NewTrie()
	trie.Insert("apple", 5)
	trie.Insert("app", 3)
	trie.Insert("banana", 2)
	if !trie.Search("apple") {
		return fmt.Errorf("expected apple to exist")
	}
	if !trie.Search("app") {
		return fmt.Errorf("expected app to exist")
	}
	if trie.Search("apricot") {
		return fmt.Errorf("apricot should not exist")
	}
	return nil
}

func testAutoComplete() error {
	trie := NewTrie()
	trie.Insert("apple", 5)
	trie.Insert("app", 3)
	trie.Insert("application", 4)
	trie.Insert("apex", 2)
	trie.Insert("banana", 1)
	suggestions := trie.AutoComplete("app", 3)
	expected := []string{"apple", "application", "app"}
	if len(suggestions) != len(expected) {
		return fmt.Errorf("expected %d suggestions, got %d", len(expected), len(suggestions))
	}
	for i, w := range expected {
		if suggestions[i] != w {
			return fmt.Errorf("expected %s at position %d, got %s", w, i, suggestions[i])
		}
	}
	return nil
}

func testEmptyTrie() error {
	trie := NewTrie()
	suggestions := trie.AutoComplete("any", 5)
	if len(suggestions) != 0 {
		return fmt.Errorf("expected 0 suggestions, got %d", len(suggestions))
	}
	return nil
}

func benchmarkAutoComplete() error {
	trie := NewTrie()
	words := []struct {
		word string
		freq int
	}{
		{"alpha", 10}, {"beta", 8}, {"gamma", 6}, {"delta", 4},
		{"epsilon", 2}, {"zeta", 1}, {"eta", 5}, {"theta", 3},
	}
	for _, w := range words {
		trie.Insert(w.word, w.freq)
	}
	start := time.Now()
	for i := 0; i < 1000; i++ {
		trie.AutoComplete("a", 5)
	}
	elapsed := time.Since(start)
	if elapsed > time.Second {
		return fmt.Errorf("benchmark took too long: %s", elapsed)
	}
	return nil
}

func main() {
	tests := []struct {
		name string
		fn   func() error
	}{
		{"TestInsertAndSearch", testInsertAndSearch},
		{"TestAutoComplete", testAutoComplete},
		{"TestEmptyTrie", testEmptyTrie},
		{"BenchmarkAutoComplete", benchmarkAutoComplete},
	}
	var passed, failed int
	for _, t := range tests {
		res := runTest(t.name, t.fn)
		if res.passed {
			fmt.Printf("[PASS] %s\n", res.name)
			passed++
		} else {
			fmt.Printf("[FAIL] %s: %s\n", res.name, res.err)
			failed++
		}
	}
	fmt.Printf("\nTotal: %d passed, %d failed\n", passed, failed)
	if failed > 0 {
		os.Exit(1)
	}
}
