package calc

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var errInvalidExpression = errors.New("Invalid Expression")

var oprData = map[string]struct {
	prec  int
	rAsoc bool // true = right // false = left
	fx    func(x, y float64) float64
}{
	"^": {4, true, func(x, y float64) float64 { return math.Pow(x, y) }},
	"*": {3, false, func(x, y float64) float64 { return x * y }},
	"/": {3, false, func(x, y float64) float64 { return x / y }},
	"+": {2, false, func(x, y float64) float64 { return x + y }},
	"-": {2, false, func(x, y float64) float64 { return x - y }},
}

var funcs = map[string]Function{
	"ln":    math.Log,
	"abs":   math.Abs,
	"cos":   cos,
	"sin":   sin,
	"tan":   tan,
	"acos":  acos,
	"asin":  asin,
	"atan":  atan,
	"sqrt":  math.Sqrt,
	"cbrt":  math.Cbrt,
	"ceil":  math.Ceil,
	"floor": math.Floor,
}

var consts = map[string]ConstFunction{
	"e":       func() float64 { return math.E },
	"pi":      func() float64 { return math.Pi },
	"phi":     func() float64 { return math.Phi },
	"sqrt2":   func() float64 { return math.Sqrt2 },
	"sqrte":   func() float64 { return math.SqrtE },
	"sqrtpi":  func() float64 { return math.SqrtPi },
	"sqrtphi": func() float64 { return math.SqrtPhi },
	"pol":     func() float64 { return 25.4 },
}

var elemNames map[string]TokenType

func init() {
	elemNames = make(map[string]TokenType)
	for k := range funcs {
		elemNames[k] = FUNCTION
	}

	for k := range consts {
		elemNames[k] = CONSTANT
	}
}

// SolvePostfix evaluates and returns the answer of the expression converted to postfix
func SolvePostfix(tokens Stack) (float64, error) {
	stack := Stack{}
	funcStack := Stack{}

	for _, v := range tokens.Values {
		switch v.Type {
		case NUMBER:
			if funcStack.IsEmpty() {
				stack.Push(v)
			} else {
				x, _ := strconv.ParseFloat(v.Value, 64)
				for !funcStack.IsEmpty() {
					f := funcs[funcStack.Pop().Value]
					x = f(x)
				}
				stack.Push(Token{NUMBER, strconv.FormatFloat(x, 'f', -1, 64)})
			}
		case FUNCTION:
			funcStack.Push(v)
		case CONSTANT:
			c, ok := consts[v.Value]
			if !ok {
				break
			}
			x := c()

			if funcStack.IsEmpty() {
				stack.Push(Token{NUMBER, strconv.FormatFloat(x, 'f', -1, 64)})
			} else {
				for !funcStack.IsEmpty() {
					f := funcs[funcStack.Pop().Value]
					x = f(x)
				}
				stack.Push(Token{NUMBER, strconv.FormatFloat(x, 'f', -1, 64)})
			}

		case OPERATOR:
			f := oprData[v.Value].fx
			var x, y float64
			y, _ = strconv.ParseFloat(stack.Pop().Value, 64)
			x, _ = strconv.ParseFloat(stack.Pop().Value, 64)
			result := f(x, y)
			stack.Push(Token{NUMBER, strconv.FormatFloat(result, 'f', -1, 64)})
		}
	}
	if len(stack.Values) != 1 {
		return -1, errInvalidExpression
	}
	return strconv.ParseFloat(stack.Values[0].Value, 64)
}

func addMissingOperator(stack Stack) Stack {
	if len(stack.Values) == 0 {
		return stack
	}

	fixed := Stack{}
	lastToken := TokenType(-1)

	for _, v := range stack.Values {
		if (lastToken == NUMBER || lastToken == RPAREN || lastToken == CONSTANT) &&
			(v.Type == NUMBER || v.Type == LPAREN || v.Type == CONSTANT || v.Type == FUNCTION) {
			fixed.Push(Token{OPERATOR, "*"})
		}

		lastToken = v.Type
		fixed.Push(v)
	}
	return fixed
}

// ContainsLetter checks if a string contains a letter
func ContainsLetter(s string) bool {
	for _, v := range s {
		if unicode.IsLetter(v) {
			return true
		}
	}
	return false
}

func Solve(s string) (float64, error) {
	stack, err := ParseExpression(s)
	if err != nil {
		return -1, err
	}

	return SolveStack(stack)
}

func ParseExpression(s string) (Stack, error) {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.Replace(s, ",", ".", -1)

	p := NewParser(strings.NewReader(s))

	stack, err := p.Parse()
	if err != nil {
		return Stack{}, err
	}
	stack = addMissingOperator(stack)

	return stack, nil
}

func SolveStack(stack Stack) (float64, error) {
	stack = ShuntingYard(stack)

	result, err := SolvePostfix(stack)

	if err != nil {
		return -1, err
	}

	return result, nil
}

func AddConstant(name string, constCreator ConstFunction) {
	consts[name] = constCreator
	elemNames[name] = CONSTANT
}
