package main

import (
	"fmt"
)

func assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func testBasic() {
	s := "ababa"
	sam := NewSAM(s)

	// Substring existence tests
	assert(sam.Contains("a"), "contains a")
	assert(sam.Contains("b"), "contains b")
	assert(sam.Contains("ab"), "contains ab")
	assert(sam.Contains("ba"), "contains ba")
	assert(sam.Contains("aba"), "contains aba")
	assert(sam.Contains("bab"), "contains bab")
	assert(sam.Contains("ababa"), "contains ababa")
	assert(!sam.Contains("c"), "does not contain c")
	assert(!sam.Contains("aa"), "does not contain aa")
	assert(!sam.Contains("abb"), "does not contain abb")

	// Distinct substrings count: for "ababa" it's 9
	assert(sam.CountDistinct() == 9, fmt.Sprintf("distinct count %d != 9", sam.CountDistinct()))
}

func testEdgeCases() {
	// Empty string
	empty := NewSAM("")
	assert(empty.CountDistinct() == 0, "empty distinct count")
	assert(empty.Contains(""), "empty contains empty")
	assert(!empty.Contains("a"), "empty does not contain a")

	// Single character
	single := NewSAM("x")
	assert(single.CountDistinct() == 1, "single distinct count")
	assert(single.Contains("x"), "single contains x")
	assert(!single.Contains("xx"), "single does not contain xx")
	assert(single.Contains(""), "single contains empty")
}

func testLonger() {
	// Repeated pattern
	s := "abcabcabc"
	sam := NewSAM(s)
	// Known distinct substrings count for this string is 15
	assert(sam.CountDistinct() == 15, fmt.Sprintf("longer distinct %d != 15", sam.CountDistinct()))
	assert(sam.Contains("abcabc"), "contains abcabc")
	assert(sam.Contains("bcab"), "contains bcab")
	assert(!sam.Contains("abcd"), "does not contain abcd")
}

func main() {
	testBasic()
	testEdgeCases()
	testLonger()
	fmt.Println("All tests passed")
}
