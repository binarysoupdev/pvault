package v2

import (
	"encoding/binary"
	"io"
	v1 "pvault/app/vault/record/record/v1"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

type rawV1 struct {
	Data []byte
	Name string
}

func (e Encoder) EncodeRawV1(w io.Writer, data []byte, name string) error {
	header := make([]byte, 2+2)
	binary.BigEndian.PutUint16(header, v1.VERSION)
	binary.BigEndian.PutUint16(header[2:], uint16(len(name)))

	w.Write(header)
	w.Write([]byte(name))
	w.Write(data)

	return nil
}

func (e Encoder) DecodeV1(r io.Reader, password string, id uuid.UUID) (v1.Record, error) {
	raw, err := e.DecodeRawV1(r)
	if err != nil {
		return v1.Record{}, err
	}

	record, err := crypt.Unmarshal[v1.Record](password, raw.Data)
	if err != nil {
		return v1.Record{}, errors.Chain(err, "error decrypting record")
	}

	record.ID = id
	record.Name = raw.Name

	return record, nil
}

func (e Encoder) DecodeRawV1(r io.Reader) (rawV1, error) {
	length := make([]byte, 2)
	r.Read(length)

	name := make([]byte, binary.BigEndian.Uint16(length))
	r.Read(name)

	data, _ := io.ReadAll(r)

	return rawV1{
		Data: data,
		Name: string(name),
	}, nil
}
