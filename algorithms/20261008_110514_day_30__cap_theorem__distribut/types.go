package main

import (
	"errors"
)

type ConsistencyMode int

const (
	Strong ConsistencyMode = iota
	Eventual
)

var (
	ErrNodeOffline = errors.New("node is offline")
	ErrNoMajority  = errors.New("no majority available")
	ErrKeyNotFound = errors.New("key not found")
)

type Node struct {
	id     int
	store  map[string]string
	online bool
}

func NewNode(id int) *Node {
	return &Node{
		id:     id,
		store:  make(map[string]string),
		online: true,
	}
}

type Cluster struct {
	nodes []*Node
	mode  ConsistencyMode
}

func NewCluster(nodeCount int, mode ConsistencyMode) *Cluster {
	nodes := make([]*Node, nodeCount)
	for i := 0; i < nodeCount; i++ {
		nodes[i] = NewNode(i)
	}
	return &Cluster{
		nodes: nodes,
		mode:  mode,
	}
}

func (c *Cluster) SetNodeOnline(id int, online bool) {
	if id >= 0 && id < len(c.nodes) {
		c.nodes[id].online = online
	}
}

func (c *Cluster) OnlineNodes() []*Node {
	online := []*Node{}
	for _, n := range c.nodes {
		if n.online {
			online = append(online, n)
		}
	}
	return online
}

func (c *Cluster) MajorityCount() int {
	return len(c.nodes)/2 + 1
}
