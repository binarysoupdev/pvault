package v3

import (
	"encoding/binary"
	"io"
	"pvault/util"
	"pvault/vault/record"
	v2 "pvault/vault/record/encoder/legacy/v2"
	record_v1 "pvault/vault/record/record/legacy/v1"
	record_v2 "pvault/vault/record/record/v2"

	"github.com/binarysoupdev/go-commando/errors"
)

const VERSION = 3

type Encoder struct{}

func (e Encoder) EncodeRecord(w io.Writer, password string, r record.Record) error {
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
	header, err := util.ReadBytes(r, 2)
	if err != nil {
		return nil, errors.Chain(err, "error reading header")
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
