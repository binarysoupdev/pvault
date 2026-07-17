package record

import (
	v2 "pvault/vault/record/version2"

	"github.com/google/uuid"
)

type Record interface {
	GetVersion() int

	GetID() uuid.UUID
	GetName() string

	Validate() error
	SaveFile(path string, password string) error

	Upgrade() v2.Record
}
