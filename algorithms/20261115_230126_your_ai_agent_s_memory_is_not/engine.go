package main

// Set stores the given value under the supplied key.
// It atomically increments the global version counter and records that
// version for the key. The new version is returned.
func (s *Store) Set(key, value string) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.globalVersion++
	ve := valueEntry{
		value:   value,
		version: s.globalVersion,
	}
	s.data[key] = ve
	return ve.version
}

// Get retrieves the value and its version for the supplied key.
// The boolean indicates whether the key existed.
func (s *Store) Get(key string) (string, uint64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ve, ok := s.data[key]
	if !ok {
		return "", 0, false
	}
	return ve.value, ve.version, true
}

// IsStale reports whether the caller's knownVersion for the key
// is older than the current version stored in the Store.
// If the key does not exist, it returns false (nothing to be stale about).
func (s *Store) IsStale(key string, knownVersion uint64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ve, ok := s.data[key]
	if !ok {
		return false
	}
	return knownVersion < ve.version
}
