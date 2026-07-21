package encoder

import (
	"io"
	"os"
	"pvault/app/vault/record"

	"github.com/binarysoupdev/go-commando/errors"
)

type Encoder interface {
	EncodeRecord(w io.Writer, password string, r record.Record) error
	DecodeRecord(r io.Reader, password string) (record.Record, error)
}

func SaveRecordFile(e Encoder, r record.Record, path string, password string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating record file")
	}
	defer file.Close()

	err = e.EncodeRecord(file, password, r)
	if err != nil {
		return errors.Chain(err, "error encoding record")
	}

	return nil
}

func LoadRecordFile(e Encoder, path string, password string) (record.Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Chain(err, "error opening record file")
	}
	defer file.Close()

	r, err := e.DecodeRecord(file, password)
	if err != nil {
		return nil, errors.Chain(err, "error decoding record")
	}

	return r, nil
}
