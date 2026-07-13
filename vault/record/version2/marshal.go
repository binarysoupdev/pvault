package v2

import (
	"encoding/binary"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
)

func Unmarshal(password string, bytes []byte) (Record, error) {
	version := binary.BigEndian.Uint16(bytes)
	if version != VERSION {
		return Record{}, errors.Format("incorrect version \"%d\"", version)
	}

	return crypt.Unmarshal[Record](password, bytes[2:])
}

func (r Record) Marshal(password string) ([]byte, error) {
	ciphertext, err := crypt.Marshal(password, r)
	if err != nil {
		return nil, err
	}

	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, VERSION)

	return append(header, ciphertext...), nil
}
