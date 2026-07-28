package v1

import (
	"io"
	"pvault/app/vault/record"
	record_v1 "pvault/app/vault/record/record/v1"

	"github.com/binarysoupdev/cryptool/crypt"
	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const VERSION = 1

type Encoder struct{}

func (e Encoder) EncodeRecord(w io.Writer, c crypt.Crypt, r record.Record) error {
	if r.GetVersion() != record_v1.VERSION {
		return errors.Format("unsupported record version \"%d\"", r.GetVersion())
	}

	ciphertext, err := crypt.Marshal(c, r)
	if err != nil {
		return errors.Chain(err, "error encrypting record")
	}

	err = e.EncodeRawV1(w, ciphertext)
	if err != nil {
		return errors.Chain(err, "error encoding record v1")
	}

	return nil
}

func (e Encoder) DecodeRecord(r io.Reader, password string) (record.Record, error) {
	return e.DecodeV1(r, password, uuid.Nil, "")
}
