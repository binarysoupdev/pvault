package vault

import (
	"fmt"
	"os"
	"pvault/chain"
	"pvault/data"
	"pvault/vault/record"

	"github.com/google/uuid"
)

func (v Vault) SaveRecord(r record.Record) error {
	existingId, ok := v.Index[r.Name]
	if ok && existingId != r.ID {
		return fmt.Errorf("name \"%s\" already exists", r.Name)
	}

	err := data.SaveJSON(r, v.RecordPath(r.ID))
	if err != nil {
		return chain.Error(err, "error saving record file")
	}

	err = v.updateIndex(r)
	if err != nil {
		return err
	}

	return nil
}

func (v Vault) LoadRecord(name string) (record.Record, error) {
	var r record.Record

	id, ok := v.Index[name]
	if !ok {
		return r, fmt.Errorf("name \"%s\" not found", name)
	}

	r, err := data.LoadJSON[record.Record](v.RecordPath(id))
	if err != nil {
		return r, chain.Error(err, "error loading record file")
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
		return uuid.Nil, chain.Error(err, "error deleting record file")
	}

	err = v.deleteIndex(name)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
