package v2

import (
	"encoding/binary"
	"os"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (r Record) GetVersion() int {
	return VERSION
}

func (r Record) GetID() uuid.UUID {
	return r.ID
}

func (r Record) GetName() string {
	return r.Name
}

func (r Record) Upgrade() Record {
	return r
}

func (r Record) Validate() error {
	errs := errors.Errors{}

	if r.ID == uuid.Nil {
		errs.AddNew("id cannot be nil (all zeroes)")
	}
	if r.Name == "" {
		errs.AddNew("name cannot be empty")
	}

	return errs.Collapse(", ")
}

func (r Record) SaveFile(path string, password string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating record file")
	}
	defer file.Close()

	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, VERSION)
	file.Write(version)

	_, err = crypt.Encode(file, password, r)
	if err != nil {
		return errors.Chain(err, "error encrypting record")
	}

	return nil
}
