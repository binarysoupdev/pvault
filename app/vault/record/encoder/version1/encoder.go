package record_v1

import (
	"io"
	"pvault/app/vault/record"
	record_v1 "pvault/app/vault/record/version1"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const (
	VERSION          = 1
	LEGACY_HASH_SIZE = 60
)

type Encoder struct{}

func (e Encoder) EncodeRecord(w io.Writer, password string, r record.Record) error {
	if r.GetVersion() != record_v1.VERSION {
		return errors.Format("unsupported record version \"%d\"", r.GetVersion())
	}

	ciphertext, err := crypt.Marshal(password, r)
	if err != nil {
		return errors.Chain(err, "error encrypting record")
	}

	err = e.EncodeRawV1(w, ciphertext)
	if err != nil {
		return err
	}

	return nil
}

func (e Encoder) EncodeRawV1(w io.Writer, data []byte) error {
	hash := make([]byte, LEGACY_HASH_SIZE)
	w.Write(hash)
	w.Write(data)
	return nil
}

func (e Encoder) DecodeRecord(r io.Reader, password string) (record.Record, error) {
	return e.DecodeV1(r, password, uuid.Nil, "")
}

func (e Encoder) DecodeV1(r io.Reader, password string, id uuid.UUID, name string) (record_v1.Record, error) {
	data, err := e.DecodeRawV1(r)
	if err != nil {
		return record_v1.Record{}, err
	}

	record, err := crypt.Unmarshal[record_v1.Record](password, []byte(data))
	if err != nil {
		return record_v1.Record{}, errors.Chain(err, "error decrypting record")
	}

	record.ID = id
	record.Name = name

	return record, nil
}

func (e Encoder) DecodeRawV1(r io.Reader) ([]byte, error) {
	hash := make([]byte, LEGACY_HASH_SIZE)
	r.Read(hash)

	return io.ReadAll(r)
}
