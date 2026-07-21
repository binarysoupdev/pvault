package v2

import (
	"encoding/binary"
	"io"
	"pvault/app/vault/index"

	"github.com/binarysoupdev/go-commando/errors"
	"github.com/google/uuid"
)

const VERSION = 2

type Encoder struct{}

func (e Encoder) EncodeIndex(w io.Writer, idx index.IndexMap) error {
	err := e.writeHeader(w, len(idx))
	if err != nil {
		return errors.Chain(err, "error writing header")
	}

	for name, id := range idx {
		err = e.writeEntry(w, id, name)
		if err != nil {
			return errors.Chain(err, "error writing entry")
		}
	}

	return nil
}

func (e Encoder) writeHeader(w io.Writer, numRecords int) error {
	header := make([]byte, 2)
	binary.BigEndian.PutUint16(header[2:], uint16(numRecords))

	_, err := w.Write(header)
	return err
}

func (e Encoder) writeEntry(w io.Writer, id uuid.UUID, name string) error {
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(id)+len(name)))

	w.Write(length)
	w.Write(id[:])
	w.Write([]byte(name))

	return nil
}

func (e Encoder) DecodeIndex(r io.Reader) (index.IndexMap, error) {
	header := make([]byte, 2)
	r.Read(header)

	entryCount := binary.BigEndian.Uint16(header[2:])
	idx := index.IndexMap{}

	for range entryCount {
		length := make([]byte, 2)
		e.decodeEntry(idx, r, int(binary.BigEndian.Uint16(length)))
	}

	return idx, nil
}

func (Encoder) decodeEntry(idx index.IndexMap, r io.Reader, length int) {
	id := uuid.UUID{}
	r.Read(id[:])

	name := make([]byte, length-len(id))
	r.Read(name)

	idx[string(name)] = id
}
