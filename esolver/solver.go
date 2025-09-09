package esolver

import (
	"errors"
	"math"
	"math/big"
	"strings"
)

var errInvalidExpression = errors.New("invalid expression")

var oprData = map[string]struct {
	prec  int
	rAsoc bool // true = right // false = left
	fx    func(x, y *big.Float) *big.Float
}{
	"^": {4, true, func(x, y *big.Float) *big.Float { return bigPow(x, y) }},
	"*": {3, false, func(x, y *big.Float) *big.Float { z := new(big.Float); return z.Mul(x, y) }},
	"/": {3, false, func(x, y *big.Float) *big.Float { z := new(big.Float); return z.Quo(x, y) }},
	"+": {2, false, func(x, y *big.Float) *big.Float { z := new(big.Float); return z.Add(x, y) }},
	"-": {2, false, func(x, y *big.Float) *big.Float { z := new(big.Float); return z.Sub(x, y) }},
}

var funcs = map[string]Function{
	"ln":    bigLog,
	"abs":   bigAbs,
	"cos":   cos,
	"sin":   sin,
	"tan":   tan,
	"acos":  acos,
	"asin":  asin,
	"atan":  atan,
	"sqrt":  bigSqrt,
	"cbrt":  bigCbrt,
	"ceil":  bigCeil,
	"floor": bigFloor,
}

var consts = map[string]ConstFunction{
	"e":       func() *big.Float { return big.NewFloat(math.E) },
	"pi":      func() *big.Float { return big.NewFloat(math.Pi) },
	"phi":     func() *big.Float { return big.NewFloat(math.Phi) },
	"sqrt2":   func() *big.Float { return big.NewFloat(math.Sqrt2) },
	"sqrte":   func() *big.Float { return big.NewFloat(math.SqrtE) },
	"sqrtpi":  func() *big.Float { return big.NewFloat(math.SqrtPi) },
	"sqrtphi": func() *big.Float { return big.NewFloat(math.SqrtPhi) },
	"in":      func() *big.Float { return big.NewFloat(25.4) },
}

type ESolver interface {
	Solve(s string) (*big.Float, error)
	SolveStack(stack Stack) (*big.Float, error)
	SolvePostfix(tokens Stack) (*big.Float, error)
	ParseExpression(s string) (Stack, error)
	AddConstant(name string, constCreator ConstFunction)
}
type esolver struct {
	elemNames  map[string]TokenType
	userConsts map[string]ConstFunction
}

func New() ESolver {
	elemNames := make(map[string]TokenType)
	for k := range funcs {
		elemNames[k] = FUNCTION
	}

	for k := range consts {
		elemNames[k] = CONSTANT
	}
	return &esolver{
		elemNames:  elemNames,
		userConsts: make(map[string]ConstFunction),
	}
}

func (e *esolver) Solve(s string) (*big.Float, error) {
	stack, err := e.ParseExpression(s)
	if err != nil {
		return nil, err
	}

	return e.SolveStack(stack)
}

func (e *esolver) SolveStack(stack Stack) (*big.Float, error) {
	stack = ShuntingYard(stack)

	result, err := e.SolvePostfix(stack)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// SolvePostfix evaluates and returns the answer of the expression converted to postfix
func (e *esolver) SolvePostfix(tokens Stack) (*big.Float, error) {
	stack := Stack{}
	funcStack := Stack{}

	for _, v := range tokens.Values {
		switch v.Type {
		case NUMBER:
			if funcStack.IsEmpty() {
				stack.Push(v)
			} else {
				x, _, err := big.ParseFloat(v.Value, 10, 256, big.ToNearestEven)
				if err != nil {
					return nil, err
				}
				for !funcStack.IsEmpty() {
					f := funcs[funcStack.Pop().Value]
					x = f(x)
				}
				stack.Push(Token{NUMBER, x.Text('f', -1)})
			}
		case FUNCTION:
			funcStack.Push(v)
		case CONSTANT:
			c, ok := e.findConst(v.Value)
			if !ok {
				break
			}
			x := c()

			if funcStack.IsEmpty() {
				stack.Push(Token{NUMBER, x.Text('f', -1)})
			} else {
				for !funcStack.IsEmpty() {
					f := funcs[funcStack.Pop().Value]
					x = f(x)
				}
				stack.Push(Token{NUMBER, x.Text('f', -1)})
			}

		case OPERATOR:
			f := oprData[v.Value].fx

			yVal, _, err := big.ParseFloat(stack.Pop().Value, 10, 256, big.ToNearestEven)
			if err != nil {
				return nil, err
			}

			xVal, _, err := big.ParseFloat(stack.Pop().Value, 10, 256, big.ToNearestEven)
			if err != nil {
				return nil, err
			}

			result := f(xVal, yVal)
			stack.Push(Token{NUMBER, result.Text('f', -1)})
		}
	}
	if len(stack.Values) != 1 {
		return nil, errInvalidExpression
	}
	result, _, err := big.ParseFloat(stack.Values[0].Value, 10, 256, big.ToNearestEven)
	return result, err
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

func (e *esolver) ParseExpression(s string) (Stack, error) {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.Replace(s, ",", ".", -1)

	p := NewParser(strings.NewReader(s), e.elemNames)

	stack, err := p.Parse()
	if err != nil {
		return Stack{}, err
	}
	stack = addMissingOperator(stack)

	return stack, nil
}

func (e *esolver) AddConstant(name string, constCreator ConstFunction) {
	e.userConsts[name] = constCreator
	e.elemNames[name] = CONSTANT
}

func (e *esolver) findConst(name string) (c ConstFunction, ok bool) {
	if c, ok = consts[name]; ok {
		return c, ok
	}
	c, ok = e.userConsts[name]
	return
}
