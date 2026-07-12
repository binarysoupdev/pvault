package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"

	"github.com/binarysoupdev/go-commando/command"
	cfg "github.com/binarysoupdev/go-commando/config"
	"github.com/binarysoupdev/go-commando/errors"
)

type UnlockCommand struct {
	command.FlagCommandBase
	cfg.Loader[config.Config]
}

func NewUnlockCommand(loader cfg.Loader[config.Config]) *UnlockCommand {
	return &UnlockCommand{
		FlagCommandBase: command.NewFlagCommandBase("unlock", "unlock a record from the vault"),
		Loader:          loader,
	}
}

func (cmd *UnlockCommand) Initialize() error {
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
