package main

type state struct {
	next   map[byte]int
	link   int
	length int
}

// SAM implements a suffix automaton for a given string.
type SAM struct {
	st   []state
	last int
}

// NewSAM builds a suffix automaton for the provided string.
func NewSAM(s string) *SAM {
	sam := &SAM{
		st:   make([]state, 1, 2*len(s)+1),
		last: 0,
	}
	sam.st[0] = state{next: make(map[byte]int), link: -1, length: 0}
	for i := 0; i < len(s); i++ {
		sam.extend(s[i])
	}
	return sam
}

// extend adds a character to the automaton.
func (sam *SAM) extend(c byte) {
	cur := len(sam.st)
	sam.st = append(sam.st, state{
		next:   make(map[byte]int),
		length: sam.st[sam.last].length + 1,
	})
	p := sam.last
	for p != -1 {
		if _, ok := sam.st[p].next[c]; !ok {
			sam.st[p].next[c] = cur
			p = sam.st[p].link
		} else {
			break
		}
	}
	if p == -1 {
		sam.st[cur].link = 0
	} else {
		q := sam.st[p].next[c]
		if sam.st[p].length+1 == sam.st[q].length {
			sam.st[cur].link = q
		} else {
			clone := len(sam.st)
			sam.st = append(sam.st, state{
				next:   make(map[byte]int),
				link:   sam.st[q].link,
				length: sam.st[p].length + 1,
			})
			for k, v := range sam.st[q].next {
				sam.st[clone].next[k] = v
			}
			for p != -1 {
				if nxt, ok := sam.st[p].next[c]; ok && nxt == q {
					sam.st[p].next[c] = clone
					p = sam.st[p].link
				} else {
					break
				}
			}
			sam.st[q].link = clone
			sam.st[cur].link = clone
		}
	}
	sam.last = cur
}

// Contains checks whether substr exists as a substring of the original string.
func (sam *SAM) Contains(substr string) bool {
	v := 0
	for i := 0; i < len(substr); i++ {
		c := substr[i]
		nxt, ok := sam.st[v].next[c]
		if !ok {
			return false
		}
		v = nxt
	}
	return true
}

// CountDistinct returns the number of distinct substrings of the original string.
func (sam *SAM) CountDistinct() int64 {
	var total int64 = 0
	for i := 1; i < len(sam.st); i++ {
		total += int64(sam.st[i].length - sam.st[sam.st[i].link].length)
	}
	return total
}
