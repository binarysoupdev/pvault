package vault

import (
	"pvault/app/flow/prompt"
	record_v2 "pvault/app/vault/record/record/v2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/got-style/style"
)

func SaveRecord(v Vault, r record_v2.Record) error {
	err := v.ValidateRecord(r)
	if err != nil {
		return errors.Chain(err, "error validating record")
	}

	password := prompt.Password("New PASSWORD: ")
	if prompt.Password("Verify PASSWORD: ") != password {
		return errors.New("passwords do not match")
	}

	err = v.SaveRecord(r, password)
	if err != nil {
		return errors.Chain(err, "error saving vault record")
	}

	style.BoldCreate.Printf("[+] Saved Record: %s\n", r.GetID().String())
	return nil
}

func LoadRecord(v Vault, name string) (record_v2.Record, error) {
	password := prompt.Password("Enter PASSWORD: ")

	r, err := v.LoadRecord(name, password)
	if err != nil {
		return record_v2.Record{}, errors.Chain(err, "error loading vault record")
	}

	style.BoldInfo.Printf("[=] Loaded Record: %s\n", r.GetID().String())
	return r.Upgrade(), nil
}

func DeleteRecord(v Vault, name string) error {
	if prompt.Prompt("Confirm NAME: ") != name {
		return errors.New("names do not match")
	}

	id, err := v.DeleteRecord(name)
	if err != nil {
		return errors.Chain(err, "error deleting vault record")
	}

	style.BoldDelete.Printf("[-] Deleted Record: %s\n", id.String())
	return nil
}
