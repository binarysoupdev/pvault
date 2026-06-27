package vault

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"pvault/errors"
	"pvault/vault/index"
	"pvault/vault/index/legacy"
)

const LEGACY_INDEX_FILE = "index.txt"

type Decoder interface {
	Version() int
	Decode(path string) (index.IndexMap, error)
}

func selectDecoder(path string) (Decoder, error) {
	// first try current format
	indexPath := filepath.Join(path, INDEX_FILE)

	_, err := os.Stat(indexPath)
	if err == nil {
		return selectCurrentDecoder(indexPath)
	}

	// else try legacy format
	legacyPath := filepath.Join(path, LEGACY_INDEX_FILE)

	_, err = os.Stat(legacyPath)
	if err == nil {
		return legacy.DecoderV0{}, nil
	}

	return nil, errors.New("index file not found")
}

func selectCurrentDecoder(path string) (Decoder, error) {
	raw := make([]byte, 2)

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Chain(err, "error opening index file")
	}
	defer file.Close()

	_, err = file.Read(raw)
	if err != nil {
		return nil, errors.Chain(err, "error reading version from header")
	}

	version := binary.BigEndian.Uint16(raw[:2])
	switch version {
	case 1:
		return index.Codec{}, nil
	default:
		return nil, errors.Format("unsupported version \"%d\"", version)
	}
}
