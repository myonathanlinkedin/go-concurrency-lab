package main

import (
	"fmt"
	"unicode"
)

// TokenType represents the category of lexical tokens.
type TokenType string

const (
	// Special tokens
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"

	// Literals
	IDENT  TokenType = "IDENT"  // variable names, function names
	NUMBER TokenType = "NUMBER" // integer literals

	// Operators
	PLUS     TokenType = "PLUS"
	MINUS    TokenType = "MINUS"
	ASTERISK TokenType = "ASTERISK"
	SLASH    TokenType = "SLASH"
	ASSIGN   TokenType = "ASSIGN"

	// Delimiters
	LPAREN TokenType = "LPAREN"
	RPAREN TokenType = "RPAREN"
	COMMA  TokenType = "COMMA"
	SEMIC  TokenType = "SEMIC"

	// Keywords
	FUNC TokenType = "FUNC"
	LET  TokenType = "LET"
)

// Token holds a lexical token with its type and literal representation.
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer implements a finite state machine based tokenizer.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position (after current char)
	ch           rune // current char under examination
}

// NewLexer creates a new Lexer for the given input string.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar advances the lexer by one Unicode code point.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // NUL signifies EOF
	} else {
		r, w := rune(l.input[l.readPosition]), 1
		// handle multibyte UTF-8
		if r >= 0x80 {
			r, w = utf8DecodeRuneInString(l.input[l.readPosition:])
		}
		l.ch = r
		l.readPosition += w
	}
	l.position = l.readPosition - 1
}

// peekChar returns the next character without consuming it.
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8DecodeRuneInString(l.input[l.readPosition:])
	return r
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var tok Token

	switch l.ch {
	case '+':
		tok = newToken(PLUS, l.ch)
	case '-':
		tok = newToken(MINUS, l.ch)
	case '*':
		tok = newToken(ASTERISK, l.ch)
	case '/':
		tok = newToken(SLASH, l.ch)
	case '=':
		tok = newToken(ASSIGN, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case ';':
		tok = newToken(SEMIC, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tt := lookupIdent(lit)
			return Token{Type: tt, Literal: lit}
		} else if unicode.IsDigit(l.ch) {
			return Token{Type: NUMBER, Literal: l.readNumber()}
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

// skipWhitespace consumes spaces, tabs, newlines.
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// readIdentifier reads a sequence of letters (and digits after first) as identifier.
func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) || unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readNumber reads a contiguous sequence of digits.
func (l *Lexer) readNumber() string {
	start := l.position
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// Helper functions

func newToken(tokenType TokenType, ch rune) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// lookupIdent checks if identifier is a keyword.
func lookupIdent(ident string) TokenType {
	switch ident {
	case "fn":
		return FUNC
	case "let":
		return LET
	default:
		return IDENT
	}
}

// utf8DecodeRuneInString is a tiny wrapper to avoid importing "unicode/utf8" in every call.
func utf8DecodeRuneInString(s string) (r rune, size int) {
	// Inline the standard library call to keep imports minimal.
	return []rune(s)[0], len(string([]rune(s)[0]))
}

// AST node definitions

type Node interface {
	String() string
}

// Expression nodes

type Expr interface {
	Node
	exprNode()
}

// Identifier node
type Identifier struct {
	Value string
}

func (i *Identifier) exprNode() {}
func (i *Identifier) String() string {
	return i.Value
}

// Number literal node
type NumberLiteral struct {
	Value string
}

func (n *NumberLiteral) exprNode() {}
func (n *NumberLiteral) String() string {
	return n.Value
}

// Binary operation node
type BinaryOp struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *BinaryOp) exprNode() {}
func (b *BinaryOp) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Operator.Literal, b.Right.String())
}

// Parser implements a Pratt parser using precedence climbing.
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

// NewParser creates a parser with the given lexer.
func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	// Load two tokens, curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances tokens.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseExpression parses an expression with the lowest precedence.
func (p *Parser) ParseExpression() Expr {
	return p.parseExpression(0)
}

// precedence levels
var precedences = map[TokenType]int{
	PLUS:     1,
	MINUS:    1,
	ASTERISK: 2,
	SLASH:    2,
}

// curPrecedence returns precedence of current token.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return 0
}

// peekPrecedence returns precedence of peek token.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return 0
}

// parseExpression parses based on precedence.
func (p *Parser) parseExpression(precedence int) Expr {
	var left Expr

	switch p.curToken.Type {
	case IDENT:
		left = &Identifier{Value: p.curToken.Literal}
	case NUMBER:
		left = &NumberLiteral{Value: p.curToken.Literal}
	case LPAREN:
		p.nextToken()
		left = p.parseExpression(0)
		if p.curToken.Type != RPAREN {
			panic("expected closing parenthesis")
		}
	default:
		panic(fmt.Sprintf("unexpected token: %s", p.curToken.Type))
	}

	for p.peekToken.Type != SEMIC && p.peekToken.Type != EOF && precedence < p.peekPrecedence() {
		op := p.peekToken
		if op.Type != PLUS && op.Type != MINUS && op.Type != ASTERISK && op.Type != SLASH {
			return left
		}
		p.nextToken() // advance to operator
		p.nextToken() // advance to right-hand side
		right := p.parseExpression(precedences[op.Type])
		left = &BinaryOp{Left: left, Operator: op, Right: right}
	}
	return left
}
