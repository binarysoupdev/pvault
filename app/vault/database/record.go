package database

import (
	"os"
	"pvault/app/vault/record"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func SaveRecord(db Encoder, path string, r record.Record, password string) error {
	file, err := os.Create(db.RecordPath(path, r.GetID()))
	if err != nil {
		return errors.Chain(err, "error creating record file")
	}
	defer file.Close()

	err = db.EncodeRecord(file, password, r)
	if err != nil {
		return errors.Chain(err, "error encoding record")
	}

	return nil
}

func LoadRecord(db Encoder, path string, id uuid.UUID, password string) (record.Record, error) {
	file, err := os.Open(db.RecordPath(path, id))
	if err != nil {
		return nil, errors.Chain(err, "error opening record file")
	}
	defer file.Close()

	r, err := db.DecodeRecord(file, password)
	if err != nil {
		return nil, errors.Chain(err, "error decoding record")
	}

	return r, nil
}

func DeleteRecord(db Encoder, path string, id uuid.UUID) error {
	err := os.Remove(db.RecordPath(path, id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}

	return nil
}
