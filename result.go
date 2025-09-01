package main

import (
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/rodcorsi/ecalc/calc"
)

type Result struct {
	Value       *big.Float
	Degree      bool
	Error       error
	Writer      io.Writer
	EngNotation bool
	Partial     bool
	Expression  string
	StackExpr   calc.Stack
}

func (c *Result) FormatValue() string {
	if c.EngNotation {
		return fmtPrompt.Sprint(c.Value.Text('e', 14))
	}
	return fmtPrompt.Sprint(c.Value.Text('f', 8))
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
		return fmtResult.Sprint(convertDMS(c.Value))
	} else if c.EngNotation {
		return fmtResult.Sprintf("%e", c.Value)
	}
	return fmtResult.Sprintf("%s", formatRecurring(c.Value, 20))
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
	parts := strings.Split(s, ".")
	integerPart := parts[0]
	fractionalPartStr := parts[1]

	for start := 0; start < len(fractionalPartStr); start++ {
		for length := 1; start+length*2 <= len(fractionalPartStr); length++ {
			pattern := fractionalPartStr[start : start+length]
			nextPart := fractionalPartStr[start+length : start+length*2]

			if pattern == nextPart && pattern != strings.Repeat("0", length) {
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
	return result
}

func addOverline(text string) string {
	var builder strings.Builder
	for _, char := range text {
		builder.WriteString("\u0305")
		builder.WriteRune(char)
	}
	return builder.String()
}
