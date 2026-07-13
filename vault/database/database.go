package database

import (
	"pvault/vault/index"
	record "pvault/vault/record"

	"github.com/google/uuid"
)

type Database interface {
	GetVersion() uint16

	Initialize(idx index.IndexMap) error
	Upgrade(idx index.IndexMap) error

	IndexPath() string
	SaveIndex(idx index.IndexMap) error
	LoadIndex() (index.IndexMap, error)

	RecordPath(id uuid.UUID) string
	SaveRecord(r record.Record, password string) error
	LoadRecord(id uuid.UUID, password string) (record.Record, error)
	DeleteRecord(id uuid.UUID) error
}
