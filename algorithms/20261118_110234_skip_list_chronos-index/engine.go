package main

import (
	"fmt"
)

// randomLevel generates a level for a new node based on the probability.
// It never exceeds the skip list's maxLevel.
func (sl *SkipList) randomLevel() int {
	lvl := 1
	for sl.randSrc.Float64() < sl.prob && lvl < sl.maxLevel {
		lvl++
	}
	return lvl
}

// Insert adds a key-value pair to the skip list.
// If the key already exists, its value is updated and false is returned.
// If the key is new, it is inserted and true is returned.
func (sl *SkipList) Insert(key int64, value interface{}) bool {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	update := make([]*Node, sl.maxLevel)
	x := sl.header

	// Find the place where the new node will be inserted.
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}

	// Check if key already exists.
	if next := x.forward[0]; next != nil && next.key == key {
		next.value = value
		return false
	}

	// Determine level for new node.
	newLevel := sl.randomLevel()
	if newLevel > sl.level {
		for i := sl.level; i < newLevel; i++ {
			update[i] = sl.header
		}
		sl.level = newLevel
	}

	newNode := &Node{
		key:     key,
		value:   value,
		forward: make([]*Node, newLevel),
	}

	// Splice new node into the list.
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
	return true
}

// Search looks for a key and returns its associated value and true if found.
func (sl *SkipList) Search(key int64) (interface{}, bool) {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()

	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
	}
	x = x.forward[0]
	if x != nil && x.key == key {
		return x.value, true
	}
	return nil, false
}

// Delete removes a key from the skip list.
// Returns true if the key was present and removed.
func (sl *SkipList) Delete(key int64) bool {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	update := make([]*Node, sl.maxLevel)
	x := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}

	target := x.forward[0]
	if target == nil || target.key != key {
		return false
	}

	for i := 0; i < sl.level; i++ {
		if update[i].forward[i] != target {
			break
		}
		update[i].forward[i] = target.forward[i]
	}

	// Reduce level if highest levels are empty.
	for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
		sl.level--
	}
	return true
}

// RangeQuery returns all values whose keys lie in [start, end] inclusive, ordered by key.
func (sl *SkipList) RangeQuery(start, end int64) []interface{} {
	if start > end {
		return []interface{}{}
	}
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()

	result := []interface{}{}
	x := sl.header

	// Locate first node >= start.
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < start {
			x = x.forward[i]
		}
	}
	x = x.forward[0]

	// Collect while within range.
	for x != nil && x.key <= end {
		if x.key >= start {
			result = append(result, x.value)
		}
		x = x.forward[0]
	}
	return result
}

// DebugString returns a string representation of the skip list (for debugging).
func (sl *SkipList) DebugString() string {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()
	s := ""
	for i := sl.level - 1; i >= 0; i-- {
		s += fmt.Sprintf("Level %d: ", i+1)
		x := sl.header.forward[i]
		for x != nil {
			s += fmt.Sprintf("%d ", x.key)
			x = x.forward[i]
		}
		s += "\n"
	}
	return s
}
