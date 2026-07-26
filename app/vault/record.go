package vault

import (
	"pvault/app/vault/database"
	"pvault/app/vault/record"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (v Vault) ValidateRecord(r record.Record) error {
	err := record.Validate(r)
	if err != nil {
		return errors.Chain(err, "record invalid")
	}

	existingId, ok := v.Map[r.GetName()]
	if ok && existingId != r.GetID() {
		return errors.Format("name \"%s\" already exists", r.GetName())
	}

	return nil
}

func (v Vault) SaveRecord(r record.Record, password string) error {
	err := v.ValidateRecord(r)
	if err != nil {
		return errors.Chain(err, "error validating record")
	}

	err = database.SaveRecord(v.Database, v.Path, r, password)
	if err != nil {
		return errors.Chain(err, "error saving record")
	}

	existingName, ok := v.Map.FindName(r.GetID())
	if ok && existingName != r.GetName() {
		delete(v.Map, existingName)
	}
	v.Map[r.GetName()] = r.GetID()

	err = database.SaveIndex(v.Database, v.Path, v.Map)
	if err != nil {
		return errors.Chain(err, "error saving index")
	}

	return nil
}

func (v Vault) LoadRecord(name string, password string) (record.Record, error) {
	id, ok := v.Map[name]
	if !ok {
		return nil, errors.Format("name \"%s\" not found", name)
	}

	r, err := database.LoadRecord(v.Database, v.Path, id, password)
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

	err := database.DeleteRecord(v.Database, v.Path, id)
	if err != nil {
		return uuid.Nil, errors.Chain(err, "error deleting record")
	}

	delete(v.Map, name)

	err = database.SaveIndex(v.Database, v.Path, v.Map)
	if err != nil {
		return id, errors.Chain(err, "error saving index map")
	}

	return id, nil
}
