package ecalc

import (
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/rodcorsi/ecalc/esolver"
)

type Result struct {
	Value       *big.Float
	Degree      bool
	Error       error
	Writer      io.Writer
	EngNotation bool
	Partial     bool
	Expression  string
	StackExpr   esolver.Stack
}

func (e *Result) FormatExpression(printer func(value string, t esolver.Token)) {
	lastType := esolver.TokenType(-1)
	for _, v := range e.StackExpr.Values {
		closeParen := false
		if lastType == esolver.FUNCTION && (v.Type == esolver.NUMBER || v.Type == esolver.CONSTANT) {
			printer("(", v)
			closeParen = true
		}

		if v.Type == esolver.FUNCTION || v.Type == esolver.CONSTANT {
			printer(v.Value, v)
		} else if v.Type == esolver.NUMBER {
			printer(v.Value, v)
		} else if v.Value == "+" || v.Value == "-" {
			printer(" "+v.Value+" ", v)
		} else {
			printer(v.Value, v)
		}
		if closeParen {
			printer(")", v)
		}
		lastType = v.Type
	}
}

func (c *Result) String() string {
	if c.Error != nil {
		return c.Error.Error()
	} else if c.Degree {
		return convertDMS(c.Value)
	} else if c.EngNotation {
		return fmt.Sprintf("%e", c.Value)
	}
	return formatRecurring(c.Value, 20)
}
func convertDMS(value *big.Float) string {
	v, _ := value.Float64()
	d := int64(v)
	m := (v - float64(d)) * 60.0
	s := (m - float64(int64(m))) * 60.0
	return fmt.Sprintf(`%vd%2d'%.5f"`, d, int64(m), s)
}

func formatRecurring(n *big.Float, precision int) string {
	if n.Cmp(big.NewFloat(0)) == 0 {
		return "0"
	}

	var sign string
	if n.Sign() < 0 {
		sign = "-"
		n.Abs(n)
	}

	s := n.Text('f', precision)

	if strings.Contains(s, ".") {
		parts := strings.Split(s, ".")
		trimmedFractional := strings.TrimRight(parts[1], "0")

		isAllNines := len(trimmedFractional) > 10
		if isAllNines {
			for _, r := range trimmedFractional {
				if r != '9' {
					isAllNines = false
					break
				}
			}
		}

		if isAllNines {
			integer, _ := new(big.Int).SetString(parts[0], 10)
			integer.Add(integer, big.NewInt(1))
			return sign + integer.String()
		}
	}

	s = strings.TrimRight(s, "0")
	if !strings.Contains(s, ".") {
		s = s + "."
	}

	parts := strings.Split(s, ".")
	integerPart := parts[0]
	fractionalPartStr := parts[1]

	if len(fractionalPartStr) == 0 {
		if integerPart == "" {
			return "0"
		}
		return sign + integerPart
	}

	for start := 0; start < len(fractionalPartStr); start++ {
		for length := 1; start+length <= len(fractionalPartStr); length++ {
			if start+length*2 > len(fractionalPartStr) {
				continue
			}

			pattern := fractionalPartStr[start : start+length]
			isRepeating := true
			for i := start + length; i < len(fractionalPartStr); i++ {
				if fractionalPartStr[i] != pattern[(i-start)%length] {
					if i == len(fractionalPartStr)-1 && fractionalPartStr[i] == pattern[(i-start)%length]+1 {
						// Rounding detected, ignore mismatch for last digit
					} else {
						isRepeating = false
						break
					}
				}
			}

			if isRepeating {
				if len(fractionalPartStr)-start < 6 {
					continue
				}
				nonRepeatingPart := fractionalPartStr[:start]
				repeatingStr := addOverline(pattern)
				return fmt.Sprintf("%s%s.%s%s", sign, integerPart, nonRepeatingPart, repeatingStr)
			}
		}
	}

	result := n.Text('f', -1)
	if strings.Contains(result, ".") {
		result = strings.TrimRight(result, "0")
		result = strings.TrimRight(result, ".")
	}
	return sign + result
}

func addOverline(text string) string {
	var builder strings.Builder
	for _, char := range text {
		builder.WriteString("\u0305")
		builder.WriteRune(char)
	}
	return builder.String()
}
