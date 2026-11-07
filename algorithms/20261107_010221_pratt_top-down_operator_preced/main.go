package main

import (
	"fmt"
	"math"
)

func approxEqual(a, b float64) bool {
	const eps = 1e-9
	return math.Abs(a-b) <= eps
}

func main() {
	type testCase struct {
		expr     string
		expected float64
	}
	tests := []testCase{
		{"1+2", 3},
		{"2*3+4", 10},
		{"2+(3*4)", 14},
		{"(2+3)*4", 20},
		{"-5+3", -2},
		{"- (4 + 1) * 2", -10},
		{"2^3^2", 512}, // right‑associative: 2^(3^2)=2^9=512
		{"3+4*2/(1-5)^2^3", 3.0001220703125},
		{"3.5*2", 7},
		{"(1+2)*(3+4)", 21},
	}
	for i, tc := range tests {
		ast := Parse(tc.expr)
		result := ast.Eval()
		if !approxEqual(result, tc.expected) {
			panic(fmt.Sprintf("test %d failed: expr=%s got=%v want=%v", i, tc.expr, result, tc.expected))
		}
	}
	// Additional edge case: multiple unary minuses
	ast := Parse("--5")
	if !approxEqual(ast.Eval(), 5) {
		panic("double unary minus failed")
	}
	// Ensure parser rejects malformed input
	func() {
		defer func() {
			if r := recover(); r == nil {
				panic("expected panic on malformed expression")
			}
		}()
		Parse("2+*3")
	}()
	fmt.Println("All tests passed.")
}
