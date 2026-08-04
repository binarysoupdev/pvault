package prompt

import (
	"bufio"
	"fmt"
	"os"
)

func Prompt(prompt string) string {
	fmt.Print(prompt)
	return readStdin()
}

func readStdin() string {
	res, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return res[:len(res)-1]
}
