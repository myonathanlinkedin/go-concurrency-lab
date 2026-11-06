package main

import (
	"errors"
	"strconv"
	"unicode"
)

// state represents the current state of the tokenizer FSM.
type state int

const (
	stateStart state = iota
	stateIdent
	stateNumber
	stateOperator
	stateWhitespace
)

// Tokenize converts an input string into a slice of Tokens using a deterministic FSM.
func Tokenize(input string) ([]Token, error) {
	var tokens []Token
	runes := []rune(input)
	start := 0
	currState := stateStart

	emit := func(t TokenType, end int) {
		if start < end {
			tok := Token{
				Type:   t,
				Lexeme: string(runes[start:end]),
				Pos:    start,
			}
			tokens = append(tokens, tok)
		}
		start = end
	}

	for i, r := range runes {
		switch currState {
		case stateStart:
			switch {
			case unicode.IsLetter(r) || r == '_':
				currState = stateIdent
			case unicode.IsDigit(r):
				currState = stateNumber
			case unicode.IsSpace(r):
				currState = stateWhitespace
			case isOperatorRune(r):
				currState = stateOperator
			default:
				return nil, errors.New("invalid character at position " + strconv.Itoa(i))
			}
			start = i

		case stateIdent:
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				// stay in identifier
			} else {
				emit(TokenIdentifier, i)
				currState = stateStart
				i-- // reprocess this rune in the new state
			}

		case stateNumber:
			if unicode.IsDigit(r) {
				// stay in number
			} else {
				emit(TokenNumber, i)
				currState = stateStart
				i--
			}

		case stateOperator:
			// Operators are single-character for this simple lexer.
			emit(TokenOperator, i+1)
			currState = stateStart
			i-- // reprocess next rune in start state

		case stateWhitespace:
			if unicode.IsSpace(r) {
				// stay in whitespace
			} else {
				emit(TokenWhitespace, i)
				currState = stateStart
				i--
			}
		}
	}

	// Emit any trailing token.
	switch currState {
	case stateIdent:
		emit(TokenIdentifier, len(runes))
	case stateNumber:
		emit(TokenNumber, len(runes))
	case stateOperator:
		emit(TokenOperator, len(runes))
	case stateWhitespace:
		emit(TokenWhitespace, len(runes))
	}

	// Append EOF token.
	tokens = append(tokens, Token{Type: TokenEOF, Lexeme: "", Pos: len(runes)})
	return tokens, nil
}

// isOperatorRune reports whether r is a recognized operator character.
func isOperatorRune(r rune) bool {
	switch r {
	case '+', '-', '*', '/', '=', ';', ',', '(', ')':
		return true
	default:
		return false
	}
}

// Parse consumes a slice of Tokens and produces a slice of assignment statements.
// Grammar (very limited):
//   stmt   := identifier "=" number ";"
//   program:= { stmt }
func Parse(tokens []Token) ([]*AssignStmt, error) {
	var stmts []*AssignStmt
	pos := 0
	expect := func(t TokenType) (Token, error) {
		if pos >= len(tokens) {
			return Token{}, errors.New("unexpected end of input")
		}
		tok := tokens[pos]
		if tok.Type != t {
			return Token{}, errors.New("unexpected token " + tok.String())
		}
		pos++
		return tok, nil
	}

	for {
		// Stop at EOF
		if pos < len(tokens) && tokens[pos].Type == TokenEOF {
			break
		}
		// identifier
		idTok, err := expect(TokenIdentifier)
		if err != nil {
			return nil, err
		}
		// optional whitespace
		if pos < len(tokens) && tokens[pos].Type == TokenWhitespace {
			pos++
		}
		// '=' operator
		eqTok, err := expect(TokenOperator)
		if err != nil || eqTok.Lexeme != "=" {
			return nil, errors.New("expected '=' operator")
		}
		// optional whitespace
		if pos < len(tokens) && tokens[pos].Type == TokenWhitespace {
			pos++
		}
		// number
		numTok, err := expect(TokenNumber)
		if err != nil {
			return nil, err
		}
		// optional whitespace
		if pos < len(tokens) && tokens[pos].Type == TokenWhitespace {
			pos++
		}
		// ';' operator
		semiTok, err := expect(TokenOperator)
		if err != nil || semiTok.Lexeme != ";" {
			return nil, errors.New("expected ';' terminator")
		}
		// optional whitespace after statement
		if pos < len(tokens) && tokens[pos].Type == TokenWhitespace {
			pos++
		}
		// Build AST node
		val, err := strconv.Atoi(numTok.Lexeme)
		if err != nil {
			return nil, errors.New("invalid integer literal")
		}
		stmt := &AssignStmt{
			Var:   idTok.Lexeme,
			Value: val,
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}
