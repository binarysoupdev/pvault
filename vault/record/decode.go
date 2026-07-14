package record

import (
	"encoding/binary"
	"io"
	v1 "pvault/vault/record/version1"
	v2 "pvault/vault/record/version2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func Decode(r io.Reader, password string, id uuid.UUID) (Record, error) {
	header := make([]byte, 2)
	r.Read(header)

	version := binary.BigEndian.Uint16(header)

	switch version {
	case 1:
		return v1.Decode(r, password, id)
	case 2:
		return v2.Decode(r, password)
	default:
		return nil, errors.Format("unsupported record version \"%d\"", version)
	}
}
