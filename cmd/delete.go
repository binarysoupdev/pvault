package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/errors"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type DeleteCommand struct {
	command.FlagCommandBase
	config.LoaderModule[config.Config]
}

func NewDeleteCommand(loader config.Loader[config.Config]) *DeleteCommand {
	return &DeleteCommand{
		FlagCommandBase: command.NewFlagCommandBase("delete", "delete a record from the vault"),
		LoaderModule:    config.NewLoaderModule(loader),
	}
}

func (cmd *DeleteCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()

	err := cmd.LoadConfig()
	if err != nil {
		return errors.Chain(err, "error loading config")
	}

	return nil
}

func (cmd DeleteCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.Flags.Parse(args)

	v, err := vault.Open(cmd.Config.VaultPath)
	if err != nil {
		return errors.Chain(err, "error opening vault")
	}

	name, err := search.Select(v)
	if err != nil {
		return err
	}

	if flow.Prompt("Confirm NAME: ") != name {
		return errors.New("names do not match")
	}

	id, err := v.DeleteRecord(name)
	if err != nil {
		return errors.Chain(err, "error deleting vault record")
	}

	style.BoldDelete.Printf("[-] Deleted Record: %s\n", id.String())
	return nil
}
