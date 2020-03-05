package main

import (
	"fmt"
	"io"
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

type Result struct {
	Value       float64
	Degree      bool
	Error       error
	Writer      io.Writer
	EngNotation bool
	Partial     bool
	Expression  string
	StackExpr   calc.Stack
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

func (c *Result) FormatValue() string {
	if c.EngNotation {
		return fmtPrompt.Sprintf("%e", c.Value)
	}
	return fmtPrompt.Sprintf("%14.8f", c.Value)
}

func (e *Result) FormatExpression() string {
	result := ""
	lastType := calc.TokenType(-1)

	for _, v := range e.StackExpr.Values {
		closeParen := false
		if lastType == calc.FUNCTION && (v.Type == calc.NUMBER || v.Type == calc.CONSTANT) {
			result += "("
			closeParen = true
		}

		if v.Type == calc.FUNCTION || v.Type == calc.CONSTANT {
			fmtFunction.Print(v.Value)
		} else if v.Type == calc.NUMBER {
			result += v.Value
		} else if v.Value == "+" || v.Value == "-" {
			result += " " + v.Value + " "
		} else {
			result += v.Value
		}
		if closeParen {
			result += ")"
		}
		lastType = v.Type
	}
	return result
}

func (c *Result) FormatResult() string {
	if c.Error != nil {
		return fmtError.Sprint("Error:", c.Error.Error())
	} else if c.Degree {
		return fmtResult.Sprintln(convertDMS(c.Value))
	} else if c.EngNotation {
		return fmtResult.Sprintf("%e\n", c.Value)
	}
	return fmtResult.Sprintf("%.12f\n", c.Value)
}

func convertDMS(value float64) string {
	d := int64(value)
	m := (value - float64(d)) * 60.0
	s := (m - float64(int64(m))) * 60.0
	return fmt.Sprintf(`%vd%2d'%.5f"`, d, int64(m), s)
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
