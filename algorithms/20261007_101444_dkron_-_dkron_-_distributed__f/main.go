package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	// Test 1: basic execution count
	scheduler := NewScheduler()
	var count int32
	jobFunc := func() error {
		atomic.AddInt32(&count, 1)
		return nil
	}
	if err := scheduler.AddJob("test1", 100*time.Millisecond, jobFunc, 0); err != nil {
		panic(err)
	}
	time.Sleep(350 * time.Millisecond)
	scheduler.Stop()
	if atomic.LoadInt32(&count) < 3 {
		panic(fmt.Sprintf("expected at least 3 executions, got %d", count))
	}

	// Test 2: retry on failure
	scheduler = NewScheduler()
	var attempts int32
	jobFunc2 := func() error {
		a := atomic.AddInt32(&attempts, 1)
		if a == 1 {
			return fmt.Errorf("intentional failure")
		}
		return nil
	}
	if err := scheduler.AddJob("retryJob", 100*time.Millisecond, jobFunc2, 1); err != nil {
		panic(err)
	}
	time.Sleep(250 * time.Millisecond)
	scheduler.Stop()
	if atomic.LoadInt32(&attempts) < 2 {
		panic(fmt.Sprintf("expected at least 2 attempts, got %d", attempts))
	}

	// Test 3: removal stops execution
	scheduler = NewScheduler()
	var cnt3 int32
	jobFunc3 := func() error {
		atomic.AddInt32(&cnt3, 1)
		return nil
	}
	if err := scheduler.AddJob("toRemove", 100*time.Millisecond, jobFunc3, 0); err != nil {
		panic(err)
	}
	time.Sleep(150 * time.Millisecond)
	if err := scheduler.RemoveJob("toRemove"); err != nil {
		panic(err)
	}
	prev := atomic.LoadInt32(&cnt3)
	time.Sleep(200 * time.Millisecond)
	scheduler.Stop()
	if atomic.LoadInt32(&cnt3) != prev {
		panic(fmt.Sprintf("job not removed, count changed from %d to %d", prev, cnt3))
	}

	fmt.Println("All tests passed")
}
