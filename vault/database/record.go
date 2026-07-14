package database

import (
	"os"
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (db Database) SaveRecord(r record.Record, password string) error {
	file, err := os.Create(db.RecordPath(r.GetID()))
	if err != nil {
		return errors.Chain(err, "error creating record file")
	}
	defer file.Close()

	err = r.Encode(file, password)
	if err != nil {
		return errors.Chain(err, "error encoding record")
	}

	return nil
}

func (db Database) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	file, err := os.Open(db.RecordPath(id))
	if err != nil {
		return nil, errors.Chain(err, "error opening record file")
	}
	defer file.Close()

	r, err := record.Decode(file, password, id)
	if err != nil {
		return nil, errors.Chain(err, "error decoding record")
	}

	return r, nil
}

func (db Database) DeleteRecord(id uuid.UUID) error {
	err := os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}
	return nil
}
