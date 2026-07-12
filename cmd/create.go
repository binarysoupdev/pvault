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
	command.FlagCommandBase

	ConfigLoader json.Loader[config.Config]
	Config       config.Config
}

func NewCreateCommand(loader json.Loader[config.Config]) *CreateCommand {
	return &CreateCommand{
		FlagCommandBase: command.NewFlagCommandBase("create", "create a new vault record"),
		ConfigLoader:    loader,
	}
}

func (cmd *CreateCommand) Initialize() error {
	var err error
	cmd.Config, err = flow.LoadConfig(cmd.ConfigLoader)
	if err != nil {
		return err
	}

	err = cmd.Config.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating config \"output_path\"")
	}

	return cmd.FlagCommandBase.Initialize()
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
