package main

func (c *PNCounter) Increment(delta uint64) {
	if _, ok := c.P[c.id]; !ok {
		c.P[c.id] = 0
	}
	c.P[c.id] += delta
}

func (c *PNCounter) Decrement(delta uint64) {
	if _, ok := c.N[c.id]; !ok {
		c.N[c.id] = 0
	}
	c.N[c.id] += delta
}

func (c *PNCounter) Value() int64 {
	var sumP uint64
	var sumN uint64
	for _, v := range c.P {
		sumP += v
	}
	for _, v := range c.N {
		sumN += v
	}
	return int64(sumP) - int64(sumN)
}

func (c *PNCounter) Merge(other *PNCounter) {
	for id, v := range other.P {
		if cur, ok := c.P[id]; !ok || v > cur {
			c.P[id] = v
		}
	}
	for id, v := range other.N {
		if cur, ok := c.N[id]; !ok || v > cur {
			c.N[id] = v
		}
	}
}
