package v3

import (
	"encoding/binary"
	"io"

	"pvault/vault/record"
	v2 "pvault/vault/record/encoder/legacy/v2"
	record_v1 "pvault/vault/record/record/legacy/v1"
	record_v2 "pvault/vault/record/record/v2"

	"github.com/binarysoupdev/go-extensions/errors"
)

const VERSION = 3

type Encoder struct{}

func (e Encoder) EncodeRecord(w io.Writer, password string, r record.Record) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header, uint16(r.GetVersion()))

	if _, err := w.Write(header); err != nil {
		return errors.Chain(err, "error writing version header")
	}

	switch r.GetVersion() {
	case record_v1.VERSION:
		return e.EncodeV1(w, password, r)
	case record_v2.VERSION:
		return v2.Encoder{}.EncodeV2(w, password, r)
	default:
		return errors.Format("unsupported record version \"%d\"", r.GetVersion())
	}
}

func (e Encoder) DecodeRecord(r io.Reader, password string) (record.Record, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, errors.Chain(err, "error reading version header")
	}

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
