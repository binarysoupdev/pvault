package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/errors"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
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
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.Flags.Parse(args)

	v, err := vault.Open(config.Global.VaultPath)
	if err != nil {
		return errors.Chain(err, "error opening vault")
	}

	return search.Display(v)
}
