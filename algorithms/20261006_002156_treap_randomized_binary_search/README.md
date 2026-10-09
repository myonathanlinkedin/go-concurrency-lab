# Treap Randomized Binary Search Tree with Heap Priorities

A clean, dependency-free **Go** reference implementation of **Treap Randomized Binary Search Tree with Heap Priorities**, focused on core algorithmic mechanics, clear memory layout, and test verification.

---

## 🏛️ Architecture & Design Decisions

This module organizes `Treap Randomized Binary Search Tree with Heap Priorities` into an isolated, self-contained unit:
* **Domain Focus**: `Balanced Hierarchical Indexing`
* **Primary Primitives**: `Node Pointers & Self-Balancing Trees`
* **Memory Strategy**: Zero external heap dependencies; designed as a pure in-memory algorithmic component.
* **Correctness Model**: State consistency is verified after mutations through assertion test coverage.

### Asymptotic Complexity

| Metric | Bound | Characteristics |
| :--- | :---: | :--- |
| **Best Case Time** | `O(1)` | Optimized fast-path execution |
| **Average / Worst Time** | `O(log N)` | Deterministic upper bound for generalized workloads |
| **Space Complexity** | `O(N)` | Strict bounds without unconstrained heap growth |

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

*Part of the Polyglot Systems Lab • Maintained by [@myonathanlinkedin](https://github.com/myonathanlinkedin)*
