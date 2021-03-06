package main

import "fmt"

type instType int

const (
	whileT instType = iota
	ifT
	assignT
	skipT
	printT
)

// integer expression
type IExpr interface {
	evaluate(s *State) int
}

type PlusExpr struct {
	x1 IExpr
	x2 IExpr
}

func (p PlusExpr) evaluate(s *State) int {
	return p.x1.evaluate(s) + p.x2.evaluate(s)
}

type MinusExpr struct {
	x1 IExpr
	x2 IExpr
}

func (m MinusExpr) evaluate(s *State) int {
	return m.x1.evaluate(s) - m.x2.evaluate(s)
}

type TimeExpr struct {
	x1 IExpr
	x2 IExpr
}

func (t TimeExpr) evaluate(s *State) int {
	return t.x1.evaluate(s) * t.x2.evaluate(s)
}

type valExpr struct {
	name string
}

func (v valExpr) evaluate(s *State) int {
	return s.val[v.name]
}

type constExpr struct {
	val int
}

func (c constExpr) evaluate(s *State) int {
	return c.val
}

// boolean expression
type BExpr interface {
	evaluate(s *State) bool
}

type TrueExpr struct{}

func (t TrueExpr) evaluate(s *State) bool {
	return true
}

type FalseExpr struct{}

func (f FalseExpr) evaluate(s *State) bool {
	return false
}

type EqualExpr struct {
	x1 IExpr
	x2 IExpr
}

func (e EqualExpr) evaluate(s *State) bool {
	return e.x1.evaluate(s) == e.x2.evaluate(s)
}

type LessEqualExpr struct {
	x1 IExpr
	x2 IExpr
}

func (l LessEqualExpr) evaluate(s *State) bool {
	return l.x1.evaluate(s) <= l.x2.evaluate(s)
}

type NotExpr struct {
	b BExpr
}

func (n NotExpr) evaluate(s *State) bool {
	return !n.b.evaluate(s)
}

type OrExpr struct {
	b1 BExpr
	b2 BExpr
}

func (o OrExpr) evaluate(s *State) bool {
	return o.b1.evaluate(s) || o.b2.evaluate(s)
}

type AndExpr struct {
	b1 BExpr
	b2 BExpr
}

func (a AndExpr) evaluate(s *State) bool {
	return a.b1.evaluate(s) && a.b2.evaluate(s)
}

type State struct {
	val map[string]int
}

type Instruction interface {
	apply(*State) *State
	getType() instType
}

type SkipInst struct{}

func (skip SkipInst) apply(s *State) *State { return s }
func (skip SkipInst) getType() instType     { return skipT }

type PrintInst struct {
	val IExpr
}

func (p PrintInst) apply(s *State) *State {
	fmt.Println(p.val.evaluate(s))
	return s
}
func (p PrintInst) getType() instType { return printT }

type AssignInst struct {
	name string
	val  IExpr
}

func (a AssignInst) apply(s *State) *State {
	s.val[a.name] = a.val.evaluate(s)
	return s
}
func (a AssignInst) getType() instType { return assignT }

type WhileInst struct {
	cond BExpr
	loop []Instruction
}

func (w WhileInst) apply(s *State) *State {
	tempP := Interpreter{s: s, i: w.loop}
	for w.cond.evaluate(s) {
		s = tempP.execute()
	}
	return s
}
func (w WhileInst) getType() instType { return whileT }

type IfInst struct {
	cond  BExpr
	path1 []Instruction
	path2 []Instruction
}

func (i IfInst) apply(s *State) *State {
	var tempP *Interpreter
	if i.cond.evaluate(s) {
		tempP = &Interpreter{s: s, i: i.path1}
	} else {
		tempP = &Interpreter{s: s, i: i.path1}
	}
	return tempP.execute()
}
func (i IfInst) getType() instType { return ifT }

type Interpreter struct {
	s *State
	i []Instruction
}

func (i *Interpreter) execute() *State {
	for _, inst := range i.i {
		i.s = inst.apply(i.s)
	}
	return i.s
}

func NewInterpreter(insts []Instruction) *Interpreter {
	s := &State{make(map[string]int)}
	return &Interpreter{s: s, i: insts}
}
