package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"pvault/chain"
	"pvault/vault"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
	"github.com/google/uuid"
)

type UnlockCommand struct {
	command.FlagCommandBase
}

func NewUnlockCommand() *UnlockCommand {
	return &UnlockCommand{
		FlagCommandBase: command.NewFlagCommandBase("unlock", "unlock a record from the vault"),
	}
}

func (cmd UnlockCommand) Run(args []string) error {
	name := cmd.Flags.String("name", "", "the name of the record")
	cmd.Flags.Parse(args)

	if *name == "" {
		return chain.New("\"name\" cannot be empty")
	}

	id, err := uuid.Parse(*name)
	if err != nil {
		return chain.Error(err, "error parsing ID")
	}

	r, err := vault.Vault{}.LoadRecord(id)
	if err != nil {
		return chain.Error(err, "error loading record")
	}

	style.BoldInfo.Printf("[=] Loaded Record: %s\n", r.ID.String())

	filename := filepath.Join(LOCAL_DIR, r.ID.String()+".json")

	file, err := os.Create(filename)
	if err != nil {
		return chain.Error(err, "error creating record file")
	}
	defer file.Close()

	e := json.NewEncoder(file)
	e.SetIndent("", "    ")

	err = e.Encode(r)
	if err != nil {
		return chain.Error(err, "error encode record JSON")
	}

	style.Create.Printf("[+] %s\n", filename)
	return nil
}
