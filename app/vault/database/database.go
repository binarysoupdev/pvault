package database

import (
	v3 "pvault/app/vault/database/database/v3"
	"pvault/app/vault/index"
	"pvault/app/vault/record"

	"github.com/google/uuid"
)

const CURRENT_VERSION = v3.VERSION

type Database interface {
	IndexPath() string
	index.Encoder

	RecordPath(id uuid.UUID) string
	record.Encoder

	Upgrade(idx index.IndexMap) (v3.Database, error)
}
