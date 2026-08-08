package record

import (
	"pvault/app/commands/base"
	"pvault/app/config"
	vault_flow "pvault/app/flow/vault"
	search_flow "pvault/app/flow/vault/search"

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
	search := search_flow.NewSearchFlow(cmd.Flags)
	cmd.ParseFlags(args)

	v, err := vault_flow.OpenCurrentVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	name, err := search.Select(v)
	if err != nil {
		return err
	}

	err = vault_flow.DeleteRecord(v, name)
	if err != nil {
		return err
	}

	return nil
}
