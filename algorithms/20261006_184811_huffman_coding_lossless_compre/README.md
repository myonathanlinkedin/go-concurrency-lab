# Huffman Coding Lossless Compression and Decompression

Core **Go** implementation for **Huffman Coding Lossless Compression and Decompression**, structured for computational clarity, explicit data structures, and deterministic unit test coverage.

### Core Highlights
* **Language & Standard**: Modern `Go` standard library conventions.
* **Architecture Pattern**: Designed for `Computational Mathematics & Transformation` using `Lookup Tables & Bitwise Bitvectors`.
* **Runtime Overhead**: Zero external heap dependencies; designed as a pure in-memory algorithmic component.
* **Concurrency & Safety**: Execution behavior is validated against nominal workflows and boundary edge cases.

---

### Complexity Analysis

| Dimension | Bound |
| :--- | :--- |
| **Time (Best Case)** | `O(N log N)` |
| **Time (Worst Case)** | `O(N * M)` |
| **Auxiliary Space** | `O(N)` |

---

### Test Suite Execution

Self-contained verification drivers are embedded directly in `main.go` to validate happy paths, boundary inputs, and invariant preservation.

```bash
go run main.go
```

---

*Part of the Polyglot Systems Lab • Maintained by [@myonathanlinkedin](https://github.com/myonathanlinkedin)*