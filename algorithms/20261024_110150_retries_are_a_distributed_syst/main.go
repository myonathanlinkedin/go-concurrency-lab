package main

import (
	"fmt"
)

func main() {
	store := NewResourceStore()
	engine := NewEngine(store)

	// Test 1: Idempotent create
	req1 := Request{IdempotencyKey: "key1", Operation: "create", Payload: "res1:Hello"}
	resp1 := engine.ProcessRequest(req1)
	if resp1.Status != 201 {
		panic(fmt.Sprintf("Test1 failed: expected 201, got %d", resp1.Status))
	}

	// Retry same request
	resp1Retry := engine.ProcessRequest(req1)
	if resp1Retry.Status != 201 || resp1Retry.Body != resp1.Body {
		panic("Test1 retry failed: response mismatch")
	}

	// Test 2: Different key same payload
	req2 := Request{IdempotencyKey: "key2", Operation: "create", Payload: "res1:Hello"}
	resp2 := engine.ProcessRequest(req2)
	if resp2.Status != 409 {
		panic(fmt.Sprintf("Test2 failed: expected 409, got %d", resp2.Status))
	}

	// Test 3: Get existing resource
	req3 := Request{IdempotencyKey: "key3", Operation: "get", Payload: "res1"}
	resp3 := engine.ProcessRequest(req3)
	if resp3.Status != 200 || resp3.Body != "resource res1: Hello" {
		panic("Test3 failed: get resource")
	}

	// Test 4: Get non-existent resource
	req4 := Request{IdempotencyKey: "key4", Operation: "get", Payload: "res2"}
	resp4 := engine.ProcessRequest(req4)
	if resp4.Status != 404 {
		panic("Test4 failed: expected 404")
	}

	// Test 5: Unknown operation
	req5 := Request{IdempotencyKey: "key5", Operation: "delete", Payload: "res1"}
	resp5 := engine.ProcessRequest(req5)
	if resp5.Status != 400 {
		panic("Test5 failed: expected 400")
	}

	// Test 6: Idempotency key reuse with different payload
	req6a := Request{IdempotencyKey: "key6", Operation: "create", Payload: "res2:World"}
	resp6a := engine.ProcessRequest(req6a)
	if resp6a.Status != 201 {
		panic("Test6a failed")
	}
	req6b := Request{IdempotencyKey: "key6", Operation: "create", Payload: "res3:World"}
	resp6b := engine.ProcessRequest(req6b)
	if resp6b.Status != 201 || resp6b.Body != resp6a.Body {
		panic("Test6b failed: idempotency key should ignore payload")
	}

	fmt.Println("All tests passed")
}
