package index

import (
	"encoding/binary"
	"os"
	"pvault/errors"

	"github.com/google/uuid"
)

func (c Codec) Decode(path string) (IndexMap, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return IndexMap{}, errors.Chain(err, "error reading index file")
	}
	header := raw[:4]

	version := binary.BigEndian.Uint16(header)
	if version != 1 {
		return IndexMap{}, errors.Format("incorrect version \"%d\"", version)
	}

	entryCount := binary.BigEndian.Uint16(header[2:])
	ptr := len(header)

	idx := IndexMap{}

	for range entryCount {
		length := int(binary.BigEndian.Uint16(raw[ptr : ptr+2]))
		ptr += 2

		c.decodeEntry(idx, raw[ptr:ptr+length])
		ptr += length
	}

	return idx, nil
}

func (Codec) decodeEntry(idx IndexMap, raw []byte) {
	id, _ := uuid.FromBytes(raw[:16])
	name := string(raw[16:])

	idx[name] = id
}
