package v2

import (
	"encoding/binary"
	"os"
	"path/filepath"
	v1 "pvault/vault/record/version/v1"
	v2 "pvault/vault/record/version/v2"

	"github.com/binarysoupdev/go-commando/errors"

	"github.com/google/uuid"
)

const RECORD_VERSION uint16 = 2

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String())
}

func (db Database) SaveRecord(r v2.Record, password string) error {
	bytes, err := r.Marshal(password)
	if err != nil {
		return errors.Chain(err, "error encrypting record")
	}

	err = os.WriteFile(db.RecordPath(r.ID), bytes, 0666)
	if err != nil {
		return errors.Chain(err, "error writing record file")
	}

	return nil
}

func (db Database) LoadRecord(id uuid.UUID, password string) (v2.Record, error) {
	bytes, err := os.ReadFile(db.RecordPath(id))
	if err != nil {
		return v2.Record{}, errors.Chain(err, "error reading record file")
	}

	version := binary.BigEndian.Uint16(bytes)

	switch version {
	case 1:
		return v1.Unmarshal(password, bytes, id)
	case 2:
		return v2.Unmarshal(password, bytes)
	default:
		return v2.Record{}, errors.Format("unsupported record version \"%d\"", version)
	}
}

func (db Database) DeleteRecord(id uuid.UUID) error {
	err := os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}
	return nil
}
