package main

import (
	"fmt"
	"io"
	"math"
	"regexp"

	"github.com/fatih/color"
	"github.com/rodcorsi/ecalc/calc"
)

type ECalc struct {
	Value       float64
	Degree      bool
	Error       error
	Writer      io.Writer
	EngNotation bool
	Partial     bool
	Expression  string
	StackExpr   calc.Stack
}

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

var Ans float64

func init() {
	calc.AddConstant("ans", func() float64 {
		return Ans
	})
}

func NewECalc(writer io.Writer) *ECalc {
	return &ECalc{
		Writer: writer,
	}
}

func (e *ECalc) PrintResult() {
	//Print Stack
	lastType := calc.TokenType(-1)

	for _, v := range e.StackExpr.Values {
		closeParen := false
		if lastType == calc.FUNCTION && (v.Type == calc.NUMBER || v.Type == calc.CONSTANT) {
			e.Print("(")
			closeParen = true
		}

		if v.Type == calc.FUNCTION || v.Type == calc.CONSTANT {
			fmtFunction.Print(v.Value)
		} else if v.Type == calc.NUMBER {
			fmt.Print(v.Value)
		} else if v.Value == "+" || v.Value == "-" {
			e.Print(" " + v.Value + " ")
		} else {
			e.Print(v.Value)
		}
		if closeParen {
			e.Print(")")
		}
		lastType = v.Type
	}

	e.Print(" = ")

	if e.Error != nil {
		//fmtError.Fprintln(e.Writer, "Error:", e.Error.Error())
		fmtError.Println("Error:", e.Error.Error())
	} else if e.Degree {
		//fmtResult.Fprintln(e.Writer, convertDMS(e.Value))
		fmtResult.Println(convertDMS(e.Value))
	} else if e.EngNotation {
		//fmtResult.Fprintf(e.Writer, "%e\n", e.Value)
		fmtResult.Printf("%e\n", e.Value)
	} else {
		//fmtResult.Fprintf(e.Writer, "%.12f\n", e.Value)
		fmtResult.Printf("%.12f\n", e.Value)
	}
}

func (e *ECalc) PrintPrompt() {
	e.Print("(ans:")
	if e.EngNotation {
		//fmtPrompt.Fprintf(e.Writer, "%e", e.Value)
		fmtPrompt.Printf("%e", e.Value)
	} else {
		//fmtPrompt.Fprintf(e.Writer, "%14.8f", e.Value)
		fmtPrompt.Printf("%14.8f", e.Value)
	}

	e.Print(") » ")
}

func (e *ECalc) Print(a ...interface{}) {
	fmt.Fprint(e.Writer, a...)
}

func (e *ECalc) Println(a ...interface{}) {
	fmt.Fprintln(e.Writer, a...)
}

func (e *ECalc) Printf(format string, a ...interface{}) {
	fmt.Fprintf(e.Writer, format, a...)
}

func convertDMS(value float64) string {
	d := int64(value)
	m := (value - float64(d)) * 60.0
	s := (m - float64(int64(m))) * 60.0
	return fmt.Sprintf(`%vd%2d'%.5f"`, d, int64(m), s)
}

func (e *ECalc) Eval(expr string) {
	e.Expression = expr
	e.Degree = re.MatchString(expr)

	stack, err := calc.ParseExpression(expr)
	if err != nil {
		e.Error = err
		return
	}

	stack, e.Partial = addANS(stack)
	e.StackExpr = stack

	e.Value, e.Error = calc.SolveStack(stack)

	if e.Error == nil {
		Ans = e.Value
	}

	if e.Value == 0 {
		e.EngNotation = false
		return
	}
	v := math.Abs(e.Value)
	e.EngNotation = (v > maxEngNotation || v < minEngNotation)
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
		stack.Values = append([]calc.Token{calc.Token{calc.CONSTANT, "ans"}}, stack.Values...)
		added = true
	} else if v := stack.Values[len(stack.Values)-1]; v.Type == calc.FUNCTION || v.Type == calc.OPERATOR {
		// last token is an operator or function add ans constant in the last position
		stack.Push(calc.Token{calc.CONSTANT, "ans"})
		added = true
	}

	return stack, added
}
