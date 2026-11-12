package main

import (
	"container/list"
	"sync"
)

// entry represents a key-value pair stored in the LRU cache.
type entry struct {
	key   string
	value string
}

// MobileStorage provides a local LRU cache backed by a simulated remote storage.
// All operations are safe for concurrent use.
type MobileStorage struct {
	capacity int
	cache    map[string]*list.Element // maps key to list element
	lru      *list.List               // most-recently used at front
	remote   map[string]string        // simulated remote storage
	mu       sync.RWMutex
}

// NewMobileStorage creates a MobileStorage with the given cache capacity.
// Capacity must be positive; otherwise the function panics.
func NewMobileStorage(capacity int) *MobileStorage {
	if capacity <= 0 {
		panic("capacity must be positive")
	}
	return &MobileStorage{
		capacity: capacity,
		cache:    make(map[string]*list.Element, capacity),
		lru:      list.New(),
		remote:   make(map[string]string),
	}
}

// Get retrieves the value for a key. It first checks the local cache; on a miss it
// fetches from remote storage, populates the cache (evicting if necessary), and returns.
func (ms *MobileStorage) Get(key string) (string, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Check local cache.
	if elem, ok := ms.cache[key]; ok {
		ms.lru.MoveToFront(elem)
		en := elem.Value.(*entry)
		return en.value, true
	}

	// Cache miss: consult remote storage.
	val, ok := ms.remote[key]
	if !ok {
		return "", false
	}
	// Populate cache with the remote value.
	ms.addToCache(key, val)
	return val, true
}

// Put stores the key-value pair in both remote storage and the local cache.
// If the cache exceeds its capacity, the least-recently used entry is evicted.
func (ms *MobileStorage) Put(key, value string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Update remote storage.
	ms.remote[key] = value

	// Update cache.
	if elem, ok := ms.cache[key]; ok {
		// Key already cached: update value and move to front.
		en := elem.Value.(*entry)
		en.value = value
		ms.lru.MoveToFront(elem)
		return
	}
	// Not in cache: add it.
	ms.addToCache(key, value)
}

// Delete removes a key from both remote storage and the local cache.
func (ms *MobileStorage) Delete(key string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	delete(ms.remote, key)

	if elem, ok := ms.cache[key]; ok {
		ms.removeElement(elem)
	}
}

// addToCache inserts a new entry into the cache, evicting the LRU entry if needed.
func (ms *MobileStorage) addToCache(key, value string) {
	if len(ms.cache) >= ms.capacity {
		// Evict least-recently used (back of list).
		if back := ms.lru.Back(); back != nil {
			ms.removeElement(back)
		}
	}
	en := &entry{key: key, value: value}
	elem := ms.lru.PushFront(en)
	ms.cache[key] = elem
}

// removeElement evicts the given list element from the cache.
func (ms *MobileStorage) removeElement(elem *list.Element) {
	ms.lru.Remove(elem)
	en := elem.Value.(*entry)
	delete(ms.cache, en.key)
}

// Size returns the number of items currently stored in the remote storage.
// This method is provided for testing purposes.
func (ms *MobileStorage) Size() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.remote)
}

// CacheSize returns the number of items currently cached locally.
// This method is provided for testing purposes.
func (ms *MobileStorage) CacheSize() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.cache)
}
