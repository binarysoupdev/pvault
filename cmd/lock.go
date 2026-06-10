package cmd

import (
	"os"
	"pvault/chain"
	"pvault/data"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type LockCommand struct {
	command.FlagCommandBase
}

func NewLockCommand() *LockCommand {
	return &LockCommand{
		FlagCommandBase: command.NewFlagCommandBase("lock", "lock a record in the vault"),
	}
}

func (cmd LockCommand) Run(args []string) error {
	path := cmd.Flags.String("path", "", "path to the record JSON")
	cmd.Flags.Parse(args)

	if *path == "" {
		return chain.New("\"path\" cannot be empty")
	}

	r, err := data.LoadJSON[vault.Record](*path)
	if err != nil {
		return chain.Error(err, "error loading source record")
	}

	err = vault.Vault{}.SaveRecord(r)
	if err != nil {
		return chain.Error(err, "error saving vault record")
	}

	style.BoldInfo.Printf("[+] Updated Record: %s\n", r.ID.String())

	err = os.Remove(*path)
	if err != nil {
		return chain.Error(err, "error removing source record")
	}

	style.Delete.Printf("[-] %s\n", *path)
	return nil
}
