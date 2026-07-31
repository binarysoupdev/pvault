package v1

import (
	"path/filepath"
	index_v1 "pvault/app/vault/index/encoder/v1"
	record_v1 "pvault/app/vault/record/encoder/v1"

	"github.com/google/uuid"
)

const VERSION = 1

type IndexEncoder = index_v1.Encoder
type RecordEncoder = record_v1.Encoder

type Database struct {
	IndexEncoder
	RecordEncoder
}

func (db Database) GetVersion() int {
	return VERSION
}

func (db Database) IndexPath(path string) string {
	return filepath.Join(path, "index.txt")
}

func (db Database) RecordPath(path string, id uuid.UUID) string {
	return filepath.Join(path, id.String()+".crypt")
}
