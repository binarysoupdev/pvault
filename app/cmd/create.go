package cmd

import (
	"pvault/app/flow"
	v2 "pvault/app/vault/record/version2"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
)

type CreateCommand struct {
	command.CommandBase
	command.FlagCommand
	ConfigCommandBase
}

func NewCreateCommand(configLoader json.Loader[config.Config]) *CreateCommand {
	return &CreateCommand{
		CommandBase:       command.NewCommandBase("create", "create a new vault record"),
		ConfigCommandBase: NewConfigCommandBase(configLoader),
	}
}

func (cmd *CreateCommand) Initialize() error {
	if err := cmd.LoadConfig(); err != nil {
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

	v, err := flow.OpenLocalVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	err = cmd.Config.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating config \"output_path\"")
	}

	r := v2.NewEmptyRecord(*name)

	err = flow.SaveVaultRecord(v, r)
	if err != nil {
		return err
	}

	err = flow.SaveOutputRecord(cmd.Config, r)
	if err != nil {
		return err
	}

	return nil
}
