package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenOperator
	TokenLParen
	TokenRParen
)

type Token struct {
	typ TokenType
	val string
}

// precedence map for binary operators
var precedence = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
	"^": 3,
}

// Node is the AST interface
type Node interface {
	Eval() float64
}

// NumberNode represents a literal number
type NumberNode struct {
	Value float64
}

func (n *NumberNode) Eval() float64 { return n.Value }

// UnaryNode represents a prefix unary operator (currently only '-')
type UnaryNode struct {
	Op    string
	Right Node
}

func (u *UnaryNode) Eval() float64 {
	switch u.Op {
	case "-":
		return -u.Right.Eval()
	default:
		panic(fmt.Sprintf("unsupported unary operator %s", u.Op))
	}
}

// BinaryNode represents a binary operator
type BinaryNode struct {
	Op    string
	Left  Node
	Right Node
}

func (b *BinaryNode) Eval() float64 {
	switch b.Op {
	case "+":
		return b.Left.Eval() + b.Right.Eval()
	case "-":
		return b.Left.Eval() - b.Right.Eval()
	case "*":
		return b.Left.Eval() * b.Right.Eval()
	case "/":
		return b.Left.Eval() / b.Right.Eval()
	case "^":
		return math.Pow(b.Left.Eval(), b.Right.Eval())
	default:
		panic(fmt.Sprintf("unsupported binary operator %s", b.Op))
	}
}

// Lexer converts an input string into a slice of Tokens
func Lex(input string) []Token {
	var tokens []Token
	i := 0
	for i < len(input) {
		ch := input[i]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			i++
		case ch >= '0' && ch <= '9' || ch == '.':
			start := i
			for i < len(input) && ((input[i] >= '0' && input[i] <= '9') || input[i] == '.') {
				i++
			}
			tokens = append(tokens, Token{typ: TokenNumber, val: input[start:i]})
		case strings.ContainsRune("+-*/^", rune(ch)):
			tokens = append(tokens, Token{typ: TokenOperator, val: string(ch)})
			i++
		case ch == '(':
			tokens = append(tokens, Token{typ: TokenLParen, val: "("})
			i++
		case ch == ')':
			tokens = append(tokens, Token{typ: TokenRParen, val: ")"})
			i++
		default:
			panic(fmt.Sprintf("invalid character %q at position %d", ch, i))
		}
	}
	tokens = append(tokens, Token{typ: TokenEOF, val: ""})
	return tokens
}

// Parser implements Pratt top‑down parsing
type Parser struct {
	toks []Token
	pos  int
}

func NewParser(toks []Token) *Parser {
	return &Parser{toks: toks, pos: 0}
}

func (p *Parser) current() Token {
	if p.pos >= len(p.toks) {
		return Token{typ: TokenEOF}
	}
	return p.toks[p.pos]
}

func (p *Parser) next() Token {
	tok := p.current()
	if p.pos < len(p.toks) {
		p.pos++
	}
	return tok
}

// parsePrimary handles numbers, parenthesized expressions, and prefix unary operators
func (p *Parser) parsePrimary() Node {
	tok := p.current()
	switch tok.typ {
	case TokenNumber:
		p.next()
		val, err := strconv.ParseFloat(tok.val, 64)
		if err != nil {
			panic(fmt.Sprintf("cannot parse number %s", tok.val))
		}
		return &NumberNode{Value: val}
	case TokenLParen:
		p.next()
		expr := p.parseExpression(0)
		if p.current().typ != TokenRParen {
			panic("missing closing parenthesis")
		}
		p.next()
		return expr
	case TokenOperator:
		// Only unary minus is supported as a prefix operator
		if tok.val == "-" {
			p.next()
			right := p.parsePrimary()
			return &UnaryNode{Op: "-", Right: right}
		}
		panic(fmt.Sprintf("unexpected operator %s", tok.val))
	default:
		panic(fmt.Sprintf("unexpected token %v", tok))
	}
}

// parseExpression parses based on minimum precedence (Pratt algorithm)
func (p *Parser) parseExpression(minPrec int) Node {
	left := p.parsePrimary()
	for {
		opTok := p.current()
		if opTok.typ != TokenOperator {
			break
		}
		prec, ok := precedence[opTok.val]
		if !ok || prec < minPrec {
			break
		}
		// Determine next minimum precedence based on associativity
		nextMin := prec + 1
		if opTok.val == "^" { // right‑associative
			nextMin = prec
		}
		p.next()
		right := p.parseExpression(nextMin)
		left = &BinaryNode{Op: opTok.val, Left: left, Right: right}
	}
	return left
}

// Parse builds an AST from the token slice
func Parse(input string) Node {
	toks := Lex(input)
	parser := NewParser(toks)
	ast := parser.parseExpression(0)
	if parser.current().typ != TokenEOF {
		panic("unexpected tokens after complete parse")
	}
	return ast
}
