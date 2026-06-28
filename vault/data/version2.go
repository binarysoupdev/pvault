package data

import (
	"encoding/binary"
	"os"
	"pvault/errors"
	"pvault/vault/index"

	"github.com/google/uuid"
)

type DatabaseV2 struct {
	Path string
}

func (DatabaseV2) Version() uint16 {
	return 2
}

func NewDatabaseV2(path string) DatabaseV2 {
	return DatabaseV2{
		Path: path,
	}
}

func (DatabaseV2) Upgrade() error {
	return nil
}

func (db DatabaseV2) SaveIndex(idx index.IndexMap) error {
	file, err := os.Create(db.Path)
	if err != nil {
		return errors.Chain(err, "error creating index file")
	}
	defer file.Close()

	err = db.writeHeader(file, len(idx))
	if err != nil {
		return errors.Chain(err, "error writing header")
	}

	for name, id := range idx {
		err = db.writeEntry(file, name, id)
		if err != nil {
			return errors.Chain(err, "error writing entry")
		}
	}

	return nil
}

func (db DatabaseV2) writeHeader(file *os.File, numRecords int) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint16(header, db.Version())
	binary.BigEndian.PutUint16(header[2:], uint16(numRecords))

	_, err := file.Write(header)
	return err
}

func (DatabaseV2) writeEntry(file *os.File, name string, id uuid.UUID) error {
	entry := make([]byte, 2+16+len(name))

	binary.BigEndian.PutUint16(entry, 16+uint16(len(name)))
	copy(entry[2:], id[:])
	copy(entry[2+16:], []byte(name))

	_, err := file.Write(entry)
	return err
}

func (db DatabaseV2) LoadIndex() (index.IndexMap, error) {
	raw, err := os.ReadFile(db.Path)
	if err != nil {
		return index.IndexMap{}, errors.Chain(err, "error reading index file")
	}
	header := raw[:4]

	version := binary.BigEndian.Uint16(header)
	if version != db.Version() {
		return index.IndexMap{}, errors.Format("incorrect version \"%d\"", version)
	}

	entryCount := binary.BigEndian.Uint16(header[2:])
	ptr := len(header)

	idx := index.IndexMap{}

	for range entryCount {
		length := int(binary.BigEndian.Uint16(raw[ptr : ptr+2]))
		ptr += 2

		db.decodeEntry(idx, raw[ptr:ptr+length])
		ptr += length
	}

	return idx, nil
}

func (DatabaseV2) decodeEntry(idx index.IndexMap, raw []byte) {
	id, _ := uuid.FromBytes(raw[:16])
	name := string(raw[16:])

	idx[name] = id
}
