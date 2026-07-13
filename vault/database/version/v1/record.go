package v1

import (
	"encoding/binary"
	"path/filepath"
	"pvault/vault/database"
	v2 "pvault/vault/record/version/v2"

	"github.com/google/uuid"
)

const (
	RECORD_VERSION   uint16 = 1
	LEGACY_HASH_SIZE uint16 = 60
)

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String()+".crypt")
}

func (db Database) SaveRecord(r v2.Record, password string) error {
	return database.NotSupportedError{}
}

func (Database) LoadRecord(id uuid.UUID, password string) (v2.Record, error) {
	return v2.Record{}, database.NotSupportedError{}
}

func (Database) DeleteRecord(id uuid.UUID) error {
	return database.NotSupportedError{}
}

func (db Database) SaveRecordV1(path string, r v2.Record, password string) error {
	v1 := record.RecordV1{
		Password: r.Password,
		Username: r.Username,
	}
	header := db.buildRecordV1Header(r.Name)

	return database.SaveEncryptedRecord(path, password, header, v1)
}

func (Database) buildRecordV1Header(name string) []byte {
	header := make([]byte, 2+2+len(name))

	binary.BigEndian.PutUint16(header, RECORD_VERSION)
	binary.BigEndian.PutUint16(header[2:], uint16(len(name)))
	copy(header[2+2:], []byte(name))

	return header
}

func (db Database) ParseRecordV1(id uuid.UUID, password string, raw []byte) (v2.Record, error) {
	length := binary.BigEndian.Uint16(raw)
	raw = raw[2:]

	name := string(raw[:length])
	raw = raw[length:]

	r, err := database.DecryptRecord[record.RecordV1](password, raw)
	if err != nil {
		return v2.Record{}, err
	}

	return r.Upgrade(id, name), nil
}

func (db Database) SaveLegacyRecord(id uuid.UUID, password string, r record.RecordV1) error {
	hash := make([]byte, LEGACY_HASH_SIZE)
	return database.SaveEncryptedRecord(db.RecordPath(id), password, hash, r)
}
