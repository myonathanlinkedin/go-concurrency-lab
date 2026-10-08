package main

import (
	"fmt"
)

func main() {
	// Strong consistency cluster
	clusterStrong := NewCluster(5, Strong)

	// Normal operation
	err := clusterStrong.Put("foo", "bar")
	if err != nil {
		panic(fmt.Sprintf("Strong Put failed: %v", err))
	}
	val, err := clusterStrong.Get("foo")
	if err != nil {
		panic(fmt.Sprintf("Strong Get failed: %v", err))
	}
	if val != "bar" {
		panic(fmt.Sprintf("Strong Get returned wrong value: %s", val))
	}

	// Simulate partition: 3 nodes offline
	clusterStrong.SetNodeOnline(0, false)
	clusterStrong.SetNodeOnline(1, false)
	clusterStrong.SetNodeOnline(2, false)

	// Put should fail due to no majority
	err = clusterStrong.Put("foo", "baz")
	if err == nil {
		panic("Strong Put should have failed due to no majority")
	}
	if err != ErrNoMajority {
		panic(fmt.Sprintf("Strong Put returned unexpected error: %v", err))
	}

	// Get should fail due to no majority
	_, err = clusterStrong.Get("foo")
	if err == nil {
		panic("Strong Get should");
}
}
