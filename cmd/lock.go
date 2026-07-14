package cmd

import (
	"os"
	"pvault/cmd/flow"
	"pvault/config"
	v2 "pvault/vault/record/version2"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

type LockCommand struct {
	command.CommandBase
	command.FlagCommand
	flow.ConfigCommand
}

func NewLockCommand(configLoader json.Loader[config.Config]) *LockCommand {
	return &LockCommand{
		CommandBase:   command.NewCommandBase("lock", "lock a record in the vault"),
		ConfigCommand: flow.NewConfigCommand(configLoader),
	}
}

func (cmd *LockCommand) Initialize() error {
	if err := cmd.LoadConfig(); err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd LockCommand) Run(args []string) error {
	path := cmd.Flags.String("path", "", "path to the record JSON")
	cmd.ParseFlags(args)

	if *path == "" {
		return errors.New("\"path\" cannot be empty")
	}

	v, err := flow.OpenVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	r, err := json.UnmarshalFile[v2.Record](*path)
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
