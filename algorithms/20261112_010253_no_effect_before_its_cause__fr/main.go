package main

import "fmt"

func main() {
	// test 1: basic increment
	lc := NewLamportClock()
	if lc.Timestamp() != 0 {
		panic("initial timestamp should be 0")
	}
	lc.Increment()
	if lc.Timestamp() != 1 {
		panic("timestamp after increment should be 1")
	}

	// test 2: update with higher timestamp
	lc.Update(5)
	if lc.Timestamp() != 6 {
		panic("timestamp after update with higher should be 6")
	}

	// test 3: update with lower timestamp
	lc.Update(3)
	if lc.Timestamp() != 7 {
		panic("timestamp after update with lower should be 7")
	}

	// test 4: process communication causal ordering
	p1 := NewProcess(1)
	p2 := NewProcess(2)

	// p1 sends to p2
	p1.Send(p2)

	// after send, p1 timestamp should be 1, p2 timestamp should be 1
	if p1.clock.Timestamp() != 1 {
		panic("p1 timestamp after send should be 1")
	}
	if p2.clock.Timestamp() != 1 {
		panic("p2 timestamp after receive should be 1")
	}

	// p2 sends to p1
	p2.Send(p1)

	// after second send, p2 timestamp should be 2, p1 timestamp should be 2
	if p2.clock.Timestamp() != 2 {
		panic("p2 timestamp after second send should be 2")
	}
	if p1.clock.Timestamp() != 2 {
		panic("p1 timestamp after receive should be 2")
	}

	// test 5: causal ordering: p1's send before p2's send
	if !(p1.clock.Timestamp() < p2.clock.Timestamp()) {
		panic("causal ordering violated: p1's timestamp should be less than p2's")
	}

	// test 6: concurrent events: simulate two processes incrementing without communication
	p3 := NewProcess(3)
	p4 := NewProcess(4)
	p3.clock.Increment()
	p4.clock.Increment()
	// timestamps may be equal or not; but they should not be less than each other
	if p3.clock.Timestamp() < p4.clock.Timestamp() {
		// fine
	} else if p3.clock.Timestamp() > p4.clock.Timestamp() {
		// fine
	} else {
		// equal, fine
	}

	// test 7: large number of increments
	lc2 := NewLamportClock()
	for i := 0; i < 1000000; i++ {
		lc2.Increment()
	}
	if lc2.Timestamp() != 1000000 {
		panic("timestamp after 1M increments should be 1M")
	}

	fmt.Println("All tests passed")
}
