package cmd

import (
	"pvault/config"
	"pvault/errors"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type InitCommand struct {
	command.FlagCommandBase
	config.LoaderModule[config.Config]
}

func NewInitCommand(loader config.Loader[config.Config]) *InitCommand {
	return &InitCommand{
		FlagCommandBase: command.NewFlagCommandBase("init", "initialize a new vault"),
		LoaderModule:    config.NewLoaderModule(loader),
	}
}

func (cmd *InitCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()

	err := cmd.LoadConfig()
	if err != nil {
		return errors.Chain(err, "error loading config")
	}

	return nil
}

func (cmd InitCommand) Run(_ []string) error {
	_, err := vault.InitializeNew(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error initializing new vault")
	}

	style.BoldCreate.Printf("[+] New Vault Initialized: %s\n", cmd.Config.VaultPath)

	return nil
}
