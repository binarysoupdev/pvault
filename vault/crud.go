package vault

import (
	"os"
	"pvault/errors"
	"pvault/vault/record"

	"github.com/google/uuid"
)

func (v Vault) SaveRecord(r record.Record, password string) error {
	err := v.ValidateRecord(r)
	if err != nil {
		return errors.Chain(err, "error validating record")
	}

	err = v.saveEncryptedRecord(r, password)
	if err != nil {
		return err
	}

	err = v.updateIndex(r)
	if err != nil {
		return err
	}

	return nil
}

func (v Vault) LoadRecord(name string, password string) (record.Record, error) {
	id, ok := v.Index[name]
	if !ok {
		return record.Record{}, errors.Format("name \"%s\" not found", name)
	}

	r, err := v.loadEncryptedRecord(id, password)
	if err != nil {
		return record.Record{}, err
	}

	return r, nil
}

func (v Vault) DeleteRecord(name string) (uuid.UUID, error) {
	id, ok := v.Index[name]
	if !ok {
		return uuid.Nil, errors.Format("name \"%s\" not found", name)
	}

	err := os.Remove(v.RecordPath(id))
	if err != nil {
		return id, errors.Chain(err, "error deleting record file")
	}

	err = v.deleteIndex(name)
	if err != nil {
		return id, err
	}

	return id, nil
}
