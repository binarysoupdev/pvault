package v3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"pvault/vault/index"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/google/uuid"
)

const VERSION = 3

type Encoder struct{}

func (e Encoder) EncodeIndex(w io.Writer, idx index.IndexMap) error {
	err := e.encodeIndexHeader(w, len(idx))
	if err != nil {
		return errors.Chain(err, "error encoding header")
	}

	entryNum := 0
	for name, id := range idx {
		err := e.encodeEntry(w, id, name)
		if err != nil {
			return errors.Chain(err, fmt.Sprintf("error encoding entry [%d]", entryNum))
		}
		entryNum++
	}

	return nil
}

func (e Encoder) encodeIndexHeader(w io.Writer, numRecords int) error {
	header := make([]byte, 2+2)
	binary.BigEndian.PutUint16(header, uint16(index.VERSION))
	binary.BigEndian.PutUint16(header[2:], uint16(numRecords))

	_, err := w.Write(header)
	return err
}

func (e Encoder) encodeEntry(w io.Writer, id uuid.UUID, name string) error {
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(id)+len(name)))

	_, err := w.Write(bytes.Join([][]byte{length, id[:], []byte(name)}, []byte{}))
	return err
}

func (e Encoder) DecodeIndex(r io.Reader) (index.IndexMap, error) {
	header, err := e.decodeIndexHeader(r)
	if err != nil {
		return nil, errors.Chain(err, "error decoding header")
	}

	entryCount := binary.BigEndian.Uint16(header[2:])
	idx := index.IndexMap{}

	for i := range entryCount {
		err := e.decodeEntry(idx, r)
		if err != nil {
			return nil, errors.Chain(err, fmt.Sprintf("error decoding entry [%d]", i))
		}
	}

	return idx, nil
}

func (Encoder) decodeIndexHeader(r io.Reader) ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}

	version := binary.BigEndian.Uint16(header)
	if version != index.VERSION {
		return nil, errors.Format("unsupported index version \"%d\"", version)
	}

	return header, nil
}

func (Encoder) decodeEntry(idx index.IndexMap, r io.Reader) error {
	header := make([]byte, 2)
	if _, err := io.ReadFull(r, header); err != nil {
		return errors.Chain(err, "error reading header")
	}

	length := int(binary.BigEndian.Uint16(header))
	if length < 16 {
		return errors.Format("length too short: %d", length)
	}

	entry := make([]byte, length)
	if _, err := io.ReadFull(r, entry); err != nil {
		return errors.Chain(err, "error decoding body")
	}

	name := ""
	if length > 16 {
		name = string(entry[16:])
	}

	idx[name] = uuid.UUID(entry[:16])
	return nil
}
