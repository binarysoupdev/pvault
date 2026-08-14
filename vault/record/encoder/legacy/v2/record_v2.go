package v2

import (
	"bytes"
	"encoding/binary"
	"io"

	"pvault/vault/record"
	v2 "pvault/vault/record/record/v2"

	"github.com/binarysoupdev/go-extensions/errors"
)

func (e Encoder) EncodeV2(w io.Writer, password string, r record.Record) error {
	bytes, err := record.Encrypt(r, password)
	if err != nil {
		return errors.Chain(err, "error encrypting record v2")
	}

	err = e.EncodeRawV2(w, bytes)
	if err != nil {
		return errors.Chain(err, "error encoding record v2")
	}

	return nil
}

func (e Encoder) EncodeRawV2(w io.Writer, data []byte) error {
	version := make([]byte, 2)
	binary.BigEndian.PutUint16(version, v2.VERSION)

	_, err := w.Write(bytes.Join([][]byte{version, data}, []byte{}))
	return err
}

func (e Encoder) DecodeV2(r io.Reader, password string) (v2.Record, error) {
	data, err := e.DecodeRawV2(r)
	if err != nil {
		return v2.Record{}, errors.Chain(err, "error decoding record v2")
	}

	record, err := record.Decrypt[v2.Record](data, password)
	if err != nil {
		return v2.Record{}, errors.Chain(err, "error decrypting record v2")
	}

	return record, nil
}

func (e Encoder) DecodeRawV2(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
