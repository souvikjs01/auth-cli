package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "auth-cli",
	Short: "CLI authentication system",
}

func Execute() {
	initApp()
	StartShell()
}
