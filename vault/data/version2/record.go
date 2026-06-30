package version2

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/data/version1"
	"pvault/vault/record"

	"github.com/google/uuid"
)

const RECORD_VERSION uint16 = 2

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String())
}

func (db Database) SaveRecord(r record.Record, password string) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, RECORD_VERSION)

	return data.SaveEncryptedRecord(db.RecordPath(r.ID), password, header, r)
}

func (db Database) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	raw, err := os.ReadFile(db.RecordPath(id))
	if err != nil {
		return record.Record{}, errors.Chain(err, "error reading record file")
	}

	version := binary.BigEndian.Uint16(raw)
	raw = raw[2:]

	switch version {
	case 1:
		return version1.NewDatabase(db.Path).ParseRecordV1(id, password, raw)
	case 2:
		return data.DecryptRecord[record.Record](password, raw)
	default:
		return record.Record{}, errors.Format("unsupported record version \"%d\"", version)
	}
}

func (db Database) DeleteRecord(id uuid.UUID) error {
	err := os.Remove(db.RecordPath(id))
	if err != nil {
		return errors.Chain(err, "error removing record file")
	}
	return nil
}
