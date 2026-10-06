# Thread-Safe Bounded Blocking Queue with Condition Variables

An in-memory reference implementation of **Thread-Safe Bounded Blocking Queue with Condition Variables** in **Go**, adhering to standard library idioms, clean data structures, and assertion test suites.

---

## 🏛️ Architecture & Design Decisions

This module organizes `Thread-Safe Bounded Blocking Queue with Condition Variables` into an isolated, self-contained unit:
* **Domain Focus**: `Low-Latency Systems & Memory Layout`
* **Primary Primitives**: `Contiguous Memory Buffer & Ring Pointers`
* **Memory Strategy**: Memory allocations are kept minimal to maintain clear data locality and predictable memory bounds.
* **Correctness Model**: State transitions follow clear ordering guarantees with explicit validation at each phase.

### Asymptotic Complexity

| Metric | Bound | Characteristics |
| :--- | :---: | :--- |
| **Best Case Time** | `O(1)` | Optimized fast-path execution |
| **Average / Worst Time** | `O(1)` | Deterministic upper bound for generalized workloads |
| **Space Complexity** | `O(N) bounded` | Strict bounds without unconstrained heap growth |

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

*Reference implementation verified by [@myonathanlinkedin](https://github.com/myonathanlinkedin)*
