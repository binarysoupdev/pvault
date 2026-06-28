package legacy

import (
	"errors"
	"pvault/vault/data"
	"pvault/vault/index"
)

type DatabaseV1 struct{}

func (DatabaseV1) Version() uint16 {
	return 1
}

func (DatabaseV1) SaveIndex(path string, idx index.IndexMap) error {
	return errors.New("not supported")
}

func (DatabaseV1) LoadIndex(path string) (index.IndexMap, error) {
	// TODO: implement
	return index.IndexMap{}, nil
}

func (DatabaseV1) Upgrade(path string, idx index.IndexMap, target data.Database) error {
	// TODO: implement
	return nil
}
