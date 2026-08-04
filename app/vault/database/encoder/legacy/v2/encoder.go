package v2

import (
	"path/filepath"
	index_v2 "pvault/app/vault/index/encoder/v2"
	record_v2 "pvault/app/vault/record/encoder/legacy/v2"

	"github.com/google/uuid"
)

const VERSION = 2

type IndexEncoder = index_v2.Encoder
type RecordEncoder = record_v2.Encoder

type Encoder struct {
	IndexEncoder
	RecordEncoder
}

func (db Encoder) GetVersion() int {
	return VERSION
}

func (idx Encoder) IndexPath(path string) string {
	return filepath.Join(path, "index.bin")
}

func (idx Encoder) RecordPath(path string, id uuid.UUID) string {
	return filepath.Join(path, id.String())
}
