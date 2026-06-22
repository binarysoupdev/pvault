package cmd

import (
	"pvault/chain"
	"pvault/config"
	"pvault/vault"
	"strings"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type SearchCommand struct {
	command.FlagCommandBase
}

func NewSearchCommand() *SearchCommand {
	return &SearchCommand{
		FlagCommandBase: command.NewFlagCommandBase("search", "search records in the vault"),
	}
}

func (cmd SearchCommand) Run(args []string) error {
	term := cmd.Flags.String("s", "", "the search term")
	cmd.Flags.Parse(args)

	v, err := vault.Open(config.Global.VaultPath)
	if err != nil {
		return chain.Error(err, "error opening vault")
	}

	matches := v.Search(*term)

	result := style.New(style.YELLOW)
	highlight := append(result, style.UNDERLINE)

	for i, match := range matches {
		start := strings.Index(strings.ToLower(match), strings.ToLower(*term))
		end := start + len(*term)

		result.Printf("[%d] %s", i+1, match[:start])
		highlight.Print(match[start:end])
		result.Println(match[end:])
	}

	return nil
}
