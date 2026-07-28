package v1

import (
	"io"
	v1 "pvault/app/vault/record/record/v1"
	"pvault/crypt"
	"pvault/util"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const HASH_SIZE = 60

func (e Encoder) EncodeRawV1(w io.Writer, data []byte) error {
	hash := make([]byte, HASH_SIZE)
	return util.WriteBytes(w, hash, data)
}

func (e Encoder) DecodeV1(r io.Reader, password string, id uuid.UUID, name string) (v1.Record, error) {
	data, err := e.DecodeRawV1(r)
	if err != nil {
		return v1.Record{}, errors.Chain(err, "error decoding record")
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
	_, err := util.ReadBytes(r, HASH_SIZE)
	if err != nil {
		return nil, errors.Chain(err, "error decoding hash prefix")
	}

	return io.ReadAll(r)
}
