package main

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
