package vault

import (
	"os"
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (v Vault) ValidateRecord(r record.Record) error {
	err := r.Validate()
	if err != nil {
		return err
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

	err = r.SaveFile(v.Index.RecordPath(r.GetID()), password)
	if err != nil {
		return errors.Chain(err, "error encoding record")
	}

	existingName, ok := v.Map.FindName(r.GetID())
	if ok && existingName != r.GetName() {
		delete(v.Map, existingName)
	}
	v.Map[r.GetName()] = r.GetID()

	err = v.Index.SaveMap(v.Map)
	if err != nil {
		return errors.Chain(err, "error saving index to database")
	}

	return nil
}

func (v Vault) LoadRecord(name string, password string) (record.Record, error) {
	id, ok := v.Map[name]
	if !ok {
		return nil, errors.Format("name \"%s\" not found", name)
	}

	r, err := record.Load(v.Index.RecordPath(id), password, id)
	if err != nil {
		return nil, errors.Chain(err, "error decoding record")
	}

	return r, nil
}

func (v Vault) DeleteRecord(name string) (uuid.UUID, error) {
	id, ok := v.Map[name]
	if !ok {
		return uuid.Nil, errors.Format("name \"%s\" not found", name)
	}

	err := os.Remove(v.Index.RecordPath(id))
	if err != nil {
		return uuid.Nil, errors.Chain(err, "error deleting record file")
	}

	delete(v.Map, name)

	err = v.Index.SaveMap(v.Map)
	if err != nil {
		return id, errors.Chain(err, "error saving index to database")
	}

	return id, nil
}
