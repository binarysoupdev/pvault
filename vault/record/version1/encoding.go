package v1

import (
	"encoding/binary"
	"io"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

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

func (r Record) Validate() error {
	// TODO: implement
	return nil
}

func (r Record) Encode(w io.Writer, password string) error {
	writeHeader(w, r.Name)

	_, err := crypt.Encode(w, password, r)
	if err != nil {
		return errors.Chain(err, "error encrypting record")
	}

	return nil
}

func writeHeader(w io.Writer, name string) {
	header := make([]byte, 2+2)
	binary.BigEndian.PutUint16(header, VERSION)
	binary.BigEndian.PutUint16(header[2:], uint16(len(name)))

	w.Write(header)
	w.Write([]byte(name))
}
