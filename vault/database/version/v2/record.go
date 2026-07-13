package v2

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"pvault/vault/database"
	v1 "pvault/vault/database/version/v1"
	"pvault/vault/record"

	"github.com/binarysoupdev/go-commando/errors"

	"github.com/google/uuid"
)

const RECORD_VERSION uint16 = 2

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String())
}

func (db Database) SaveRecord(r record.RecordV2, password string) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, RECORD_VERSION)

	return database.SaveEncryptedRecord(db.RecordPath(r.ID), password, header, r)
}

func (db Database) LoadRecord(id uuid.UUID, password string) (record.RecordV2, error) {
	raw, err := os.ReadFile(db.RecordPath(id))
	if err != nil {
		return record.RecordV2{}, errors.Chain(err, "error reading record file")
	}

	version := binary.BigEndian.Uint16(raw)
	raw = raw[2:]

	switch version {
	case 1:
		return v1.New(db.Path).ParseRecordV1(id, password, raw)
	case 2:
		return database.DecryptRecord[record.RecordV2](password, raw)
	default:
		return record.RecordV2{}, errors.Format("unsupported record version \"%d\"", version)
	}
}

func (db Database) DeleteRecord(id uuid.UUID) error {
	err := os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}
	return nil
}
