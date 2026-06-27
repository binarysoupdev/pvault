package legacy

import (
	"pvault/vault/index"
)

type DecoderV0 struct{}

func (DecoderV0) Version() int {
	return 0
}

func (DecoderV0) Decode(path string) (index.IndexMap, error) {
	return index.IndexMap{}, nil
}
