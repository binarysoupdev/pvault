package cmd

import (
	"pvault/chain"
	"pvault/cmd/flow"
	"pvault/config"
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
		return chain.Error(err, "error opening vault")
	}

	search.Display(v)
	return nil
}
