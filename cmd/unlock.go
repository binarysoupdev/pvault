package cmd

import (
	"path/filepath"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/data"
	"pvault/errors"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type UnlockCommand struct {
	command.FlagCommandBase
	config.Loader[config.Config]
}

func NewUnlockCommand(loader config.Loader[config.Config]) *UnlockCommand {
	return &UnlockCommand{
		FlagCommandBase: command.NewFlagCommandBase("unlock", "unlock a record from the vault"),
		Loader:          loader,
	}
}

func (cmd *UnlockCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()

	err := cmd.LoadConfig()
	if err != nil {
		return errors.Chain(err, "error loading config")
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

	v, err := vault.Open(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error opening vault")
	}

	name, err := search.Select(v)
	if err != nil {
		return err
	}

	password := flow.PromptPassword("Enter PASSWORD: ")

	r, err := v.LoadRecord(name, password)
	if err != nil {
		return errors.Chain(err, "error loading vault record")
	}

	style.BoldInfo.Printf("[=] Loaded Record: %s\n", r.ID.String())

	path := filepath.Join(cmd.Config.OutputPath, r.ID.String()+".json")
	err = data.SaveJSON(r, path)
	if err != nil {
		return errors.Chain(err, "error creating output record")
	}

	style.Create.Printf("[+] %s\n", path)
	return nil
}
