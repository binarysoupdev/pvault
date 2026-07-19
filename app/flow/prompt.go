package flow

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

func Prompt(prompt string) string {
	fmt.Print(prompt)
	return readStdin()
}

func PromptPassword(prompt string) string {
	fmt.Print(prompt)

	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return readStdin()
	}

	fmt.Println()
	return string(password)
}

func readStdin() string {
	res, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return res[:len(res)-1]
}
