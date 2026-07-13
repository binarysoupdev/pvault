package v1

import (
	"encoding/binary"
	"pvault/crypt"
	v2 "pvault/vault/record/version/v2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func Unmarshal(password string, bytes []byte, id uuid.UUID) (v2.Record, error) {
	version := binary.BigEndian.Uint16(bytes)
	if version != VERSION {
		return v2.Record{}, errors.Format("incorrect version \"%d\"", version)
	}

	length := binary.BigEndian.Uint16(bytes[2:])
	name := string(bytes[2+2 : length])

	r, err := crypt.Unmarshal[Record](password, bytes[2+2+length:])
	if err != nil {
		return v2.Record{}, errors.Chain(err, "error decrypting record")
	}

	return r.Upgrade(id, name), nil
}
