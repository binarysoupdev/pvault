package cmd

import (
	"pvault/app/cmd/base"
	"pvault/app/flow"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/json"
)

type DeleteCommand struct {
	command.CommandBase
	command.FlagCommand
	base.ConfigCommand
}

func NewDeleteCommand(configLoader json.Loader[config.Config]) *DeleteCommand {
	return &DeleteCommand{
		CommandBase:   command.NewCommandBase("delete", "delete a record from the vault"),
		ConfigCommand: base.NewConfigCommand(configLoader),
	}
}

func (cmd *DeleteCommand) Initialize() error {
	if err := cmd.LoadConfig(); err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd DeleteCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.ParseFlags(args)

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
