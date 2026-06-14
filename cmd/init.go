package cmd

import (
	"pvault/chain"
	"pvault/config"
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
	err := vault.InitializeNew(config.Global.VaultPath)
	if err != nil {
		return chain.Error(err, "error initializing new vault")
	}

	style.BoldCreate.Printf("[+] New Vault Initialized: %s\n", config.Global.VaultPath)

	return nil
}
