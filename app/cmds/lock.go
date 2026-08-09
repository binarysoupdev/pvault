package cmds

import (
	"os"
	"pvault/app/config"
	config_flow "pvault/app/flow/config"
	vault_flow "pvault/app/flow/vault"
	record_v2 "pvault/vault/record/record/v2"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/go-commando/json"
	"github.com/binarysoupdev/got-style/style"
)

type LockCommand struct {
	command.CommandBase
	command.FlagCommand
	command.ConfigCommand[config.Config]
}

func NewLockCommand(configLoader json.Loader[config.Config]) *LockCommand {
	return &LockCommand{
		CommandBase:   command.NewCommandBase("lock", "lock a record in the vault"),
		ConfigCommand: command.NewConfigCommand(configLoader),
	}
}

func (cmd *LockCommand) Initialize() error {
	err := config_flow.LoadConfig(&cmd.ConfigCommand)
	if err != nil {
		return err
	}

	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd LockCommand) Run(args []string) error {
	path := cmd.Flags.String("path", "", "path to the record JSON")
	cmd.ParseFlags(args)

	v, err := vault_flow.OpenCurrentVault(cmd.Config.VaultPath)
	if err != nil {
		return err
	}

	if *path == "" {
		return errors.New("\"path\" cannot be empty")
	}

	r, err := json.UnmarshalFile[record_v2.Record](*path)
	if err != nil {
		return errors.Chain(err, "error loading source record")
	}

	err = vault_flow.SaveRecord(v, r)
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
