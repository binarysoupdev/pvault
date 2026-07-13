package v2

import (
	"encoding/binary"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
)

func (r Record) Marshal(password string) ([]byte, error) {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, VERSION)

	ciphertext, err := crypt.Marshal(password, r)
	if err != nil {
		return nil, err
	}

	return append(header, ciphertext...), nil
}

func Unmarshal(password string, bytes []byte) (Record, error) {
	version := binary.BigEndian.Uint16(bytes)
	if version != VERSION {
		return Record{}, errors.Format("incorrect version \"%d\"", version)
	}

	return crypt.Unmarshal[Record](password, bytes[2:])
}
