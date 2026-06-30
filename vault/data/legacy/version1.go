package legacy

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/index"
	"pvault/vault/record"
	"pvault/vault/record/legacy"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/google/uuid"
)

type DatabaseV1 struct {
	Path string
}

func NewDatabaseV1(path string) DatabaseV1 {
	return DatabaseV1{
		Path: path,
	}
}

func (DatabaseV1) GetVersion() uint16 {
	return 1
}

func (DatabaseV1) Upgrade(target data.Database) error {
	// TODO: implement
	return nil
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

func (DatabaseV1) SaveTestRecord(path string, r record.Record, password string) error {
	const LEGACY_RECORD_VERSION = 1

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

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, LEGACY_RECORD_VERSION)
	file.Write(version)

	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(r.Name)))
	file.Write(length)
	file.Write([]byte(r.Name))

	file.Write(salt)
	file.Write(ciphertext)

	return nil
}
