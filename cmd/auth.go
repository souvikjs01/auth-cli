package cmd

import (
	"fmt"

	"github.com/souvikjs01/auth-cli/internals/service"
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
		email, err := promptInput("Email: ")
		if err != nil {
			return err
		}

		password, err := promptPassword("Password: ")
		if err != nil {
			return err
		}

		user, session, err := application.Login(
			email,
			password,
		)

		if err != nil {
			return err
		}

		fmt.Println("✓ Login successful")
		fmt.Printf("Welcome, %s\n", user.Name)
		fmt.Printf(
			"Session expires: %s\n",
			session.ExpiresAt.Format("15:04:05"),
		)

		return nil
	},
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new user account",

	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := promptInput("Name: ")
		if err != nil {
			return err
		}

		email, err := promptInput("Email: ")
		if err != nil {
			return err
		}

		password, err := promptPassword("Password: ")
		if err != nil {
			return err
		}

		confirmPassword, err := promptPassword("Confirm password: ")
		if err != nil {
			return err
		}

		if password != confirmPassword {
			return fmt.Errorf("passwords do not match")
		}

		user, err := application.AuthService.Register(
			service.RegisterInput{
				Name:     name,
				Email:    email,
				Password: password,
			},
		)

		if err != nil {
			return err
		}

		fmt.Printf(
			"✓ Account created successfully for %s\n",
			user.Email,
		)

		return nil
	},
}

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show available commands",

	Run: func(cmd *cobra.Command, args []string) {
		if application.CurrentSessionID == "" {
			fmt.Println("Available commands:")
			fmt.Println("  register    Create a new account")
			fmt.Println("  login       Login to your account")
			fmt.Println("  help        Show available commands")
			fmt.Println("  exit        Exit the application")

			return
		}

		fmt.Println("Available commands:")
		fmt.Println("  whoami      Show current user")
		fmt.Println("  enable-2fa  Enable MFA")
		fmt.Println("  disable-2fa Disable MFA")
		fmt.Println("  logout      Logout")
		fmt.Println("  help        Show available commands")
		fmt.Println("  exit        Exit the application")
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user details",

	RunE: func(cmd *cobra.Command, args []string) error {
		session, err := application.CurrentSession()
		if err != nil {
			return err
		}

		user := &session.User

		fmt.Printf("Name: %s\n", user.Name)
		fmt.Printf("Email: %s\n", user.Email)
		fmt.Printf(
			"Session expires: %s\n",
			session.ExpiresAt.Format("2006-01-02 15:04:05"),
		)

		if user.MFAEnabled {
			fmt.Println("MFA: Enabled")
		} else {
			fmt.Println("MFA: Disabled")
		}

		if user.LastLoginAt != nil {
			fmt.Printf(
				"Last login: %s\n",
				user.LastLoginAt.Format("2006-01-02 15:04:05"),
			)
		}

		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from the current session",

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := application.Logout(); err != nil {
			return err
		}

		fmt.Println("✓ Logged out successfully")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(helpCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(logoutCmd)

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
