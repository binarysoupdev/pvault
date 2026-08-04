package cmd

import (
	"pvault/app/cmd/base"
	"pvault/app/config"
	"pvault/app/flow"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
)

type UnlockCommand struct {
	command.CommandBase
	command.FlagCommand
	base.ConfigCommand
}

func NewUnlockCommand(configLoader json.Loader[config.Config]) *UnlockCommand {
	return &UnlockCommand{
		CommandBase:   command.NewCommandBase("unlock", "unlock a record from the vault"),
		ConfigCommand: base.NewConfigCommand(configLoader),
	}
}

func (cmd *UnlockCommand) Initialize() error {
	if err := cmd.LoadConfig(); err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd UnlockCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.ParseFlags(args)

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	err = cmd.Config.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating config \"output_path\"")
	}

	name, err := search.Select(v)
	if err != nil {
		return err
	}

	r, err := flow.LoadVaultRecord(v, name)
	if err != nil {
		return err
	}

	err = flow.SaveOutputRecord(cmd.Config, r)
	if err != nil {
		return err
	}

	return nil
}
