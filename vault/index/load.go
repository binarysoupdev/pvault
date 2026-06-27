package index

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"pvault/errors"

	"github.com/google/uuid"
)

func LoadIndex(path string) (IndexMap, error) {
	raw, version, err := loadRaw(path)
	if err != nil {
		return IndexMap{}, err
	}

	if version > VERSION || version < 0 {
		return IndexMap{}, errors.Format("unsupported version \"%d\"", version)
	} else if version < VERSION {
		return IndexMap{}, newOutOfDateError(version)
	}

	idx := IndexMap{}
	header := raw[:4]

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

func loadRaw(path string) ([]byte, int, error) {
	// first try current format
	indexPath := filepath.Join(path, INDEX_FILE)

	_, err := os.Stat(indexPath)
	if err == nil {
		raw, err := os.ReadFile(indexPath)
		if err != nil {
			return nil, -1, errors.Chain(err, "error reading index file")
		}

		version := binary.BigEndian.Uint16(raw[:2])
		return raw, int(version), nil
	}

	// else try legacy format
	legacyPath := filepath.Join(path, LEGACY_INDEX_FILE)

	_, err = os.Stat(legacyPath)
	if err == nil {
		raw, err := os.ReadFile(legacyPath)
		if err != nil {
			return nil, -1, errors.Chain(err, "error reading legacy index file")
		}

		return raw, 0, nil
	}

	return nil, -1, errors.New("index file not found")
}

func (idx IndexMap) loadEntry(raw []byte) {
	id, _ := uuid.FromBytes(raw[:16])
	name := string(raw[16:])

	idx[name] = id
}
