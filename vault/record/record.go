package record

import (
	v2 "pvault/vault/record/version2"

	"github.com/google/uuid"
)

type Record interface {
	GetID() uuid.UUID
	Marshal(password string) ([]byte, error)
	Convert() v2.Record
}
