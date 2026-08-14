package cmd

func promptInput(prompt string) (string, error) {
	saved := shell.Config.Prompt
	shell.SetPrompt(prompt)
	defer shell.SetPrompt(saved)

	return shell.Readline()
}

func promptPassword(prompt string) (string, error) {
	pw, err := shell.ReadPassword(prompt)
	return string(pw), err
}
