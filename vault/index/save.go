package index

import (
	"encoding/binary"
	"os"
	"pvault/chain"

	"github.com/google/uuid"
)

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
	entry := make([]byte, 2+16+len(name))

	binary.BigEndian.PutUint16(entry, 16+uint16(len(name)))
	copy(entry[2:], id[:])
	copy(entry[2+16:], []byte(name))

	_, err := file.Write(entry)
	return err
}
