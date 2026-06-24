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
	config.LoaderModule[config.Config]
}

func NewCreateCommand(loader config.Loader[config.Config]) *CreateCommand {
	return &CreateCommand{
		FlagCommandBase: command.NewFlagCommandBase("create", "create a new vault record"),
		LoaderModule:    config.NewLoaderModule(loader),
	}
}

func (cmd *CreateCommand) Initialize() error {
	_ = cmd.FlagCommandBase.Initialize()

	err := cmd.LoadConfig()
	if err != nil {
		return errors.Chain(err, "error loading config")
	}

	return nil
}

func (cmd CreateCommand) Run(args []string) error {
	name := cmd.Flags.String("name", "", "name of the record")
	cmd.Flags.Parse(args)

	if *name == "" {
		return errors.New("\"name\" cannot be empty")
	}

	v, err := vault.Open(cmd.Config.VaultPath)
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

	path := filepath.Join(cmd.Config.OutputPath, r.ID.String()+".json")
	err = data.SaveJSON(r, path)
	if err != nil {
		return errors.Chain(err, "error creating output record")
	}

	style.Create.Printf("[+] %s\n", path)
	return nil
}
