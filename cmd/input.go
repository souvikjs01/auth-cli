package cmd

// promptInput temporarily changes the readline prompt to the given text
// and reads a line of input. The original prompt is restored afterwards.
func promptInput(prompt string) (string, error) {
	saved := shell.Config.Prompt
	shell.SetPrompt(prompt)
	defer shell.SetPrompt(saved)

	return shell.Readline()
}

// promptPassword reads a password from the user with masked input.
// Uses readline's built-in ReadPassword which hides characters as they are typed.
func promptPassword(prompt string) (string, error) {
	pw, err := shell.ReadPassword(prompt)
	return string(pw), err
}
