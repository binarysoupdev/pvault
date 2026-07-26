package v2

import (
	"path/filepath"
	index_v2 "pvault/app/vault/index/encoder/v2"
	record_v2 "pvault/app/vault/record/encoder/v2"

	"github.com/google/uuid"
)

const (
	VERSION    = 2
	INDEX_FILE = "index.bin"
)

type IndexEncoder = index_v2.Encoder
type RecordEncoder = record_v2.Encoder

type Database struct {
	IndexEncoder
	RecordEncoder
}

func (idx Database) IndexPath(path string) string {
	return filepath.Join(path, INDEX_FILE)
}

func (idx Database) RecordPath(path string, id uuid.UUID) string {
	return filepath.Join(path, id.String())
}
