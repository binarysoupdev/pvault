package index

import (
	"io"
)

type Encoder interface {
	EncodeIndex(w io.Writer, idx IndexMap) error
	DecodeIndex(r io.Reader) (IndexMap, error)
}
