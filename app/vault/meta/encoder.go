package meta

import (
	"encoding/binary"
	"io"
	"os"
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
)

func SaveMetadata(m Metadata) error {
	file, err := os.Create(m.Path)
	if err != nil {
		return errors.Chain(err, "error creating metadata file")
	}
	defer file.Close()

	err = EncodeMetadata(file, m)
	if err != nil {
		return errors.Chain(err, "error encoding metadata")
	}

	return nil
}

func EncodeMetadata(w io.Writer, m Metadata) error {
	header := make([]byte, 2+2)
	binary.BigEndian.PutUint16(header, VERSION)
	binary.BigEndian.PutUint16(header, uint16(m.DatabaseVersion))

	w.Write(header)
	return nil
}

func LoadMetadata(path string) (Metadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return Metadata{}, errors.Chain(err, "error opening metadata file")
	}
	defer file.Close()

	m, err := DecodeMetadata(file)
	if err != nil {
		return Metadata{}, errors.Chain(err, "error decoding metadata")
	}

	m.Path = path
	return m, nil
}

func DecodeMetadata(r io.Reader) (Metadata, error) {
	header := make([]byte, 2+2)
	r.Read(header)

	version := binary.BigEndian.Uint16(header)
	if version != index.VERSION {
		return Metadata{}, errors.Format("unsupported metadata version \"%d\"", version)
	}

	dbVersion := binary.BigEndian.Uint16(header[2:])

	return Metadata{
		DatabaseVersion: int(dbVersion),
	}, nil
}
