package main

import (
	"bufio"
	"fmt"
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

func (i *Interpreter) scan() {

}

func main() {
	fmt.Println("vim-go")
}
