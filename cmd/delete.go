package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/json"
)

type DeleteCommand struct {
	command.FlagCommandBase

	ConfigLoader json.Loader[config.Config]
	Config       config.Config
}

func NewDeleteCommand(loader json.Loader[config.Config]) *DeleteCommand {
	return &DeleteCommand{
		FlagCommandBase: command.NewFlagCommandBase("delete", "delete a record from the vault"),
		ConfigLoader:    loader,
	}
}

func (cmd *DeleteCommand) Initialize() error {
	var err error
	cmd.Config, err = flow.LoadConfig(cmd.ConfigLoader)
	if err != nil {
		return err
	}

	return cmd.FlagCommandBase.Initialize()
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
