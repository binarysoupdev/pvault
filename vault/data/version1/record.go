package version1

import (
	"encoding/binary"
	"path/filepath"
	"pvault/vault/data"
	"pvault/vault/record"
	"pvault/vault/record/legacy"

	"github.com/google/uuid"
)

const (
	RECORD_VERSION   uint16 = 1
	LEGACY_HASH_SIZE uint16 = 60
)

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String())
}

func (db Database) SaveRecord(r record.Record, password string) error {
	return data.NotSupportedError{}
}

func (Database) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	return record.Record{}, data.NotSupportedError{}
}

func (Database) DeleteRecord(id uuid.UUID) error {
	return data.NotSupportedError{}
}

func (db Database) SaveRecordV1(r record.Record, password string) error {
	v1 := legacy.RecordV1{
		Password: r.Password,
		Username: r.Username,
	}
	header := db.buildRecordV1Header(r.Name)

	return data.SaveEncryptedRecord(db.RecordPath(r.ID), password, header, v1)
}

func (Database) buildRecordV1Header(name string) []byte {
	header := make([]byte, 2+2+len(name))

	binary.BigEndian.PutUint16(header, RECORD_VERSION)
	binary.BigEndian.PutUint16(header[2:], uint16(len(name)))
	copy(header[2+2:], []byte(name))

	return header
}

func (db Database) ParseRecordV1(id uuid.UUID, password string, raw []byte) (record.Record, error) {
	length := binary.BigEndian.Uint16(raw)
	raw = raw[2:]

	name := string(raw[:length])
	raw = raw[length:]

	r, err := data.DecryptRecord[legacy.RecordV1](password, raw)
	if err != nil {
		return record.Record{}, err
	}

	return r.Upgrade(id, name), nil
}

func (db Database) LegacyRecordPath(id uuid.UUID) string {
	return db.RecordPath(id) + ".crypt"
}

func (db Database) SaveLegacyRecord(id uuid.UUID, password string, r legacy.RecordV1) error {
	hash := make([]byte, LEGACY_HASH_SIZE)
	return data.SaveEncryptedRecord(db.LegacyRecordPath(id), password, hash, r)
}
