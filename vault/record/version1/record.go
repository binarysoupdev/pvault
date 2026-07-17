package v1

import (
	"encoding/binary"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const VERSION = 1

type Record struct {
	Password      string   `json:"password"`
	Username      string   `json:"username"`
	URL           string   `json:"url"`
	RecoveryCodes []string `json:"recovery_codes"`

	ID   uuid.UUID
	Name string
}

func Unmarshal(bytes []byte, password string, id uuid.UUID) (Record, error) {
	version := binary.BigEndian.Uint16(bytes)
	if version != VERSION {
		return Record{}, errors.Format("incorrect version \"%d\"", version)
	}

	length := binary.BigEndian.Uint16(bytes[2:])
	name := string(bytes[2+2 : 2+2+length])

	record, err := crypt.Unmarshal[Record](password, bytes[2+2+length:])
	if err != nil {
		return Record{}, errors.Chain(err, "error decrypting record")
	}

	record.ID = id
	record.Name = name

	return record, nil
}
