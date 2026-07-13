package record

import (
	"encoding/binary"
	v1 "pvault/vault/record/version1"
	v2 "pvault/vault/record/version2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func UnmarshalGeneric(password string, bytes []byte, id uuid.UUID) (Record, error) {
	version := binary.BigEndian.Uint16(bytes)

	switch version {
	case 1:
		return v1.Unmarshal(password, bytes, id)
	case 2:
		return v2.Unmarshal(password, bytes)
	default:
		return nil, errors.Format("unsupported record version \"%d\"", version)
	}
}
