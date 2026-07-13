package vault

import (
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/errors"
)

func (v Vault) ValidateRecord(r record.RecordV2) error {
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
