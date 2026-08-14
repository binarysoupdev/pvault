package vault

import (
	"pvault/vault/database"
	"pvault/vault/record"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/google/uuid"
)

func (v Vault) RecordPath(id uuid.UUID) string {
	return v.DatabaseEncoder.RecordPath(v.Path, id)
}

func (v Vault) ValidateRecord(r record.Record) error {
	err := record.Validate(r)
	if err != nil {
		return err
	}

	existingId, ok := v.Map[r.GetName()]
	if ok && existingId != r.GetID() {
		return errors.Format("name \"%s\" already exists for a different record", r.GetName())
	}

	return nil
}

func (v Vault) SaveRecord(r record.Record, password string) error {
	err := v.ValidateRecord(r)
	if err != nil {
		return errors.Chain(err, "error validating record")
	}

	err = database.SaveRecord(v.DatabaseEncoder, v.Path, r, password)
	if err != nil {
		return errors.Chain(err, "error saving record")
	}

	existingName, ok := v.Map.FindName(r.GetID())
	if ok && existingName != r.GetName() {
		delete(v.Map, existingName)
	}
	v.Map[r.GetName()] = r.GetID()

	err = v.SaveIndex()
	if err != nil {
		return err
	}

	return nil
}

func (v Vault) LoadRecord(name string, password string) (record.Record, error) {
	id, ok := v.Map[name]
	if !ok {
		return nil, errors.Format("name \"%s\" not found", name)
	}

	r, err := database.LoadRecord(v.DatabaseEncoder, v.Path, id, password)
	if err != nil {
		return nil, errors.Chain(err, "error loading record")
	}

	return r, nil
}

func (v Vault) DeleteRecord(name string) (uuid.UUID, error) {
	id, ok := v.Map[name]
	if !ok {
		return uuid.Nil, errors.Format("name \"%s\" not found", name)
	}

	err := database.DeleteRecord(v.DatabaseEncoder, v.Path, id)
	if err != nil {
		return uuid.Nil, errors.Chain(err, "error deleting record")
	}

	delete(v.Map, name)

	err = v.SaveIndex()
	if err != nil {
		return id, err
	}

	return id, nil
}
