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
	config.Loader[config.Config]
}

func NewSearchCommand(loader config.Loader[config.Config]) *SearchCommand {
	return &SearchCommand{
		FlagCommandBase: command.NewFlagCommandBase("search", "search records in the vault"),
		Loader:          loader,
	}
}

func (cmd *SearchCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()

	err := cmd.LoadConfig()
	if err != nil {
		return errors.Chain(err, "error loading config")
	}

	return nil
}

func (cmd SearchCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.Flags.Parse(args)

	v, err := vault.Open(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error opening vault")
	}

	return search.Display(v)
}
