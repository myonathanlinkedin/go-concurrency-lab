package main

import (
	"fmt"
)

func main() {
	// Test cases for the tokenizer and parser.
	tests := []struct {
		input          string
		expectedTokens []Token
		expectedStmts  []*AssignStmt
	}{
		{
			input: "x=42;",
			expectedTokens: []Token{
				{Type: TokenIdentifier, Lexeme: "x", Pos: 0},
				{Type: TokenOperator, Lexeme: "=", Pos: 1},
				{Type: TokenNumber, Lexeme: "42", Pos: 2},
				{Type: TokenOperator, Lexeme: ";", Pos: 4},
				{Type: TokenEOF, Lexeme: "", Pos: 5},
			},
			expectedStmts: []*AssignStmt{
				{Var: "x", Value: 42},
			},
		},
		{
			input: "var123 =  99 ;",
			expectedTokens: []Token{
				{Type: TokenIdentifier, Lexeme: "var123", Pos: 0},
				{Type: TokenWhitespace, Lexeme: " ", Pos: 6},
				{Type: TokenOperator, Lexeme: "=", Pos: 7},
				{Type: TokenWhitespace, Lexeme: "  ", Pos: 8},
				{Type: TokenNumber, Lexeme: "99", Pos: 10},
				{Type: TokenWhitespace, Lexeme: " ", Pos: 12},
				{Type: TokenOperator, Lexeme: ";", Pos: 13},
				{Type: TokenEOF, Lexeme: "", Pos: 14},
			},
			expectedStmts: []*AssignStmt{
				{Var: "var123", Value: 99},
			},
		},
		{
			input: "a=1;b=2;",
			expectedTokens: []Token{
				{Type: TokenIdentifier, Lexeme: "a", Pos: 0},
				{Type: TokenOperator, Lexeme: "=", Pos: 1},
				{Type: TokenNumber, Lexeme: "1", Pos: 2},
				{Type: TokenOperator, Lexeme: ";", Pos: 3},
				{Type: TokenIdentifier, Lexeme: "b", Pos: 4},
				{Type: TokenOperator, Lexeme: "=", Pos: 5},
				{Type: TokenNumber, Lexeme: "2", Pos: 6},
				{Type: TokenOperator, Lexeme: ";", Pos: 7},
				{Type: TokenEOF, Lexeme: "", Pos: 8},
			},
			expectedStmts: []*AssignStmt{
				{Var: "a", Value: 1},
				{Var: "b", Value: 2},
			},
		},
	}

	for idx, tc := range tests {
		toks, err := Tokenize(tc.input)
		if err != nil {
			panic(fmt.Sprintf("test %d: tokenization error: %v", idx, err))
		}
		if len(toks) != len(tc.expectedTokens) {
			panic(fmt.Sprintf("test %d: token count mismatch: got %d, want %d", idx, len(toks), len(tc.expectedTokens)))
		}
		for i := range toks {
			exp := tc.expectedTokens[i]
			got := toks[i]
			if got.Type != exp.Type || got.Lexeme != exp.Lexeme || got.Pos != exp.Pos {
				panic(fmt.Sprintf("test %d: token %d mismatch: got %+v, want %+v", idx, i, got, exp))
			}
		}
		stmts, err := Parse(toks)
		if err != nil {
			panic(fmt.Sprintf("test %d: parse error: %v", idx, err))
		}
		if len(stmts) != len(tc.expectedStmts) {
			panic(fmt.Sprintf("test %d: stmt count mismatch: got %d, want %d", idx, len(stmts), len(tc.expectedStmts)))
		}
		for i := range stmts {
			exp := tc.expectedStmts[i]
			got := stmts[i]
			if got.Var != exp.Var || got.Value != exp.Value {
				panic(fmt.Sprintf("test %d: stmt %d mismatch: got %+v, want %+v", idx, i, got, exp))
			}
		}
	}
	fmt.Println("All tests passed.")
}
