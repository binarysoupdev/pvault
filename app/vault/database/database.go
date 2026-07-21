package database

import (
	v3 "pvault/app/vault/database/version3"
	idx "pvault/app/vault/index"
	index "pvault/app/vault/index/encoder"
	record "pvault/app/vault/record/encoder"

	"github.com/google/uuid"
)

type Database interface {
	GetVersion() int

	GetPath() string
	IndexPath() string
	RecordPath(id uuid.UUID) string

	index.Encoder
	record.Encoder

	Upgrade(idx idx.IndexMap) (v3.Database, error)
}
