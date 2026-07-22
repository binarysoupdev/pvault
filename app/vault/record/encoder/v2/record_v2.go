package v2

import (
	"encoding/binary"
	"io"
	v2 "pvault/app/vault/record/record/v2"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
)

func (e Encoder) EncodeRawV2(w io.Writer, data []byte) error {
	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, v2.VERSION)

	w.Write(version)

	return nil
}

func (e Encoder) DecodeV2(r io.Reader, password string) (v2.Record, error) {
	data, err := e.DecodeRawV2(r)
	if err != nil {
		return v2.Record{}, err
	}

	record, err := crypt.Unmarshal[v2.Record](password, data)
	if err != nil {
		return v2.Record{}, errors.Chain(err, "error decrypting record")
	}

	return record, nil
}

func (e Encoder) DecodeRawV2(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
