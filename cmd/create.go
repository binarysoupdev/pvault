package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"pvault/chain"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

const LOCAL_DIR = "tmp/local"

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
		return err
	}
	style.BoldCreate.Printf("[+] New Record: %s\n", r.ID.String())

	filename := filepath.Join(LOCAL_DIR, r.ID.String()+".json")
	file, err := os.Create(filename)
	if err != nil {
		return chain.Error(err, "error creating record file")
	}

	e := json.NewEncoder(file)
	e.SetIndent("", "    ")

	err = e.Encode(r)
	if err != nil {
		return chain.Error(err, "error encode record JSON")
	}

	style.Create.Printf("[+] %s\n", filename)
	return nil
}
