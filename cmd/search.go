package cmd

import (
	"pvault/chain"
	"pvault/config"
	"pvault/vault"

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

	for i, match := range matches {
		style.New(style.YELLOW).Printf("[%d] %s\n", i, match)
	}

	return nil
}
