package main

import (
	"math"
	"regexp"

	"github.com/fatih/color"
	"github.com/rodcorsi/ecalc/calc"
)

var (
	fmtPrompt = color.New(color.FgCyan)
	fmtResult = color.New(color.FgYellow)
	fmtError  = color.New(color.FgRed)

	fmtFunction = color.New(color.FgYellow)
	//fmtNumber   = color.New(color.FgRed)

	re = regexp.MustCompile(`d|'|"|atan|acos|asin`)
)

const minEngNotation = float64(0.00000001)
const maxEngNotation = float64(9999999999999.0)

var lastAnswer float64

func init() {
	calc.AddConstant("ans", func() float64 {
		return lastAnswer
	})
}

type ECalc struct {
	Result *Result
}

func NewECalc() *ECalc {
	ecalc := &ECalc{}
	ecalc.Eval("0")
	return ecalc
}

func (e *ECalc) Eval(expr string) *Result {
	c := &Result{
		Expression: expr,
		Degree:     re.MatchString(expr),
	}
	e.Result = c

	stack, err := calc.ParseExpression(expr)
	if err != nil {
		c.Error = err
		return c
	}

	stack, c.Partial = addANS(stack)
	c.StackExpr = stack

	c.Value, c.Error = calc.SolveStack(stack)

	if c.Error == nil {
		lastAnswer = c.Value
	}

	if c.Value == 0 {
		c.EngNotation = false
		return c
	}
	v := math.Abs(c.Value)
	c.EngNotation = (v > maxEngNotation || v < minEngNotation)

	return c
}

func (e *ECalc) AddConstant(name string, value float64) {
	calc.AddConstant(name, func() float64 {
		return value
	})
}

func addANS(stack calc.Stack) (calc.Stack, bool) {
	if len(stack.Values) == 0 {
		return stack, false
	}

	added := false

	if stack.Values[0].Type == calc.OPERATOR {
		// first token is an operator add ans constant first
		stack.Values = append([]calc.Token{calc.Token{Type: calc.CONSTANT, Value: "ans"}}, stack.Values...)
		added = true
	} else if v := stack.Values[len(stack.Values)-1]; v.Type == calc.FUNCTION || v.Type == calc.OPERATOR {
		// last token is an operator or function add ans constant in the last position
		stack.Push(calc.Token{Type: calc.CONSTANT, Value: "ans"})
		added = true
	}

	return stack, added
}
