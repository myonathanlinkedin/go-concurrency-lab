package main

type GCounter map[string]uint64

type PNCounter struct {
	id string
	P  GCounter
	N  GCounter
}

func NewPNCounter(id string) *PNCounter {
	return &PNCounter{
		id: id,
		P:  make(GCounter),
		N:  make(GCounter),
	}
}
