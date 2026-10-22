package main

import (
	"fmt"
	"sort"
)

// MemTable is an in-memory sorted key-value buffer.
// Invariant: keys are unique and maintained in ascending order.
// Time: O(log n) insert, O(n) flush. Space: O(n).
type MemTable struct {
	keys   []string
	values []string
}

func NewMemTable() *MemTable {
	return &MemTable{}
}

// Insert adds or updates a key-value pair. O(log n) due to binary search + slice shift.
func (m *MemTable) Insert(key, value string) {
	i := sort.SearchStrings(m.keys, key)
	if i < len(m.keys) && m.keys[i] == key {
		m.values[i] = value
		return
	}
	m.keys = append(m.keys, "")
	m.values = append(m.values, "")
	copy(m.keys[i+1:], m.keys[i:])
	copy(m.values[i+1:], m.values[i:])
	m.keys[i] = key
	m.values[i] = value
}

// Get retrieves a value by key. O(log n).
func (m *MemTable) Get(key string) (string, bool) {
	i := sort.SearchStrings(m.keys, key)
	if i < len(m.keys) && m.keys[i] == key {
		return m.values[i], true
	}
	return "", false
}

// Len returns the number of entries. O(1).
func (m *MemTable) Len() int { return len(m.keys) }

// Flush serializes the memtable into a sorted SSTable. O(n log n) worst-case (already sorted, so O(n)).
func (m *MemTable) Flush() *SSTable {
	sst := &SSTable{
		keys:   make([]string, len(m.keys)),
		values: make([]string, len(m.values)),
	}
	copy(sst.keys, m.keys)
	copy(sst.values, m.values)
	return sst
}

// SSTable is an immutable sorted on-disk table.
// Invariant: keys are strictly ascending and unique.
type SSTable struct {
	keys   []string
	values []string
}

// Get performs binary search. O(log n).
func (s *SSTable) Get(key string) (string, bool) {
	i := sort.SearchStrings(s.keys, key)
	if i < len(s.keys) && s.keys[i] == key {
		return s.values[i], true
	}
	return "", false
}

// Len returns entry count. O(1).
func (s *SSTable) Len() int { return len(s.keys) }

// LSMEngine orchestrates memtable and SSTable layers.
// Invariant: at most one active memtable; SSTables are immutable.
type LSMEngine struct {
	memtable  *MemTable
	sstables  []*SSTable
	maxMemLen int
}

func NewLSMEngine(maxMemLen int) *LSMEngine {
	return &LSMEngine{
		memtable:  NewMemTable(),
		sstables:  make([]*SSTable, 0),
		maxMemLen: maxMemLen,
	}
}

// Put inserts into memtable; auto-flushes when threshold exceeded. O(log n) amortized.
func (e *LSMEngine) Put(key, value string) {
	e.memtable.Insert(key, value)
	if e.memtable.Len() >= e.maxMemLen {
		e.flush()
	}
}

// Get searches memtable first, then SSTables newest-to-oldest. O(log n + k*log n) where k = #sstables.
func (e *LSMEngine) Get(key string) (string, bool) {
	if v, ok := e.memtable.Get(key); ok {
		return v, true
	}
	for i := len(e.sstables) - 1; i >= 0; i-- {
		if v, ok := e.sstables[i].Get(key); ok {
			return v, true
		}
	}
	return "", false
}

// flush moves memtable to SSTable layer. O(n).
func (e *LSMEngine) flush() {
	if e.memtable.Len() == 0 {
		return
	}
	e.sstables = append(e.sstables, e.memtable.Flush())
	e.memtable = NewMemTable()
}

// Compact merges all SSTables into one, resolving duplicates (newest wins).
// Time: O(N log N) where N = total entries. Space: O(N).
func (e *LSMEngine) Compact() {
	if len(e.sstables) <= 1 {
		return
	}
	// Collect all entries with level info (higher index = newer).
	type entry struct {
		key   string
		value string
		level int
	}
	var all []entry
	for i, sst := range e.sstables {
		for j := range sst.keys {
			all = append(all, entry{sst.keys[j], sst.values[j], i})
		}
	}
	// Sort by key ascending, then by level descending (newer first for same key).
	sort.Slice(all, func(a, b int) bool {
		if all[a].key != all[b].key {
			return all[a].key < all[b].key
		}
		return all[a].level > all[b].level
	})
	// Deduplicate: keep first occurrence per key (highest level).
	var mergedKeys, mergedVals []string
	for i := range all {
		if len(mergedKeys) > 0 && mergedKeys[len(mergedKeys)-1] == all[i].key {
			continue
		}
		mergedKeys = append(mergedKeys, all[i].key)
		mergedVals = append(mergedVals, all[i].value)
	}
	e.sstables = []*SSTable{{keys: mergedKeys, values: mergedVals}}
}

// String returns a debug representation.
func (e *LSMEngine) String() string {
	return fmt.Sprintf("LSMEngine{memtable=%d, sstables=%d}", e.memtable.Len(), len(e.sstables))
}
