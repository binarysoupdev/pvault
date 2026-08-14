package cmds

import (
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	output_flow "pvault/app/flow/output"
	vault_flow "pvault/app/flow/vault"
	search_flow "pvault/app/flow/vault/search"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/go-extensions/errors"
)

type UnlockCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]
}

func NewUnlockCommand(configLoader json.Loader[config.Config]) *UnlockCommand {
	return &UnlockCommand{
		CommandBase:   command.NewCommandBase("unlock", "Unlock a record from the vault"),
		ConfigCommand: command.NewConfigCommand(configLoader),
	}
}

func (cmd *UnlockCommand) Initialize() error {
	err := config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
		return err
	}

	err = cmd.Config.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating \"config.output_path\"")
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd UnlockCommand) Run(args []string) error {
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

	r, err := vault_flow.LoadRecord(v, name)
	if err != nil {
		return err
	}

	err = output_flow.SaveRecord(cmd.Config, r)
	if err != nil {
		return err
	}

	return nil
}
