package prompt

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func Password(prompt string) string {
	fmt.Print(prompt)

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return readStdin()
	}

	fmt.Println()
	return string(password)
}
