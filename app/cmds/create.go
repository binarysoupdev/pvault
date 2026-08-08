package cmds

import (
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	output_flow "pvault/app/flow/output"
	vault_flow "pvault/app/flow/vault"
	v2 "pvault/app/vault/record/record/v2"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
)

type CreateCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]
}

func NewCreateCommand(configLoader json.Loader[config.Config]) *CreateCommand {
	return &CreateCommand{
		CommandBase:   command.NewCommandBase("create", "create a new vault record"),
		ConfigCommand: command.NewConfigCommand(configLoader),
	}
}

func (cmd *CreateCommand) Initialize() error {
	err := config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd CreateCommand) Run(args []string) error {
	name := cmd.Flags.String("name", "", "name of the record")
	cmd.ParseFlags(args)

	if *name == "" {
		return errors.New("\"name\" cannot be empty")
	}

	v, err := vault_flow.OpenCurrentVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	err = cmd.Config.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating config \"output_path\"")
	}

	r := v2.NewEmptyRecord(*name)

	err = vault_flow.SaveRecord(v, r)
	if err != nil {
		return err
	}

	err = output_flow.SaveRecord(cmd.Config, r)
	if err != nil {
		return err
	}

	return nil
}
