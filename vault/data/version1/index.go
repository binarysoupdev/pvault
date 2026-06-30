package version1

import (
	"path/filepath"
	"pvault/vault/data"
	"pvault/vault/index"
)

const INDEX_FILE = "index.txt"

func (db Database) IndexPath() string {
	return filepath.Join(db.Path, INDEX_FILE)
}

func (Database) SaveIndex(idx index.IndexMap) error {
	return data.NotSupportedError{}
}

func (Database) LoadIndex() (index.IndexMap, error) {
	// TODO: implement
	return index.IndexMap{}, nil
}
