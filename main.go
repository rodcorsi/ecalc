package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
)

func main() {

	fmt.Println(hello)

	reader := bufio.NewReader(os.Stdin)
	ctrlCAgain := false
	ecalc := NewECalc(os.Stdout)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	go func() {
		for sig := range c {
			if sig == syscall.SIGINT {
				if ctrlCAgain {
					exit(ecalc, "")
				}

				err := clipboard.WriteAll(strconv.FormatFloat(ecalc.Value, 'f', -1, 64))
				if err != nil {
					ecalc.Println("\nCan't copy to clipboard! Terminating...")
					exit(ecalc, "")
				}

				ecalc.Println("\nResult copied to clipboard!\nIf you want to close press CTRL+C again or type 'exit'.")
				ctrlCAgain = true
				ecalc.PrintPrompt()
			}
		}
	}()

	for {
		ctrlCAgain = false
		ecalc.PrintPrompt()

		text, _ := reader.ReadString('\n')
		text = strings.ToLower(strings.TrimSpace(text))

		if len(text) == 0 && ecalc.Partial && ecalc.Error == nil {
			text = ecalc.Expression
		}

		exprs := strings.Split(text, "|")

		for _, exp := range exprs {
			exp = strings.TrimSpace(exp)
			args := strings.SplitN(exp, " ", 2)
			cmdExecuted := false
			for k, v := range cmds {
				if k == args[0] {
					if len(args) == 1 {
						v(ecalc, "")
					} else {
						v(ecalc, args[1])
					}

					cmdExecuted = true
					break
				}
			}

			if cmdExecuted {
				continue
			}

			if len(exp) != 0 {
				ecalc.Eval(exp)
				ecalc.PrintResult()
			}
		}
	}
}
