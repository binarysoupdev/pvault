package vault

import (
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) ValidateRecord(r record.Record) error {
	err := r.Validate()
	if err != nil {
		return err
	}

	existingId, ok := v.Index[r.GetName()]
	if ok && existingId != r.GetID() {
		return errors.Format("name \"%s\" already exists", r.GetName())
	}

	return nil
}
