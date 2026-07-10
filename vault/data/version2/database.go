package version2

import (
	"pvault/errors"
	"pvault/vault/data"
	"pvault/vault/index"
)

const VERSION = 2

type Database struct {
	Path string
}

func New(path string) Database {
	return Database{
		Path: path,
	}
}

func (Database) GetVersion() uint16 {
	return VERSION
}

func (db Database) Initialize(idx index.IndexMap) error {
	err := db.SaveIndex(idx)
	if err != nil {
		return errors.Chain(err, "error saving index file")
	}

	return nil
}

func (Database) Upgrade(idx index.IndexMap, target data.Database) error {
	return data.NotSupportedError{}
}
