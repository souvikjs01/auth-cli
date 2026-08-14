package cmd

import (
	"fmt"
)

func promptInput(prompt string) (string, error) {
	fmt.Print(prompt)

	return shell.Readline()
}
func promptPassword(prompt string) (string, error) {
	shell.SetMaskRune('*')
	defer shell.SetMaskRune(0)

	fmt.Print(prompt)

	return shell.Readline()
}
