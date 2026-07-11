package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/errors"

	"github.com/binarysoupdev/go-commando/command"
)

type CreateCommand struct {
	command.FlagCommandBase
	config.Loader[config.Config]
}

func NewCreateCommand(loader config.Loader[config.Config]) *CreateCommand {
	return &CreateCommand{
		FlagCommandBase: command.NewFlagCommandBase("create", "create a new vault record"),
		Loader:          loader,
	}
}

func (cmd *CreateCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()

	err := flow.LoadConfig(&cmd.Loader)
	if err != nil {
		return err
	}

	err = cmd.Config.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating config \"output_path\"")
	}

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
