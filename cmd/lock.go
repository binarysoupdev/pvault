package cmd

import (
	"os"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/data"
	"pvault/errors"
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type LockCommand struct {
	command.FlagCommandBase
	config.Loader[config.Config]
}

func NewLockCommand(loader config.Loader[config.Config]) *LockCommand {
	return &LockCommand{
		FlagCommandBase: command.NewFlagCommandBase("lock", "lock a record in the vault"),
		Loader:          loader,
	}
}

func (cmd *LockCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()
	return flow.LoadConfig(&cmd.Loader)
}

func (cmd LockCommand) Run(args []string) error {
	path := cmd.Flags.String("path", "", "path to the record JSON")
	cmd.Flags.Parse(args)

	if *path == "" {
		return errors.New("\"path\" cannot be empty")
	}

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	r, err := data.LoadJSON[record.Record](*path)
	if err != nil {
		return errors.Chain(err, "error loading source record")
	}

	err = v.ValidateRecord(r)
	if err != nil {
		return errors.Chain(err, "error validating record")
	}

	password := flow.PromptPassword("New PASSWORD: ")
	if flow.PromptPassword("Verify PASSWORD: ") != password {
		return errors.New("passwords do not match")
	}

	err = v.SaveRecord(r, password)
	if err != nil {
		return errors.Chain(err, "error saving vault record")
	}

	style.BoldCreate.Printf("[+] Updated Record: %s\n", r.ID.String())

	err = os.Remove(*path)
	if err != nil {
		return errors.Chain(err, "error removing source record")
	}

	style.Delete.Printf("[-] %s\n", *path)
	return nil
}
