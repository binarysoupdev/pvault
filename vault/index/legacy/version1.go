package legacy

import (
	"pvault/vault/index"
)

type DecoderV1 struct{}

func (DecoderV1) Decode(path string) (index.IndexMap, error) {
	return index.IndexMap{}, nil
}
