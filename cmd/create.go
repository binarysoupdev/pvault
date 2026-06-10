package cmd

import (
	"path/filepath"
	"pvault/cfg"
	"pvault/chain"
	"pvault/data"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

type CreateCommand struct {
	command.FlagCommandBase
}

func NewCreateCommand() *CreateCommand {
	return &CreateCommand{
		FlagCommandBase: command.NewFlagCommandBase("create", "create a new vault record"),
	}
}

func (cmd CreateCommand) Run(args []string) error {
	name := cmd.Flags.String("name", "", "name of the record")
	cmd.Flags.Parse(args)

	if *name == "" {
		return chain.New("\"name\" cannot be empty")
	}

	r, err := vault.Vault{}.NewRecord(*name)
	if err != nil {
		return chain.Error(err, "error creating vault record")
	}

	style.BoldCreate.Printf("[+] New Record: %s\n", r.ID.String())

	path := filepath.Join(cfg.Global.OutputPath, r.ID.String()+".json")
	err = data.SaveJSON(r, path)
	if err != nil {
		return chain.Error(err, "error creating output record")
	}

	style.Create.Printf("[+] %s\n", path)
	return nil
}
