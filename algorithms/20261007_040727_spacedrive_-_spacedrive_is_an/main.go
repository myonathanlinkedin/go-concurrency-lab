package main

import (
	"fmt"
	"sort"
)

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := append([]string(nil), a...)
	bCopy := append([]string(nil), b...)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}

func main() {
	root := NewNode("", true)

	paths := []string{
		"docs/readme.md",
		"docs/tutorial/intro.md",
		"src/main.go",
		"src/utils/helpers.go",
		"src/utils/helpers_test.go",
		"assets/images/logo.png",
	}

	for _, p := range paths {
		if err := root.AddPath(p); err != nil {
			panic(fmt.Sprintf("AddPath error: %v", err))
		}
	}

	// Test listing of docs directory
	docs, err := root.Find("docs")
	if err != nil || docs == nil {
		panic("docs directory not found")
	}
	list := docs.List()
	expectedDocs := []string{"readme.md", "tutorial"}
	if !equalStringSlices(list, expectedDocs) {
		panic(fmt.Sprintf("docs list mismatch: got %v, want %v", list, expectedDocs))
	}

	// Test finding a file
	file, err := root.Find("src/utils/helpers.go")
	if err != nil || file == nil || file.isDir {
		panic("helpers.go not found or is dir")
	}

	// Test non-existent path
	_, err = root.Find("nonexistent/file.txt")
	if err == nil {
		panic("expected error for nonexistent path")
	}

	// Simple benchmark loop
	var sum int
	for i := 0; i < 100000; i++ {
		n, _ := root.Find("src/main.go")
		if n != nil {
			sum++
		}
	}
	fmt.Println("Benchmark sum:", sum)

	fmt.Println("All assertions passed.")
}
