package v2

import (
	"path/filepath"
	index_v2 "pvault/app/vault/index/encoder/version2"
	record_v2 "pvault/app/vault/record/encoder/version2"

	"github.com/google/uuid"
)

const (
	VERSION    = 2
	INDEX_FILE = "index.bin"
)

type IndexEncoder = index_v2.Encoder
type RecordEncoder = record_v2.Encoder

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

func (idx Database) IndexPath() string {
	return filepath.Join(idx.Path, INDEX_FILE)
}

func (idx Database) RecordPath(id uuid.UUID) string {
	return filepath.Join(idx.Path, id.String())
}
