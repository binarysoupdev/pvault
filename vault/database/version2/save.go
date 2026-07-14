package v2

import (
	"encoding/binary"
	"os"
	"pvault/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

func (idx Index) SaveIndex(m index.IndexMap) error {
	file, err := os.Create(idx.Filepath())
	if err != nil {
		return errors.Chain(err, "error creating index file")
	}
	defer file.Close()

	err = idx.writeHeader(file, len(m))
	if err != nil {
		return errors.Chain(err, "error writing header")
	}

	for name, id := range m {
		err = idx.writeEntry(file, name, id)
		if err != nil {
			return errors.Chain(err, "error writing entry")
		}
	}

	return nil
}

func (idx Index) writeHeader(file *os.File, numRecords int) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint16(header, uint16(idx.GetVersion()))
	binary.BigEndian.PutUint16(header[2:], uint16(numRecords))

	_, err := file.Write(header)
	return err
}

func (Index) writeEntry(file *os.File, name string, id uuid.UUID) error {
	entry := make([]byte, 2+16+len(name))

	binary.BigEndian.PutUint16(entry, 16+uint16(len(name)))
	copy(entry[2:], id[:])
	copy(entry[2+16:], []byte(name))

	_, err := file.Write(entry)
	return err
}
