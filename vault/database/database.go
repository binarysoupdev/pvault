package database

import (
	"pvault/vault/index"
	v2 "pvault/vault/record/version/v2"

	"github.com/google/uuid"
)

type Database interface {
	GetVersion() uint16

	Initialize(idx index.IndexMap) error
	Upgrade(idx index.IndexMap, target Database) error

	IndexPath() string
	SaveIndex(idx index.IndexMap) error
	LoadIndex() (index.IndexMap, error)

	RecordPath(id uuid.UUID) string
	SaveRecord(r v2.Record, password string) error
	LoadRecord(id uuid.UUID, password string) (v2.Record, error)
	DeleteRecord(id uuid.UUID) error
}
