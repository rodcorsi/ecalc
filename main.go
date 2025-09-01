package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/atotto/clipboard"
)

func main() {
	shell := ishell.New()
	ecalc := NewECalc()

	AddCommands(shell, ecalc)

	shell.NotFound(func(c *ishell.Context) {
		result := ecalc.Eval(strings.Join(c.Args, " "))
		c.Println(ResultLine(result))
		c.SetPrompt(Prompt(lastAnswer))
	})

	shell.SetPrompt(Prompt(ecalc.Result))
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

func Prompt(result *Result) string {
	return fmt.Sprintf("(ans:%v) » ", result.FormatValue())
}

func ResultLine(result *Result) string {
	expr := result.FormatExpression()
	value := result.FormatResult()
	return fmt.Sprintf("%v = %v", expr, value)
}

func copyToClipboard(c *ishell.Context, ecalc *ECalc) (string, error) {
	value := ecalc.Result.Value.Text('f', -1)
	return value, clipboard.WriteAll(value)
}
