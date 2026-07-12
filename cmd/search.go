package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	cfg "github.com/binarysoupdev/go-commando/config"
)

type SearchCommand struct {
	command.FlagCommandBase
	cfg.Loader[config.Config]
}

func NewSearchCommand(loader cfg.Loader[config.Config]) *SearchCommand {
	return &SearchCommand{
		FlagCommandBase: command.NewFlagCommandBase("search", "search records in the vault"),
		Loader:          loader,
	}
}

func (cmd *SearchCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()
	return flow.LoadConfig(&cmd.Loader)
}

func (cmd SearchCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.Flags.Parse(args)

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	return search.Display(v)
}
