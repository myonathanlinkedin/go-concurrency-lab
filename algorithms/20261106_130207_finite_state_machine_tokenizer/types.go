package main

import "fmt"

// TokenType enumerates the categories of lexical tokens.
type TokenType int

const (
	TokenInvalid TokenType = iota
	TokenIdentifier
	TokenNumber
	TokenOperator
	TokenWhitespace
	TokenEOF
)

// Token represents a lexical token with its type, literal text, and start position.
type Token struct {
	Type   TokenType
	Lexeme string
	Pos    int // byte offset in the original input
}

// String implements fmt.Stringer for debugging.
func (t Token) String() string {
	return fmt.Sprintf("{Type:%d Lexeme:%q Pos:%d}", t.Type, t.Lexeme, t.Pos)
}

// AssignStmt is a minimal AST node representing an assignment like `x = 42;`.
type AssignStmt struct {
	Var   string
	Value int
}

// String implements fmt.Stringer for debugging.
func (a *AssignStmt) String() string {
	return fmt.Sprintf("%s = %d;", a.Var, a.Value)
}
