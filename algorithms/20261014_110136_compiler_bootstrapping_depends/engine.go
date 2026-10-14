package main

import (
	"fmt"
	"strconv"
	"unicode"
)

// ---------- Lexer ----------
type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(src string) *Lexer {
	return &Lexer{input: []rune(src), pos: 0}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	return ch
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.peek()) {
		l.next()
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	ch := l.peek()
	if ch == 0 {
		return Token{Typ: EOF}
	}
	// Single-char tokens
	switch ch {
	case '+':
		l.next()
		return Token{Typ: Plus, Lit: "+"}
	case '-':
		l.next()
		return Token{Typ: Minus, Lit: "-"}
	case '*':
		l.next()
		return Token{Typ: Mul, Lit: "*"}
	case '/':
		l.next()
		return Token{Typ: Div, Lit: "/"}
	case '=':
		l.next()
		return Token{Typ: Assign, Lit: "="}
	case '(':
		l.next()
		return Token{Typ: LParen, Lit: "("}
	case ')':
		l.next()
		return Token{Typ: RParen, Lit: ")"}
	case ';':
		l.next()
		return Token{Typ: Semicolon, Lit: ";"}
	}
	// Identifier or integer
	if unicode.IsLetter(ch) {
		start := l.pos
		for unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_' {
			l.next()
		}
		lit := string(l.input[start:l.pos])
		return Token{Typ: Ident, Lit: lit}
	}
	if unicode.IsDigit(ch) {
		start := l.pos
		for unicode.IsDigit(l.peek()) {
			l.next()
		}
		lit := string(l.input[start:l.pos])
		return Token{Typ: Int, Lit: lit}
	}
	panic(fmt.Sprintf("unexpected character: %q", ch))
}

// ---------- Parser ----------
type Parser struct {
	lex   *Lexer
	cur   Token
	peek  Token
	ready bool
}

func NewParser(src string) *Parser {
	lex := NewLexer(src)
	p := &Parser{lex: lex}
	p.next()
	p.next()
	return p
}

func (p *Parser) next() {
	if p.ready {
		p.cur = p.peek
		p.ready = false
	} else {
		p.cur = p.lex.NextToken()
	}
	p.peek = p.lex.NextToken()
}

func (p *Parser) expect(t TokenType) Token {
	if p.cur.Typ != t {
		panic(fmt.Sprintf("expected token %v, got %v", t, p.cur.Typ))
	}
	tok := p.cur
	p.next()
	return tok
}

func (p *Parser) Parse() []Stmt {
	var stmts []Stmt
	for p.cur.Typ != EOF {
		stmt := p.parseStmt()
		stmts = append(stmts, stmt)
		if p.cur.Typ == Semicolon {
			p.next()
		}
	}
	return stmts
}

func (p *Parser) parseStmt() Stmt {
	if p.cur.Typ == Ident && p.peek.Typ == Assign {
		// assignment
		name := p.cur.Lit
		p.next() // ident
		p.next() // '='
		val := p.parseExpr()
		return AssignStmt{Name: name, Value: val}
	}
	// expression statement
	expr := p.parseExpr()
	return ExprStmt{Expr: expr}
}

func (p *Parser) parseExpr() Expr {
	left := p.parseTerm()
	for p.cur.Typ == Plus || p.cur.Typ == Minus {
		op := p.cur.Typ
		p.next()
		right := p.parseTerm()
		left = BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseTerm() Expr {
	left := p.parseFactor()
	for p.cur.Typ == Mul || p.cur.Typ == Div {
		op := p.cur.Typ
		p.next()
		right := p.parseFactor()
		left = BinaryExpr{Op: op, Left: left, Right: right}
	}
	return left
}

func (p *Parser) parseFactor() Expr {
	switch p.cur.Typ {
	case Int:
		val, err := strconv.Atoi(p.cur.Lit)
		if err != nil {
			panic(err)
		}
		p.next()
		return IntLiteral{Value: val}
	case Ident:
		name := p.cur.Lit
		p.next()
		return IdentExpr{Name: name}
	case LParen:
		p.next()
		expr := p.parseExpr()
		p.expect(RParen)
		return expr
	default:
		panic(fmt.Sprintf("unexpected token in factor: %v", p.cur))
	}
}

// ---------- Interpreter ----------
type Interpreter struct {
	vars map[string]int
}

func NewInterpreter() *Interpreter {
	return &Interpreter{vars: make(map[string]int)}
}

func (i *Interpreter) Eval(stmts []Stmt) int {
	var last int
	for _, s := range stmts {
		switch stmt := s.(type) {
		case AssignStmt:
			val := i.evalExpr(stmt.Value)
			i.vars[stmt.Name] = val
			last = val
		case ExprStmt:
			last = i.evalExpr(stmt.Expr)
		default:
			panic("unknown statement type")
		}
	}
	return last
}

func (i *Interpreter) evalExpr(e Expr) int {
	switch expr := e.(type) {
	case IntLiteral:
		return expr.Value
	case IdentExpr:
		val, ok := i.vars[expr.Name]
		if !ok {
			panic(fmt.Sprintf("undefined variable %s", expr.Name))
		}
		return val
	case BinaryExpr:
		left := i.evalExpr(expr.Left)
		right := i.evalExpr(expr.Right)
		switch expr.Op {
		case Plus:
			return left + right
		case Minus:
			return left - right
		case Mul:
			return left * right
		case Div:
			if right == 0 {
				panic("division by zero")
			}
			return left / right
		default:
			panic("unknown binary operator")
		}
	default:
		panic("unknown expression type")
	}
}

// ---------- Compiler ----------
type Compiler struct {
	instructions []Instruction
	varIndices   map[string]int
	varNames     []string
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions: []Instruction{},
		varIndices:   make(map[string]int),
		varNames:     []string{},
	}
}

func (c *Compiler) Compile(stmts []Stmt) []Instruction {
	for _, s := range stmts {
		switch stmt := s.(type) {
		case AssignStmt:
			c.compileExpr(stmt.Value)
			idx := c.varIndex(stmt.Name)
			c.emit(OpStore, idx)
		case ExprStmt:
			c.compileExpr(stmt.Expr)
			// result stays on stack
		}
	}
	c.emit(OpHalt, 0)
	return c.instructions
}

func (c *Compiler) varIndex(name string) int {
	if idx, ok := c.varIndices[name]; ok {
		return idx
	}
	idx := len(c.varNames)
	c.varNames = append(c.varNames, name)
	c.varIndices[name] = idx
	return idx
}

func (c *Compiler) emit(op OpCode, arg int) {
	c.instructions = append(c.instructions, Instruction{Op: op, Arg: arg})
}

func (c *Compiler) compileExpr(e Expr) {
	switch expr := e.(type) {
	case IntLiteral:
		c.emit(OpConst, expr.Value)
	case IdentExpr:
		idx := c.varIndex(expr.Name)
		c.emit(OpLoad, idx)
	case BinaryExpr:
		c.compileExpr(expr.Left)
		c.compileExpr(expr.Right)
		switch expr.Op {
		case Plus:
			c.emit(OpAdd, 0)
		case Minus:
			c.emit(OpSub, 0)
		case Mul:
			c.emit(OpMul, 0)
		case Div:
			c.emit(OpDiv, 0)
		default:
			panic("unknown binary operator in compiler")
		}
	default:
		panic("unknown expression type in compiler")
	}
}

// ---------- Virtual Machine ----------
type VM struct {
	instructions []Instruction
	pc           int
	stack        []int
	vars         []int // indexed by var index
}

func NewVM(instr []Instruction, varCount int) *VM {
	return &VM{
		instructions: instr,
		pc:           0,
		stack:        []int{},
		vars:         make([]int, varCount),
	}
}

func (vm *VM) Run() int {
	for {
		if vm.pc >= len(vm.instructions) {
			panic("program counter out of range")
		}
		ins := vm.instructions[vm.pc]
		vm.pc++
		switch ins.Op {
		case OpConst:
			vm.stack = append(vm.stack, ins.Arg)
		case OpLoad:
			if ins.Arg >= len(vm.vars) {
				panic("load from undefined variable")
			}
			vm.stack = append(vm.stack, vm.vars[ins.Arg])
		case OpStore:
			if len(vm.stack) == 0 {
				panic("store with empty stack")
			}
			val := vm.stack[len(vm.stack)-1]
			vm.stack = vm.stack[:len(vm.stack)-1]
			if ins.Arg >= len(vm.vars) {
				panic("store to undefined variable")
			}
			vm.vars[ins.Arg] = val
		case OpAdd, OpSub, OpMul, OpDiv:
			if len(vm.stack) < 2 {
				panic("binary op with insufficient stack")
			}
			b := vm.stack[len(vm.stack)-1]
			a := vm.stack[len(vm.stack)-2]
			vm.stack = vm.stack[:len(vm.stack)-2]
			var res int
			switch ins.Op {
			case OpAdd:
				res = a + b
			case OpSub:
				res = a - b
			case OpMul:
				res = a * b
			case OpDiv:
				if b == 0 {
					panic("division by zero in VM")
				}
				res = a / b
			}
			vm.stack = append(vm.stack, res)
		case OpHalt:
			if len(vm.stack) == 0 {
				return 0
			}
			return vm.stack[len(vm.stack)-1]
		default:
			panic(fmt.Sprintf("unknown opcode %v", ins.Op))
		}
	}
}
