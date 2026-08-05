package v1

import (
	"encoding/binary"
	"io"
	"path/filepath"
	"pvault/app/vault/meta"
	"pvault/util"
	"time"

	"github.com/binarysoupdev/go-commando/errors"
)

const HEADER_SIZE = 2 + 2 + 2 + 2

type Encoder struct{}

func (Encoder) MetadataPath(path string) string {
	return filepath.Join(path, "VAULT")
}

func (Encoder) EncodeMetadata(w io.Writer, m meta.Metadata) error {
	header := make([]byte, HEADER_SIZE)
	binary.BigEndian.PutUint16(header, meta.VERSION)
	binary.BigEndian.PutUint16(header[2:], uint16(m.DatabaseVersion))
	binary.BigEndian.PutUint16(header[2+2:], uint16(len(m.Nickname)))

	dateBytes, err := m.CreationDate.MarshalBinary()
	if err != nil {
		return errors.Chain(err, "error marshaling creation date")
	}
	binary.BigEndian.PutUint16(header[2+2+2:], uint16(len(dateBytes)))

	return util.WriteBytes(w, header, []byte(m.Nickname), dateBytes)
}

func (Encoder) DecodeMetadata(r io.Reader) (meta.Metadata, error) {
	header, err := util.ReadBytes(r, HEADER_SIZE)
	if err != nil {
		return meta.Metadata{}, errors.Chain(err, "error reading header")
	}

	version := binary.BigEndian.Uint16(header)
	if version != meta.VERSION {
		return meta.Metadata{}, errors.Format("unsupported metadata version \"%d\"", version)
	}
	dbVersion := binary.BigEndian.Uint16(header[2:])

	name, err := util.ReadBytes(r, int(binary.BigEndian.Uint16(header[2+2:])))
	if err != nil {
		return meta.Metadata{}, errors.Chain(err, "error reading nickname")
	}

	dateBytes, err := util.ReadBytes(r, int(binary.BigEndian.Uint16(header[2+2+2:])))
	if err != nil {
		return meta.Metadata{}, errors.Chain(err, "error reading creation date")
	}

	creationDate := time.Time{}
	err = creationDate.UnmarshalBinary(dateBytes)
	if err != nil {
		return meta.Metadata{}, errors.Chain(err, "error unmarshaling creation date")
	}

	return meta.Metadata{
		DatabaseVersion: int(dbVersion),
		Nickname:        string(name),
		CreationDate:    creationDate,
	}, nil
}
