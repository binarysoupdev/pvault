package v3

import (
	"encoding/binary"
	"io"
	"pvault/app/vault/record"
	v2 "pvault/app/vault/record/encoder/v2"
	record_v1 "pvault/app/vault/record/record/v1"
	record_v2 "pvault/app/vault/record/record/v2"
	"pvault/crypt"

	"github.com/binarysoupdev/go-commando/errors"
)

const VERSION = 3

type Encoder struct{}

func (e Encoder) EncodeRecord(w io.Writer, password string, r record.Record) error {
	ciphertext, err := crypt.Marshal(password, r)
	if err != nil {
		return errors.Chain(err, "error encrypting record")
	}

	switch r.GetVersion() {
	case record_v1.VERSION:
		return e.EncodeRawV1(w, ciphertext, r.GetID(), r.GetName())
	case record_v2.VERSION:
		return v2.Encoder{}.EncodeRawV2(w, ciphertext)
	default:
		return errors.Format("unsupported record version \"%d\"", r.GetVersion())
	}
}

func (e Encoder) DecodeRecord(r io.Reader, password string) (record.Record, error) {
	header := make([]byte, 2)
	r.Read(header)

	version := binary.BigEndian.Uint16(header)

	switch version {
	case record_v1.VERSION:
		return e.DecodeV1(r, password)
	case record_v2.VERSION:
		return v2.Encoder{}.DecodeV2(r, password)
	default:
		return nil, errors.Format("unsupported record version \"%d\"", version)
	}
}
