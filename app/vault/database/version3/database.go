package v3

import (
	"path/filepath"
	"pvault/app/vault/index"
	index_v2 "pvault/app/vault/index/encoder/version2"
	record_v3 "pvault/app/vault/record/encoder/version3"

	"github.com/google/uuid"
)

const (
	VERSION    = 3
	INDEX_FILE = "INDEX"
)

type IndexEncoder = index_v2.Encoder
type RecordEncoder = record_v3.Encoder

type Database struct {
	Path string
	IndexEncoder
	RecordEncoder
}

func NewDatabase(path string) Database {
	return Database{
		Path:          path,
		IndexEncoder:  IndexEncoder{},
		RecordEncoder: RecordEncoder{},
	}
}

func (Database) GetVersion() int {
	return VERSION
}

func (db Database) GetPath() string {
	return db.Path
}

func (db Database) IndexPath() string {
	return filepath.Join(db.Path, INDEX_FILE)
}

func (db Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(db.Path, id.String())
}

func (db Database) Upgrade(idx index.IndexMap) (Database, error) {
	return db, nil
}
