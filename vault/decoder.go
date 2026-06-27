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

func (v Vault) selectDecoder(path string) (Decoder, error) {
	// first try modern format
	_, err := os.Stat(v.IndexPath())
	if err == nil {
		return v.selectDecoderFromHeader()
	}

	// else try legacy format
	_, err = os.Stat(filepath.Join(path, LEGACY_INDEX_FILE))
	if err == nil {
		return legacy.DecoderV0{}, nil
	}

	return nil, errors.New("index file not found")
}

func (v Vault) selectDecoderFromHeader() (Decoder, error) {
	header := make([]byte, 2)

	file, err := os.Open(v.IndexPath())
	if err != nil {
		return nil, errors.Chain(err, "error opening index file")
	}
	defer file.Close()

	_, err = file.Read(header)
	if err != nil {
		return nil, errors.Chain(err, "error reading version header")
	}
	version := binary.BigEndian.Uint16(header)

	switch version {
	case 1:
		return index.Codec{}, nil
	default:
		return nil, errors.Format("unsupported version \"%d\"", version)
	}
}
