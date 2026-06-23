package flow

import (
	"bufio"
	"fmt"
	"os"
)

func Prompt(prompt string) string {
	fmt.Print(prompt)

	res, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return res[:len(res)-1]
}
