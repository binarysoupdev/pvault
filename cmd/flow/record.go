package flow

import (
	"path/filepath"
	"pvault/config"
	"pvault/errors"
	"pvault/json"
	"pvault/vault"
	"pvault/vault/record"

	"github.com/binarysoupdev/got-style/style"
)

func SaveVaultRecord(v vault.Vault, r record.Record) error {
	err := v.ValidateRecord(r)
	if err != nil {
		return errors.Chain(err, "error validating record")
	}

	password := PromptPassword("New PASSWORD: ")
	if PromptPassword("Verify PASSWORD: ") != password {
		return errors.New("passwords do not match")
	}

	err = v.SaveRecord(r, password)
	if err != nil {
		return errors.Chain(err, "error saving vault record")
	}

	style.BoldCreate.Printf("[+] Saved Record: %s\n", r.ID.String())
	return nil
}

func LoadVaultRecord(v vault.Vault, name string) (record.Record, error) {
	password := PromptPassword("Enter PASSWORD: ")

	r, err := v.LoadRecord(name, password)
	if err != nil {
		return record.Record{}, errors.Chain(err, "error loading vault record")
	}

	style.BoldInfo.Printf("[=] Loaded Record: %s\n", r.ID.String())
	return r, nil
}

func DeleteVaultRecord(v vault.Vault, name string) error {
	if Prompt("Confirm NAME: ") != name {
		return errors.New("names do not match")
	}

	id, err := v.DeleteRecord(name)
	if err != nil {
		return errors.Chain(err, "error deleting vault record")
	}

	style.BoldDelete.Printf("[-] Deleted Record: %s\n", id.String())
	return nil
}

func SaveOutputRecord(cfg config.Config, r record.Record) error {
	err := cfg.ValidateOutputPath()
	if err != nil {
		return errors.Chain(err, "error validating output path")
	}

	path := filepath.Join(cfg.OutputPath, r.ID.String()+".json")

	err = json.MarshalFilePretty(r, path, "    ")
	if err != nil {
		return errors.Chain(err, "error creating output record")
	}

	style.Create.Printf("[+] %s\n", path)
	return nil
}
