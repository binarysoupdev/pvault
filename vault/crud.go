package vault

import (
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/errors"

	"github.com/google/uuid"
)

func (v Vault) SaveRecord(r record.RecordV2, password string) error {
	err := v.ValidateRecord(r)
	if err != nil {
		return errors.Chain(err, "error validating record")
	}

	err = v.Database.SaveRecord(r, password)
	if err != nil {
		return errors.Chain(err, "error saving record to database")
	}

	existingName, ok := v.Index.FindName(r.ID)
	if ok && existingName != r.Name {
		delete(v.Index, existingName)
	}
	v.Index[r.Name] = r.ID

	err = v.Database.SaveIndex(v.Index)
	if err != nil {
		return errors.Chain(err, "error saving index to database")
	}

	return nil
}

func (v Vault) LoadRecord(name string, password string) (record.RecordV2, error) {
	id, ok := v.Index[name]
	if !ok {
		return record.RecordV2{}, errors.Format("name \"%s\" not found", name)
	}

	r, err := v.Database.LoadRecord(id, password)
	if err != nil {
		return record.RecordV2{}, errors.Chain(err, "error loading record from database")
	}

	return r, nil
}

func (v Vault) DeleteRecord(name string) (uuid.UUID, error) {
	id, ok := v.Index[name]
	if !ok {
		return uuid.Nil, errors.Format("name \"%s\" not found", name)
	}

	err := v.Database.DeleteRecord(id)
	if err != nil {
		return id, errors.Chain(err, "error deleting record from database")
	}

	delete(v.Index, name)

	err = v.Database.SaveIndex(v.Index)
	if err != nil {
		return id, errors.Chain(err, "error saving index to database")
	}

	return id, nil
}
