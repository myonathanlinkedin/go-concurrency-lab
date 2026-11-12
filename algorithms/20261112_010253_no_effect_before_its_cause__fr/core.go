package main

type LamportClock struct {
	time uint64
}

func NewLamportClock() *LamportClock {
	return &LamportClock{time: 0}
}

func (lc *LamportClock) Increment() {
	lc.time++
}

func (lc *LamportClock) Update(received uint64) {
	if received >= lc.time {
		lc.time = received + 1
	} else {
		lc.time++
	}
}

func (lc *LamportClock) Timestamp() uint64 {
	return lc.time
}

type Process struct {
	id    int
	clock *LamportClock
}

func NewProcess(id int) *Process {
	return &Process{id: id, clock: NewLamportClock()}
}

func (p *Process) Send(to *Process) {
	p.clock.Increment()
	to.Receive(p.clock.Timestamp())
}

func (p *Process) Receive(timestamp uint64) {
	p.clock.Update(timestamp)
}
