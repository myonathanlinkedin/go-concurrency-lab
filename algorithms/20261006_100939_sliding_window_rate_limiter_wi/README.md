# Sliding Window Rate Limiter with Distributed Token Bucket in Go

An in-memory reference implementation of **Sliding Window Rate Limiter with Distributed Token Bucket** in **Go**, adhering to standard library idioms, clean data structures, and assertion test suites.

## Implementation Details

* **Category**: `Distributed Consensus & State Machine`
* **Data Structure Foundation**: `Append-Only State Log & Version Matrix`
* **Allocation Pattern**: Buffer boundaries and collection indices are explicitly validated to prevent out-of-bounds access.
* **Invariant Integrity**: State consistency is verified after mutations through assertion test coverage.

## Performance Characteristics

* **Time**: `$O(\log N) or O(1)$` average, with `$O(1)$` best-case response under ideal conditions.
* **Space**: `$O(N) state log$` memory usage.

## Test Harness

To compile and execute the test assertions for this module:

```bash
go run main.go
```

---

*Source code released under the MIT License • [@myonathanlinkedin](https://github.com/myonathanlinkedin)*