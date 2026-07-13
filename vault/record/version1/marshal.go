package v1

import (
	"encoding/binary"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func Unmarshal(password string, bytes []byte, id uuid.UUID) (Record, error) {
	version := binary.BigEndian.Uint16(bytes)
	if version != VERSION {
		return Record{}, errors.Format("incorrect version \"%d\"", version)
	}

	length := binary.BigEndian.Uint16(bytes[2:])
	name := string(bytes[2+2 : 2+2+length])

	r, err := crypt.Unmarshal[Record](password, bytes[2+2+length:])
	if err != nil {
		return Record{}, errors.Chain(err, "error decrypting record")
	}

	r.ID = id
	r.Name = name

	return r, nil
}

func (r Record) Marshal(password string) ([]byte, error) {
	ciphertext, err := crypt.Marshal(password, r)
	if err != nil {
		return nil, err
	}

	return append(buildHeader(r.Name), ciphertext...), nil
}

func buildHeader(name string) []byte {
	header := make([]byte, 2+2+len(name))

	binary.BigEndian.PutUint16(header, VERSION)
	binary.BigEndian.PutUint16(header[2:], uint16(len(name)))
	copy(header[2+2:], []byte(name))

	return header
}
