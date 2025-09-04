package ecalc

import (
	"math/big"
	"regexp"

	"github.com/rodcorsi/ecalc/esolver"
)

var reDegree = regexp.MustCompile(`d|'|"|atan|acos|asin`)
var minEngNotation = big.NewFloat(0.00000001)
var maxEngNotation = big.NewFloat(9999999999999.0)

type ECalc struct {
	solver     esolver.ESolver
	Result     *Result
	LastAnswer *Result
}

func NewECalc() *ECalc {
	e := &ECalc{
		solver: esolver.New(),
	}
	e.solver.AddConstant("ans", func() *big.Float {
		return e.LastAnswer.Value
	})
	e.Eval("0")
	return e
}

func (e *ECalc) Eval(expr string) *Result {
	c := &Result{
		Expression: expr,
		Degree:     reDegree.MatchString(expr),
	}
	e.Result = c

	stack, err := e.solver.ParseExpression(expr)
	if err != nil {
		c.Error = err
		return c
	}

	stack, c.Partial = addANS(stack)
	c.StackExpr = stack

	c.Value, c.Error = e.solver.SolveStack(stack)

	if c.Error != nil {
		return c
	}
	e.LastAnswer = c

	if c.Value.Cmp(big.NewFloat(0)) == 0 {
		c.EngNotation = false
		return c
	}
	v := new(big.Float).Abs(c.Value)
	c.EngNotation = (v.Cmp(maxEngNotation) > 0 || v.Cmp(minEngNotation) < 0)

	return c
}

func (e *ECalc) AddConstant(name string, value *big.Float) {
	e.solver.AddConstant(name, func() *big.Float {
		return value
	})
}

func addANS(stack esolver.Stack) (esolver.Stack, bool) {
	if len(stack.Values) == 0 {
		return stack, false
	}

	added := false

	if stack.Values[0].Type == esolver.OPERATOR {
		// first token is an operator add ans constant first
		stack.Values = append([]esolver.Token{esolver.Token{Type: esolver.CONSTANT, Value: "ans"}}, stack.Values...)
		added = true
	} else if v := stack.Values[len(stack.Values)-1]; v.Type == esolver.FUNCTION || v.Type == esolver.OPERATOR {
		// last token is an operator or function add ans constant in the last position
		stack.Push(esolver.Token{Type: esolver.CONSTANT, Value: "ans"})
		added = true
	}

	return stack, added
}
