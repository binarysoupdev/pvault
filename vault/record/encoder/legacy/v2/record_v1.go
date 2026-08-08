package v2

import (
	"encoding/binary"
	"io"
	"pvault/util"
	"pvault/vault/record"
	v1 "pvault/vault/record/record/legacy/v1"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

type rawV1 struct {
	Data []byte
	Name string
}

func (e Encoder) EncodeV1(w io.Writer, password string, r record.Record) error {
	bytes, err := record.Encrypt(r, password)
	if err != nil {
		return errors.Chain(err, "error encrypting record v1")
	}

	err = e.EncodeRawV1(w, bytes, r.GetName())
	if err != nil {
		return errors.Chain(err, "error encoding record v1")
	}

	return nil
}

func (e Encoder) EncodeRawV1(w io.Writer, data []byte, name string) error {
	header := make([]byte, 2+2)
	binary.BigEndian.PutUint16(header, v1.VERSION)
	binary.BigEndian.PutUint16(header[2:], uint16(len(name)))

	return util.WriteBytes(w, header, []byte(name), data)
}

func (e Encoder) DecodeV1(r io.Reader, password string, id uuid.UUID) (v1.Record, error) {
	raw, err := e.DecodeRawV1(r)
	if err != nil {
		return v1.Record{}, errors.Chain(err, "error decoding record v1")
	}

	record, err := record.Decrypt[v1.Record](raw.Data, password)
	if err != nil {
		return v1.Record{}, errors.Chain(err, "error decrypting record v1")
	}

	record.ID = id
	record.Name = raw.Name

	return record, nil
}

func (e Encoder) DecodeRawV1(r io.Reader) (rawV1, error) {
	length, err := util.ReadBytes(r, 2)
	if err != nil {
		return rawV1{}, errors.Chain(err, "error reading length")
	}

	name, err := util.ReadBytes(r, int(binary.BigEndian.Uint16(length)))
	if err != nil {
		return rawV1{}, errors.Chain(err, "error reading name")
	}

	data, _ := io.ReadAll(r)

	return rawV1{
		Data: data,
		Name: string(name),
	}, nil
}
