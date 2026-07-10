package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
)

type DeleteCommand struct {
	command.FlagCommandBase
	config.Loader[config.Config]
}

func NewDeleteCommand(loader config.Loader[config.Config]) *DeleteCommand {
	return &DeleteCommand{
		FlagCommandBase: command.NewFlagCommandBase("delete", "delete a record from the vault"),
		Loader:          loader,
	}
}

func (cmd *DeleteCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()
	return flow.LoadConfig(&cmd.Loader)
}

func (cmd DeleteCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.Flags.Parse(args)

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	name, err := search.Select(v)
	if err != nil {
		return err
	}

	err = flow.DeleteVaultRecord(v, name)
	if err != nil {
		return err
	}

	return nil
}
