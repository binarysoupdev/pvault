package cmd

import (
	"os"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

type LockCommand struct {
	command.FlagCommandBase

	ConfigLoader json.Loader[config.Config]
	Config       config.Config
}

func NewLockCommand(loader json.Loader[config.Config]) *LockCommand {
	return &LockCommand{
		FlagCommandBase: command.NewFlagCommandBase("lock", "lock a record in the vault"),
		ConfigLoader:    loader,
	}
}

func (cmd *LockCommand) Initialize() error {
	var err error
	cmd.Config, err = flow.LoadConfig(cmd.ConfigLoader)
	if err != nil {
		return err
	}

	return cmd.FlagCommandBase.Initialize()
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

	r, err := json.UnmarshalFile[record.Record](*path)
	if err != nil {
		return errors.Chain(err, "error loading source record")
	}

	err = flow.SaveVaultRecord(v, r)
	if err != nil {
		return err
	}

	err = os.Remove(*path)
	if err != nil {
		return errors.Chain(err, "error removing source record")
	}

	style.Delete.Printf("[-] %s\n", *path)
	return nil
}
