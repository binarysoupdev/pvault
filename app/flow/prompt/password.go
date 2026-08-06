package prompt

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func Password(prompt string) string {
	for {
		fmt.Print(prompt)
		if password := readPassword(); password != "" {
			return password
		}
	}
}

func readPassword() string {
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return readStdin()
	}

	fmt.Println()
	return string(password)
}
