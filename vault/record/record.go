package record

import (
	"io"
	v2 "pvault/vault/record/version2"

	"github.com/binarysoupdev/go-commando/json"
	"github.com/google/uuid"
)

type Record interface {
	GetID() uuid.UUID
	GetName() string

	Validate() error
	Encode(w io.Writer, password string) error
	Upgrade() v2.Record
}

func New(name string) v2.Record {
	return v2.NewEmptyRecord(name)
}

func LoadFromFile(path string) (v2.Record, error) {
	return json.UnmarshalFile[v2.Record](path)
}
