package encoder

import (
	"io"
	"pvault/app/vault/index"
)

type Encoder interface {
	EncodeIndex(w io.Writer, idx index.IndexMap) error
	DecodeIndex(r io.Reader) (index.IndexMap, error)
}
