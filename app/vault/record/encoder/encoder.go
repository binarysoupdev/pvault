package encoder

import (
	"io"
	"pvault/app/vault/record"
)

type Encoder interface {
	EncodeRecord(w io.Writer, password string, r record.Record) error
	DecodeRecord(r io.Reader, password string) (record.Record, error)
}
