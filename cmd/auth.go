// Package cmd implements the interactive CLI shell and all user-facing commands.
// Commands are split into pre-login (register, login, help, exit) and
// post-login (whoami, enable-2fa, disable-2fa, logout, help, exit).
package cmd

import (
	"fmt"

	"github.com/souvikjs01/auth-cli/internals/service"
	"github.com/spf13/cobra"
)

// loginCmd authenticates a user with username and password.
// On success it creates a session and displays the user's profile details
// including registration date, MFA status, session expiration, and last login.
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your account",

	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := promptInput("Username: ")
		if err != nil {
			return err
		}

		password, err := promptPassword("Password: ")
		if err != nil {
			return err
		}

		user, session, err := application.Login(
			username,
			password,
		)

		if err != nil {
			return err
		}

		// Auto-display user details after login (PRD section 5).
		fmt.Println("✓ Login successful")
		fmt.Println()
		fmt.Printf("  Username:           %s\n", user.Username)
		fmt.Printf("  Registered:         %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))

		if user.MFAEnabled {
			fmt.Println("  MFA:                Enabled")
		} else {
			fmt.Println("  MFA:                Disabled")
		}

		fmt.Printf("  Session expires:    %s\n",
			session.ExpiresAt.Format("2006-01-02 15:04:05"),
		)

		if user.LastLoginAt != nil {
			fmt.Printf("  Last login:         %s\n",
				user.LastLoginAt.Format("2006-01-02 15:04:05"),
			)
		}

		return nil
	},
}

// registerCmd creates a new user account with username and password.
// Password confirmation is required to prevent typos.
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a new user account",

	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := promptInput("Username: ")
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
				Username: username,
				Password: password,
			},
		)

		if err != nil {
			return err
		}

		fmt.Printf(
			"✓ Account created successfully for %s\n",
			user.Username,
		)

		return nil
	},
}

// helpCmd displays available commands based on the current authentication state.
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show available commands",

	Run: func(cmd *cobra.Command, args []string) {
		// Show different command sets depending on login state.
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

// whoamiCmd displays the current authenticated user's profile details.
var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user details",

	RunE: func(cmd *cobra.Command, args []string) error {
		session, err := application.CurrentSession()
		if err != nil {
			return err
		}

		user := &session.User

		fmt.Printf("Username:           %s\n", user.Username)
		fmt.Printf("Registered:         %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))

		if user.MFAEnabled {
			fmt.Println("MFA:                Enabled")
		} else {
			fmt.Println("MFA:                Disabled")
		}

		fmt.Printf(
			"Session expires:    %s\n",
			session.ExpiresAt.Format("2006-01-02 15:04:05"),
		)

		if user.LastLoginAt != nil {
			fmt.Printf(
				"Last login:         %s\n",
				user.LastLoginAt.Format("2006-01-02 15:04:05"),
			)
		}

		return nil
	},
}

// logoutCmd ends the current session and clears session state.
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

// enable2FACmd sets up TOTP-based multi-factor authentication.
// It generates a secret key, displays it for the user to add to an
// authenticator app, then verifies a code before enabling MFA.
var enable2FACmd = &cobra.Command{
	Use:   "enable-2fa",
	Short: "Enable TOTP-based MFA",

	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := application.CurrentUser()
		if err != nil {
			return err
		}

		if user.MFAEnabled {
			return service.ErrMFAAlreadyEnabled
		}

		key, err := application.TokenService.GenerateKey(
			user.Username,
		)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("Add this account to Google Authenticator.")
		fmt.Println()
		fmt.Println("Secret:")
		fmt.Println(key.Secret())
		fmt.Println()
		fmt.Println("OTP URL:")
		fmt.Println(key.URL())
		fmt.Println()

		code, err := promptInput("Enter the 6-digit code: ")
		if err != nil {
			return err
		}

		if !application.TokenService.Verify(
			code,
			key.Secret(),
		) {
			return service.ErrInvalidTOTP
		}

		if err := application.AuthService.EnableMFA(
			user,
			key.Secret(),
		); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("✓ MFA enabled successfully")

		return nil
	},
}

// disable2FACmd removes TOTP-based MFA from the user's account.
var disable2FACmd = &cobra.Command{
	Use:   "disable-2fa",
	Short: "Disable TOTP-based MFA",

	RunE: func(cmd *cobra.Command, args []string) error {
		user, err := application.CurrentUser()
		if err != nil {
			return err
		}

		if err := application.AuthService.DisableMFA(user); err != nil {
			return err
		}

		fmt.Println("✓ MFA disabled successfully")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(helpCmd)
	rootCmd.AddCommand(whoamiCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(enable2FACmd)
	rootCmd.AddCommand(disable2FACmd)
}
