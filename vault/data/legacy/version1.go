package legacy

import (
	"errors"
	"pvault/vault/index"
)

type DatabaseV1 struct {
	Path string
}

func NewDatabaseV1(path string) DatabaseV1 {
	return DatabaseV1{
		Path: path,
	}
}

func (DatabaseV1) Version() uint16 {
	return 1
}

func (DatabaseV1) SaveIndex(idx index.IndexMap) error {
	return errors.New("save index not supported")
}

func (DatabaseV1) LoadIndex() (index.IndexMap, error) {
	// TODO: implement
	return index.IndexMap{}, nil
}

func (DatabaseV1) Upgrade() error {
	// TODO: implement
	return nil
}
