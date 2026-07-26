package v3

import (
	"path/filepath"
	index_v2 "pvault/app/vault/index/encoder/v2"
	record_v3 "pvault/app/vault/record/encoder/v3"

	"github.com/google/uuid"
)

const (
	VERSION    = 3
	INDEX_FILE = "INDEX"
)

type IndexEncoder = index_v2.Encoder
type RecordEncoder = record_v3.Encoder

type Database struct {
	IndexEncoder
	RecordEncoder
}

func (db Database) IndexPath(path string) string {
	return filepath.Join(path, INDEX_FILE)
}

func (db Database) RecordPath(path string, id uuid.UUID) string {
	return filepath.Join(path, id.String())
}

func (db Database) Upgrade(_ string) (Database, error) {
	return db, nil
}
