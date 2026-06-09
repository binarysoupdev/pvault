package cmd

import (
	"encoding/json"
	"os"
	"pvault/chain"
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

	file, err := os.Open(*path)
	if err != nil {
		return chain.Error(err, "error opening local record")
	}
	defer file.Close()

	var r vault.Record
	err = json.NewDecoder(file).Decode(&r)
	if err != nil {
		return chain.Error(err, "error decoding record JSON")
	}

	err = vault.Vault{}.SaveRecord(r)
	if err != nil {
		return chain.Error(err, "error saving record")
	}
	style.BoldInfo.Printf("[+] Updated Record: %s\n", r.ID.String())

	err = os.Remove(*path)
	if err != nil {
		return chain.Error(err, "error removing local record")
	}

	style.Delete.Printf("[-] %s\n", *path)
	return nil
}
