package version1

import (
	"pvault/vault/data"
	"pvault/vault/index"
	"pvault/vault/record"

	"github.com/google/uuid"
)

type Database struct {
	Path string
}

func NewDatabase(path string) Database {
	return Database{
		Path: path,
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
