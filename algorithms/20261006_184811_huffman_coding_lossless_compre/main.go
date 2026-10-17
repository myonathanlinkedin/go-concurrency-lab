package main

import (
	"fmt"
	"time"
)

func main() {
	// Test 1: Basic example
	data := []byte("this is an example for huffman encoding")
	freq := make(map[byte]int)
	for _, b := range data {
		freq[b]++
	}
	root := BuildTree(freq)
	encoded, _ := Encode(data, root)
	decoded := Decode(encoded, root)
	if string(decoded) != string(data) {
		panic("Test 1 failed: decoded data does not match original")
	}
	fmt.Println("Test 1 passed")

	// Test 2: Single character
	data2 := []byte("aaaaaa")
	freq2 := make(map[byte]int)
	for _, b := range data2 {
		freq2[b]++
	}
	root2 := BuildTree(freq2)
	encoded2, _ := Encode(data2, root2)
	decoded2 := Decode(encoded2, root2)
	if string(decoded2) != string(data2) {
		panic("Test 2 failed: single character decoding mismatch")
	}
	fmt.Println("Test 2 passed")

	// Test 3: Empty data
	data3 := []byte{}
	freq3 := make(map[byte]int)
	root3 := BuildTree(freq3)
	encoded3, _ := Encode(data3, root3)
	decoded3 := Decode(encoded3, root3)
	if len(decoded3) != 0 {
		panic("Test 3 failed: empty data decoding mismatch")
	}
	fmt.Println("Test 3 passed")

	// Benchmark: encode/decode performance
	benchmarkData := []byte("benchmarking huffman coding with a larger dataset to measure performance")
	start := time.Now()
	for i := 0; i < 1000; i++ {
		freqB := make(map[byte]int)
		for _, b := range benchmarkData {
			freqB[b]++
		}
		rootB := BuildTree(freqB)
		encB, _ := Encode(benchmarkData, rootB)
		decB := Decode(encB, rootB)
		if string(decB) != string(benchmarkData) {
			panic("Benchmark failed: data mismatch")
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("Benchmark completed in %s\n", elapsed)
}
