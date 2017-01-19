package main

import (
	"fmt"
	"strconv"
)

// stack makes shunting yard easier to write
type stack []Token

func (s stack) push(t Token) stack {
	return append(s, t)
}

func (s stack) pop() (stack, Token) {
	l := len(s)
	return s[:l-1], s[l-1]
}

// precedence store precedence value for the precede function
var precedence = map[tokenType]int{
	Time:  0,
	Plus:  1,
	Minus: 1,
	And:   0,
	Or:    0,
	Not:   1,
}

func precede(op1 tokenType, op2 tokenType) bool {
	val1, ok := precedence[op1]
	if !ok {
		panic("Unexpected operator while parsing IExpr")
	}
	val2, ok := precedence[op1]
	if !ok {
		panic("Unexpected operator while parsing IExpr")
	}
	return val1 <= val2
}

// parser is the struct used to generate instructions
type parser struct {
	tokens []Token
	inst   []Instruction
}

func (p *parser) pop() Token {
	t := p.tokens[0]
	p.tokens = p.tokens[1:]
	return t
}

func (p *parser) peek() Token {
	return p.tokens[0]
}

func addIOperator(out []IExpr, op Token) []IExpr {
	l := len(out)
	if l < 2 { // that would mean two operator in a row
		panic("Unexpected token while add Ioperator")
	}

	// Apply the operator to the two IExpr in out
	x1 := out[l-1]
	x2 := out[l-2]
	out = out[:l-2]
	switch op.ttype {
	case Plus:
		out = append(out, PlusExpr{x2, x1})
	case Minus:
		out = append(out, MinusExpr{x2, x1})
	case Time:
		out = append(out, TimeExpr{x2, x1})
	}
	return out
}

// parseIExpr parses an expression that evaluate to int
func (p *parser) parseIExpr() IExpr {
	out := make([]IExpr, 0)
	operators := stack(make([]Token, 0))
tokenLoop:
	for {
		t := p.tokens[0]
		switch t.ttype {
		case Int:
			val, err := strconv.Atoi(t.lexeme)
			if err != nil {
				panic(fmt.Sprintf("Incorrect value line %d", t.line))
			}
			out = append(out, constExpr{val})
		case Variable:
			out = append(out, valExpr{t.lexeme})
		case Plus, Minus, Time:
			for len(operators) > 0 {
				operators, op := operators.pop()
				if precede(op.ttype, t.ttype) {
					out = addIOperator(out, op)
				} else {
					operators = operators.push(op)
					break
				}
			}
			operators = operators.push(t)
		case LeftParen:
			p.pop()
			expr := p.parseIExpr()
			out = append(out, expr)
			if p.tokens[0].ttype != RightParen {
				panic(fmt.Sprintf("Unmatched parentheses line %d", t.line))
			}
		default:
			break tokenLoop
		}
		p.pop()
	}
	for _, op := range operators {
		out = addIOperator(out, op)
	}
	if len(out) != 1 {
		panic("Can't parse IExpr")
	}
	return out[0]
}

func addBOperator(out []BExpr, op Token) []BExpr {
	// vase wher it's an unary operator
	if op.ttype == Not {
		l := len(out)
		if l < 1 { // that would mean two operator in a row
			panic("Unexpected token while add Ioperator")
		}
		b := out[l-1]
		out = out[:l-1]
		out = append(out, NotExpr{b})
		return out
	}

	l := len(out)
	if l < 2 { // that would mean two operator in a row
		panic("Unexpected token while adding Boperator")
	}
	// Apply the operator to the two BExpr in out
	b1 := out[l-1]
	b2 := out[l-2]
	out = out[:l-2]
	switch op.ttype {
	case Or:
		out = append(out, OrExpr{b1, b2})
	case And:
		out = append(out, AndExpr{b1, b2})
	}
	return out
}

// parseBExpr parses an expression that evaluate to bool
func (p *parser) parseBExpr() BExpr {
	out := make([]BExpr, 0)
	operators := stack(make([]Token, 0))
tokenLoop:
	for {
		t := p.tokens[0]
		switch t.ttype {
		case True:
			out = append(out, TrueExpr{})
			p.pop()
		case False:
			out = append(out, FalseExpr{})
			p.pop()
		case Or, And, Not:
			for len(operators) > 0 {
				operators, op := operators.pop()
				fmt.Println(op.lexeme)
				if precede(op.ttype, t.ttype) {
					out = addBOperator(out, op)
				} else {
					operators = operators.push(op)
					break
				}
			}
			operators = operators.push(t)
			p.pop()
		case Variable, Int:
			expr1 := p.parseIExpr()
			op := p.pop()
			if op.ttype != Equal && op.ttype != LessEqual {
				panic(fmt.Sprintf("Unexpected token %s line %d", op.lexeme, op.line))
			}
			expr2 := p.parseIExpr()
			switch op.ttype {
			case Equal:
				out = append(out, EqualExpr{expr1, expr2})
			case LessEqual:
				out = append(out, LessEqualExpr{expr1, expr2})
			}
		case LeftParen:
			p.pop()
			expr := p.parseBExpr()
			out = append(out, expr)
			if p.tokens[0].ttype != RightParen {
				panic(fmt.Sprintf("Unmatched parentheses line %d", t.line))
			}
			p.pop()
		default:
			break tokenLoop
		}
	}

	for _, op := range operators {
		out = addBOperator(out, op)
	}
	if len(out) != 1 {
		panic("Can't parse BExpr")
	}
	return out[0]
}

// parseInst parses a logical block of instruction
func (p *parser) parseInst() []Instruction {
	insts := make([]Instruction, 0)
	insts = append(insts, p.parse())
	for p.peek().ttype == Semicolon {
		p.pop()
		if p.peek().ttype == Eof {
			break
		}
		insts = append(insts, p.parse())
	}
	return insts
}

// Parse will parse the next instruction
func (p *parser) parse() Instruction {
	t := p.pop()
	switch {
	case t.ttype == Skip:
		return SkipInst{}
	case t.ttype == Variable:
		eq := p.pop()
		if eq.ttype != Assign {
			panic(fmt.Sprintf("Unexpected token %s line %d", eq.lexeme, eq.line))
		}
		expr := p.parseIExpr()
		return AssignInst{name: t.lexeme, val: expr}
	case t.ttype == Print:
		expr := p.parseIExpr()
		return PrintInst{expr}
	case t.ttype == While:
		expr := p.parseBExpr()
		do := p.pop()
		if do.ttype != Do {
			panic(fmt.Sprintf("Unexpected token %s line %d", do.lexeme, do.line))
		}
		loop := p.parseInst()
		od := p.pop()
		if od.ttype != Od {
			panic(fmt.Sprintf("Unexpected token %s line %d", od.lexeme, od.line))
		}
		return WhileInst{cond: expr, loop: loop}
	case t.ttype == If:
		expr := p.parseBExpr()
		then := p.pop()
		if then.ttype != Then {
			panic(fmt.Sprintf("Unexpected token %s line %d", then.lexeme, then.line))
		}
		path1 := p.parseInst()
		if_ := p.pop()
		if if_.ttype != Else {
			panic(fmt.Sprintf("Unexpected token %s line %d", if_.lexeme, if_.line))
		}
		path2 := p.parseInst()
		fi := p.pop()
		if fi.ttype != Fi {
			panic(fmt.Sprintf("Unexpected token %s line %d", fi.lexeme, fi.line))
		}
		return IfInst{cond: expr, path1: path1, path2: path2}
	default:
		panic(fmt.Sprintf("Unexpected token %s line %d", t.lexeme, t.line))
	}
}
