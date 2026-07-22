package v3

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
	ID   uuid.UUID
	Name string
}

func (e Encoder) EncodeRawV1(w io.Writer, data []byte, id uuid.UUID, name string) error {
	header := make([]byte, 2+2)
	binary.BigEndian.PutUint16(header, v1.VERSION)
	binary.BigEndian.PutUint16(header[2:], uint16(len(name)))

	w.Write(header)
	w.Write(id[:])
	w.Write([]byte(name))
	w.Write(data)

	return nil
}

func (e Encoder) DecodeV1(r io.Reader, password string) (v1.Record, error) {
	raw, err := e.DecodeRawV1(r)
	if err != nil {
		return v1.Record{}, err
	}

	record, err := crypt.Unmarshal[v1.Record](password, raw.Data)
	if err != nil {
		return v1.Record{}, errors.Chain(err, "error decrypting record")
	}

	record.ID = raw.ID
	record.Name = raw.Name

	return record, nil
}

func (e Encoder) DecodeRawV1(r io.Reader) (rawV1, error) {
	length := make([]byte, 2)
	r.Read(length)

	id := uuid.UUID{}
	r.Read(id[:])

	name := make([]byte, binary.BigEndian.Uint16(length))
	r.Read(name)

	data, _ := io.ReadAll(r)

	return rawV1{
		Data: data,
		ID:   id,
		Name: string(name),
	}, nil
}
