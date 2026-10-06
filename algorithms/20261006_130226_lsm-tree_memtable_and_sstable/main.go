package main

import (
	"fmt"
	"time"
)

func main() {
	tree := NewLSMTree(3)

	// Basic Put/Get
	tree.Put("a", "1")
	tree.Put("b", "2")
	tree.Put("c", "3")
	if v, ok := tree.Get("a"); !ok || v != "1" {
		panic("assertion failed: get a")
	}
	if v, ok := tree.Get("b"); !ok || v != "2" {
		panic("assertion failed: get b")
	}
	if v, ok := tree.Get("c"); !ok || v != "3" {
		panic("assertion failed: get c")
	}

	// Trigger flush by adding another key
	tree.Put("d", "4")
	if len(tree.sstables) != 1 {
		panic("assertion failed: flush not triggered")
	}
	if v, ok := tree.Get("a"); !ok || v != "1" {
		panic("assertion failed: get a after flush")
	}
	if v, ok := tree.Get("d"); !ok || v != "4" {
		panic("assertion failed: get d after flush")
	}

	// Delete key
	tree.Delete("b")
	if _, ok := tree.Get("b"); ok {
		panic("assertion failed: delete b")
	}

	// Flush after delete
	tree.Put("e", "5")
	tree.Put("f", "6")
	if len(tree.sstables) != 2 {
		panic("assertion failed: second flush")
	}
	if _, ok := tree.Get("b"); ok {
		panic("assertion failed: b should still be deleted")
	}
	if v, ok := tree.Get("e"); !ok || v != "5" {
		panic("assertion failed: get e")
	}

	// Compact
	tree.Compact()
	if len(tree.sstables) != 1 {
		panic("assertion failed: compact did not reduce to one sstable")
	}
	if _, ok := tree.Get("b"); ok {
		panic("assertion failed: b should be removed after compaction")
	}
	if v, ok := tree.Get("f"); !ok || v != "6" {
		panic("assertion failed: get f after compaction")
	}

	// Benchmark simple bulk insert
	start := time.Now()
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key%d", i)
		tree.Put(key, "value")
	}
	elapsed := time.Since(start)
	fmt.Printf("Inserted 10000 keys in %s\n", elapsed)

	fmt.Println("All assertions passed.")
}
