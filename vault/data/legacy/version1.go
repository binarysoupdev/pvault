package legacy

import (
	"pvault/vault/data"
	"pvault/vault/index"
	"pvault/vault/record"

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

func (DatabaseV1) Upgrade() error {
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
