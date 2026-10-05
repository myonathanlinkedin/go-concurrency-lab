# Show HN: Glashütte Trash Clock – A 30-minute pendulum clock made from trash

A clean, dependency-free **Go** implementation of **Show HN: Glashütte Trash Clock – A 30-minute pendulum clock made from trash**, focused on predictable latency, strict memory layout, and deterministic execution.

---

## 🏛️ Architecture & Design Decisions

This module organizes `Show HN: Glashütte Trash Clock – A 30-minute pendulum clock made from trash` into an isolated, self-contained unit:
* **Domain Focus**: `Low-Latency Systems & Memory Layout`
* **Primary Primitives**: `Contiguous Memory Buffer & Ring Pointers`
* **Memory Strategy**: Zero superfluous dynamic allocations; structured for mechanical sympathy with the host runtime.
* **Correctness Model**: Designed with reentrancy and thread isolation in mind, preventing data races under parallel workloads.

### Asymptotic Complexity

| Metric | Bound | Characteristics |
| :--- | :---: | :--- |
| **Best Case Time** | `$O(1)$` | Optimized fast-path execution |
| **Average / Worst Time** | `$O(1)$` | Deterministic upper bound for generalized workloads |
| **Space Complexity** | `$O(N) bounded$` | Strict bounds without unconstrained heap growth |

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