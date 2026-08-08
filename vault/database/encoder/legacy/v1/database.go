package v1

import (
	"path/filepath"
	index_v1 "pvault/vault/index/encoder/legacy/v1"
	record_v1 "pvault/vault/record/encoder/legacy/v1"

	"github.com/google/uuid"
)

const VERSION = 1

type IndexEncoder = index_v1.Encoder
type RecordEncoder = record_v1.Encoder

type Encoder struct {
	IndexEncoder
	RecordEncoder
}

func (db Encoder) GetVersion() int {
	return VERSION
}

func (db Encoder) IndexPath(path string) string {
	return filepath.Join(path, "index.txt")
}

func (db Encoder) RecordPath(path string, id uuid.UUID) string {
	return filepath.Join(path, id.String()+".crypt")
}
