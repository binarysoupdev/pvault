package record

import (
	"io"
	v2 "pvault/vault/record/version2"

	"github.com/google/uuid"
)

type Record interface {
	GetVersion() int

	GetID() uuid.UUID
	GetName() string

	Validate() error
	Encode(w io.Writer, password string) error

	Upgrade() v2.Record
}
