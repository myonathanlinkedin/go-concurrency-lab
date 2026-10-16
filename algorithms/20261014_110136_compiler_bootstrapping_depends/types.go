package main

import "fmt"

type TokenType int

const (
	EOF TokenType = iota
	Ident
	Int
	Plus
	Minus
	Mul
	Div
	Assign
	LParen
	RParen
	Semicolon
)

type Token struct {
	Typ TokenType
	Lit string
}

// AST node interfaces
type Stmt interface {
	isStmt()
}
type Expr interface {
	isExpr()
}

// Statements
type AssignStmt struct {
	Name  string
	Value Expr
}

func (AssignStmt) isStmt() {}

type ExprStmt struct {
	Expr Expr
}

func (ExprStmt) isStmt() {}

// Expressions
type BinaryExpr struct {
	Op    TokenType
	Left  Expr
	Right Expr
}

func (BinaryExpr) isExpr() {}

type IdentExpr struct {
	Name string
}

func (IdentExpr) isExpr() {}

type IntLiteral struct {
	Value int
}

func (IntLiteral) isExpr() {}

// Bytecode definitions
type OpCode int

const (
	OpConst OpCode = iota
	OpLoad
	OpStore
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpHalt
)

type Instruction struct {
	Op  OpCode
	Arg int // meaning depends on Op
}

// Helper for debugging
func (op OpCode) String() string {
	switch op {
	case OpConst:
		return "CONST"
	case OpLoad:
		return "LOAD"
	case OpStore:
		return "STORE"
	case OpAdd:
		return "ADD"
	case OpSub:
		return "SUB"
	case OpMul:
		return "MUL"
	case OpDiv:
		return "DIV"
	case OpHalt:
		return "HALT"
	default:
		return fmt.Sprintf("OP(%d)", int(op))
	}
}
