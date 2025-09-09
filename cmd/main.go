package main

import (
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/rodcorsi/ecalc"
	"github.com/rodcorsi/ecalc/esolver"
)

var (
	fmtPrompt   = color.New(color.FgCyan)
	fmtResult   = color.New(color.FgYellow)
	fmtError    = color.New(color.FgRed)
	fmtFunction = color.New(color.FgYellow)
)

func main() {
	shell := ishell.New()
	ecalc := ecalc.NewECalc()

	addCommands(shell, ecalc)

	shell.NotFound(func(c *ishell.Context) {
		result := ecalc.Eval(strings.Join(c.Args, " "))
		c.Println(resultLine(result))
		c.SetPrompt(prompt(ecalc.LastAnswer))
	})

	shell.SetPrompt(prompt(ecalc.Result))
	shell.Println(hello)
	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		if count == 1 {
			if clipboard.Unsupported {
				os.Exit(0)
			} else {
				value, err := copyToClipboard(c, ecalc)
				if err != nil {
					c.Println("\nCan't copy to clipboard! Terminating...")
					os.Exit(1)
				}
				c.Printf("\n'%v' copied to clipboard! Press CTRL+C again to exit.\n", value)
			}
			return
		}
		os.Exit(0)
	})
	shell.IgnoreCase(true)
	shell.SetHomeHistoryPath(".ecalc_history")
	shell.Run()
}

func prompt(result *ecalc.Result) string {
	return fmt.Sprintf("(ans:%v) » ", fmtPrompt.Sprintf("%-10s", formatValue(result.Value)))
}

func formatValue(value *big.Float) string {
	const maxLen = 10

	if value.IsInt() {
		formattedInt := value.Text('f', 0)
		if len(formattedInt) <= maxLen {
			return fmt.Sprint(formattedInt)
		}
	}

	var integerPart big.Int
	value.Int(&integerPart)

	integerPartLength := len(integerPart.String())
	if integerPart.Sign() == 0 {
		integerPartLength = 1
	}

	precision := maxLen - integerPartLength - 1
	if precision < 0 {
		precision = 0
	}

	formattedFloat := value.Text('f', precision)
	if precision > 0 {
		formattedFloat = strings.TrimRight(formattedFloat, "0")
		formattedFloat = strings.TrimRight(formattedFloat, ".")
	}

	isTooLong := len(formattedFloat) > maxLen
	isRoundedToZero := formattedFloat == "0" && value.Sign() != 0

	if isTooLong || isRoundedToZero {
		return fmt.Sprint(value.Text('e', 4))
	}

	return fmt.Sprint(formattedFloat)
}

func resultLine(result *ecalc.Result) string {
	var sb strings.Builder
	result.FormatExpression(func(value string, t esolver.Token) {
		if t.Type == esolver.FUNCTION || t.Type == esolver.CONSTANT {
			sb.WriteString(fmtFunction.Sprint(value))
		} else {
			sb.WriteString(value)
		}
	})
	return fmt.Sprintf("%v = %v", sb.String(), formatResult(result))
}

func formatResult(c *ecalc.Result) string {
	if c.Error != nil {
		return fmtError.Sprint("Error:", c.Error.Error())
	}
	return fmtResult.Sprint(c.String())
}

func copyToClipboard(c *ishell.Context, ecalc *ecalc.ECalc) (string, error) {
	value := ecalc.Result.Value.Text('f', -1)
	return value, clipboard.WriteAll(value)
}
