package flow

import (
	"pvault/app/vault/record"
	v2 "pvault/app/vault/record/version2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/binarysoupdev/got-style/style"
	"github.com/google/uuid"
)

type Vault interface {
	GetVersion() int

	SearchNames(term string) []string

	ValidateRecord(r record.Record) error
	SaveRecord(r record.Record, password string) error
	LoadRecord(name string, password string) (record.Record, error)
	DeleteRecord(name string) (uuid.UUID, error)
}

func SaveVaultRecord(v Vault, r v2.Record) error {
	err := v.ValidateRecord(r)
	if err != nil {
		return errors.Chain(err, "error validating record")
	}

	password := promptPassword("New PASSWORD: ")
	if promptPassword("Verify PASSWORD: ") != password {
		return errors.New("passwords do not match")
	}

	err = v.SaveRecord(r, password)
	if err != nil {
		return errors.Chain(err, "error saving vault record")
	}

	style.BoldCreate.Printf("[+] Saved Record: %s\n", r.GetID().String())
	return nil
}

func LoadVaultRecord(v Vault, name string) (v2.Record, error) {
	password := promptPassword("Enter PASSWORD: ")

	r, err := v.LoadRecord(name, password)
	if err != nil {
		return v2.Record{}, errors.Chain(err, "error loading vault record")
	}

	style.BoldInfo.Printf("[=] Loaded Record: %s\n", r.GetID().String())
	return r.Upgrade(), nil
}

func DeleteVaultRecord(v Vault, name string) error {
	if prompt("Confirm NAME: ") != name {
		return errors.New("names do not match")
	}

	id, err := v.DeleteRecord(name)
	if err != nil {
		return errors.Chain(err, "error deleting vault record")
	}

	style.BoldDelete.Printf("[-] Deleted Record: %s\n", id.String())
	return nil
}
