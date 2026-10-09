package main

import (
	"fmt"
	"time"
)

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func testIncrement() {
	vc := NewVectorClock()
	vc.Inc(1)
	vc.Inc(2)
	vc.Inc(1)
	assert(vc[1] == 2, "node 1 count")
	assert(vc[2] == 1, "node 2 count")
}

func testMerge() {
	a := NewVectorClock()
	b := NewVectorClock()
	a.Inc(1) // a: {1:1}
	a.Inc(2) // a: {1:1,2:1}
	b.Inc(2) // b: {2:1}
	b.Inc(3) // b: {2:1,3:1}
	a.Merge(b)
	assert(a[1] == 1, "merge keep 1")
	assert(a[2] == 1, "merge max 2")
	assert(a[3] == 1, "merge add 3")
}

func testCompare() {
	x := NewVectorClock()
	y := NewVectorClock()
	assert(x.Compare(y) == Equal, "empty equal")
	x.Inc(1) // x:{1:1}
	assert(x.Compare(y) == HappensAfter, "x after y")
	assert(y.Compare(x) == HappensBefore, "y before x")
	y.Inc(1) // y:{1:1}
	assert(x.Compare(y) == Equal, "now equal")
	y.Inc(2) // y:{1:1,2:1}
	assert(x.Compare(y) == HappensBefore, "x before y")
	x.Inc(2) // x:{1:1,2:1}
	assert(x.Compare(y) == Equal, "equal again")
	// concurrent case
	a := NewVectorClock()
	b := NewVectorClock()
	a.Inc(1) // a:{1:1}
	b.Inc(2) // b:{2:1}
	assert(a.Compare(b) == Concurrent, "concurrent")
	assert(b.Compare(a) == Concurrent, "concurrent reverse")
}

func benchmarkMerge(iter int) time.Duration {
	a := NewVectorClock()
	b := NewVectorClock()
	for i := 0; i < iter; i++ {
		a.Inc(i % 5)
		b.Inc((i + 2) % 5)
	}
	start := time.Now()
	a.Merge(b)
	return time.Since(start)
}

func main() {
	fmt.Println("Running VectorClock tests...")
	testIncrement()
	testMerge()
	testCompare()
	fmt.Println("All assertions passed.")
	dur := benchmarkMerge(100000)
	fmt.Printf("Benchmark merge 100k ops: %v\n", dur)
}
