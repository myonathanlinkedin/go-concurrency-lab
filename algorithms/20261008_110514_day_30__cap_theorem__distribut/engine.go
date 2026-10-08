package main

import (
	"errors"
)

func (c *Cluster) Put(key, value string) error {
	online := c.OnlineNodes()
	if len(online) == 0 {
		return ErrNodeOffline
	}
	switch c.mode {
	case Strong:
		if len(online) < c.MajorityCount() {
			return ErrNoMajority
		}
		for _, n := range online {
			n.store[key] = value
		}
		return nil
	case Eventual:
		for _, n := range online {
			n.store[key] = value
		}
		return nil
	default:
		return errors.New("unknown consistency mode")
	}
}

func (c *Cluster) Get(key string) (string, error) {
	online := c.OnlineNodes()
	if len(online) == 0 {
		return "", ErrNodeOffline
	}
	switch c.mode {
	case Strong:
		if len(online) < c.MajorityCount() {
			return "", ErrNoMajority
		}
		val, ok := online[0].store[key]
		if !ok {
			return "", ErrKeyNotFound
		}
		return val, nil
	case Eventual:
		val, ok := online[0].store[key]
		if !ok {
			return "", ErrKeyNotFound
		}
		return val, nil
	default:
		return "", errors.New("unknown consistency mode")
	}
}
