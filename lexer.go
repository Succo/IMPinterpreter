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

	// Separator
	LeftParen
	RightParen

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
	Od
	If
	Fi
	Then
	Else
	Print

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

type Lexer struct {
	sigma   map[string]int
	prog    *bufio.Reader
	scanned []Token

	line int
}

func NewLexer(in *bufio.Reader) *Lexer {
	sigma := make(map[string]int)
	scanned := make([]Token, 0)
	return &Lexer{sigma: sigma, prog: in, scanned: scanned, line: 0}
}

func (l *Lexer) read() rune {
	ch, _, err := l.prog.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (l *Lexer) unread() {
	_ = l.prog.UnreadRune()
}

func (l *Lexer) peek() rune {
	ch, _, err := l.prog.ReadRune()
	if err != nil {
		ch = eof
	}
	l.prog.UnreadRune()
	return ch
}

// addToken add a parsed token to the token list
func (l *Lexer) addToken(ttype tokenType, lexeme string) {
	t := Token{ttype, lexeme, l.line}
	l.scanned = append(l.scanned, t)
}

// scanWord scan a word ans returns its value
func (l *Lexer) scanWord(r rune) string {
	var buf bytes.Buffer
	buf.WriteRune(r)
	r = l.read()
	for isValidCharacter(r) {
		buf.WriteRune(r)
		r = l.read()
	}
	l.unread()
	return buf.String()
}

// scanInt scan an int and return its value
func (l *Lexer) scanInt(r rune) string {
	var buf bytes.Buffer
	buf.WriteRune(r)
	r = l.read()
	for unicode.IsDigit(r) {
		buf.WriteRune(r)
		r = l.read()
	}
	l.unread()
	return buf.String()
}

// scan scans the next lexeme
func (l *Lexer) scan() bool {
	r := l.read()
	switch {
	case r == '\n':
		l.line++
	case unicode.IsSpace(r):
	case r == '+':
		l.addToken(Plus, "+")
	case r == '-':
		l.addToken(Minus, "-")
	case r == '*':
		l.addToken(Time, "*")
	case r == '=':
		l.addToken(Equal, "=")
	case r == '(':
		l.addToken(LeftParen, "(")
	case r == ')':
		l.addToken(RightParen, ")")
	case r == ';':
		l.addToken(Semicolon, ";")
	case r == '!': // ¬ is too weird of a character
		l.addToken(Not, "!")
	case r == '|': // ∧ is too weird of a character
		l.read()
		l.addToken(Or, "||")
	case r == '&': // ∨ is too weird of a character
		l.read()
		l.addToken(Or, "&&")
	case r == '<':
		l.read() // No weird checking because < is reserved and only used here
		l.addToken(LessEqual, "<=")
	case r == ':':
		l.read() // No weird checking because : is reserved and only used here
		l.addToken(Assign, ":=")
	case isValidCharacter(r):
		w := l.scanWord(r)
		switch w {
		case "true":
			l.addToken(True, "true")
		case "false":
			l.addToken(False, "false")
		case "skip":
			l.addToken(Skip, "skip")
		case "while":
			l.addToken(While, "while")
		case "do":
			l.addToken(Do, "do")
		case "od":
			l.addToken(Od, "od")
		case "if":
			l.addToken(If, "if")
		case "fi":
			l.addToken(Fi, "fi")
		case "then":
			l.addToken(Then, "then")
		case "else":
			l.addToken(Else, "else")
		case "print":
			l.addToken(Print, "print")
		default:
			l.addToken(Variable, w)
		}
	case unicode.IsNumber(r):
		w := l.scanInt(r)
		l.addToken(Int, w)
	case r == eof:
		l.addToken(Eof, "")
		return false
	default:
		panic(fmt.Sprintf("Unknown character %s on line %d", string(r), l.line))
	}
	return true
}

// Scan the input stream and build the list of scanned token
func (l *Lexer) Scan() parser {
	for l.scan() {
	}
	return parser{l.scanned, make([]Instruction, 0)}
}

// isValidCharacter returns true is the rune is an acceptable character (only letter)
func isValidCharacter(r rune) bool {
	return ('a' <= r && 'z' >= r) ||
		('A' <= r && 'Z' >= r)
}
