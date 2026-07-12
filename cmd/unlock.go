package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
)

type UnlockCommand struct {
	command.FlagCommandBase

	ConfigLoader json.Loader[config.Config]
	Config       config.Config
}

func NewUnlockCommand(loader json.Loader[config.Config]) *UnlockCommand {
	return &UnlockCommand{
		FlagCommandBase: command.NewFlagCommandBase("unlock", "unlock a record from the vault"),
		ConfigLoader:    loader,
	}
}

func (cmd *UnlockCommand) Initialize() error {
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

func (cmd UnlockCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.Flags.Parse(args)

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return err
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
