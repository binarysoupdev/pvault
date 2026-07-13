package version2

import (
	"os"
	"path/filepath"
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const RECORD_VERSION uint16 = 2

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String())
}

func (db Database) SaveRecord(r record.Record, password string) error {
	bytes, err := r.Marshal(password)
	if err != nil {
		return errors.Chain(err, "error encrypting record")
	}

	err = os.WriteFile(db.RecordPath(r.GetID()), bytes, 0666)
	if err != nil {
		return errors.Chain(err, "error writing record file")
	}

	return nil
}

func (db Database) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	bytes, err := os.ReadFile(db.RecordPath(id))
	if err != nil {
		return nil, errors.Chain(err, "error reading record file")
	}

	r, err := record.UnmarshalGeneric(password, bytes, id)
	if err != nil {
		return nil, errors.Chain(err, "error unmarshaling record")
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
