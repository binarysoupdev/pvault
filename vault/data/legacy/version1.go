package legacy

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/index"
	"pvault/vault/record"
	"pvault/vault/record/legacy"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/google/uuid"
)

type DatabaseV1 struct {
	IndexPath string
}

func NewDatabaseV1(path string) DatabaseV1 {
	return DatabaseV1{
		IndexPath: path,
	}
}

func (DatabaseV1) SaveIndex(idx index.IndexMap) error {
	return data.NotSupportedError{}
}

func (DatabaseV1) LoadIndex() (index.IndexMap, error) {
	// TODO: implement
	return index.IndexMap{}, nil
}

func (DatabaseV1) SaveRecord(r record.Record, password string) error {
	return data.NotSupportedError{}
}

func (DatabaseV1) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	return record.Record{}, data.NotSupportedError{}
}

func (DatabaseV1) DeleteRecord(id uuid.UUID) error {
	return data.NotSupportedError{}
}

//=====================================

func (DatabaseV1) GetVersion() uint16 {
	return 1
}

func (db DatabaseV1) RecordPath(id uuid.UUID) string {
	return filepath.Join(filepath.Dir(db.IndexPath), id.String()+".crypt")
}

func (db DatabaseV1) SaveTestRecord(path string, r record.Record, password string) error {
	old := legacy.RecordV1{
		Password: r.Password,
		Username: r.Username,
	}

	c, salt := crypt.NewFromPassword(password)

	plaintext, err := json.Marshal(old)
	if err != nil {
		return errors.Chain(err, "error marshaling json")
	}

	ciphertext := c.Encrypt(plaintext)

	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating record file")
	}
	defer file.Close()

	db.writeRecordMeta(file, r.Name)
	file.Write(salt)
	file.Write(ciphertext)

	return nil
}

func (db DatabaseV1) writeRecordMeta(file *os.File, name string) {
	const LEGACY_RECORD_VERSION = 1

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, LEGACY_RECORD_VERSION)
	file.Write(version)

	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(name)))
	file.Write(length)
	file.Write([]byte(name))
}

func (db DatabaseV1) Upgrade(idx index.IndexMap, target data.DatabaseV2) error {
	const LEGACY_HASH_SIZE = 60

	for name, id := range idx {
		oldFile := db.RecordPath(id)

		raw, err := os.ReadFile(oldFile)
		if err != nil {
			continue
		}

		file, err := os.Create(target.RecordPath(id))
		if err != nil {
			return errors.Chain(err, "error creating converted record file")
		}
		defer file.Close()

		db.writeRecordMeta(file, name)
		file.Write(raw[LEGACY_HASH_SIZE:])

		_ = os.Remove(oldFile)
	}

	_ = os.Remove(db.IndexPath)
	return nil
}
