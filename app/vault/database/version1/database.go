package v1

import (
	"path/filepath"
	index_v1 "pvault/app/vault/index/encoder/version1"
	record_v1 "pvault/app/vault/record/encoder/version1"

	"github.com/google/uuid"
)

const (
	VERSION    = 1
	INDEX_FILE = "index.txt"
)

type IndexEncoder = index_v1.Encoder
type RecordEncoder = record_v1.Encoder

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
	return filepath.Join(db.Path, id.String()+".crypt")
}
