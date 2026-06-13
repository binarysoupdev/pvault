package vault

import (
	"encoding/binary"
	"fmt"
	"os"
	"pvault/chain"

	"github.com/google/uuid"
)

const INDEX_VERSION = 1

type IndexMap map[string]uuid.UUID

func (idx IndexMap) Load(path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return chain.Error(err, "error reading index file")
	}

	header := raw[:4]

	version := binary.BigEndian.Uint16(header)
	if version < INDEX_VERSION {
		return fmt.Errorf("unsupported version \"%d\"", version)
	}

	entryCount := binary.BigEndian.Uint16(header[2:])
	ptr := uint32(len(header))

	for range entryCount {
		length := binary.BigEndian.Uint32(raw[ptr : ptr+4])
		ptr += 4

		idx.loadEntry(raw[ptr : ptr+length])
		ptr += length
	}

	return nil
}

func (idx IndexMap) loadEntry(raw []byte) {
	id, _ := uuid.FromBytes(raw[:16])
	name := string(raw[16:])

	idx[name] = id
}

func (idx IndexMap) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return chain.Error(err, "error creating index file")
	}
	defer file.Close()

	err = idx.writeHeader(file)
	if err != nil {
		return chain.Error(err, "error writing header")
	}

	for name, id := range idx {
		err = idx.writeEntry(file, name, id)
		if err != nil {
			return chain.Error(err, "error writing entry")
		}
	}

	return nil
}

func (idx IndexMap) writeHeader(file *os.File) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint16(header, uint16(INDEX_VERSION))
	binary.BigEndian.PutUint16(header[2:], uint16(len(idx)))

	_, err := file.Write(header)
	return err
}

func (IndexMap) writeEntry(file *os.File, name string, id uuid.UUID) error {
	entry := make([]byte, 4+16+len(name))

	binary.BigEndian.PutUint32(entry, 16+uint32(len(name)))
	copy(entry[4:], id[:])
	copy(entry[4+16:], []byte(name))

	_, err := file.Write(entry)
	return err
}
