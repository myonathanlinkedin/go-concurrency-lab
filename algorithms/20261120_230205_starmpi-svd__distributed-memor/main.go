package main

import (
	"fmt"
	"math"
)

func main() {
	// Test 1: Simple 2×3 matrix.
	A := Matrix{
		{3, 1, 1},
		{-1, 3, 1},
	}
	U, S, Vt := SVD(A)

	// Reconstruct Â = U·S·Vᵀ
	US := U.Mul(S)
	Recon := US.Mul(Vt)

	// Verify reconstruction accuracy.
	const tol = 1e-9
	if !Recon.ApproxEqual(A, tol) {
		fmt.Printf("Reconstruction error exceeds tolerance.\nOriginal: %v\nReconstructed: %v\n", A, Recon)
		panic("assertion failed")
	}

	// Verify singular values are non‑increasing.
	sVals := []float64{S[0][0], S[1][1]}
	if sVals[0]+tol < sVals[1] {
		fmt.Printf("Singular values not sorted descending: %v\n", sVals)
		panic("assertion failed")
	}

	// Verify orthogonality of U (Uᵀ·U ≈ I)
	UtU := U.Transpose().Mul(U)
	I := Identity(U.Rows())
	if !UtU.ApproxEqual(I, 1e-9) {
		fmt.Printf("U is not orthogonal. Uᵀ·U = %v\n", UtU)
		panic("assertion failed")
	}

	// Verify orthogonality of V (Vᵀ·V ≈ I)
	V := Vt.Transpose()
	VtV := Vt.Mul(V)
	if !VtV.ApproxEqual(I, 1e-9) {
		fmt.Printf("V is not orthogonal. Vᵀ·V = %v\n", VtV)
		panic("assertion failed")
	}

	// Edge case: Zero matrix.
	Z := NewMatrix(3, 3)
	_, Sz, _ := SVD(Z)
	for i := 0; i < Sz.Rows(); i++ {
		for j := 0; j < Sz.Cols(); j++ {
			if math.Abs(Sz[i][j]) > tol {
				panic("assertion failed: singular values of zero matrix should be zero")
			}
		}
	}

	fmt.Println("All assertions passed.")
}
