package vault

import (
	"encoding/binary"
	"os"
	"pvault/chain"

	"github.com/google/uuid"
)

const INDEX_VERSION = 1

type IndexMap map[string]uuid.UUID

func (idx IndexMap) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return chain.Error(err, "error creating file")
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
	entry := make([]byte, 16+4+len(name))
	copy(entry, id[:])
	binary.BigEndian.PutUint32(entry[16:], uint32(len(name)))
	copy(entry[16+4:], []byte(name))

	_, err := file.Write(entry)
	return err
}
