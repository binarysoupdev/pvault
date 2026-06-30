package version1

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

type Database struct {
	IndexPath string
}

func NewDatabase(path string) Database {
	return Database{
		IndexPath: path,
	}
}

func (Database) GetVersion() uint16 {
	return 1
}

func (Database) SaveIndex(idx index.IndexMap) error {
	return data.NotSupportedError{}
}

func (Database) LoadIndex() (index.IndexMap, error) {
	// TODO: implement
	return index.IndexMap{}, nil
}

func (Database) SaveRecord(r record.Record, password string) error {
	return data.NotSupportedError{}
}

func (Database) LoadRecord(id uuid.UUID, password string) (record.Record, error) {
	return record.Record{}, data.NotSupportedError{}
}

func (Database) DeleteRecord(id uuid.UUID) error {
	return data.NotSupportedError{}
}

//=====================================

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(filepath.Dir(db.IndexPath), id.String()+".crypt")
}

func (db Database) SaveTestRecord(path string, r record.Record, password string) error {
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

func (db Database) writeRecordMeta(file *os.File, name string) {
	const LEGACY_RECORD_VERSION = 1

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, LEGACY_RECORD_VERSION)
	file.Write(version)

	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(name)))
	file.Write(length)
	file.Write([]byte(name))
}
