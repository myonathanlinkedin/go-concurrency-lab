package main

import (
	"math"
)

// eps for iterative algorithms.
const eps = 1e-12
const maxIter = 100

// SVD computes the singular value decomposition of A.
// Returns U (m×m), Σ (m×n) diagonal matrix of singular values, and Vt (n×n) = Vᵀ.
// The implementation follows the classic approach:
//   1. Compute B = Aᵀ·A (n×n, symmetric).
//   2. Eigen‑decompose B → eigenvalues λ_i and eigenvectors v_i.
//   3. Singular values σ_i = sqrt(λ_i) (sorted descending).
//   4. V = [v_1 … v_n].
//   5. U = A·V·Σ⁻¹ (columns corresponding to zero σ are left as zero).
func SVD(A Matrix) (U, S, Vt Matrix) {
	m, n := A.Rows(), A.Cols()
	if m == 0 || n == 0 {
		panic("empty matrix")
	}

	// Step 1: B = Aᵀ·A
	At := A.Transpose()
	B := At.Mul(A) // n×n symmetric

	// Step 2: Eigen‑decompose B using Jacobi rotations.
	eigVals, eigVecs := JacobiEigen(B)

	// Step 3: Singular values and sorting.
	type sv struct {
		val  float64
		vec  []float64
		idx  int
	}
	svs := make([]sv, n)
	for i := 0; i < n; i++ {
		svs[i] = sv{val: eigVals[i], vec: make([]float64, n), idx: i}
		for j := 0; j < n; j++ {
			svs[i].vec[j] = eigVecs[j][i] // column i of eigenvector matrix
		}
	}
	// Sort descending by singular value (sqrt of eigenvalue)
	// Using simple bubble sort because n is tiny in tests; replace with sort if needed.
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if svs[i].val < svs[j].val {
				svs[i], svs[j] = svs[j], svs[i]
			}
		}
	}
	// Build Σ and V
	S = NewMatrix(m, n)
	V := NewMatrix(n, n)
	for i := 0; i < n; i++ {
		sig := math.Sqrt(math.Max(svs[i].val, 0.0))
		S[i][i] = sig
		for j := 0; j < n; j++ {
			V[j][i] = svs[i].vec[j]
		}
	}
	Vt = V.Transpose()

	// Step 5: Compute U = A·V·Σ⁻¹
	AV := A.Mul(V) // m×n
	U = NewMatrix(m, m)
	// Build Σ⁻¹ (n×n) diagonal
	SInv := NewMatrix(n, n)
	for i := 0; i < n; i++ {
		if S[i][i] > eps {
			SInv[i][i] = 1.0 / S[i][i]
		}
	}
	AVSInv := AV.Mul(SInv) // m×n
	// Extend to full m×m orthogonal matrix.
	// For rank‑deficient cases we fill remaining columns with an orthonormal basis via Gram‑Schmidt.
	// Here we simply copy AVSInv into the first n columns of U and then orthonormalize.
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			U[i][j] = AVSInv[i][j]
		}
	}
	// Orthonormalize columns of U (Gram‑Schmidt)
	for i := 0; i < m; i++ {
		// Normalize column i
		norm := 0.0
		for k := 0; k < m; k++ {
			norm += U[k][i] * U[k][i]
		}
		if norm > eps {
			inv := 1.0 / math.Sqrt(norm)
			for k := 0; k < m; k++ {
				U[k][i] *= inv
			}
		}
		// Orthogonalize subsequent columns
		for j := i + 1; j < m; j++ {
			dot := 0.0
			for k := 0; k < m; k++ {
				dot += U[k][i] * U[k][j]
			}
			for k := 0; k < m; k++ {
				U[k][j] -= dot * U[k][i]
			}
		}
	}
	return
}

// JacobiEigen computes eigenvalues and eigenvectors of a symmetric matrix using the classic Jacobi method.
// Returns eigenvalues as a slice and eigenvectors as a matrix whose columns are the eigenvectors.
func JacobiEigen(A Matrix) (vals []float64, vecs Matrix) {
	n := A.Rows()
	if n != A.Cols() {
		panic("JacobiEigen requires a square matrix")
	}
	// Initialize eigenvector matrix as identity.
	vecs = Identity(n)
	// Work on a copy to avoid mutating the caller.
	D := A.Copy()

	for iter := 0; iter < maxIter; iter++ {
		// Find the largest off‑diagonal absolute element.
		p, q := -1, -1
		maxOff := 0.0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if math.Abs(D[i][j]) > maxOff {
					maxOff = math.Abs(D[i][j])
					p, q = i, j
				}
			}
		}
		if maxOff < eps {
			break // converged
		}
		// Compute rotation angle.
		if D[p][p] == D[q][q] {
			theta := math.Pi / 4
			cos, sin := math.Cos(theta), math.Sin(theta)
			applyJacobiRotation(D, vecs, p, q, cos, sin)
			continue
		}
		tau := (D[q][q] - D[p][p]) / (2.0 * D[p][q])
		t := 1.0 / (math.Abs(tau) + math.Sqrt(1.0+tau*tau))
		if tau < 0 {
			t = -t
		}
		cos := 1.0 / math.Sqrt(1.0+t*t)
		sin := t * cos
		applyJacobiRotation(D, vecs, p, q, cos, sin)
	}
	// Extract eigenvalues from diagonal.
	vals = make([]float64, n)
	for i := 0; i < n; i++ {
		vals[i] = D[i][i]
	}
	return
}

// applyJacobiRotation applies a single Jacobi rotation to matrix D and accumulates it into eigenvector matrix V.
func applyJacobiRotation(D, V Matrix, p, q int, c, s float64) {
	n := D.Rows()
	// Rotate rows and columns of D.
	for i := 0; i < n; i++ {
		if i != p && i != q {
			dp := D[i][p]
			dq := D[i][q]
			D[i][p] = c*dp - s*dq
			D[p][i] = D[i][p]
			D[i][q] = s*dp + c*dq
			D[q][i] = D[i][q]
		}
	}
	// Update diagonal elements.
	pp := D[p][p]
	qq := D[q][q]
	pq := D[p][q]
	D[p][p] = c*c*pp - 2*c*s*pq + s*s*qq
	D[q][q] = s*s*pp + 2*c*s*pq + c*c*qq
	D[p][q] = 0.0
	D[q][p] = 0.0

	// Accumulate rotation into V.
	for i := 0; i < n; i++ {
		viP := V[i][p]
		viQ := V[i][q]
		V[i][p] = c*viP - s*viQ
		V[i][q] = s*viP + c*viQ
	}
}
