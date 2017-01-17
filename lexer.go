package main

import (
	"bufio"
	"bytes"
	"fmt"
	"unicode"
)

type tokenType int

const (
	// Operator on int
	Plus tokenType = iota
	Minus
	Time

	// Operator on bool
	True
	False
	Equal
	LessEqual
	Not
	Or
	And

	// Instructions keyword
	Skip
	Assign
	Semicolon
	While
	Do
	If

	// Identifier
	Int
	Variable

	Eof
)

var eof = rune(0)

type Token struct {
	ttype  tokenType
	lexeme string
	line   int
}

type Interpreter struct {
	sigma   map[string]int
	prog    bufio.Reader
	scanned []Token

	line int
}

func newInterpreter(in bufio.Reader) *Interpreter {
	sigma := make(map[string]int)
	scanned := make([]Token, 0)
	return &Interpreter{sigma: sigma, prog: in, scanned: scanned, line: 0}
}

func (i *Interpreter) read() rune {
	ch, _, err := i.prog.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (i *Interpreter) unread() {
	_ = i.prog.UnreadRune()
}

func (i *Interpreter) peek() rune {
	ch, _, err := i.prog.ReadRune()
	if err != nil {
		ch = eof
	}
	i.prog.UnreadRune()
	return ch
}

// addToken add a parsed token to the token list
func (i *Interpreter) addToken(ttype tokenType, lexeme string) {
	t := Token{ttype, lexeme, i.line}
	i.scanned = append(i.scanned, t)
}

// scanWord scan a word ans returns its value
func (i *Interpreter) scanWord(r rune) string {
	var buf bytes.Buffer
	buf.WriteRune(r)
	r = i.read()
	for isValidCharacter(r) {
		buf.WriteRune(r)
		r = i.read()
	}
	if !unicode.IsSpace(r) {
		panic(fmt.Sprintf("Unknown character %s on line %d", r, i.line))
	}
	i.unread()
	return buf.String()
}

// scanInt scan an int and return its value
func (i *Interpreter) scanInt(r rune) string {
	var buf bytes.Buffer
	buf.WriteRune(r)
	r = i.read()
	for unicode.IsDigit(r) {
		buf.WriteRune(r)
		r = i.read()
	}
	if !unicode.IsSpace(r) {
		panic(fmt.Sprintf("Unknown character %s on line %d", r, i.line))
	}
	i.unread()
	return buf.String()
}

// scan scans the next lexeme
func (i *Interpreter) scan() bool {
	r := i.read()
	switch {
	case unicode.IsSpace(r):
	case r == '\n':
		i.line++
	case r == '+':
		i.addToken(Plus, "+")
	case r == '-':
		i.addToken(Minus, "-")
	case r == '*':
		i.addToken(Time, "*")
	case r == '=':
		i.addToken(Equal, "=")
	case r == '!': // ¬ is too weird of a character
		i.addToken(Not, "!")
	case r == '|': // ∧ is too weird of a character
		i.read()
		i.addToken(Or, "||")
	case r == '&': // ∨ is too weird of a character
		i.read()
		i.addToken(Or, "&&")
	case r == '<':
		i.read() // No weird checking because < is reserved and only used here
		i.addToken(LessEqual, "<=")
	case r == ':':
		i.read() // No weird checking because : is reserved and only used here
		i.addToken(Assign, ":=")
	case isValidCharacter(r):
		w := i.scanWord(r)
		switch w {
		case "true":
			i.addToken(True, "true")
		case "false":
			i.addToken(False, "false")
		case "skip":
			i.addToken(Skip, "skip")
		case "while":
			i.addToken(While, "while")
		case "do":
			i.addToken(Do, "do")
		case "if":
			i.addToken(If, "if")
		default:
			i.addToken(Variable, w)
		}
	case unicode.IsNumber(r):
		w := i.scanInt(r)
		i.addToken(Int, w)
	case r == eof:
		i.addToken(Eof, "")
		return false
	default:
		panic(fmt.Sprintf("Unknown character %s on line %d", r, i.line))
	}
	return true
}

func (i *Interpreter) Scan() {
	for i.scan() {
	}
}

// isValidCharacter returns true is the rune is an acceptable character (only letter)
func isValidCharacter(r rune) bool {
	return ('a' <= r && 'z' >= r) ||
		('A' <= r && 'Z' >= r)
}

func main() {
	fmt.Println("vim-go")
}
