package vault

import (
	"fmt"
	"os"
	"pvault/chain"
	"pvault/vault/record"

	"github.com/google/uuid"
)

func (v Vault) SaveRecord(r record.Record) error {
	existingId, ok := v.Index[r.Name]
	if ok && existingId != r.ID {
		return fmt.Errorf("name \"%s\" already exists", r.Name)
	}

	err := v.saveEncryptedRecord(r)
	if err != nil {
		return err
	}

	err = v.updateIndex(r)
	if err != nil {
		return err
	}

	return nil
}

func (v Vault) LoadRecord(name string) (record.Record, error) {
	id, ok := v.Index[name]
	if !ok {
		return record.Record{}, fmt.Errorf("name \"%s\" not found", name)
	}

	r, err := v.loadEncryptedRecord(id)
	if err != nil {
		return record.Record{}, err
	}

	return r, nil
}

func (v Vault) DeleteRecord(name string) (uuid.UUID, error) {
	id, ok := v.Index[name]
	if !ok {
		return uuid.Nil, fmt.Errorf("name \"%s\" not found", name)
	}

	err := os.Remove(v.RecordPath(id))
	if err != nil {
		return id, chain.Error(err, "error deleting record file")
	}

	err = v.deleteIndex(name)
	if err != nil {
		return id, err
	}

	return id, nil
}
