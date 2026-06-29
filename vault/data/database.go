package data

import (
	"pvault/vault/index"
	"pvault/vault/record"

	"github.com/google/uuid"
)

type Database interface {
	Version() uint16
	Upgrade() error

	SaveIndex(idx index.IndexMap) error
	LoadIndex() (index.IndexMap, error)

	SaveRecord(r record.Record, password string) error
	LoadRecord(id uuid.UUID, password string) (record.Record, error)
	DeleteRecord(id uuid.UUID) error
}
