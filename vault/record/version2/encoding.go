package v2

import (
	"encoding/binary"
	"io"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (r Record) Validate() error {
	errs := errors.Errors{}

	if r.ID == uuid.Nil {
		errs.AddNew("\"ID\" cannot be nil (all zeroes)")
	}
	if r.Name == "" {
		errs.AddNew("\"Name\" cannot be empty")
	}

	return errs.Collapse(", ")
}

func Decode(r io.Reader, password string) (Record, error) {
	record, err := crypt.Decode[Record](r, password)
	if err != nil {
		return Record{}, errors.Chain(err, "error decrypting record")
	}
	return record, nil
}

func (r Record) Encode(w io.Writer, password string) error {
	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, VERSION)
	w.Write(version)

	_, err := crypt.Encode(w, password, r)
	if err != nil {
		return errors.Chain(err, "error encrypting record")
	}

	return nil
}
