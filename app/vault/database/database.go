package database

import (
	v3 "pvault/app/vault/database/database/v3"
	"pvault/app/vault/index"
	"pvault/app/vault/record"

	"github.com/google/uuid"
)

const CURRENT_VERSION = v3.VERSION

type Database interface {
	GetVersion() int

	IndexPath(path string) string
	index.Encoder

	RecordPath(path string, id uuid.UUID) string
	record.Encoder

	Upgrade(path string) (v3.Database, error)
}
