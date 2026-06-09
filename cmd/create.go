package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"pvault/chain"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
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
	r := vault.NewRecord("foobar")

	file, err := os.Create(filepath.Join("tmp", r.ID.String()+".json"))
	if err != nil {
		return chain.Error(err, "error creating record file")
	}

	err = json.NewEncoder(file).Encode(r)
	if err != nil {
		return chain.Error(err, "error encode record JSON")
	}

	return nil
}
