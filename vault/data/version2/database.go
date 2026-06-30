package version2

import (
	"pvault/vault/data"
	"pvault/vault/index"
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
	return 2
}

func (Database) UpgradeToVersion2(idx index.IndexMap) error {
	return data.NotSupportedError{}
}
