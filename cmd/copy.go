package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/errors"
	"pvault/vault"

	"github.com/atotto/clipboard"
	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type CopyCommand struct {
	command.FlagCommandBase
}

func NewCopyCommand() *CopyCommand {
	return &CopyCommand{
		FlagCommandBase: command.NewFlagCommandBase("copy", "copy password/username of a record"),
	}
}

func (cmd CopyCommand) Run(args []string) error {
	search := flow.NewSearchFlow(cmd.Flags)
	cmd.Flags.Parse(args)

	v, err := vault.Open(config.Global.VaultPath)
	if err != nil {
		return errors.Chain(err, "error opening vault")
	}

	name, err := search.Select(v)
	if err != nil {
		return err
	}

	password := flow.PromptPassword("Enter PASSWORD: ")

	r, err := v.LoadRecord(name, password)
	if err != nil {
		return errors.Chain(err, "error loading vault record")
	}

	style.BoldInfo.Printf("[=] Loaded Record: %s\n", r.ID.String())

	err = clipboard.WriteAll(r.Password)
	if err != nil {
		return errors.Chain(err, "error copying to clipboard")
	}

	style.Info.Printf("[=] PASSWORD copied to clipboard\n")

	return nil
}
