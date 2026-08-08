package database

import (
	v3 "pvault/vault/database/encoder/v3"
	"pvault/vault/index"
	"pvault/vault/record"

	"github.com/google/uuid"
)

const CURRENT_VERSION = v3.VERSION

type Encoder interface {
	GetVersion() int

	IndexPath(path string) string
	index.Encoder

	RecordPath(path string, id uuid.UUID) string
	record.Encoder

	Upgrade(path string) (v3.Encoder, error)
}
