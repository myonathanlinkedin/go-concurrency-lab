package main

import (
	"fmt"
	"sort"
)

func assertEqual(a, b interface{}, msg string) {
	if fmt.Sprintf("%v", a) != fmt.Sprintf("%v", b) {
		panic(fmt.Sprintf("assertion failed: %s (got %v, want %v)", msg, a, b))
	}
}

func assertTrue(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("assertion failed: %s", msg))
	}
}

func testMemTableBasics() {
	mt := NewMemTable()
	mt.Insert("b", "2")
	mt.Insert("a", "1")
	mt.Insert("c", "3")
	assertEqual(mt.Len(), 3, "memtable length")
	v, ok := mt.Get("a")
	assertTrue(ok, "key a exists")
	assertEqual(v, "1", "value of a")
	v, ok = mt.Get("b")
	assertTrue(ok, "key b exists")
	assertEqual(v, "2", "value of b")
	_, ok = mt.Get("z")
	assertTrue(!ok, "key z absent")
	// Update existing key
	mt.Insert("b", "20")
	v, _ = mt.Get("b")
	assertEqual(v, "20", "updated value of b")
	assertEqual(mt.Len(), 3, "length unchanged after update")
}

func testMemTableFlush() {
	mt := NewMemTable()
	mt.Insert("x", "10")
	mt.Insert("y", "20")
	mt.Insert("z", "30")
	sst := mt.Flush()
	assertEqual(sst.Len(), 3, "sstable length after flush")
	v, ok := sst.Get("y")
	assertTrue(ok, "sstable has y")
	assertEqual(v, "20", "sstable value y")
	// Keys must be sorted
	assertTrue(sort.StringsAreSorted(sst.keys), "sstable keys sorted")
}

func testLSMEnginePutGet() {
	e := NewLSMEngine(3)
	e.Put("k1", "v1")
	e.Put("k2", "v2")
	e.Put("k3", "v3") // triggers flush (len >= 3)
	assertEqual(len(e.sstables), 1, "one sstable after flush")
	assertEqual(e.memtable.Len(), 0, "memtable empty after flush")
	v, ok := e.Get("k1")
	assertTrue(ok, "k1 found")
	assertEqual(v, "v1", "k1 value")
	v, ok = e.Get("k2")
	assertTrue(ok, "k2 found")
	assertEqual(v, "v2", "k2 value")
	v, ok = e.Get("k3")
	assertTrue(ok, "k3 found")
	assertEqual(v, "v3", "k3 value")
	_, ok = e.Get("k4")
	assertTrue(!ok, "k4 not found")
}

func testLSMEngineOverwriteAcrossFlush() {
	e := NewLSMEngine(2)
	e.Put("a", "old")
	e.Put("b", "old") // flush
	e.Put("a", "new") // in memtable
	v, ok := e.Get("a")
	assertTrue(ok, "a found")
	assertEqual(v, "new", "a should be new (memtable wins)")
	v, _ = e.Get("b")
	assertEqual(v, "old", "b from sstable")
}

func testLSMEngineCompact() {
	e := NewLSMEngine(2)
	e.Put("a", "1")
	e.Put("b", "2") // flush -> sstable[0]
	e.Put("a", "3")
	e.Put("c", "4") // flush -> sstable[1]
	assertEqual(len(e.sstables), 2, "two sstables before compact")
	e.Compact()
	assertEqual(len(e.sstables), 1, "one sstable after compact")
	v, _ := e.Get("a")
	assertEqual(v, "3", "a should be 3 after compact (newest wins)")
	v, _ = e.Get("b")
	assertEqual(v, "2", "b should be 2")
	v, _ = e.Get("c")
	assertEqual(v, "4", "c should be 4")
}

func testLSMEngineCompactEmpty() {
	e := NewLSMEngine(5)
	e.Compact()
	assertEqual(len(e.sstables), 0, "no sstables when empty")
}

func testLSMEngineCompactSingle() {
	e := NewLSMEngine(2)
	e.Put("x", "1")
	e.Put("y", "2") // flush
	assertEqual(len(e.sstables), 1, "one sstable")
	e.Compact()
	assertEqual(len(e.sstables), 1, "still one sstable")
	v, _ := e.Get("x")
	assertEqual(v, "1", "x intact after compact")
}

func testLSMEngineMultipleFlushes() {
	e := NewLSMEngine(2)
	for i := 0; i < 10; i++ {
		e.Put(fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
	}
	assertEqual(len(e.sstables), 5, "five sstables after 10 inserts with cap 2")
	for i := 0; i < 10; i++ {
		v, ok := e.Get(fmt.Sprintf("key%d", i))
		assertTrue(ok, fmt.Sprintf("key%d found", i))
		assertEqual(v, fmt.Sprintf("val%d", i), fmt.Sprintf("key%d value", i))
	}
	e.Compact()
	assertEqual(len(e.sstables), 1, "one sstable after compact")
	for i := 0; i < 10; i++ {
		v, ok := e.Get(fmt.Sprintf("key%d", i))
		assertTrue(ok, fmt.Sprintf("key%d found after compact", i))
		assertEqual(v, fmt.Sprintf("val%d", i), fmt.Sprintf("key%d value after compact", i))
	}
}

func testSSTableBinarySearch() {
	sst := &SSTable{
		keys:   []string{"a", "c", "e", "g", "i"},
		values: []string{"1", "3", "5", "7", "9"},
	}
	v, ok := sst.Get("e")
	assertTrue(ok, "e found")
	assertEqual(v, "5", "e value")
	_, ok = sst.Get("d")
	assertTrue(!ok, "d not found")
	v, ok = sst.Get("a")
	assertTrue(ok, "a found (first)")
	assertEqual(v, "1", "a value")
	v, ok = sst.Get("i")
	assertTrue(ok, "i found (last)")
	assertEqual(v, "9", "i value")
}

func testMemTableEdgeCases() {
	mt := NewMemTable()
	// Empty get
	_, ok := mt.Get("anything")
	assertTrue(!ok, "empty memtable get returns false")
	// Single element
	mt.Insert("only", "1")
	assertEqual(mt.Len(), 1, "single element")
	v, ok := mt.Get("only")
	assertTrue(ok, "single element found")
	assertEqual(v, "1", "single element value")
	// Duplicate insert
	mt.Insert("only", "2")
	assertEqual(mt.Len(), 1, "still one element after dup insert")
	v, _ = mt.Get("only")
	assertEqual(v, "2", "updated single element")
}

func main() {
	testMemTableBasics()
	testMemTableFlush()
	testLSMEnginePutGet()
	testLSMEngineOverwriteAcrossFlush()
	testLSMEngineCompact()
	testLSMEngineCompactEmpty()
	testLSMEngineCompactSingle()
	testLSMEngineMultipleFlushes()
	testSSTableBinarySearch()
	testMemTableEdgeCases()
	fmt.Println("All LSM-Tree tests passed.")
}
