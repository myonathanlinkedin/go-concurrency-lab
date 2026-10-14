package main

import (
	"fmt"
)

func main() {
	// Test programs
	tests := []struct {
		src      string
		expected int
	}{
		{
			src:      "a=3+4*2; b=a*5; b",
			expected: 55,
		},
		{
			src:      "x=10; y= (x-3) * 2; y + 1",
			expected: 15,
		},
		{
			src:      "n=7; n* n",
			expected: 49,
		},
		{
			src:      "p=8/2; q=p+1; q*3",
			expected: 15,
		},
	}

	for idx, tt := range tests {
		// Parse
		parser := NewParser(tt.src)
		stmts := parser.Parse()

		// Interpret
		interp := NewInterpreter()
		interpResult := interp.Eval(stmts)

		// Compile & Run VM
		compiler := NewCompiler()
		instr := compiler.Compile(stmts)
		vm := NewVM(instr, len(compiler.varNames))
		vmResult := vm.Run()

		// Assertions
		if interpResult != tt.expected {
			panic(fmt.Sprintf("test %d interpreter failed: got %d, want %d", idx, interpResult, tt.expected))
		}
		if vmResult != tt.expected {
			panic(fmt.Sprintf("test %d VM failed: got %d, want %d", idx, vmResult, tt.expected))
		}
		// Ensure both paths agree
		if interpResult != vmResult {
			panic(fmt.Sprintf("test %d mismatch between interpreter and VM: %d vs %d", idx, interpResult, vmResult))
		}
		fmt.Printf("Test %d passed: result=%d\n", idx, tt.expected)
	}
	// Edge case: division by zero should panic (uncomment to verify)
	// parser := NewParser("1/0")
	// stmts := parser.Parse()
	// interp := NewInterpreter()
	// _ = interp.Eval(stmts) // should panic
}
