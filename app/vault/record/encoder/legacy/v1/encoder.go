package v1

import (
	"io"
	"pvault/app/vault/record"
	record_v1 "pvault/app/vault/record/record/legacy/v1"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const VERSION = 1

type Encoder struct{}

func (e Encoder) EncodeRecord(w io.Writer, password string, r record.Record) error {
	if r.GetVersion() != record_v1.VERSION {
		return errors.Format("unsupported record version \"%d\"", r.GetVersion())
	}

	return e.EncodeV1(w, password, r)
}

func (e Encoder) DecodeRecord(r io.Reader, password string) (record.Record, error) {
	return e.DecodeV1(r, password, uuid.Nil, "")
}
