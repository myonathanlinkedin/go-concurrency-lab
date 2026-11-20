package main

import (
	"math"
)

// Matrix is a simple 2‑dimensional slice of float64.
type Matrix [][]float64

// NewMatrix creates a rows×cols matrix initialized to zero.
func NewMatrix(rows, cols int) Matrix {
	m := make(Matrix, rows)
	for i := range m {
		m[i] = make([]float64, cols)
	}
	return m
}

// Rows returns the number of rows.
func (m Matrix) Rows() int { return len(m) }

// Cols returns the number of columns (0 if no rows).
func (m Matrix) Cols() int {
	if len(m) == 0 {
		return 0
	}
	return len(m[0])
}

// Copy returns a deep copy of the matrix.
func (m Matrix) Copy() Matrix {
	r, c := m.Rows(), m.Cols()
	n := NewMatrix(r, c)
	for i := 0; i < r; i++ {
		copy(n[i], m[i])
	}
	return n
}

// Transpose returns the transpose of the matrix.
func (m Matrix) Transpose() Matrix {
	r, c := m.Rows(), m.Cols()
	t := NewMatrix(c, r)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			t[j][i] = m[i][j]
		}
	}
	return t
}

// Mul returns the matrix product m * n.
// Panics if dimensions are incompatible.
func (m Matrix) Mul(n Matrix) Matrix {
	if m.Cols() != n.Rows() {
		panic("matrix dimension mismatch for multiplication")
	}
	r, c, inner := m.Rows(), n.Cols(), m.Cols()
	res := NewMatrix(r, c)
	for i := 0; i < r; i++ {
		for k := 0; k < inner; k++ {
			tmp := m[i][k]
			for j := 0; j < c; j++ {
				res[i][j] += tmp * n[k][j]
			}
		}
	}
	return res
}

// Scale returns a new matrix where each element is multiplied by s.
func (m Matrix) Scale(s float64) Matrix {
	r, c := m.Rows(), m.Cols()
	res := NewMatrix(r, c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			res[i][j] = m[i][j] * s
		}
	}
	return res
}

// Identity creates an identity matrix of size n×n.
func Identity(n int) Matrix {
	id := NewMatrix(n, n)
	for i := 0; i < n; i++ {
		id[i][i] = 1.0
	}
	return id
}

// FrobeniusNorm returns the Frobenius norm of the matrix.
func (m Matrix) FrobeniusNorm() float64 {
	sum := 0.0
	for i := 0; i < m.Rows(); i++ {
		for j := 0; j < m.Cols(); j++ {
			sum += m[i][j] * m[i][j]
		}
	}
	return math.Sqrt(sum)
}

// ApproxEqual checks whether two matrices are equal within eps.
func (a Matrix) ApproxEqual(b Matrix, eps float64) bool {
	if a.Rows() != b.Rows() || a.Cols() != b.Cols() {
		return false
	}
	for i := 0; i < a.Rows(); i++ {
		for j := 0; j < a.Cols(); j++ {
			if math.Abs(a[i][j]-b[i][j]) > eps {
				return false
			}
		}
	}
	return true
}
