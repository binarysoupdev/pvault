package cmd

import (
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/errors"
	"pvault/tools/clipboard"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type CopyCommand struct {
	command.FlagCommandBase
	clipboard clipboard.Clipboard
}

func NewCopyCommand(clipboard clipboard.Clipboard) *CopyCommand {
	return &CopyCommand{
		FlagCommandBase: command.NewFlagCommandBase("copy", "copy password/username of a record"),
		clipboard:       clipboard,
	}
}

func (cmd CopyCommand) resolveDependencies() error {
	err := cmd.clipboard.CheckUnsupported()
	if err != nil {
		return errors.Chain(err, "clipboard unsupported")
	}

	return nil
}

func (cmd CopyCommand) Run(args []string) error {
	err := cmd.resolveDependencies()
	if err != nil {
		return err
	}

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

	err = cmd.clipboard.Write(r.Password)
	if err != nil {
		return errors.Chain(err, "error copying to clipboard")
	}

	style.Info.Printf("[=] PASSWORD copied to clipboard\n")

	return nil
}
