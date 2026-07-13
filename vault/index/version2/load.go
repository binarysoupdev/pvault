package v2

import (
	"encoding/binary"
	"os"
	"pvault/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (idx Index) Load() (index.IndexMap, error) {
	raw, err := os.ReadFile(idx.filepath())
	if err != nil {
		return index.IndexMap{}, errors.Chain(err, "error reading index file")
	}
	header := raw[:4]

	version := binary.BigEndian.Uint16(header)
	if version != idx.GetVersion() {
		return index.IndexMap{}, errors.Format("incorrect version \"%d\"", version)
	}

	entryCount := binary.BigEndian.Uint16(header[2:])
	ptr := len(header)

	m := index.IndexMap{}

	for range entryCount {
		length := int(binary.BigEndian.Uint16(raw[ptr : ptr+2]))
		ptr += 2

		idx.decodeEntry(m, raw[ptr:ptr+length])
		ptr += length
	}

	return m, nil
}

func (Index) decodeEntry(idx index.IndexMap, raw []byte) {
	id, _ := uuid.FromBytes(raw[:16])
	name := string(raw[16:])

	idx[name] = id
}
