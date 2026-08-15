package cmds

import (
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	vault_flow "pvault/app/flow/vault"
	search_flow "pvault/app/flow/vault/search"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/json"
)

type DeleteCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]
}

func NewDeleteCommand(configLoader json.Loader[config.Config]) *DeleteCommand {
	return &DeleteCommand{
		CommandBase:   command.NewCommandBase("delete", "Delete a record from the vault"),
		ConfigCommand: command.NewConfigCommand(configLoader),
	}
}

func (cmd *DeleteCommand) Initialize() error {
	err := config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
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
