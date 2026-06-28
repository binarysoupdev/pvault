package data

import "pvault/vault/index"

type Database interface {
	Version() uint16
	SaveIndex(path string, idx index.IndexMap) error
	LoadIndex(path string) (index.IndexMap, error)
	Upgrade(path string, idx index.IndexMap, target Database) error
}

type CurrentDatabase = DatabaseV2
