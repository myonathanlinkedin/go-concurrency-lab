package main

import (
	"fmt"
)

func main() {
	// Unit test 1: Tokenizer
	testTokenizer()

	// Unit test 2: Parser
	testParser()

	fmt.Println("All assertions passed.")
}

// testTokenizer validates token stream for a sample input.
func testTokenizer() {
	input := "let x = 42 + y;"
	lexer := NewLexer(input)

	expected := []Token{
		{Type: LET, Literal: "let"},
		{Type: IDENT, Literal: "x"},
		{Type: ASSIGN, Literal: "="},
		{Type: NUMBER, Literal: "42"},
		{Type: PLUS, Literal: "+"},
		{Type: IDENT, Literal: "y"},
		{Type: SEMIC, Literal: ";"},
		{Type: EOF, Literal: ""},
	}

	for i, exp := range expected {
		tok := lexer.NextToken()
		if tok.Type != exp.Type || tok.Literal != exp.Literal {
			panic(fmt.Sprintf("Tokenizer test failed at index %d: expected %v, got %v", i, exp, tok))
		}
	}
}

// testParser validates parsing of an arithmetic expression.
func testParser() {
	input := "a + 3 * (b - 2);"
	lexer := NewLexer(input)
	parser := NewParser(lexer)

	// Consume tokens until first expression (ignore trailing semicolon)
	expr := parser.ParseExpression()

	// Expected AST string: "(a + (3 * (b - 2)))"
	expected := "(a + (3 * (b - 2)))"
	if expr.String() != expected {
		panic(fmt.Sprintf("Parser test failed: expected %s, got %s", expected, expr.String()))
	}
}

// End of file.
