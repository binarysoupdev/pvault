package v2

import (
	"io"
	"pvault/util"
	"pvault/vault/index"
	v3 "pvault/vault/index/encoder/v3"

	"github.com/binarysoupdev/go-commando/errors"
)

const VERSION = 2

type Encoder struct{}

func (e Encoder) EncodeIndex(w io.Writer, idx index.IndexMap) error {
	err := util.WriteBytes(w, make([]byte, 2))
	if err != nil {
		return errors.Chain(err, "error encoding null header")
	}

	return v3.Encoder{}.EncodeIndex(w, idx)
}

func (e Encoder) DecodeIndex(r io.Reader) (index.IndexMap, error) {
	_, err := util.ReadBytes(r, 2)
	if err != nil {
		return nil, errors.Chain(err, "error decoding null header")
	}

	return v3.Encoder{}.DecodeIndex(r)
}
