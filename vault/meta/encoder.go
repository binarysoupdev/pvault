package meta

import (
	"io"
	"os"

	"github.com/binarysoupdev/go-extensions/errors"
)

type Encoder interface {
	MetadataPath(path string) string
	EncodeMetadata(w io.Writer, m Metadata) error
	DecodeMetadata(r io.Reader) (Metadata, error)
}

func SaveMetadata(e Encoder, path string, m Metadata) error {
	file, err := os.Create(e.MetadataPath(path))
	if err != nil {
		return errors.Chain(err, "error creating metadata file")
	}
	defer file.Close()

	err = e.EncodeMetadata(file, m)
	if err != nil {
		return errors.Chain(err, "error encoding metadata")
	}

	return nil
}

func LoadMetadata(e Encoder, path string) (Metadata, error) {
	file, err := os.Open(e.MetadataPath(path))
	if err != nil {
		return Metadata{}, errors.Chain(err, "error opening metadata file")
	}
	defer file.Close()

	m, err := e.DecodeMetadata(file)
	if err != nil {
		return Metadata{}, errors.Chain(err, "error decoding metadata")
	}

	return m, nil
}
