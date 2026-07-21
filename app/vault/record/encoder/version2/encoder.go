package version2

import (
	"encoding/binary"
	"io"
	"pvault/app/vault/record"
	record_v1 "pvault/app/vault/record/version1"
	record_v2 "pvault/app/vault/record/version2"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

type rawV1 struct {
	Data []byte
	Name string
}

type Encoder struct{}

func (e Encoder) EncodeRecord(w io.Writer, password string, r record.Record) error {
	ciphertext, err := crypt.Marshal(password, r)
	if err != nil {
		return errors.Chain(err, "error encrypting record")
	}

	switch r.GetVersion() {
	case record_v1.VERSION:
		return e.EncodeRawV1(w, ciphertext, r.GetName())
	case record_v2.VERSION:
		return e.EncodeRawV2(w, ciphertext)
	default:
		return errors.Format("unsupported record version \"%d\"", r.GetVersion())
	}
}

func (e Encoder) EncodeRawV1(w io.Writer, data []byte, name string) error {
	header := make([]byte, 2+2)
	binary.BigEndian.PutUint16(header, record_v1.VERSION)
	binary.BigEndian.PutUint16(header[2:], uint16(len(name)))

	w.Write(header)
	w.Write([]byte(name))
	w.Write(data)

	return nil
}

func (e Encoder) EncodeRawV2(w io.Writer, data []byte) error {
	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, record_v2.VERSION)

	w.Write(version)

	return nil
}

func (e Encoder) DecodeRecord(r io.Reader, password string) (record.Record, error) {
	header := make([]byte, 2)
	r.Read(header)

	version := binary.BigEndian.Uint16(header)

	switch version {
	case record_v1.VERSION:
		return e.DecodeV1(r, password, uuid.Nil)
	case record_v2.VERSION:
		return e.DecodeV2(r, password)
	default:
		return nil, errors.Format("unsupported record version \"%d\"", version)
	}
}

func (e Encoder) DecodeV1(r io.Reader, password string, id uuid.UUID) (record_v1.Record, error) {
	raw, err := e.DecodeRawV1(r)
	if err != nil {
		return record_v1.Record{}, err
	}

	record, err := crypt.Unmarshal[record_v1.Record](password, raw.Data)
	if err != nil {
		return record_v1.Record{}, errors.Chain(err, "error decrypting record")
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

func (e Encoder) DecodeV2(r io.Reader, password string) (record_v2.Record, error) {
	data, err := e.DecodeRawV2(r)
	if err != nil {
		return record_v2.Record{}, err
	}

	record, err := crypt.Unmarshal[record_v2.Record](password, data)
	if err != nil {
		return record_v2.Record{}, errors.Chain(err, "error decrypting record")
	}

	return record, nil
}

func (e Encoder) DecodeRawV2(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
