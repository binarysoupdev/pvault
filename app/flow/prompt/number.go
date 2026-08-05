package prompt

import "strconv"

func Number(prompt string, fallback int) int {
	num, err := strconv.Atoi(Prompt(prompt))
	if err != nil {
		return fallback
	}
	return num
}
