package v1

import (
	"encoding/binary"
	"io"
	"os"
	"pvault/crypt"
	v2 "pvault/vault/record/version2"

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

func (r Record) Validate() error {
	return nil
}

func (r Record) SaveFile(path string, password string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating record file")
	}
	defer file.Close()

	writeHeader(file, r.Name)

	_, err = crypt.Encode(file, password, r)
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

func (r Record) Upgrade() v2.Record {
	other := map[string]any{}

	if r.URL != "" {
		other["url"] = r.URL
	}
	if len(r.RecoveryCodes) > 0 {
		other["recovery_codes"] = r.RecoveryCodes
	}

	return v2.Record{
		ID:       r.ID,
		Name:     r.Name,
		Username: r.Username,
		Password: r.Password,
		Other:    other,
	}
}
