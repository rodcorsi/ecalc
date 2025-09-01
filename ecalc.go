package main

import (
	"math/big"
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

var minEngNotation = big.NewFloat(0.00000001)
var maxEngNotation = big.NewFloat(9999999999999.0)

var lastAnswer = &Result{Value: new(big.Float)}

func init() {
	calc.AddConstant("ans", func() *big.Float {
		return lastAnswer.Value
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

	if c.Error != nil {
		return c
	}
	lastAnswer = c

	if c.Value.Cmp(big.NewFloat(0)) == 0 {
		c.EngNotation = false
		return c
	}
	v := new(big.Float).Abs(c.Value)
	c.EngNotation = (v.Cmp(maxEngNotation) > 0 || v.Cmp(minEngNotation) < 0)

	return c
}

func (e *ECalc) AddConstant(name string, value *big.Float) {
	calc.AddConstant(name, func() *big.Float {
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
