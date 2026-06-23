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
}

func NewInitCommand() *InitCommand {
	return &InitCommand{
		FlagCommandBase: command.NewFlagCommandBase("init", "initialize a new vault"),
	}
}

func (cmd InitCommand) Run(_ []string) error {
	_, err := vault.InitializeNew(config.Global.VaultPath)
	if err != nil {
		return errors.Chain(err, "error initializing new vault")
	}

	style.BoldCreate.Printf("[+] New Vault Initialized: %s\n", config.Global.VaultPath)

	return nil
}
