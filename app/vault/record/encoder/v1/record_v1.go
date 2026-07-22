package v1

import (
	"io"
	v1 "pvault/app/vault/record/record/v1"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const HASH_SIZE = 60

func (e Encoder) EncodeRawV1(w io.Writer, data []byte) error {
	hash := make([]byte, HASH_SIZE)
	w.Write(hash)
	w.Write(data)
	return nil
}

func (e Encoder) DecodeV1(r io.Reader, password string, id uuid.UUID, name string) (v1.Record, error) {
	data, err := e.DecodeRawV1(r)
	if err != nil {
		return v1.Record{}, err
	}

	record, err := crypt.Unmarshal[v1.Record](password, []byte(data))
	if err != nil {
		return v1.Record{}, errors.Chain(err, "error decrypting record")
	}

	record.ID = id
	record.Name = name

	return record, nil
}

func (e Encoder) DecodeRawV1(r io.Reader) ([]byte, error) {
	hash := make([]byte, HASH_SIZE)
	r.Read(hash)

	return io.ReadAll(r)
}
