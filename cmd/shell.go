// Package cmd implements the interactive CLI shell and all user-facing commands.
package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
)

var shell *readline.Instance

// StartShell initializes the interactive readline prompt and enters the
// main command loop. It supports command history and handles interrupt/EOF
// signals gracefully.
func StartShell() {
	var err error
	shell, err = readline.NewEx(&readline.Config{
		Prompt:          "auth> ",
		HistoryFile:     ".auth_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})

	if err != nil {
		fmt.Printf("failed to start CLI: %v\n", err)
		return
	}

	defer shell.Close()

	fmt.Println("CLI Authentication System")
	fmt.Println("Type 'help' to see available commands.")
	fmt.Println()

	for {
		line, err := shell.Readline()

		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			}

			continue
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("input error: %v\n", err)
			break
		}

		command := strings.TrimSpace(line)

		if command == "" {
			continue
		}

		if command == "exit" {
			break
		}

		handleCommand(command)
	}

	fmt.Println("Goodbye!")
}

// handleCommand parses user input and dispatches it to the appropriate
// cobra subcommand.
func handleCommand(command string) {
	args := strings.Fields(command)

	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
