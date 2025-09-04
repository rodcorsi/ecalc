package main

import (
	"regexp"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/rodcorsi/ecalc"
)

const version = "v0.4"

const hello = "Ecalc (" + version + `) - Engineer command calculator
type 'help' for more informations or 'exit' to leave
`

const helpText = `Ecalc - Engineer command calculator
Expressions example:
    5+2
    15*pi
    tan45
    (4+5)*cos45d25m33.15s
    *2
	/3
CTRL+C:
    Copy result to clipboard
Commands:
    exit  terminate this
    help  show this text
    dms   print last result to Degree Minutes Seconds
    clear clear all screen
    cls   same as clear command
    set   define variable with last value
    cp    copy to clipboard
Operator:
    + - * / ^
Functions:
    ln abs cos sin tan acos asin atan sqrt cbrt ceil floor
Constants:
    e pi phi sqrt2 sqrte sqrtpi sqrtphi ans
ANS:
    you can use an special variable 'ans' to use the last result on your expression
`

var reValidSetName = regexp.MustCompile(`^[a-zA-Z]+$`)

func addCommands(shell *ishell.Shell, ecalc *ecalc.ECalc) {
	shell.AddCmd(&ishell.Cmd{
		Name: "help",
		Help: "Help",
		Func: func(c *ishell.Context) {
			c.Println(helpText)
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "cls",
		Help: "clear screen",
		Func: func(c *ishell.Context) {
			c.ClearScreen()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "dms",
		Help: "print last result to Degree Minutes Seconds",
		Func: func(c *ishell.Context) {
			ecalc.Result.Degree = true
			c.Println(resultLine(ecalc.Result))
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "set",
		Help: "define variable with last value",
		Func: func(c *ishell.Context) {
			expr := strings.Join(c.Args, " ")
			varName := strings.TrimSpace(expr)
			if varName == "" {
				varName = "x"
			} else if !reValidSetName.MatchString(varName) {
				c.Printf("Invalid variable name '%v'. Must be just letters\n", expr)
				return
			}
			value := ecalc.Result.Value
			ecalc.AddConstant(varName, value)
			c.Printf("%v => %v\n", varName, formatResult(ecalc.Result))
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "cp",
		Help: "Copy to clipboard",
		Func: func(c *ishell.Context) {
			value, err := copyToClipboard(c, ecalc)
			if err != nil {
				c.Println("\nCan't copy to clipboard! ")
				return
			}
			c.Printf("\n'%v' copied to clipboard!\n", value)
		},
	})
}
