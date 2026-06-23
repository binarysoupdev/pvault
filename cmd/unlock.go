package cmd

import (
	"path/filepath"
	"pvault/chain"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/data"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type UnlockCommand struct {
	command.FlagCommandBase
}

func NewUnlockCommand() *UnlockCommand {
	return &UnlockCommand{
		FlagCommandBase: command.NewFlagCommandBase("unlock", "unlock a record from the vault"),
	}
}

func (cmd UnlockCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.Flags.Parse(args)

	v, err := vault.Open(config.Global.VaultPath)
	if err != nil {
		return chain.Error(err, "error opening vault")
	}

	name, err := search.Select(v)
	if err != nil {
		return err
	}

	password := flow.PromptPassword("Enter PASSWORD: ")

	r, err := v.LoadRecord(name, password)
	if err != nil {
		return chain.Error(err, "error loading vault record")
	}

	style.BoldInfo.Printf("[=] Loaded Record: %s\n", r.ID.String())

	path := filepath.Join(config.Global.OutputPath, r.ID.String()+".json")
	err = data.SaveJSON(r, path)
	if err != nil {
		return chain.Error(err, "error creating output record")
	}

	style.Create.Printf("[+] %s\n", path)
	return nil
}
