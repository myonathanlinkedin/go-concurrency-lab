package main

import (
	"fmt"
	"strings"
	"sync"
)

type Engine struct {
	store          *ResourceStore
	idempotencyMap map[string]Response
	mu             sync.Mutex
}

func NewEngine(store *ResourceStore) *Engine {
	return &Engine{
		store:          store,
		idempotencyMap: make(map[string]Response),
	}
}

func (e *Engine) ProcessRequest(req Request) Response {
	e.mu.Lock()
	defer e.mu.Unlock()

	if resp, exists := e.idempotencyMap[req.IdempotencyKey]; exists {
		return resp
	}

	var resp Response

	switch req.Operation {
	case "create":
		parts := strings.SplitN(req.Payload, ":", 2)
		if len(parts) != 2 {
			resp = Response{Status: 400, Body: "invalid payload"}
		} else {
			id, val := parts[0], parts[1]
			err := e.store.Create(id, val)
			if err != nil {
				resp = Response{Status: 409, Body: err.Error()}
			} else {
				resp = Response{Status: 201, Body: fmt.Sprintf("resource %s created", id)}
			}
		}
	case "get":
		id := req.Payload
		if val, ok := e.store.Get(id); ok {
			resp = Response{Status: 200, Body: fmt.Sprintf("resource %s: %s", id, val)}
		} else {
			resp = Response{Status: 404, Body: "resource not found"}
		}
	default:
		resp = Response{Status: 400, Body: "unknown operation"}
	}

	e.idempotencyMap[req.IdempotencyKey] = resp
	return resp
}
