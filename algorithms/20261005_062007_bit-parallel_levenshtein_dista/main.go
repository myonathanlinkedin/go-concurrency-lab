package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Levenshtein returns the Levenshtein distance between a and b.
// For strings up to 64 characters it uses Myers' bit‑parallel algorithm,
// otherwise it falls back to the classic dynamic programming approach.
func Levenshtein(a, b string) int {
	if len(a) > len(b) {
		a, b = b, a
	}
	if len(a) == 0 {
		return len(b)
	}
	if len(a) <= 64 {
		return levenshteinBitParallel(a, b)
	}
	return levenshteinDP(a, b)
}

// Myers' bit‑parallel algorithm for strings where len(a) <= 64.
func levenshteinBitParallel(a, b string) int {
	m := len(a)
	// Build Peq: for each byte, a bitmask of positions where it occurs in a.
	var peq [256]uint64
	for i := 0; i < m; i++ {
		peq[a[i]] |= 1 << uint(i)
	}

	var vp uint64 = ^uint64(0) // all 1s
	var vn uint64 = 0
	var curr uint64
	var d = len(b)

	for i := 0; i < len(b); i++ {
		// Eq contains 1s where pattern char matches text char.
		eq := peq[b[i]]

		// Step 1: compute D0.
		x := eq | vn
		tmp := ((x & vp) + vp) ^ vp
		d0 := tmp | x

		// Step 2: compute HP and HN.
		hp := vn | ^(d0 | vp)
		hn := d0 & vp

		// Step 3: shift HP and HN.
		hp = (hp << 1) | 1
		hn = hn << 1

		// Step 4: update VP and VN.
		vp = hn | ^(d0 | hp)
		vn = hp & d0

		// Step 5: adjust distance.
		if (hp>>uint(m-1))&1 == 1 {
			d++
		}
		if (hn>>uint(m-1))&1 == 1 {
			d--
		}
		_ = curr // silence unused warning in older Go versions
	}
	return d
}

// Classic DP implementation for longer strings.
func levenshteinDP(a, b string) int {
	m, n := len(a), len(b)
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}
	prev := make([]int, n+1)
	curr := make([]int, n+1)
	for j := 0; j <= n; j++ {
		prev[j] = j
	}
	for i := 1; i <= m; i++ {
		curr[0] = i
		ai := a[i-1]
		for j := 1; j <= n; j++ {
			cost := 0
			if ai != b[j-1] {
				cost = 1
			}
			del := prev[j] + 1
			ins := curr[j-1] + 1
			sub := prev[j-1] + cost
			curr[j] = min(del, ins, sub)
		}
		prev, curr = curr, prev
	}
	return prev[n]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// Simple assertion helper.
func assertEqual(got, want int, msg string) {
	if got != want {
		panic(fmt.Sprintf("FAIL: %s - got %d, want %d", msg, got, want))
	}
}

// Unit tests executed from main.
func runTests() {
	type test struct {
		a, b string
		dist int
	}
	tests := []test{
		{"", "", 0},
		{"", "abc", 3},
		{"kitten", "sitting", 3},
		{"flaw", "lawn", 2},
		{"intention", "execution", 5},
		{"abcdefghij", "abcdefghij", 0},
		{"abcdefghij", "abcdxfghij", 1},
		{"short", "a very long string that exceeds sixty-four characters for testing fallback", 71},
	}
	for _, tc := range tests {
		got := Levenshtein(tc.a, tc.b)
		assertEqual(got, tc.dist, fmt.Sprintf("%q vs %q", tc.a, tc.b))
	}
	// Randomized consistency check between DP and bit‑parallel for short strings.
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		m := rand.Intn(20) // length up to 20
		n := rand.Intn(20)
		a := randomString(m)
		b := randomString(n)
		d1 := levenshteinBitParallel(a, b)
		d2 := levenshteinDP(a, b)
		if d1 != d2 {
			panic(fmt.Sprintf("Mismatch on %q vs %q: bit=%d dp=%d", a, b, d1, d2))
		}
	}
	fmt.Println("All tests passed.")
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	runTests()
}
