package main

import "fmt"

func main() {
	tree := &RBTree{}

	// Insert a known sequence and validate after each insertion.
	keys := []int{5, 3, 7, 2, 4, 6, 8, 1, 9}
	for _, k := range keys {
		tree.Insert(k, fmt.Sprintf("val%d", k))
		if !tree.Validate() {
			panic(fmt.Sprintf("validation failed after inserting key %d", k))
		}
	}

	// Verify retrieval of all inserted keys.
	for _, k := range keys {
		if v, ok := tree.Get(k); !ok || v != fmt.Sprintf("val%d", k) {
			panic(fmt.Sprintf("lookup failed for key %d, got (%s,%v)", k, v, ok))
		}
	}

	// Delete keys in the same order, validating after each deletion.
	for _, k := range keys {
		tree.Delete(k)
		if !tree.Validate() {
			panic(fmt.Sprintf("validation failed after deleting key %d", k))
		}
		if _, ok := tree.Get(k); ok {
			panic(fmt.Sprintf("key %d still present after deletion", k))
		}
	}

	// Tree should be empty now.
	if tree.root != nil {
		panic("tree not empty after all deletions")
	}

	fmt.Println("All red‑black tree tests passed")
}
