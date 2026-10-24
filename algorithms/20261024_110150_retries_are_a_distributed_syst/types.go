package main

import (
	"fmt"
	"sync"
)

type Request struct {
	IdempotencyKey string
	Operation      string
	Payload        string
}

type Response struct {
	Status int
	Body   string
}

type ResourceStore struct {
	mu        sync.Mutex
	resources map[string]string
}

func NewResourceStore() *ResourceStore {
	return &ResourceStore{
		resources: make(map[string]string),
	}
}

func (rs *ResourceStore) Create(id string, value string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if _, exists := rs.resources[id]; exists {
		return fmt.Errorf("resource exists")
	}
	rs.resources[id] = value
	return nil
}

func (rs *ResourceStore) Get(id string) (string, bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	val, ok := rs.resources[id]
	return val, ok
}
