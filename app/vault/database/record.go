package database

import (
	"os"
	"pvault/app/vault/record"
	"pvault/app/vault/record/encoder"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func SaveRecord(db Database, r record.Record, password string) error {
	return encoder.SaveRecordFile(db, r, db.RecordPath(r.GetID()), password)
}

func LoadRecord(db Database, id uuid.UUID, password string) (record.Record, error) {
	return encoder.LoadRecordFile(db, db.RecordPath(id), password)
}

func DeleteRecord(db Database, id uuid.UUID) error {
	err := os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}

	return nil
}
