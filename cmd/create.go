package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
)

type CreateCommand struct {
	command.CommandBase
	command.FlagCommand
	flow.ConfigCommand
}

func NewCreateCommand(configLoader json.Loader[config.Config]) *CreateCommand {
	return &CreateCommand{
		CommandBase:   command.NewCommandBase("create", "create a new vault record"),
		ConfigCommand: flow.NewConfigCommand(configLoader),
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
	cmd.Flags.Parse(args)

	if *name == "" {
		return errors.New("\"name\" cannot be empty")
	}

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	err = cmd.Config.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating config \"output_path\"")
	}

	r := record.NewFromName(*name)

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
