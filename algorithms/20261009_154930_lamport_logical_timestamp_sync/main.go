package main

import (
	"fmt"
	"sync"
)

func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: " + msg)
	}
}

func testTick() {
	clk := NewLamportClock()
	t1 := clk.Tick()
	t2 := clk.Tick()
	assert(t2 == t1+1, "Tick should increment by one")
	assert(clk.Value() == t2, "Value should reflect last tick")
}

func testUpdate() {
	clk := NewLamportClock()
	clk.Tick() // counter = 1
	ts := clk.Update(5)
	assert(ts == 6, "Update with larger remote should set counter to remote+1")
	assert(clk.Value() == 6, "Value after update must match")
	ts2 := clk.Update(3)
	assert(ts2 == 7, "Update with smaller remote still increments")
}

func testEngineOrdering() {
	eng := NewEngine()
	e1 := eng.LocalEvent("a")
	e2 := eng.RemoteEvent("b", 10)
	e3 := eng.LocalEvent("c")
	log := eng.Log()
	assert(log[0].Timestamp == e1.Timestamp, "first event order")
	assert(log[1].Timestamp == e2.Timestamp, "second event order")
	assert(log[2].Timestamp == e3.Timestamp, "third event order")
	assert(Compare(e1.Timestamp, e2.Timestamp) == -1, "e1 < e2")
	assert(Compare(e2.Timestamp, e3.Timestamp) == -1, "e2 < e3")
}

func testConcurrency() {
	eng := NewEngine()
	var wg sync.WaitGroup
	// Spawn local events.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			eng.LocalEvent(fmt.Sprintf("local-%d", i))
		}(i)
	}
	// Spawn remote events with varying timestamps.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			remoteTS := uint64(i * 2)
			eng.RemoteEvent(fmt.Sprintf("remote-%d", i), remoteTS)
		}(i)
	}
	wg.Wait()
	log := eng.Log()
	assert(len(log) == 200, "expected 200 events after concurrency test")
	for i := 1; i < len(log); i++ {
		assert(log[i].Timestamp > log[i-1].Timestamp, "timestamps must be strictly monotonic")
	}
}

func main() {
	fmt.Println("Running Lamport Clock unit tests...")
	testTick()
	testUpdate()
	testEngineOrdering()
	testConcurrency()
	fmt.Println("All tests passed.")
}
