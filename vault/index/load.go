package index

import (
	"encoding/binary"
	"fmt"
	"os"
	"pvault/chain"

	"github.com/google/uuid"
)

func LoadIndex(path string) (IndexMap, error) {
	idx := IndexMap{}

	raw, err := os.ReadFile(path)
	if err != nil {
		return idx, chain.Error(err, "error reading index file")
	}

	header := raw[:4]

	version := binary.BigEndian.Uint16(header)
	if version > INDEX_VERSION {
		return idx, fmt.Errorf("unsupported version \"%d\"", version)
	}

	entryCount := binary.BigEndian.Uint16(header[2:])
	ptr := len(header)

	for range entryCount {
		length := int(binary.BigEndian.Uint16(raw[ptr : ptr+2]))
		ptr += 2

		idx.loadEntry(raw[ptr : ptr+length])
		ptr += length
	}

	return idx, nil
}

func (idx IndexMap) loadEntry(raw []byte) {
	id, _ := uuid.FromBytes(raw[:16])
	name := string(raw[16:])

	idx[name] = id
}
