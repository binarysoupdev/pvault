package encoder

import (
	"io"
	"os"
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

type Encoder interface {
	EncodeIndex(w io.Writer, idx index.IndexMap) error
	DecodeIndex(r io.Reader) (index.IndexMap, error)
}

func SaveIndexFile(e Encoder, idx index.IndexMap, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Chain(err, "error creating index file")
	}
	defer file.Close()

	err = e.EncodeIndex(file, idx)
	if err != nil {
		return errors.Chain(err, "error encoding index")
	}

	return nil
}

func LoadIndexFile(e Encoder, path string) (index.IndexMap, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Chain(err, "error opening index file")
	}
	defer file.Close()

	idx, err := e.DecodeIndex(file)
	if err != nil {
		return nil, errors.Chain(err, "error decoding index")
	}

	return idx, nil
}
