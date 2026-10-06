package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Basic broadcast test
	net := NewNetwork()
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		p := NewProcess(i, nil)
		net.AddProcess(p)
		p.Start(&wg)
	}
	for i := 0; i < 5; i++ {
		net.Broadcast(Message{From: 0, Body: "inc"})
	}
	time.Sleep(100 * time.Millisecond)
	for _, p := range net.processes {
		if p.Counter() != 5 {
			panic(fmt.Sprintf("process %d counter = %d, want 5", p.ID, p.Counter()))
		}
	}
	// Direct send test
	net.Send(2, Message{From: 0, To: 2, Body: "inc"})
	time.Sleep(50 * time.Millisecond)
	if net.processes[2].Counter() != 6 {
		panic(fmt.Sprintf("process 2 counter after direct send = %d, want 6", net.processes[2].Counter()))
	}
	net.Close()
	wg.Wait()
	fmt.Println("All functional tests passed")

	// Simple benchmark
	benchmarkBroadcast()
}

// benchmarkBroadcast measures broadcast throughput for a larger network.
func benchmarkBroadcast() {
	net := NewNetwork()
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		p := NewProcess(i, nil)
		net.AddProcess(p)
		p.Start(&wg)
	}
	const msgs = 10000
	start := time.Now()
	for i := 0; i < msgs; i++ {
		net.Broadcast(Message{From: 0, Body: "inc"})
	}
	// Allow processing to catch up
	time.Sleep(200 * time.Millisecond)
	elapsed := time.Since(start)
	fmt.Printf("Broadcast %d messages to %d processes in %v\n", msgs, len(net.processes), elapsed)
	net.Close()
	wg.Wait()
}
