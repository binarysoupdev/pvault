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
	header := make([]byte, 2+2)
	binary.BigEndian.PutUint16(header, uint16(index.VERSION))
	binary.BigEndian.PutUint16(header[2:], uint16(numRecords))

	_, err := w.Write(header)
	return err
}

func (e Encoder) writeEntry(w io.Writer, id uuid.UUID, name string) error {
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(id)+len(name)))

	// TODO: handle error
	w.Write(length)
	w.Write(id[:])
	w.Write([]byte(name))

	return nil
}

func (e Encoder) DecodeIndex(r io.Reader) (index.IndexMap, error) {
	header := make([]byte, 4)
	r.Read(header)

	version := binary.BigEndian.Uint16(header)
	if version != index.VERSION {
		return nil, errors.Format("unsupported index version \"%d\"", version)
	}

	entryCount := binary.BigEndian.Uint16(header[2:])
	idx := index.IndexMap{}

	for range entryCount {
		err := e.decodeEntry(idx, r)
		if err != nil {
			return nil, errors.Chain(err, "error decoding entry")
		}
	}

	return idx, nil
}

func (Encoder) decodeEntry(idx index.IndexMap, r io.Reader) error {
	length := make([]byte, 2)
	if _, err := r.Read(length); err != nil {
		return errors.Chain(err, "error reading length")
	}

	id := uuid.UUID{}
	if _, err := r.Read(id[:]); err != nil {
		return errors.Chain(err, "error reading id")
	}

	name := make([]byte, int(binary.BigEndian.Uint16(length))-len(id))
	if _, err := r.Read(name); err != nil {
		return errors.Chain(err, "error reading name")
	}

	idx[string(name)] = id
	return nil
}
