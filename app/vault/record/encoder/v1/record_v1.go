package v1

import (
	"io"
	"pvault/app/vault/record"
	v1 "pvault/app/vault/record/record/v1"
	"pvault/util"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const HASH_SIZE = 60

func (e Encoder) EncodeV1(w io.Writer, password string, r record.Record) error {
	bytes, err := record.Encrypt(r, password)
	if err != nil {
		return errors.Chain(err, "error encrypting record v1")
	}

	err = e.EncodeRawV1(w, bytes)
	if err != nil {
		return errors.Chain(err, "error encoding record v1")
	}

	return nil
}

func (e Encoder) EncodeRawV1(w io.Writer, data []byte) error {
	hash := make([]byte, HASH_SIZE)
	return util.WriteBytes(w, hash, data)
}

func (e Encoder) DecodeV1(r io.Reader, password string, id uuid.UUID, name string) (v1.Record, error) {
	bytes, err := e.DecodeRawV1(r)
	if err != nil {
		return v1.Record{}, errors.Chain(err, "error decoding record v1")
	}

	r1, err := record.Decrypt[v1.Record](bytes, password)
	if err != nil {
		return v1.Record{}, errors.Chain(err, "error decrypting record v1")
	}

	r1.ID = id
	r1.Name = name

	return r1, nil
}

func (e Encoder) DecodeRawV1(r io.Reader) ([]byte, error) {
	_, err := util.ReadBytes(r, HASH_SIZE)
	if err != nil {
		return nil, errors.Chain(err, "error decoding hash prefix")
	}

	return io.ReadAll(r)
}
