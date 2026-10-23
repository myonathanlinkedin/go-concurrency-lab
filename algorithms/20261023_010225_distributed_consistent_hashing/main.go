package main

import (
    "fmt"
)

func main() {
    // Create consistent hash with 3 virtual nodes per physical node
    ch := NewConsistentHash(3)

    // Test empty ring
    if ch.GetNode("anykey") != "" {
        panic("expected empty string for GetNode on empty ring")
    }

    // Add nodes
    ch.AddNode("node1")
    ch.AddNode("node2")
    ch.AddNode("node3")

    // Verify ring length
    expectedRingLen := 3 * 3 // 3 nodes * 3 virtual nodes
    if len(ch.ring) != expectedRingLen {
        panic(fmt.Sprintf("expected ring length %d, got %d", expectedRingLen, len(ch.ring)))
    }

    // Verify mapping consistency
    keys := []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}
    mapping := make(map[string]string)
    for _, key := range keys {
        node := ch.GetNode(key)
        if node == "" {
            panic(fmt.Sprintf("GetNode returned empty for key %s", key))
        }
        mapping[key] = node
    }

    // Recompute mapping to ensure determinism
    for _, key := range keys {
        node := ch.GetNode(key)
        if node != mapping[key] {
            panic(fmt.Sprintf("determinism failed for key %s: expected %s, got %s", key, mapping[key], node))
        }
    }

    // Remove a node and verify redistribution
    ch.RemoveNode("node2")

    // After removal, ring length should be 2*3
    if len(ch.ring) != 2*3 {
        panic(fmt.Sprintf("expected ring length %d after removal, got %d", 2*3, len(ch.ring)))
    }

    // Keys that previously mapped to node2 should now map to node1 or node3
    for _, key := range keys {
        node := ch.GetNode(key)
        if node == "" {
            panic(fmt.Sprintf("GetNode returned empty after removal for key %s", key))
        }
        if mapping[key] == "node2" && node == "node2" {
            panic(fmt.Sprintf("key %s still mapped to removed node2", key))
        }
    }

    // Add duplicate node and ensure no duplication
    ch.AddNode("node1")
    if len(ch.ring) != 3*3 { // should be back to 3 nodes * 3 virtual nodes
        panic(fmt.Sprintf("ring length after adding duplicate node1 should remain %d, got %d", 3*3, len(ch.ring)))
    }

    // Remove non-existent node
    ch.RemoveNode("nodeX")
    if len(ch.ring) != 3*3 {
        panic(fmt.Sprintf("ring length after removing non-existent node should remain %d, got %d", 3*3, len(ch.ring)))
    }

    // Verify that all nodes in mapping are still present
    for _, node := range mapping {
        if node != "node1" && node != "node3" {
            panic(fmt.Sprintf("unexpected node %s in mapping after removal", node))
        }
    }

    // Verify ring sorted property
    for i := 1; i < len(ch.ring); i++ {
        if ch.ring[i-1] > ch.ring[i] {
            panic("ring is not sorted")
        }
    }

    // Verify that GetNode for a key that was not present before still works
    newKey := "honeydew"
    node := ch.GetNode(newKey)
    if node == "" {
        panic("GetNode returned empty for new key");
}
}
