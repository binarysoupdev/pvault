package v1

import (
	"encoding/binary"
	"io"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func Decode(r io.Reader, password string, id uuid.UUID) (Record, error) {
	length := make([]byte, 2)
	r.Read(length)

	name := make([]byte, binary.BigEndian.Uint16(length))
	r.Read(name)

	record, err := crypt.Decode[Record](r, password)
	if err != nil {
		return Record{}, errors.Chain(err, "error decrypting record")
	}

	record.ID = id
	record.Name = string(name)

	return record, nil
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
