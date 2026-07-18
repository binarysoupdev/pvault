package vault

import (
	"pvault/app/vault/record"

	"github.com/google/uuid"
)

type Vault interface {
	GetVersion() int

	Load() error
	SearchNames(term string) []string

	ValidateRecord(r record.Record) error
	SaveRecord(r record.Record, password string) error
	LoadRecord(name string, password string) (record.Record, error)
	DeleteRecord(name string) (uuid.UUID, error)
}
