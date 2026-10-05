# Hopcroft-Karp Bipartite Matching Algorithm

High-performance **Hopcroft-Karp Bipartite Matching Algorithm** primitive implemented in idiomatic **Go**. Built from scratch using standard library constructs with zero external dependencies.

---

## 🏛️ Architecture & Design Decisions

This module organizes `Hopcroft-Karp Bipartite Matching Algorithm` into an isolated, self-contained unit:
* **Domain Focus**: `Algorithmic Engineering`
* **Primary Primitives**: `Standard Memory Primitives`
* **Memory Strategy**: Memory allocations are kept minimal to avoid allocator contention and preserve CPU cache locality.
* **Correctness Model**: State transitions adhere to strict ordering guarantees with explicit synchronization fences where necessary.

### Asymptotic Complexity

| Metric | Bound | Characteristics |
| :--- | :---: | :--- |
| **Best Case Time** | `$O(1)$` | Optimized fast-path execution |
| **Average / Worst Time** | `$O(N)$` | Deterministic upper bound for generalized workloads |
| **Space Complexity** | `$O(N)$` | Strict bounds without unconstrained heap growth |

---

## 🧪 Verification Suite

The accompanying `main.go` driver executes self-contained verification tests:
1. **Nominal Flow**: Validates baseline correctness under typical real-world inputs.
2. **Boundary Conditions**: Exercises extreme edge cases (empty inputs, singletons, capacity limits).
3. **Invariant Preservation**: Validates internal state consistency throughout mutation lifecycles.

### Running Locally

```bash
go run main.go
```

---

*Authored & verified by [@myonathanlinkedin](https://github.com/myonathanlinkedin) • Systems Engineering Portfolio*