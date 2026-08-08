package record

import (
	"io"
)

type Encoder interface {
	EncodeRecord(w io.Writer, password string, r Record) error
	DecodeRecord(r io.Reader, password string) (Record, error)
}
