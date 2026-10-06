package main

import "fmt"

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func testIncrementAndCompare() {
	nodes := []int{1, 2, 3}
	vcA := NewVectorClock(nodes)
	vcB := NewVectorClock(nodes)

	vcA.Increment(1)
	vcA.Increment(1) // A: {1:2}
	vcB.Increment(2) // B: {2:1}

	assert(vcA.Compare(vcB) == Concurrent, "A and B should be concurrent")
	assert(vcA.Compare(vcA) == Equal, "A should equal itself")
}

func testMergeAndHappensBefore() {
	nodes := []int{1, 2}
	vc1 := NewVectorClock(nodes)
	vc2 := NewVectorClock(nodes)

	vc1.Increment(1) // {1:1}
	vc2.Increment(1) // {1:1}
	vc2.Increment(2) // {1:1,2:1}

	// vc1 happens before vc2 after merge
	vc1.Merge(vc2)
	assert(vc1.Compare(vc2) == Equal, "After merge, clocks should be equal")
	assert(vc1.Compare(vc1) == Equal, "Self comparison after merge")
}

func testConcurrentDetection() {
	nodes := []int{1, 2}
	c1 := NewVectorClock(nodes)
	c2 := NewVectorClock(nodes)

	c1.Increment(1) // {1:1}
	c2.Increment(2) // {2:1}

	assert(c1.Compare(c2) == Concurrent, "c1 and c2 should be concurrent")
	c1.Merge(c2) // c1 becomes {1:1,2:1}
	assert(c1.Compare(c2) == HappensAfter, "c1 should happen after c2 after merge")
}

func testCloneIndependence() {
	nodes := []int{1}
	orig := NewVectorClock(nodes)
	orig.Increment(1) // {1:1}
	clone := orig.Clone()
	clone.Increment(1) // clone {1:2}
	assert(orig.Compare(clone) == HappensBefore, "orig should happen before clone")
	assert(clone.Compare(orig) == HappensAfter, "clone should happen after orig")
}

func main() {
	testIncrementAndCompare()
	testMergeAndHappensBefore()
	testConcurrentDetection()
	testCloneIndependence()
	fmt.Println("All vector clock tests passed.")
}
