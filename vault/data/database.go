package data

import (
	"pvault/vault/index"
	"pvault/vault/record"

	"github.com/google/uuid"
)

type Database interface {
	GetVersion() uint16
	UpgradeToVersion2(idx index.IndexMap) error

	SaveIndex(idx index.IndexMap) error
	LoadIndex() (index.IndexMap, error)

	SaveRecord(r record.Record, password string) error
	LoadRecord(id uuid.UUID, password string) (record.Record, error)
	DeleteRecord(id uuid.UUID) error
}
