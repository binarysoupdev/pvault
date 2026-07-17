package v2

import (
	"encoding/binary"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const VERSION = 2

type Record struct {
	ID       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Other    map[string]any `json:"other"`
}

func NewEmptyRecord(name string) Record {
	return Record{
		ID:       uuid.New(),
		Name:     name,
		Username: "",
		Password: "",
		Other:    map[string]interface{}{},
	}
}

func Unmarshal(bytes []byte, password string) (Record, error) {
	version := binary.BigEndian.Uint16(bytes)
	if version != VERSION {
		return Record{}, errors.Format("incorrect version \"%d\"", version)
	}

	record, err := crypt.Unmarshal[Record](password, bytes[2:])
	if err != nil {
		return Record{}, errors.Chain(err, "error decrypting record")
	}
	return record, nil
}
