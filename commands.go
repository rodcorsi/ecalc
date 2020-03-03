package main

import "os"
import "regexp"
import "strings"

const version = "v0.3"

var cmds = map[string]func(*ECalc, string){
	"exit":  exit,
	"help":  help,
	"dms":   dms,
	"clear": clear,
	"cls":   clear,
	"set":   set,
}

const hello = "Ecalc (" + version + `) - Engineer command calculator
type 'help' for more informations or 'exit' to leave
`

const helpText = `Ecalc - Engineer command calculator
Expressions:
	5+2
	15*pi
	tan45
	(4+5)*cos45d25'33.15"
CTRL+C:
	Copy result to clipboard
Commands:
	exit	terminate this
	help	show this text
	dms		print last result to Degree Minutes Seconds
	clear	clear all screen
	cls		same as clear command
Operator:
	+ - * / ^
Functions:
	ln abs cos sin tan acos asin atan sqrt cbrt ceil floor
Constants:
	e pi phi sqrt2 sqrte sqrtpi sqrtphi ans
ANS:
	you can use an special constant 'ans' to put last result in your expression
`

func exit(e *ECalc, args string) {
	e.Println("Bye!")
	os.Exit(0)
}

func help(e *ECalc, args string) {
	e.Println(helpText)
}

func dms(e *ECalc, args string) {
	e.Degree = true
	e.PrintResult()
}

func clear(e *ECalc, args string) {
	for i := 0; i < 30; i++ {
		e.Println()
	}
}

var reSet = regexp.MustCompile(`[a-zA-Z]+`)

func set(e *ECalc, args string) {
	varName := strings.TrimSpace(strings.ToLower(args))
	if varName == "" {
		varName = "x"
	} else if !reSet.MatchString(varName) {
		e.Printf("Invalid variable name '%v'. Must be just letters\n", args)
		return
	}

	e.AddConstant(varName, e.Value)
	e.Printf("%v => %.12f\n", varName, e.Value)
}
