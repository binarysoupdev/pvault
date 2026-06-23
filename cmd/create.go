package cmd

import (
	"path/filepath"
	"pvault/cmd/flow"
	"pvault/config"
	"pvault/data"
	"pvault/errors"
	"pvault/vault"
	"pvault/vault/record"

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
		return errors.New("\"name\" cannot be empty")
	}

	v, err := vault.Open(config.Global.VaultPath)
	if err != nil {
		return errors.Chain(err, "error opening vault")
	}

	r := record.NewFromName(*name)

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

	style.BoldCreate.Printf("[+] New Record: %s\n", r.ID.String())

	path := filepath.Join(config.Global.OutputPath, r.ID.String()+".json")
	err = data.SaveJSON(r, path)
	if err != nil {
		return errors.Chain(err, "error creating output record")
	}

	style.Create.Printf("[+] %s\n", path)
	return nil
}
