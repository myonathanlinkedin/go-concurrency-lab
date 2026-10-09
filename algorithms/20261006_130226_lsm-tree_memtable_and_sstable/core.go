package main

import (
	"sort"
)

type Entry struct {
	key   string
	value string
}

type MemTable struct {
	data map[string]string
}

func NewMemTable() *MemTable {
	return &MemTable{data: make(map[string]string)}
}

func (m *MemTable) Put(key, value string) {
	m.data[key] = value
}

func (m *MemTable) Delete(key string) {
	m.data[key] = ""
}

func (m *MemTable) Size() int {
	return len(m.data)
}

func (m *MemTable) Flush() []*Entry {
	entries := make([]*Entry, 0, len(m.data))
	for k, v := range m.data {
		entries = append(entries, &Entry{key: k, value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})
	m.data = make(map[string]string)
	return entries
}

type SSTable struct {
	entries []*Entry
}

func NewSSTable(entries []*Entry) *SSTable {
	return &SSTable{entries: entries}
}

func (s *SSTable) Find(key string) (string, bool) {
	i := sort.Search(len(s.entries), func(i int) bool {
		return s.entries[i].key >= key
	})
	if i < len(s.entries) && s.entries[i].key == key {
		if s.entries[i].value == "" {
			return "", false
		}
		return s.entries[i].value, true
	}
	return "", false
}

type LSMTree struct {
	mem          *MemTable
	sstables     []*SSTable
	memThreshold int
}

func NewLSMTree(threshold int) *LSMTree {
	return &LSMTree{
		mem:          NewMemTable(),
		sstables:     []*SSTable{},
		memThreshold: threshold,
	}
}

func (t *LSMTree) Put(key, value string) {
	t.mem.Put(key, value)
	if t.mem.Size() >= t.memThreshold {
		t.Flush()
	}
}

func (t *LSMTree) Delete(key string) {
	t.mem.Delete(key)
	if t.mem.Size() >= t.memThreshold {
		t.Flush()
	}
}

func (t *LSMTree) Get(key string) (string, bool) {
	if val, ok := t.mem.data[key]; ok {
		if val == "" {
			return "", false
		}
		return val, true
	}
	for i := len(t.sstables) - 1; i >= 0; i-- {
		if val, ok := t.sstables[i].Find(key); ok {
			return val, true
		}
	}
	return "", false
}

func (t *LSMTree) Flush() {
	entries := t.mem.Flush()
	if len(entries) == 0 {
		return
	}
	t.sstables = append(t.sstables, NewSSTable(entries))
}

func (t *LSMTree) Compact() {
	if len(t.sstables) <= 1 {
		return
	}
	merged := make(map[string]string)
	for i := len(t.sstables) - 1; i >= 0; i-- {
		for _, e := range t.sstables[i].entries {
			if _, exists := merged[e.key]; !exists {
				merged[e.key] = e.value
			}
		}
	}
	entries := make([]*Entry, 0, len(merged))
	for k, v := range merged {
		entries = append(entries, &Entry{key: k, value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].key < entries[j].key
	})
	t.sstables = []*SSTable{NewSSTable(entries)}
}
