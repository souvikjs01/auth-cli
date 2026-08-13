package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	loginEmail    string
	loginPassword string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your account",

	RunE: func(cmd *cobra.Command, args []string) error {
		user, session, err := application.Login(
			loginEmail,
			loginPassword,
		)

		if err != nil {
			return err
		}

		fmt.Printf("✓ Login successful\n")
		fmt.Printf("Welcome, %s\n", user.Name)
		fmt.Printf("Session expires: %s\n",
			session.ExpiresAt.Format("15:04:05"),
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(
		&loginEmail,
		"email",
		"e",
		"",
		"Email address",
	)

	loginCmd.Flags().StringVarP(
		&loginPassword,
		"password",
		"p",
		"",
		"Password",
	)
}
