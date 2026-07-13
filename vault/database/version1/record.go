package version1

import (
	"encoding/binary"
	"path/filepath"
	"pvault/vault/database"
	"pvault/vault/record"

	"github.com/google/uuid"
)

const (
	RECORD_VERSION   uint16 = 1
	LEGACY_HASH_SIZE uint16 = 60
)

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String()+".crypt")
}

func (db Database) SaveRecord(r record.RecordV2, password string) error {
	return database.NotSupportedError{}
}

func (Database) LoadRecord(id uuid.UUID, password string) (record.RecordV2, error) {
	return record.RecordV2{}, database.NotSupportedError{}
}

func (Database) DeleteRecord(id uuid.UUID) error {
	return database.NotSupportedError{}
}

func (db Database) SaveRecordV1(path string, r record.RecordV2, password string) error {
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

func (db Database) ParseRecordV1(id uuid.UUID, password string, raw []byte) (record.RecordV2, error) {
	length := binary.BigEndian.Uint16(raw)
	raw = raw[2:]

	name := string(raw[:length])
	raw = raw[length:]

	r, err := database.DecryptRecord[record.RecordV1](password, raw)
	if err != nil {
		return record.RecordV2{}, err
	}

	return r.Upgrade(id, name), nil
}

func (db Database) SaveLegacyRecord(id uuid.UUID, password string, r record.RecordV1) error {
	hash := make([]byte, LEGACY_HASH_SIZE)
	return database.SaveEncryptedRecord(db.RecordPath(id), password, hash, r)
}
