package record

import (
	"encoding/binary"
	"os"
	v1 "pvault/vault/record/version1"
	v2 "pvault/vault/record/version2"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func Load(path string, password string, id uuid.UUID) (Record, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Chain(err, "error reading record file")
	}

	version := binary.BigEndian.Uint16(bytes)

	switch version {
	case 1:
		return v1.Unmarshal(bytes, password, id)
	case 2:
		return v2.Unmarshal(bytes, password)
	default:
		return nil, errors.Format("unsupported record version \"%d\"", version)
	}
}
