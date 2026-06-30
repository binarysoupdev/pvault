package version1

import (
	"encoding/binary"
	"path/filepath"
	"pvault/vault/data"
	"pvault/vault/record"
	"pvault/vault/record/legacy"

	"github.com/google/uuid"
)

const RECORD_VERSION uint16 = 1

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String())
}

func (db Database) SaveRecord(r record.Record, password string) error {
	old := legacy.RecordV1{
		Password: r.Password,
		Username: r.Username,
	}

	header := make([]byte, 2+len(r.Name))
	binary.BigEndian.PutUint16(header, uint16(len(r.Name)))
	copy(header[2:], []byte(r.Name))

	return data.SaveEncryptedRecord(db.RecordPath(r.ID), password, RECORD_VERSION, header, old)
}

func (Database) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	return record.Record{}, data.NotSupportedError{}
}

func (Database) DeleteRecord(id uuid.UUID) error {
	return data.NotSupportedError{}
}
