# LSM-Tree MemTable and SSTable Flush Compaction Engine

Self-contained **LSM-Tree MemTable and SSTable Flush Compaction Engine** algorithmic primitive written in idiomatic **Go**. Built from scratch using standard library constructs with zero external dependencies.

### Core Highlights
* **Language & Standard**: Modern `Go` standard library conventions.
* **Architecture Pattern**: Designed for `Balanced Hierarchical Indexing` using `Node Pointers & Self-Balancing Trees`.
* **Runtime Overhead**: Contiguous memory layouts and standard collections are favored for straightforward iteration and access.
* **Concurrency & Safety**: State consistency is verified after mutations through assertion test coverage.

---

### Complexity Analysis

| Dimension | Bound |
| :--- | :--- |
| **Time (Best Case)** | `O(1)` |
| **Time (Worst Case)** | `O(log N)` |
| **Auxiliary Space** | `O(N)` |

---

### Test Suite Execution

Self-contained verification drivers are embedded directly in `main.go` to validate happy paths, boundary inputs, and invariant preservation.

```bash
go run main.go
```

---

<sub>Standard Go reference implementation • Maintained by [@myonathanlinkedin](https://github.com/myonathanlinkedin)</sub>