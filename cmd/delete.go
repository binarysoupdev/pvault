package cmd

import (
	"pvault/chain"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type DeleteCommand struct {
	command.FlagCommandBase
}

func NewDeleteCommand() *DeleteCommand {
	return &DeleteCommand{
		FlagCommandBase: command.NewFlagCommandBase("delete", "delete a record from the vault"),
	}
}

func (cmd DeleteCommand) Run(args []string) error {
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

	if flow.Prompt("Confirm NAME: ") != name {
		return chain.New("names do not match")
	}

	id, err := v.DeleteRecord(name)
	if err != nil {
		return chain.Error(err, "error deleting vault record")
	}

	style.BoldInfo.Printf("[-] Deleted Record: %s\n", id.String())
	return nil
}
