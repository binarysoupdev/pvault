package vault

import (
	v2 "pvault/vault/record/version2"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) ValidateRecord(r v2.Record) error {
	err := r.Validate()
	if err != nil {
		return err
	}

	existingId, ok := v.Index[r.Name]
	if ok && existingId != r.ID {
		return errors.Format("name \"%s\" already exists", r.Name)
	}

	return nil
}
