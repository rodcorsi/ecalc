package esolver

import (
	"math/big"
	"reflect"
	"testing"
)

func Test_esolver_Solve(t *testing.T) {
	assert := func(value string, expected *big.Float, isError bool) {

		x, err := New().Solve(value)
		if (err != nil) != isError {
			t.Errorf("Expected error(%v) value:`%v` result:`%v` Error:%v", isError, value, x, err)
			return
		}

		diff := new(big.Float).Sub(x, expected)
		diff.Abs(diff)
		if diff.Cmp(big.NewFloat(0.0000000001)) > 0 {
			t.Errorf("Value:`%v` expected:`%v` result:`%v`", value, expected, x)
		}
	}

	assert(".", big.NewFloat(0), false)
	assert("1", big.NewFloat(1), false)
	assert(".1", big.NewFloat(0.1), false)
	assert("1.", big.NewFloat(1), false)
	assert("1.5", big.NewFloat(1.5), false)
	assert("1.5+1", big.NewFloat(2.5), false)
	assert("1.5*2", big.NewFloat(3), false)
	assert("1.5*2+2", big.NewFloat(5), false)
	assert("1+3-2^3/5*3+2*3", big.NewFloat(1+3-8.0/5*3+2*3), false)

	assert(`45d15'25"`, degToDec(`45d15'25"`), false)
	assert(`45d15'`, degToDec(`45d15'`), false)
	assert(`45d25"`, degToDec(`45d25"`), false)
	assert(`15'25"`, degToDec(`15'25"`), false)
	assert(`45d`, degToDec(`45d`), false)
	assert(`45'`, degToDec(`45'`), false)
	assert(`45"`, degToDec(`45"`), false)

	assert(`45d15m25s`, degToDec(`45d15m25s`), false)
	assert(`45d15m`, degToDec(`45d15m`), false)
	assert(`45d25s`, degToDec(`45d25s`), false)
	assert(`15m25s`, degToDec(`15m25s`), false)
	assert(`45m`, degToDec(`45m`), false)
	assert(`45s`, degToDec(`45s`), false)

	assert(`5sin90`, big.NewFloat(5), false)

	assert(`tan45d15'25"`, tan(degToDec(`45d15'25"`)), false)
	assert(`tan45d15'25"+5`, new(big.Float).Add(tan(degToDec(`45d15'25"`)), big.NewFloat(5)), false)
	assert(`tan45d15'25"*5`, new(big.Float).Mul(tan(degToDec(`45d15'25"`)), big.NewFloat(5)), false)
	assert(`5tan45`, new(big.Float).Mul(big.NewFloat(5), tan(degToDec(`45`))), false)
	assert(`tan(tan30)`, tan(tan(big.NewFloat(30))), false)
	assert(`pi3`, new(big.Float).Mul(consts["pi"](), big.NewFloat(3)), false)
	assert(`pi*pi`, new(big.Float).Mul(consts["pi"](), consts["pi"]()), false)
	assert(`5(5)`, big.NewFloat(25), false)
	assert(`(5)(5)`, big.NewFloat(25), false)
}

func Test_esolver_ParseExpression(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    Stack
		wantErr bool
	}{
		{"literal", "1", Stack{[]Token{{NUMBER, "1"}}}, false},
		{"literal more numbers", "100", Stack{[]Token{{NUMBER, "100"}}}, false},
		{"addition", "1+2", Stack{[]Token{{NUMBER, "1"}, {OPERATOR, "+"}, {NUMBER, "2"}}}, false},
		{"subtraction", "3-1", Stack{[]Token{{NUMBER, "3"}, {OPERATOR, "-"}, {NUMBER, "1"}}}, false},
		{"multiplication", "2*3", Stack{[]Token{{NUMBER, "2"}, {OPERATOR, "*"}, {NUMBER, "3"}}}, false},
		{"division", "4/2", Stack{[]Token{{NUMBER, "4"}, {OPERATOR, "/"}, {NUMBER, "2"}}}, false},
		{"exponent", "2^3", Stack{[]Token{{NUMBER, "2"}, {OPERATOR, "^"}, {NUMBER, "3"}}}, false},
		{"parentheses", "(1+2)", Stack{[]Token{{LPAREN, "("}, {NUMBER, "1"}, {OPERATOR, "+"}, {NUMBER, "2"}, {RPAREN, ")"}}}, false},
		{"function", "sin(90)", Stack{[]Token{{FUNCTION, "sin"}, {LPAREN, "("}, {NUMBER, "90"}, {RPAREN, ")"}}}, false},
		{"constant", "pi", Stack{[]Token{{CONSTANT, "pi"}}}, false},
		{"unary minus", "-1", Stack{[]Token{{OPERATOR, "-"}, {NUMBER, "1"}}}, false},
		{"unary plus", "+1", Stack{[]Token{{OPERATOR, "+"}, {NUMBER, "1"}}}, false},
		{"complex expression", "5.5*pi+sin(90-10)", Stack{[]Token{
			{NUMBER, "5.5"},
			{OPERATOR, "*"},
			{CONSTANT, "pi"},
			{OPERATOR, "+"},
			{FUNCTION, "sin"},
			{LPAREN, "("},
			{NUMBER, "90"},
			{OPERATOR, "-"},
			{NUMBER, "10"},
			{RPAREN, ")"},
		}}, false},
		{"implicit multiplication constant", "3pi", Stack{[]Token{{NUMBER, "3"}, {OPERATOR, "*"}, {CONSTANT, "pi"}}}, false},
		{"implicit multiplication function", "3sin(90)", Stack{[]Token{{NUMBER, "3"}, {OPERATOR, "*"}, {FUNCTION, "sin"}, {LPAREN, "("}, {NUMBER, "90"}, {RPAREN, ")"}}}, false},
		{"implicit multiplication parentheses", "3(4)", Stack{[]Token{{NUMBER, "3"}, {OPERATOR, "*"}, {LPAREN, "("}, {NUMBER, "4"}, {RPAREN, ")"}}}, false},
		{"implicit multiplication parentheses 2", "(3)4", Stack{[]Token{{LPAREN, "("}, {NUMBER, "3"}, {RPAREN, ")"}, {OPERATOR, "*"}, {NUMBER, "4"}}}, false},
		{"implicit multiplication parentheses 3", "(3)(4)", Stack{[]Token{{LPAREN, "("}, {NUMBER, "3"}, {RPAREN, ")"}, {OPERATOR, "*"}, {LPAREN, "("}, {NUMBER, "4"}, {RPAREN, ")"}}}, false},
		{"dms", `45d15'25"`, Stack{[]Token{{NUMBER, `45.2569444444444444444444444444444444444444444444444444444444444444444444444444`}}}, false},
		{"tokenizer test", "1+*2", Stack{[]Token{{NUMBER, "1"}, {OPERATOR, "+"}, {OPERATOR, "*"}, {NUMBER, "2"}}}, false},
	}
	e := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := e.ParseExpression(tt.s)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseExpression() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ParseExpression() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}
