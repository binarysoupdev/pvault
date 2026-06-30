package data

import "pvault/vault/index"

type DatabaseV2 struct {
	Path string
}

func NewDatabaseV2(path string) DatabaseV2 {
	return DatabaseV2{
		Path: path,
	}
}

func (DatabaseV2) GetVersion() uint16 {
	return 2
}

func (DatabaseV2) Upgrade(idx index.IndexMap, target DatabaseV2) error {
	return NotSupportedError{}
}
