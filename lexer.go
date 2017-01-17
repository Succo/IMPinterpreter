package main

import (
	"fmt"
	"io"
)

type token int

const (
	// Operator on int
	Plus token = iota
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

type interpreter struct {
	sigma map[string]int
	prog  io.Reader
}

func newInterpreter(in io.Reader) *interpreter {
	sigma := make(map[string]int)
	return &interpreter{sigma: sigma, prog: in}
}

func main() {
	fmt.Println("vim-go")
}
